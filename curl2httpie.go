package curl2httpie

import (
	"bytes"
	"errors"
	"flag"
)

var cl *flag.FlagSet

type DataFlag [][]byte

// http://stackoverflow.com/questions/36854408/how-to-append-to-a-slice-pointer-receiver
func (f *DataFlag) String() string {
	if len(*f) > 0 {
		return "'" + string(bytes.Join(*f, []byte("' '"))) + "'"
	}
	return ""
}

func (f *DataFlag) Set(value string) error {
	*f = append(*f, []byte(value))
	return nil
}

type HeaderFlag map[string][]string

func (h HeaderFlag) String() string {
	bs := make([]byte, 1024)
	first := true
	for k, v := range h {
		if first {
			first = false
		} else {
			bs = append(bs, ' ')
		}

		bs = append(bs, []byte("\""+k+":")...)
		for i, vv := range v {
			if i > 0 {
				bs = append(bs, []byte(", ")...)
			}
			bs = append(bs, []byte(vv)...)
		}
		bs = append(bs, '"')
	}

	return string(bs)
}
func (h HeaderFlag) Set(value string) error {
	kv := bytes.SplitN([]byte(value), []byte(":"), 2)
	bk := bytes.TrimSpace(kv[0])
	bv := bytes.TrimSpace(kv[1])

	k, v := string(bk), string(bv)

	if _, ok := h[k]; !ok {
		h[k] = []string{v}
	} else {
		h[k] = append(h[k], v)
	}

	return nil
}

var osargs []string

var method string
var data *DataFlag
var user string
var url string
var header HeaderFlag
var path string
var urlParameters map[string][]string
var verbose bool
var userAgent string

func emptyVariables() {
	osargs = make([]string, 0, 64)
	method = ""
	tData := make(DataFlag, 0, 0)
	data = &tData
	user = ""
	url = ""
	header = make(HeaderFlag)
	path = ""
	urlParameters = make(map[string][]string)
	verbose = false
	userAgent = ""
}

func Httpie(args []string) (string, error) {
	emptyVariables()
	osargs = args
	l := len(osargs)
	curlCommandExists := l >= 2 && osargs[1] == "curl"
	if !curlCommandExists {
		return "", errors.New("no curl command")
	}
	if l < 3 {
		return "", errors.New("not enough arguments")
	}

	idx := 0
	if osargs[2][0] == '-' {
		idx = 2
	} else {
		idx = 3
	}
	err := parseFlag(idx)
	if idx == 2 && url == "" {
		url = osargs[l-1]
	} else if idx == 3 && url == "" {
		url = osargs[2]
	}
	if err != nil {
		return "", err
	}

	if url == "" {
		return "", errors.New("could not find url")
	}
	{
		t := bytes.SplitN([]byte(url), []byte{'?'}, 2)
		path = string(t[0])
		if len(t) > 1 {
			params := bytes.Split(t[1], []byte{'&'})
			for _, param := range params {
				byteKV := bytes.SplitN(param, []byte{'='}, 2)
				k, v := string(byteKV[0]), string(byteKV[1])

				if _, ok := urlParameters[k]; !ok {
					urlParameters[k] = []string{v}
				} else {
					urlParameters[k] = append(urlParameters[k], v)
				}
			}
		}
	}

	return string(bytes.TrimSpace([]byte(Convert2Httpie()))), nil
}

func Convert2Httpie() string {
	user := user
	if user != "" {
		user = "-a " + user
	}

	bs := make([]byte, 1024)
	first := true
	for k, v := range urlParameters {
		for _, vv := range v {
			if first {
				first = false
			} else {
				bs = append(bs, ' ')
			}

			bs = append(bs, []byte("'"+k+"=="+vv+"'")...)
		}
	}
	params := bs

	v := ""
	if verbose {
		v = "-v"
	}

	method = string(bytes.ToUpper([]byte(method)))
	if method == "GET" && data.String() != "" {
		method = "POST"
	}

	contentType := ""
	switch method {
	case "POST", "PUT", "DELETE":
		contentType = "--form"
		if contentT, ok := header["Content-Type"]; ok {
			for _, v := range contentT {
				if bytes.Contains([]byte(v), []byte("json")) {
					contentType = "--json"
				}
			}
		}
	}

	if userAgent != "" {
		header["User-Agent"] = []string{userAgent}
	}
	return "http" + appendSpace(contentType, user, method, path, header.String(), string(params), v, data.String())
}

func appendSpace(vs ...string) (val string) {
	for _, v := range vs {
		if v != "" {
			val += " " + v
		}
	}
	return
}

func parseFlag(idx int) error {
	cl = flag.NewFlagSet("curl", flag.ExitOnError)

	stringVarAliasNames(&method, []string{"X", "request"}, "GET", "Specify request command to use")
	stringVarAliasNames(&user, []string{"u", "user"}, "", "USER[:PASSWORD]  Server user and password")
	stringVarAliasNames(&url, []string{"url"}, "", "URL to work with")
	stringVarAliasNames(&userAgent, []string{"A", "user-agent"}, "", "Send User-Agent STRING to server")

	cl.Var(data, "data", "HTTP POST data")
	cl.Var(data, "d", "HTTP POST data")

	cl.Var(&header, "H", "Pass custom header LINE to server")
	cl.Var(&header, "header", "Pass custom header LINE to server")
	cl.BoolVar(&verbose, "v", false, "Print the whole request as well as the response.")
	cl.BoolVar(&verbose, "verbose", false, "Print the whole request as well as the response.")
	return cl.Parse(osargs[idx:])
}

func stringVarAliasNames(p *string, names []string, value string, usage string) {
	for _, name := range names {
		cl.StringVar(p, name, value, usage)
	}
}

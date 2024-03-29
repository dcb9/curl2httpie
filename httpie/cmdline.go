package httpie

import (
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"strings"

	"github.com/dcb9/curl2httpie/shellwords"
)

type CmdLine struct {
	Flags         []*Flag
	Method        *Method
	URL           string
	IsHttps       bool
	Items         []*Item
	HasBody       bool
	ContentType   string
	DirectedInput io.ReadCloser
}

func (cl *CmdLine) AddFlag(f *Flag) {
	cl.Flags = append(cl.Flags, f)
}

func (cl *CmdLine) SetMethod(m *Method) {
	cl.Method = m
}

func (cl *CmdLine) SetURL(url string) {
	cl.URL = url
}

func (cl *CmdLine) AddItem(i *Item) {
	cl.Items = append(cl.Items, i)
}

func (cl *CmdLine) String() string {
	// slice
	s := make([]string, 0, len(cl.Flags)+len(cl.Items)+3) // http method url

	if strings.HasPrefix(cl.URL, "https://") {
		s = append(s, "https")
		cl.URL = strings.TrimPrefix(cl.URL, "https://")
	} else {
		s = append(s, "http")
		cl.URL = strings.TrimPrefix(cl.URL, "http://")
	}

	// flags

	// default flag
	for _, v := range cl.Flags {
		s = append(s, v.String())
	}

	if cl.Method == nil {
		if cl.HasBody {
			cl.Method = NewMethod("POST") // use POST as default, if it has body and no specified method
		} else {
			cl.Method = NewMethod("GET")
		}
	}

	if cl.ContentType == "" {
		if cl.Method.String() != "GET" {
			s = append(s, "--form") // use form as default
		}
	} else {
		s = append(s, "--"+cl.ContentType)
	}

	s = append(s, cl.Method.String(), shellwords.AddQuoteIfNeeded(cl.URL))

	for _, v := range cl.Items {
		s = append(s, v.String())
	}

	if cl.DirectedInput != nil && cl.HasBody {
		bytes, err := ioutil.ReadAll(cl.DirectedInput)
		if err != nil {
			fmt.Println("Skipped: Read DirectedInput error", err.Error())
		} else {
			s = append([]string{"echo", fmt.Sprintf("'%s'", bytes), "|"}, s...)
			cl.DirectedInput.Close()
		}
	}

	return strings.Join(s, " ")
}

func NewCmdLine() *CmdLine {
	return &CmdLine{
		Flags: make([]*Flag, 0),
		Items: make([]*Item, 0),
	}
}

func NewCmdLineByArgs(args []string) (*CmdLine, error) {
	cmdLine := NewCmdLine()
	if args[0] == "https" {
		cmdLine.IsHttps = true
	}
	args = args[1:]
	if len(args) == 1 {
		cmdLine.URL = args[0]
		return cmdLine, nil
	}

	var err error
	var pureArgs []string
	cmdLine.Flags, pureArgs, err = removeFlags(args)
	if err != nil {
		return nil, fmt.Errorf("NewCmdLineByArgs: %w", err)
	}

	cmdLine.Method, cmdLine.URL, cmdLine.Items, err = getMethodURLAndItems(pureArgs)
	return cmdLine, nil
}

func getMethodURLAndItems(args []string) (method *Method, url string, items []*Item, err error) {
	method = NewMethod("")

	possibleMethodIndex := 0

	urlIndex := possibleMethodIndex
	possibleMethod := strings.ToUpper(args[possibleMethodIndex])
	if inStringSlice(httpMethods, possibleMethod) {
		method = NewMethod(possibleMethod)
		urlIndex = possibleMethodIndex + 1
	}
	url = args[urlIndex]

	if len(args) > urlIndex+1 {
		items, err = parseItems(args[urlIndex+1:])
	}
	return
}

func parseItems(args []string) ([]*Item, error) {
	items := make([]*Item, 0, len(args))
	for _, arg := range args {
		item, err := getItemByArg(arg)
		if err != nil {
			return nil, err
		}

		if item != nil {
			items = append(items, item)
		}
	}
	return items, nil
}

func getItemByArg(arg string) (*Item, error) {
	for i, r := range arg {
		if i > 0 && arg[i-1] == '\\' {
			continue
		}

		switch r {
		case '@':
			return NewFileField(arg[:i], arg[i+1:]), nil
		case ':':
			if arg[i+1] == '=' {
				return NewJSONField(arg[:i], arg[i+2:]), nil
			}
			return NewHeader(arg[:i], arg[i+1:]), nil
		case '=':
			if arg[i+1] == '=' {
				return NewURLParam(arg[:i], arg[i+2:]), nil
			}
			return NewDataField(arg[:i], arg[i+1:]), nil
		}
	}

	return nil, fmt.Errorf("unknown item")
}

var httpMethods = []string{
	http.MethodDelete,
	http.MethodGet,
	http.MethodHead,
	http.MethodPost,
	http.MethodPut,
	http.MethodPatch,
	http.MethodDelete,
	http.MethodConnect,
	http.MethodOptions,
	http.MethodTrace,
}

func inStringSlice(slice []string, target string) bool {
	for _, v := range slice {
		if v == target {
			return true
		}
	}

	return false
}

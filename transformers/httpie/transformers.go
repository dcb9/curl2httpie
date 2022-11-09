package httpie

import (
	"encoding/json"
	"errors"
	"io/ioutil"
	"strings"

	"github.com/dcb9/curl2httpie/curl"
	"github.com/dcb9/curl2httpie/httpie"
)

type Transformer func(cl *httpie.CmdLine, o *curl.Option)

func Header(cl *httpie.CmdLine, o *curl.Option) {
	s := strings.SplitN(o.Arg, ":", 2)

	if len(s) != 2 {
		return
	}
	k, v := strings.TrimSpace(s[0]), strings.TrimSpace(s[1])
	lk := strings.ToLower(k)
	if lk == "content-type" {
		if strings.HasSuffix(v, "json") {
			cl.ContentType = "json"
			return
		}
	} else if lk == "accept" {
		if strings.HasSuffix(v, "json") {
			cl.ContentType = "json"
			return
		}
	}

	h := httpie.NewHeader(k, v)

	cl.AddItem(h)
}

func Method(cl *httpie.CmdLine, o *curl.Option) {
	m := httpie.NewMethod(o.Arg)
	cl.SetMethod(m)
}

var ErrUnknownDataType = errors.New("unknown data type")

func Data(cl *httpie.CmdLine, o *curl.Option) {
	args := strings.Split(o.Arg, "&")
	var urlEncoded bool
	for _, arg := range args {
		s := strings.SplitN(arg, "=", 2)
		if len(s) == 2 {
			i := httpie.NewDataField(s[0], s[1])
			cl.AddItem(i)
			cl.HasBody = true
			urlEncoded = true
		}
	}
	if urlEncoded {
		return
	}

	// try RAW JSON
	var js json.RawMessage
	err := json.Unmarshal([]byte(o.Arg), &js)
	if err != nil {
		panic(ErrUnknownDataType)
	}

	cl.DirectedInput = ioutil.NopCloser(strings.NewReader(o.Arg))
	cl.HasBody = true
	return
}

func URL(cl *httpie.CmdLine, o *curl.Option) {
	cl.SetURL(o.Arg)
}

func User(cl *httpie.CmdLine, o *curl.Option) {
	f := httpie.NewFlag("auth")
	f.SetArg(o.Arg)

	cl.AddFlag(f)
}

func UserAgent(cl *httpie.CmdLine, o *curl.Option) {
	h := httpie.NewHeader("User-Agent", o.Arg)
	cl.AddItem(h)
}

func Verbose(cl *httpie.CmdLine, o *curl.Option) {
	f := httpie.NewFlag("verbose")

	cl.AddFlag(f)
}

func Referer(cl *httpie.CmdLine, o *curl.Option) {
	h := httpie.NewHeader("Referer", o.Arg)
	cl.AddItem(h)
}

func Cookie(cl *httpie.CmdLine, o *curl.Option) {
	h := httpie.NewHeader("Cookie", o.Arg)
	cl.AddItem(h)
}

func Noop(cl *httpie.CmdLine, o *curl.Option) {
}

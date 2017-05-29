package transformer

import (
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
			f := httpie.NewFlag("json")
			cl.AddFlag(f)
			return
		}
	} else if lk == "accept" {
		if strings.HasSuffix(v, "json") {
			f := httpie.NewFlag("json")
			cl.AddFlag(f)
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

func Data(cl *httpie.CmdLine, o *curl.Option) {
	s := strings.SplitN(o.Arg, "=", 2)
	if len(s) != 2 {
		panic(s)
	}

	i := httpie.NewDataField(s[0], s[1])
	cl.AddItem(i)
	cl.HasBody = true
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

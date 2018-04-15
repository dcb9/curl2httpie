package curl

import (
	"encoding/json"
	"strings"
	"errors"

	"github.com/dcb9/curl2httpie/curl"
	"github.com/dcb9/curl2httpie/httpie"
	"io/ioutil"
)

type ItemTransformer func(*curl.CmdLine, *httpie.Item)
type FlagTransformer func(*curl.CmdLine, *httpie.Flag)

func Method(cl *curl.CmdLine, method *httpie.Method) {
	cl.Options = append(cl.Options, curl.NewMethod(string(*method)))
}

func Auth(cl *curl.CmdLine,  flag *httpie.Flag) {
	cl.Options = append(cl.Options, curl.NewUser(flag.Arg))
}

var ErrUnknownDataType = errors.New("unknown data type")
func Data(cl *httpie.CmdLine, o *curl.Option) {
	s := strings.SplitN(o.Arg, "=", 2)
	if len(s) == 2 {
		i := httpie.NewDataField(s[0], s[1])
		cl.AddItem(i)
		cl.HasBody = true
		return
	}

	// try RAW JSON
	var js json.RawMessage
	err := json.Unmarshal([]byte(o.Arg[1:len(o.Arg)-1]), &js)
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

func Noop(cl *curl.CmdLine, o *httpie.Flag) {
}

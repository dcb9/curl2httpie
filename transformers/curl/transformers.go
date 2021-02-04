package curl

import (
	"fmt"
	"github.com/dcb9/curl2httpie/curl"
	"github.com/dcb9/curl2httpie/httpie"
	"regexp"
)

type ItemTransformer func(*curl.CmdLine, *httpie.Item)
type FlagTransformer func(*curl.CmdLine, *httpie.Flag)

// TransURL supports HTTPie url shortcuts for localhost
func TransURL(url string) string {
	if url[0] == ':' {
		if len(url) == 1 {
			return "localhost/"
		}

		portRe := regexp.MustCompile(`^[0-9]+`)
		port := portRe.Find([]byte(url[1:]))
		if len(port) == 0 {
			return "localhost" + url[1:]
		}

		return fmt.Sprintf("%s:%s%s", "localhost", port, url[len(port)+1:])
	}

	return url
}

func Method(cl *curl.CmdLine, method *httpie.Method) {
	cl.Options = append(cl.Options, curl.NewMethod(string(*method)))
}

func Auth(cl *curl.CmdLine, flag *httpie.Flag) {
	cl.Options = append(cl.Options, curl.NewUser(flag.Arg))
}

func AuthType(cl *curl.CmdLine, flag *httpie.Flag) {
	cl.Options = append(cl.Options, curl.NewNoArgOption(flag.Arg, 0))
}

func Proxy(cl *curl.CmdLine, flag *httpie.Flag) {
	cl.Options = append(cl.Options, curl.NewProxy(flag.Arg))
}

func Follow(cl *curl.CmdLine, flag *httpie.Flag) {
	cl.Options = append(cl.Options, curl.NewLocation())
}

func MaxRedirects(cl *curl.CmdLine, flag *httpie.Flag) {
	cl.Options = append(cl.Options, curl.NewMaxRedirs(flag.Arg))
}

func Timeout(cl *curl.CmdLine, flag *httpie.Flag) {
	cl.Options = append(cl.Options, curl.NewMaxTime(flag.Arg))
}

func Noop(cl *curl.CmdLine, o *httpie.Flag) {
}

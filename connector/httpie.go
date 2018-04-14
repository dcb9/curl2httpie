package connector

import (
	"fmt"

	"github.com/dcb9/curl2httpie/httpie"
	"github.com/dcb9/curl2httpie/curl"
	"strings"
	curlTransformer "github.com/dcb9/curl2httpie/transformers/curl"
)

func Httpie2Curl(args []string) (cmdStringer fmt.Stringer, warningMessages []WarningMessage, err error) {
	httpieInstance, err := httpie.NewCmdLineByArgs(args)
	if err != nil {
		return
	}

	curlCmdLine := curl.NewCmdLine()
	curlCmdLine.URL = httpieInstance.URL
	items := make([]httpie.Itemer, 0, len(httpieInstance.Items))
	queries := make([]string, 0)
	for _, i := range httpieInstance.Items {
		if queryItemer, ok := i.(httpie.QueryItemer); ok {
			queries = append(queries, queryItemer.ToQuery())
			continue
		}

		items = append(items, i)
	}
	if len(queries) > 0 {
		curlCmdLine.URL = fmt.Sprintf("%s?%s", curlCmdLine.URL, strings.Join(queries, "&"))
	}


	// flags
	for _, f := range httpieInstance.Flags {
		t, ok := httpieFlag2CurlOptionTransformerMap[f.Long]
		if !ok {
			warningMessages = append(
				warningMessages,
				WarningMessage(fmt.Sprintf(`skipped: flag "%s" is not supported.`, f.Long)),
			)

			t = curlTransformer.Noop
		}
		t(curlCmdLine, f)
	}

	cmdStringer = curlCmdLine.NewStringer(true)

	return
}

var httpieFlag2CurlOptionTransformerMap = map[string]curlTransformer.FlagTransformer{
	"auth":       curlTransformer.Auth,
}

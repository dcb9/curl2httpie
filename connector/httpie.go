package connector

import (
	"fmt"

	"encoding/json"
	"github.com/dcb9/curl2httpie/curl"
	"github.com/dcb9/curl2httpie/httpie"
	curlTransformer "github.com/dcb9/curl2httpie/transformers/curl"
	"io/ioutil"
	"os"
	"strings"
)

func Httpie2Curl(args []string) (cmdStringer fmt.Stringer, warningMessages []WarningMessage, err error) {
	httpieInstance, err := httpie.NewCmdLineByArgs(args)
	if err != nil {
		return
	}

	curlCmdLine := curl.NewCmdLine()
	if len(httpieInstance.URL) > 0 {
		curlCmdLine.URL = curlTransformer.TransURL(httpieInstance.URL)
	}
	if method := httpieInstance.Method; method != nil && string(*method) != "GET" {
		curlCmdLine.Options = append(curlCmdLine.Options, curl.NewRequest(string(*method)))
	}

	isJSONContentType := true
	// flags
	for _, f := range httpieInstance.Flags {
		if f.Long == "form" {
			isJSONContentType = false
			continue
		}

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

	data := make(map[string]interface{})
	hasData := false
	queries := make([]string, 0)
	fieldOrder := make([]string, 0)
	for _, i := range httpieInstance.Items {
		switch i.S {
		case httpie.SEP_URL_PARAM:
			queries = append(queries, fmt.Sprintf("%s=%s", i.K, i.V))
		case httpie.SEP_HEADER:
			curlCmdLine.Options = append(curlCmdLine.Options, curl.NewHeader(i.K, i.V))
		case httpie.SEP_FILE:
			if isJSONContentType {
				warningMessages = append(
					warningMessages,
					WarningMessage("skipped: file field could not be with JSON content-type"),
				)
			} else {
				hasData = true
				curlCmdLine.Options = append(curlCmdLine.Options, curl.NewForm(fmt.Sprintf(`%s=@"%s"`, i.K, i.V)))
			}
		case httpie.SEP_DATA:
			hasData = true
			fieldOrder = append(fieldOrder, i.K)
			if i.V[0] == '@' {
				var bytes []byte
				bytes, err = getFileContent(i.V[1:])
				if err != nil {
					return
				}
				data[i.K] = string(bytes)
			} else {
				data[i.K] = i.V
			}
		case httpie.SEP_JSON:
			hasData = true
			fieldOrder = append(fieldOrder, i.K)
			if i.V[0] == '@' {
				var bytes []byte
				bytes, err = getFileContent(i.V[1:])
				if err != nil {
					return
				}
				data[i.K] = json.RawMessage(bytes)
			} else {
				data[i.K] = json.RawMessage(i.V)
			}
		}
	}
	if len(queries) > 0 {
		curlCmdLine.URL = fmt.Sprintf("%s?%s", curlCmdLine.URL, strings.Join(queries, "&"))
	}

	if hasData {
		if isJSONContentType {
			var bs []byte
			bs, err = json.Marshal(data)
			if err != nil {
				return
			}
			curlCmdLine.Options = append(curlCmdLine.Options, curl.NewJSONHeader(), curl.NewData(string(bs)))
		} else {
			fields := make([]string, 0, len(data))
			for _, key := range fieldOrder {
				fields = append(fields, fmt.Sprintf("%s=%s", key, data[key]))
			}
			curlCmdLine.Options = append(curlCmdLine.Options, curl.NewData(strings.Join(fields, "&")))
		}
	}

	cmdStringer = curlCmdLine.NewStringer(true)

	return
}

var httpieFlag2CurlOptionTransformerMap = map[string]curlTransformer.FlagTransformer{
	"auth":          curlTransformer.Auth,
	"auth-type":     curlTransformer.AuthType,
	"proxy":         curlTransformer.Proxy,
	"follow":        curlTransformer.Follow,
	"max-redirects": curlTransformer.MaxRedirects,
	"timeout":       curlTransformer.Timeout,
}

func getFileContent(filename string) ([]byte, error) {
	f, err := os.Open(filename)
	if err != nil {
		return nil, err
	}

	return ioutil.ReadAll(f)
}

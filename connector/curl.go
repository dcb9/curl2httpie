package connector

import (
	"fmt"

	"github.com/dcb9/curl2httpie/curl"
	"github.com/dcb9/curl2httpie/httpie"

	httpieTransformer "github.com/dcb9/curl2httpie/transformers/httpie"
)

var curl2HttpieTransformerMap = map[curl.LongName]httpieTransformer.Transformer{
	"header":     httpieTransformer.Header,
	"request":    httpieTransformer.Method,
	"data":       httpieTransformer.Data,
	"url":        httpieTransformer.URL,
	"user":       httpieTransformer.User,
	"user-agent": httpieTransformer.UserAgent,
	"verbose":    httpieTransformer.Verbose,
	"referer":    httpieTransformer.Referer,
	"cookie":     httpieTransformer.Cookie,
}

func Curl2Httpie(args []string) (cmdStringer fmt.Stringer, warningMessages []WarningMessage, err error) {
	defer func() {
		if recoverErr := recover(); recoverErr != nil {
			err = recoverErr.(error)
			return
		}
	}()

	warningMessages = make([]WarningMessage, 0)
	url, curlOptions := "", make([]*curl.Option, 0, len(args))
	cmdline := httpie.NewCmdLine()

	url, curlOptions, err = curl.URLAndOptions(args)
	if err != nil {
		return
	}

	cmdline.SetURL(url)
	if curlOptions != nil {
		for _, o := range curlOptions {
			t, ok := curl2HttpieTransformerMap[o.Long]
			if !ok {
				warningMessages = append(
					warningMessages,
					WarningMessage(fmt.Sprintf(`skipped: option "%s" is not supported.`, o.Long)),
				)

				t = httpieTransformer.Noop
			}
			t(cmdline, o)
		}
	}

	return cmdline, warningMessages, nil
}

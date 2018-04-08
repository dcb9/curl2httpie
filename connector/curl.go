package connector

import (
	"fmt"
	"errors"

	"github.com/dcb9/curl2httpie/curl"
	"github.com/dcb9/curl2httpie/httpie"
	"github.com/dcb9/curl2httpie/transformer"
)

var curl2HttpieTransformerMap = map[curl.LongName]transformer.Transformer{
	"header":     transformer.Header,
	"request":    transformer.Method,
	"data":       transformer.Data,
	"url":        transformer.URL,
	"user":       transformer.User,
	"user-agent": transformer.UserAgent,
	"verbose":    transformer.Verbose,
	"referer":    transformer.Referer,
	"cookie":     transformer.Cookie,
}

var ErrNoEnoughArgs = errors.New("no enough args")

func Curl2Httpie(args []string) (cmdStringer fmt.Stringer, warningMessages []WarningMessage, err error) {
	defer func() {
		if recoverErr := recover(); recoverErr!=nil {
			err = recoverErr.(error)
			return
		}
	}()

	warningMessages = make([]WarningMessage, 0)
	url, curlOptions := "", make([]*curl.Option, 0, len(args))
	cmdline := httpie.NewCmdLine()
	if len(args) < 2 {
		return nil, warningMessages, ErrNoEnoughArgs
	} else if len(args) == 2 {
		cmdline.SetURL(args[1])
		return cmdline, warningMessages, nil
	} else {
		args = args[1:]
		url, curlOptions , err = curl.URLAndOptions(args)
		if err != nil {
			return nil, warningMessages, err
		}
	}

	cmdline.SetURL(url)
	for _, o := range curlOptions {
		t, ok := curl2HttpieTransformerMap[o.Long]
		if !ok {
			warningMessages = append(
				warningMessages,
				WarningMessage(fmt.Sprintf(`skipped: option "%s" is not supported.`, o.Long)),
			)

			t = transformer.Noop
		}
		t(cmdline, o)
	}

	return cmdline, warningMessages, nil
}


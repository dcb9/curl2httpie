package connector

import (
	"github.com/dcb9/curl2httpie/curl"
	"github.com/dcb9/curl2httpie/httpie"
	"github.com/dcb9/curl2httpie/transformer"
)

var transformerMap = map[curl.LongName]transformer.Transformer{
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

func Curl2Httpie(args []string) string {
	url, curlOptions := "", make([]*curl.Option, 0, len(args))
	cmdline := httpie.NewCmdLine()
	if len(args) < 2 {
		panic("No enough args")
	} else if len(args) == 2 {
		cmdline.SetURL(args[1])
		return cmdline.String()
	} else {
		args = args[1:]
		url, curlOptions = curl.URLAndOptions(args)
	}

	cmdline.SetURL(url)
	for _, o := range curlOptions {
		transformerMap[o.Long](cmdline, o)
	}

	return cmdline.String()
}

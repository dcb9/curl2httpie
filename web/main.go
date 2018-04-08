package main

import (
	"github.com/gopherjs/gopherjs/js"
	"github.com/lmika/shellwords"
	"github.com/dcb9/curl2httpie/connector"
)

func Httpie(cmd string) string {
	cmdStringer, _, _ := connector.Curl2Httpie(shellwords.Split(cmd))
	return cmdStringer.String()
}

func main() {
	js.Global.Set("curl2httpie", map[string]interface{}{
		"Do": Httpie,
	})
}

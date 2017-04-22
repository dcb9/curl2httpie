package main

import (
	"github.com/dcb9/curl2httpie"
	"github.com/gopherjs/gopherjs/js"
	"github.com/lmika/shellwords"
)

func Httpie(cmd string) string {
	str, err := curl2httpie.Httpie(shellwords.Split("web " + cmd))
	if err != nil {
		js.Global.Call("alert", err.Error())
		return ""
	}
	return str
}

func main() {
	js.Global.Set("curl2httpie", map[string]interface{}{
		"Do": Httpie,
	})
}

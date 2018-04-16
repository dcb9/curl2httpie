package main

import (
	"github.com/dcb9/curl2httpie/connector"
	"github.com/gopherjs/gopherjs/js"
	"github.com/lmika/shellwords"
)

func Httpie(cmd string) map[string]interface{} {
	stringer, warnings, err := connector.Convert(shellwords.Split(cmd))
	if err != nil {
		return map[string]interface{}{
			"cmd":      "",
			"warnings": warnings,
			"error":    err.Error(),
		}
	}

	return map[string]interface{}{
		"cmd":      stringer.String(),
		"warnings": warnings,
		"error":    "",
	}
}

func main() {
	js.Global.Set("curl2httpie", map[string]interface{}{
		"Do": Httpie,
	})
}

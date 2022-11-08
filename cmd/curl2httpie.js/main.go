package main

import (
	"fmt"
	"regexp"

	"github.com/dcb9/curl2httpie/connector"
	"github.com/dcb9/curl2httpie/constant"
	"github.com/dcb9/curl2httpie/shellwords"
	"github.com/gopherjs/gopherjs/js"
)

var toOneLineReg = regexp.MustCompile(`[\\]{1}\W*\n`)

func Version() string {
	return fmt.Sprintf("Build at <i>%s</i> Version: <i>%s</i> Commit: <a href=\"https://github.com/dcb9/curl2httpie/commit/%s\">%s</a>", constant.BuildAt, constant.Version, constant.Commit, constant.Commit[:7])
}

func Httpie(cmd string) map[string]interface{} {
	cmdBytes := toOneLineReg.ReplaceAll([]byte(cmd), nil)
	stringer, warnings, err := connector.Convert(shellwords.New(string(cmdBytes)).Split())
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
		"Do":      Httpie,
		"Version": Version,
	})
}

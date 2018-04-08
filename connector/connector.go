package connector

import (
	"errors"
	"fmt"
)

var ErrUnknownCommandType = errors.New("the command must begin with 'curl'")

type WarningMessage string

func Convert(args []string) (cmdStringer fmt.Stringer, warningMessages []WarningMessage, err error) {
	switch args[0] {
	case "curl":
		return Curl2Httpie(args)
	}

	err =	ErrUnknownCommandType
	return
}

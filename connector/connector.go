package connector

import (
	"errors"
	"fmt"
)

var (
	ErrUnknownCommandType = errors.New("the command must begin with 'curl'")
	ErrNoEnoughArgs       = errors.New("no enough args")
)

type WarningMessage string

func Convert(args []string) (cmdStringer fmt.Stringer, warningMessages []WarningMessage, err error) {
	if len(args) < 2 {
		err = ErrNoEnoughArgs
		return
	}

	switch args[0] {
	case "curl":
		return Curl2Httpie(args[1:])
	case "http":
		return Httpie2Curl(args[1:])
	}

	err = ErrUnknownCommandType
	return
}

package curl

import (
	"fmt"
	"strings"
)

type CmdLine struct {
	Options []*Option
	URL     string
}

func (cmdline *CmdLine) NewStringer(useLongName bool) fmt.Stringer {
	return &CmdLineStringer{
		cmdline,
		useLongName,
	}
}

type CmdLineStringer struct {
	*CmdLine
	useLongName bool
}

func (cmdlineStringer *CmdLineStringer) String() string {
	parts := []string{"curl"}

	options := make([]string, len(cmdlineStringer.Options))
	for i, o := range cmdlineStringer.Options {
		options[i] = o.String(cmdlineStringer.useLongName)
	}

	if len(cmdlineStringer.Options) > 0 {
		parts = append(parts, strings.Join(options, " "))
	}
	parts = append(parts, cmdlineStringer.URL)

	return strings.Join(parts, " ")
}

func NewCmdLine() *CmdLine {
	return &CmdLine{
		Options: make([]*Option, 0),
	}
}

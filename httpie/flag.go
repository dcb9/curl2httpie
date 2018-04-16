package httpie

import (
	"fmt"
	"github.com/pkg/errors"
	flag "github.com/spf13/pflag"
)

type Flag struct {
	Long   string
	Short  byte
	HasArg bool
	Arg    string
}

func (f *Flag) SetShort(s byte) {
	f.Short = s
}

func (f *Flag) SetArg(arg string) {
	f.HasArg = true
	f.Arg = arg
}

func (f *Flag) String() string {
	arg := ""
	if f.HasArg {
		arg = fmt.Sprintf(` "%s"`, f.Arg)
	}
	return fmt.Sprintf("--%s%s", f.Long, arg)
}

func NewFlag(l string) *Flag {
	return &Flag{Long: l}
}

func AuthFlagWithArg(auth string) *Flag {
	f := &Flag{}
	*f = *AuthFlag
	f.Arg = auth
	return f
}

var JSONFlag = &Flag{ Long: "json", Short: 'j'}
var FormFlag = &Flag{ Long: "form", Short: 'f'}
var AuthFlag = &Flag{ Long: "auth", Short: 'a', HasArg: true}
var AuthTypeFlag = &Flag{ Long: "auth-type", Short: 'A', HasArg: true}
var ProxyFlag = &Flag{ Long: "proxy", HasArg: true}
var FollowFlag = &Flag{ Long: "follow", Short:'F', HasArg: false}
var MaxRedirectsFlag = &Flag{ Long: "max-redirects", HasArg: true}
var TimeoutFlag = &Flag{ Long: "timeout", HasArg: true}

var AllFlags = []*Flag{
	JSONFlag,
	FormFlag,
	AuthFlag,
	AuthTypeFlag,
	ProxyFlag,
	FollowFlag,
	MaxRedirectsFlag,
	TimeoutFlag,
}

func getFlagsByArgs(args []string) ([]*Flag, error) {
	CommandLine := flag.NewFlagSet("httpie", flag.ContinueOnError)
	boolValues := make([]*bool, len(AllFlags))
	stringValues := make([]*string, len(AllFlags))
	for i, f := range AllFlags {
		if f.HasArg {
			if f.Short != 0 {
				stringValues[i] = CommandLine.StringP(f.Long, string(f.Short), "", "")
			} else {
				stringValues[i] = CommandLine.String(f.Long, "", "")
			}
		} else {
			if f.Short != 0 {
				boolValues[i] = CommandLine.BoolP(f.Long, string(f.Short), false, "")
			} else {
				boolValues[i] = CommandLine.Bool(f.Long, false, "")
			}
		}
	}
	err := CommandLine.Parse(args)
	if err != nil {
		return nil, errors.Wrap(err, "GetFlagsByArgs")
	}
	flags := make([]*Flag, 0, len(args))
	for i, f := range AllFlags {
		if f.HasArg {
			if val := *stringValues[i]; val != "" {
				f.Arg = val
				flags = append(flags, f)
			}
		} else {
			if *boolValues[i] == true {
				flags = append(flags, f)
			}
		}
	}
	return flags, nil
}

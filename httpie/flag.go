package httpie

import "fmt"

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

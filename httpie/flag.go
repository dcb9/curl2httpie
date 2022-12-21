package httpie

import (
	"fmt"

	flag "github.com/spf13/pflag"
)

type Flag struct {
	Long      string
	Short     byte
	HasArg    bool
	Arg       string
	Separator string
}

func (f *Flag) SetShort(s byte) {
	f.Short = s
}

func (f *Flag) SetArg(arg string, separator string) {
	f.HasArg = true
	f.Arg = arg
	f.Separator = separator
}

func (f *Flag) String() string {
	arg := ""
	if f.HasArg {
		if f.Separator == "" {
			f.Separator = " " // Use whitespace as default separator
		}
		arg = fmt.Sprintf(`%s%s`, f.Separator, addQuoteIfNeeded(f.Arg))
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

var (
	JSONFlag            = &Flag{Long: "json", Short: 'j'}
	FormFlag            = &Flag{Long: "form", Short: 'f'}
	PrettyFlag          = &Flag{Long: "pretty", HasArg: true}
	StyleFlag           = &Flag{Long: "style", Short: 's', HasArg: true}
	PrintFlag           = &Flag{Long: "print", Short: 'p', HasArg: true}
	HeadersFlag         = &Flag{Long: "headers", Short: 'h'}
	BodyFlag            = &Flag{Long: "body", Short: 'b'}
	VerboseFlag         = &Flag{Long: "verbose", Short: 'v'}
	AllFlag             = &Flag{Long: "all"}
	HistoryPrintFlag    = &Flag{Long: "history-print", Short: 'P'}
	StreamFlag          = &Flag{Long: "stream", Short: 'S'}
	OutputFlag          = &Flag{Long: "output", Short: 'o'}
	DownloadFlag        = &Flag{Long: "download", Short: 'd'}
	ContinueFlag        = &Flag{Long: "continue", Short: 'c'}
	SessionFlag         = &Flag{Long: "session", HasArg: true}
	SessionReadOnlyFlag = &Flag{Long: "session-read-only", HasArg: true}
	CheckStatusFlag     = &Flag{Long: "check-status"}
	IgnoreStdinFlag     = &Flag{Long: "ignore-stdin", Short: 'I'}
	HelpFlag            = &Flag{Long: "help"}
	VersionFlag         = &Flag{Long: "version"}
	TracebackFlag       = &Flag{Long: "traceback"}
	DefaultSchemeFlag   = &Flag{Long: "default-scheme", HasArg: true}
	DebugFlag           = &Flag{Long: "debug"}
	AuthFlag            = &Flag{Long: "auth", Short: 'a', HasArg: true}
	AuthTypeFlag        = &Flag{Long: "auth-type", Short: 'A', HasArg: true}
	ProxyFlag           = &Flag{Long: "proxy", HasArg: true}
	FollowFlag          = &Flag{Long: "follow", Short: 'F', HasArg: false}
	MaxRedirectsFlag    = &Flag{Long: "max-redirects", HasArg: true}
	TimeoutFlag         = &Flag{Long: "timeout", HasArg: true}
	VerifyFlag          = &Flag{Long: "verify", HasArg: true}
	SSLFlag             = &Flag{Long: "ssl", HasArg: true}
	CertFlag            = &Flag{Long: "cert", HasArg: true}
	CertKeyFlag         = &Flag{Long: "cert-key", HasArg: true}
)

var AllFlags = []*Flag{
	JSONFlag,
	FormFlag,
	PrettyFlag,
	StyleFlag,
	PrintFlag,
	HeadersFlag,
	BodyFlag,
	VerboseFlag,
	AllFlag,
	HistoryPrintFlag,
	StreamFlag,
	OutputFlag,
	DownloadFlag,
	ContinueFlag,
	SessionFlag,
	SessionReadOnlyFlag,
	CheckStatusFlag,
	IgnoreStdinFlag,
	HelpFlag,
	VersionFlag,
	TracebackFlag,
	DefaultSchemeFlag,
	DebugFlag,
	AuthFlag,
	AuthTypeFlag,
	ProxyFlag,
	FollowFlag,
	MaxRedirectsFlag,
	TimeoutFlag,
	VerifyFlag,
	SSLFlag,
	CertFlag,
	CertKeyFlag,
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
		return nil, fmt.Errorf("GetFlagsByArgs: %w", err)
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

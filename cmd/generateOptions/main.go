package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/dcb9/curl2httpie/curl"
)

func main() {
	path := flag.String("path", "", "curl source code path; you could clone from github.com (curl/curl)")
	flag.Parse()

	*path = strings.TrimRight(*path, "/") + "/docs/cmdline-opts/"
	options := curl.GenerateHTTPOptions(*path)
	if len(options) < 50 {
		fmt.Fprintf(os.Stderr, "Too few options found: %d\n", len(options))
		os.Exit(1)
	}

	optionsCode := ""

	for _, o := range options {
		hasArg := "false"
		if o.HasArg {
			hasArg = "true"
		}
		tags := ""
		for _, s := range o.Tags {
			tags += fmt.Sprintf("`%s`,", s)
		}

		protocols := ""
		for _, s := range o.Protocols {
			protocols += fmt.Sprintf("`%s`,", s)
		}

		mutexed := ""
		for _, s := range o.Mutexed {
			mutexed += fmt.Sprintf("`%s`,", s)
		}

		requires := ""
		for _, s := range o.Requires {
			requires += fmt.Sprintf("`%s`,", s)
		}

		optionsCode += fmt.Sprintf(`{
Short: %d,
Long: "%s",
HasArg: %s,
Arg: "%s",
Magic: "%s",
Tags:      []Tag{%s},
Protocols: []Protocol{%s},
Added:     %s,
Mutexed:   []LongName{%s},
Requires:  []Feature{%s},
},`, o.Short, o.Long, hasArg, o.Arg, o.Magic, tags, protocols,
			"`"+o.Added+"`", mutexed, requires)
	}

	code := fmt.Sprintf(`package curl

func init() {
    optionList = []*Option{ %s }
}
`, optionsCode)

	if err := ioutil.WriteFile("./curl/optionList.go", []byte(code), 0644); err != nil {
		fmt.Fprintf(os.Stderr, "fail to persist curl options: %v\n", err)
		os.Exit(1)
	}
}

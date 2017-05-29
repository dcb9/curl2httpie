package main

import (
	"strings"
	"flag"
	"github.com/dcb9/curl2httpie/curl"
	"encoding/json"
	"fmt"
	"os"
	"io/ioutil"
)

func main() {
	path := flag.String("path", "", "curl source code path; you could clone from github.com (curl/curl)")
	flag.Parse()

	*path = strings.TrimRight(*path, "/") + "/docs/cmdline-opts/"
	options := curl.GenerateHTTPOptions(*path)
	bs, err := json.Marshal(options)
	if err != nil {
		panic(err)
	}
	if len(options) < 50 {
		fmt.Fprintf(os.Stderr,"Too few options found: %d\n", len(options))
		os.Exit(1)
	}

	err = ioutil.WriteFile("./data/options.json", bs, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "writer json file: %v\n", err)
		os.Exit(1)
	}
}

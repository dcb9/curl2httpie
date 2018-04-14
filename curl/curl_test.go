package curl

import (
	"encoding/json"
	"fmt"
	"testing"
)

func TestDotD2Option(t *testing.T) {
	content := `
Short: 0
Long: http1.0
Tags: Versions
Protocols: HTTP
Added:
Mutexed: http1.1 http2
Help: Use HTTP 1.0
---
Tells curl to use HTTP version 1.0 instead of using its internally preferred
HTTP version.
`
	option := DotD2Option(content)
	bytes, err := json.Marshal(option)
	fmt.Println(string(bytes))
	fmt.Println(err)
}

func TestFiles(t *testing.T) {
	options, err := HTTPOptions()
	if err != nil {
		t.Fatalf("HTTPOptions error: %s", err.Error())
	}

	for _, o := range options {
		fmt.Printf("%s\n", o.dumpStructure())
	}
}

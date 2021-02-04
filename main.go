package main

import (
	"fmt"
	"github.com/dcb9/curl2httpie/constant"
	"log"
	"os"

	"github.com/dcb9/curl2httpie/connector"
)

var errLogger = log.New(os.Stderr, "", 0)

func main() {
	if os.Args[1] == "-v" || os.Args[1] == "--version" || os.Args[1] == "version" {
		fmt.Printf("curl2httpie version: %s commit %s\n", constant.Version, constant.Commit)
		return
	}

	cmdStr, warningMessages, err := connector.Convert(os.Args[1:])
	if len(warningMessages) > 0 {
		errLogger.Println("warnings:")
		for i, message := range warningMessages {
			errLogger.Printf("\t%d. %s\n", i+1, message)
		}
		errLogger.Println()
	}

	if err != nil {
		errLogger.Fatalf("convert error: %s\n", err.Error())
	}

	fmt.Println(cmdStr)
}

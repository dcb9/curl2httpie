package main

import (
	"fmt"
	"log"
	"os"

	"github.com/dcb9/curl2httpie/connector"
)

var errLogger = log.New(os.Stderr, "", 0)

func main() {
	fmt.Println(os.Args)
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

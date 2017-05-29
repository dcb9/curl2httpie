package main

import (
	"fmt"
	"os"

	"github.com/dcb9/curl2httpie/connector"
)

func main() {
	httpie := connector.Curl2Httpie(os.Args[1:])
	fmt.Println(httpie)
}

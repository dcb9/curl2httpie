package main

import (
	"fmt"
	"os"

	"github.com/dcb9/curl2httpie"
)

func main() {
	httpie, err := curl2httpie.Httpie(os.Args)
	if err != nil {
		fmt.Println("error: ", err)
		os.Exit(2)
	}
	fmt.Println(httpie)
}

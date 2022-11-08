package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/dcb9/curl2httpie/constant"
)

func main() {
	// var Version = "v0.0.0"
	// var Commit = "00000"
	// var BuildAt = "yyy-mm-dd HH:ii:ss"
	data, err := ioutil.ReadFile("constant/constant.go")
	if err != nil {
		fmt.Println("Read constant from file error:", err)
		os.Exit(-1)
	}

	data = bytes.Replace(data, []byte(`Version = "v0.0.0"`), []byte(fmt.Sprintf(`Version = "%s"`, constant.Version)), 1)
	data = bytes.Replace(data, []byte(`Commit = "00000"`), []byte(fmt.Sprintf(`Commit = "%s"`, constant.Commit)), 1)
	data = bytes.Replace(data, []byte(`BuildAt = "yyy-mm-dd HH:ii:ss"`), []byte(fmt.Sprintf(`BuildAt = "%s"`, constant.BuildAt)), 1)

	f, err := os.OpenFile("constant/constant.go", os.O_TRUNC|os.O_WRONLY, 0644)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	fmt.Println(string(data))
	_, err = f.Write(data)
	if err != nil {
		panic(err)
	}
}

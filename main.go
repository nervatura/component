package main

import (
	"fmt"

	"github.com/nervatura/component/pkg/demo"
)

var (
	version = "dev"
)

func main() {
	fmt.Printf("Version: %s\n", version)
	err := demo.New(version)
	if err != nil {
		fmt.Println(err.Error())
	}

}

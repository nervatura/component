package main

import (
	"fmt"
	"os"

	"github.com/nervatura/component/pkg/demo"
	ut "github.com/nervatura/component/pkg/util"
)

const httpPort = int64(5000)

var (
	version = "dev"
)

func main() {
	fmt.Printf("Version: %s\n", version)
	port := httpPort
	if len(os.Args) > 1 {
		port = ut.ToInteger(os.Args[1], httpPort)
	}
	demo.New(version, port)
}

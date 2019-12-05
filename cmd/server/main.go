package main

import (
	"fmt"
	cmd "github.com/donggangcj/go-grpc-http-rest-microservice-tutorial/pkg/cmd/server"
	"os"
)

func main() {
	if err := cmd.RunServer(); err != nil {
		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}

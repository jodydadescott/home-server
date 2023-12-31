package main

import (
	"fmt"
	"os"

	"github.com/jodydadescott/home-server/cmd"
)

func main() {
	err := cmd.Execute()

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error())
		os.Exit(1)
	}
}

package main

import (
	"github.com/magazin/command"
	"os"
)

func main() {
	if err := command.Main(os.Args); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"os"
)

func main() {
	command := generate.Main()
	args := os.Args

	if args != nil {
		args = args[1:]
		command.SetArgs(args)
	}

	err := command.Execute()
	if err != nil {
		os.Exit(1)
	}
}

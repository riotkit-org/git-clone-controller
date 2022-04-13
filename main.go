package main

import (
	"github.com/riotkit-org/git-clone-operator/cmd"
	"os"
)

func main() {
	command := cmd.Main()
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

package main

import (
	"os"

	"github.com/ankitiscracked/jmp/cmd/jmp/commands"
)

func main() {
	if err := commands.Execute(); err != nil {
		if code := commands.ExitCode(err); code != 0 {
			os.Exit(code)
		}
		os.Exit(1)
	}
}

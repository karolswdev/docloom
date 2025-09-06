package main

import (
	"os"

	"github.com/karolswdev/docloom/tools/claude-code-cli/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		os.Exit(1)
	}
}

package main

import (
	"fmt"
	"os"

	"github.com/karolswdev/docloom/internal/cli"
	"github.com/karolswdev/docloom/internal/version"
)

func main() {
	// Check for --version flag early
	for _, arg := range os.Args[1:] {
		if arg == "--version" || arg == "-version" {
			fmt.Println(version.Info())
			os.Exit(0)
		}
	}
	
	if err := cli.Execute(); err != nil {
		os.Exit(1)
	}
}
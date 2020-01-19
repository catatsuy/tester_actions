package main

import (
	"os"

	"github.com/catatsuy/bento/cli"
)

func main() {
	cli := cli.NewCLI(os.Stdout, os.Stderr)
	os.Exit(cli.Run(os.Args))
}

package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/padok-team/tfautomv/internal/tfautomv"
)

var (
	debug = flag.Bool("debug", false, "enable verbose output")
)

func main() {
	if err := run(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v", err)
		os.Exit(1)
	}
}

func run() error {
	flag.Parse()

	fmt.Fprintln(os.Stderr, "Generating moved blocks...")

	n, err := tfautomv.GenerateMovedBlocks(".")
	if err != nil {
		return err
	}

	fmt.Fprintf(os.Stderr, "Generated %d moved blocks.\n", n)

	return nil
}

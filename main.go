package main

import (
	"fmt"
	"os"

	"github.com/akornatskyy/precompress/cmd/cli"
	"github.com/akornatskyy/precompress/walker"
)

func main() {
	opts, err := cli.ParseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing arguments:", err)
		os.Exit(1)
	}

	w, err := walker.New(
		walker.MinSize(opts.MinSize),
		walker.MaxDepth(opts.MaxDepth),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error initializing walker:", err)
		os.Exit(1)
	}

	if err = w.Walk(opts.Paths); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

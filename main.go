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

	err = walker.Walk(
		walker.MinSize(opts.MinSize),
		walker.MaxDepth(opts.MaxDepth),
		walker.Paths(opts.Args),
	)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

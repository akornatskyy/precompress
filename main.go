package main

import (
	"fmt"
	"os"

	"github.com/akornatskyy/precompress/cmd/cli"
	"github.com/akornatskyy/precompress/compressors/brotli"
	"github.com/akornatskyy/precompress/compressors/gzip"
	"github.com/akornatskyy/precompress/compressors/zstd"
	"github.com/akornatskyy/precompress/precompressor"
	"github.com/akornatskyy/precompress/walker"
)

func main() {
	opts, err := cli.ParseArgs()
	if err != nil {
		fmt.Fprintln(os.Stderr, "Error parsing arguments:", err)
		os.Exit(1)
	}

	if err = Run(&opts); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

func Run(opts *cli.Options) error {
	w, err := walker.New(
		walker.MinSize(opts.MinSize),
		walker.MaxDepth(opts.MaxDepth),
		walker.Exclude([]string{".gz", ".br", ".zst"}),
	)
	if err != nil {
		return err
	}

	p, err := precompressor.New(
		precompressor.BrotliCompressionLevel(brotli.DefaultCompression),
		precompressor.GzipCompressionLevel(gzip.DefaultCompression),
		precompressor.ZstdCompressionLevel(zstd.DefaultCompression),
	)
	if err != nil {
		return err
	}

	return w.Walk(opts.Paths, p.Precompress)
}

package walker

import (
	"fmt"
	"slices"

	"github.com/akornatskyy/precompress/compressors"
	"github.com/akornatskyy/precompress/compressors/brotli"
	"github.com/akornatskyy/precompress/compressors/gzip"
	"github.com/akornatskyy/precompress/compressors/zstd"
)

type config struct {
	minSize   int64
	maxDepth  int
	providers map[string]compressors.CompressorProvider
	exclude   []string
	paths     []string
}

type Option func(c *config) error

func errorOption(err error) Option {
	return func(_ *config) error {
		return err
	}
}

func MinSize(size int64) Option {
	return func(c *config) error {
		if size < 0 {
			return fmt.Errorf("minimum size can not be negative: %d", size)
		}
		c.minSize = size
		return nil
	}
}

func MaxDepth(depth int) Option {
	return func(c *config) error {
		if depth < 0 {
			return fmt.Errorf("max depth can not be negative: %d", depth)
		}
		c.maxDepth = depth
		return nil
	}
}

func Exclude(exclude []string) Option {
	return func(c *config) error {
		for _, ext := range exclude {
			if !slices.Contains(c.exclude, ext) {
				c.exclude = append(c.exclude, ext)
			}
		}

		return nil
	}
}

func Paths(paths []string) Option {
	return func(c *config) error {
		if len(paths) == 0 {
			return fmt.Errorf("no input files or directories provided")
		}
		c.paths = paths
		return nil
	}
}

func BrotliCompressionLevel(level int) Option {
	p, err := brotli.New(brotli.Options{Quality: level})
	if err != nil {
		return errorOption(err)
	}

	return func(c *config) error {
		c.providers[".br"] = p
		return nil
	}
}

func GzipCompressionLevel(level int) Option {
	p, err := gzip.New(gzip.Options{Level: level})
	if err != nil {
		return errorOption(err)
	}

	return func(c *config) error {
		c.providers[".gz"] = p
		return nil
	}
}

func ZstdCompressionLevel(level int) Option {
	p, err := zstd.New(zstd.Options{Level: level})
	if err != nil {
		return errorOption(err)
	}

	return func(c *config) error {
		c.providers[".zst"] = p
		return nil
	}
}

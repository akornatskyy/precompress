package walker

import (
	"fmt"
	"slices"

	"github.com/akornatskyy/precompress/compressors"
	"github.com/akornatskyy/precompress/compressors/brotli"
	"github.com/akornatskyy/precompress/compressors/gzip"
	"github.com/akornatskyy/precompress/compressors/zstd"
)

type walker struct {
	minSize   int64
	maxDepth  int
	providers map[string]compressors.CompressorProvider
	exclude   []string
}

type Option func(c *walker) error

func MinSize(size int64) Option {
	return func(w *walker) error {
		if size < 0 {
			return fmt.Errorf("minimum size can not be negative: %d", size)
		}
		w.minSize = size
		return nil
	}
}

func MaxDepth(depth int) Option {
	return func(w *walker) error {
		if depth < 0 {
			return fmt.Errorf("max depth can not be negative: %d", depth)
		}
		w.maxDepth = depth
		return nil
	}
}

func Exclude(exclude []string) Option {
	return func(w *walker) error {
		for _, ext := range exclude {
			if !slices.Contains(w.exclude, ext) {
				w.exclude = append(w.exclude, ext)
			}
		}

		return nil
	}
}

func BrotliCompressionLevel(level int) Option {
	return func(w *walker) error {
		p, err := brotli.New(brotli.Options{Quality: level})
		if err != nil {
			return err
		}
		w.providers[".br"] = p
		return nil
	}
}

func GzipCompressionLevel(level int) Option {
	return func(w *walker) error {
		p, err := gzip.New(gzip.Options{Level: level})
		if err != nil {
			return err
		}
		w.providers[".gz"] = p
		return nil
	}
}

func ZstdCompressionLevel(level int) Option {
	return func(w *walker) error {
		p, err := zstd.New(zstd.Options{Level: level})
		if err != nil {
			return err
		}
		w.providers[".zst"] = p
		return nil
	}
}

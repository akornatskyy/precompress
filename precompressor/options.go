package precompressor

import (
	"github.com/akornatskyy/precompress/compressors/brotli"
	"github.com/akornatskyy/precompress/compressors/gzip"
	"github.com/akornatskyy/precompress/compressors/zstd"
)

type Option func(c *Precompressor) error

func BrotliCompressionLevel(level int) Option {
	return func(w *Precompressor) error {
		p, err := brotli.New(brotli.Options{Quality: level})
		if err != nil {
			return err
		}
		w.providers[".br"] = p
		return nil
	}
}

func GzipCompressionLevel(level int) Option {
	return func(w *Precompressor) error {
		p, err := gzip.New(gzip.Options{Level: level})
		if err != nil {
			return err
		}
		w.providers[".gz"] = p
		return nil
	}
}

func ZstdCompressionLevel(level int) Option {
	return func(w *Precompressor) error {
		p, err := zstd.New(zstd.Options{Level: level})
		if err != nil {
			return err
		}
		w.providers[".zst"] = p
		return nil
	}
}

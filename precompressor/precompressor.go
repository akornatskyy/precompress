package precompressor

import (
	"bytes"
	"errors"
	"io"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/akornatskyy/precompress/compressors"
)

type Precompressor struct {
	providers map[string]compressors.CompressorProvider
}

func New(opts ...Option) (*Precompressor, error) {
	w := &Precompressor{providers: map[string]compressors.CompressorProvider{}}
	for _, o := range opts {
		if err := o(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (p *Precompressor) Precompress(path string, fi fs.FileInfo) error {
	ts := fi.ModTime()

	var source []byte
	var compressed = &bytes.Buffer{}

	for ext, provider := range p.providers {
		target_path := path + ext
		tfi, err := os.Stat(target_path)
		if err == nil {
			if tfi.ModTime().After(ts) {
				continue
			}
		} else if !os.IsNotExist(err) {
			return err
		}

		if source == nil {
			if source, err = os.ReadFile(path); err != nil {
				return err
			}
		}

		compressed.Reset()
		compressor := provider.Get(compressed)
		if _, err = compressor.Write(source); err != nil {
			return err
		}
		if err = compressor.Close(); err != nil {
			return err
		}

		if tfi != nil && tfi.Size() == int64(compressed.Len()) {
			existing, err := os.ReadFile(target_path)
			if err != nil {
				return err
			}

			if bytes.Equal(existing, compressed.Bytes()) {
				continue
			}
		}

		if err = writeToFile(target_path, compressed); err != nil {
			return err
		}
	}

	return nil
}

func writeToFile(path string, buffer *bytes.Buffer) error {
	f, err := os.CreateTemp(filepath.Dir(path), "precompress-*.tmp")
	if err != nil {
		return err
	}

	if _, err = io.Copy(f, buffer); err != nil {
		return errors.Join(err, os.Remove(f.Name()))
	}

	if err = f.Close(); err != nil {
		return err
	}

	return os.Rename(f.Name(), path)
}

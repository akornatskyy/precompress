package walker

import (
	"bytes"
	"io"
	"io/fs"
	"os"
)

func process(c *config, path string, fi fs.FileInfo) error {
	ts := fi.ModTime()

	var source []byte
	var compressed = &bytes.Buffer{}

	for ext, provider := range c.providers {
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
		if _, err := compressor.Write(source); err != nil {
			return err
		}
		if err := compressor.Close(); err != nil {
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

		f, err := os.Create(target_path)
		if err != nil {
			return err
		}
		defer f.Close()
		if _, err = io.Copy(f, compressed); err != nil {
			return err
		}
	}

	return nil
}

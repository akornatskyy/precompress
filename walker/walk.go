package walker

import (
	"fmt"
	"io/fs"
	"log"
	"path/filepath"
	"runtime"
	"strings"
	"sync"

	"github.com/akornatskyy/precompress/compressors"
	"github.com/akornatskyy/precompress/compressors/brotli"
	"github.com/akornatskyy/precompress/compressors/gzip"
	"github.com/akornatskyy/precompress/compressors/zstd"
)

func Walk(opts ...Option) error {
	defaults := []Option{
		MinSize(1024),
		Exclude([]string{".gz", ".br", ".zst"}),
		BrotliCompressionLevel(brotli.DefaultCompression),
		GzipCompressionLevel(gzip.DefaultCompression),
		ZstdCompressionLevel(zstd.DefaultCompression),
	}
	opts = append(defaults, opts...)

	c := config{providers: map[string]compressors.CompressorProvider{}}
	for _, o := range opts {
		if err := o(&c); err != nil {
			return err
		}
	}

	sem := make(chan struct{}, max(runtime.NumCPU()/2, 2))
	var wg sync.WaitGroup
	for _, path := range c.paths {
		rootDepth := strings.Count(filepath.Clean(path), string(filepath.Separator))
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return fmt.Errorf("error accessing path %q: %v", path, err)
			}

			if d.IsDir() {
				if c.maxDepth != 0 &&
					(strings.Count(path, string(filepath.Separator))-rootDepth) >=
						c.maxDepth {
					return fs.SkipDir
				}
			} else {
				for _, ext := range c.exclude {
					if strings.HasSuffix(path, ext) {
						return nil
					}
				}

				fi, err := d.Info()
				if err != nil {
					return fmt.Errorf("error getting path info %q: %v", path, err)
				}
				if fi.Size() < c.minSize {
					return nil
				}

				sem <- struct{}{}
				wg.Add(1)
				go func(path string, fi fs.FileInfo) {
					defer func() {
						<-sem
						wg.Done()
					}()
					if err := process(&c, path, fi); err != nil {
						log.Printf("error processing %q: %v", path, err)
					}
				}(path, fi)
			}
			return nil
		})

		if err != nil {
			return err
		}
	}

	wg.Wait()
	return nil
}

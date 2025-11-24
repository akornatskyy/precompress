package walker

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"sync"
)

type WalkFunc func(path string, fi fs.FileInfo) error

type Walker interface {
	Walk(paths []string, fn WalkFunc) error
}

type walker struct {
	minSize  int64
	maxDepth int
	exclude  []string
}

func New(opts ...Option) (Walker, error) {
	w := &walker{}
	for _, o := range opts {
		if err := o(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}

func (w *walker) Walk(paths []string, fn WalkFunc) error {
	sem := make(chan struct{}, max(runtime.NumCPU()/2, 2))
	var wg sync.WaitGroup
	for _, path := range paths {
		rootDepth := strings.Count(filepath.Clean(path), string(filepath.Separator))
		err := filepath.WalkDir(path, func(path string, d fs.DirEntry, err error) error {
			if err != nil {
				return err
			}

			if d.IsDir() {
				if w.maxDepth != 0 &&
					(strings.Count(path, string(filepath.Separator))-rootDepth) >=
						w.maxDepth {
					return fs.SkipDir
				}
			} else {
				for _, ext := range w.exclude {
					if strings.HasSuffix(path, ext) {
						return nil
					}
				}

				fi, err := d.Info()
				if err != nil {
					return err
				}
				if fi.Size() < w.minSize {
					return nil
				}

				sem <- struct{}{}
				wg.Add(1)
				go func(path string, fi fs.FileInfo) {
					defer func() {
						<-sem
						wg.Done()
					}()
					if err := fn(path, fi); err != nil {
						// TODO: better error propagation?
						fmt.Fprintf(os.Stderr, "Error processing %q: %v\n", path, err)
						os.Exit(1)
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

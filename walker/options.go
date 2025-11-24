package walker

import (
	"fmt"
	"slices"
)

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

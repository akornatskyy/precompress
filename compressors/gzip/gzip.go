package gzip

import (
	"compress/gzip"
	"io"
	"sync"
)

const (
	DefaultCompression = gzip.DefaultCompression
)

type Options struct {
	Level int
}

type compressor struct {
	pool sync.Pool
	opt  Options
}

func New(opt Options) (*compressor, error) {
	c := &compressor{opt: opt}
	return c, nil
}

func (c *compressor) Get(w io.Writer) io.WriteCloser {
	if gw, ok := c.pool.Get().(*writer); ok {
		gw.Reset(w)
		return gw
	}
	gw, _ := gzip.NewWriterLevel(w, c.opt.Level)
	return &writer{
		Writer: gw,
		c:      c,
	}
}

type writer struct {
	*gzip.Writer
	c *compressor
}

func (w *writer) Close() error {
	err := w.Writer.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}

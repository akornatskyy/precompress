package zstd

import (
	"io"
	"sync"

	"github.com/klauspost/compress/zstd"
)

const (
	DefaultCompression = int(zstd.SpeedDefault)
)

type Options struct {
	Level int
}

type compressor struct {
	pool sync.Pool
	opts []zstd.EOption
}

func New(opts Options) (c *compressor, err error) {
	c = &compressor{
		opts: []zstd.EOption{
			zstd.WithEncoderLevel(zstd.EncoderLevel(opts.Level)),
		},
	}
	return c, nil
}

func (c *compressor) Get(w io.Writer) io.WriteCloser {
	if gw, ok := c.pool.Get().(*writer); ok {
		gw.Reset(w)
		return gw
	}
	gw, _ := zstd.NewWriter(w, c.opts...)
	return &writer{
		Encoder: gw,
		c:       c,
	}
}

type writer struct {
	*zstd.Encoder
	c *compressor
}

func (w *writer) Close() error {
	err := w.Encoder.Close()
	w.Reset(nil)
	w.c.pool.Put(w)
	return err
}

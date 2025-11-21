package zstd_test

import (
	"bytes"
	"io"
	"testing"

	zst "github.com/klauspost/compress/zstd"

	"github.com/akornatskyy/precompress/compressors"
	"github.com/akornatskyy/precompress/compressors/zstd"
)

var _ compressors.CompressorProvider = &zstd.Compressor{}

func TestGzip(t *testing.T) {
	t.Parallel()

	s := []byte("Hello, World!")

	p, err := zstd.New(zstd.Options{Level: zstd.DefaultCompression})
	if err != nil {
		t.Fatal(err)
	}
	b := &bytes.Buffer{}
	w := p.Get(b)
	if _, err = w.Write(s); err != nil {
		t.Fatal(err)
	}

	if err = w.Close(); err != nil {
		t.Fatal(err)
	}

	r, err := zst.NewReader(b)
	if err != nil {
		t.Fatal(err)
	}
	d, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(s, d) {
		t.Fatalf("reader mismatch\ngot: %q\nexp: %q", d, s)
	}
}

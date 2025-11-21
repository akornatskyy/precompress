package gzip_test

import (
	"bytes"
	"io"
	"testing"

	gz "compress/gzip"

	"github.com/akornatskyy/precompress/compressors"
	"github.com/akornatskyy/precompress/compressors/gzip"
)

var _ compressors.CompressorProvider = &gzip.Compressor{}

func TestGzip(t *testing.T) {
	t.Parallel()

	s := []byte("Hello, World!")

	p, err := gzip.New(gzip.Options{Level: gzip.DefaultCompression})
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

	r, err := gz.NewReader(b)
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

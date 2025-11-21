package brotli_test

import (
	"bytes"
	"io"
	"testing"

	br "github.com/andybalholm/brotli"

	"github.com/akornatskyy/precompress/compressors"
	"github.com/akornatskyy/precompress/compressors/brotli"
)

var _ compressors.CompressorProvider = &brotli.Compressor{}

func TestBrotli(t *testing.T) {
	t.Parallel()

	s := []byte("Hello, World!")

	p, err := brotli.New(brotli.Options{Quality: brotli.DefaultCompression})
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

	r := br.NewReader(b)
	d, err := io.ReadAll(r)
	if err != nil {
		t.Fatal(err)
	}
	if !bytes.Equal(s, d) {
		t.Fatalf("reader mismatch\ngot: %q\nexp: %q", d, s)
	}
}

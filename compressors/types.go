package compressors

import (
	"io"
)

type CompressorProvider interface {
	Get(w io.Writer) (compressor io.WriteCloser)
}

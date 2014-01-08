package readcaster

import (
	"io"
)

type ReadCaster struct {
	In io.Reader
}

// New creates a new ReadCaster that will allow multiple io.Readers to read
// from the specified source.
func New(source io.Reader) *ReadCaster {
	return &ReadCaster{In: source}
}

// NewReader creates a new io.Reader capable of reading from the source
// of the ReadCaster.
func (c *ReadCaster) NewReader() *Reader {
	return &Reader{}
}

package readcaster

import (
	"io"
)

type ReadCaster struct {
	// In represents the source io.Reader where this ReadCaster will read from.
	In io.Reader
	// readers are all the Readers that will be reading from this ReadCaster.
	readers []*chanReader
}

// New creates a new ReadCaster that will allow multiple io.Readers to read
// from the specified source.
func New(source io.Reader) *ReadCaster {
	return &ReadCaster{In: source}
}

// NewReader creates a new io.Reader capable of reading from the source
// of the ReadCaster.
func (c *ReadCaster) NewReader() io.Reader {
	sourceChan := make(chan []byte)
	r := newChanReader(sourceChan)
	c.readers = append(c.readers, r)
	return r
}

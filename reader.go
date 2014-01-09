package readcaster

import (
	"io"
	"sync"
)

type chanReader struct {
	caster *ReadCaster
	once   sync.Once
	source chan []byte
	buf    []byte
}

// NewReader creates a new Reader using the specified source channel to
// read its data from.
func newChanReader(caster *ReadCaster) *chanReader {
	return &chanReader{caster: caster, source: make(chan []byte, caster.backlogSize)}
}

// Read satisfies io.Reader and writes data from the source into the
// specified byte slice.
func (r *chanReader) Read(to []byte) (int, error) {

	// make sure we have begun reading so the channels get filled up
	r.once.Do(r.caster.beginReading)

	if len(r.buf) == 0 {
		// this will block until we get data
		//
		// @tylerb: this is OK unless the buffer gets shrunk down
		// by the stuff around line 55.  In that case, it's finished but
		// there's no way for it to know.
		r.buf = <-r.source
	}

	// are we finished?
	if len(r.buf) == 0 {
		return 0, io.EOF
	}

	// if our destination is bigger than the buffer (or the same size)
	// then we're finished with the buffer
	if len(to) >= len(r.buf) && len(r.buf) != 0 {
		// fill the destination with the entire buffer
		count := copy(to, r.buf)
		r.buf = nil
		return count, nil
	}

	// if our buffer is bigger than the destination, then just copy the
	// subset.
	if len(to) < len(r.buf) && len(r.buf) != 0 {

		// fill the destination with data from the buffer
		count := copy(to, r.buf[:len(to)])

		// shrink the buffer down since we just read some
		r.buf = r.buf[len(to):]

		if len(r.buf) == 0 {
			r.buf = nil
		}

		return count, nil
	}

	return 0, io.EOF
}

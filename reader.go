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
	return &chanReader{caster: caster, source: make(chan []byte, channelBacklogSize)}
}

// Read satisfies io.Reader and writes data from the source into the
// specified byte slice.
func (r *chanReader) Read(to []byte) (int, error) {

	r.once.Do(r.caster.read)

	if len(r.buf) == 0 || r.buf == nil {
		// this will block until we get data
		r.buf = <-r.source
	}

	if len(r.buf) == 0 || r.buf == nil {
		return 0, io.EOF
	}

	if len(to) >= len(r.buf) && len(r.buf) != 0 {
		count := copy(to, r.buf)
		r.buf = nil
		return count, nil
	}

	if len(to) < len(r.buf) && len(r.buf) != 0 {
		count := copy(to, r.buf[:len(to)])
		return count, nil
	}

	return 0, io.EOF
}

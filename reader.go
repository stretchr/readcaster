package readcaster

import (
	"io"
)

// chanReader is an io.Reader that reads from the channel receiving buffers
// of data from the caster ReadCaster.
type chanReader struct {
	// caster is the ReadCaster that this chanReader will receive
	// buffers of data from.
	caster *ReadCaster
	// source is the channel on which this reader will receive buffers
	// of data from the caster ReadCaster.
	source chan []byte
	// buf is the most recent buffer of data received on the source channel.
	buf []byte
}

// NewReader creates a new Reader using the specified source channel to
// read its data from.
func newChanReader(caster *ReadCaster) *chanReader {
	return &chanReader{caster: caster, source: make(chan []byte, caster.backlogSize)}
}

// Read satisfies io.Reader and writes data from the source into the
// specified byte slice.
//
// The number of bytes read will be returned, or an error if something
// goes wrong.  As per the io.Reader interface, Read will return an io.EOF error
// when there is no more data to come.
func (r *chanReader) Read(to []byte) (int, error) {

	// make sure we have begun reading so the channels get filled up
	r.caster.beginReading()

	if len(r.buf) == 0 {
		// this will block until we get data
		r.buf = <-r.source
	}

	// are we finished?
	if len(r.buf) == 0 {
		// we're done
		return 0, io.EOF
	}

	// if our destination is bigger than the buffer (or the same size)
	// then we're finished with the buffer
	if len(to) >= len(r.buf) && len(r.buf) != 0 {
		// fill the destination with the entire buffer
		count := copy(to, r.buf)
		r.buf = nil
		// we've read some, but there is more
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

		// we've read some, but there is more
		return count, nil
	}

	// we're done
	return 0, io.EOF
}

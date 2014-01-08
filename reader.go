package readcaster

type chanReader struct {
	source chan []byte
}

// NewReader creates a new Reader using the specified source channel to
// read its data from.
func newChanReader(source chan []byte) *chanReader {
	return &chanReader{source: source}
}

func (r *chanReader) Read(to []byte) (int, error) {
	return 0, nil
}

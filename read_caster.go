package readcaster

import (
	"io"
	"sync"
)

const (
	defaultBufferSize  int = 4096
	channelBacklogSize int = 10
)

type ReadCaster struct {
	// In represents the source io.Reader where this ReadCaster will read from.
	In io.Reader
	// readers are all the Readers that will be reading from this ReadCaster.
	readers []*chanReader
	// once is used to control the initiation of the reading process
	once sync.Once
	// bufferSize is the length of the internal buffer.
	bufferSize int
}

// New creates a new ReadCaster that will allow multiple io.Readers to read
// from the specified source.
func New(source io.Reader) *ReadCaster {
	return NewSize(source, defaultBufferSize)
}

// NewSize creates a new ReadCaster that will allow multiple io.Readers to read
// from the specified source, while also setting the size of the internal buffer.
func NewSize(source io.Reader, bufferSize int) *ReadCaster {
	return &ReadCaster{In: source, bufferSize: bufferSize}
}

// NewReader creates a new io.Reader capable of reading from the source
// of the ReadCaster.
func (c *ReadCaster) NewReader() io.Reader {
	r := newChanReader(c)
	c.readers = append(c.readers, r)
	return r
}

// BufferSize gets the size of the internal buffer that is used to hold
// the content from the source.
func (c *ReadCaster) BufferSize() int {
	return c.bufferSize
}

// ApproxMemoryUse calculates the maximum amount of memory (in bytes)
// that will be used by this ReadCaster and its Readers.
//
// It is calcualted by finding the product of the buffer size, the
// channelBacklogSize and the number of Readers created by a call to
// NewReader.
func (c *ReadCaster) ApproxMemoryUse() int {
	return c.bufferSize * channelBacklogSize * len(c.readers)
}

// read begins reading data from In and sending it to the channels
func (c *ReadCaster) read() {
	c.once.Do(func() {
		go func() {
			for {
				buf := make([]byte, c.bufferSize)
				n, err := c.In.Read(buf)
				if err != nil || n == 0 {
					for _, reader := range c.readers {
						close(reader.source)
					}
					break
				}
				for _, reader := range c.readers {
					reader.source <- buf[:n]
				}
			}
		}()
	})
}

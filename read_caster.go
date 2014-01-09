package readcaster

import (
	"io"
	"sync"
)

const (
	// defaultBufferSize is the default size (in bytes) of the buffer that is used
	// to store data read from the source, before it is read by the readers.
	defaultBufferSize int = 4096
	// defaultBacklogSize is the number of buffers that will be queued up ready for
	// the readers to read.  This allows readers to read at their own pace before
	// the reads get blocked waiting for other readers to catch up.
	defaultBacklogSize int = 10
)

// ReadCaster allows you to spawn many io.Readers using NewReader() that may each
// read, at their own pace, from the same io.Reader source.
//
// The BufferSize and and BacklogSize (set with NewSize) allow you to limit the
// amount of memory used by each reader, allowing you to keep the memory footprint
// of your application under control.
type ReadCaster struct {
	// in represents the source io.Reader where this ReadCaster will read from.
	in io.Reader
	// readers are all the Readers that will be reading from this ReadCaster.
	readers []*chanReader
	// once is used to control the initiation of the reading process
	once sync.Once
	// bufferSize is the length of the internal buffer.
	bufferSize int
	// backlogSize is the number of buffers to keep in a queue ready for the
	// readers to read.
	backlogSize int
}

// New creates a new ReadCaster that will allow multiple io.Readers to read
// from the specified source.
//
// It will use the default sizes for the buffer and backlog.
func New(source io.Reader) *ReadCaster {
	return NewSize(source, defaultBufferSize, defaultBacklogSize)
}

// NewSize creates a new ReadCaster that will allow multiple io.Readers to read
// from the specified source, while also setting the BufferSize() and BacklogSize().
func NewSize(source io.Reader, bufferSize, backlogSize int) *ReadCaster {

	if bufferSize < 1 {
		panic("readcaster: bufferSize must be greater than zero.")
	}

	return &ReadCaster{in: source, bufferSize: bufferSize, backlogSize: backlogSize}
}

// NewReader creates a new io.Reader capable of reading from the source
// of the ReadCaster.
//
// The readers returned from this method must be passed into a go routine
// in order for reading to commence to avoid the chance of deadlock.
func (c *ReadCaster) NewReader() io.Reader {
	r := newChanReader(c)
	c.readers = append(c.readers, r)
	return r
}

// BufferSize gets the size of the internal buffer that is used to hold
// the content from the source.
//
// For best performance, calls to the Read method of the readers should
// try and read the same number of bytes in this buffer.
func (c *ReadCaster) BufferSize() int {
	return c.bufferSize
}

// BacklogSize gets the number of buffers that will be queued up ready for
// the readers to read.  This allows readers to read at their own pace before
// the reads get blocked waiting for other readers to catch up.
func (c *ReadCaster) BacklogSize() int {
	return c.backlogSize
}

// ApproxMemoryUse calculates the maximum amount of memory (in bytes)
// that will be used by this ReadCaster and its Readers.
//
// It is calcualted by finding the product of the BufferSize(), the
// BacklogSize() and the number of Readers created by calling
// NewReader.
func (c *ReadCaster) ApproxMemoryUse() int {
	return c.bufferSize * c.backlogSize * len(c.readers)
}

// beginReading begins reading data from In and sending it to the channels
func (c *ReadCaster) beginReading() {
	c.once.Do(func() {
		go func() {
			for {
				buf := make([]byte, c.bufferSize)
				n, err := c.in.Read(buf)
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

package readcaster

import (
	"io"
	"sync"
	"time"
)

const (
	// defaultBufferSize is the default size (in bytes) of the buffer that is used
	// to store data read from the source, before it is read by the readers.
	defaultBufferSize int = 4096
	// defaultBacklogSize is the number of buffers that will be queued up ready for
	// the readers to read.  This allows readers to read at their own pace before
	// the reads get blocked waiting for other readers to catch up.
	defaultBacklogSize int = 10
	// defaultReaderTimeout is the default duration that the caster will wait
	// before killing a slow/dead reader.
	defaultReaderTimeout time.Duration = 1 * time.Second
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
	// startedReading is whether one or more readers have already started reading
	// from this ReadCaster.  If they have, configuration becomes locked down.
	startedReading bool
	// bufferSize is the length of the internal buffer.
	bufferSize int
	// backlogSize is the number of buffers to keep in a queue ready for the
	// readers to read.
	backlogSize int
	// readerTimeout is the duration that the caster will wait
	// before killing a slow/dead reader.
	readerTimeout time.Duration
	// Progress is the channel on which read progress is sent. Each
	// message to the channel is the total size read at that point.
	Progress chan int
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

	return &ReadCaster{
		in:            source,
		bufferSize:    bufferSize,
		backlogSize:   backlogSize,
		readerTimeout: defaultReaderTimeout,
		Progress:      make(chan int),
	}
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

// ensureCanChangeConfig panics if one or more readers have begun reading
// from this ReadCaster.
func (c *ReadCaster) ensureCanChangeConfig() {
	if c.startedReading {
		panic("readcaster: Cannot change configuration of a ReadCaster once Readers have begun reading.  Ensure all configuration occurrs before generating readers with NewReader() calls.")
	}
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

// MaxMemoryUse calculates the maximum amount of memory (in bytes)
// that will be used by this ReadCaster and its Readers.
//
// It is calcualted by finding the product of the BufferSize(), the
// BacklogSize() and the number of Readers created by calling
// NewReader.
func (c *ReadCaster) MaxMemoryUse() int {
	return c.bufferSize * c.backlogSize * len(c.readers)
}

func (c *ReadCaster) ReaderTimeout() time.Duration {
	return c.readerTimeout
}

func (c *ReadCaster) SetReaderTimeout(duration time.Duration) {
	c.ensureCanChangeConfig()
	c.readerTimeout = duration
}

// beginReading begins reading data from In and sending it to the channels
func (c *ReadCaster) beginReading() {
	c.once.Do(func() {
		c.startedReading = true
		go func() {
			totalBytesRead := 0
			for {
				// make a buffer
				buf := make([]byte, c.bufferSize)
				n, err := c.in.Read(buf)

				totalBytesRead += n

				select {
				case c.Progress <- totalBytesRead:
				default:
				}

				if err != nil || n == 0 {
					// close the channels - we're done
					for _, reader := range c.readers {
						if !reader.hasTimedOut {
							close(reader.source)
						}
					}
					close(c.Progress)
					break
				}

				// send the content from the buffer to the channels
				timeout := time.NewTimer(c.readerTimeout)
				for _, reader := range c.readers {
					if reader.hasTimedOut {
						continue
					}
					select {
					case reader.source <- buf[:n]:
						continue
					case <-timeout.C:
						reader.hasTimedOut = true
						close(reader.source)
					}

				}

			}

		}()
	})
}

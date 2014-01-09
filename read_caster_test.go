package readcaster

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
	"time"
)

func TestNew(t *testing.T) {

	source := strings.NewReader("Hello from Stretchr")
	c := New(source)

	assert.NotNil(t, c)
	assert.Equal(t, c.in, source)
	assert.Equal(t, c.bufferSize, defaultBufferSize)

}

func TestNewSize(t *testing.T) {

	source := strings.NewReader("Hello from Stretchr")
	c := NewSize(source, 25, 5)

	assert.NotNil(t, c)
	assert.Equal(t, c.in, source)
	assert.Equal(t, c.bufferSize, 25)
	assert.Equal(t, c.backlogSize, 5)
	assert.Equal(t, c.readerTimeout, defaultReaderTimeout, "readerTimeout should be set to default")

	assert.Panics(t, func() {
		NewSize(source, 0, 1)
	}, "Zero buffer size")
	assert.Panics(t, func() {
		NewSize(source, -1, 5)
	}, "Nagative buffer size")

}

func TestSizeGetters(t *testing.T) {

	source := strings.NewReader("Hello from Stretchr")
	c := NewSize(source, 25, 5)

	assert.NotNil(t, c)
	assert.Equal(t, c.in, source)
	assert.Equal(t, c.BufferSize(), 25)
	assert.Equal(t, c.BacklogSize(), 5)

}

func TestReaderTimeout(t *testing.T) {

	source := strings.NewReader("Hello from Stretchr")
	c := NewSize(source, 25, 5)

	c.SetReaderTimeout(10 * time.Second)
	assert.Equal(t, c.ReaderTimeout(), 10*time.Second)

	c.startedReading = true

	assert.Panics(t, func() {
		c.SetReaderTimeout(5 * time.Second)
	}, "Should panic when trying to set ReaderTimeout after reading has started")

}

func TestReadCasterEnsureCanChangeConfig(t *testing.T) {

	source := strings.NewReader("Hello from Stretchr")
	c := NewSize(source, 25, 5)

	c.ensureCanChangeConfig()

	c.startedReading = true

	assert.Panics(t, func() {
		c.ensureCanChangeConfig()
	}, "Should panic when calling ensureCanChangeConfig after reading has started")

}

func TestApproxMemoryUse(t *testing.T) {

	source := strings.NewReader("Hello from Stretchr")
	c := NewSize(source, 25, 5)

	c.NewReader()
	c.NewReader()
	c.NewReader()

	assert.Equal(t, c.MaxMemoryUse(), 25*5*3)

}

func TestNewReader(t *testing.T) {

	source := strings.NewReader("Hello from Stretchr")
	c := New(source)

	reader := c.NewReader().(*chanReader)
	assert.NotNil(t, reader)

	// ensure the reader was added to the readers array
	if assert.Equal(t, 1, len(c.readers)) {
		assert.Equal(t, reader, c.readers[0])
	}

}

package readcaster

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
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

func TestApproxMemoryUse(t *testing.T) {

	source := strings.NewReader("Hello from Stretchr")
	c := NewSize(source, 25, 5)

	c.NewReader()
	c.NewReader()
	c.NewReader()

	assert.Equal(t, c.ApproxMemoryUse(), 25*5*3)

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

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
	assert.Equal(t, c.In, source)

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

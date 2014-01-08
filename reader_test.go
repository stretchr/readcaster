package readcaster

import (
	"github.com/stretchr/testify/assert"
	"io"
	"strings"
	"testing"
)

func TestReaderInterface(t *testing.T) {

	var ioreader io.Reader = new(chanReader)
	assert.NotNil(t, ioreader)

}

func TestReaderNewReader(t *testing.T) {

	sourceReader := strings.NewReader("Test")
	readcaster := New(sourceReader)
	reader := newChanReader(readcaster)

	if assert.NotNil(t, reader) {
		assert.NotNil(t, reader.source)
	}

}

package readcaster

import (
	"github.com/stretchr/testify/assert"
	"io"
	"testing"
)

func TestReaderInterface(t *testing.T) {

	var ioreader io.Reader = new(chanReader)
	assert.NotNil(t, ioreader)

}

func TestReaderNewReader(t *testing.T) {

	reader := newChanReader()

	if assert.NotNil(t, reader) {
		assert.NotNil(t, reader.source)
	}

}

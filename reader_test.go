package readcaster

import (
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReaderInterface(t *testing.T) {

	var ioreader = new(chanReader)
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

func TestReadWhenTimedOut(t *testing.T) {

	sourceReader := strings.NewReader("Test")
	readcaster := New(sourceReader)
	reader := newChanReader(readcaster)

	var wg sync.WaitGroup
	wg.Add(1)

	reader.hasTimedOut = true
	go func() {
		_, err := reader.Read(nil)
		if assert.Error(t, err) {
			assert.Equal(t, "", err.Error())
		}
		wg.Done()
	}()

}

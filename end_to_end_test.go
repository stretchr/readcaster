package readcaster

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"strings"
	"sync"
	"testing"
)

// TestEndToEndAllBuffers tests the concept with an ever-increasing buffer
// size, starting at 10 and counting up by 10 until it passes 1024.
func TestEndToEndAllBuffers(t *testing.T) {

	bufferSize := 0

	for {

		bufferSize = bufferSize + 10
		source := "Hello from Stretchr."
		sourceReader := strings.NewReader(source)
		caster := NewSize(sourceReader, bufferSize)

		r1 := caster.NewReader()
		r2 := caster.NewReader()

		var r1bytes []byte
		var r2bytes []byte

		// read all in all readers
		var allread sync.WaitGroup
		allread.Add(2)
		go func() {
			var err error
			r1bytes, err = ioutil.ReadAll(r1)
			assert.NoError(t, err)
			allread.Done()
		}()
		go func() {
			var err error
			r2bytes, err = ioutil.ReadAll(r2)
			assert.NoError(t, err)
			allread.Done()
		}()
		allread.Wait() // wait for all readers to finish

		// make sure all bytes are present and correct
		assert.Equal(t, source, string(r1bytes), "r1bytes")
		assert.Equal(t, source, string(r2bytes), "r2bytes")

		// test with a buffer size of up to 1024
		if bufferSize > 1024 {
			break
		}

	}
}

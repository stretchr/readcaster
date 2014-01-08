package readcaster

import (
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"strings"
	"sync"
	"testing"
)

// TestEndToEndAllBuffers tests the concept with an ever-increasing buffer
// size, starting at 10 and counting up by 10 until it passes 1024.
func TestEndToEndAllBuffers(t *testing.T) {

	bufferSize := 0
	backlogSize := 1

	for backlogSize < 10 {

		for {

			bufferSize = bufferSize + 10
			source := "Hello from Stretchr."
			sourceReader := strings.NewReader(source)
			caster := NewSize(sourceReader, bufferSize, backlogSize)

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
		backlogSize++
	}

}

func TestEndToEndReadAByteAtATime(t *testing.T) {

	source := "Hello from Stretchr."
	sourceReader := strings.NewReader(source)
	caster := New(sourceReader)

	r1 := caster.NewReader()
	r2 := caster.NewReader()

	var r1bytes []byte
	var r2bytes []byte

	// read all in all readers
	var allread sync.WaitGroup
	allread.Add(2)
	go func() {
		var err error
		var n int
		var count int
		for {
			buf := make([]byte, 1)
			n, err = r1.Read(buf)
			log.Printf("Read: %s", buf)
			count += n
			if count > len(source) {
				break
			}
			if assert.NoError(t, err) {
				assert.Equal(t, n, 1)
				r1bytes = append(r1bytes, buf[:]...)
			}
		}
		allread.Done()
	}()

	go func() {
		var err error
		var n int
		var count int
		for {
			buf := make([]byte, 1)
			n, err = r2.Read(buf)
			log.Printf("Read: %s", buf)
			count += n
			if count > len(source) {
				break
			}
			if assert.NoError(t, err) {
				assert.Equal(t, n, 1)
				r2bytes = append(r2bytes, buf[:]...)
			}
		}
		allread.Done()
	}()

	allread.Wait() // wait for all readers to finish

	// make sure all bytes are present and correct
	assert.Equal(t, source, string(r1bytes), "Reading a byte at a time should still work.")
	assert.Equal(t, source, string(r2bytes), "Reading a byte at a time should still work.")

}

package watcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

const WatchPath = "./test-watcher"

var origBytes = []byte("orig")
var updatedBytes = []byte("updated")

func writeFile(bts []byte) error {
	return ioutil.WriteFile(WatchPath, bts, 0664)
}

func TestNewWatcher(t *testing.T) {
	err := writeFile(origBytes)
	assert.NoError(t, err)

	w := New(WatchPath)

	err = w.Init()
	assert.NoError(t, err)

	doneChan := make(chan bool)
	once := sync.Once
	closeChan := func() {
		once.Do(func() { close(doneChan) })
	}

	w.Register(func(bts []byte) {
		// check updated bytes on write
		assert.Equal(t, fmt.Sprintf("%s", updatedBytes), fmt.Sprintf("%s", bts))
		closeChan()
	})

	go w.Watch()

	// let the file watcher catch up
	time.Sleep(10 * time.Millisecond)

	err = writeFile(updatedBytes)
	assert.NoError(t, err)

	<-doneChan

	err = os.Remove(WatchPath)
	assert.NoError(t, err)
}

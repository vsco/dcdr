package watcher

import (
	"testing"

	"io/ioutil"
	"sync"

	"time"

	"fmt"

	"os"

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

	w := NewWatcher(WatchPath)

	err = w.Init()
	assert.NoError(t, err)

	var wg sync.WaitGroup
	wg.Add(2)

	updateCount := 0

	w.Register(func(bts []byte) {
		if updateCount == 0 {
			// check original bytes when watch is started
			assert.Equal(t, fmt.Sprintf("%s", origBytes), fmt.Sprintf("%s", bts))
		} else {
			// check updated bytes on write
			assert.Equal(t, fmt.Sprintf("%s", updatedBytes), fmt.Sprintf("%s", bts))
		}

		updateCount++
		wg.Done()
	})

	go w.Watch()

	// let the file watcher catch up
	time.Sleep(10 * time.Millisecond)

	err = writeFile(updatedBytes)
	assert.NoError(t, err)

	wg.Wait()

	err = os.Remove(WatchPath)
	assert.NoError(t, err)
}

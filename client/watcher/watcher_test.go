package watcher

import (
	"fmt"
	"io/ioutil"
	"os"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/cli/api/ioutil2"
)

const WatchPath = "./test-watcher"
const AtomicWatchPath = "./test-watcher"

var origBytes = []byte("orig")
var updatedBytes = []byte("updated")

func writeFile(bts []byte) error {
	return ioutil.WriteFile(WatchPath, bts, 0664)
}

func writeFileAtomic(bts []byte) error {
	return ioutil2.WriteFileAtomic(AtomicWatchPath, bts, 0664)
}

func TestNewWatcher(t *testing.T) {
	err := writeFile(origBytes)
	assert.NoError(t, err)

	w := New(WatchPath)

	err = w.Init()
	assert.NoError(t, err)

	doneChan := make(chan bool)
	var once sync.Once
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

func TestNewWatcherAtomicWrites(t *testing.T) {
	err := writeFileAtomic(origBytes)
	assert.NoError(t, err)

	w := New(AtomicWatchPath)

	err = w.Init()
	assert.NoError(t, err)

	doneChan := make(chan bool)
	var once sync.Once
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

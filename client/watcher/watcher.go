package watcher

import (
	"errors"
	"os"

	"io/ioutil"
	"sync"

	"github.com/fsnotify/fsnotify"
	"github.com/vsco/dcdr/cli/printer"
)

// IFace interface for the the file system watcher.
type IFace interface {
	Init() error
	Watch()
	Register(func(bts []byte))
	UpdateBytes() error
	ReadFile() ([]byte, error)
}

// Watcher is a wrapper for `fsnotify` that provides the
// registration of a callback for WRITE events.
type Watcher struct {
	path          string
	writeCallback func(bts []byte)
	watcher       *fsnotify.Watcher
	mu            sync.Mutex
}

// New initializes a Watcher and verifies that `path` exists.
func New(path string) (w *Watcher) {
	_, err := os.Stat(path)

	if err != nil {
		printer.LogErrf("could not start watcher: %v", err)
		return nil
	}

	printer.Logf("watching path`: %s", path)

	w = &Watcher{
		path: path,
	}

	return
}

// Init creates a new `fsnotify` watcher observing `path`.
func (w *Watcher) Init() error {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}

	err = watcher.Add(w.path)

	if err != nil {
		return err
	}

	w.mu.Lock()
	defer w.mu.Unlock()
	w.watcher = watcher

	return nil
}

func (w *Watcher) Watch() {
	done := make(chan bool)
	go func() {
		for {
			w.mu.Lock()
			select {
			case event := <-w.watcher.Events:
				printer.Logf("event log: %s", event.String())
				if event.Op&fsnotify.Write == fsnotify.Write ||
					event.Op&fsnotify.Create == fsnotify.Create ||
					event.Op&fsnotify.Chmod == fsnotify.Chmod {
					err := w.UpdateBytes()
					if err != nil {
						printer.LogErrf("UpdateBytes error: %v", err)
					}

					// Rewatch the path
					err = w.watcher.Remove(w.path)
					if err != nil {
						printer.LogErrf("fsnotify Remove error: %v", err)
					}
					err = w.watcher.Add(w.path)
					if err != nil {
						printer.LogErrf("fsnotify Add error: %v", err)
					}
				}
			case err := <-w.watcher.Errors:
				printer.LogErrf("watch error: %v", err)
			}
			w.mu.Unlock()
		}
	}()

	defer w.Close()

	<-done
}

func (w *Watcher) Close() {
	w.watcher.Close()
}

// Register assigns the WRITE event callback.
func (w *Watcher) Register(cb func(bts []byte)) {
	w.writeCallback = cb
}

// UpdateBytes reads the contents of `path` and passes
// the bytes to `writeCallback`.
func (w *Watcher) UpdateBytes() error {
	bts, err := w.ReadFile()

	if err != nil {
		return err
	}

	if len(bts) == 0 {
		return errors.New("Empty file read.")
	}

	w.writeCallback(bts)

	return nil
}

// ReadFile reads the contents of `path`.
func (w *Watcher) ReadFile() ([]byte, error) {
	bts, err := ioutil.ReadFile(w.path)

	return bts, err
}

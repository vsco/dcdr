package watcher

import (
	"errors"
	"os"

	"io/ioutil"

	"github.com/vsco/dcdr/cli/printer"
	"gopkg.in/fsnotify.v1"
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
}

// New initializes a Watcher and verifies that `path` exists.
func New(path string) (w *Watcher) {
	_, err := os.Stat(path)

	if err != nil {
		printer.LogErrf("could not start watcher: %v", err)
		return nil
	}

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

	w.watcher = watcher

	return nil
}

// Watch observes WRITE events, forwarding them to `Updated`
func (w *Watcher) Watch() {
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write || event.Op&fsnotify.Create == fsnotify.Create {
					err := w.UpdateBytes()
					if err != nil {
						printer.LogErrf("[dcdr] watch error: %v", err)
					}
				}
			case err := <-w.watcher.Errors:
				printer.LogErrf("[dcdr] watch error: %v", err)
			}
		}
	}()

	defer w.watcher.Close()

	<-done
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

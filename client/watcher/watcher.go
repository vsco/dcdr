package watcher

import (
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
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
	watchedPath   string
	writeCallback func(bts []byte)
	watcher       *fsnotify.Watcher
	mu            sync.Mutex
}

// New initializes a Watcher and verifies that `filepath` exists.
func New(filepath string) (w *Watcher) {
	filepath = path.Clean(filepath)

	_, err := os.Stat(filepath)
	if err != nil {
		printer.LogErrf("could not start watcher: %v", err)
		return nil
	}

	printer.Logf("watching path`: %s", filepath)

	w = &Watcher{
		path: filepath,
	}

	return
}

// Init creates a new `fsnotify` watcher observing `path`.
func (w *Watcher) Init() error {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return err
	}
	//watch the parent dir
	if err = watcher.Add(path.Dir(w.path)); err != nil {
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
		defer close(done)

		for {
			w.mu.Lock()
			select {
			case event, ok := <-w.watcher.Events:
				if !ok {
					w.mu.Unlock()
					return
				}

				correctFile := event.Name != "" && strings.Contains(w.path, path.Clean(event.Name))
				isWriteEvent := (event.Op&fsnotify.Write == fsnotify.Write) || (event.Op&fsnotify.Create == fsnotify.Create) (event.Op&fsnotify.Rename == fsnotify.Rename)

				printer.Say("received fsnotify event: %v %v. Path: %v Correct file: &v, Write: %v", event.Op, event.Name, w.path, correctFile, isWriteEvent)

				if correctFile && isWriteEvent {
					printer.Say("handling fsnotify event: %v %v", event.Op, event.Name)
					err := w.UpdateBytes()
					if err != nil {
						printer.Err("UpdateBytes error: %v", err)
					}

					// Rewatch the path
					err = w.watcher.Add(path.Dir(w.path))
					if err != nil {
						printer.Err("fsnotify Add error: %v", err)
					}
				}
			case err, ok := <-w.watcher.Errors:
				if ok {
					printer.Err("watch error: %v", err)
				}
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

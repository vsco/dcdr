package watcher

import (
	"os"

	"io/ioutil"

	"log"

	"gopkg.in/fsnotify.v1"
)

type WatcherIFace interface {
	Init() error
	Watch()
	Register(func(bts []byte))
	ReadFile() ([]byte, error)
	Updated() error
}

type Watcher struct {
	path    string
	cb      func(bts []byte)
	watcher *fsnotify.Watcher
}

func NewWatcher(path string) (w *Watcher) {
	_, err := os.Stat(path)

	if err != nil {
		panic(err)
	}

	w = &Watcher{
		path: path,
	}

	return
}

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

func (w *Watcher) Watch() {
	done := make(chan bool)
	go func() {
		for {
			select {
			case event := <-w.watcher.Events:
				if event.Op&fsnotify.Write == fsnotify.Write {
					w.Updated()
				}
			case err := <-w.watcher.Errors:
				log.Printf("[dcdr] watch error: %v", err)
			}
		}
	}()

	w.Updated()
	defer w.watcher.Close()

	<-done
}

func (w *Watcher) Register(cb func(bts []byte)) {
	w.cb = cb
}

func (w *Watcher) ReadFile() ([]byte, error) {
	bts, err := ioutil.ReadFile(w.path)

	return bts, err
}

func (w *Watcher) Updated() error {
	bts, err := w.ReadFile()

	if err != nil {
		return err
	}

	w.cb(bts)

	return nil
}

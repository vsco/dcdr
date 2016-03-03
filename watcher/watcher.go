package watcher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"sync"

	"github.com/vsco/decider-api-rf/config"
	"github.com/zenazn/goji/web"
	"gopkg.in/fsnotify.v1"
)

const (
	// CurrentSHAKey current sha key
	CurrentSHAKey = "current_sha"
)

type Features map[string]interface{}

type FeatureMap struct {
	Decider Features `json:"decider"`
}

func GetFeatureMap(c *web.C) FeatureMap {
	if fm, ok := c.Env[config.DeciderEnvKey]; ok {
		return fm.(FeatureMap)
	}

	return *NewFeatureMap()
}

func NewFeatureMap() (fm *FeatureMap) {
	fm = &FeatureMap{
		Decider: Features{},
	}

	return
}

type Watcher struct {
	sync.RWMutex
	watcher      *fsnotify.Watcher
	configPath   string
	FeatureBytes []byte
	CurrentSHA   string
}

func NewWatcher(cfg *Config) (l *Watcher, err error) {
	watcher, err := fsnotify.NewWatcher()

	if err != nil {
		return
	}

	l = &Watcher{
		configPath: cfg.ConfigPath,
		watcher:    watcher,
	}

	err = l.loadFeatures()

	if err != nil {
		return nil, err
	}

	return
}

func InitEnv(c *web.C) {
	if c.Env == nil {
		c.Env = map[interface{}]interface{}{}
	}

	if c.Env[config.DeciderEnvKey] == nil {
		c.Env[config.DeciderEnvKey] = FeatureMap{
			Decider: Features{},
		}
	}
}

func AppendFeature(k string, v interface{}, c *web.C) {
	InitEnv(c)

	GetFeatureMap(c).Decider[k] = v
}

func (l *Watcher) MergeFeatures(fm *FeatureMap) (*FeatureMap, error) {
	fts, err := l.Features()

	if err != nil {
		return nil, err
	}

	merged := NewFeatureMap()

	for k, v := range fts.Decider {
		merged.Decider[k] = v
	}

	for k, v := range fm.Decider {
		merged.Decider[k] = v
	}

	return merged, nil
}

func (w *Watcher) Features() (*FeatureMap, error) {
	var fm FeatureMap

	err := json.Unmarshal(w.FeatureBytes, &fm)

	if err != nil {
		return nil, err
	}

	return &fm, nil
}

func (l *Watcher) setFeatureBytes(b []byte) {
	l.FeatureBytes = b
}

func (l *Watcher) loadFeatures() error {
	b, err := ioutil.ReadFile(l.configPath)

	if err != nil {
		return err
	}

	l.setFeatureBytes(b)

	return nil
}

func (l *Watcher) Expired(etag string) bool {
	if etag == l.CurrentSHA {
		return false
	}

	return true
}

func (l *Watcher) removeWatcher() {
	l.watcher.Remove(l.configPath)
}

func (l *Watcher) watchHandler(watcher *fsnotify.Watcher) {
	for {
		select {
		case event := <-watcher.Events:
			log.Println(event)
			if event.Op&fsnotify.Write == fsnotify.Write ||
				event.Op&fsnotify.Chmod == fsnotify.Chmod {
				log.Println("[DECIDER] reloading features")
				l.loadFeatures()
			} else if event.Op&fsnotify.Remove == fsnotify.Remove ||
				event.Op&fsnotify.Rename == fsnotify.Rename {
				log.Println("[DECIDER] removing watcher")
				l.removeWatcher()
				l.WatchConfig()
			}
		case err := <-watcher.Errors:
			log.Fatal("[DECIDER] watch error:", err)
		}
	}
}

func (l *Watcher) WatchConfig() error {
	log.Printf("Decider started watching: %s", l.configPath)

	done := make(chan bool)
	defer l.watcher.Close()
	go l.watchHandler(l.watcher)

	err := l.watcher.Add(l.configPath)

	if err != nil {
		return err
	}

	<-done

	return nil
}

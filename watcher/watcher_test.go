package watcher

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

var (
	features = &FeatureMap{
		Decider: Features{},
	}

	featureBytes, _ = json.Marshal(features)
	config_path, _  = filepath.Abs("../config/decider_fixtures.json")
)

func writeFeatures(features *FeatureMap) {
	b, _ := json.Marshal(features)

	err := ioutil.WriteFile(config_path, b, 0644)
	if err != nil {
		log.Println(err)
	}
}

func TestSetFeatures(t *testing.T) {
	w, _ := NewWatcher(TestConfig())
	w.setFeatureBytes(featureBytes)

	fts, err := w.Features()

	assert.NoError(t, err)
	assert.EqualValues(t, fts, features)
}

func TestMergeFeatures(t *testing.T) {
	w, err := NewWatcher(TestConfig())
	assert.NoError(t, err)

	fm := NewFeatureMap()
	fm.Decider["new_key"] = true

	merged, err := w.MergeFeatures(fm)
	assert.NoError(t, err)

	assert.True(t, true, merged.Decider["bool"].(bool))
	assert.True(t, true, merged.Decider["new_key"].(bool))
}

func TestWatchFeatures(t *testing.T) {
	writeFeatures(features)
	w, _ := NewWatcher(TestConfig())
	w.loadFeatures()

	go w.WatchConfig()

	fts, err := w.Features()

	assert.NoError(t, err)
	assert.EqualValues(t, fts, features)

	ticker := time.NewTicker(100 * time.Millisecond)
	quit := make(chan struct{})

	features.Decider["bool"] = false
	features.Decider["float"] = 0.5

	reload := 0
	for {
		select {
		case <-ticker.C:
			if reload%2 != 0 {
				fts, err := w.Features()

				assert.NoError(t, err)
				assert.False(t, fts.Decider["bool"].(bool))
				assert.Equal(t, 0.5, fts.Decider["float"].(float64))
				close(quit)
			} else {
				writeFeatures(features)
			}
			reload++
		case <-quit:
			ticker.Stop()
			return
		}
	}
}

package models

import (
	"encoding/json"
	"path/filepath"
	"testing"

	"io/ioutil"

	"github.com/stretchr/testify/assert"
)

var FixturePath, _ = filepath.Abs("../config/decider_fixtures.json")

func FixtureBytes() []byte {
	bts, err := ioutil.ReadFile(FixturePath)

	if err != nil {
		panic(err)
	}

	return bts
}

func FixtureMap() *FeatureMap {
	fm, err := NewFeatureMap(FixtureBytes())

	if err != nil {
		panic(err)
	}

	return fm
}

func TestUnmarshalling(t *testing.T) {
	fm := FixtureMap()
	err := json.Unmarshal(FixtureBytes(), &fm)
	assert.NoError(t, err)

	assert.Equal(t, "abcde", fm.Dcdr.Info.CurrentSHA)
}

func TestInScope(t *testing.T) {
	fm := FixtureMap()

	assert.Equal(t, float64(0), fm.Dcdr.InScope("default")["float"].(float64), "default: float")
	assert.False(t, fm.Dcdr.InScope("default")["bool"].(bool), "default: bool")
	assert.True(t, fm.Dcdr.InScope("cc/cn")["bool"].(bool), "cc/cn: bool")
	assert.Equal(t, float64(0.5), fm.Dcdr.InScope("cc/cn")["float"].(float64), "cc/cn: float")
	assert.Equal(t, nil, fm.Dcdr.InScope("cc/cn/foo/bar")["bool"], "cc/cn/foo/bar")
}

func TestMergedScopes(t *testing.T) {
	json := []byte(`{
  "dcdr": {
    "info": {
      "current_sha": "abcde"
    },
    "features": {
      "ab": {
        "float": 0.2
      },
      "cc": {
        "cn": {
          "float": 0.1
        }
      },
      "default": {
        "float": 0,
        "default_bool": true
      }
    }
  }
}`)
	fm, err := NewFeatureMap(json)
	assert.NoError(t, err)

	// Test Defaults
	assert.Equal(t, float64(0), fm.Dcdr.Defaults()["float"])
	assert.Equal(t, true, fm.Dcdr.Defaults()["default_bool"])

	// Test Nested Scopes
	assert.Equal(t, true, fm.Dcdr.MergedScopes("cc/cn")["default_bool"])
	assert.Equal(t, 0.1, fm.Dcdr.MergedScopes("cc/cn")["float"])

	// Test Scope Override Order
	assert.Equal(t, 0.2, fm.Dcdr.MergedScopes("ab", "cc/cn")["float"])
	assert.Equal(t, true, fm.Dcdr.MergedScopes("ab", "cc/cn")["default_bool"])
}

func TestToJson(t *testing.T) {
	fm := FixtureMap()

	bts, err := fm.ToJSON()
	assert.NoError(t, err)

	assert.Equal(t, FixtureBytes(), bts)
}

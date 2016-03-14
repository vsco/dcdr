package client

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/client/models"
	"github.com/vsco/dcdr/config"
)

func NewTestClient() (c *Client) {
	c = New(&config.Config{})

	return
}

var JsonBytes = []byte(`{
  "dcdr": {
    "features": {
      "ab": {
        "float": 0.5,
        "bool": false
      },
      "cc": {
        "cn": {
          "float": 1,
          "bool": true
        }
      },
      "default": {
        "float": 0,
        "bool_false": false,
        "bool": true,
        "default_float": 0.5
      }
    },
    "info": {
      "current_sha": "abcde"
    }
  }
}`)

var UpdatedJsonBytes = []byte(`{
  "dcdr": {
    "features": {
      "ab": {
      	"new_ab_feature":0.5,
        "float": 0.0,
        "bool": false
      },
      "cc": {
        "cn": {
          "float": 1,
          "bool": true
        }
      },
      "default": {
      	"new_feature": 0.5,
        "float": 1,
        "bool_false": false,
        "bool": false,
        "default_float": 1
      }
    },
    "info": {
      "current_sha": "abcde"
    }
  }
}`)

func MockFeatureMap() *models.FeatureMap {
	var fm models.FeatureMap

	err := json.Unmarshal(JsonBytes, &fm)

	if err != nil {
		panic(err)
	}

	return &fm
}

func TestNewClient(t *testing.T) {
	NewTestClient()
}

func TestWithScope(t *testing.T) {
	c := NewTestClient().WithScopes("cc/cn", "ab")

	assert.Equal(t, []string{"cc/cn", "ab"}, c.Scopes())
}

func TestSetFeatureMap(t *testing.T) {
	m := MockFeatureMap()
	c := NewTestClient().SetFeatureMap(m)

	assert.Equal(t, m.Dcdr.Defaults(), c.Features())
}

func TestEmptyFeatureMap(t *testing.T) {
	c, err := NewTestClient().Watch()
	assert.NoError(t, err)

	// ensure nil pointer guards
	c.WithScopes("scope").ScopedMap().ToJson()
}

func TestScopedFeaturesCreateNewInstance(t *testing.T) {
	scopes := []string{"ab", "cc/cn"}
	m := MockFeatureMap()
	c := NewTestClient().SetFeatureMap(m)
	c2 := c.WithScopes(scopes...)

	assert.Equal(t, m.Dcdr.Defaults(), c.Features())
	assert.Equal(t, m.Dcdr.MergedScopes(scopes...), c2.Features())
}

func TestScopedMap(t *testing.T) {
	scopes := []string{"ab", "cc/cn"}
	m := MockFeatureMap()
	c := NewTestClient().SetFeatureMap(m)
	c2 := c.WithScopes(scopes...)

	assert.False(t, c2.ScopedMap().Dcdr.Features["bool"].(bool))
	assert.Equal(t, 0.5, c2.ScopedMap().Dcdr.Features["float"])
}

func TestFeatureExists(t *testing.T) {
	m := MockFeatureMap()
	c := NewTestClient().SetFeatureMap(m)

	assert.True(t, c.FeatureExists("bool_false"))
	assert.False(t, c.FeatureExists("nope"))
}

func TestIsAvailable(t *testing.T) {
	m := MockFeatureMap()
	c := NewTestClient().SetFeatureMap(m)

	assert.True(t, c.IsAvailable("bool"))
	assert.False(t, c.IsAvailable("bool_false"))
	assert.False(t, c.IsAvailable("nope"))
}

func TestIsAvailableScoped(t *testing.T) {
	m := MockFeatureMap()
	c := NewTestClient().SetFeatureMap(m).WithScopes("ab")

	assert.False(t, c.IsAvailable("bool"))
}

func TestIsAvailableForId(t *testing.T) {
	m := MockFeatureMap()
	c := NewTestClient().SetFeatureMap(m)

	assert.False(t, c.IsAvailableForId("float", 1))
	assert.False(t, c.IsAvailableForId("float", 100))
	assert.False(t, c.IsAvailableForId("bool", 100))

	assert.True(t, c.IsAvailableForId("default_float", 10))
	assert.True(t, c.IsAvailableForId("default_float", 5))
}

func TestScaleValue(t *testing.T) {
	m := MockFeatureMap()
	c := NewTestClient().SetFeatureMap(m)

	assert.Equal(t, float64(5), c.ScaleValue("default_float", 0, 10))
	assert.Equal(t, float64(7.5), c.ScaleValue("default_float", 5, 10))
}

// ruby -e "require 'zlib';puts Zlib::crc32('some_feature123');"
// => 1706325722
// php -r "echo crc32('some_feature123');"
// => 1706325722
func TestCrc32(t *testing.T) {
	c := NewTestClient()
	id := int(c.crc(123, "some_feature"))
	expected := 1706325722

	assert.Equal(t, expected, id)
}

func TestUpdateFeatures(t *testing.T) {
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

	update := []byte(`{
	  "dcdr": {
		"info": {
		  "current_sha": "updated"
		},
		"features": {
		  "ab": {
			"float": 0.3,
			"new_ab_feature": true
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

	cfg := config.DefaultConfig()
	cfg.Git.RepoPath = "/tmp"
	c := New(cfg)
	c.UpdateFeatures(json)

	scoped := c.WithScopes("ab")

	assert.False(t, scoped.FeatureExists("new_ab_feature"))
	assert.Equal(t, float64(0.2), scoped.Features()["float"])
	assert.Equal(t, true, scoped.Features()["default_bool"])

	scoped.UpdateFeatures(update)

	assert.Equal(t, float64(0.3), scoped.Features()["float"])
	assert.True(t, scoped.FeatureExists("new_ab_feature"))
	assert.Equal(t, true, scoped.Features()["default_bool"])
}

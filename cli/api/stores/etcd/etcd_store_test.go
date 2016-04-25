package etcd

import (
	"encoding/json"
	"testing"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/client"
	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

type MockAPI struct {
	Response *client.Response
	Error    error
}

func (m *MockAPI) Get(ctx context.Context, key string, opts *client.GetOptions) (*client.Response, error) {
	return m.Response, m.Error
}

func (m *MockAPI) Set(ctx context.Context, key, value string, opts *client.SetOptions) (*client.Response, error) {
	return m.Response, m.Error
}

func (m *MockAPI) Delete(ctx context.Context, key string, opts *client.DeleteOptions) (*client.Response, error) {
	return m.Response, m.Error
}

func (m *MockAPI) Create(ctx context.Context, key, value string) (*client.Response, error) {
	return m.Response, m.Error
}

func (m *MockAPI) CreateInOrder(ctx context.Context, dir, value string, opts *client.CreateInOrderOptions) (*client.Response, error) {
	return m.Response, m.Error
}

func (m *MockAPI) Update(ctx context.Context, key, value string) (*client.Response, error) {
	return m.Response, m.Error
}

func (m *MockAPI) Watcher(key string, opts *client.WatcherOptions) client.Watcher {

	return nil
}

func TestGet(t *testing.T) {
	r := &client.Response{
		Node: &client.Node{
			Key:   "k",
			Value: "v",
		},
	}
	e := &Store{
		kv: &MockAPI{
			Response: r,
		},
		config: config.TestConfig(),
	}

	kv, err := e.Get("k")
	assert.NoError(t, err)
	assert.Equal(t, "k", kv.Key)
	assert.Equal(t, []byte("v"), kv.Bytes)
}

func TestSet(t *testing.T) {
	r := &client.Response{}
	e := &Store{
		kv: &MockAPI{
			Response: r,
		},
		config: config.TestConfig(),
	}

	err := e.Set("k", []byte("v"))
	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	r := &client.Response{
		Node: &client.Node{
			Key:   "k",
			Value: "v",
		},
	}
	e := &Store{
		kv: &MockAPI{
			Response: r,
		},
		config: config.TestConfig(),
	}

	kv, err := e.List("k")
	assert.NoError(t, err)
	assert.Equal(t, "k", kv[0].Key)
	assert.Equal(t, []byte("v"), kv[0].Bytes)
}

var jsonBytes = []byte(`{
  "action": "get",
  "node": {
    "key": "/dcdr",
    "dir": true,
    "nodes": [
      {
        "key": "/dcdr/features",
        "dir": true,
        "nodes": [
          {
            "key": "/dcdr/features/default",
            "dir": true,
            "nodes": [
              {
                "key": "/dcdr/features/default/test",
                "value": "{\"feature_type\":\"boolean\",\"key\":\"test\",\"namespace\":\"dcdr\",\"scope\":\"default\",\"value\":false,\"comment\":\"\",\"updated_by\":\"chrisb\"}",
                "modifiedIndex": 4,
                "createdIndex": 4
              }
            ],
            "modifiedIndex": 4,
            "createdIndex": 4
          }
        ],
        "modifiedIndex": 4,
        "createdIndex": 4
      }
    ],
    "modifiedIndex": 4,
    "createdIndex": 4
  }
}`)

func TestUpdated(t *testing.T) {
	s := New(config.TestConfig())

	s.Register(func(kvb stores.KVBytes) {
		assert.Equal(t, 1, len(kvb))

		var f models.Feature
		err := json.Unmarshal(kvb[0].Bytes, &f)

		assert.NoError(t, err)
		assert.Equal(t, "test", f.Key)
		assert.Equal(t, "dcdr", f.Namespace)
		assert.Equal(t, "default", f.Scope)
		assert.Equal(t, false, f.Value.(bool))
		assert.Equal(t, "dcdr/features/default/test", kvb[0].Key)
	})

	var resp client.Response
	err := json.Unmarshal(jsonBytes, &resp)
	assert.NoError(t, err)
	s.Updated(resp.Node)
}

package etcd

import (
	"encoding/json"
	"testing"

	"github.com/coreos/etcd/client"
	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

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
	kvw := New(config.DefaultConfig())

	kvw.Register(func(kvb stores.KVBytes) {
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
	kvw.Updated(resp.Node)
}

package consul

import (
	"encoding/json"
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

var jsonBytes = []byte(`[
  {
    "CreateIndex": 4398,
    "ModifyIndex": 4398,
    "LockIndex": 0,
    "Key": "dcdr/features/cc/b",
    "Flags": 0,
    "Value": "eyJmZWF0dXJlX3R5cGUiOiJwZXJjZW50aWxlIiwia2V5IjoiYiIsIm5hbWVzcGFjZSI6ImRjZHIvZmVhdHVyZXMiLCJzY29wZSI6ImNjIiwidmFsdWUiOjEsImNvbW1lbnQiOiIiLCJ1cGRhdGVkX2J5IjoiY2hyaXNiIn0="
  }]`)

func TestUpdated(t *testing.T) {
	kvw := New(config.DefaultConfig())

	kvw.Register(func(kvb stores.KVBytes) {
		assert.Equal(t, 1, len(kvb))

		var f models.Feature
		err := json.Unmarshal(kvb[0].Bytes, &f)

		assert.NoError(t, err)
		assert.Equal(t, "b", f.Key)
		assert.Equal(t, "dcdr/features", f.Namespace)
		assert.Equal(t, "cc", f.Scope)
		assert.Equal(t, 1.0, f.Value.(float64))
		assert.Equal(t, "dcdr/features/cc/b", kvb[0].Key)
	})

	var kvp api.KVPairs
	err := json.Unmarshal(jsonBytes, &kvp)
	assert.NoError(t, err)
	kvw.Updated(kvp)
}

package models

import (
	"encoding/json"
	"testing"

	"fmt"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/kv/stores"
)

var ExpectedJSON = `{
  "dcdr": {
    "features": {
      "cc": {
        "b": 1
      },
      "default": {
        "a": 1
      }
    },
    "info": {
      "current_sha": "43d4b9e7de8ed48a954f3594e6bd33e6d69b5516"
    }
  }
}`

var ExportJson = []byte(`[
  {
    "CreateIndex": 4398,
    "ModifyIndex": 4398,
    "LockIndex": 0,
    "Key": "dcdr/features/cc/b",
    "Flags": 0,
    "Value": "eyJmZWF0dXJlX3R5cGUiOiJwZXJjZW50aWxlIiwia2V5IjoiYiIsIm5hbWVzcGFjZSI6ImRjZHIvZmVhdHVyZXMiLCJzY29wZSI6ImNjIiwidmFsdWUiOjEsImNvbW1lbnQiOiIiLCJ1cGRhdGVkX2J5IjoiY2hyaXNiIn0="
  },
  {
    "CreateIndex": 4400,
    "ModifyIndex": 4400,
    "LockIndex": 0,
    "Key": "dcdr/features/default/a",
    "Flags": 0,
    "Value": "eyJmZWF0dXJlX3R5cGUiOiJwZXJjZW50aWxlIiwia2V5IjoiYSIsIm5hbWVzcGFjZSI6ImRjZHIvZmVhdHVyZXMiLCJzY29wZSI6ImRlZmF1bHQiLCJ2YWx1ZSI6MSwiY29tbWVudCI6IiIsInVwZGF0ZWRfYnkiOiJjaHJpc2IifQ=="
  },
  {
    "CreateIndex": 4399,
    "ModifyIndex": 4401,
    "LockIndex": 0,
    "Key": "dcdr/info",
    "Flags": 0,
    "Value": "eyJjdXJyZW50X3NoYSI6IjQzZDRiOWU3ZGU4ZWQ0OGE5NTRmMzU5NGU2YmQzM2U2ZDY5YjU1MTYifQ=="
  }
]
`)

func TestGetFeatureTypeFromValue(t *testing.T) {
	percentiles := []string{"1", "1.0", "0.0", "0", "0.5"}

	for _, v := range percentiles {
		_, ft := ParseValueAndFeatureType(v)
		assert.Equal(t, Percentile, ft, v)
	}
}

func TestMarshaling(t *testing.T) {
	f := &Feature{
		Key:         "test",
		Value:       true,
		FeatureType: Boolean,
		Comment:     "testing",
	}

	ff := &Feature{}

	js, _ := json.Marshal(f)
	json.Unmarshal(js, &ff)

	assert.EqualValues(t, f, ff)
}

func TestTypes(t *testing.T) {
	pf := NewFeature("key", 0.1, "comment", "user", "scope", "n")
	assert.Equal(t, Percentile, pf.FeatureType)
	assert.Equal(t, 0.1, pf.FloatValue())

	pf = NewFeature("key", true, "comment", "user", "scope", "n")
	assert.Equal(t, Boolean, pf.FeatureType)
	assert.Equal(t, true, pf.BoolValue())
}

func TestFeaturesToKVMapToJSON(t *testing.T) {
	kvp, err := stores.KvPairsBytesToKvBytes(ExportJson)
	assert.NoError(t, err)

	fts, err := KVsToFeatureMap(kvp)
	assert.NoError(t, err)

	json, err := json.MarshalIndent(fts, "", "  ")

	assert.NoError(t, err)
	assert.Equal(t, fmt.Sprintf("%s", ExpectedJSON), fmt.Sprintf("%s", json))
}

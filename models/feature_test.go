package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

var ConsulJSON = []byte(`[{"Key":"dcdr/features/cn/test","CreateIndex":319,"ModifyIndex":319,"LockIndex":0,
"Flags":0,"Value":"eyJmZWF0dXJlX3R5cGUiOiJwZXJjZW50aWxlIiwia2V5IjoidGVzdCIsIm5hbWVzcGFjZSI6ImRjZHIvZmVhdHVyZXMiLCJzY29wZSI6ImNuIiwidmFsdWUiOjAuNSwiY29tbWVudCI6IiIsInVwZGF0ZWRfYnkiOiJjaHJpc2IifQ==",
"Session":""}]`)

var ExpectedJSON = `{
	"dcdr": {
		"features": {
			"cn": {
				"test": 0.5
			}
		}
	}
}`

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
	pf := NewFeature("key", 0.1, "comment", "user", "scope")
	assert.Equal(t, Percentile, pf.FeatureType)
	assert.Equal(t, 0.1, pf.FloatValue())

	pf = NewFeature("key", true, "comment", "user", "scope")
	assert.Equal(t, Boolean, pf.FeatureType)
	assert.Equal(t, true, pf.BoolValue())
}

func TestParseFeatures(t *testing.T) {
	fts, err := KVsToFeatures(ConsulJSON)

	assert.NoError(t, err)
	assert.Equal(t, 1, len(fts))
	assert.Equal(t, "test", fts[0].Key)
	assert.Equal(t, "cn", fts[0].Scope)
	assert.Equal(t, 0.5, fts[0].Value)
}

func TestFeaturesToKVMapToJSON(t *testing.T) {
	fts, err := KVsToFeatures(ConsulJSON)

	assert.NoError(t, err)

	bts, err := fts.ToJSON()
	assert.NoError(t, err)
	assert.Equal(t, ExpectedJSON, string(bts[:]))
}

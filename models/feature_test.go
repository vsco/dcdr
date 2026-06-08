package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetFeatureTypeFromValue(t *testing.T) {
	percentiles := []string{"1", "1.0", "0.0", "0", "0.5"}

	for _, v := range percentiles {
		_, ft := ParseValueAndFeatureType(v)
		assert.Equal(t, Percentile, ft, v)
	}

	booleans := []string{"true", "false"}

	for _, v := range booleans {
		_, ft := ParseValueAndFeatureType(v)
		assert.Equal(t, Boolean, ft, v)
	}

	strings := []string{"debug", "info", "some-string"}

	for _, v := range strings {
		val, ft := ParseValueAndFeatureType(v)
		assert.Equal(t, String, ft, v)
		assert.Equal(t, v, val, v)
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

	pf = NewFeature("key", "debug", "comment", "user", "scope", "n")
	assert.Equal(t, String, pf.FeatureType)
	assert.Equal(t, "debug", pf.StringValue())
}

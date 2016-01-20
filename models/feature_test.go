package models

import (
	"encoding/json"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMarshaling(t *testing.T) {
	f := &Feature{
		Name:        "test",
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
	pf := PercentileFeature("foo", 0.1, "testing", "me")
	assert.Equal(t, Percentile, pf.FeatureType)
	assert.Equal(t, 0.1, pf.FloatValue())

	bf := BooleanFeature("foo", false, "testing", "me")
	assert.Equal(t, Boolean, bf.FeatureType)
	assert.Equal(t, false, bf.BoolValue())
}

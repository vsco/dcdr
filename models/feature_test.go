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

func TestParseFeatureType(t *testing.T) {
	cases := map[string]FeatureType{
		"boolean":    Boolean,
		"percentile": Percentile,
		"string":     String,
	}

	for in, want := range cases {
		ft, ok := ParseFeatureType(in)
		assert.True(t, ok, in)
		assert.Equal(t, want, ft, in)
	}

	for _, in := range []string{"", "bool", "pct", "decimal", "nope"} {
		ft, ok := ParseFeatureType(in)
		assert.False(t, ok, in)
		assert.Equal(t, Invalid, ft, in)
	}
}

func TestParseValueForType(t *testing.T) {
	v, err := ParseValueForType("true", Boolean)
	assert.NoError(t, err)
	assert.Equal(t, true, v)

	_, err = ParseValueForType("notabool", Boolean)
	assert.Error(t, err)

	v, err = ParseValueForType("0.5", Percentile)
	assert.NoError(t, err)
	assert.Equal(t, 0.5, v)

	_, err = ParseValueForType("notanumber", Percentile)
	assert.Error(t, err)

	// String stores the literal value verbatim, even bool/number-looking input.
	for _, s := range []string{"debug", "true", "0.5"} {
		v, err = ParseValueForType(s, String)
		assert.NoError(t, err)
		assert.Equal(t, s, v)
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

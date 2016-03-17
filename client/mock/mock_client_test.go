package mock

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestMockClient(t *testing.T) {
	d := New()

	d.EnableBoolFeature("bool")
	assert.True(t, d.IsAvailable("bool"))
	d.DisableBoolFeature("bool")
	assert.False(t, d.IsAvailable("bool"))

	d.EnablePercentileFeature("float")
	assert.True(t, d.IsAvailableForID("float", 2))
	d.DisablePercentileFeature("float")
	assert.False(t, d.IsAvailableForID("float", 8))
}

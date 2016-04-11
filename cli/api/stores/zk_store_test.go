package stores

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestZkGet(t *testing.T) {
	mc := NewMockZk("a", MockKVBytes, nil)
	cs := NewZkStore(mc)

	bts, err := cs.Get("a")

	assert.NoError(t, err)
	assert.Equal(t, MockKVBytes[0], bts)
}

func TestZkList(t *testing.T) {
	mc := NewMockZk("n", MockKVBytes, nil)
	cs := NewZkStore(mc)

	bts, err := cs.List("n")

	assert.NoError(t, err)
	assert.Equal(t, MockKVBytes, bts)
}

func TestZkPut(t *testing.T) {
	mc := NewMockZk("n", MockKVBytes, nil)
	cs := NewZkStore(mc)

	err := cs.Put("n", MockBytes)

	assert.NoError(t, err)
}

func TestZkDelete(t *testing.T) {
	mc := NewMockZk("n", MockKVBytes, nil)
	cs := NewZkStore(mc)

	err := cs.Delete("n")

	assert.NoError(t, err)
}

package consul

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/cli/api/stores"
)

var MockBytes = []byte("asdf")
var MockKVBytes = stores.KVBytes{
	&stores.KVByte{
		Key:   "a",
		Bytes: MockBytes,
	},
}

func TestConsulGet(t *testing.T) {
	mc := NewMockConsul("a", MockKVBytes, nil)
	cs := New(mc)

	bts, err := cs.Get("a")

	assert.NoError(t, err)
	assert.Equal(t, MockKVBytes[0], bts)
}

func TestConsulList(t *testing.T) {
	mc := NewMockConsul("n", MockKVBytes, nil)
	cs := New(mc)

	bts, err := cs.List("n")

	assert.NoError(t, err)
	assert.Equal(t, MockKVBytes, bts)
}

func TestConsulPut(t *testing.T) {
	mc := NewMockConsul("n", MockKVBytes, nil)
	cs := New(mc)

	err := cs.Set("n", MockBytes)

	assert.NoError(t, err)
}

func TestConsulDelete(t *testing.T) {
	mc := NewMockConsul("n", MockKVBytes, nil)
	cs := New(mc)

	err := cs.Delete("n")

	assert.NoError(t, err)
}

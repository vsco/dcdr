package stores

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

var MockBytes = []byte("asdf")
var MockKVBytes = KVBytes{
	&KVByte{
		Key:   "a",
		Bytes: MockBytes,
	},
}

type MockConsul struct {
	Item  *KVByte
	Items KVBytes
	Err   error
}

func NewMockConsul(key string, kvb KVBytes, err error) (mc *MockConsul) {
	mc = &MockConsul{
		Err: err,
	}

	if len(kvb) != 0 {
		mc.Item = kvb[0]
		mc.Items = kvb
	}

	return
}

func (mc *MockConsul) get(key string) *api.KVPair {
	return &api.KVPair{
		Key:   key,
		Value: mc.Item.Bytes,
	}
}

func (mc *MockConsul) List(prefix string, qo *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error) {
	items := api.KVPairs{&api.KVPair{
		Key:   mc.Item.Key,
		Value: mc.Item.Bytes,
	},
	}
	return items, nil, mc.Err
}

func (mc *MockConsul) Get(key string, qo *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error) {
	return mc.get(key), nil, mc.Err
}

func (mc *MockConsul) Put(p *api.KVPair, qo *api.WriteOptions) (*api.WriteMeta, error) {
	return nil, mc.Err
}

func (mc *MockConsul) Delete(key string, w *api.WriteOptions) (*api.WriteMeta, error) {
	return nil, mc.Err
}

func TestConsulGet(t *testing.T) {
	mc := NewMockConsul("a", MockKVBytes, nil)
	cs := NewConsulStore(mc)

	bts, err := cs.Get("a")

	assert.NoError(t, err)
	assert.Equal(t, MockKVBytes[0], bts)
}

func TestConsulList(t *testing.T) {
	mc := NewMockConsul("n", MockKVBytes, nil)
	cs := NewConsulStore(mc)

	bts, err := cs.List("n")

	assert.NoError(t, err)
	assert.Equal(t, MockKVBytes, bts)
}

func TestConsulPut(t *testing.T) {
	mc := NewMockConsul("n", MockKVBytes, nil)
	cs := NewConsulStore(mc)

	err := cs.Put("n", MockBytes)

	assert.NoError(t, err)
}

func TestConsulDelete(t *testing.T) {
	mc := NewMockConsul("n", MockKVBytes, nil)
	cs := NewConsulStore(mc)

	err := cs.Delete("n")

	assert.NoError(t, err)
}

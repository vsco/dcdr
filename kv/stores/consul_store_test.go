package stores

import (
	"testing"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
)

var MockBytes = []byte("asdf")

type MockConsul struct {
	Item  []byte
	Items api.KVPairs
	Err   error
}

func NewMockConsul(key string, bts []byte, err error) (mc *MockConsul) {
	mc = &MockConsul{
		Item: bts,
		Err:  err,
	}

	if len(bts) != 0 {
		mc.Items = api.KVPairs{
			{
				Key:   key,
				Value: bts,
			},
		}
	}

	return
}

func (mc *MockConsul) get(key string) *api.KVPair {
	return &api.KVPair{
		Key:   key,
		Value: mc.Item,
	}
}

func (mc *MockConsul) List(prefix string, qo *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error) {
	return mc.Items, nil, mc.Err
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
	mc := NewMockConsul("n", MockBytes, nil)
	cs := NewConsulStore(mc)

	bts, err := cs.Get("n")

	assert.NoError(t, err)
	assert.Equal(t, MockBytes, bts)
}

func TestConsulList(t *testing.T) {
	mc := NewMockConsul("n", MockBytes, nil)
	cs := NewConsulStore(mc)

	bts, err := cs.List("n")

	assert.NoError(t, err)
	assert.Equal(t, [][]byte{MockBytes}, bts)
}

func TestConsulPut(t *testing.T) {
	mc := NewMockConsul("n", MockBytes, nil)
	cs := NewConsulStore(mc)

	err := cs.Put("n", MockBytes)

	assert.NoError(t, err)
}

func TestConsulDelete(t *testing.T) {
	mc := NewMockConsul("n", MockBytes, nil)
	cs := NewConsulStore(mc)

	err := cs.Delete("n")

	assert.NoError(t, err)
}

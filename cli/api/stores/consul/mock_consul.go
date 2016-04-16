package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/vsco/dcdr/cli/api/stores"
)

type MockConsul struct {
	Item  *stores.KVByte
	Items stores.KVBytes
	Err   error
}

func NewMockConsul(key string, kvb stores.KVBytes, err error) (mc *MockConsul) {
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

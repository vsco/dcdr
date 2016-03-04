package kv

import (
	"encoding/json"

	"github.com/hashicorp/consul/api"
	"github.com/vsco/dcdr/models"
)

type MockConsul struct {
	Item  *models.Feature
	Items api.KVPairs
	Err   error
}

func MockConsulStore(ft *models.Feature, err error) (cs *ConsulStore) {
	mc := &MockConsul{
		Item: ft,
		Err:  err,
	}

	if ft != nil {
		bts, _ := json.Marshal(mc.Item)
		mc.Items = api.KVPairs{
			{
				Key:   ft.Key,
				Value: bts,
			},
		}
	}

	cs = &ConsulStore{
		kv: mc,
	}

	return
}

func (mc *MockConsul) get(key string) *api.KVPair {
	bts, _ := json.Marshal(mc.Item)

	return &api.KVPair{
		Key:   key,
		Value: bts,
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

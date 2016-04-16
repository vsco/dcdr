package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/vsco/dcdr/cli/api/stores"
)

type ConsulKVIFace interface {
	List(prefix string, q *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error)
	Get(key string, q *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error)
	Put(p *api.KVPair, q *api.WriteOptions) (*api.WriteMeta, error)
	Delete(key string, w *api.WriteOptions) (*api.WriteMeta, error)
}

type ConsulStore struct {
	kv ConsulKVIFace
	qo *api.QueryOptions
	wo *api.WriteOptions
}

func NewDefault() (stores.StoreIFace, error) {
	client, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		return nil, err
	}

	return &ConsulStore{
		kv: client.KV(),
		qo: nil,
		wo: nil,
	}, nil
}

func New(cn ConsulKVIFace) stores.StoreIFace {
	return &ConsulStore{
		kv: cn,
		qo: nil,
		wo: nil,
	}
}

func (cs *ConsulStore) Get(key string) (*stores.KVByte, error) {
	kv, _, err := cs.kv.Get(key, cs.qo)

	k := &stores.KVByte{}

	if err != nil || kv == nil {
		return nil, err
	}

	k.Key = kv.Key
	k.Bytes = kv.Value

	return k, nil
}

func (cs *ConsulStore) Set(key string, bts []byte) error {
	p := &api.KVPair{
		Key:   key,
		Value: bts,
	}

	_, err := cs.kv.Put(p, cs.wo)

	return err
}

func (cs *ConsulStore) Delete(key string) error {
	_, err := cs.kv.Delete(key, cs.wo)

	return err
}

func (cs *ConsulStore) List(prefix string) (stores.KVBytes, error) {
	kvs, _, err := cs.kv.List(prefix, cs.qo)

	kvb := make(stores.KVBytes, len(kvs))

	if err != nil {
		return kvb, err
	}

	for i := 0; i < len(kvs); i++ {
		kvb[i] = &stores.KVByte{
			Key:   kvs[i].Key,
			Bytes: kvs[i].Value,
		}
	}

	return kvb, err
}

func KvPairsToKvBytes(kvp api.KVPairs) (stores.KVBytes, error) {
	kvb := make(stores.KVBytes, len(kvp))

	for i := 0; i < len(kvp); i++ {
		kvb[i] = &stores.KVByte{
			Key:   kvp[i].Key,
			Bytes: kvp[i].Value,
		}
	}

	return kvb, nil
}

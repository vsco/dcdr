package stores

import (
	"encoding/json"

	"github.com/hashicorp/consul/api"
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

func DefaultConsulStore() (StoreIFace, error) {
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

func NewConsulStore(cn ConsulKVIFace) StoreIFace {
	return &ConsulStore{
		kv: cn,
		qo: nil,
		wo: nil,
	}
}

func (cs *ConsulStore) Get(key string) (*KVByte, error) {
	kv, _, err := cs.kv.Get(key, cs.qo)

	k := &KVByte{}

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

func (cs *ConsulStore) List(prefix string) (KVBytes, error) {
	kvs, _, err := cs.kv.List(prefix, cs.qo)

	kvb := make(KVBytes, len(kvs))

	if err != nil {
		return kvb, err
	}

	for i := 0; i < len(kvs); i++ {
		kvb[i] = &KVByte{
			Key:   kvs[i].Key,
			Bytes: kvs[i].Value,
		}
	}

	return kvb, err
}

func (cs *ConsulStore) Put(key string, bts []byte) error {
	p := &api.KVPair{
		Key:   key,
		Value: bts,
	}

	_, err := cs.kv.Put(p, cs.wo)

	return err
}

func KvPairsToKvBytes(kvp api.KVPairs) (KVBytes, error) {
	kvb := make(KVBytes, len(kvp))

	for i := 0; i < len(kvp); i++ {
		kvb[i] = &KVByte{
			Key:   kvp[i].Key,
			Bytes: kvp[i].Value,
		}
	}

	return kvb, nil
}

func KvPairsBytesToKvBytes(bts []byte) (KVBytes, error) {
	var kvp api.KVPairs

	err := json.Unmarshal(bts, &kvp)

	if err != nil {
		return make(KVBytes, 0), nil
	}

	kvb, err := KvPairsToKvBytes(kvp)

	return kvb, nil
}

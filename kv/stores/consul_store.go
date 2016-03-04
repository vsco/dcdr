package stores

import "github.com/hashicorp/consul/api"

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

func (cs *ConsulStore) Get(key string) ([]byte, error) {
	kv, _, err := cs.kv.Get(key, cs.qo)

	if err != nil || kv == nil {
		return []byte{}, err
	}

	return kv.Value, nil
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

func (cs *ConsulStore) List(prefix string) ([][]byte, error) {
	kvs, _, err := cs.kv.List(prefix, cs.qo)

	if err != nil {
		return [][]byte{}, err
	}

	res := make([][]byte, len(kvs))

	for i := 0; i < len(kvs); i++ {
		res[i] = kvs[i].Value
	}

	return res, err
}

func (cs *ConsulStore) Put(key string, bts []byte) error {
	p := &api.KVPair{
		Key:   key,
		Value: bts,
	}

	_, err := cs.kv.Put(p, cs.wo)

	return err
}

package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/config"
)

type ConsulKVIFace interface {
	List(prefix string, q *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error)
	Get(key string, q *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error)
	Put(p *api.KVPair, q *api.WriteOptions) (*api.WriteMeta, error)
	Delete(key string, w *api.WriteOptions) (*api.WriteMeta, error)
}

type Store struct {
	kv  ConsulKVIFace
	qo  *api.QueryOptions
	wo  *api.WriteOptions
	cb  func(kvb stores.KVBytes)
	cfg *config.Config
}

func NewDefault(cfg *config.Config) (stores.IFace, error) {
	client, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		return nil, err
	}

	return &Store{
		cfg: cfg,
		kv:  client.KV(),
		qo:  nil,
		wo:  nil,
	}, nil
}

func New(cfg *config.Config, cn ConsulKVIFace) stores.IFace {
	return &Store{
		cfg: cfg,
		kv:  cn,
		qo:  nil,
		wo:  nil,
	}
}

func (cs *Store) Get(key string) (*stores.KVByte, error) {
	kv, _, err := cs.kv.Get(key, cs.qo)

	k := &stores.KVByte{}

	if err != nil || kv == nil {
		return nil, err
	}

	k.Key = kv.Key
	k.Bytes = kv.Value

	return k, nil
}

func (cs *Store) Set(key string, bts []byte) error {
	p := &api.KVPair{
		Key:   key,
		Value: bts,
	}

	_, err := cs.kv.Put(p, cs.wo)

	return err
}

func (cs *Store) Delete(key string) error {
	_, err := cs.kv.Delete(key, cs.wo)

	return err
}

func (cs *Store) List(prefix string) (stores.KVBytes, error) {
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

func (cs *Store) Register(cb func(kvb stores.KVBytes)) {
	cs.cb = cb
}

func (cs *Store) Updated(kvs interface{}) {
	kvp := kvs.(api.KVPairs)
	kvb, err := KvPairsToKvBytes(kvp)

	if err != nil {
		printer.LogErrf("%v", err)
		return
	}

	cs.cb(kvb)
}

func (cs *Store) Watch() error {
	params := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": cs.cfg.Namespace,
	}

	wp, err := watch.Parse(params)
	defer wp.Stop()

	if err != nil {
		printer.LogErrf("%v", err)
	}

	wp.Handler = func(idx uint64, data interface{}) {
		cs.Updated(data)
	}

	if err := wp.Run(""); err != nil {
		printer.LogErrf("Error querying Consul agent: %s", err)
	}

	return nil
}

func (s *Store) Close() {}

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

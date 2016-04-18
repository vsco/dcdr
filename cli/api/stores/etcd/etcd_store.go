package etcd

import (
	"log"
	"time"

	"strings"

	"github.com/coreos/etcd/client"
	"github.com/fsouza/go-dockerclient/external/golang.org/x/net/context"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/config"
)

type ETCDStore struct {
	kv         client.KeysAPI
	ctx        context.Context
	getOpts    *client.GetOptions
	setOpts    *client.SetOptions
	deleteOpts *client.DeleteOptions
}

var DefaultEndpoints = []string{"http://127.0.0.1:2379"}

func DefaultETCDlStore(cfg *config.Config) (client.KeysAPI, error) {
	endpoints := DefaultEndpoints

	if len(cfg.Etcd.Endpoints) > 0 {
		endpoints = cfg.Etcd.Endpoints
	}

	ecfg := client.Config{
		Endpoints:               endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(ecfg)

	if err != nil {
		log.Fatal(err)
	}

	return client.NewKeysAPI(c), nil
}

func New(cfg *config.Config) stores.StoreIFace {
	kv, _ := DefaultETCDlStore(cfg)

	es := &ETCDStore{
		kv:  kv,
		ctx: context.Background(),
	}

	return es
}

func (s *ETCDStore) Get(key string) (*stores.KVByte, error) {
	resp, err := s.kv.Get(s.ctx, key, s.getOpts)

	if err != nil {
		return nil, etcdError(err)
	}

	return toKVByte(resp.Node), nil
}

func (s *ETCDStore) Set(key string, bts []byte) error {
	_, err := s.kv.Set(s.ctx, key, string(bts), s.setOpts)

	return err
}

func (s *ETCDStore) Delete(key string) error {
	_, err := s.kv.Delete(s.ctx, key, s.deleteOpts)

	return etcdError(err)
}

func (s *ETCDStore) List(prefix string) (stores.KVBytes, error) {
	opts := &client.GetOptions{
		Recursive: true,
		Sort:      true,
		Quorum:    true,
	}

	resp, err := s.kv.Get(s.ctx, prefix, opts)

	if err != nil {
		return nil, etcdError(err)
	}

	kvbs := FlattenToKVBytes(resp.Node, make(stores.KVBytes, 0))

	return kvbs, nil
}

func FlattenToKVBytes(n *client.Node, nodes stores.KVBytes) stores.KVBytes {
	if n.Dir {
		for _, nd := range n.Nodes {
			nodes = FlattenToKVBytes(nd, nodes)
		}
	} else {
		nodes = append(nodes, toKVByte(n))
	}

	return nodes
}

func toKVByte(n *client.Node) *stores.KVByte {
	return &stores.KVByte{
		// remove leading slash as it adds an empty
		// hash entry when exploded to JSON.
		Key:   strings.TrimPrefix(n.Key, "/"),
		Bytes: []byte(n.Value),
	}
}

func etcdError(err error) error {
	switch err.(type) {
	case client.Error:
		if err.(client.Error).Code == client.ErrorCodeKeyNotFound {
			return nil
		}

		return err
	default:
		return err
	}
}

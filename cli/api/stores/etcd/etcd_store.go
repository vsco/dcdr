package etcd

import (
	"log"
	"time"

	"strings"

	"github.com/coreos/etcd/client"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/config"
	"golang.org/x/net/context"
)

type Store struct {
	kv         client.KeysAPI
	ctx        context.Context
	getOpts    *client.GetOptions
	setOpts    *client.SetOptions
	deleteOpts *client.DeleteOptions
	cb         func(kvb stores.KVBytes)
	config     *config.Config
}

var (
	DefaultEndpoints = []string{"http://127.0.0.1:2379"}
	ReconnectTime    = 2 * time.Second
)

func DefaultStore(cfg *config.Config) (client.KeysAPI, error) {
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

func New(cfg *config.Config) stores.IFace {
	kv, _ := DefaultStore(cfg)

	es := &Store{
		config: cfg,
		kv:     kv,
		ctx:    context.Background(),
	}

	return es
}

func (s *Store) Get(key string) (*stores.KVByte, error) {
	resp, err := s.kv.Get(s.ctx, key, s.getOpts)

	if err != nil {
		return nil, etcdError(err)
	}

	return toKVByte(resp.Node), nil
}

func (s *Store) Set(key string, bts []byte) error {
	_, err := s.kv.Set(s.ctx, key, string(bts), s.setOpts)

	return err
}

func (s *Store) Delete(key string) error {
	_, err := s.kv.Delete(s.ctx, key, s.deleteOpts)

	return etcdError(err)
}

func (s *Store) List(prefix string) (stores.KVBytes, error) {
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

func (s *Store) Register(cb func(kvb stores.KVBytes)) {
	s.cb = cb
}

func (s *Store) Updated(kvs interface{}) {
	kvb := FlattenToKVBytes(kvs.(*client.Node), make(stores.KVBytes, 0))

	s.cb(kvb)
}

func (s *Store) Watch() error {
	watcherOpts := client.WatcherOptions{AfterIndex: 0, Recursive: true}
	w := s.kv.Watcher(s.config.Namespace, &watcherOpts)

	s.Init()

	for {
		r, err := w.Next(context.Background())
		if err != nil {
			printer.LogErrf("Error occurred: %e", err)
			time.Sleep(ReconnectTime)
			continue
		}

		switch r.Action {
		case "set", "update", "create":
			s.Updated(r.Node)
		case "delete":
			resp, err := s.kv.Get(context.Background(), s.config.Namespace, nil)
			if err != nil {
				printer.LogErrf("Error occurred: %e", err)
			}

			s.Updated(resp.Node)
		}

	}
}

// Init etcd watches do not fire an initial event. This method triggers
// a write to the file systems of the entire keyspace.
func (s *Store) Init() {
	opts := &client.GetOptions{
		Recursive: true,
		Sort:      true,
		Quorum:    true,
	}

	resp, err := s.kv.Get(context.Background(), s.config.Namespace, opts)

	if etcdError(err) != nil {
		printer.LogErrf("Error occurred: %e", err)
	}

	if resp != nil {
		s.Updated(resp.Node)
	}
}

func (s *Store) Close() {}

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

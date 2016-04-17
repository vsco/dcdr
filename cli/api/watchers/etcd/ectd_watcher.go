package etcd

import (
	"time"

	"golang.org/x/net/context"

	"github.com/coreos/etcd/client"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/cli/api/stores/etcd"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/config"
)

type Watcher struct {
	config *config.Config
	cb     func(kvb stores.KVBytes)
	api    client.KeysAPI
}

func New(cfg *config.Config) (cw *Watcher) {
	ecfg := client.Config{
		Endpoints:               cfg.Etcd.Endpoints,
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}

	c, err := client.New(ecfg)

	if err != nil {
		printer.LogErrf("watch error: %v", err)
	}

	cw = &Watcher{
		config: cfg,
		api:    client.NewKeysAPI(c),
	}

	return
}

func (cw *Watcher) Register(cb func(kvb stores.KVBytes)) {
	cw.cb = cb
}

func (cw *Watcher) Updated(kvs interface{}) {
	kvb := etcd.FlattenToKVBytes(kvs.(*client.Node), make(stores.KVBytes, 0))

	cw.cb(kvb)
}

func (cw *Watcher) Watch() {
	watcherOpts := client.WatcherOptions{AfterIndex: 0, Recursive: true}
	w := cw.api.Watcher(cw.config.Namespace, &watcherOpts)

	cw.Init()

	for {
		r, err := w.Next(context.Background())
		if err != nil {
			printer.LogErrf("Error occurred: %e", err)
		}

		switch r.Action {
		case "set", "update", "create":
			cw.Updated(r.Node)
		case "delete":
			resp, err := cw.api.Get(context.Background(), cw.config.Namespace, nil)
			if err != nil {
				printer.LogErrf("Error occurred: %e", err)
			}
			cw.Updated(resp.Node)
		}

	}
}

// Init etcd watches do not fire an initial event. This method triggers
// a write to the file systems of the entire keyspace.
func (cw *Watcher) Init() {
	opts := &client.GetOptions{
		Recursive: true,
		Sort:      true,
		Quorum:    true,
	}

	resp, err := cw.api.Get(context.Background(), cw.config.Namespace, opts)

	if err != nil {
		printer.LogErrf("Error occurred: %e", err)
	}

	cw.Updated(resp.Node)
}

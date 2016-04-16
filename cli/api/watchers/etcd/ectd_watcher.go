package etcd

import (
	"log"
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
}

func New(cfg *config.Config) (cw *Watcher) {
	cw = &Watcher{
		config: cfg,
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
	cfg := client.Config{
		Endpoints:               []string{"http://127.0.0.1:2379"},
		Transport:               client.DefaultTransport,
		HeaderTimeoutPerRequest: time.Second,
	}
	c, err := client.New(cfg)
	if err != nil {
		log.Fatal(err)
	}

	kapi := client.NewKeysAPI(c)

	watcherOpts := client.WatcherOptions{AfterIndex: 0, Recursive: true}
	w := kapi.Watcher(cw.config.Namespace, &watcherOpts)

	for {
		r, err := w.Next(context.Background())
		if err != nil {
			printer.LogErrf("Error occurred: %e", err)
		}

		switch r.Action {
		case "set", "update", "create":
			cw.Updated(r.Node)
		case "delete":
			resp, err := kapi.Get(context.Background(), cw.config.Namespace, nil)
			if err != nil {
				printer.LogErrf("Error occurred: %e", err)
			}
			cw.Updated(resp.Node)
		}

	}
}

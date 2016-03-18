package consul

import (
	"github.com/hashicorp/consul/api"
	"github.com/hashicorp/consul/watch"
	"github.com/vsco/dcdr/cli/api/stores"
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
	kvp := kvs.(api.KVPairs)
	kvb, err := stores.KvPairsToKvBytes(kvp)

	if err != nil {
		printer.LogErrf("%v", err)
		return
	}

	cw.cb(kvb)
}

func (cw *Watcher) Watch() {
	params := map[string]interface{}{
		"type":   "keyprefix",
		"prefix": cw.config.Namespace,
	}

	wp, err := watch.Parse(params)
	defer wp.Stop()

	if err != nil {
		printer.LogErrf("%v", err)
	}

	wp.Handler = func(idx uint64, data interface{}) {
		cw.Updated(data)
	}

	if err := wp.Run(""); err != nil {
		printer.LogErrf("Error querying Consul agent: %s", err)
	}
}

package watchers

import "github.com/vsco/dcdr/cli/api/stores"

type KVWatcherIFace interface {
	Register(func(kvb stores.KVBytes))
	Watch()
	Updated(kvs interface{})
}

package resolver

import (
	"log"

	"github.com/vsco/dcdr/cli/api/stores"
	consul_store "github.com/vsco/dcdr/cli/api/stores/consul"
	etcd_store "github.com/vsco/dcdr/cli/api/stores/etcd"
	"github.com/vsco/dcdr/cli/api/watchers"
	consul_watcher "github.com/vsco/dcdr/cli/api/watchers/consul"
	etcd_watcher "github.com/vsco/dcdr/cli/api/watchers/etcd"
	"github.com/vsco/dcdr/config"
)

func LoadWatcher(cfg *config.Config) watchers.KVWatcherIFace {
	switch cfg.Storage {
	case "consul":
		return consul_watcher.New(cfg)
	case "etcd":
		return etcd_watcher.New(cfg)
	default:
		log.Fatalf("invalid storage type %s", cfg.Storage)
		return nil
	}
}

func LoadStore(cfg *config.Config) stores.StoreIFace {
	switch cfg.Storage {
	case "consul":
		c, err := consul_store.NewDefault()

		if err != nil {
			log.Fatalf("could not load consul: %v", err)
		}

		return c
	case "etcd":
		return etcd_store.New(cfg)
	default:
		log.Fatalf("invalid storage type %s", cfg.Storage)
		return nil
	}
}

package resolver

import (
	"log"

	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/cli/api/stores/consul"
	"github.com/vsco/dcdr/cli/api/stores/redis"
	"github.com/vsco/dcdr/config"
)

func LoadStore(cfg *config.Config) stores.IFace {
	switch cfg.Storage {
	case "etcd":
		log.Fatal("etcd is no longer supported")
		return nil
	case "redis":
		r, err := redis.New(cfg)

		if err != nil {
			log.Fatalf("could not load redis: %v", err)
		}

		return r
	default:
		c, err := consul.NewDefault(cfg)

		if err != nil {
			log.Fatalf("could not load consul: %v", err)
		}

		return c
	}
}

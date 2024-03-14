package main

import (
	"fmt"
	"os"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/vsco/dcdr/cli"
	"github.com/vsco/dcdr/cli/api"
	"github.com/vsco/dcdr/cli/api/resolver"
	"github.com/vsco/dcdr/cli/controller"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/cli/repo"
	"github.com/vsco/dcdr/config"
)

func main() {
	cfg := config.LoadConfig()
	store := resolver.LoadStore(cfg)

	rp := repo.New(cfg)

	var stats statsd.ClientInterface

	if cfg.StatsEnabled() {
		var err error
		stats, err = statsd.New(fmt.Sprintf("%s:%d", cfg.Stats.Host, cfg.Stats.Port))

		if err != nil {
			printer.SayErr("%v", err)
			os.Exit(1)
		}
	} else {
		stats = &statsd.NoOpClient{}
	}

	kv := api.New(store, rp, cfg, stats)
	ctrl := controller.New(cfg, kv)

	dcdr := cli.New(ctrl)
	dcdr.Run()
}

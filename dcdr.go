package main

import (
	"os"

	"github.com/PagerDuty/godspeed"
	"github.com/vsco/dcdr/cli"
	"github.com/vsco/dcdr/cli/api"
	"github.com/vsco/dcdr/cli/api/loader"
	"github.com/vsco/dcdr/cli/controller"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/cli/repo"
	"github.com/vsco/dcdr/config"
)

func main() {
	cfg := config.LoadConfig()
	store := loader.LoadStore(cfg)

	rp := repo.New(cfg)

	var gs *godspeed.Godspeed
	var err error
	if cfg.StatsEnabled() {
		gs, err = godspeed.New(cfg.Stats.Host, cfg.Stats.Port, false)

		if err != nil {
			printer.SayErr("%v", err)
			os.Exit(1)
		}
	}

	w := loader.LoadWatcher(cfg)

	kv := api.New(store, rp, w, cfg, gs)
	ctrl := controller.New(cfg, kv)

	dcdr := cli.New(ctrl)
	dcdr.Run()
}

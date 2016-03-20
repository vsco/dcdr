package main

import (
	"fmt"

	"os"

	"github.com/PagerDuty/godspeed"
	"github.com/vsco/dcdr/cli"
	"github.com/vsco/dcdr/cli/api"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/cli/api/watchers/consul"
	"github.com/vsco/dcdr/cli/controller"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/cli/repo"
	"github.com/vsco/dcdr/config"
)

func main() {
	cfg := config.LoadConfig()
	store, err := stores.DefaultConsulStore()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rp := repo.New(cfg)

	var gs *godspeed.Godspeed

	if cfg.StatsEnabled() {
		gs, err = godspeed.New(cfg.Stats.Host, cfg.Stats.Port, false)

		if err != nil {
			printer.SayErr("%v", err)
			os.Exit(1)
		}
	}

	w := consul.New(cfg)

	kv := api.New(store, rp, w, cfg, gs)
	ctrl := controller.New(cfg, kv)

	dcdr := cli.New(ctrl)
	dcdr.Run()
}

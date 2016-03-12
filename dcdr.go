package main

import (
	"fmt"

	"os"

	"github.com/PagerDuty/godspeed"
	"github.com/vsco/dcdr/cli"
	"github.com/vsco/dcdr/cli/api"
	"github.com/vsco/dcdr/cli/api/stores"
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

	kv := api.New(store, rp, cfg.Namespace, gs)
	ctrl := cli.NewController(cfg, kv)

	dcdr := cli.New(ctrl)
	dcdr.Run()
}

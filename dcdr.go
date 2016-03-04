package main

import (
	"fmt"

	"os"

	"github.com/PagerDuty/godspeed"
	"github.com/vsco/dcdr/cli"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/kv"
	"github.com/vsco/dcdr/kv/stores"
	"github.com/vsco/dcdr/repo"
)

func main() {
	cfg := config.LoadConfig()
	store, err := stores.DefaultConsulStore()

	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	rp := repo.New(cfg)

	cmd := ""

	if len(os.Args) > 1 {
		cmd = os.Args[1]
	}

	if cmd != "init" && cmd != "watch" && rp.Enabled() && !rp.Exists() {
		fmt.Printf("%s has not been cloned to %s. see `dcdr help init` for usage\n", cfg.Git.RepoURL, cfg.Git.RepoPath)
		os.Exit(1)
	}

	var gs *godspeed.Godspeed

	if cfg.StatsEnabled() {
		gs, err = godspeed.New(cfg.Stats.Host, cfg.Stats.Port, false)

		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	}

	kv := kv.New(store, rp, cfg.Namespace, gs)
	ctrl := cli.NewController(cfg, kv)

	dcdr := cli.New(ctrl)
	dcdr.Run()
}

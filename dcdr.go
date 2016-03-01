package main

import (
	"fmt"

	"os"

	"github.com/vsco/dcdr/cli"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/kv"
	"github.com/vsco/dcdr/repo"
)

func main() {
	cfg := config.LoadConfig()
	store, err := kv.DefaultConsulStore()

	if err != nil {
		fmt.Println(err)
		return
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

	kv := kv.New(store, rp, cfg.Namespace)
	ctrl := cli.NewController(cfg, kv)

	dcdr := cli.New(ctrl)
	dcdr.Run()
}

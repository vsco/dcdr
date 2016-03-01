package main

import (
	"fmt"

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

	kv := kv.New(store, cfg.Namespace)
	rp := repo.New(cfg)
	ctrl := cli.NewController(cfg, kv, rp)

	dcdr := cli.New(ctrl)
	dcdr.Run()
}

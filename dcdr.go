package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/vsco/dcdr/cli"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/kv"
	"github.com/vsco/dcdr/repo"
)

func main() {
	cfg := config.LoadConfig()
	kv := kv.New(api.DefaultConfig(), cfg.Namespace)
	rp := repo.New(cfg)
	ctrl := cli.NewController(cfg, kv, rp)

	dcdr := cli.New(ctrl)
	dcdr.Run()
}

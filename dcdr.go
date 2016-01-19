package main

import (
	"github.com/hashicorp/consul/api"
	"github.com/vsco/dcdr/cli"
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/git"
	"github.com/vsco/dcdr/models"
)

func main() {
	cfg := models.LoadConfig()
	c := client.New(api.DefaultConfig(), cfg.Namespace)
	g := git.New(cfg)
	ctrl := cli.NewCommandController(cfg, c, g)

	dcdr := cli.New(ctrl)
	dcdr.Run()
}

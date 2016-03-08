package main

import (
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/server"
	"github.com/zenazn/goji"
)

const FixturePath = "./decider.json"

func main() {
	cfg := config.DefaultConfig()
	//cfg.FeatureMapPath = FixturePath

	c, err := client.New(cfg).Watch()

	if err != nil {
		panic(err)
	}

	s := server.NewServer(cfg, goji.DefaultMux, c)
	s.Serve()
}

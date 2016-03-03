package main

import (
	"github.com/vsco/dcdr/server"
	"github.com/vsco/dcdr/watcher"
)

func main() {
	cfg := watcher.TestConfig()
	cfg.ConfigPath = "../../config/decider_fixtures.json"

	server.NewWithConfig(cfg).Init().Serve()
}

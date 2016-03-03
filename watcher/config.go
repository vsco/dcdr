package watcher

import (
	"flag"

	"log"

	decider "github.com/vsco/decider-go"
)

const (
	DefaultConfigPath = "/etc/config/decider_fixtures.json"
	DefaultEndPoint   = "/decider.json"
	DeciderEnvKey     = "decider"
)

type Config struct {
	Decider         decider.DeciderIFace
	ConfigPath      string
	FeatureEndpoint string
}

func DefaultConfig() (c *Config) {
	path := flag.String("cfg", DefaultConfigPath, "path to decider config")
	endpoint := flag.String("endpoint", DefaultEndPoint, "API endpoint path")

	flag.Parse()

	log.Printf("[DECIDER] watching %s. Served on GET %s", *path, *endpoint)

	c = &Config{
		FeatureEndpoint: *endpoint,
		ConfigPath:      *path,
	}

	return
}

func TestConfig() (c *Config) {
	c = &Config{
		FeatureEndpoint: DefaultEndPoint,
		ConfigPath:      "../config/decider_fixtures.json",
	}

	return
}

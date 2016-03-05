package client

const DefaultConfigPath = "/etc/dcdr/decider.json"

type Config struct {
	WatchPath string
}

func DefaultConfig() (c *Config) {
	c = &Config{
		WatchPath: DefaultConfigPath,
	}

	return
}

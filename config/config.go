package config

import (
	"io/ioutil"
	"os"
	"os/user"

	"github.com/hashicorp/hcl"
	"github.com/vsco/dcdr/cli/printer"
)

const (
	DefaultNamespace     = "dcdr"
	DefaultInfoNamespace = DefaultNamespace + "/" + "info"
	DefaultUsername      = "unknown"
	ConfigPath           = "/etc/dcdr/config.hcl"
	DefaultFilePath      = "/etc/dcdr/decider.json"
	DefaultEndpoint      = "/dcdr.json"
	DefaultHost          = ":8000"
)

type Server struct {
	Endpoint string
	Host     string
	JsonRoot string
}

type Stats struct {
	Namespace string
	Host      string
	Port      int
}

type Git struct {
	RepoPath string
	RepoURL  string
}

type Config struct {
	Username       string
	Namespace      string
	FeatureMapPath string
	Git            Git
	Stats          Stats
	Server         Server
}

func (c *Config) GitEnabled() bool {
	return c.Git.RepoURL != ""
}

func (c *Config) StatsEnabled() bool {
	return c.Stats.Host != ""
}

func DefaultConfig() *Config {
	uname := DefaultUsername
	u, _ := user.Current()

	if u != nil {
		uname = u.Username
	}

	return &Config{
		Username:       uname,
		Namespace:      DefaultNamespace,
		FeatureMapPath: DefaultFilePath,
		Server: Server{
			Endpoint: DefaultEndpoint,
			Host:     DefaultHost,
			JsonRoot: DefaultNamespace,
		},
	}
}

func readConfig() *Config {
	bts, err := ioutil.ReadFile(ConfigPath)

	if err != nil {
		printer.SayErr("Could not read %s", ConfigPath)
		os.Exit(1)
	}

	var cfg *Config
	defaults := DefaultConfig()

	err = hcl.Decode(&cfg, string(bts[:]))

	if err != nil {
		printer.SayErr("[dcdr] config parse error %+v", err)
		os.Exit(1)
	}

	if cfg.Namespace == "" {
		cfg.Namespace = defaults.Namespace
	}

	if cfg.Username == "" {
		cfg.Username = defaults.Username
	}

	if cfg.FeatureMapPath == "" {
		cfg.FeatureMapPath = defaults.FeatureMapPath
	}

	if cfg.Server.Host == "" {
		cfg.Server.Host = defaults.Server.Host
	}

	if cfg.Server.Endpoint == "" {
		cfg.Server.Endpoint = defaults.Server.Endpoint
	}

	if cfg.Server.JsonRoot == "" {
		cfg.Server.JsonRoot = defaults.Server.JsonRoot
	}

	return cfg
}

func LoadConfig() *Config {
	if _, err := os.Stat(ConfigPath); err == nil {
		return readConfig()
	} else {
		return DefaultConfig()
	}
}

package config

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/user"

	"github.com/hashicorp/hcl"
)

const (
	DefaultNamespace     = "dcdr"
	DefaultInfoNamespace = DefaultNamespace + "/" + "info"
	DefaultUsername      = "unknown"
	ConfigPath           = "/etc/dcdr/config.hcl"
	DefaultFilePath      = "/etc/dcdr/decider.json"
)

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
	Username  string
	Namespace string
	FilePath  string
	Git       Git
	Stats     Stats
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
		Username:  uname,
		Namespace: DefaultNamespace,
		FilePath:  DefaultFilePath,
	}
}

func readConfig() *Config {
	bts, err := ioutil.ReadFile(ConfigPath)

	if err != nil {
		fmt.Printf("Could not read %s", ConfigPath)
		os.Exit(1)
	}

	var cfg *Config
	defaults := DefaultConfig()

	err = hcl.Decode(&cfg, string(bts[:]))

	if err != nil {
		fmt.Printf("[dcdr] config parse error %+v", err)
		os.Exit(1)
	}

	if cfg.Namespace == "" {
		cfg.Namespace = defaults.Namespace
	}

	if cfg.Username == "" {
		cfg.Username = defaults.Username
	}

	if cfg.FilePath == "" {
		cfg.FilePath = defaults.FilePath
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

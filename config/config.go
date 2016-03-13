package config

import (
	"io/ioutil"
	"os"
	"os/user"

	"fmt"

	"github.com/hashicorp/hcl"
	"github.com/vsco/dcdr/cli/printer"
)

const (
	ConfigFileName       = "config.hcl"
	OutputFileName       = "decider.json"
	DefaultNamespace     = "dcdr"
	DefaultInfoNamespace = DefaultNamespace + "/" + "info"
	DefaultUsername      = "unknown"
	EnvConfigDirOverride = "DCDR_CONFIG_DIR"
	DefaultEndpoint      = "/dcdr.json"
	DefaultHost          = ":8000"
)

var ConfigDir = "/etc/dcdr"

func ConfigPath() string {
	return fmt.Sprintf("%s/%s", ConfigDir, ConfigFileName)
}

func OutputPath() string {
	return fmt.Sprintf("%s/%s", ConfigDir, OutputFileName)
}

var ExampleConfig = []byte(`
// Username = "dcdr admin"
// Namespace = "dcdr"

// Watcher {
//   OutputPath = "/etc/dcdr/decider.json"
// }

// Server {
//   JsonRoot = "dcdr"
//   Endpoint = "/dcdr.json"
// }

// Git {
//   RepoURL = "git@github.com:vsco/decider-test-config.git"
//   RepoPath = "/etc/dcdr/audit"
// }

// Stats {
//   Namespace = "decider"
//   Host = "127.0.0.1"
//   Port = 8126
// }`)

type Server struct {
	Endpoint string
	Host     string
	JsonRoot string
}

type Watcher struct {
	OutputPath string
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
	Username  string
	Namespace string
	Watcher   Watcher
	Git       Git
	Stats     Stats
	Server    Server
}

func (c *Config) GitEnabled() bool {
	return c.Git.RepoPath != ""
}

func (c *Config) PushEnabled() bool {
	return c.Git.RepoURL != ""
}

func (c *Config) StatsEnabled() bool {
	return c.Stats.Host != ""
}

func TestConfig() *Config {
	cfg := DefaultConfig()
	cfg.Watcher.OutputPath = ""

	return cfg
}

func DefaultConfig() *Config {
	uname := DefaultUsername
	u, err := user.Current()

	if err != nil {
		uname = DefaultUsername
	}

	if u != nil {
		uname = u.Username
	}

	return &Config{
		Username:  uname,
		Namespace: DefaultNamespace,
		Watcher: Watcher{
			OutputPath: OutputPath(),
		},
		Server: Server{
			Endpoint: DefaultEndpoint,
			Host:     DefaultHost,
			JsonRoot: DefaultNamespace,
		},
	}
}

func readConfig() *Config {
	bts, err := ioutil.ReadFile(ConfigPath())

	if err != nil {
		printer.SayErr("Could not read %s", ConfigPath())
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

	if cfg.Watcher.OutputPath == "" {
		cfg.Watcher.OutputPath = defaults.Watcher.OutputPath
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
	if v := os.Getenv(EnvConfigDirOverride); v != "" {
		ConfigDir = v
	}

	if _, err := os.Stat(ConfigPath()); err == nil {
		return readConfig()
	} else {
		return DefaultConfig()
	}
}

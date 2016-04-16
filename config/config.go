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
	configFileName       = "config.hcl"
	defaultNamespace     = "dcdr"
	defaultUsername      = "unknown"
	defaultStorage       = "consul"
	envConfigDirOverride = "DCDR_CONFIG_DIR"
	defaultHost          = ":8000"
	defaultEndpoint      = "/dcdr.json"

	// OutputFileName name used for output path.
	OutputFileName = "decider.json"
	// DefaultInfoNamespace path for the info key.
	DefaultInfoNamespace = defaultNamespace + "/" + "info"
)

// ConfigDir default config directory.
var ConfigDir = "/etc/dcdr"

// Path path to config.hcl
func Path() string {
	return fmt.Sprintf("%s/%s", ConfigDir, configFileName)
}

// OutputPath path to write `OutputFileName`
func OutputPath() string {
	return fmt.Sprintf("%s/%s", ConfigDir, OutputFileName)
}

// ExampleConfig an example config written by `dcdr init`
var ExampleConfig = []byte(`
// Username = "dcdr admin"
// Namespace = "dcdr"
// Storage = "consul"

// Etcd {
//	 Endpoints = ["http://127.0.0.1:2379"]
// }

// Consul {
//	 Address = "127.0.0.1:8500"
// }

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

// Server config struct for `dcdr server`
type Server struct {
	Endpoint string
	Host     string
	JSONRoot string
}

// Consul config struct for the consul store. Most of consul
// configuration is handled by environment variables
type Consul struct {
	Address string
}

// Etcd config struct for the etcd store.
type Etcd struct {
	Endpoints []string
}

// Watcher config struct for `dcdr watch`
type Watcher struct {
	OutputPath string
}

// Stats config struct for statsd
type Stats struct {
	Namespace string
	Host      string
	Port      int
}

// Git config struct for the audit repo
type Git struct {
	RepoPath string
	RepoURL  string
}

// Config config struct for the `CLI`, `Client`, and `Server`
type Config struct {
	Username  string
	Namespace string
	Storage   string
	Consul    Consul
	Etcd      Etcd
	Watcher   Watcher
	Git       Git
	Stats     Stats
	Server    Server
}

// GitEnabled checks if a git repo has been configured.
func (c *Config) GitEnabled() bool {
	return c.Git.RepoPath != ""
}

// PushEnabled checks if a git repo has a remote origin.
func (c *Config) PushEnabled() bool {
	return c.Git.RepoURL != ""
}

// StatsEnabled checks if statsd has been configured.
func (c *Config) StatsEnabled() bool {
	return c.Stats.Host != ""
}

// TestConfig used for testing.
func TestConfig() *Config {
	cfg := DefaultConfig()
	cfg.Watcher.OutputPath = ""

	return cfg
}

// DefaultConfig returns a `Config` with default values.
func DefaultConfig() *Config {
	uname := defaultUsername
	u, err := user.Current()

	if err != nil {
		uname = defaultUsername
	}

	if u != nil {
		uname = u.Username
	}

	return &Config{
		Username:  uname,
		Namespace: defaultNamespace,
		Storage:   defaultStorage,
		Watcher: Watcher{
			OutputPath: OutputPath(),
		},
		Server: Server{
			Endpoint: defaultEndpoint,
			Host:     defaultHost,
			JSONRoot: defaultNamespace,
		},
	}
}

// LoadConfig reads config.hcl if found and merges with `DefaultConfig`.
func LoadConfig() *Config {
	if v := os.Getenv(envConfigDirOverride); v != "" {
		ConfigDir = v
	}

	if _, err := os.Stat(Path()); err == nil {
		return readConfig()
	}

	return DefaultConfig()
}

func readConfig() *Config {
	bts, err := ioutil.ReadFile(Path())

	if err != nil {
		printer.SayErr("Could not read %s", Path())
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

	if cfg.Server.JSONRoot == "" {
		cfg.Server.JSONRoot = defaults.Server.JSONRoot
	}

	return cfg
}

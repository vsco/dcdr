package models

import "os/user"

const (
	DefaultNamespace = "decider/features"
	DefaultUsername  = "unknown"
)

type Tunnel struct {
	Host string
	Port int
}

type Git struct {
	RepoPath string
	RepoURL  string
}

type Config struct {
	Username  string
	Namespace string
	Git       Git
	Tunnel    Tunnel
}

func (c *Config) UseGit() bool {
	return c.Git.RepoURL != ""
}

func (c *Config) UseTunnel() bool {
	return c.Tunnel.Host != ""
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
	}
}

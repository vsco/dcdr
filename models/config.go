package models

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/user"

	"github.com/hashicorp/hcl"
)

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

func configPath() string {
	usr, err := user.Current()

	if err != nil {
		log.Fatal(err)
	}

	return usr.HomeDir + "/.dcdr"
}

func readConfig() *Config {
	bts, err := ioutil.ReadFile(configPath())

	if err != nil {
		fmt.Printf("Could not read %s", configPath())
		os.Exit(1)
	}

	cfg := DefaultConfig()

	err = hcl.Decode(&cfg, string(bts[:]))

	if err != nil {
		fmt.Printf("parse error %+v", err)
		os.Exit(1)
	}

	return cfg
}

func LoadConfig() *Config {
	if _, err := os.Stat(configPath()); err == nil {
		return readConfig()
	} else {
		return DefaultConfig()
	}
}

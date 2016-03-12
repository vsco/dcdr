package config

import (
	"os/user"
	"testing"

	"os"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	cfg := LoadConfig()

	assert.Equal(t, cfg.Namespace, DefaultNamespace)

	user, err := user.Current()
	assert.NoError(t, err)

	assert.Equal(t, ConfigPath(), "/etc/dcdr/config.hcl")
	assert.Equal(t, OutputPath(), "/etc/dcdr/decider.json")

	assert.Equal(t, cfg.Username, user.Username)
	assert.Equal(t, cfg.Watcher.OutputPath, OutputPath())
	assert.Equal(t, cfg.Server.Endpoint, DefaultEndpoint)
	assert.Equal(t, cfg.Server.Host, DefaultHost)
	assert.Equal(t, cfg.Server.JsonRoot, DefaultNamespace)
	assert.Equal(t, cfg.Git.RepoPath, "")
	assert.Equal(t, cfg.Git.RepoURL, "")
	assert.Equal(t, cfg.Stats.Host, "")
	assert.Equal(t, cfg.Stats.Port, 0)
	assert.Equal(t, cfg.Stats.Namespace, "")
}

func TestEnvOverride(t *testing.T) {
	os.Setenv(EnvConfigDirOverride, "/tmp/dcdr")
	cfg := LoadConfig()

	assert.Equal(t, ConfigPath(), fmt.Sprintf("%s/%s", os.Getenv(EnvConfigDirOverride), ConfigFileName))
	assert.Equal(t, OutputPath(), fmt.Sprintf("%s/%s", os.Getenv(EnvConfigDirOverride), OutputFileName))
	assert.Equal(t, cfg.Watcher.OutputPath, OutputPath())
}

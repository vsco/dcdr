package config

import (
	"os/user"
	"testing"

	"os"

	"fmt"

	"github.com/stretchr/testify/assert"
)

func TestDefaultConfig(t *testing.T) {
	ConfigDir = "/tmp/does/not/exist"
	cfg := LoadConfig()

	assert.Equal(t, cfg.Namespace, defaultNamespace)

	user, err := user.Current()
	assert.NoError(t, err)

	assert.Equal(t, Path(), fmt.Sprintf("%s/%s", ConfigDir, configFileName))
	assert.Equal(t, OutputPath(), fmt.Sprintf("%s/%s", ConfigDir, OutputFileName))

	assert.False(t, cfg.GitEnabled())
	assert.False(t, cfg.PushEnabled())
	assert.False(t, cfg.StatsEnabled())

	assert.Equal(t, cfg.Username, user.Username)
	assert.Equal(t, cfg.Watcher.OutputPath, OutputPath())
	assert.Equal(t, cfg.Server.Endpoint, defaultEndpoint)
	assert.Equal(t, cfg.Server.Host, defaultHost)
	assert.Equal(t, cfg.Server.JSONRoot, defaultNamespace)
	assert.Equal(t, cfg.Git.RepoPath, "")
	assert.Equal(t, cfg.Git.RepoURL, "")
	assert.Equal(t, cfg.Stats.Host, "")
	assert.Equal(t, cfg.Stats.Port, 0)
	assert.Equal(t, cfg.Stats.Namespace, "")
}

func TestEnvOverride(t *testing.T) {
	os.Setenv(envConfigDirOverride, "/tmp/dcdr")
	cfg := LoadConfig()

	assert.Equal(t, Path(), fmt.Sprintf("%s/%s", os.Getenv(envConfigDirOverride), configFileName))
	assert.Equal(t, OutputPath(), fmt.Sprintf("%s/%s", os.Getenv(envConfigDirOverride), OutputFileName))
	assert.Equal(t, cfg.Watcher.OutputPath, OutputPath())
}

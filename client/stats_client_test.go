package client

import (
	"testing"

	"strings"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/client/stats"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

type MockStatter struct {
	Count map[string]int
}

func (ms *MockStatter) Incr(key string) {
	ms.Count[key]++
}

func (ms *MockStatter) Tags() []string {
	return []string{"tag"}
}

func NewMockStatter() (ms *MockStatter) {
	ms = &MockStatter{
		Count: make(map[string]int),
	}

	return
}

func TestStatsClientIsAvailable(t *testing.T) {
	ft := "feature"
	ms := NewMockStatter()
	c, err := NewStatsClient(&config.Config{}, ms)
	assert.NoError(t, err)

	enabled := c.IsAvailable(ft)
	key := c.statKey(ft, enabled)
	assert.Equal(t, 1, ms.Count[key])
}

func TestStatsClientIsAvailableForID(t *testing.T) {
	ft := "feature-2"
	ms := NewMockStatter()
	c, err := NewStatsClient(&config.Config{}, ms)
	assert.NoError(t, err)

	enabled := c.IsAvailableForID(ft, 1)
	key := c.statKey(ft, enabled)
	assert.Equal(t, 1, ms.Count[key])
}

func TestFormatKey(t *testing.T) {
	ft := "feature-2"
	ms := NewMockStatter()
	cfg := &config.Config{
		Namespace: "test",
	}
	c, err := NewStatsClient(cfg, ms)
	assert.NoError(t, err)

	expected := strings.Join([]string{cfg.Namespace, models.DefaultScope, ft, stats.Enabled}, stats.JoinWith)
	assert.Equal(t, expected, c.statKey(ft, true))

	c.scopes = []string{"a", "b/c"}

	expected = strings.Join([]string{cfg.Namespace, "a.b.c", ft, stats.Enabled}, stats.JoinWith)
	assert.Equal(t, expected, c.statKey(ft, true))

	expected = strings.Join([]string{cfg.Namespace, "a.b.c", ft, stats.Disabled}, stats.JoinWith)
	assert.Equal(t, expected, c.statKey(ft, false))
}

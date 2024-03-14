package client

import (
	"strings"
	"testing"
	"time"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

type MockStatter struct {
	count map[string]int
}

func (ms *MockStatter) Gauge(name string, value float64, tags []string, rate float64) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) GaugeWithTimestamp(name string, value float64, tags []string, rate float64, timestamp time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) Count(name string, value int64, tags []string, rate float64) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) CountWithTimestamp(name string, value int64, tags []string, rate float64, timestamp time.Time) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) Histogram(name string, value float64, tags []string, rate float64) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) Distribution(name string, value float64, tags []string, rate float64) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) Decr(name string, tags []string, rate float64) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) Incr(name string, tags []string, rate float64) error {
	//TODO implement me
	ms.count[name]++
	return nil
}

func (ms *MockStatter) Set(name string, value string, tags []string, rate float64) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) Timing(name string, value time.Duration, tags []string, rate float64) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) TimeInMilliseconds(name string, value float64, tags []string, rate float64) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) Event(e *statsd.Event) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) SimpleEvent(title, text string) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) ServiceCheck(sc *statsd.ServiceCheck) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) SimpleServiceCheck(name string, status statsd.ServiceCheckStatus) error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) Close() error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) Flush() error {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) IsClosed() bool {
	//TODO implement me
	panic("implement me")
}

func (ms *MockStatter) GetTelemetry() statsd.Telemetry {
	//TODO implement me
	panic("implement me")
}

func NewMockStatter() (ms *MockStatter) {
	ms = &MockStatter{
		count: make(map[string]int),
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
	assert.Equal(t, 1, ms.count[key])
}

func TestStatsClientIsAvailableForID(t *testing.T) {
	ft := "feature-2"
	ms := NewMockStatter()
	c, err := NewStatsClient(&config.Config{}, ms)
	assert.NoError(t, err)

	enabled := c.IsAvailableForID(ft, 1)
	key := c.statKey(ft, enabled)
	assert.Equal(t, 1, ms.count[key])
}

func TestFormatKey(t *testing.T) {
	ft := "feature-2"
	ms := NewMockStatter()
	cfg := &config.Config{
		Namespace: "test",
	}
	c, err := NewStatsClient(cfg, ms)
	assert.NoError(t, err)

	expected := strings.Join([]string{cfg.Namespace, models.DefaultScope, ft, "enabled"}, ".")
	assert.Equal(t, expected, c.statKey(ft, true))

	c.scopes = []string{"a", "b/c"}

	expected = strings.Join([]string{cfg.Namespace, "a.b.c", ft, "enabled"}, ".")
	assert.Equal(t, expected, c.statKey(ft, true))

	expected = strings.Join([]string{cfg.Namespace, "a.b.c", ft, "disabled"}, ".")
	assert.Equal(t, expected, c.statKey(ft, false))
}

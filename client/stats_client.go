package client

import (
	"strings"

	"github.com/DataDog/datadog-go/v5/statsd"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

// StatsClient delegates `Client` methods with metrics.
type StatsClient struct {
	Client
	stats statsd.ClientInterface
}

// NewStatsClient creates a new client.
func NewStatsClient(cfg *config.Config, stats statsd.ClientInterface) (sc *StatsClient, err error) {
	sc = &StatsClient{
		stats: stats,
	}

	c, err := New(cfg)

	sc.Client = *c
	sc.Client.Watch()

	return
}

// NewStatsClient creates a new client.
func NewStatsDefault(stats statsd.ClientInterface) (sc *StatsClient, err error) {
	sc = &StatsClient{
		stats: stats,
	}

	c, err := NewDefault()

	if err != nil {
		return sc, err
	}

	sc.Client = *c
	sc.Client.Watch()

	return
}

// IsAvailable delegates `IsAvailable` and increments the provided `feature` status.
func (sc *StatsClient) IsAvailable(feature string) bool {
	enabled := sc.Client.IsAvailable(feature)
	defer sc.Incr(feature, enabled, 1)

	return enabled
}

// IsAvailableForID delegates `IsAvailableForID` and increments the provided `feature` status.
func (sc *StatsClient) IsAvailableForID(feature string, id uint64) bool {
	enabled := sc.Client.IsAvailableForID(feature, id)
	defer sc.Incr(feature, enabled, 1)

	return enabled
}

// ScaleValue delegates `ScaleValue`.
func (sc *StatsClient) ScaleValue(feature string, min float64, max float64) float64 {
	return sc.Client.ScaleValue(feature, min, max)
}

// UpdateFeatures delegates `UpdateFeatures`.
func (sc *StatsClient) UpdateFeatures(bts []byte) {
	sc.Client.UpdateFeatures(bts)
}

// FeatureExists delegates `FeatureExists`.
func (sc *StatsClient) FeatureExists(feature string) bool {
	return sc.Client.FeatureExists(feature)
}

// Features delegates `Features`.
func (sc *StatsClient) Features() models.FeatureScopes {
	return sc.Client.Features()
}

// Scopes delegates `Scopes`.
func (sc *StatsClient) Scopes() []string {
	return sc.Client.Scopes()
}

// Incr increments the formatted `statKey`.
func (sc *StatsClient) Incr(feature string, enabled bool, sampleRate float64) {
	key := sc.statKey(feature, enabled)
	sc.stats.Incr(key, []string{}, sampleRate)
}

func (sc *StatsClient) statKey(feature string, enabled bool) string {
	status := "enabled"

	if enabled == false {
		status = "disabled"
	}

	scopes := models.DefaultScope

	if len(sc.Client.Scopes()) > 0 {
		scopes = strings.Replace(strings.Join(sc.Client.Scopes(), "."), "/", ".", -1)
	}

	return strings.Join([]string{sc.config.Namespace, scopes, feature, status}, ".")
}

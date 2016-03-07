package client

import (
	"strings"

	"github.com/vsco/dcdr/client/models"
	"github.com/vsco/dcdr/config"
)

type Statter interface {
	Incr(feature string)
	Tags() []string
}

const (
	Enabled  = "enabled"
	Disabled = "disabled"
	JoinWith = "."
)

type statsClient struct {
	Client
	stats Statter
}

func NewStatsClient(cfg *config.Config, stats Statter) (sc *statsClient) {
	sc = &statsClient{
		stats: stats,
	}

	sc.Client.config = cfg

	return
}

func (sc *statsClient) IsAvailable(feature string) bool {
	enabled := sc.Client.IsAvailable(feature)
	defer sc.Incr(feature, enabled)

	return enabled
}

func (sc *statsClient) IsAvailableForId(feature string, id uint64) bool {
	enabled := sc.Client.IsAvailableForId(feature, id)
	defer sc.Incr(feature, enabled)

	return enabled
}

func (sc *statsClient) ScaleValue(feature string, min float64, max float64) float64 {
	return sc.Client.ScaleValue(feature, min, max)
}

func (sc *statsClient) UpdateFeatures(bts []byte) {
	sc.Client.UpdateFeatures(bts)
}

func (sc *statsClient) FeatureExists(feature string) bool {
	return sc.Client.FeatureExists(feature)
}

func (sc *statsClient) Features() models.Features {
	return sc.Client.Features()
}

func (sc *statsClient) Scopes() []string {
	return sc.Client.Scopes()
}

func (sc *statsClient) Incr(feature string, enabled bool) {
	key := sc.statKey(feature, enabled)
	sc.stats.Incr(key)
}

func (sc *statsClient) statKey(feature string, enabled bool) string {
	status := Enabled

	if enabled == false {
		status = Disabled
	}

	scopes := models.DefaultScope

	if len(sc.Client.Scopes()) > 0 {
		scopes = strings.Replace(strings.Join(sc.Client.Scopes(), JoinWith), "/", JoinWith, -1)
	}

	return strings.Join([]string{sc.config.Namespace, scopes, feature, status}, JoinWith)
}

package mock

import (
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

// New creates a `Client` with an empty `FeatureMap` and `Config`.
func New() (d *Client) {
	c, _ := client.New(&config.Config{
		Watcher: config.Watcher{
			OutputPath: "",
		},
	})

	d = &Client{
		Client:     *c,
		featureMap: models.EmptyFeatureMap(),
	}

	d.Client.SetFeatureMap(d.featureMap)
	return
}

// Client mock `Client` for testing.
type Client struct {
	client.Client
	featureMap *models.FeatureMap
}

// EnableBoolFeature set a boolean feature to true
func (d *Client) EnableBoolFeature(feature string) {
	d.FeatureMap().Dcdr.Defaults()[feature] = true
	d.MergeScopes()
}

// DisableBoolFeature set a boolean feature to false
func (d *Client) DisableBoolFeature(feature string) {
	d.Client.FeatureMap().Dcdr.Defaults()[feature] = false
	d.MergeScopes()
}

// SetPercentileFeature set a percentile feature to an arbitrary value
func (d *Client) SetPercentileFeature(feature string, val float64) {
	d.Client.FeatureMap().Dcdr.Defaults()[feature] = val
	d.MergeScopes()
}

// EnablePercentileFeature set a percentile feature to true
func (d *Client) EnablePercentileFeature(feature string) {
	d.SetPercentileFeature(feature, 1.0)
}

// DisablePercentileFeature set a percentile feature to false
func (d *Client) DisablePercentileFeature(feature string) {
	d.SetPercentileFeature(feature, 0.0)
}

// Features `features` accessor
func (c *Client) Features() models.FeatureScopes {
	return c.features
}

// Watch noop for tests.
func (d *Client) Watch() *Client {
	return d
}

package mock

import (
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/client/models"
	"github.com/vsco/dcdr/config"
)

func New() (d *Client) {
	d = &Client{
		Client:     *client.New(&config.Config{}),
		featureMap: models.EmptyFeatureMap(),
	}

	d.Client.SetFeatureMap(d.featureMap)
	return
}

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

// EnablePercentileFeature set a percentile feature to true
func (d *Client) EnablePercentileFeature(feature string) {
	d.Client.FeatureMap().Dcdr.Defaults()[feature] = 1.0
	d.MergeScopes()
}

// DisablePercentileFeature set a percentile feature to false
func (d *Client) DisablePercentileFeature(feature string) {
	d.Client.FeatureMap().Dcdr.Defaults()[feature] = 0.0
	d.MergeScopes()
}

func (d *Client) Watch() *Client {
	return d
}

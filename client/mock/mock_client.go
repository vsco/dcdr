package mock

import (
	"github.com/vsco/dcdr/client"
	"github.com/vsco/dcdr/client/models"
	"github.com/vsco/dcdr/config"
)

func New() (d *mockDecider) {
	d = &mockDecider{
		Client:     *client.New(&config.Config{}),
		featureMap: models.EmptyFeatureMap(),
	}

	d.Client.SetFeatureMap(d.featureMap)
	return
}

type mockDecider struct {
	client.Client
	featureMap *models.FeatureMap
}

// EnableBoolFeature set a boolean feature to true
func (d *mockDecider) EnableBoolFeature(feature string) {
	d.FeatureMap().Dcdr.Defaults()[feature] = true
	d.MergeScopes()
}

// DisableBoolFeature set a boolean feature to false
func (d *mockDecider) DisableBoolFeature(feature string) {
	d.Client.FeatureMap().Dcdr.Defaults()[feature] = false
	d.MergeScopes()
}

// EnablePercentileFeature set a percentile feature to true
func (d *mockDecider) EnablePercentileFeature(feature string) {
	d.Client.FeatureMap().Dcdr.Defaults()[feature] = 1.0
	d.MergeScopes()
}

// DisablePercentileFeature set a percentile feature to false
func (d *mockDecider) DisablePercentileFeature(feature string) {
	d.Client.FeatureMap().Dcdr.Defaults()[feature] = 0.0
	d.MergeScopes()
}

func (d *mockDecider) Watch() *mockDecider {
	return d
}

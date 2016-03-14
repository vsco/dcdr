package client

import (
	"hash/crc32"
	"strconv"

	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/client/models"
	"github.com/vsco/dcdr/client/watcher"
	"github.com/vsco/dcdr/config"
)

// ClientIFace interface for Decider Clients
type ClientIFace interface {
	IsAvailable(feature string) bool
	IsAvailableForId(feature string, id uint64) bool
	ScaleValue(feature string, min float64, max float64) float64
	UpdateFeatures(bts []byte)
	FeatureExists(feature string) bool
	Features() models.Features
	FeatureMap() *models.FeatureMap
	ScopedMap() *models.FeatureMap
	Scopes() []string
	CurrentSha() string
	WithScopes(scopes ...string) *Client
}

// Client handles access to the FeatureMap
type Client struct {
	featureMap *models.FeatureMap
	config     *config.Config
	watcher    watcher.WatcherIFace
	features   models.Features
	scopes     []string
}

// New creates a new Client with a custom Config
func New(cfg *config.Config) (c *Client) {
	c = &Client{
		config: cfg,
	}

	if c.config.Watcher.OutputPath != "" {
		c.watcher = watcher.NewWatcher(c.config.Watcher.OutputPath)
		printer.Say("started watching %s", c.config.Watcher.OutputPath)
	}

	return
}

// NewDefault creates a new default Client
func NewDefault() (c *Client) {
	c = New(config.LoadConfig())

	return
}

// WithScopes creates a new Client from an existing one that is "scoped"
// to the provided scopes param. `scopes` are provided in priority order.
// For example, when given WithScopes("a", "b", "c"). Keys found in "a"
// will override the same keys found in "b" and so on for "c".
//
// The provided `scopes` are appended to the existing Client's `scopes`,
// merged, and then a new `Watcher` is assigned to the new `Client` so
// that future changes to the `FeatureMap` will be observed.
func (c *Client) WithScopes(scopes ...string) *Client {
	if len(scopes) == 0 {
		return c
	}

	if len(scopes) == 1 && scopes[0] == "" {
		return c
	}

	newScopes := append(scopes, c.scopes...)

	newClient := &Client{
		featureMap: c.FeatureMap(),
		scopes:     newScopes,
		config:     c.config,
	}

	newClient.MergeScopes()

	return newClient
}

// MergeScopes delegates merging to the underlying `FeatureMap`
func (c *Client) MergeScopes() {
	if c.featureMap != nil {
		c.features = c.featureMap.Dcdr.MergedScopes(c.scopes...)
	}
}

// Scopes `scopes` accessor
func (c *Client) Scopes() []string {
	return c.scopes
}

// SetFeatureMap assigns a `FeatureMap` and merges the current
// scopes. When git is enabled a new `FeatureMap` will not be
// assigned unless its `CurrentSha` is different from the one
// currently found in `CurrentSha()`.
func (c *Client) SetFeatureMap(fm *models.FeatureMap) *Client {
	if c.config.GitEnabled() && c.CurrentSha() == fm.Dcdr.CurrentSha() {
		return c
	}

	c.featureMap = fm

	c.MergeScopes()

	return c
}

// FeatureMap `featureMap` accessor
func (c *Client) FeatureMap() *models.FeatureMap {
	if c.featureMap != nil {
		return c.featureMap
	} else {
		return models.EmptyFeatureMap()
	}
}

// ScopedMap a `FeatureMap` containing only merged features.
// Mostly used for JSON output.
func (c *Client) ScopedMap() *models.FeatureMap {
	fm := models.EmptyFeatureMap()
	fm.Dcdr.Features = c.Features()
	fm.Dcdr.Info = c.FeatureMap().Dcdr.Info

	return fm
}

// Features `features` accessor
func (c *Client) Features() models.Features {
	return c.features
}

// CurrentSha accessor for the underlying `CurrentSha` from
// the `FeatureMap`
func (c *Client) CurrentSha() string {
	return c.FeatureMap().Dcdr.Info.CurrentSha
}

// FeatureExists checks the existence of a key
func (c *Client) FeatureExists(feature string) bool {
	_, exists := c.Features()[feature]

	return exists
}

// IsAvailable used to check features with boolean values.
func (c *Client) IsAvailable(feature string) bool {
	val, exists := c.Features()[feature]

	switch val.(type) {
	case bool:
		return exists && val.(bool)
	default:
		return false
	}
}

// IsAvailableForId used to check features with float values between 0.0-1.0.
func (c *Client) IsAvailableForId(feature string, id uint64) bool {
	val, exists := c.Features()[feature]

	switch val.(type) {
	case float64, int:
		return exists && c.withinPercentile(id, val.(float64), feature)
	default:
		return false
	}
}

// UpdateFeatures creates and assigns a new `FeatureMap` from a
// Marshalled JSON byte array
func (c *Client) UpdateFeatures(bts []byte) {
	fm, err := models.NewFeatureMap(bts)

	if err != nil {
		printer.SayErr("parse error: %v", err)
		return
	}

	c.SetFeatureMap(fm)
}

// ScaleValue returns a value scaled between min and max
// given the current value of the feature.
//
// Given the K/V dcdr/features/scalar => 0.5
// ScaleValue("scalar", 0, 10) => 5
func (c *Client) ScaleValue(feature string, min float64, max float64) float64 {
	val, exists := c.Features()[feature]

	if !exists {
		return min
	}

	switch val.(type) {
	case float64, int:
		return min + (max-min)*val.(float64)
	default:
		return min
	}
}

// Watch initializes the `Watcher`, registers the `UpdateFeatures`
// method with it and spawns the watch in a go routine returning the
// `Client` for a fluent interface.
func (c *Client) Watch() (*Client, error) {
	if c.watcher != nil {
		err := c.watcher.Init()

		if err != nil {
			return nil, err
		}

		c.watcher.Register(c.UpdateFeatures)
		go c.watcher.Watch()
	}

	return c, nil
}

func (c *Client) withinPercentile(id uint64, val float64, feature string) bool {
	uid := c.crc(id, feature)
	percentage := uint32(val * 100)

	return uid%100 < percentage
}

func (c *Client) crc(id uint64, feature string) uint32 {
	b := []byte(feature + strconv.FormatInt(int64(id), 10))

	return crc32.ChecksumIEEE(b)
}

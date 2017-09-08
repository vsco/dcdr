package client

import (
	"hash/crc32"
	"strconv"

	"os"

	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/client/watcher"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

// IFace interface for Decider Clients
type IFace interface {
	IsAvailable(feature string) bool
	IsAvailableForID(feature string, id uint64) bool
	ScaleValue(feature string, min float64, max float64) float64
	UpdateFeatures(bts []byte)
	FeatureExists(feature string) bool
	Features() models.FeatureScopes
	FeatureMap() *models.FeatureMap
	SetFeatureMap(fm *models.FeatureMap) *Client
	ScopedMap() *models.FeatureMap
	Scopes() []string
	Info() *models.Info
	WithScopes(scopes ...string) *Client
}

// Client handles access to the `FeatureMap`
type Client struct {
	featureMap *models.FeatureMap
	config     *config.Config
	watcher    watcher.IFace
	features   models.FeatureScopes
	scopes     []string
}

// New creates a new Client with a custom Config
func New(cfg *config.Config) (c *Client, err error) {
	c = &Client{
		config: cfg,
	}

	if c.config.Watcher.OutputPath != "" {
		_, err = os.Stat(c.config.Watcher.OutputPath)

		if err != nil {
			return
		}

		c.watcher = watcher.New(c.config.Watcher.OutputPath)
		_, err = c.Watch()
	}

	return
}

// NewDefault creates a new default Client
func NewDefault() (c *Client, err error) {
	c, err = New(config.LoadConfig())

	return
}

// WithScopes creates a new Client from `c` that is "scoped"
// to the provided scopes argument. `scopes` are provided in priority order.
// For example, when given WithScopes("a", "b", "c"). Keys found in "a"
// will override the same keys found in "b" and so on for "c".
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
// assigned unless its `CurrentSHA` is different from the one
// currently found in `CurrentSHA()`.
func (c *Client) SetFeatureMap(fm *models.FeatureMap) *Client {
	if c.config.GitEnabled() && c.Info().CurrentSHA == fm.Dcdr.CurrentSHA() {
		return c
	}

	c.featureMap = fm

	c.MergeScopes()

	return c
}

// FeatureMap `featureMap` accessor. Returns an empty `FeatureMap`
// if the `featureMap` is nil.
func (c *Client) FeatureMap() *models.FeatureMap {
	if c.featureMap != nil {
		return c.featureMap
	}

	return models.EmptyFeatureMap()
}

// ScopedMap a `FeatureMap` containing only merged features and `Info`.
func (c *Client) ScopedMap() *models.FeatureMap {
	fm := models.EmptyFeatureMap()
	fm.Dcdr.FeatureScopes = c.Features()
	fm.Dcdr.Info = c.FeatureMap().Dcdr.Info

	return fm
}

// Features `features` accessor
func (c *Client) Features() models.FeatureScopes {
	return c.features
}

// Info accessor for the underlying `Info` from `FeatureMap`
func (c *Client) Info() *models.Info {
	return c.FeatureMap().Dcdr.Info
}

// FeatureExists checks the existence of a key
func (c *Client) FeatureExists(feature string) bool {
	_, exists := c.Features()[feature]

	return exists
}

// IsAvailable used to check features with boolean values. Returns false
// if a non-boolean type `feature` is passed.
func (c *Client) IsAvailable(feature string) bool {
	val, exists := c.Features()[feature]

	switch val.(type) {
	case bool:
		return exists && val.(bool)
	default:
		return false
	}
}

// IsAvailableForID used to check features with float values between 0.0-1.0.
// Returns false if a non-percentile type `feature` is passed.
func (c *Client) IsAvailableForID(feature string, id uint64) bool {
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
		printer.SayErr("parse error: %v, feature payload: %s", err, bts)
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
// method, and spawns the watch in a go routine returning the
// `Client` for a fluent interface.
func (c *Client) Watch() (*Client, error) {
	if c.watcher != nil {
		err := c.watcher.Init()

		if err != nil {
			return nil, err
		}

		c.watcher.Register(c.UpdateFeatures)

		// Load initial values into `FeatureMap`
		c.watcher.UpdateBytes()
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

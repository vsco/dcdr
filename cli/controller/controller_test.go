package controller

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/tucnak/climax"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

const (
	Error   = 1
	Success = 0
)

type MockClient struct {
	Features models.Features
	Feature  *models.Feature
	Error    error
}

func NewMockClient(f *models.Feature, fts models.Features, err error) (m *MockClient) {
	m = &MockClient{
		Features: fts,
		Feature:  f,
		Error:    err,
	}

	return
}

func (m *MockClient) Get(key string, v interface{}) error {
	return m.Error
}

func (m *MockClient) Set(ft *models.Feature) error {
	return m.Error
}

func (m *MockClient) Delete(key string, scope string) error {
	return m.Error
}

func (m *MockClient) Namespace() string {
	return "dcdr"
}

func (m *MockClient) List(prefix string, scope string) (models.Features, error) {
	return m.Features, m.Error
}

func (m *MockClient) GetInfo() (*models.Info, error) {
	return nil, m.Error
}

func (m *MockClient) Commit(ft *models.Feature, deleted bool) error {
	return m.Error
}

func (m *MockClient) Push() error {
	return m.Error
}

func (m *MockClient) UpdateCurrentSHA() (string, error) {
	return "", m.Error
}

func (m *MockClient) InitRepo(create bool) error {
	return m.Error
}

func (m *MockClient) Watch() {}

func TestListEmptyFeatures(t *testing.T) {
	cfg := config.DefaultConfig()
	c := NewMockClient(nil, nil, nil)
	ctl := New(cfg, c)

	ctx := climax.Context{
		Variable: map[string]string{},
	}

	code := ctl.List(ctx)

	assert.Equal(t, Error, code)
}

func TestListFeatures(t *testing.T) {
	cfg := config.DefaultConfig()
	fts := models.Features{
		models.Feature{
			Key:   "test",
			Value: true,
		},
	}
	c := NewMockClient(nil, fts, nil)
	ctl := New(cfg, c)

	ctx := climax.Context{
		Variable: map[string]string{},
	}

	code := ctl.List(ctx)

	assert.Equal(t, Success, code)
}

func TestSet(t *testing.T) {
	cfg := config.DefaultConfig()
	fts := models.Features{
		models.Feature{
			Key:   "test",
			Value: true,
		},
	}
	c := NewMockClient(nil, fts, nil)
	ctl := New(cfg, c)

	ctx := climax.Context{
		Variable: map[string]string{"name": "null-test"},
	}

	code := ctl.Set(ctx)

	assert.Equal(t, Success, code)
}

func newParseController() *Controller {
	return New(config.DefaultConfig(), NewMockClient(nil, nil, nil))
}

func parseCtx(vars map[string]string) climax.Context {
	return climax.Context{Variable: vars}
}

func TestParseContextExplicitTypes(t *testing.T) {
	ctl := newParseController()

	f, err := ctl.ParseContext(parseCtx(map[string]string{
		"name": "min-log-level", "value": "debug", "type": "string",
	}))
	assert.NoError(t, err)
	assert.Equal(t, models.String, f.FeatureType)
	assert.Equal(t, "debug", f.Value)

	// string type stores bool/number-looking values verbatim
	f, err = ctl.ParseContext(parseCtx(map[string]string{
		"name": "literal", "value": "true", "type": "string",
	}))
	assert.NoError(t, err)
	assert.Equal(t, models.String, f.FeatureType)
	assert.Equal(t, "true", f.Value)

	f, err = ctl.ParseContext(parseCtx(map[string]string{
		"name": "flag", "value": "true", "type": "boolean",
	}))
	assert.NoError(t, err)
	assert.Equal(t, models.Boolean, f.FeatureType)
	assert.Equal(t, true, f.Value)

	f, err = ctl.ParseContext(parseCtx(map[string]string{
		"name": "flag", "value": "0.5", "type": "percentile",
	}))
	assert.NoError(t, err)
	assert.Equal(t, models.Percentile, f.FeatureType)
	assert.Equal(t, 0.5, f.Value)
}

func TestParseContextExplicitTypeErrors(t *testing.T) {
	ctl := newParseController()

	_, err := ctl.ParseContext(parseCtx(map[string]string{
		"name": "flag", "value": "debug", "type": "nope",
	}))
	assert.Equal(t, errInvalidType, err)

	_, err = ctl.ParseContext(parseCtx(map[string]string{
		"name": "flag", "value": "notabool", "type": "boolean",
	}))
	assert.Equal(t, errInvalidBool, err)

	_, err = ctl.ParseContext(parseCtx(map[string]string{
		"name": "flag", "value": "notanumber", "type": "percentile",
	}))
	assert.Equal(t, errInvalidRange, err)

	_, err = ctl.ParseContext(parseCtx(map[string]string{
		"name": "flag", "value": "2.0", "type": "percentile",
	}))
	assert.Equal(t, errInvalidRange, err)
}

func TestParseContextInference(t *testing.T) {
	ctl := newParseController()

	// bool/percentile still work without -type
	f, err := ctl.ParseContext(parseCtx(map[string]string{"name": "flag", "value": "false"}))
	assert.NoError(t, err)
	assert.Equal(t, models.Boolean, f.FeatureType)

	f, err = ctl.ParseContext(parseCtx(map[string]string{"name": "flag", "value": "0.25"}))
	assert.NoError(t, err)
	assert.Equal(t, models.Percentile, f.FeatureType)

	// a string value without -type is rejected
	_, err = ctl.ParseContext(parseCtx(map[string]string{"name": "flag", "value": "debug"}))
	assert.Equal(t, errTypeRequiredForString, err)

	// a numeric typo infers to string and is likewise rejected
	_, err = ctl.ParseContext(parseCtx(map[string]string{"name": "flag", "value": "0.5x"}))
	assert.Equal(t, errTypeRequiredForString, err)

	// out-of-range percentile without -type
	_, err = ctl.ParseContext(parseCtx(map[string]string{"name": "flag", "value": "2.0"}))
	assert.Equal(t, errInvalidRange, err)
}

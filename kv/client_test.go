package kv

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/models"
)

type MockStore struct {
	Item  *models.Feature
	Items models.Features
	Err   error
}

func (ms *MockStore) List(prefix string) (models.Features, error) {
	return ms.Items, ms.Err
}

func (ms *MockStore) Get(key string) (*models.Feature, error) {
	return ms.Item, ms.Err
}

func (ms *MockStore) Set(f *models.Feature) error {
	return ms.Err
}

func (ms *MockStore) Delete(key string) error {
	return ms.Err
}

type MockRepo struct {
	error   error
	sha     string
	exists  bool
	enabled bool
}

func (mr *MockRepo) Clone() error {
	return mr.error
}

func (mr *MockRepo) Commit(features models.Features, msg string) error {
	return mr.error
}

func (mr *MockRepo) Create() error {
	return mr.error
}

func (mr *MockRepo) Exists() bool {
	return mr.exists
}

func (mr *MockRepo) Enabled() bool {
	return mr.enabled
}

func (mr *MockRepo) Pull() error {
	return mr.error
}

func (mr *MockRepo) CurrentSha() (string, error) {
	return mr.sha, mr.error
}

func (mr *MockRepo) Init() {
}

func TestSetRequestToFeature(t *testing.T) {
	sr := &SetRequest{
		Key:     "key",
		Scope:   "scope",
		Value:   0.5,
		Comment: "comment",
		User:    "user",
	}

	ft, err := sr.ToFeature()

	assert.NoError(t, err)
	assert.Equal(t, sr.Key, ft.Key)
	assert.Equal(t, sr.Value, ft.Value)
	assert.Equal(t, sr.Comment, ft.Comment)
	assert.Equal(t, sr.User, ft.UpdatedBy)
	assert.Equal(t, models.Percentile, ft.FeatureType)
}

func TestClientSet(t *testing.T) {
	sr := &SetRequest{
		Key:       "key",
		Scope:     "scope",
		Namespace: "namespace",
		Value:     0.5,
		Comment:   "comment",
		User:      "user",
	}

	c := New(&MockStore{}, &MockRepo{}, sr.Namespace, nil)

	err := c.Set(sr)

	assert.NoError(t, err)
}

func TestClientSetExisting(t *testing.T) {
	sr := &SetRequest{
		Key:       "key",
		Scope:     "scope",
		Namespace: "namespace",
		User:      "user",
	}

	orig := &models.Feature{
		Key:       "key",
		Scope:     "scope",
		Namespace: "namespace",
		Value:     0.5,
		Comment:   "comment",
		UpdatedBy: "user",
	}

	c := New(&MockStore{
		Item: orig,
	}, &MockRepo{}, sr.Namespace, nil)

	err := c.Set(sr)

	assert.NoError(t, err)
}

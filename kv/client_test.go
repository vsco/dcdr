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

	c := New(&MockStore{}, sr.Namespace)

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
	}, sr.Namespace)

	err := c.Set(sr)

	assert.NoError(t, err)
}

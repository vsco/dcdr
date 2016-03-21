package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

func TestClientSet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")

	c := New(NewMockStore(ft, nil), &MockRepo{}, nil, config.DefaultConfig(), nil)

	err := c.Set(ft)

	assert.NoError(t, err)
}

func TestClientSetExisting(t *testing.T) {
	update := models.NewFeature("test", nil, "c", "u", "s", "n")
	orig := models.NewFeature("test", 0.5, "c", "u", "s", "n")

	c := New(NewMockStore(orig, nil), &MockRepo{}, nil, config.DefaultConfig(), nil)

	err := c.Set(update)

	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, nil, config.DefaultConfig(), nil)

	fts, err := c.List("test", "")

	assert.Nil(t, err)
	assert.Equal(t, models.Features{*ft}, fts)
}

func TestGet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, nil, config.DefaultConfig(), nil)

	var f *models.Feature
	err := c.Get("test", &f)

	assert.Nil(t, err)
	assert.Equal(t, f, ft)
}

func TestNilGet(t *testing.T) {
	cs := NewMockStore(nil, nil)
	c := New(cs, &MockRepo{}, nil, config.DefaultConfig(), nil)

	var f *models.Feature
	err := c.Get("test", &f)

	assert.EqualError(t, err, "dcdr/test not found")
	assert.Nil(t, f)
}

func TestSet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, nil, config.DefaultConfig(), nil)

	err := c.Set(ft)

	assert.Nil(t, err)
}

func TestSetErrorOnNilValue(t *testing.T) {
	ft := models.NewFeature("test", nil, "c", "u", "s", "n")
	cs := NewMockStore(nil, nil)
	c := New(cs, &MockRepo{}, nil, config.DefaultConfig(), nil)

	err := c.Set(ft)

	assert.Equal(t, ErrNilValue, err)
}

func TestTypeChangeErrorSet(t *testing.T) {
	orig := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	bad := models.NewFeature("test", false, "c", "u", "s", "n")

	cs := NewMockStore(orig, nil)
	c := New(cs, nil, nil, config.DefaultConfig(), nil)

	err := c.Set(bad)
	assert.Equal(t, ErrTypeChange, err)
}

func TestSetWithError(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	e := errors.New("")
	cs := NewMockStore(ft, e)
	c := New(cs, nil, nil, config.DefaultConfig(), nil)

	err := c.Set(ft)

	assert.Equal(t, e, err)
}

func TestDelete(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, nil, config.DefaultConfig(), nil)

	err := c.Delete(ft.Key, "")

	assert.Nil(t, err)
}

func TestDeleteWithError(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	e := errors.New("")
	cs := NewMockStore(ft, e)
	c := New(cs, &MockRepo{}, nil, config.DefaultConfig(), nil)

	err := c.Delete(ft.Key, "")

	assert.Equal(t, e, err)
}

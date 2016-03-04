package kv

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/models"
)

var ExportJson = []byte(`[
  {
    "CreateIndex": 4398,
    "ModifyIndex": 4398,
    "LockIndex": 0,
    "Key": "dcdr/features/cc/b",
    "Flags": 0,
    "Value": "eyJmZWF0dXJlX3R5cGUiOiJwZXJjZW50aWxlIiwia2V5IjoiYiIsIm5hbWVzcGFjZSI6ImRjZHIvZmVhdHVyZXMiLCJzY29wZSI6ImNjIiwidmFsdWUiOjEsImNvbW1lbnQiOiIiLCJ1cGRhdGVkX2J5IjoiY2hyaXNiIn0="
  },
  {
    "CreateIndex": 4400,
    "ModifyIndex": 4400,
    "LockIndex": 0,
    "Key": "dcdr/features/default/a",
    "Flags": 0,
    "Value": "eyJmZWF0dXJlX3R5cGUiOiJwZXJjZW50aWxlIiwia2V5IjoiYSIsIm5hbWVzcGFjZSI6ImRjZHIvZmVhdHVyZXMiLCJzY29wZSI6ImRlZmF1bHQiLCJ2YWx1ZSI6MSwiY29tbWVudCI6IiIsInVwZGF0ZWRfYnkiOiJjaHJpc2IifQ=="
  },
  {
    "CreateIndex": 4399,
    "ModifyIndex": 4401,
    "LockIndex": 0,
    "Key": "dcdr/info",
    "Flags": 0,
    "Value": "eyJjdXJyZW50X3NoYSI6IjQzZDRiOWU3ZGU4ZWQ0OGE5NTRmMzU5NGU2YmQzM2U2ZDY5YjU1MTYifQ=="
  }
]
`)

type MockStore struct {
	Item  []byte
	Items [][]byte
	Err   error
}

func NewMockStore(ft *models.Feature, err error) (ms *MockStore) {
	bts, _ := ft.ToJson()

	ms = &MockStore{
		Item:  bts,
		Items: [][]byte{bts},
		Err:   err,
	}

	return
}

func (ms *MockStore) List(prefix string) ([][]byte, error) {
	return ms.Items, ms.Err
}

func (ms *MockStore) Get(key string) ([]byte, error) {
	return ms.Item, ms.Err
}

func (ms *MockStore) Set(key string, bts []byte) error {
	return ms.Err
}

func (ms *MockStore) Delete(key string) error {
	return ms.Err
}

func (ms *MockStore) Put(key string, bts []byte) error {
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

func (mr *MockRepo) Push() error {
	return mr.error
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
	sr := &models.Feature{
		Key:       "key",
		Scope:     "scope",
		Namespace: "namespace",
		Value:     0.5,
		Comment:   "comment",
		UpdatedBy: "user",
	}

	c := New(&MockStore{}, &MockRepo{}, sr.Namespace, nil)

	err := c.Set(sr)

	assert.NoError(t, err)
}

func TestClientSetExisting(t *testing.T) {
	sr := &models.Feature{
		Key:         "key",
		Scope:       "scope",
		Namespace:   "namespace",
		UpdatedBy:   "user",
		FeatureType: models.GetFeatureType("percentile"),
	}

	orig := &models.Feature{
		Key:         "key",
		Scope:       "scope",
		Namespace:   "namespace",
		Value:       0.5,
		Comment:     "comment",
		UpdatedBy:   "user",
		FeatureType: models.GetFeatureType("percentile"),
	}

	c := New(NewMockStore(orig, nil), &MockRepo{}, sr.Namespace, nil)

	err := c.Set(sr)

	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, "", nil)

	fts, err := c.List("test", "")

	assert.Nil(t, err)
	assert.Equal(t, models.Features{*ft}, fts)
}

func TestGet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, "", nil)

	var f *models.Feature
	err := c.Get("test", &f)

	assert.Nil(t, err)
	assert.Equal(t, f, ft)
}

func TestNilGet(t *testing.T) {
	cs := NewMockStore(nil, nil)
	c := New(cs, &MockRepo{}, "", nil)

	var f *models.Feature
	err := c.Get("test", &f)

	assert.Nil(t, err)
	assert.Nil(t, f)
}

func TestSet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, "", nil)

	err := c.Set(ft)

	assert.Nil(t, err)
}

func TestTypeChangeErrorSet(t *testing.T) {
	orig := models.NewFeature("test", 0.5, "c", "u", "s")
	bad := models.NewFeature("test", false, "c", "u", "s")

	cs := NewMockStore(orig, nil)
	c := New(cs, nil, "", nil)

	err := c.Set(bad)
	assert.Equal(t, TypeChangeError, err)
}

func TestSetWithError(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	e := errors.New("")
	cs := NewMockStore(ft, e)
	c := New(cs, nil, "", nil)

	err := c.Set(ft)

	assert.Equal(t, e, err)
}

func TestDelete(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, "", nil)

	err := c.Delete(ft.Key, "")

	assert.Nil(t, err)
}

func TestDeleteWithError(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	e := errors.New("")
	cs := NewMockStore(ft, e)
	c := New(cs, &MockRepo{}, "", nil)

	err := c.Delete(ft.Key, "")

	assert.Equal(t, e, err)
}

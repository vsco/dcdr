package api

import (
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/cli/models"
)

type MockStore struct {
	Item  *stores.KVByte
	Items stores.KVBytes
	Err   error
}

func NewMockStore(ft *models.Feature, err error) (ms *MockStore) {
	bts, _ := ft.ToJson()

	ms = &MockStore{
		Err: err,
	}

	if ft != nil {
		kvb := stores.KVBytes{
			&stores.KVByte{
				Key:   ft.Key,
				Bytes: bts,
			},
		}

		ms.Item = kvb[0]
		ms.Items = kvb
	}

	return
}

func (ms *MockStore) List(prefix string) (stores.KVBytes, error) {
	return ms.Items, ms.Err
}

func (ms *MockStore) Get(key string) (*stores.KVByte, error) {
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

func (mr *MockRepo) Commit(bts []byte, msg string) error {
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

func TestClientSet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")

	c := New(NewMockStore(ft, nil), &MockRepo{}, ft.GetNamespace(), nil)

	err := c.Set(ft)

	assert.NoError(t, err)
}

func TestClientSetExisting(t *testing.T) {
	update := models.NewFeature("test", nil, "c", "u", "s", "n")
	orig := models.NewFeature("test", 0.5, "c", "u", "s", "n")

	c := New(NewMockStore(orig, nil), &MockRepo{}, update.GetNamespace(), nil)

	err := c.Set(update)

	assert.NoError(t, err)
}

func TestList(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, "", nil)

	fts, err := c.List("test", "")

	assert.Nil(t, err)
	assert.Equal(t, models.Features{*ft}, fts)
}

func TestGet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
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

	assert.EqualError(t, err, "/test not found")
	assert.Nil(t, f)
}

func TestSet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, "", nil)

	err := c.Set(ft)

	assert.Nil(t, err)
}

func TestTypeChangeErrorSet(t *testing.T) {
	orig := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	bad := models.NewFeature("test", false, "c", "u", "s", "n")

	cs := NewMockStore(orig, nil)
	c := New(cs, nil, "", nil)

	err := c.Set(bad)
	assert.Equal(t, TypeChangeError, err)
}

func TestSetWithError(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	e := errors.New("")
	cs := NewMockStore(ft, e)
	c := New(cs, nil, "", nil)

	err := c.Set(ft)

	assert.Equal(t, e, err)
}

func TestDelete(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	cs := NewMockStore(ft, nil)
	c := New(cs, &MockRepo{}, "", nil)

	err := c.Delete(ft.Key, "")

	assert.Nil(t, err)
}

func TestDeleteWithError(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s", "n")
	e := errors.New("")
	cs := NewMockStore(ft, e)
	c := New(cs, &MockRepo{}, "", nil)

	err := c.Delete(ft.Key, "")

	assert.Equal(t, e, err)
}

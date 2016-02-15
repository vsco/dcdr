package kv

import (
	"errors"
	"testing"

	"encoding/json"

	"github.com/hashicorp/consul/api"
	"github.com/stretchr/testify/assert"
	"github.com/vsco/dcdr/models"
)

type MockConsul struct {
	Item  *models.Feature
	Items api.KVPairs
	Err   error
}

func MockConsulStore(ft *models.Feature, err error) (cs *ConsulStore) {
	mc := &MockConsul{
		Item: ft,
		Err:  err,
	}

	if ft != nil {
		bts, _ := json.Marshal(mc.Item)
		mc.Items = api.KVPairs{
			{
				Key:   ft.Key,
				Value: bts,
			},
		}
	}

	cs = &ConsulStore{
		kv: mc,
	}

	return
}

func (mc *MockConsul) get(key string) *api.KVPair {
	bts, _ := json.Marshal(mc.Item)

	return &api.KVPair{
		Key:   key,
		Value: bts,
	}
}

func (mc *MockConsul) List(prefix string, qo *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error) {
	return mc.Items, nil, mc.Err
}

func (mc *MockConsul) Get(key string, qo *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error) {
	return mc.get(key), nil, mc.Err
}

func (mc *MockConsul) Put(p *api.KVPair, qo *api.WriteOptions) (*api.WriteMeta, error) {
	return nil, mc.Err
}

func (mc *MockConsul) Delete(key string, w *api.WriteOptions) (*api.WriteMeta, error) {
	return nil, mc.Err
}

func TestList(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	cs := MockConsulStore(ft, nil)

	fts, err := cs.List("test")

	assert.Nil(t, err)
	assert.Equal(t, models.Features{*ft}, fts)
}

func TestGet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	cs := MockConsulStore(ft, nil)

	f, err := cs.Get("test")

	assert.Nil(t, err)
	assert.Equal(t, f, ft)
}

func TestNilGet(t *testing.T) {
	cs := MockConsulStore(nil, nil)

	f, err := cs.Get("test")

	assert.Nil(t, err)
	assert.Nil(t, f)
}

func TestSet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	cs := MockConsulStore(ft, nil)

	err := cs.Set(ft)

	assert.Nil(t, err)
}

func TestSetWithError(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	e := errors.New("")
	cs := MockConsulStore(ft, e)

	err := cs.Set(ft)

	assert.Equal(t, e, err)
}

func TestTypeChangeErrorSet(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	bt := models.NewFeature("test", false, "c", "u", "s")

	cs := MockConsulStore(ft, nil)

	err := cs.Set(bt)
	assert.Equal(t, TypeChangeError, err)
}

func TestDelete(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	cs := MockConsulStore(ft, nil)

	err := cs.Delete(ft.Key)

	assert.Nil(t, err)
}

func TestDeleteWithError(t *testing.T) {
	ft := models.NewFeature("test", 0.5, "c", "u", "s")
	e := errors.New("")
	cs := MockConsulStore(ft, e)

	err := cs.Delete(ft.Key)

	assert.Equal(t, e, err)
}

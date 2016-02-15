package kv

import (
	"encoding/json"

	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/vsco/dcdr/models"
)

type ConsulKVIFace interface {
	List(prefix string, q *api.QueryOptions) (api.KVPairs, *api.QueryMeta, error)
	Get(key string, q *api.QueryOptions) (*api.KVPair, *api.QueryMeta, error)
	Put(p *api.KVPair, q *api.WriteOptions) (*api.WriteMeta, error)
	Delete(key string, w *api.WriteOptions) (*api.WriteMeta, error)
}

type ConsulStore struct {
	kv ConsulKVIFace
	qo *api.QueryOptions
	wo *api.WriteOptions
}

func DefaultConsulStore() (StoreIFace, error) {
	client, err := api.NewClient(api.DefaultConfig())

	if err != nil {
		return nil, err
	}

	return &ConsulStore{
		kv: client.KV(),
		qo: nil,
		wo: nil,
	}, nil
}

func (cs *ConsulStore) Get(key string) (*models.Feature, error) {
	kv, _, err := cs.kv.Get(key, cs.qo)

	if err != nil || kv == nil {
		return nil, err
	}

	var f *models.Feature

	err = json.Unmarshal(kv.Value, &f)

	if err != nil {
		return f, err
	}

	return f, nil
}

func (cs *ConsulStore) Set(f *models.Feature) error {
	org, err := cs.Get(f.ScopedKey())

	if err != nil {
		return err
	}

	if org != nil {
		if f.FeatureType != org.FeatureType {
			return TypeChangeError
		}

		if f.Comment == "" {
			f.Comment = org.Comment
		}
	}

	bts, err := json.Marshal(f)

	if err != nil {
		return err
	}

	p := &api.KVPair{
		Key:   f.ScopedKey(),
		Value: bts,
	}

	_, err = cs.kv.Put(p, cs.wo)

	return err
}

func (cs *ConsulStore) Delete(key string) error {
	_, err := cs.kv.Delete(key, cs.wo)

	return err
}

func (cs *ConsulStore) List(prefix string) (models.Features, error) {
	fmt.Printf("listing:  %s", prefix)
	kvs, _, err := cs.kv.List(prefix, cs.qo)

	var fts models.Features

	if err != nil {
		return fts, err
	}

	for _, v := range kvs {
		var f models.Feature

		err := json.Unmarshal(v.Value, &f)

		if err != nil {
			return fts, err
		}

		fts = append(fts, f)
	}

	return fts, err
}

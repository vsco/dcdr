package kv

import (
	"log"

	"encoding/json"

	"fmt"

	"github.com/hashicorp/consul/api"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

const CurrentShaKey = "current_sha"

type Client struct {
	consul    *api.Client
	Namespace string
}

func New(cfg *api.Config, namespace string) (c *Client) {
	cl, err := api.NewClient(cfg)

	if err != nil {
		log.Fatal(err)
	}

	c = &Client{
		consul:    cl,
		Namespace: namespace,
	}

	return
}

func (c *Client) List(prefix string) (models.Features, error) {
	key := fmt.Sprintf("%s/%s", c.Namespace, prefix)
	kvs, _, err := c.consul.KV().List(key, nil)
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

func (c *Client) SetPercentile(k string, v float64, cmt string, user string) {
	f := models.PercentileFeature(k, v, cmt, user)
	c.set(f)
}

func (c *Client) SetBoolean(k string, v bool, cmt string, user string) {
	f := models.BooleanFeature(k, v, cmt, user)

	c.set(f)
}

func (c *Client) Delete(key string) error {
	err := c.delete(key)

	return err
}

func (c *Client) SetCurrentSha(sha string) error {
	return c.putWithNamespace(CurrentShaKey, config.DefaultInfoNamespace, []byte(sha))
}

func (c *Client) Get(key string) (*models.Feature, error) {
	kv, err := c.get(key)

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

func (c *Client) get(key string) (*api.KVPair, error) {
	kv, _, err := c.consul.KV().Get(fmt.Sprintf("%s/%s", c.Namespace, key), nil)

	return kv, err
}

func (c *Client) set(f *models.Feature) {
	kv, err := c.get(f.Name)

	if err != nil {
		log.Fatal(err)
	}

	if kv != nil {
		var existing models.Feature
		err := json.Unmarshal(kv.Value, &existing)

		if err != nil {
			log.Fatal(err)
		}

		if f.FeatureType != existing.FeatureType {
			log.Fatal("cannot change existing feature types.")
		}

		if f.Comment == "" {
			f.Comment = existing.Comment
		}
	}

	bts, err := json.Marshal(f)

	if err != nil {
		log.Fatal(err)
	}

	err = c.put(f.Name, bts)

	if err != nil {
		log.Fatal(err)
	}
}

func (c *Client) putWithNamespace(key string, ns string, bts []byte) error {
	p := &api.KVPair{Key: fmt.Sprintf("%s/%s", ns, key), Value: bts}
	_, err := c.consul.KV().Put(p, nil)

	return err
}

func (c *Client) put(key string, bts []byte) error {
	return c.putWithNamespace(key, c.Namespace, bts)
}

func (c *Client) delete(key string) error {
	_, err := c.consul.KV().Delete(fmt.Sprintf("%s/%s", c.Namespace, key), nil)

	return err
}

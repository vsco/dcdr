package client

import (
	"log"
	"sort"

	"encoding/json"

	"github.com/hashicorp/consul/api"
	"github.com/vsco/decider-cli/models"
)

type Client struct {
	consul    *api.Client
	namespace string
}

type Features []models.Feature

func New(cfg *api.Config, namespace string) (c *Client) {
	cl, err := api.NewClient(cfg)

	if err != nil {
		log.Fatal(err)
	}

	c = &Client{
		consul:    cl,
		namespace: namespace,
	}

	return
}

func (c *Client) List(prefix string) (Features, error) {
	kvs, _, err := c.consul.KV().List(c.namespace+"/"+prefix, nil)
	var fts Features

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

	sort.Sort(models.ByName(fts))

	return fts, err
}

func (c *Client) SetPercentile(k string, v float64, cmt string) {
	f := models.PercentileFeature(k, v, cmt)
	c.set(f)
}

func (c *Client) SetBoolean(k string, v bool, cmt string) {
	f := models.BooleanFeature(k, v, cmt)

	c.set(f)
}

func (c *Client) SetScalar(k string, v float64, cmt string) {
	f := models.ScalarFeature(k, v, cmt)

	c.set(f)
}

func (c *Client) Delete(k string) error {
	_, err := c.consul.KV().Delete(c.namespace+"/"+k, nil)

	return err
}

func (c *Client) set(f *models.Feature) {
	kv, _, err := c.consul.KV().Get(c.namespace+"/"+f.Name, nil)

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
			log.Fatal("cannot change feature types.")
		}

		if f.Comment == "" {
			f.Comment = existing.Comment
		}
	}

	bts, err := json.Marshal(f)

	if err != nil {
		log.Fatal(err)
	}

	p := &api.KVPair{Key: c.namespace + "/" + f.Name, Value: bts}
	_, err = c.consul.KV().Put(p, nil)

	if err != nil {
		log.Fatal(err)
	}
}

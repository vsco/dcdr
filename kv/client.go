package kv

import (
	"errors"
	"fmt"

	"github.com/vsco/dcdr/models"
)

const (
	CurrentShaKey = "current_sha"
)

func ValidationError(n string) error {
	return errors.New(fmt.Sprintf("%s is required"))
}

type ClientIFace interface {
	List(prefix string, scope string) (models.Features, error)
	Set(sr *SetRequest) error
	Delete(key string, scope string) error
	Namespace() string
}

type Client struct {
	Store     StoreIFace
	namespace string
}

type SetRequest struct {
	Key       string
	Value     interface{}
	Scope     string
	Namespace string
	Comment   string
	User      string
}

func (sr *SetRequest) ToFeature() (*models.Feature, error) {
	var ft models.FeatureType

	if sr.Key == "" {
		return nil, ValidationError("name")
	}

	switch sr.Value.(type) {
	case bool:
		ft = models.Boolean
	default:
		ft = models.Percentile
	}

	return &models.Feature{
		Key:         sr.Key,
		Scope:       sr.Scope,
		Namespace:   sr.Namespace,
		Value:       sr.Value,
		Comment:     sr.Comment,
		UpdatedBy:   sr.User,
		FeatureType: ft,
	}, nil
}

func New(st StoreIFace, namespace string) (c *Client) {
	c = &Client{
		Store:     st,
		namespace: namespace,
	}

	return
}
func (c *Client) Namespace() string {
	return c.namespace
}

func (c *Client) List(prefix string, scope string) (models.Features, error) {
	prefix = fmt.Sprintf("%s/%s/%s", c.Namespace(), scope, prefix)
	fts, err := c.Store.List(prefix)

	return fts, err
}

func (c *Client) Set(sr *SetRequest) error {
	ft, err := sr.ToFeature()

	if err != nil {
		return err
	}

	existing, err := c.Store.Get(ft.ScopedKey())

	if existing != nil {
		if ft.Comment == "" {
			ft.Comment = existing.Comment
		}
		if ft.Value == nil {
			ft.Value = existing.Value
		}
	}

	return c.Store.Set(ft)
}

func (c *Client) Delete(key string, scope string) error {
	key = fmt.Sprintf("%s/%s/%s", c.Namespace(), scope, key)

	err := c.Store.Delete(key)

	return err
}

//func (c *Client) SetCurrentSha(sha string) error {
//	return c.putWithNamespace(CurrentShaKey, config.DefaultInfoNamespace, []byte(sha))
//}

//func (c *Client) Get(key string, scope string) (*models.Feature, error) {
//	kv, err := c.get(key)
//
//	if err != nil || kv == nil {
//		return nil, err
//	}
//
//	var f *models.Feature
//
//	err = json.Unmarshal(kv.Value, &f)
//
//	if err != nil {
//		return f, err
//	}
//
//	return f, nil
//}

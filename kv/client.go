package kv

import (
	"errors"
	"fmt"

	"github.com/PagerDuty/godspeed"
	"github.com/vsco/dcdr/models"
	"github.com/vsco/dcdr/repo"
)

const (
	CurrentShaKey = "current_sha"
)

func ValidationError(n string) error {
	return errors.New(fmt.Sprintf("%s is required", n))
}

func KeyNotFoundError(n string) error {
	return errors.New(fmt.Sprintf("%s not found", n))
}

type ClientIFace interface {
	List(prefix string, scope string) (models.Features, error)
	Set(sr *SetRequest) error
	Delete(key string, scope string) error
	InitRepo(create bool) error
	Namespace() string
}

type Client struct {
	Store     StoreIFace
	Repo      repo.RepoIFace
	Stats     *godspeed.Godspeed
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

func New(st StoreIFace, rp repo.RepoIFace, namespace string, stats *godspeed.Godspeed) (c *Client) {
	c = &Client{
		Store:     st,
		Repo:      rp,
		Stats:     stats,
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

	err = c.Store.Set(ft)

	if err != nil {
		return err
	}

	if c.Repo.Enabled() {
		err := c.CommitFeatures(ft, false)

		if err != nil {
			return err
		}
	}

	err = c.SendStatEvent(ft, false)

	return nil
}

func (c *Client) SendStatEvent(f *models.Feature, delete bool) error {
	if c.Stats == nil {
		return nil
	}

	var text string
	title := "Decider Change"

	if delete {
		text = fmt.Sprintf("deleted %s", f.ScopedKey(), f.Value)
	} else {
		text = fmt.Sprintf("set %s: %v", f.ScopedKey(), f.Value)
	}

	optionals := make(map[string]string)
	optionals["alert_type"] = "info"
	optionals["source_type_name"] = "dcdr"
	tags := []string{"source_type:dcdr"}

	return c.Stats.Event(title, text, optionals, tags)
}

func (c *Client) Delete(key string, scope string) error {
	key = fmt.Sprintf("%s/%s/%s", c.Namespace(), scope, key)

	existing, err := c.Store.Get(key)

	if err != nil {
		return err
	}

	if existing == nil {
		return KeyNotFoundError(key)
	}

	err = c.Store.Delete(key)

	if err != nil {
		return err
	}

	if c.Repo.Exists() {
		err := c.CommitFeatures(existing, true)

		if err != nil {
			return err
		}
	}

	err = c.SendStatEvent(existing, true)

	return err
}

func (c *Client) CommitFeatures(ft *models.Feature, deleted bool) error {
	fts, err := c.List("", "")

	if err != nil {
		return err
	}

	var msg string

	if deleted {
		msg = fmt.Sprintf("%s deleted %s", ft.UpdatedBy, ft.ScopedKey())
	} else {
		msg = fmt.Sprintf("%s set %s to %v", ft.UpdatedBy, ft.ScopedKey(), ft.Value)
	}

	return c.Repo.Commit(fts, msg)
}

func (c *Client) InitRepo(create bool) error {
	if create {
		return c.Repo.Create()
	}

	return c.Repo.Clone()
}

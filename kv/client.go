package kv

import (
	"errors"
	"fmt"

	"encoding/json"

	"github.com/PagerDuty/godspeed"
	"github.com/vsco/dcdr/kv/stores"
	"github.com/vsco/dcdr/models"
	"github.com/vsco/dcdr/repo"
)

const (
	CurrentShaKey = "current_sha"
	InfoNameSpace = "info"
)

var TypeChangeError = errors.New("cannot change existing feature types.")

func ValidationError(n string) error {
	return errors.New(fmt.Sprintf("%s is required", n))
}

func KeyNotFoundError(n string) error {
	return errors.New(fmt.Sprintf("%s not found", n))
}

type ClientIFace interface {
	List(prefix string, scope string) (models.Features, error)
	Set(ft *models.Feature) error
	Get(key string, v interface{}) error
	Delete(key string, scope string) error
	GetInfo() (*models.Info, error)
	InitRepo(create bool) error
	Namespace() string
}

type Client struct {
	Store     stores.StoreIFace
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

func New(st stores.StoreIFace, rp repo.RepoIFace, namespace string, stats *godspeed.Godspeed) (c *Client) {
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
	res, err := c.Store.List(prefix)
	fts := make(models.Features, len(res))

	if err != nil {
		return fts, err
	}

	for i := 0; i < len(res); i++ {
		var f models.Feature
		err := json.Unmarshal(res[i].Bytes, &f)

		if err != nil {
			return fts, err
		}

		fts[i] = f
	}

	return fts, nil
}

func (c *Client) Set(ft *models.Feature) error {
	var existing *models.Feature

	kvb, err := c.Store.Get(ft.ScopedKey())

	if err != nil {
		return err
	}

	if kvb != nil {
		err = json.Unmarshal(kvb.Bytes, &existing)

		if err != nil {
			return err
		}
	}

	if err != nil {
		return err
	}

	if existing != nil {
		if ft.Comment == "" {
			ft.Comment = existing.Comment
		}
		if ft.Value == nil {
			ft.Value = existing.Value
		}
		if ft.FeatureType != existing.FeatureType && ft.FeatureType != "" {
			return TypeChangeError
		}
		if ft.FeatureType == "" {
			ft.FeatureType = existing.FeatureType
		}
	}

	bts, err := ft.ToJson()

	if err != nil {
		return err
	}

	err = c.Store.Put(ft.ScopedKey(), bts)

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

func (c *Client) Get(key string, v interface{}) error {
	key = fmt.Sprintf("%s/%s", c.namespace, key)

	bts, err := c.Store.Get(key)

	if err != nil {
		return err
	}

	if bts == nil {
		return KeyNotFoundError(key)
	}

	return json.Unmarshal(bts.Bytes, &v)
}

func (c *Client) SendStatEvent(f *models.Feature, delete bool) error {
	if c.Stats == nil {
		return nil
	}

	var text string
	title := "Decider Change"

	if delete {
		text = fmt.Sprintf("deleted %s", f.ScopedKey())
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
	var existing *models.Feature

	key = fmt.Sprintf("%s/%s/%s", c.Namespace(), scope, key)
	kv, err := c.Store.Get(key)

	if err != nil {
		return err
	}

	err = json.Unmarshal(kv.Bytes, &existing)

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
	kvb, err := c.Store.List("dcdr/features")

	if err != nil {
		return err
	}

	fm, err := models.KVsToFeatureMap(kvb)

	if err != nil {
		return err
	}

	bts, err := json.MarshalIndent(fm, "", "  ")

	if err != nil {
		return err
	}

	var msg string

	if deleted {
		msg = fmt.Sprintf("%s deleted %s", ft.UpdatedBy, ft.ScopedKey())
	} else {
		msg = fmt.Sprintf("%s set %s to %v", ft.UpdatedBy, ft.ScopedKey(), ft.Value)
	}

	fmt.Println("[dcdr] commiting changes")
	err = c.Repo.Commit(bts, msg)

	if err != nil {
		return err
	}

	sha, err := c.Repo.CurrentSha()

	if err != nil {
		return err
	}

	err = c.SetCurrentSha(sha)

	if err != nil {
		return err
	}

	fmt.Println("[dcdr] pushing commit to origin")
	err = c.Repo.Push()

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) GetInfo() (*models.Info, error) {
	key := fmt.Sprintf("dcdr/%s", InfoNameSpace)

	var info *models.Info

	kv, err := c.Store.Get(key)

	if err != nil {
		return nil, err
	}

	if len(kv.Bytes) == 0 {
		return &models.Info{}, nil
	}

	err = json.Unmarshal(kv.Bytes, &info)

	if err != nil {
		return nil, err
	}

	return info, err
}

func (c *Client) SetCurrentSha(sha string) error {
	key := fmt.Sprintf("dcdr/%s", InfoNameSpace)

	info := &models.Info{
		CurrentSha: sha,
	}

	bts, err := json.Marshal(info)

	if err != nil {
		return err
	}

	return c.Store.Put(key, bts)
}

func (c *Client) InitRepo(create bool) error {
	if create {
		return c.Repo.Create()
	}

	return c.Repo.Clone()
}

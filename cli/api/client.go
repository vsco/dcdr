package api

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"encoding/json"

	"time"

	"github.com/PagerDuty/godspeed"
	"github.com/vsco/dcdr/cli/api/ioutil2"
	"github.com/vsco/dcdr/cli/api/stores"
	"github.com/vsco/dcdr/cli/printer"
	"github.com/vsco/dcdr/cli/repo"
	"github.com/vsco/dcdr/config"
	"github.com/vsco/dcdr/models"
)

const InfoNameSpace = "info"

var ErrTypeChange = errors.New("cannot change existing feature types")
var ErrRepoExists = errors.New("repository already exists")
var ErrNilValue = errors.New("value cannot be nil")

func KeyNotFoundError(n string) error {
	return fmt.Errorf("%s not found", n)
}

type ClientIFace interface {
	List(prefix string, scope string) (models.Features, error)
	Set(ft *models.Feature) error
	Get(key string, v interface{}) error
	Delete(key string, scope string) error
	GetInfo() (*models.Info, error)
	InitRepo(create bool) error
	Commit(ft *models.Feature, deleted bool) error
	Push() error
	UpdateCurrentSHA() (string, error)
	Watch()
	Namespace() string
}

type Client struct {
	Store  stores.IFace
	Repo   repo.IFace
	Stats  *godspeed.Godspeed
	config *config.Config
}

func New(st stores.IFace, rp repo.IFace, cfg *config.Config, stats *godspeed.Godspeed) (c *Client) {
	c = &Client{
		Store:  st,
		Repo:   rp,
		Stats:  stats,
		config: cfg,
	}

	return
}
func (c *Client) Namespace() string {
	return c.config.Namespace
}

func (c *Client) List(prefix string, scope string) (models.Features, error) {
	defer c.Store.Close()

	if prefix == "" {
		prefix = fmt.Sprintf("%s/features/%s", c.Namespace(), scope)
	} else {
		prefix = fmt.Sprintf("%s/features/%s/%s", c.Namespace(), scope, prefix)
	}

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
	defer c.Store.Close()

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
			return ErrTypeChange
		}
		if ft.FeatureType == "" {
			ft.FeatureType = existing.FeatureType
		}
	} else {
		if ft.Value == nil {
			return ErrNilValue
		}
	}

	bts, err := ft.ToJSON()

	if err != nil {
		return err
	}

	err = c.Store.Set(ft.ScopedKey(), bts)

	if err != nil {
		return err
	}

	err = c.SendStatEvent(ft, false)

	return nil
}

func (c *Client) Get(key string, v interface{}) error {
	defer c.Store.Close()

	key = fmt.Sprintf("%s/%s", c.Namespace(), key)

	bts, err := c.Store.Get(key)

	if err != nil {
		return err
	}

	if bts == nil {
		return KeyNotFoundError(key)
	}

	return json.Unmarshal(bts.Bytes, &v)
}

func (c *Client) Delete(key string, scope string) error {
	defer c.Store.Close()

	var existing *models.Feature

	key = fmt.Sprintf("%s/features/%s/%s", c.Namespace(), scope, key)
	kv, err := c.Store.Get(key)

	if err != nil {
		return err
	}

	if kv != nil {
		err = json.Unmarshal(kv.Bytes, &existing)

		if err != nil {
			return err
		}
	}

	if existing == nil {
		return KeyNotFoundError(key)
	}

	err = c.Store.Delete(key)

	if err != nil {
		return err
	}

	err = c.SendStatEvent(existing, true)

	return err
}

func (c *Client) Commit(ft *models.Feature, deleted bool) error {
	if !c.Repo.Exists() {
		err := c.Repo.Clone()

		if err != nil {
			return err
		}
	}

	kvb, err := c.Store.List(fmt.Sprintf("%s/features", c.Namespace()))

	if err != nil {
		return err
	}

	fm, err := c.KVsToFeatureMap(kvb)

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

	err = c.Repo.Commit(bts, msg)

	if err != nil {
		return err
	}

	return nil
}

func (c *Client) Push() error {
	return c.Repo.Push()
}

func (c *Client) GetInfo() (*models.Info, error) {
	defer c.Store.Close()

	key := fmt.Sprintf("%s/%s", c.Namespace(), InfoNameSpace)

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

func (c *Client) UpdateCurrentSHA() (string, error) {
	defer c.Store.Close()

	sha, err := c.Repo.CurrentSHA()

	if err != nil {
		return sha, err
	}

	key := fmt.Sprintf("%s/%s", c.Namespace(), InfoNameSpace)

	info := &models.Info{
		CurrentSHA:       sha,
		LastModifiedDate: time.Now().UTC().Unix(),
	}

	bts, err := json.Marshal(info)

	if err != nil {
		return sha, err
	}

	return sha, c.Store.Set(key, bts)
}

func (c *Client) InitRepo(create bool) error {
	if c.Repo.Exists() {
		return ErrRepoExists
	}

	if create {
		return c.Repo.Create()
	}

	return c.Repo.Clone()
}

func (c *Client) Watch() {
	c.Store.Register(c.WriteOutputFile)
	c.Store.Watch()
}

func (c *Client) WriteOutputFile(kvb stores.KVBytes) {
	fts, err := c.KVsToFeatureMap(kvb)

	if err != nil {
		printer.LogErrf("parse features error: %v", err)
		os.Exit(1)
	}

	bts, err := json.MarshalIndent(fts, "", "  ")

	if err != nil {
		printer.LogErrf("%v", err)
		os.Exit(1)
	}

	err = ioutil2.WriteFileAtomic(c.config.Watcher.OutputPath, bts, 0644)

	if err != nil {
		printer.LogErrf("%v", err)
		os.Exit(1)
	}

	printer.Logf("wrote changes to: %s", c.config.Watcher.OutputPath)
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

// KVsToFeatures helper for unmarshalling `KVBytes` to a `FeatureMap`
func (c *Client) KVsToFeatureMap(kvb stores.KVBytes) (*models.FeatureMap, error) {
	fm := models.EmptyFeatureMap()

	for _, v := range kvb {
		var key string
		var value interface{}

		if v.Key == config.DefaultInfoNamespace {
			var info models.Info
			err := json.Unmarshal(v.Bytes, &info)

			if err != nil {
				return fm, err
			}

			fm.Dcdr.Info = &info
		} else {
			var ft models.Feature
			err := json.Unmarshal(v.Bytes, &ft)

			if err != nil {
				printer.SayErr("%s: %s", v.Key, v.Bytes)
				return fm, err
			}

			key = strings.Replace(v.Key, fmt.Sprintf("%s/features/", c.Namespace()), "", 1)
			value = ft.Value
		}

		explode(fm.Dcdr.FeatureScopes, key, value)
	}

	return fm, nil
}

func explode(m models.FeatureScopes, k string, v interface{}) {
	if strings.Contains(k, "/") {
		pts := strings.Split(k, "/")
		top := pts[0]
		key := strings.Join(pts[1:], "/")

		if _, ok := m[top]; !ok {
			m[top] = make(map[string]interface{})
		}

		explode(m[top].(map[string]interface{}), key, v)
	} else {
		if k != "" {
			m[k] = v
		}
	}
}

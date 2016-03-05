package models

import (
	"encoding/json"
	"strings"
	"sync"
)

const DefaultScope = "default"

type Info struct {
	CurrentSha string `json:"current_sha"`
}

type FeatureMap struct {
	Dcdr Root `json:"dcdr"`
}

type Features map[string]interface{}

type Root struct {
	sync.RWMutex
	Info     Info     `json:"info"`
	Features Features `json:"features"`
}

func NewFeatureMap(bts []byte) (*FeatureMap, error) {
	var fm *FeatureMap
	err := json.Unmarshal(bts, &fm)

	if err != nil {
		return nil, err
	}

	return fm, nil
}

func (fm *FeatureMap) ToJson() ([]byte, error) {
	bts, err := json.MarshalIndent(fm, "", "  ")

	if err != nil {
		return bts, err
	}

	return bts, nil
}

func (d *Root) InScope(scope string) Features {
	scopes := strings.Split(scope, "/")

	d.RLock()
	defer d.RUnlock()

	top := d.Features
	for _, s := range scopes {
		if m, ok := top[s]; ok {
			top = m.(map[string]interface{})
		} else {
			return make(map[string]interface{})
		}
	}

	return top
}

func (d *Root) Defaults() Features {
	return d.InScope(DefaultScope)
}

func rev(a []string) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

func (d *Root) MergedScopes(scopes ...string) Features {
	scopes = append(scopes, DefaultScope)
	mrg := make(Features)

	rev(scopes)
	for _, scope := range scopes {
		fts := d.InScope(scope)

		for k, v := range fts {
			mrg[k] = v
		}
	}

	return mrg
}

package models

import (
	"encoding/json"
	"strings"
	"sync"
)

// DefaultScope the key in which default features are set.
// By default `dcdr/features/<DefaultScope>`/<feature-name>.
const DefaultScope = "default"

// FeatureMap contains a nested map of `Features` and `Info`.
type FeatureMap struct {
	Dcdr Root `json:"dcdr"`
}

// Info contains the metadata for the current `FeatureMap`.
type Info struct {
	CurrentSHA       string `json:"current_sha,omitempty"`
	LastModifiedDate int64  `json:"last_modified_date,omitempty"`
}

// FeatureScopes the map of percentile and boolean K/Vs.
type FeatureScopes map[string]interface{}

// Root wrapper struct for `Info` and `Features`.
type Root struct {
	sync.RWMutex
	Info          *Info         `json:"info"`
	FeatureScopes FeatureScopes `json:"features"`
}

// EmptyFeatureMap helper method for constructing an empty `FeatureMap`.
func EmptyFeatureMap() (fm *FeatureMap) {
	fm = &FeatureMap{
		Dcdr: Root{
			Info: &Info{
				CurrentSHA: "",
			},
			FeatureScopes: FeatureScopes{
				DefaultScope: make(map[string]interface{}),
			},
		},
	}

	return
}

// NewFeatureMap constructs a `FeatureMap` for marshalled JSON.
func NewFeatureMap(bts []byte) (*FeatureMap, error) {
	var fm *FeatureMap
	err := json.Unmarshal(bts, &fm)

	if err != nil {
		return nil, err
	}

	return fm, nil
}

// ToJSON helper method for marshalling a `FeatureMap` to JSON.
func (fm *FeatureMap) ToJSON() ([]byte, error) {
	bts, err := json.MarshalIndent(fm, "", "  ")

	if err != nil {
		return bts, err
	}

	return bts, nil
}

// InScope returns the `Features` found within `scope`.
func (d *Root) InScope(scope string) FeatureScopes {
	scopes := strings.Split(scope, "/")

	d.RLock()
	defer d.RUnlock()

	top := d.FeatureScopes
	for _, s := range scopes {
		if m, ok := top[s]; ok {
			top = m.(map[string]interface{})
		} else {
			return make(map[string]interface{})
		}
	}

	return top
}

// Defaults returns `Features` within the 'default' scope.
func (d *Root) Defaults() FeatureScopes {
	return d.InScope(DefaultScope)
}

// MergedScopes given a slice of scopes in priority order will return a
// merged set of `Featured` including the 'default' scope.
func (d *Root) MergedScopes(scopes ...string) FeatureScopes {
	scopes = append(scopes, DefaultScope)
	mrg := make(FeatureScopes)

	rev(scopes)
	for _, scope := range scopes {
		if scope != "" {
			fts := d.InScope(scope)

			for k, v := range fts {
				mrg[k] = v
			}
		}
	}

	return mrg
}

func rev(a []string) {
	for i := len(a)/2 - 1; i >= 0; i-- {
		opp := len(a) - 1 - i
		a[i], a[opp] = a[opp], a[i]
	}
}

// CurrentSHA accessor for the underlying `CurrentSHA` found in `Info`.
func (d *Root) CurrentSHA() string {
	if d.Info == nil {
		return ""
	}

	return d.Info.CurrentSHA
}

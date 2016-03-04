package models

import (
	"encoding/json"
	"fmt"

	"strings"

	"errors"

	"github.com/hashicorp/consul/api"
)

type Info struct {
	CurrentSha string `json:"current_sha"`
}

type DcdrMap struct {
	Info     *Info                  `json:"info"`
	Features map[string]interface{} `json:"features"`
}

// Features a Feature result set
type Features []Feature

func (a Features) Len() int           { return len(a) }
func (a Features) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Features) Less(i, j int) bool { return a[i].Key < a[j].Key }

func InvalidJsonResponse(msg string) error {
	return errors.New(fmt.Sprintf("invalid json: %s", msg))
}

// FeatureType accepted feature types
type FeatureType string

const (
	Percentile   FeatureType = "percentile"
	Boolean      FeatureType = "boolean"
	Invalid      FeatureType = "invalid"
	DefaultScope             = "default"
)

// GetFeatureType string to type helper
func GetFeatureType(t string) FeatureType {
	switch t {
	case "percentile":
		return Percentile
	case "boolean":
		return Boolean
	default:
		return Invalid
	}
}

// GetFeatureTypeFromValue interface to type helper
func GetFeatureTypeFromValue(v interface{}) FeatureType {
	switch v.(type) {
	case bool:
		return Boolean
	default:
		return Percentile
	}
}

// Feature KV model for feature flags
type Feature struct {
	FeatureType FeatureType `json:"feature_type"`
	Key         string      `json:"key"`
	Namespace   string      `json:"namespace"`
	Scope       string      `json:"scope"`
	Value       interface{} `json:"value"`
	Comment     string      `json:"comment"`
	UpdatedBy   string      `json:"updated_by"`
}

func (f *Feature) GetScope() string {
	if f.Scope == "" {
		f.Scope = DefaultScope
	}

	return f.Scope
}

// ScopedKey expanded key with namespace and scope
func (f *Feature) ScopedKey() string {
	return fmt.Sprintf("%s/%s/%s", f.Namespace, f.GetScope(), f.Key)
}

// NewFeature init a Feature
func NewFeature(name string, value interface{}, comment string, user string, scope string) (f *Feature) {
	var ft FeatureType

	switch value.(type) {
	case float64:
		ft = Percentile
	case bool:
		ft = Boolean
	}

	f = &Feature{
		Key:         name,
		Value:       value,
		FeatureType: ft,
		Comment:     comment,
		UpdatedBy:   user,
		Scope:       scope,
	}

	return
}

// FloatValue cast Value to float64
func (f *Feature) FloatValue() float64 {
	return f.Value.(float64)
}

// BoolValue cast Value to bool
func (f *Feature) BoolValue() bool {
	return f.Value.(bool)
}

// ToJson marshal feature to json
func (f *Feature) ToJson() ([]byte, error) {
	return json.Marshal(f)
}

// ExplodeToMap explode feature namespaces and scopes to nested maps
func (fts Features) ExplodeToMap() map[string]interface{} {
	m := make(map[string]interface{})

	for _, f := range fts {
		explode(m, f.ScopedKey(), f.Value)
	}

	return m
}

// ToJSON exploded JSON
func (fts Features) ToJSON() ([]byte, error) {
	m := fts.ExplodeToMap()
	return json.MarshalIndent(m, "", "  ")
}

// KVsToFeatures helper for unmarshalling consul result
// sets into Features
func KVsToFeatures(bts []byte) (Features, error) {
	var kvs api.KVPairs

	err := json.Unmarshal(bts, &kvs)

	if err != nil {
		return nil, InvalidJsonResponse(fmt.Sprintf("%s", bts))
	}

	var fts Features

	for _, v := range kvs {
		var f Feature

		err := json.Unmarshal(v.Value, &f)

		if err != nil {
			return fts, err
		}

		fts = append(fts, f)
	}

	return fts, nil
}

func explode(m map[string]interface{}, k string, v interface{}) {
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

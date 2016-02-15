package models

import (
	"encoding/json"
	"fmt"

	"strings"

	"github.com/hashicorp/consul/api"
)

type KV map[string]interface{}
type Scopes map[string]KV
type Fts map[string]Scopes
type FeatureMap struct {
	Root Fts
}

type Features []Feature

func (a Features) Len() int           { return len(a) }
func (a Features) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Features) Less(i, j int) bool { return a[i].Key < a[j].Key }

type FeatureType string

const (
	Percentile   FeatureType = "percentile"
	Boolean      FeatureType = "boolean"
	Invalid      FeatureType = "invalid"
	DefaultScope             = "default"
)

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
func (f *Feature) ScopedKey() string {
	return fmt.Sprintf("%s/%s/%s", f.Namespace, f.GetScope(), f.Key)
}

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

func (f *Feature) FloatValue() float64 {
	return f.Value.(float64)
}

func (f *Feature) BoolValue() bool {
	return f.Value.(bool)
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

func (fts Features) ToKVMap() map[string]interface{} {
	m := make(map[string]interface{})

	for _, f := range fts {
		explode(m, f.ScopedKey(), f.Value)
	}

	return m
}

func (fts Features) ToJSON() ([]byte, error) {
	m := fts.ToKVMap()
	return json.MarshalIndent(m, "", "	")
}

func ParseFeatures(bts []byte) (Features, error) {
	var kvs api.KVPairs

	err := json.Unmarshal(bts, &kvs)

	if err != nil {
		return nil, err
	}

	var fts Features

	if err != nil {
		return fts, err
	}

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

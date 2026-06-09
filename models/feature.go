package models

import (
	"encoding/json"
	"fmt"
	"strconv"
)

// Features a Feature result set
type Features []Feature

func (a Features) Len() int           { return len(a) }
func (a Features) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Features) Less(i, j int) bool { return a[i].Key < a[j].Key }

// FeatureType accepted feature types
type FeatureType string

const (
	// Percentile percentile `FeatureType`
	Percentile FeatureType = "percentile"
	// Boolean boolean `FeatureType`
	Boolean FeatureType = "boolean"
	// String string `FeatureType`
	String FeatureType = "string"
	// Invalid invalid `FeatureType`
	Invalid FeatureType = "invalid"
	// FeatureScope scoping for feature keys
	FeatureScope = "features"
)

// ParseValueAndFeatureType string to type helper
func ParseValueAndFeatureType(v string) (interface{}, FeatureType) {
	b, err := strconv.ParseBool(v)

	if err == nil && v != "0" && v != "1" {
		return b, Boolean
	}

	f, err := strconv.ParseFloat(v, 64)

	if err == nil {
		return f, Percentile
	}

	i, err := strconv.ParseInt(v, 10, 64)

	if err == nil {
		return i, Percentile
	}

	// Any value that is not a bool or a number is treated as a free-form
	// string feature, so this function never returns Invalid. The Invalid
	// constant is retained for callers that compare against it explicitly.
	return v, String
}

// ParseFeatureType resolves a user-supplied `--type` value to a FeatureType.
// Only the canonical names are accepted. The second return value is false when
// the type is unrecognized.
func ParseFeatureType(s string) (FeatureType, bool) {
	switch s {
	case string(Boolean):
		return Boolean, true
	case string(Percentile):
		return Percentile, true
	case string(String):
		return String, true
	default:
		return Invalid, false
	}
}

// ParseValueForType parses `val` into the concrete value for an explicitly
// provided `ft`. Unlike ParseValueAndFeatureType it does not infer the type,
// so a String stores the literal value verbatim (e.g. "true" or "0.5").
func ParseValueForType(val string, ft FeatureType) (interface{}, error) {
	switch ft {
	case Boolean:
		return strconv.ParseBool(val)
	case Percentile:
		return strconv.ParseFloat(val, 64)
	case String:
		return val, nil
	default:
		return nil, fmt.Errorf("unsupported feature type %q", ft)
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

// GetScope scope accessor
func (f *Feature) GetScope() string {
	if f.Scope == "" {
		f.Scope = DefaultScope
	}

	return f.Scope
}

// GetNamespace formats the fully scoped namespace
func (f *Feature) GetNamespace() string {
	return fmt.Sprintf("%s/%s", f.Namespace, FeatureScope)
}

// ScopedKey expanded key with namespace and scope
func (f *Feature) ScopedKey() string {
	return fmt.Sprintf("%s/%s/%s", f.GetNamespace(), f.GetScope(), f.Key)
}

// NewFeature create a Feature
func NewFeature(name string, value interface{}, comment string, user string, scope string, ns string) (f *Feature) {
	var ft FeatureType

	switch value.(type) {
	case float64:
		ft = Percentile
	case bool:
		ft = Boolean
	case string:
		ft = String
	}

	f = &Feature{
		Key:         name,
		Value:       value,
		FeatureType: ft,
		Comment:     comment,
		UpdatedBy:   user,
		Scope:       scope,
		Namespace:   ns,
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

// StringValue cast Value to string
func (f *Feature) StringValue() string {
	return f.Value.(string)
}

// ToJSON marshal feature to json
func (f *Feature) ToJSON() ([]byte, error) {
	return json.Marshal(f)
}

package models

type Features []Feature

func (a Features) Len() int           { return len(a) }
func (a Features) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a Features) Less(i, j int) bool { return a[i].Name < a[j].Name }

type FeatureType string

const (
	Percentile FeatureType = "percentile"
	Boolean    FeatureType = "boolean"
	Invalid    FeatureType = "invalid"
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
	Name        string      `json:"name"`
	Value       interface{} `json:"value"`
	Comment     string      `json:"comment"`
}

func PercentileFeature(name string, value float64, comment string) (f *Feature) {
	f = &Feature{
		Name:        name,
		Value:       value,
		FeatureType: Percentile,
		Comment:     comment,
	}

	return
}

func BooleanFeature(name string, value bool, comment string) (f *Feature) {
	f = &Feature{
		Name:        name,
		Value:       value,
		FeatureType: Boolean,
		Comment:     comment,
	}

	return
}

func (f *Feature) FloatValue() float64 {
	return f.Value.(float64)
}

func (f *Feature) BoolValue() bool {
	return f.Value.(bool)
}

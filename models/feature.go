package models

type ByName []Feature

func (a ByName) Len() int           { return len(a) }
func (a ByName) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByName) Less(i, j int) bool { return a[i].Name < a[j].Name }

type FeatureType string

const (
	Percentile FeatureType = "percentile"
	Boolean    FeatureType = "boolean"
	Scalar     FeatureType = "scalar"
)

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

func ScalarFeature(name string, value float64, comment string) (f *Feature) {
	f = &Feature{
		Name:        name,
		Value:       value,
		FeatureType: Scalar,
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

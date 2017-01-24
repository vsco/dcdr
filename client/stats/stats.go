package stats

type IFace interface {
	Incr(feature string, sampleRate float64)
	Tags() []string
}

const (
	Enabled  = "enabled"
	Disabled = "disabled"
	JoinWith = "."
)

package stats

type IFace interface {
	Incr(feature string)
	Tags() []string
}

const (
	Enabled  = "enabled"
	Disabled = "disabled"
	JoinWith = "."
)

package stats

import "github.com/PagerDuty/godspeed"

// Godspeed stats adapter for Godspeed.
type Godspeed struct {
	gs   *godspeed.Godspeed
	tags []string
}

// New creates a new Godspeed stats adapter.
func New(gs *godspeed.Godspeed, tags []string) (g *Godspeed) {
	g = &Godspeed{
		gs: gs,
	}

	return
}

// Incr increments a key.
func (g *Godspeed) Incr(key string) {
	g.gs.Incr(key, g.tags)
}

// Tags accessor for `tags`.
func (g *Godspeed) Tags() []string {
	return g.tags
}

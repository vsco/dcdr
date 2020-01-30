package stats

import "github.com/theckman/godspeed"

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
func (g *Godspeed) Incr(key string, sampleRate float64) {
	g.gs.Send(key, "c", 1, sampleRate, g.tags)
}

// Tags accessor for `tags`.
func (g *Godspeed) Tags() []string {
	return g.tags
}

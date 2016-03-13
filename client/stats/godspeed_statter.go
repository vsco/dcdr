package stats

import "github.com/PagerDuty/godspeed"

type Godspeed struct {
	gs   *godspeed.Godspeed
	tags []string
}

func NewGodspeedStatter(gs *godspeed.Godspeed, tags []string) (g *Godspeed) {
	g = &Godspeed{
		gs: gs,
	}

	return
}

func (g *Godspeed) Incr(key string) {
	g.gs.Incr(key, g.tags)
}

func (g *Godspeed) Tags() []string {
	return g.tags
}

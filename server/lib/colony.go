package colony

import (
	"encoding/gob"
	"time"
)

var _ ProducerObject = &Colony{}

func init() {
	gob.Register(&Colony{})
}

var COLONY_CYCLE = 3

type Colony struct {
	O      Owner
	Point  Point
	P      bool
	Bucket int
	Age    int64
	Cycle  int
}

func NewColony(o Owner, p Point) *Colony {
	return &Colony{
		O:      o,
		Point:  p,
		Bucket: 5,
	}
}

func (c *Colony) Owner() Owner {
	return c.O
}

func (c *Colony) Center() Point {
	return c.Point
}

func (c *Colony) Tick() {
	c.Cycle = (c.Cycle + 1) % COLONY_CYCLE
}

func (c *Colony) Dead() bool {
	if c.Age == 0 {
		c.Touch()
	}
	// Colonies live for 30 hours without activity
	now := time.Now().Unix()
	return now > c.Age+(30*60*60)
}

func (c *Colony) Touch() {
	c.Age = time.Now().Unix()
}

func (c *Colony) Reclaim(o Object) {
	if a, ok := o.(*Ant); ok && a.Owner() == c.Owner() {
		c.Bucket = c.Bucket + 1
	}
}

func (c *Colony) Produce() (Object, bool) {
	if c.P {
		c.P = false
		if c.Bucket < 1 {
			return nil, false
		}
		c.Bucket = c.Bucket - 1
		switch c.Cycle {
		case 2:
			return NewAnt(c.O, 3), true
		default:
			return NewAnt(c.O, 1), true
		}
	}
	return nil, false
}

func (c *Colony) View(o Owner) *ObjectView {
	return &ObjectView{
		Type:     "colony",
		Mine:     o == c.O,
		Strength: c.Bucket,
	}
}

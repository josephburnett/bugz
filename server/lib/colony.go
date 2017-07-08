package colony

import (
	"encoding/gob"
)

var _ ProducerObject = &Colony{}

func init() {
	gob.Register(&Colony{})
}

type Colony struct {
	O      Owner
	Point  Point
	P      bool
	Bucket int
}

func (c *Colony) Owner() Owner {
	return c.O
}

func (c *Colony) Center() Point {
	return c.Point
}

func (c *Colony) Tick() {
	if c.P == false && c.Bucket < 15 {
		c.Bucket = c.Bucket + 1
	}
}

func (c *Colony) Dead() bool {
	return false
}

func (c *Colony) Reclaim(_ Object) {
	return
}

func (c *Colony) Produce() (Object, bool) {
	if c.P {
		c.Bucket = c.Bucket - 1
		if c.Bucket < 1 {
			c.P = false
		}
		return &Ant{
			O:         c.O,
			Direction: RandomDirection(D_AROUND),
			S:         1,
			Endurance: 40,
		}, true
	}
	return nil, false
}

func (c *Colony) View(o Owner) *ObjectView {
	return &ObjectView{
		Type: "colony",
		Mine: o == c.O,
	}
}

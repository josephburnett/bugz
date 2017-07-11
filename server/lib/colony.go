package colony

import (
	"encoding/gob"
	"time"
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
	Age    int64
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

func (c *Colony) Reclaim(_ Object) {
	return
}

func (c *Colony) Produce() (Object, bool) {
	if c.P {
		c.Bucket = c.Bucket - 1
		if c.Bucket < 1 {
			c.P = false
		}
		return NewAnt(c.O), true
	}
	return nil, false
}

func (c *Colony) View(o Owner) *ObjectView {
	return &ObjectView{
		Type: "colony",
		Mine: o == c.O,
	}
}

package colony

var _ ProducerObject = &Colony{}

type Colony struct {
	owner   Owner
	point   Point
	produce bool
	bucket  int
}

func (c *Colony) Owner() Owner {
	return c.owner
}

func (c *Colony) Center() Point {
	return c.point
}

func (c *Colony) Tick() {
	if c.produce == false && c.bucket < 15 {
		c.bucket = c.bucket + 1
	}
}

func (c *Colony) Dead() bool {
	return false
}

func (c *Colony) Reclaim(_ Object) {
	return
}

func (c *Colony) Produce() (Object, bool) {
	if c.produce {
		c.bucket = c.bucket - 1
		if c.bucket < 1 {
			c.produce = false
		}
		return &Ant{
			owner:     c.owner,
			direction: RandomDirection(D_AROUND),
			strength:  1,
			endurance: 40,
		}, true
	}
	return nil, false
}

func (c *Colony) View(o Owner) *ObjectView {
	return &ObjectView{
		Type: "colony",
		Mine: o == c.owner,
	}
}

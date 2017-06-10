package "colony"

type Colony struct {
	Owner string
	Point Point
	Direction Direction
	Speed int
}

func (c *Colony) Owner() string {
	return c.Owner
}

func (c *Colony) Point() Point {
	return c.Point
}

func (c *Colony) Move(_ *Surroudings, _ *Phermones) Point {
	return c.Point
}

func (c *Colony) Resolve(o *Object) {
	// implement rules of engagement
}
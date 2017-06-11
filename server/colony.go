package colony

type Colony struct {
	owner     Owner
	point     Point
	direction Direction
	speed     int
}

func (c *Colony) Owner() Owner {
	return c.owner
}

func (c *Colony) Point() Point {
	return c.point
}

func (c *Colony) Produce(o Objects, p Phermones) (*Ant, bool) {
	return nil, false
}
package colony

var _ Object = &Colony{}

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

func (c *Colony) Move(_ Surroundings, _ Phermones) Point {
	return c.point
}

func (c *Colony) Fight(o *Object) bool {
	return true
}

func (c *Colony) Dead() bool {
	return false
}

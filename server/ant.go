package colony

var _ Object = &Ant{}

type Ant struct {
	owner     Owner
	point     Point
	direction Direction
	speed     int
	strength  int
	endurance int
}

func (a *Ant) Owner() Owner {
	return a.owner
}

func (a *Ant) Point() Point {
	return a.point
}

func (a *Ant) Move(s Surroundings, p Phermones) Point {
	return a.point
}

func (a *Ant) Fight(o *Object) bool {
	return true
}

func (a *Ant) Dead() bool {
	return false
}

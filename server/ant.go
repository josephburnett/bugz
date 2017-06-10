package "colony"

var _ Object = Ant{}

type Ant struct {
	Owner Owner
	Point Point
	Direction Direction
	Speed int
	Strength int
	Endurance int
}

func (a *Ant) Owner() Owner {
	return a.Owner
}

func (a *Ant) Point() Point {
	return a.Point
}

func (a *Ant) Move(s *Surroundings, p *Phermones) Point {
	// implement movements
}

func (a *Ant) Resolve(o *Object) {
	// implement rules of engagement
}
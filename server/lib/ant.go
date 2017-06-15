package colony

var _ Object = &Ant{}
var _ AnimateObject = &Ant{}

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

func (a *Ant) Move(o Objects, p Phermones) Point {
	return a.point
}

func (a *Ant) Attack(o Object) bool {
	switch ao := o.(type) {
	default:
		return false
	case AnimateObject:
		defense := ao.Strength()
		ao.TakeDamage(a.strength)
		a.TakeDamage(defense)
		return a.Dead()
	}
}

func (a *Ant) Strength() int {
	return a.strength
}

func (a *Ant) TakeDamage(d int) {
	a.strength = a.strength - d
}

func (a *Ant) Dead() bool {
	return a.strength > 0
}

func (a *Ant) View() *ObjectView {
	return &ObjectView{
		direction: a.direction,
	}
}

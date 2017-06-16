package colony

import "log"

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
	possible := func(d Direction) bool {
		newPoint := a.point.Plus(d)
		object, exists := o[newPoint]
		if !exists || object.Owner() != a.owner {
			return true
		}
		return false
	}
	options := make([]Point, 0, 8)
	randomChoice := func() Point {
		a.point = RandomPoint(options)
		return a.point
	}
	for _, d := range a.direction.InFront() {
		target := a.point.Plus(d)
		if _, hasPhermone := p[target]; hasPhermone && possible(d) {
			options = append(options, target)
		}
	}
	if len(options) > 0 {
		log.Println("following a phermone")
		return randomChoice()
	}
	for _, d := range a.direction.InFront() {
		target := a.point.Plus(d)
		if possible(d) {
			options = append(options, target)
		}
	}
	if len(options) > 0 {
		log.Println("following momentum")
		return randomChoice()
	}
	for _, d := range a.direction.Around() {
		target := a.point.Plus(d)
		if possible(d) {
			options = append(options, target)
		}
	}
	if len(options) > 0 {
		log.Println("choosing random direction")
		return randomChoice()
	}
	log.Println("boxed in")
	return a.point
}

func (a *Ant) Attack(o Object) bool {
	switch ao := o.(type) {
	default:
		return false
	// TODO: kill a colony
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
	return a.strength <= 0
}

func (a *Ant) View() *ObjectView {
	return &ObjectView{
		Direction: a.direction,
	}
}

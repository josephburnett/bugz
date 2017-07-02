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
	cycle     int
}

var CYCLE = 9

func (a *Ant) Owner() Owner {
	return a.owner
}

func (a *Ant) Type() string {
	return "ant"
}

func (a *Ant) Point() Point {
	return a.point
}

func (a *Ant) Tick() {
	// Advance behavior cycle
	a.cycle = (a.cycle + 1) % CYCLE
	a.endurance = a.endurance - 1
	if a.endurance == 0 {
		a.strength = 0
	}
}

func (a *Ant) Move(o map[Point]Object, p Phermones, friends map[Owner]bool) Point {
	possible := func(d Direction) bool {
		newPoint := a.point.Plus(d)
		object, exists := o[newPoint]
		if !exists {
			return true
		}
		if object.Owner() == a.owner {
			return false
		}
		if friend, ok := friends[object.Owner()]; ok && friend {
			return false
		}
		return true // attack
	}
	options := make([]Direction, 0, 8)
	move := func() Point {
		d := RandomDirection(options)
		a.point = a.point.Plus(d)
		a.direction = d
		return a.point
	}
	// Die
	if a.endurance == 0 {
		a.strength = 0
		return a.point
	}
	// Follow a phermone, in front
	for _, d := range a.direction.InFront() {
		target := a.point.Plus(d)
		if _, hasPhermone := p[target]; hasPhermone && possible(d) {
			options = append(options, d)
		}
	}
	if len(options) > 0 {
		return move()
	}
	// Follow a phermone, near by
	for _, d := range Around() {
		target := a.point.Plus(d)
		if _, hasPhermone := p[target]; hasPhermone && possible(d) {
			options = append(options, d)
		}
	}
	if len(options) > 0 {
		return move()
	}
	switch a.cycle {
	default:
		// Follow momentum
		for _, d := range a.direction.InFront() {
			if possible(d) {
				options = append(options, d)
			}
		}
		if len(options) > 0 {
			return move()
		}
	case CYCLE - 2:
		// Stay put
		return a.point
	case CYCLE - 1:
		// Random move
		for _, d := range Around() {
			if possible(d) {
				options = append(options, d)
			}
		}
		if len(options) > 0 {
			return move()
		}
	}
	// Boxed in
	return a.point
}

func (a *Ant) Attack(o Object) bool {
	switch ao := o.(type) {
	default:
		return false
	case AnimateObject:
		defense := ao.Strength()
		attack := a.strength
		if defense > attack {
			a.TakeDamage(defense)
		} else {
			ao.TakeDamage(attack)
		}
		return !a.Dead()
	case Object:
		log.Println(a.Owner(), a.Type(), "eats", o.Type())
		return true
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

func (a *Ant) View(o Owner) *ObjectView {
	return &ObjectView{
		Type:      a.Type(),
		Direction: a.direction,
		Mine:      o == a.owner,
	}
}

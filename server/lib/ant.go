package colony

import "log"

var _ AnimateObject = &Ant{}

type Ant struct {
	owner     Owner
	direction Direction
	strength  int
	endurance int
	cycle     int
}

var CYCLE = 9

func (a *Ant) Owner() Owner {
	return a.owner
}

func (a *Ant) Tick() {
	// Advance behavior cycle
	a.cycle = (a.cycle + 1) % CYCLE
	a.endurance = a.endurance - 1
	if a.endurance == 0 {
		a.strength = 0
	}
}

func (a *Ant) Move(point Point, o map[Direction]Object, p Phermones, friends Friends) Point {
	possible := func(d Direction) bool {
		object, exists := o[d]
		if !exists {
			return true // vacant
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
		a.direction = d
		return point.Plus(d)
	}
	// Die
	if a.endurance == 0 {
		a.strength = 0
		return point
	}
	// Follow a phermone, in front
	for _, d := range a.direction.InFront() {
		target := point.Plus(d)
		if _, hasPhermone := p[target]; hasPhermone && possible(d) {
			options = append(options, d)
		}
	}
	if len(options) > 0 {
		return move()
	}
	// Follow a phermone, near by
	for _, d := range Around() {
		target := point.Plus(d)
		if _, hasPhermone := p[target]; hasPhermone && possible(d) {
			options = append(options, d)
		}
	}
	if len(options) > 0 {
		return move()
	}
	// Non-phemone based, cyclic movement
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
		return point
	case CYCLE - 1:
		// Turn
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
	return point
}

func (a *Ant) Attack(o Object) bool {
	switch o := o.(type) {
	default:
		return false
	case AnimateObject:
		defense := o.Strength()
		attack := a.strength
		if defense > attack {
			a.TakeDamage(defense)
		} else {
			o.TakeDamage(attack)
		}
		return !a.Dead()
	case Object:
		log.Printf("%v ant eats %v %T", a.Owner(), o.Owner(), o)
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
		Type:      "ant",
		Direction: a.direction,
		Mine:      o == a.owner,
	}
}

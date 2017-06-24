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
	cycle     int
}

var CYCLE = 9

func (a *Ant) Owner() Owner {
	return a.owner
}

func (a *Ant) Point() Point {
	return a.point
}

func (a *Ant) Move(o Objects, p Phermones) Point {
	// Die
	if a.endurance == 0 {
		a.strength = 0
		return a.point
	}
	// Advance behavior cycle
	a.cycle = (a.cycle + 1) % CYCLE
	possible := func(d Direction) bool {
		newPoint := a.point.Plus(d)
		object, exists := o[newPoint]
		if !exists || object.Owner() != a.owner {
			return true
		}
		return false
	}
	options := make([]Direction, 0, 8)
	move := func() Point {
		d := RandomDirection(options)
		a.point = a.point.Plus(d)
		a.direction = d
		a.endurance = a.endurance - 1
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

package colony

import (
	"encoding/gob"
	"log"
)

var _ AnimateObject = &Ant{}

func init() {
	gob.Register(&Ant{})
}

type Ant struct {
	O         Owner
	Direction Direction
	S         int
	Endurance int
	Cycle     int
}

func NewAnt(o Owner, strength int) *Ant {
	return &Ant{
		O:         o,
		Direction: RandomDirection(D_AROUND),
		S:         strength,
		Endurance: 600, // 5 minutes at 2 fps
	}
}

var CYCLE = 9

func (a *Ant) Owner() Owner {
	return a.O
}

func (a *Ant) Tick() {
	// Advance behavior cycle
	a.Cycle = (a.Cycle + 1) % CYCLE
	a.Endurance = a.Endurance - 1
	if a.Endurance == 0 {
		a.S = 0
	}
}

var A_MAX_DISTANCE int = 15

func (a *Ant) Move(point Point, home Point, o map[Direction]Object, p Phermones, friends Friends) Point {
	possible := func(d Direction) bool {
		object, exists := o[d]
		if !exists {
			return true // vacant
		}
		if object.Owner() == a.O {
			return false
		}
		if friend, ok := friends[object.Owner()]; ok && friend {
			return false
		}
		return true // attack
	}
	// Die
	if a.Endurance == 0 {
		a.S = 0
		return point
	}
	// Don't stray too far from home
	options := make([]Direction, 0, 8)
	if point.DistanceFrom(home) >= A_MAX_DISTANCE {
		for _, d := range Around() {
			if point.Plus(d).DistanceFrom(home) < A_MAX_DISTANCE {
				options = append(options, d)
			}
		}
		if len(options) >= 0 {
			// Choose a random direction that leads closer to home
			d := RandomDirection(options)
			a.Direction = d
			return point.Plus(d)
		} else {
			// Stay put
			return point
		}
	}
	// Follow a phermone, directly ahead
	ahead := point.Plus(a.Direction)
	if _, hasPhermone := p[ahead]; hasPhermone && possible(a.Direction) {
		return ahead
	}
	// Follow a phermone, in front somewhere
	for _, d := range a.Direction.InFront() {
		target := point.Plus(d)
		if _, hasPhermone := p[target]; hasPhermone && possible(d) {
			options = append(options, d)
		}
	}
	if len(options) > 0 {
		d := RandomDirection(options)
		a.Direction = d
		return point.Plus(d)
	}
	// Follow momentum
	if possible(a.Direction) {
		return ahead
	}
	// Turn around
	if possible(a.Direction.Opposite()) {
		a.Direction = a.Direction.Opposite()
		return point.Plus(a.Direction)
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
		attack := a.Strength()
		a.TakeDamage(defense)
		o.TakeDamage(attack)
		return !a.Dead()
	case Object:
		log.Printf("%v ant eats %v %T", a.Owner(), o.Owner(), o)
		return true
	}
}

func (a *Ant) Strength() int {
	return a.S
}

func (a *Ant) TakeDamage(d int) {
	a.S = a.S - d
}

func (a *Ant) Dead() bool {
	return a.S <= 0
}

func (a *Ant) View(o Owner) *ObjectView {
	return &ObjectView{
		Type:      "ant",
		Direction: a.Direction,
		Mine:      o == a.O,
		Strength:  a.S,
	}
}

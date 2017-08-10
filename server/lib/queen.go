package colony

import (
	"encoding/gob"
)

var _ AnimateObject = &Queen{}

func init() {
	gob.Register(&Queen{})
}

type Queen struct {
	Ant
	Colony *Colony
}

func NewQueen(c *Colony) *Queen {
	return &Queen{
		Ant:    *NewAnt(c.O, 7),
		Colony: c,
	}
}

func (q *Queen) Owner() Owner {
	return q.Ant.Owner()
}

func (q *Queen) Tick() {
	q.Ant.Tick()
}

func (q *Queen) Dead() bool {
	return q.Ant.Dead()
}

func (q *Queen) View(o Owner) *ObjectView {
	view := q.Ant.View(o)
	view.Type = "queen"
	return view
}

func (q *Queen) Move(p Point, d map[Direction]Object, ph Phermones, f Friends) Point {
	return q.Ant.Move(p, d, ph, f)
}

func (q *Queen) Attack(o Object) bool {
	return q.Ant.Attack(o)
}

func (q *Queen) TakeDamage(d int) {
	q.Ant.TakeDamage(d)
}

func (q *Queen) Strength() int {
	return q.Ant.Strength()
}

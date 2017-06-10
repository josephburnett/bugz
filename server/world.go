package colony

import "math/rand"

type Owner string
type Point [2]int
type Direction [2]int
type Phermones map[Point]bool
type Surroundings map[Point]*Object

func (p Point) Plus(d Direction) Point {
	return Point{p[0] + d[0], p[1] + d[1]}
}

func (p1 Point) Equals(p2 Point) bool {
	if p1[0] == p2[0] && p1[1] == p2[1] {
		return true
	}
	return false
}

type Object interface {
	Owner() Owner
	Point() Point
	Move(Surroundings, Phermones) Point
	Fight(*Object) bool
	Dead() bool
}

type World struct {
	Map       Surroundings
	Phermones map[Owner]Phermones
}

func (w *World) Tick() {
	objects := make([]*Object, len(w.Map))
	for _, o := range w.Map {
		objects = append(objects, o)
	}
	perm := rand.Perm(len(objects))
	for _, i := range perm {
		o := *objects[i]
		if o.Dead() {
			continue
		}
		fromPoint := o.Point()
		toPoint := o.Move(w.Map, w.Phermones[o.Owner()])
		if fromPoint.Equals(toPoint) {
			continue
		}
		target, occupied := w.Map[toPoint]
		if occupied {
			win := o.Fight(target)
			if win {
				w.Map[toPoint] = &o
			}
			delete(w.Map, fromPoint)
		}
	}
}

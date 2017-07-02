package colony

import (
	"log"
	"math/rand"
)

type Owner string
type Point [2]int
type Direction [2]int
type Phermones map[Point]bool
type Friends map[Owner]bool

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
	Type() string
	Point() Point
	Tick()
	Dead() bool
	View(Owner) *ObjectView
}

type AnimateObject interface {
	Move(map[Point]Object, Phermones, map[Owner]bool) Point
	Attack(Object) bool
	TakeDamage(int)
	Strength() int
}

type World struct {
	owners    map[Owner]*Colony
	friends   map[Owner]Friends
	phermones map[Owner]Phermones
	objects   map[Point]Object
	colonies  map[Point]*Colony
	soil      map[Point]*Soil
}

func NewWorld() *World {
	return &World{
		owners:    make(map[Owner]*Colony),
		friends:   make(map[Owner]Friends),
		phermones: make(map[Owner]Phermones),
		objects:   make(map[Point]Object),
		colonies:  make(map[Point]*Colony),
		soil:      make(map[Point]*Soil),
	}
}

func (w *World) NewColony(o Owner) {
	var p Point
	for {
		p = Point{
			rand.Intn(40) - 20,
			rand.Intn(40) - 20,
		}
		if _, occupied := w.colonies[p]; !occupied {
			break
		}
	}
	c := &Colony{
		owner: o,
		point: p,
	}
	w.owners[o] = c
	w.phermones[o] = make(Phermones)
	w.colonies[p] = c
	log.Println("Created new colony " + o)
}

func (w *World) KillColony(o Owner) {
	c, ok := w.owners[o]
	if !ok {
		return
	}
	delete(w.owners, o)
	delete(w.phermones, o)
	delete(w.colonies, c.Point())
}

func (w *World) Friend(a Owner, b Owner) {
	friendsA, ok := w.friends[a]
	if !ok {
		friendsA = make(Friends)
		w.friends[a] = friendsA
	}
	friendsA[b] = true
	friendsB, ok := w.friends[b]
	if !ok {
		friendsB = make(Friends)
		w.friends[b] = friendsB
	}
	friendsB[a] = true
}

func (w *World) Unfriend(a Owner, b Owner) {
	friendsA, ok := w.friends[a]
	if ok {
		delete(friendsA, b)
	}
	friendsB, ok := w.friends[b]
	if ok {
		delete(friendsB, a)
	}
}

func (w *World) Advance() {
	livingObjects := make([]Object, 0, len(w.objects))
	// Age objects
	for point, o := range w.objects {
		o.Tick()
		if o.Dead() {
			log.Println(o.Owner(), o.Type(), "dies of natural causes")
			if _, enriched := w.soil[point]; !enriched {
				w.soil[point] = &Soil{}
			}
			w.soil[point].Enrich()
			delete(w.objects, point)
		} else {
			livingObjects = append(livingObjects, o)
		}
	}
	// Age soil
	for _, soil := range w.soil {
		soil.Tick()
	}
	// Visit objects in random order
	perm := rand.Perm(len(livingObjects))
	// Move objects
	for _, i := range perm {
		o := livingObjects[i]
		if o.Dead() {
			// Killed by another moving object
			continue
		}
		if ao, ok := o.(AnimateObject); ok {
			fromPoint := o.Point()
			toPoint := ao.Move(w.objects, w.phermones[o.Owner()], w.friends[o.Owner()])
			if fromPoint.Equals(toPoint) {
				// No move
				continue
			}
			target, occupied := w.objects[toPoint]
			if occupied {
				win := ao.Attack(target)
				if win {
					log.Println(o.Owner(), " ", o.Type(), " kills ", target.Owner(), " ", target.Type())
					w.objects[toPoint] = o
				} else {
					log.Println(o.Owner(), " ", o.Type(), " is killed by ", target.Owner(), " ", target.Type())
				}
				delete(w.objects, fromPoint)
			} else {
				w.objects[toPoint] = o
				delete(w.objects, fromPoint)
			}
		}
	}
	// Produce objects
	for _, c := range w.colonies {
		ant, produced := c.Produce(w.objects)
		if produced {
			w.objects[ant.Point()] = ant
		}
	}
	for point, s := range w.soil {
		if object, produced := s.Produce(); produced {
			if _, occupied := w.objects[point]; !occupied {
				w.objects[point] = object
			}
		}
	}
	// Remove the dead stuff
	for _, o := range w.objects {
		if o.Dead() {
			delete(w.objects, o.Point())
		}
	}
	for point, s := range w.soil {
		if s.Dead() {
			delete(w.soil, point)
		}
	}
}

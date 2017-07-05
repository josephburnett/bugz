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
	Tick()
	Dead() bool
	View(Owner) *ObjectView
}

type AnimateObject interface {
	Object
	Move(Point, map[Direction]Object, Phermones, Friends) Point
	Attack(Object) bool
	TakeDamage(int)
	Strength() int
}

type ProducerObject interface {
	Object
	Produce() (Object, bool)
	Reclaim(Object)
}

type World struct {
	// Players
	colonies map[Owner]*Colony
	friends  map[Owner]Friends
	// Layers
	phermones map[Owner]Phermones
	objects   map[Point]Object
	earth     map[Point]ProducerObject
}

func NewWorld() *World {
	return &World{
		colonies:  make(map[Owner]*Colony),
		friends:   make(map[Owner]Friends),
		phermones: make(map[Owner]Phermones),
		objects:   make(map[Point]Object),
		earth:     make(map[Point]ProducerObject),
	}
}

func (w *World) NewColony(o Owner) {
	var p Point
	for {
		p = Point{
			rand.Intn(20),
			rand.Intn(20),
		}
		if _, occupied := w.earth[p]; !occupied {
			break
		}
	}
	c := &Colony{
		owner: o,
		point: p,
	}
	w.colonies[o] = c
	w.phermones[o] = make(Phermones)
	w.earth[p] = c
	log.Println("Created new colony", o)
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

func (w *World) Reclaim(p Point, o Object) {
	if _, ok := w.earth[p]; !ok {
		w.earth[p] = &Soil{}
	}
	w.earth[p].Reclaim(o)
}

func (w *World) Drop(o Owner, what string) {
	var p Point
	for {
		d := Direction{
			rand.Intn(20) - 10,
			rand.Intn(20) - 10,
		}
		p = w.colonies[o].Center().Plus(d)
		if _, colony := w.earth[p].(*Colony); !colony {
			break
		}
	}
	switch what {
	case "rock":
		log.Println(o, "drops a rock")
		if _, ok := w.objects[p]; ok {
			log.Printf("%v %T is destroyed by a rock")
		}
		w.objects[p] = NewRock()
	default:
		log.Println(o, "tries to drop a", what)
	}
}

func (w *World) Advance() {
	// Age earth
	for point, o := range w.earth {
		o.Tick()
		if o.Dead() {
			log.Printf("%v %T fades away", o.Owner(), o)
			delete(w.earth, point)
		}
	}
	// Age objects
	livingObjects := make([]Object, 0, len(w.objects))
	livingObjectPoints := make([]Point, 0, len(w.objects))
	for point, o := range w.objects {
		o.Tick()
		if o.Dead() {
			log.Printf("%v %T dies of natural causes\n", o.Owner(), o)
			w.Reclaim(point, o)
			delete(w.objects, point)
		} else {
			livingObjects = append(livingObjects, o)
			livingObjectPoints = append(livingObjectPoints, point)
		}
	}
	// Visit objects in random order
	perm := rand.Perm(len(livingObjects))
	// Move animate objects
	for _, i := range perm {
		o := livingObjects[i]
		point := livingObjectPoints[i]
		if o.Dead() {
			// Killed by another moving object
			continue
		}
		if ao, ok := o.(AnimateObject); ok {
			surroundings := make(map[Direction]Object)
			for _, d := range Surrounding() {
				if so, ok := w.objects[point.Plus(d)]; ok {
					surroundings[d] = so
				}
			}
			destination := ao.Move(point, surroundings, w.phermones[o.Owner()], w.friends[o.Owner()])
			if point.Equals(destination) {
				// No move
				continue
			}
			target, occupied := w.objects[destination]
			if occupied {
				if win := ao.Attack(target); win {
					log.Println("%v %T kills %v %T\n", o.Owner(), o, target.Owner(), target)
					w.Reclaim(destination, target)
					w.objects[destination] = o
				} else {
					log.Println("%v %T is killed by %v %T\n", o.Owner(), o, target.Owner(), target)
					w.Reclaim(point, o)
				}
				delete(w.objects, point)
			} else {
				w.objects[destination] = o
				delete(w.objects, point)
			}
		}
	}
	// Produce objects
	for point, producer := range w.earth {
		if _, obstructed := w.objects[point]; !obstructed {
			o, produced := producer.Produce()
			if produced {
				w.objects[point] = o
			}
		}
	}
	// Remove the dead stuff
	for point, o := range w.objects {
		if o.Dead() {
			delete(w.objects, point)
		}
	}
	for point, producer := range w.earth {
		if producer.Dead() {
			delete(w.earth, point)
		}
	}
}

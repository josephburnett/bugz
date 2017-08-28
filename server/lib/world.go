package colony

import (
	"bytes"
	"encoding/gob"
	"io/ioutil"
	"log"
	"math/rand"
	"os"
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
	Move(Point, Point, map[Direction]Object, Phermones, Friends) Point
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
	Colonies map[Owner]Point
	Friends  map[Owner]Friends
	// Layers
	Phermones map[Owner]Phermones
	Objects   map[Point]Object
	Earth     map[Point]ProducerObject
}

func NewWorld() *World {
	return &World{
		Colonies:  make(map[Owner]Point),
		Friends:   make(map[Owner]Friends),
		Phermones: make(map[Owner]Phermones),
		Objects:   make(map[Point]Object),
		Earth:     make(map[Point]ProducerObject),
	}
}

func LoadWorld(filename string) (*World, error) {
	w := NewWorld()
	if _, err := os.Stat(filename); os.IsNotExist(err) {
		log.Println("creating world file")
		err := w.SaveWorld(filename)
		if err != nil {
			return nil, err
		}
		return w, nil
	}
	b, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}
	buf := bytes.NewBuffer(b)
	dec := gob.NewDecoder(buf)
	err = dec.Decode(w)
	if err != nil {
		return nil, err
	}
	log.Println("loaded the world")
	return w, nil
}

func (w *World) SaveWorld(filename string) error {
	var b bytes.Buffer
	enc := gob.NewEncoder(&b)
	err := enc.Encode(w)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(filename, b.Bytes(), 0644)
	if err != nil {
		return err
	}
	log.Println("saved the world")
	return nil
}

func (w *World) NewColony(o Owner) {
	var p Point
	for {
		p = Point{
			rand.Intn(20),
			rand.Intn(20),
		}
		if _, occupied := w.Earth[p]; !occupied {
			break
		}
	}
	c := &Colony{
		O:     o,
		Point: p,
	}
	w.Colonies[o] = p
	w.Phermones[o] = make(Phermones)
	w.Earth[p] = c
	log.Println("created new colony", o)
}

func (w *World) FindColony(o Owner) (*Colony, bool) {
	if point, ok := w.Colonies[o]; ok {
		if colony, ok := w.Earth[point].(*Colony); ok {
			return colony, true
		}
	}
	return nil, false
}

func (w *World) Friend(a Owner, b Owner) {
	friendsA, ok := w.Friends[a]
	if !ok {
		friendsA = make(Friends)
		w.Friends[a] = friendsA
	}
	friendsA[b] = true
	friendsB, ok := w.Friends[b]
	if !ok {
		friendsB = make(Friends)
		w.Friends[b] = friendsB
	}
	friendsB[a] = true
}

func (w *World) Unfriend(a Owner, b Owner) {
	friendsA, ok := w.Friends[a]
	if ok {
		delete(friendsA, b)
	}
	friendsB, ok := w.Friends[b]
	if ok {
		delete(friendsB, a)
	}
}

func (w *World) Reclaim(p Point, o Object) {
	// A Queen becomes the colony that she carries
	if queen, ok := o.(*Queen); ok {
		colony := queen.Colony
		colony.Point = p
		w.Earth[p] = colony
		w.Colonies[queen.Owner()] = p
		return
	}
	// Everything else returns to dust
	if _, ok := w.Earth[p]; !ok {
		w.Earth[p] = &Soil{}
	}
	w.Earth[p].Reclaim(o)
}

func (w *World) Drop(o Owner, what string) {
	var p Point
	for {
		d := Direction{
			rand.Intn(20) - 10,
			rand.Intn(20) - 10,
		}
		p = w.Colonies[o].Plus(d)
		if _, colony := w.Earth[p].(*Colony); !colony {
			break
		}
	}
	switch what {
	case "rock":
		log.Println(o, "drops a rock")
		if _, ok := w.Objects[p]; ok {
			log.Printf("%v %T is destroyed by a rock")
		}
		w.Objects[p] = NewRock()
	case "food":
		log.Println(o, "drop some food")
		if _, ok := w.Objects[p]; ok {
			log.Printf("%v %T is destroyed by falling food")
		}
		w.Objects[p] = NewFruit()
	default:
		log.Println(o, "tries to drop a", what)
	}
}

func (w *World) Advance() {
	// Age earth
	for point, o := range w.Earth {
		o.Tick()
		if o.Dead() {
			log.Printf("%v %T fades away", o.Owner(), o)
			delete(w.Earth, point)
			if _, ok := o.(*Colony); ok {
				delete(w.Colonies, o.Owner())
				delete(w.Friends, o.Owner())
				delete(w.Phermones, o.Owner())
			}
		}
	}
	// Age objects
	livingObjects := make([]Object, 0, len(w.Objects))
	livingObjectPoints := make([]Point, 0, len(w.Objects))
	for point, o := range w.Objects {
		o.Tick()
		if o.Dead() {
			log.Printf("%v %T dies of natural causes\n", o.Owner(), o)
			w.Reclaim(point, o)
			delete(w.Objects, point)
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
		move := func(to Point) {
			w.Objects[to] = o
			// Colonies follow the Queen
			if queen, ok := o.(*Queen); ok {
				w.Colonies[queen.Colony.Owner()] = to
			}
		}
		if ao, ok := o.(AnimateObject); ok {
			surroundings := make(map[Direction]Object)
			for _, d := range Surrounding() {
				if so, ok := w.Objects[point.Plus(d)]; ok {
					surroundings[d] = so
				}
			}
			destination := ao.Move(point, w.Colonies[ao.Owner()], surroundings, w.Phermones[o.Owner()], w.Friends[o.Owner()])
			if point.Equals(destination) {
				// No move
				continue
			}
			target, occupied := w.Objects[destination]
			if occupied {
				if win := ao.Attack(target); win {
					log.Printf("%v %T kills %v %T\n", o.Owner(), o, target.Owner(), target)
					w.Reclaim(destination, target)
					move(destination)
				} else {
					log.Printf("%v %T is killed by %v %T\n", o.Owner(), o, target.Owner(), target)
					w.Reclaim(point, o)
				}
				delete(w.Objects, point)
			} else {
				move(destination)
				delete(w.Objects, point)
			}
		}
	}
	// Produce objects
	for point, producer := range w.Earth {
		if _, obstructed := w.Objects[point]; !obstructed {
			o, produced := producer.Produce()
			if produced {
				w.Objects[point] = o
			}
		}
	}
	// Remove the dead stuff
	for point, o := range w.Objects {
		if o.Dead() {
			delete(w.Objects, point)
		}
	}
	for point, producer := range w.Earth {
		if producer.Dead() {
			delete(w.Earth, point)
		}
	}
}

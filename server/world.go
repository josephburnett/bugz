package "colony"

type Owner Owner
type Point int[2]
type Direction int[2]
type Phermones map[Point]bool
type Surroundings map[Point]*Object

const (
    NONE = Direction{0,0}
	UP = Direction{0,1}
	UP_RIGHT = Direction{1,1}
	RIGHT = Direction{0,1}
	DOWN_RIGHT = Direction{1,-1}
	DOWN = Direction{0,-1}
	DOWN_LEFT = Direction{-1,-1}
	LEFT = Direction{-1,0}
	UP_LEFT = Direction{-1,1}
)

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
	Move(*Surroundings, *Phermones) Point
	Resolve(*Object) bool
}

type World struct {
	Objects []*Object
	Map map[Point]*Object
	Phermones map[Owner]*Phermones
}

func (w *World) Tick(t int) {
	for o := w.Objects {
		owner := o.Owner()
		intent := o.Move(w.Map, Phermones[owner])
		if !intend.Equals(o.Point) {
			target, occupied := w.Map[]
			if !occupied {
				
			}
			dead := o.Resolve()
		}
	}
}
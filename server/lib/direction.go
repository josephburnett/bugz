package colony

import "math/rand"

var D_UP = Direction{0, 1}
var D_UP_RIGHT = Direction{1, 1}
var D_RIGHT = Direction{1, 0}
var D_DOWN_RIGHT = Direction{1, -1}
var D_DOWN = Direction{0, -1}
var D_DOWN_LEFT = Direction{-1, -1}
var D_LEFT = Direction{-1, 0}
var D_UP_LEFT = Direction{-1, 1}

var D_IN_FRONT = map[Direction][]Direction{
	D_UP:         []Direction{D_UP_LEFT, D_UP, D_UP_RIGHT},
	D_UP_RIGHT:   []Direction{D_UP, D_UP_RIGHT, D_RIGHT},
	D_RIGHT:      []Direction{D_UP_RIGHT, D_RIGHT, D_DOWN_RIGHT},
	D_DOWN_RIGHT: []Direction{D_RIGHT, D_DOWN_RIGHT, D_DOWN},
	D_DOWN:       []Direction{D_DOWN_RIGHT, D_DOWN, D_DOWN_LEFT},
	D_DOWN_LEFT:  []Direction{D_DOWN, D_DOWN_LEFT, D_LEFT},
	D_LEFT:       []Direction{D_DOWN_LEFT, D_LEFT, D_UP_LEFT},
	D_UP_LEFT:    []Direction{D_LEFT, D_UP_LEFT, D_UP},
}

var D_OPPOSITE = map[Direction]Direction{
	D_UP:         D_DOWN,
	D_UP_RIGHT:   D_DOWN_LEFT,
	D_RIGHT:      D_LEFT,
	D_DOWN_RIGHT: D_UP_LEFT,
	D_DOWN:       D_UP,
	D_DOWN_LEFT:  D_UP_RIGHT,
	D_LEFT:       D_RIGHT,
	D_UP_LEFT:    D_DOWN_RIGHT,
}

var D_AROUND []Direction = make([]Direction, 0, 8)

func init() {
	for d, _ := range D_IN_FRONT {
		D_AROUND = append(D_AROUND, d)
	}
}

var D_SURROUNDING []Direction = make([]Direction, 0, 9)

func init() {
	for d, _ := range D_IN_FRONT {
		D_SURROUNDING = append(D_SURROUNDING, d)
	}
	D_SURROUNDING = append(D_SURROUNDING, Direction{0, 0})
}

func (d Direction) InFront() []Direction {
	return D_IN_FRONT[d]
}

func (d Direction) Opposite() Direction {
	return D_OPPOSITE[d]
}

func Around() []Direction {
	return D_AROUND
}

func Surrounding() []Direction {
	return D_SURROUNDING
}

func RandomDirection(d []Direction) Direction {
	if len(d) == 1 {
		return d[0]
	}
	return d[rand.Intn(len(d))]
}

func (p1 Point) DistanceFrom(p2 Point) int {
	xDistance := p1[0] - p2[0]
	if xDistance < 0 {
		xDistance = -xDistance
	}
	yDistance := p1[1] - p2[1]
	if yDistance < 0 {
		yDistance = -yDistance
	}
	if xDistance > yDistance {
		return xDistance
	} else {
		return yDistance
	}
}

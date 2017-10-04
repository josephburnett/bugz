package colony

import (
	"encoding/gob"

	"github.com/josephburnett/colony/server/proto/view"
)

var _ Object = &Fruit{}

func init() {
	gob.Register(&Fruit{})
}

type Fruit struct {
	Freshness int
	Time      int
}

func NewFruit() *Fruit {
	return &Fruit{
		Freshness: 3,
		Time:      0,
	}
}

func (f *Fruit) Owner() Owner {
	return Owner("")
}

func (f *Fruit) Tick() {
	if f.Time == 100 {
		f.Freshness = f.Freshness - 1
		f.Time = 0
	} else {
		f.Time = f.Time + 1
	}
}

func (f *Fruit) Dead() bool {
	return f.Freshness < 1
}

func (f *Fruit) View(Owner) *view.Object {
	return &view.Object{
		Type:      "fruit",
		Direction: &view.Coordinate{0, 0},
		Mine:      false,
	}
}

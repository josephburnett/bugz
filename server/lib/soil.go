package colony

import (
	"encoding/gob"
	"log"
)

var _ ProducerObject = &Soil{}

func init() {
	gob.Register(&Soil{})
}

type Soil struct {
	Richness int
	Time     int
}

func (s *Soil) Owner() Owner {
	return Owner("")
}

func (s *Soil) Reclaim(_ Object) {
	if s.Richness < 3 {
		log.Println("soil is enriched")
		s.Richness = s.Richness + 1
		s.Time = 0
	}
}

func (s *Soil) Tick() {
	if s.Time == 120 {
		log.Println("soil decays")
		s.Richness = s.Richness - 1
		s.Time = 0
	} else {
		s.Time = s.Time + 1
	}
}

func (s *Soil) Dead() bool {
	return s.Richness < 1
}

func (s *Soil) Produce() (Object, bool) {
	if s.Richness == 3 && (s.Time%10 == 0) {
		log.Println("soild produces")
		return NewFruit(), true
	}
	return nil, false
}

func (s *Soil) View(_ Owner) *ObjectView {
	return &ObjectView{
		Type: "soil",
	}
}

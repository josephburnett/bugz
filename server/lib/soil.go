package colony

import "log"

var _ ProducerObject = &Soil{}

type Soil struct {
	richness int
	time     int
}

func (s *Soil) Owner() Owner {
	return Owner("")
}

func (s *Soil) Reclaim(_ Object) {
	if s.richness < 3 {
		log.Println("soil is enriched")
		s.richness = s.richness + 1
		s.time = 0
	}
}

func (s *Soil) Tick() {
	if s.time == 120 {
		log.Println("soil decays")
		s.richness = s.richness - 1
		s.time = 0
	} else {
		s.time = s.time + 1
	}
}

func (s *Soil) Dead() bool {
	return s.richness < 1
}

func (s *Soil) Produce() (Object, bool) {
	if s.richness == 3 && (s.time%10 == 0) {
		log.Println("soild produces")
		return &Fruit{
			freshness: 3,
			time:      0,
		}, true
	}
	return nil, false
}

func (s *Soil) View(_ Owner) *ObjectView {
	return &ObjectView{
		Type: "soil",
	}
}

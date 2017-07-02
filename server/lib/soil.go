package colony

import "log"

type Soil struct {
	richness int
	time     int
}

func (s *Soil) Enrich() {
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
	if s.richness == 3 && s.time == 10 {
		log.Println("soild produces")
		s.richness = 1
		return &Fruit{
			freshness: 3,
			time:      0,
		}, true
	}
	return nil, false
}

package colony

type Rock struct {
	lifetime int
}

func NewRock() *Rock {
	return &Rock{lifetime: 1000}
}

func (r *Rock) Owner() Owner {
	return Owner("")
}

func (r *Rock) Tick() {
	r.lifetime = r.lifetime - 1
}

func (r *Rock) Dead() bool {
	return r.lifetime == 0
}

func (r *Rock) View(o Owner) *ObjectView {
	return &ObjectView{
		Type: "rock",
	}
}

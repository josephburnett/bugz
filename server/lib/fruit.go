package colony

var _ Object = &Fruit{}

type Fruit struct {
	point     Point
	freshness int
	time      int
}

func (f *Fruit) Owner() Owner {
	return Owner("")
}

func (f *Fruit) Tick() {
	if f.time == 100 {
		f.freshness = f.freshness - 1
		f.time = 0
	} else {
		f.time = f.time + 1
	}
}

func (f *Fruit) Dead() bool {
	return f.freshness < 1
}

func (f *Fruit) View(Owner) *ObjectView {
	return &ObjectView{
		Type:      "fruit",
		Direction: Direction{0, 0},
		Mine:      false,
	}
}

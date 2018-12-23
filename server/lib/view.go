package colony

type ObjectView struct {
	Type      string
	Direction Direction
	Mine      bool
	Strength  int
}

type PointView struct {
	Point    Point
	Phermone bool
	Object   *ObjectView
	Earth    *ObjectView
}

type WorldView struct {
	Points     []*PointView
	Friends    map[Owner]bool
	LowerLeft  Point
	UpperRight Point
}

func (w *World) View(owner Owner) *WorldView {
	phermones := w.Phermones[owner]
	center := w.Colonies[owner]
	wv := &WorldView{
		Points:     make([]*PointView, 0),
		Friends:    make(map[Owner]bool),
		LowerLeft:  Point{center[0] - 19, center[1] - 19},
		UpperRight: Point{center[0] + 19, center[1] + 19},
	}
	for y := wv.UpperRight[1]; y >= wv.LowerLeft[1]; y-- {
		for x := wv.LowerLeft[0]; x <= wv.UpperRight[0]; x++ {
			point := Point{x, y}
			pv := &PointView{
				Point: point,
			}
			nonEmpty := false
			if object, exists := w.Objects[point]; exists {
				pv.Object = object.View(owner)
				nonEmpty = true
			}
			if producer, exists := w.Earth[point]; exists {
				pv.Earth = producer.View(owner)
				nonEmpty = true
			}
			if _, present := phermones[point]; present {
				pv.Phermone = true
				nonEmpty = true
			}
			if nonEmpty {
				wv.Points = append(wv.Points, pv)
			}
		}
	}
	friends, ok := w.Friends[owner]
	if !ok {
		friends = make(map[Owner]bool)
	}
	for o := range w.Colonies {
		if o == owner {
			continue
		}
		if _, friend := friends[o]; friend {
			wv.Friends[o] = true
		} else {
			wv.Friends[o] = false
		}
	}
	return wv
}

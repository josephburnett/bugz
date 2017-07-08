package colony

type ObjectView struct {
	Type      string
	Direction Direction
	Mine      bool
}

type PointView struct {
	Point    Point
	Phermone bool
	Object   *ObjectView
	Earth    *ObjectView
}

type WorldView struct {
	Points  [][]*PointView
	Friends map[Owner]bool
}

func (w *World) View(owner Owner) *WorldView {
	phermones := w.Phermones[owner]
	center := w.Colonies[owner]
	lowerLeft := &Point{center[0] - 19, center[1] - 19}
	upperRight := &Point{center[0] + 19, center[1] + 19}
	wv := &WorldView{
		Points:  make([][]*PointView, 0, 39),
		Friends: make(map[Owner]bool),
	}
	for y := upperRight[1]; y >= lowerLeft[1]; y-- {
		row := make([]*PointView, 0, 39)
		wv.Points = append(wv.Points, row)
		for x := lowerLeft[0]; x <= upperRight[0]; x++ {
			point := Point{x, y}
			pv := &PointView{
				Point: point,
			}
			if object, exists := w.Objects[point]; exists {
				pv.Object = object.View(owner)
			}
			if producer, exists := w.Earth[point]; exists {
				pv.Earth = producer.View(owner)
			}
			if _, present := phermones[point]; present {
				pv.Phermone = true
			}
			row = append(row, pv)
		}
		wv.Points = append(wv.Points, row)
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

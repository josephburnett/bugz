package colony

type ObjectView struct {
	Type      string
	Direction Direction
	Color     string
	Mine      bool
}

type PointView struct {
	Point    Point
	Object   *ObjectView
	Soil     int
	Phermone bool
	Colony   bool
}

type WorldView struct {
	Points  [][]*PointView
	Friends map[Owner]bool
}

func (w *World) View(owner Owner) *WorldView {
	phermones := w.phermones[owner]
	center := w.owners[owner].Point()
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
			if object, exists := w.objects[point]; exists {
				pv.Object = object.View(owner)
			}
			if soil, exists := w.soil[point]; exists {
				pv.Soil = soil.richness
			}
			if _, present := phermones[point]; present {
				pv.Phermone = true
			}
			if _, exists := w.colonies[point]; exists {
				pv.Colony = true
			}
			row = append(row, pv)
		}
		wv.Points = append(wv.Points, row)
	}
	friends, ok := w.friends[owner]
	if !ok {
		friends = make(map[Owner]bool)
	}
	for o := range w.owners {
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

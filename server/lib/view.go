package colony

type ObjectView struct {
	Direction Direction
	Color     string
	Mine      bool
}

type PointView struct {
	Point    Point
	Object   *ObjectView
	Phermone bool
	Colony   bool
}

type WorldView struct {
	Points [][]*PointView
}

type FriendsView struct {
	Friends map[Owner]bool
}

func (w *World) View(owner Owner) *WorldView {
	phermones := w.phermones[owner]
	center := w.owners[owner].Point()
	lowerLeft := &Point{center[0] - 19, center[1] - 19}
	upperRight := &Point{center[0] + 19, center[1] + 19}
	wv := &WorldView{
		Points: make([][]*PointView, 0, 39),
	}
	for y := upperRight[1]; y >= lowerLeft[1]; y-- {
		row := make([]*PointView, 0, 39)
		wv.Points = append(wv.Points, row)
		for x := lowerLeft[0]; x <= upperRight[0]; x++ {
			point := Point{x, y}
			pv := &PointView{
				Point: point,
			}
			object, exists := w.objects[point]
			if exists {
				pv.Object = object.View(owner)
			}
			_, present := phermones[point]
			if present {
				pv.Phermone = true
			}
			_, exists = w.colonies[point]
			if exists {
				pv.Colony = true
			}
			row = append(row, pv)
		}
		wv.Points = append(wv.Points, row)
	}
	return wv
}

func (w *World) FriendsView(owner Owner) *FriendsView {
	fv := &FriendsView{
		Friends: make(map[Owner]bool),
	}
	friends, ok := w.friends[owner]
	if !ok {
		friends = make(map[Owner]bool)
	}
	for o := range w.owners {
		if _, friend := friends[o]; friend {
			fv.Friends[o] = true
		} else {
			fv.Friends[o] = false
		}
	}
	return fv
}

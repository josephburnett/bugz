package colony

type ObjectView struct {
	Direction Direction
	Color     string
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
				pv.Object = object.View()
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

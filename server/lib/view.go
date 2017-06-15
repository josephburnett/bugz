package colony

type ObjectView struct {
	direction Direction
	color     string
}

type PointView struct {
	point    Point
	object   *ObjectView
	phermone bool
	colony   bool
}

type WorldView struct {
	points [][]*PointView
}

func (w *World) View(owner Owner) *WorldView {
	phermones := w.phermones[owner]
	center := w.owners[owner].Point()
	lowerLeft := &Point{center[0] - 19, center[1] - 19}
	upperRight := &Point{center[0] + 19, center[1] + 19}
	wv := &WorldView{
		points: make([][]*PointView, 39),
	}
	for y := upperRight[1]; y <= lowerLeft[1]; y-- {
		row := make([]*PointView, 39)
		wv.points = append(wv.points, row)
		for x := lowerLeft[0]; x <= upperRight[0]; x++ {
			point := Point{x, y}
			pv := &PointView{
				point: point,
			}
			object, exists := w.objects[point]
			if exists {
				pv.object = object.View()
			}
			_, present := phermones[point]
			if present {
				pv.phermone = true
			}
			colony, exists := w.colonies[point]
			if exists && colony.owner == owner {
				pv.colony = true
			}
			row = append(row, pv)
		}
	}
	return wv
}

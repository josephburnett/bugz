package colony

type AntView struct {
	direction Direction
	color     string
}

func (a *Ant) View() *AntView {
	return &AntView{
		direction: a.direction
		color: "#0f0"
	}
}

type PointView struct {
	point    Point
	ant      AntView
	phermone bool
	colony   bool
}

type WorldView struct {
	points [][]*PointView	
}

func (w *World) View(owner Owner) *WorldView {
	phermones := w.phermones[owner]
	center := w.owners[owner].Point()
	lowerLeft := &Point{center[0]-19,center[1]-19}
	upperRight := &Point{center[0]+19,center[1]+19}
	
	// create a world view and 2d array of point views
	for x := lowerLeft[0]; x <= upperRight[0]; x++ {
		for y:= lowerLeft[1]; y <= upporRight[1]; y++ {
			object, exists := w.objects[Point{x,y}]
			if exists {
				
			}
		}
	}
}

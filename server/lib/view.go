package colony

import (
	"github.com/josephburnett/colony/server/proto/view"
)

func (w *World) View(owner Owner) *view.World {
	phermones := w.Phermones[owner]
	center := w.Colonies[owner]
	lowerLeft := &Point{center[0] - 19, center[1] - 19}
	upperRight := &Point{center[0] + 19, center[1] + 19}
	wv := &view.World{
		LowerLeft: &view.Coordinate{
			int32(lowerLeft[0]),
			int32(lowerLeft[1]),
		},
		UpperRight: &view.Coordinate{
			int32(upperRight[0]),
			int32(upperRight[1]),
		},
		Points:   make([]*view.Point, 0, 39*39),
		Colonies: make([]*view.Colony, 0),
	}
	for y := upperRight[1]; y >= lowerLeft[1]; y-- {
		for x := lowerLeft[0]; x <= upperRight[0]; x++ {
			pv := &view.Point{
				Point: &view.Coordinate{
					X: int32(x),
					Y: int32(y),
				},
			}
			point := Point{x, y}
			if object, exists := w.Objects[point]; exists {
				pv.Object = object.View(owner)
			}
			if producer, exists := w.Earth[point]; exists {
				pv.Earth = producer.View(owner)
			}
			if _, present := phermones[point]; present {
				pv.Phermone = true
			}
			wv.Points = append(wv.Points, pv)
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
		c := &view.Colony{
			Owner: string(o),
		}
		if _, friend := friends[o]; friend {
			c.Friend = true
		} else {
			c.Friend = false
		}
		wv.Colonies = append(wv.Colonies, c)
	}
	return wv
}

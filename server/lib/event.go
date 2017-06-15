package colony

import "log"

type Event interface {
	isEvent()
}

func EventLoop(w *World) chan Event {
	ch := make(chan Event, 100)
	go func() {
		for {
			event, ok := <-ch
			if !ok {
				return
			}
			switch e := event.(type) {
			default:
				log.Println("[ERROR] unknown event")
			case *TickEvent:
				log.Println("tick")
				w.Produce()
				w.Advance()
			case *UiProduceEvent:
				c := w.owners[e.owner]
				c.produce = true
			case *UiPhermoneEvent:
				p := w.phermones[e.owner]
				if e.state {
					p[e.point] = e.state
				} else {
					delete(p, e.point)
				}
			}
		}
	}()
	return ch
}

type TickEvent struct{}

func (e *TickEvent) isEvent() {}

type UiProduceEvent struct {
	owner Owner
}

func (e *UiProduceEvent) isEvent() {}

type UiPhermoneEvent struct {
	owner Owner
	point Point
	state bool
}

func (e *UiPhermoneEvent) isEvent() {}

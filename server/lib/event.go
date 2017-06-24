package colony

import (
	"log"
)

type EventType string

const (
	E_UI_PRODUCE   = EventType("ui-produce")
	E_UI_PHERMONE  = EventType("ui-phermone")
	E_TIME_TICK    = EventType("time-tick")
	E_VIEW_REQUEST = EventType("view-request")
	E_VIEW_UPDATE  = EventType("view-update")
)

type Event interface {
	eventType() EventType
}

type EventLoop struct {
	C       chan Event
	World   *World
	viewers map[Owner][]chan *WorldView
}

func (e *EventLoop) View(o Owner, c chan *WorldView) {
	viewers, ok := e.viewers[o]
	if !ok {
		viewers = make([]chan *WorldView, 0)
	}
	e.viewers[o] = append(viewers, c)
}

func (e *EventLoop) Unview(o Owner, c chan *WorldView) {
	viewers, ok := e.viewers[o]
	if !ok {
		return
	}
	for i, ch := range viewers {
		if ch == c {
			e.viewers[o] = append(viewers[:i], viewers[i+1:]...)
			return
		}
	}
}

func (e *EventLoop) BroadcastView() {
	for o, viewers := range e.viewers {
		view := e.World.View(o)
		for _, viewer := range viewers {
			viewer <- view
		}
	}
}

func NewEventLoop(w *World) (e *EventLoop) {
	e = &EventLoop{
		C:       make(chan Event, 100),
		World:   w,
		viewers: make(map[Owner][]chan *WorldView),
	}
	go func() {
		for {
			event, ok := <-e.C
			if !ok {
				return
			}
			switch event := event.(type) {
			default:
				log.Println("[ERROR] unknown event")
			case *TimeTickEvent:
				w.Advance()
				e.BroadcastView()
			case *UiProduceEvent:
				colony := w.owners[event.Owner]
				colony.produce = true
			case *UiPhermoneEvent:
				p := w.phermones[event.Owner]
				if event.State {
					p[event.Point] = event.State
				} else {
					delete(p, event.Point)
				}
			}
		}
	}()
	return
}

type TimeTickEvent struct{}

func (e *TimeTickEvent) eventType() EventType { return E_TIME_TICK }

type UiProduceEvent struct {
	Owner Owner
}

func (e *UiProduceEvent) eventType() EventType { return E_UI_PRODUCE }

type UiPhermoneEvent struct {
	Owner Owner
	Point Point
	State bool
}

func (e *UiPhermoneEvent) eventType() EventType { return E_UI_PHERMONE }

type ViewUpdateEvent struct {
	Owner     Owner
	WorldView *WorldView
}

func (e *ViewUpdateEvent) eventType() EventType { return E_VIEW_UPDATE }

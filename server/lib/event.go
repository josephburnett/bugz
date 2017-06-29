package colony

import (
	"errors"
	"fmt"
	"log"
	"time"
)

type EventType string

const (
	E_UI_PRODUCE  = EventType("ui-produce")
	E_UI_PHERMONE = EventType("ui-phermone")
	E_UI_CONNECT  = EventType("ui-connect")
	E_UI_FRIEND   = EventType("ui-friend")
	E_TIME_TICK   = EventType("time-tick")
	E_VIEW_UPDATE = EventType("view-update")
)

type Event interface {
	eventType() EventType
}

type EventLoop struct {
	C       chan Event
	World   *World
	viewers map[Owner]chan *WorldView
}

func (e *EventLoop) View(o Owner, c chan *WorldView) {
	e.viewers[o] = c
}

func (e *EventLoop) Unview(o Owner) {
	if c, ok := e.viewers[o]; ok {
		close(c)
	}
	delete(e.viewers, o)
}

func (e *EventLoop) BroadcastView() {
	for o, viewer := range e.viewers {
		if _, exists := e.World.owners[o]; !exists {
			// TODO: move this to a viewer cleanup routine
			delete(e.viewers, o)
			return
		}
		viewer <- e.World.View(o)
	}
}

func NewEventLoop(w *World) (e *EventLoop) {
	e = &EventLoop{
		C:       make(chan Event, 100),
		World:   w,
		viewers: make(map[Owner]chan *WorldView),
	}
	go func() {
		for {
			event, ok := <-e.C
			if !ok {
				log.Println("[ERROR] event loop channel closed")
				return
			}
			switch event := event.(type) {
			default:
				log.Println("[ERROR] unknown event")
			case *TimeTickEvent:
				w.Advance()
				e.BroadcastView()
			case *UiConnectEvent:
				if _, exists := w.owners[event.Owner]; !exists {
					w.NewColony(event.Owner)
				}
			case *UiProduceEvent:
				colony, ok := w.owners[event.Owner]
				if ok {
					colony.produce = true
				}
			case *UiPhermoneEvent:
				p, ok := w.phermones[event.Owner]
				if !ok {
					log.Println("[ERROR] phermone event for unknown owner: ", event.Owner)
					continue
				}
				if event.State {
					p[event.Point] = event.State
				} else {
					delete(p, event.Point)
				}
			case *UiFriendEvent:
				if event.State {
					w.Friend(event.Owner, event.Friend)
				} else {
					w.Unfriend(event.Owner, event.Friend)
				}
			}
		}
	}()
	go func() {
		t := time.NewTicker(1000 * time.Millisecond)
		defer t.Stop()
		for {
			_, ok := <-t.C
			if !ok {
				log.Println("[ERROR] time ticker channel closed")
				return
			}
			e.C <- &TimeTickEvent{}
		}
	}()
	return
}

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

type UiConnectEvent struct {
	Owner Owner
}

func (e *UiConnectEvent) eventType() EventType { return E_UI_CONNECT }

type UiFriendEvent struct {
	Owner  Owner
	Friend Owner
	State  bool
}

func (e *UiFriendEvent) eventType() EventType { return E_UI_FRIEND }

type TimeTickEvent struct{}

func (e *TimeTickEvent) eventType() EventType { return E_TIME_TICK }

type ViewUpdateEvent struct {
	Owner     Owner
	WorldView *WorldView
}

func (e *ViewUpdateEvent) eventType() EventType { return E_VIEW_UPDATE }

func UnmarshalEvent(t EventType, event map[string]interface{}) (Event, error) {
	switch t {
	case E_UI_CONNECT:
		owner, ok := event["Owner"].(string)
		if !ok {
			return nil, errors.New("Connect event from client does not have owner")
		}
		return &UiConnectEvent{Owner: Owner(owner)}, nil
	case E_UI_PRODUCE:
		owner, ok := event["Owner"].(string)
		if !ok {
			return nil, errors.New("Produce event from client does not have owner")
		}
		return &UiProduceEvent{Owner: Owner(owner)}, nil
	case E_UI_PHERMONE:
		owner, ok := event["Owner"].(string)
		if !ok {
			return nil, errors.New("Phermone event from client does not have owner")
		}
		point, ok := event["Point"].([]interface{})
		if !ok {
			return nil, errors.New("Phermone event from client does not have point")
		}
		if len(point) != 2 {
			return nil, errors.New("Point must have exactly two values")
		}
		x, ok := point[0].(float64)
		if !ok {
			return nil, fmt.Errorf("Point x must be an number (%T)", point[0])
		}
		y, ok := point[1].(float64)
		if !ok {
			return nil, fmt.Errorf("Point y must be an number (%T)", point[0])
		}
		state, ok := event["State"].(bool)
		if !ok {
			return nil, errors.New("Phermone event from client does not have state")
		}
		return &UiPhermoneEvent{
			Owner: Owner(owner),
			Point: Point([2]int{int(x), int(y)}),
			State: state,
		}, nil
	case E_UI_FRIEND:
		owner, ok := event["Owner"].(string)
		if !ok {
			return nil, errors.New("Friend event from client does not have owner")
		}
		friend, ok := event["Friend"].(string)
		if !ok {
			return nil, errors.New("Friend event from client does not have friend")
		}
		state, ok := event["State"].(bool)
		if !ok {
			return nil, errors.New("Friend event from client does not have state")
		}
		return &UiFriendEvent{
			Owner:  Owner(owner),
			Friend: Owner(friend),
			State:  state,
		}, nil
	default:
		return nil, errors.New("Unknown message type from client")
	}
}

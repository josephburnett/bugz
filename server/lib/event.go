package colony

import (
	"log"
	"time"

	"github.com/josephburnett/colony/server/proto/event"
	"github.com/josephburnett/colony/server/proto/view"
)

type EventLoop struct {
	C       chan *event.Event
	World   *World
	viewers map[Owner]chan *view.World
}

func (e *EventLoop) View(o Owner, c chan *view.World) {
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
		if _, exists := e.World.Colonies[o]; !exists {
			// TODO: move this to a viewer cleanup routine
			delete(e.viewers, o)
			return
		}
		viewer <- e.World.View(o)
	}
}

func NewEventLoop(w *World) (e *EventLoop) {
	e = &EventLoop{
		C:       make(chan *event.Event, 100),
		World:   w,
		viewers: make(map[Owner]chan *view.World),
	}
	go func() {
		for {
			ev, ok := <-e.C
			if !ok {
				log.Println("event loop channel closed")
				return
			}
			owner := Owner(ev.Owner)
			if owner == "" {
				log.Println("invalid event, missing owner")
				continue
			}
			switch ev.GetEvent().(type) {
			default:
				log.Println("unknown event")
			// System events
			case *event.Event_Tick:
				w.Advance()
				e.BroadcastView()
			case *event.Event_SaveWorld:
				saveWorld := ev.GetSaveWorld()
				err := w.SaveWorld(saveWorld.Filename)
				if err != nil {
					log.Println("error saving the world", err.Error())
				}
			// User events
			case *event.Event_Connect:
				if _, exists := w.Colonies[owner]; !exists {
					w.NewColony(owner)
				}
			case *event.Event_Produce:
				if colony, ok := w.FindColony(owner); ok {
					colony.Touch()
					colony.P = true
				} else {
					log.Println("ignoring produce event for colony not on the earth")
				}
			case *event.Event_Move:
				if colony, ok := w.FindColony(owner); ok {
					colony.Touch()
					delete(w.Earth, colony.Center())
					w.Objects[colony.Center()] = NewQueen(colony)
				} else {
					log.Println("ignoring move event for colony not on the earth")
				}
			case *event.Event_Phermone:
				phermone := ev.GetPhermone()
				if colony, ok := w.FindColony(owner); ok {
					colony.Touch()
				}
				p, ok := w.Phermones[owner]
				if !ok {
					log.Println("phermone event for unknown colony", owner)
					continue
				}
				if phermone.Point == nil {
					log.Println("phermone event missing coordinate")
					continue
				}
				point := Point{int(phermone.Point.X), int(phermone.Point.Y)}
				if phermone.State {
					p[point] = phermone.State
				} else {
					delete(p, point)
				}
			case *event.Event_PhermoneClear:
				w.Phermones[owner] = make(Phermones)
				log.Println("clearing phermones for " + owner)
			case *event.Event_Friend:
				friend := ev.GetFriend()
				if friend.Friend == "" {
					log.Println("friend event missing friend")
					continue
				}
				if colony, ok := w.FindColony(owner); ok {
					colony.Touch()
				}
				if friend.State {
					w.Friend(owner, Owner(friend.Friend))
				} else {
					w.Unfriend(owner, Owner(friend.Friend))
				}
			case *event.Event_Drop:
				drop := ev.GetDrop()
				if drop.What == "" {
					log.Println("drop event missing what")
					continue
				}
				if colony, ok := w.FindColony(owner); ok {
					colony.Touch()
				}
				if _, ok := w.Colonies[owner]; ok {
					w.Drop(owner, drop.What)
				} else {
					log.Println("drop event for unknown colony", owner)
				}
			}
		}
	}()
	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()
		for {
			_, ok := <-t.C
			if !ok {
				log.Println("time ticker channel closed")
				return
			}
			e.C <- &event.Event{
				Owner: "",
				Event: &event.Event_Tick{},
			}
		}
	}()
	return
}

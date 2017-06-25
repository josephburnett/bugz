package colony

import (
	"time"
)

func Serve(w *World) {
	e := NewEventLoop(w)
	c := NewClients(e)
	c.Serve("0.0.0.0:8080")
	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()
		for {
			_, ok := <-t.C
			if !ok {
				return
			}
			e.C <- &TimeTickEvent{}
		}
	}()
	done := make(chan struct{})
	<-done
}

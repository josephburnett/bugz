package colony

import (
	"net/http"
	"time"
)

func Serve(w *World) {
	e := EventLoop(w)
	defer close(e)
	// Testing
	go func() {
		t := time.NewTicker(2000 * time.Millisecond)
		defer t.Stop()
		for {
			_, ok := <-t.C
			if !ok {
				return
			}
			e <- &UiProduceEvent{
				owner: Owner("joe"),
			}
		}
	}()
	c := NewClients()
	c.Serve(w, "0.0.0.0:8080")
}

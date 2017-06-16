package colony

import (
	"net/http"
	"time"
)

func Serve(w *World) {
	e := EventLoop(w)
	defer close(e)
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
	go func() {
		t := time.NewTicker(500 * time.Millisecond)
		defer t.Stop()
		for {
			_, ok := <-t.C
			if !ok {
				return
			}
			e <- &TickEvent{}
		}
	}()
	http.HandleFunc("/ws", ClientWebsocket)
	http.ListenAndServe("0.0.0.0:8081", nil)
}

package colony

import "time"

func Serve(w *World) {
	e := EventLoop(w)
	defer close(e)
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
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
	for {
		_, ok := <-t.C
		if !ok {
			return
		}
		e <- &TickEvent{}
	}
}

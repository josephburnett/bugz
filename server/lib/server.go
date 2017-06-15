package colony

import "time"

func Serve() {
	w := NewWorld()
	e := EventLoop(w)
	defer close(e)
	t := time.NewTicker(500 * time.Millisecond)
	defer t.Stop()
	for {
		_, ok := <-t.C
		if !ok {
			return
		}
		e <- &TickEvent{}
	}
}

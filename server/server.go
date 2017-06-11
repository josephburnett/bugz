package colony

import "time"

func serve() {
	w := NewWorld()
	e := EventLoop(w)
	defer close(e)
	t := time.NewTicker(250 * time.Millisecond)
	defer t.Stop()
	for {
		_, ok := <-t.C
		if !ok {
			return
		}
		e <- &TickEvent{}
	}
}

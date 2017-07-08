package main

import (
	"flag"
	"time"

	colony "github.com/josephburnett/colony/server/lib"
)

var worldFile = flag.String("world_file", "", "File for persistent world state.")

func main() {
	flag.Parse()
	var w *colony.World
	var err error
	if *worldFile != "" {
		w, err = colony.LoadWorld(*worldFile)
		if err != nil {
			panic(err)
		}
	} else {
		w = colony.NewWorld()
	}
	e := colony.NewEventLoop(w)
	if *worldFile != "" {
		go func() {
			t := time.NewTicker(30 * time.Second)
			defer t.Stop()
			for {
				_, ok := <-t.C
				if !ok {
					panic("Error with save world timer.")
				}
				e.C <- &colony.SaveWorldEvent{Filename: *worldFile}
			}
		}()
	}
	c := colony.NewClients(e)
	c.Serve("0.0.0.0:8080")
	done := make(chan struct{})
	<-done
}

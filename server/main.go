package main

import colony "github.com/josephburnett/colony/server/lib"

func main() {
	w := colony.NewWorld()
	e := colony.NewEventLoop(w)
	c := colony.NewClients(e)
	c.Serve("0.0.0.0:8080")
	done := make(chan struct{})
	<-done
}

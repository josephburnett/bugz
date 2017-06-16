package main

import colony "github.com/josephburnett/colony/server/lib"

func main() {
	w := colony.NewWorld()
	w.NewColony(colony.Owner("joe"))
	colony.Serve(w)
}

package main

import colony "github.com/josephburnett/colony/server/lib"

func main() {
	w := colony.NewWorld()
	colony.Serve(w)
}

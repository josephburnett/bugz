package main

import (
	"flag"
	"net/http"
	"time"

	colony "github.com/josephburnett/colony/server/lib"
)

func main() {
	flag.Parse()
	var w *colony.World
	var err error
	if *colony.Config.WorldFile != "" {
		w, err = colony.LoadWorld(*colony.Config.WorldFile)
		if err != nil {
			panic(err)
		}
	} else {
		w = colony.NewWorld()
	}
	e := colony.NewEventLoop(w)
	if *colony.Config.WorldFile != "" {
		go func() {
			t := time.NewTicker(30 * time.Second)
			defer t.Stop()
			for {
				_, ok := <-t.C
				if !ok {
					panic("Error with save world timer.")
				}
				e.C <- &colony.SaveWorldEvent{Filename: *colony.Config.WorldFile}
			}
		}()
	}
	c := colony.NewClients(e)
	c.Serve(*colony.Config.Ip+":"+*colony.Config.Port, colony.Handler(AssetHandler))
	done := make(chan struct{})
	<-done
}

// This is only here because go-bindata generates Asset in the main package.
func AssetHandler(suffix, contentType string) func(w http.ResponseWriter, r *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		data, err := Asset(r.URL.Path[1:] + suffix)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		w.Header().Add("Content-Type", contentType)
		w.Write(data)
	}
}

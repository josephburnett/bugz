package main

import (
	"flag"
	"net/http"
	"time"

	colony "github.com/josephburnett/colony/server/lib"
)

var worldFile = flag.String("world_file", "", "File for persistent world state.")
var ip = flag.String("ip", "0.0.0.0", "HTTP server ip.")
var port = flag.String("port", "8080", "HTTP server port.")

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
	c.Serve(*ip+":"+*port, colony.Handler(AssetHandler))
	done := make(chan struct{})
	<-done
}

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

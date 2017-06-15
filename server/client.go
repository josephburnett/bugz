package colony

import (
	"log"
	"net/http"
)

type ClientMessage struct {
	uiProduceEvent  UiProduceEvent
	uiPhermoneEVent UiPhermoneEvent
}

type ServerMessage struct {
	
}

var upgrader = websocket.Upgrader{}

func conn(w http.ResponseWriter, r *http.Request) {
	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("upgrade:", err)
		return
	}
	defer c.Close()
	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}
		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			break
		}
	}
}
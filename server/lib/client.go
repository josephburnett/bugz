package colony

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Clients struct {
	clients        map[Owner][]chan *ServerMessage
	clientMessages chan *ClientMessage
}

func NewClients() *Clients {
	return &Clients{
		clients:        make(map[Owner][]chan *ServerMessage),
		clientMessages: make(chan *ClientMessage),
	}
}

func (c *Clients) BroadcastAll(msg *ServerMessage) {
	for _, clients := range c.clients {
		for client := range clients {
			client <- msg
		}
	}
}

func (c *Clients) BroadcastOwner(o Owner) {
	clients, exist := c.clients[o]
	if exist {
		for client := range clients {
			client <- msg
		}
	}
}

func (c *Clients) Serve(w *World, addr string) {
	ClientWebsocket := func(w http.ResponseWriter, r *http.Request) {
		c, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println("Error while upgrading: ", err)
			return
		}
		defer c.Close()
		done := make(chan struct{})
		clientChan = make(chan *ServerMessage, 100)
		defer close(clientChan)
		clients[clientChan] = make(struct{})
		defer delete(clients, clientChan)
		go func() {
			for {
				msg, ok := <-serverMessage
				if !ok {
					close(done)
					return
				}
				err = c.WriteJson(msg)
				if err != nil {
					log.Println("Error while writing to client: ", err)
					close(done)
					return
				}
			}
		}()
		go func() {
			for {
				msg := &ClientMessage{}
				err := c.ReadJson(msg)
				if err != nil {
					log.Println("Error while reading from client: ", err)
					close(done)
					return
				}
				clientMessages <- msg
			}
		}()
		<-done
	}
	http.HandleFunc("/ws", ClientWebsocket)
	http.ListenAndServe("0.0.0.0:8081", nil)
}

// go func() {
// 	for {
// 		msg, ok := <-ClientMessages
// 		if !ok {
// 			return
// 		}
// 		if msg.UiProduceEvent != nil {
// 			ch <- msg.UiProduceEvent
// 		}
// 		if msg.UiPhermoneEvent != nil {
// 			ch <- msg.UiPhermoneEvent
// 		}
// 	}
// }()
// return ch

type ClientMessage struct {
	UiProduceEvent  *UiProduceEvent
	UiPhermoneEVent *UiPhermoneEvent
}

type ServerMessage struct {
	WorldView *WorldView
}

var upgrader = websocket.Upgrader{}

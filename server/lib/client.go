package colony

import (
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Clients struct {
	clients        map[Owner][]chan *Message
	clientMessages chan *Message
	eventLoop      *EventLoop
}

func NewClients(e *EventLoop) *Clients {
	return &Clients{
		clients:        make(map[Owner][]chan *Message),
		clientMessages: make(chan *Message),
		eventLoop:      e,
	}
	// TODO listen to view updates
}

func (c *Clients) BroadcastAll(msg *Message) {
	for _, clients := range c.clients {
		for _, client := range clients {
			client <- msg
		}
	}
}

func (c *Clients) BroadcastOwner(o Owner, msg *Message) {
	clients, exist := c.clients[o]
	if exist {
		for _, client := range clients {
			client <- msg
		}
	}
}

func (c *Clients) Connect(o Owner, ch chan *Message) {
	clients, exists := c.clients[o]
	if !exists {
		clients = make([]chan *Message, 0)
	}
	c.clients[o] = append(clients, ch)
}

func (c *Clients) Disconnect(o Owner, ch chan *Message) {
	clients, exists := c.clients[o]
	if !exists {
		return
	}
	for i, client := range clients {
		if client == ch {
			c.clients[o] = append(clients[:i], clients[i+1:]...)
			return
		}
	}
}

func (c *Clients) Serve(addr string) {
	go func() {
		ClientWebsocket := func(w http.ResponseWriter, r *http.Request) {
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				log.Println("Error while upgrading: ", err)
				return
			}
			defer conn.Close()
			clientChan := make(chan *Message, 10)
			defer close(clientChan)
			c.Connect("joe", clientChan)
			defer c.Disconnect("joe", clientChan)
			done := make(chan struct{})
			go func() {
				for {
					msg, ok := <-clientChan
					if !ok {
						close(done)
						return
					}
					err = conn.WriteJSON(msg)
					if err != nil {
						log.Println("Error while writing to client: ", err)
						close(done)
						return
					}
				}
			}()
			go func() {
				for {
					msg := &Message{}
					err := conn.ReadJSON(msg)
					if err != nil {
						log.Println("Error while reading from client: ", err)
						close(done)
						return
					}
					c.eventLoop.C <- *(msg.Event)
				}
			}()
			<-done
		}
		http.HandleFunc("/ws", ClientWebsocket)
		http.ListenAndServe("0.0.0.0:8081", nil)
	}()
}

type Message struct {
	Type  EventType
	Event *Event
}

var upgrader = websocket.Upgrader{}

package colony

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
)

type Clients struct {
	clients       map[Owner][]chan *Message
	viewingOwners map[Owner]chan *WorldView
	eventLoop     *EventLoop
}

func NewClients(e *EventLoop) *Clients {
	return &Clients{
		clients:       make(map[Owner][]chan *Message),
		viewingOwners: make(map[Owner]chan *WorldView),
		eventLoop:     e,
	}
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
	if _, ok := c.viewingOwners[o]; !ok {
		// TODO: discard overflow to prevent slow clients from blocking engine
		ch := make(chan *WorldView, 10)
		go func() {
			for {
				view := <-ch
				event := &ViewUpdateEvent{
					Owner:     o,
					WorldView: view,
				}
				msg := &Message{
					Type:  event.eventType(),
					Event: event,
				}
				for _, client := range c.clients[o] {
					client <- msg
				}
			}
		}()
		c.eventLoop.View(o, ch)
	}
}

func (c *Clients) Disconnect(o Owner, ch chan *Message) {
	clients, exists := c.clients[o]
	if !exists {
		return
	}
	for i, client := range clients {
		if client == ch {
			c.clients[o] = append(clients[:i], clients[i+1:]...)
			if len(c.clients[o]) == 0 {
				c.eventLoop.Unview(o, c.viewingOwners[o])
			}
			return
		}
	}
}

func (c *Clients) Serve(addr string) {
	go func() {
		ClientWebsocket := func(w http.ResponseWriter, r *http.Request) {
			log.Println("Connected client")
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
			done := make(chan bool)
			go func() {
				for {
					msg, ok := <-clientChan
					if !ok {
						done <- true
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
					_, message, err := conn.ReadMessage()
					if err != nil {
						log.Println("Error while reading from client: ", err)
						continue
					}
					var msg interface{}
					if err = json.Unmarshal(message, &msg); err != nil {
						log.Println("Error while deserializing message from client: ", string(message))
						continue
					}
					msgMap, ok := msg.(map[string]interface{})
					if !ok {
						log.Println("Unexpected structure of message from client: ", string(message))
						continue
					}
					msgType, ok := msgMap["Type"].(string)
					if !ok {
						log.Println("Message from client does not have Type: ", string(message))
						continue
					}
					msgEvent, ok := msgMap["Event"].(map[string]interface{})
					if !ok {
						log.Println("Message from client does not have event: ", string(message))
						continue
					}
					event, err := UnmarshalEvent(EventType(msgType), msgEvent)
					if err != nil {
						log.Println(err.Error, ": ", string(message))
						continue
					}
					c.eventLoop.C <- event
				}
			}()
			<-done
			log.Println("Disconnected client")
		}
		http.HandleFunc("/ws", ClientWebsocket)
		http.ListenAndServe(addr, nil)
	}()
}

type Message struct {
	Type  EventType
	Event Event
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

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
				view, ok := <-ch
				if !ok {
					return
				}
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
				c.eventLoop.Unview(o)
			}
			return
		}
	}
}

type Handler func(string, string) func(http.ResponseWriter, *http.Request)

func (c *Clients) Serve(addr string, assetHandler Handler) {
	go func() {
		ClientWebsocket := func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					log.Println("Recovered from panic in websocker handler: ", r)
				}
			}()
			owner := r.URL.Path[len("/ws/owner/"):]
			if owner == "" {
				http.Error(w, "Owner is a required parameter", http.StatusBadRequest)
				return
			}
			conn, err := upgrader.Upgrade(w, r, nil)
			if err != nil {
				http.Error(w, "Error while upgrading to websocket: "+err.Error(), http.StatusInternalServerError)
				return
			}
			log.Println("Connected client: ", owner)
			defer conn.Close()
			// Create Colony if necessary
			c.eventLoop.C <- &UiConnectEvent{
				Owner: Owner(owner),
			}
			// Register channel for communication with the client
			clientChan := make(chan *Message, 10)
			defer close(clientChan)
			c.Connect(Owner(owner), clientChan)
			defer c.Disconnect(Owner(owner), clientChan)
			done := make(chan bool)
			// Pipeline messages to client
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
			// Unwrap message to server
			go func() {
				for {
					_, message, err := conn.ReadMessage()
					if err != nil {
						log.Println("Error while reading from client: ", err)
						return
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
						log.Println(err.Error(), ": ", string(message))
						continue
					}
					c.eventLoop.C <- event
				}
			}()
			<-done
			log.Println("Disconnected client")
		}
		http.HandleFunc("/", assetHandler("index.html", "text/html"))
		http.HandleFunc("/css/", assetHandler("", "text/css"))
		http.HandleFunc("/js/config.js", ClientConfig)
		http.HandleFunc("/js/", assetHandler("", "text/javascript"))
		http.HandleFunc("/ws/owner/", ClientWebsocket)
		http.ListenAndServe(addr, nil)
	}()
}

func ClientConfig(w http.ResponseWriter, r *http.Request) {
	config, err := ConfigJson()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Content-Type", "text/javascript")
	w.Write([]byte("CONFIG = "))
	w.Write(config)
}

type Message struct {
	Type  EventType
	Event Event
}

var upgrader = websocket.Upgrader{}

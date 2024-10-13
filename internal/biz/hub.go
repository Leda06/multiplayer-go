package biz

import (
	"encoding/json"
)

type WSMessageType string

const (
	WSMessageTypePing             WSMessageType = "ping"
	WSMessageTypeGetClientCount   WSMessageType = "get-member-count"
	WSMessageTypeClientRegister   WSMessageType = "client-register"
	WSMessageTypeClientUnregister WSMessageType = "client-unregister"
)

type WSMessageData map[string]interface{}

type WSMessage[T WSMessageData] struct {
	Type WSMessageType `json:"type"`
	Data T             `json:"data"`
}

type Hub struct {
	clients    map[*Client]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	count      int
}

func NewHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]*Client),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		count:      0,
	}
}

func (h *Hub) Register(client *Client) {
	h.register <- client
}

func (h *Hub) Run() {
	go func() {
		for {
			if msg, ok := <-h.broadcast; ok {
				for client := range h.clients {
					if !client.status.closed {
						client.send <- msg
					}
				}
			}
		}
	}()
	for {
		select {
		case client := <-h.register:
			h.clients[client] = client
			// log.Println("added client", h.clients, client)
			h.count++
			foo, _ := json.Marshal(WSMessage[WSMessageData]{
				Type: WSMessageTypeClientRegister,
				Data: WSMessageData{"count": h.count},
			})
			h.broadcast <- foo
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				// log.Println("removed client", h.clients, client)
				h.count--
				foo, _ := json.Marshal(WSMessage[WSMessageData]{
					Type: WSMessageTypeClientUnregister,
					Data: WSMessageData{"count": h.count},
				})
				h.broadcast <- foo
			}
			// case msg := <-h.broadcast:
			// 	log.Println("broadcasting message", msg)
			// 	// marshalData, _ := json.Marshal(map[string]int{"count": h.count})
			// 	log.Println("Hub count", h.count)
			// 	for client := range h.clients {
			// 		log.Println("client", client)
			// 		client.send <- msg
			// 		// client.send <- marshalData
			// 		log.Println("sent message to client", h.clients, client, msg)
			// 		// select {
			// 		// case client.send <- msg:
			// 		// 	log.Println("sent message to client", h.clients, client, msg)
			// 		// so dump bro~~~~~~~~~~~~
			// 		// it is random if hit both match case
			// 		// default:
			// 		// 	close(client.send)
			// 		// 	log.Println("disconnecting client", h.clients, client)
			// 		// 	delete(h.clients, client)
			// 		// }
			// 	}
		}
	}
}

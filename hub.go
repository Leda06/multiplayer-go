package main

import (
	"log"
)

type Hub struct {
	clients    map[*Client]bool
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
	count      int
}

func newHub() *Hub {
	return &Hub{
		clients:    make(map[*Client]bool),
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		count:      0,
	}
}

func (h *Hub) run() {
	for {
		select {
		case client := <-h.register:
			h.clients[client] = true
			log.Println("added client", h.clients, client)
			h.count++
		case client := <-h.unregister:
			if _, ok := h.clients[client]; ok {
				delete(h.clients, client)
				close(client.send)
				log.Println("removed client", h.clients, client)
				h.count--
			}
		case msg := <-h.broadcast:
			log.Println("broadcasting message", msg)
			// marshalData, _ := json.Marshal(map[string]int{"count": h.count})
			log.Println("Hub count", h.count)
			for client := range h.clients {
				log.Println("client", client)
				client.send <- msg
				log.Println("sent message to client", h.clients, client, msg)
				// select {
				// case client.send <- msg:
				// 	log.Println("sent message to client", h.clients, client, msg)
				// so dump bro~~~~~~~~~~~~
				// it is random if hit both match case
				// default:
				// 	close(client.send)
				// 	log.Println("disconnecting client", h.clients, client)
				// 	delete(h.clients, client)
				// }
			}
		}
	}
}

package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// Time allowed to write a message to the peer.
	writeWait = 10 * time.Second
	// Time allowed to read the next pong message from the peer.
	pongWait = 5 * time.Second
	// Send pings to peer with this period. Must be less than pongWait.
	pingPeriod = (pongWait * 9) / 10
	// Maximum message size allowed from peer.
	maxMessageSize = 256
)

var (
	newline = []byte{'\n'}
	space   = []byte{' '}
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:   256,
	WriteBufferSize:  256,
	HandshakeTimeout: 5 * time.Second,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type Client struct {
	hub *Hub
	// The websocket connection.
	conn *websocket.Conn
	// Buffered channel of outbound messages.
	send chan []byte
}

// https://stackoverflow.com/questions/37696527/go-gorilla-websockets-on-ping-pong-fail-user-disconnct-call-function
// https://stackoverflow.com/questions/10585355/sending-websocket-ping-pong-frame-from-browser
// browser can't send a ping frame to server
// but we can simulate by send message like {type: 'ping'} to simulate ping frame
// c.conn.SetPingHandler(func(string) error { log.Println("ping handler"); return nil })
func (c *Client) readMessage() {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close()
	}()
	c.conn.SetReadLimit(maxMessageSize)
	c.conn.SetReadDeadline(time.Now().Add(pongWait))
	c.conn.SetPongHandler(func(string) error {
		// log.Println("pong handler", s)
		c.conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("error: %v", err)
			}

			return
		}
		// Print the message to the console
		fmt.Printf("%s sent: %s\n", c.conn.RemoteAddr(), string(message))

		// message = bytes.TrimSpace(bytes.Replace(message, newline, space, -1))
		c.hub.broadcast <- message
	}
}

func (c *Client) writeMessage() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close()
	}()
	for {
		select {
		case message, ok := <-c.send:
			log.Println("writeMessage", string(message))
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				log.Println("writeMessage not ok")
				return
			}
			w, err := c.conn.NextWriter(websocket.TextMessage)
			if err != nil {
				log.Println("Error getting writer", err)
				return
			}

			// w.Write(message)
			//
			// // Add queued chat messages to the current websocket message.
			// n := len(c.send)
			// for i := 0; i < n; i++ {
			// 	w.Write(newline)
			// 	w.Write(<-c.send)
			// }

			if _, err := w.Write(message); err != nil {
				log.Println("Error writing message", err)
			}
			if err := w.Close(); err != nil {
				return
			}
		case <-ticker.C:
			// log.Println("ticker.C")
			c.conn.SetWriteDeadline(time.Now().Add(writeWait))
			// pingMessage := []byte("ping")
			if err := c.conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

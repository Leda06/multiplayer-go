package main

import (
	"flag"
	"log"
	"net/http"
	"time"
)

var addr = flag.String("addr", ":8081", "http service address")

func main() {
	flag.Parse()
	hub := newHub()
	go hub.run()
	http.HandleFunc("/", rootHandler)
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		for k, v := range r.Header {
			log.Println(k, v)
		}
		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			log.Print(err)
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		}
		client := &Client{
			hub:  hub,
			conn: conn,
			send: make(chan []byte),
		}
		hub.register <- client

		go client.readMessage()
		go client.writeMessage()
	})

	server := &http.Server{
		Addr:              *addr,
		ReadHeaderTimeout: 3 * time.Second,
	}
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
}

func rootHandler(w http.ResponseWriter, r *http.Request) {
	log.Println(r.URL)
	if r.URL.Path != "/" {
		http.Error(w, "Not found", http.StatusNotFound)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	http.ServeFile(w, r, "index.html")
}

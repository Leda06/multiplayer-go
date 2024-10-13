package biz

import (
	"log"
	"net/http"
	"time"
)

func NewServer(addr *string, hub *Hub) *http.Server {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /ping", pingHandler)
	mux.HandleFunc("GET /", rootHandler)
	mux.HandleFunc("GET /ws", wsHandler(hub))

	// Custom 404 handler
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		http.Error(w, "Not found", http.StatusNotFound)
	})
	server := &http.Server{
		Addr:              *addr,
		Handler:           mux,
		ReadHeaderTimeout: 3 * time.Second,
	}

	return server
}

func wsHandler(hub *Hub) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		// for k, v := range r.Header {
		// 	log.Println(k, v)
		// }
		conn, err := upgrader.Upgrade(w, r, w.Header())
		if err != nil {
			log.Print(err)
			http.Error(w, "Could not open websocket connection", http.StatusBadRequest)
		}
		client := NewClient(hub, conn)
		hub.Register(client)

		go client.ReadMessage()
		go client.WriteMessage()
	}
}

func pingHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}
	w.Write([]byte("pong\n"))
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
	// w.Write([]byte("Welcome to multiplayer-go\n"))
	http.ServeFile(w, r, "index.html")
}

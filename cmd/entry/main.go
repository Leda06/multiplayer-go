package main

import (
	"flag"
	"log"
	"multiplayer-go/internal/biz"
)

var addr = flag.String("addr", ":8081", "http service address")

func main() {
	flag.Parse()
	hub := biz.NewHub()
	go hub.Run()
	server := biz.NewServer(addr, hub)
	err := server.ListenAndServe()
	if err != nil {
		log.Fatal("ListenAndServe Error: ", err)
	}
	log.Println("Server ListenAndServe:", *addr)
}

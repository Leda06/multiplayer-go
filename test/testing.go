package test

import (
	"multiplayer-go/internal/biz"
	"testing"

	"github.com/gorilla/websocket"
)

func writeWSJSONMessage(t testing.TB, conn *websocket.Conn, message biz.WSMessage[biz.WSMessageData]) {
	t.Helper()
	if err := conn.WriteJSON(message); err != nil {
		t.Fatalf("could not send json message over ws connection %v", err)
	}
}

// func writeWSMessage(t testing.TB, conn *websocket.Conn, message string) {
// 	t.Helper()
// 	if err := conn.WriteMessage(websocket.TextMessage, []byte(message)); err != nil {
// 		t.Fatalf("could not send message over ws connection %v", err)
// 	}
// }

func assertContentType(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("handler returned wrong content type: got %v want %v", got, want)
	}
}

func assertString(t testing.TB, got, want string) {
	t.Helper()
	if got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

package test

import (
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"

	"multiplayer-go/internal/biz"

	"github.com/gorilla/websocket"
)

func TestPingHandler(t *testing.T) {
	hub := biz.NewHub()
	server := biz.NewServer(flag.String("addr", ":8082", "test server address"), hub)

	req, err := http.NewRequest("GET", "/ping", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	server.Handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	expected := "pong\n"
	assertString(t, res.Body.String(), expected)
}

func TestWebSocketConnection(t *testing.T) {
	hub := biz.NewHub()
	go hub.Run()
	server := httptest.NewServer(biz.NewServer(flag.String("addr", ":8082", "test server address"), hub).Handler)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/ws"
	ws, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("Failed to connect to WebSocket server: %v", err)
	}
	defer ws.Close()

	clientCountMessage := biz.WSMessage[biz.WSMessageData]{
		Type: biz.WSMessageTypeGetClientCount,
	}
	writeWSJSONMessage(t, ws, clientCountMessage)

	var clientRegisterResponse biz.WSMessage[biz.WSMessageData]
	err = ws.ReadJSON(&clientRegisterResponse)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	var clientCountResponse biz.WSMessage[biz.WSMessageData]
	err = ws.ReadJSON(&clientCountResponse)
	if err != nil {
		t.Fatalf("Failed to read response: %v", err)
	}

	log.Println("---------------------")
	log.Println(reflect.TypeOf(clientRegisterResponse.Data["count"]))
	log.Println(clientRegisterResponse.Type, clientCountMessage.Type, clientRegisterResponse.Data["count"])
	log.Println(clientCountResponse.Type, clientCountMessage.Type, clientCountResponse.Data["count"])
	log.Println("---------------------")
	if clientCountResponse.Type != clientCountMessage.Type || clientCountResponse.Data["count"] != float64(1) {
		t.Errorf("Received unexpected response: got %v want %v", clientCountResponse, clientCountMessage)
	}
}

func TestRootHandler(t *testing.T) {
	hub := biz.NewHub()
	server := biz.NewServer(flag.String("addr", ":8082", "test server address"), hub)

	req, err := http.NewRequest(http.MethodGet, "/", nil)
	if err != nil {
		t.Fatal(err)
	}

	res := httptest.NewRecorder()
	server.Handler.ServeHTTP(res, req)

	if status := res.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v", status, http.StatusOK)
	}

	contentType := res.Header().Get("Content-Type")
	assertContentType(t, contentType, "text/plain; charset=utf-8")

	expected := "Welcome to multiplayer-go\n"
	assertString(t, res.Body.String(), expected)
}

func TestMain(m *testing.M) {
	// Setup
	flag.Parse()

	// Run tests
	code := m.Run()

	// Teardown (if needed)

	os.Exit(code)
}

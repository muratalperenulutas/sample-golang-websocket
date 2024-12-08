package server

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"websocket/internal/controller"

	"github.com/gorilla/websocket"
)

func TestWebSocketHandler(t *testing.T) {
	s := &Server{
		clients: controller.NewClientController(),
	}
	server := httptest.NewServer(http.HandlerFunc(s.websocketHandler))
	defer server.Close()

	wsURL := "ws" + server.URL[len("http"):]

	conn, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	if err != nil {
		t.Fatalf("error connecting to websocket server: %v", err)
	}
	defer conn.Close()

	err = conn.WriteMessage(websocket.TextMessage, []byte("test message"))
	if err != nil {
		t.Fatalf("error writing message to websocket: %v", err)
	}

	_, message, err := conn.ReadMessage()
	if err != nil {
		t.Fatalf("error reading message from websocket: %v", err)
	}

	expected := "server timestamp: "
	if string(message[:len(expected)]) != expected {
		t.Errorf("expected message to start with %v; got %v", expected, string(message))
	}
}

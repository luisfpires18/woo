package websocket_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	ws "github.com/coder/websocket"
	wws "github.com/luisfpires18/woo/internal/websocket"
)

// fakeValidator returns a TokenValidator that always succeeds with the given playerID.
func fakeValidator(playerID int64) wws.TokenValidator {
	return func(token string) (int64, string, error) {
		return playerID, "player", nil
	}
}

func TestHub_ClientCount(t *testing.T) {
	hub := wws.NewHub()
	go hub.Run()

	if hub.ClientCount() != 0 {
		t.Errorf("initial client count: got %d, want 0", hub.ClientCount())
	}
}

func TestWebSocket_ConnectAndReceiveReady(t *testing.T) {
	hub := wws.NewHub()
	go hub.Run()

	handler := wws.NewHandler(hub, fakeValidator(42), "")
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=test-token"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := ws.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close(ws.StatusNormalClosure, "")

	// Should receive connection_ready message
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var msg wws.Message
	if err := json.Unmarshal(data, &msg); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}

	if msg.Type != wws.MsgConnectionReady {
		t.Errorf("type: got %q, want %q", msg.Type, wws.MsgConnectionReady)
	}

	// Verify hub registered the client
	time.Sleep(50 * time.Millisecond)
	if hub.ClientCount() != 1 {
		t.Errorf("client count: got %d, want 1", hub.ClientCount())
	}
}

func TestWebSocket_MissingToken(t *testing.T) {
	hub := wws.NewHub()
	go hub.Run()

	handler := wws.NewHandler(hub, fakeValidator(1), "")
	server := httptest.NewServer(handler)
	defer server.Close()

	// No token query param
	resp, err := http.Get(server.URL)
	if err != nil {
		t.Fatalf("get: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", resp.StatusCode, http.StatusUnauthorized)
	}
}

func TestWebSocket_PingPong(t *testing.T) {
	hub := wws.NewHub()
	go hub.Run()

	handler := wws.NewHandler(hub, fakeValidator(1), "")
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=t"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := ws.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close(ws.StatusNormalClosure, "")

	// Read connection_ready
	conn.Read(ctx)

	// Send ping
	ping := wws.Message{Type: "ping"}
	data, _ := json.Marshal(ping)
	if err := conn.Write(ctx, ws.MessageText, data); err != nil {
		t.Fatalf("write ping: %v", err)
	}

	// Should receive pong
	_, respData, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read pong: %v", err)
	}

	var pong wws.Message
	json.Unmarshal(respData, &pong)
	if pong.Type != wws.MsgPong {
		t.Errorf("type: got %q, want %q", pong.Type, wws.MsgPong)
	}
}

func TestWebSocket_Subscribe(t *testing.T) {
	hub := wws.NewHub()
	go hub.Run()

	handler := wws.NewHandler(hub, fakeValidator(1), "")
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=t"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := ws.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close(ws.StatusNormalClosure, "")

	// Read connection_ready
	conn.Read(ctx)

	// Subscribe to a topic
	sub := wws.Message{
		Type: "subscribe",
		Data: wws.SubscribeData{Topics: []string{"village:123"}},
	}
	data, _ := json.Marshal(sub)
	if err := conn.Write(ctx, ws.MessageText, data); err != nil {
		t.Fatalf("write subscribe: %v", err)
	}

	// Should receive subscription_confirmed
	_, respData, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var confirm wws.Message
	json.Unmarshal(respData, &confirm)
	if confirm.Type != wws.MsgSubscriptionConfirmed {
		t.Errorf("type: got %q, want %q", confirm.Type, wws.MsgSubscriptionConfirmed)
	}
}

func TestHub_SendToPlayer(t *testing.T) {
	hub := wws.NewHub()
	go hub.Run()

	handler := wws.NewHandler(hub, fakeValidator(99), "")
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=t"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := ws.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close(ws.StatusNormalClosure, "")

	// Read connection_ready
	conn.Read(ctx)
	time.Sleep(50 * time.Millisecond) // Let the hub register

	// Send a message to the player via hub
	hub.SendToPlayer(99, &wws.Message{
		Type: wws.MsgBuildComplete,
		Data: wws.BuildCompleteData{VillageID: 5, BuildingType: "food_1", NewLevel: 2},
	})

	// Should receive it
	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var msg wws.Message
	json.Unmarshal(data, &msg)
	if msg.Type != wws.MsgBuildComplete {
		t.Errorf("type: got %q, want %q", msg.Type, wws.MsgBuildComplete)
	}
}

func TestHub_BroadcastAll(t *testing.T) {
	hub := wws.NewHub()
	go hub.Run()

	handler := wws.NewHandler(hub, fakeValidator(1), "")
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=t"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := ws.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close(ws.StatusNormalClosure, "")

	conn.Read(ctx) // connection_ready
	time.Sleep(50 * time.Millisecond)

	hub.BroadcastAll(&wws.Message{
		Type: wws.MsgAnnouncement,
		Data: map[string]string{"title": "Server restarting"},
	})

	_, data, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var msg wws.Message
	json.Unmarshal(data, &msg)
	if msg.Type != wws.MsgAnnouncement {
		t.Errorf("type: got %q, want %q", msg.Type, wws.MsgAnnouncement)
	}
}

func TestWebSocket_UnknownMessageType(t *testing.T) {
	hub := wws.NewHub()
	go hub.Run()

	handler := wws.NewHandler(hub, fakeValidator(1), "")
	server := httptest.NewServer(handler)
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "?token=t"

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, _, err := ws.Dial(ctx, wsURL, nil)
	if err != nil {
		t.Fatalf("dial: %v", err)
	}
	defer conn.Close(ws.StatusNormalClosure, "")

	conn.Read(ctx) // connection_ready

	// Send unknown type
	unknown := wws.Message{Type: "hack_server"}
	data, _ := json.Marshal(unknown)
	conn.Write(ctx, ws.MessageText, data)

	// Should receive error
	_, respData, err := conn.Read(ctx)
	if err != nil {
		t.Fatalf("read: %v", err)
	}

	var msg wws.Message
	json.Unmarshal(respData, &msg)
	if msg.Type != wws.MsgError {
		t.Errorf("type: got %q, want %q", msg.Type, wws.MsgError)
	}
}

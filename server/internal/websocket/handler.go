package websocket

import (
	"context"
	"log/slog"
	"net/http"
	"time"

	ws "github.com/coder/websocket"
)

// TokenValidator is a function that validates a JWT and returns the playerID.
type TokenValidator func(tokenString string) (playerID int64, role string, err error)

// Handler handles WebSocket upgrade requests.
type Handler struct {
	hub           *Hub
	validator     TokenValidator
	allowedOrigin string
}

// NewHandler creates a new WebSocket handler.
// If allowedOrigin is empty, all origins are accepted (dev mode).
func NewHandler(hub *Hub, validator TokenValidator, allowedOrigin string) *Handler {
	return &Handler{hub: hub, validator: validator, allowedOrigin: allowedOrigin}
}

// ServeHTTP upgrades the HTTP connection to WebSocket and registers the client.
// Auth is done via ?token=<JWT> query parameter.
func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// Extract token from query param.
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}

	// Validate JWT.
	playerID, _, err := h.validator(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	// Accept the WebSocket upgrade.
	opts := &ws.AcceptOptions{}
	if h.allowedOrigin != "" {
		opts.OriginPatterns = []string{h.allowedOrigin}
	} else {
		opts.InsecureSkipVerify = true
	}
	conn, err := ws.Accept(w, r, opts)
	if err != nil {
		slog.Error("websocket accept failed", "error", err)
		return
	}

	client := newClient(h.hub, conn, playerID)
	h.hub.register <- client

	// Send connection_ready.
	client.sendJSON(&Message{
		Type: MsgConnectionReady,
		Data: ConnectionReadyData{
			PlayerID:   playerID,
			ServerTime: time.Now().UTC().Format(time.RFC3339),
		},
	})

	// Start read/write pumps.
	ctx, cancel := context.WithCancel(r.Context())
	defer cancel()

	go client.writePump(ctx)
	client.readPump(ctx) // Blocks until disconnect
}

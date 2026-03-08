package websocket

import (
	"context"
	"encoding/json"
	"log/slog"
	"sync"
	"time"

	ws "github.com/coder/websocket"
)

// Client represents a single WebSocket connection for an authenticated player.
type Client struct {
	hub      *Hub
	conn     *ws.Conn
	playerID int64
	send     chan []byte

	// topics this client is subscribed to (e.g. "village:123")
	mu     sync.RWMutex
	topics map[string]bool
}

const (
	// writeWait is the time allowed to write a message.
	writeWait = 10 * time.Second
	// pongWait is the time allowed to read the next pong.
	pongWait = 60 * time.Second
	// pingPeriod sends pings at this interval (must be < pongWait).
	pingPeriod = 50 * time.Second
	// maxMessageSize limits inbound messages.
	maxMessageSize = 4096
	// sendBufferSize is the outbound channel buffer per client.
	sendBufferSize = 256
)

// newClient creates a new Client.
func newClient(hub *Hub, conn *ws.Conn, playerID int64) *Client {
	return &Client{
		hub:      hub,
		conn:     conn,
		playerID: playerID,
		send:     make(chan []byte, sendBufferSize),
		topics:   make(map[string]bool),
	}
}

// readPump reads messages from the WebSocket connection.
// Runs in its own goroutine per client.
func (c *Client) readPump(ctx context.Context) {
	defer func() {
		c.hub.unregister <- c
		c.conn.Close(ws.StatusNormalClosure, "")
	}()

	c.conn.SetReadLimit(maxMessageSize)

	for {
		_, data, err := c.conn.Read(ctx)
		if err != nil {
			if ws.CloseStatus(err) == ws.StatusNormalClosure || ws.CloseStatus(err) == ws.StatusGoingAway {
				slog.Debug("client disconnected gracefully", "player_id", c.playerID)
			} else {
				slog.Debug("client read error", "player_id", c.playerID, "error", err)
			}
			return
		}

		var msg Message
		if err := json.Unmarshal(data, &msg); err != nil {
			c.sendError("INVALID_JSON", "malformed message")
			continue
		}

		c.handleMessage(ctx, &msg)
	}
}

// writePump writes messages from the send channel to the WebSocket connection.
// Runs in its own goroutine per client.
func (c *Client) writePump(ctx context.Context) {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		c.conn.Close(ws.StatusNormalClosure, "")
	}()

	for {
		select {
		case <-ctx.Done():
			return
		case message, ok := <-c.send:
			if !ok {
				// Hub closed the channel.
				c.conn.Close(ws.StatusNormalClosure, "server shutting down")
				return
			}

			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Write(writeCtx, ws.MessageText, message)
			cancel()
			if err != nil {
				slog.Debug("client write error", "player_id", c.playerID, "error", err)
				return
			}

		case <-ticker.C:
			writeCtx, cancel := context.WithTimeout(ctx, writeWait)
			err := c.conn.Ping(writeCtx)
			cancel()
			if err != nil {
				slog.Debug("client ping failed", "player_id", c.playerID, "error", err)
				return
			}
		}
	}
}

// handleMessage routes an inbound message.
func (c *Client) handleMessage(ctx context.Context, msg *Message) {
	switch msg.Type {
	case MsgPing:
		c.sendJSON(&Message{Type: MsgPong})

	case MsgSubscribe:
		c.handleSubscribe(msg)

	case MsgUnsubscribe:
		c.handleUnsubscribe(msg)

	default:
		c.sendError("UNKNOWN_TYPE", "unknown message type: "+msg.Type)
	}
}

// handleSubscribe subscribes to topics.
func (c *Client) handleSubscribe(msg *Message) {
	raw, _ := json.Marshal(msg.Data)
	var sub SubscribeData
	if err := json.Unmarshal(raw, &sub); err != nil {
		c.sendError("INVALID_DATA", "invalid subscribe payload")
		return
	}

	c.mu.Lock()
	for _, topic := range sub.Topics {
		c.topics[topic] = true
	}
	c.mu.Unlock()

	c.sendJSON(&Message{
		Type: MsgSubscriptionConfirmed,
		Data: SubscribeData{Topics: sub.Topics},
	})
}

// handleUnsubscribe removes topic subscriptions.
func (c *Client) handleUnsubscribe(msg *Message) {
	raw, _ := json.Marshal(msg.Data)
	var sub SubscribeData
	if err := json.Unmarshal(raw, &sub); err != nil {
		c.sendError("INVALID_DATA", "invalid unsubscribe payload")
		return
	}

	c.mu.Lock()
	for _, topic := range sub.Topics {
		delete(c.topics, topic)
	}
	c.mu.Unlock()
}

// sendJSON marshals and queues a message for sending.
func (c *Client) sendJSON(msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to marshal message", "error", err)
		return
	}

	select {
	case c.send <- data:
	default:
		slog.Warn("client send buffer full, dropping message", "player_id", c.playerID)
	}
}

// sendError sends an error message to the client.
func (c *Client) sendError(code, message string) {
	c.sendJSON(&Message{
		Type: MsgError,
		Data: ErrorData{Code: code, Message: message},
	})
}

// IsSubscribed returns true if the client is subscribed to the given topic.
func (c *Client) IsSubscribed(topic string) bool {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.topics[topic]
}

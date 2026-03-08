package websocket

import (
	"encoding/json"
	"fmt"
	"log/slog"
	"sync"
)

// Hub maintains the set of active WebSocket clients and broadcasts messages.
type Hub struct {
	// All connected clients, keyed by playerID.
	mu      sync.RWMutex
	clients map[int64]*Client

	// Channel-based registration/unregistration.
	register   chan *Client
	unregister chan *Client
}

// NewHub creates a new Hub.
func NewHub() *Hub {
	return &Hub{
		clients:    make(map[int64]*Client),
		register:   make(chan *Client),
		unregister: make(chan *Client),
	}
}

// Run starts the hub's event loop. Should be called as a goroutine.
func (h *Hub) Run() {
	for {
		select {
		case client := <-h.register:
			h.mu.Lock()
			// Close existing connection for this player (single-session enforcement).
			if existing, ok := h.clients[client.playerID]; ok {
				close(existing.send)
				delete(h.clients, client.playerID)
				slog.Info("replaced existing connection", "player_id", client.playerID)
			}
			h.clients[client.playerID] = client
			h.mu.Unlock()
			slog.Info("client connected", "player_id", client.playerID, "total", h.ClientCount())

		case client := <-h.unregister:
			h.mu.Lock()
			// Only remove if this is still the current client for this player.
			if current, ok := h.clients[client.playerID]; ok && current == client {
				close(client.send)
				delete(h.clients, client.playerID)
			}
			h.mu.Unlock()
			slog.Info("client disconnected", "player_id", client.playerID, "total", h.ClientCount())
		}
	}
}

// ClientCount returns the number of connected clients.
func (h *Hub) ClientCount() int {
	h.mu.RLock()
	defer h.mu.RUnlock()
	return len(h.clients)
}

// SendToPlayer sends a message to a specific player.
func (h *Hub) SendToPlayer(playerID int64, msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to marshal message", "error", err)
		return
	}

	h.mu.RLock()
	client, ok := h.clients[playerID]
	h.mu.RUnlock()

	if !ok {
		return // Player not connected, message dropped
	}

	select {
	case client.send <- data:
	default:
		slog.Warn("send buffer full", "player_id", playerID)
	}
}

// SendToTopic sends a message to all clients subscribed to a topic.
func (h *Hub) SendToTopic(topic string, msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to marshal message", "error", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		if client.IsSubscribed(topic) {
			select {
			case client.send <- data:
			default:
				slog.Warn("send buffer full for topic delivery", "player_id", client.playerID, "topic", topic)
			}
		}
	}
}

// BroadcastAll sends a message to every connected client.
func (h *Hub) BroadcastAll(msg *Message) {
	data, err := json.Marshal(msg)
	if err != nil {
		slog.Error("failed to marshal broadcast", "error", err)
		return
	}

	h.mu.RLock()
	defer h.mu.RUnlock()

	for _, client := range h.clients {
		select {
		case client.send <- data:
		default:
			slog.Warn("send buffer full for broadcast", "player_id", client.playerID)
		}
	}
}

// BroadcastTrainComplete notifies the owning player that a troop finished training.
func (h *Hub) BroadcastTrainComplete(playerID, villageID int64, troopType string, newTotal int) {
	h.SendToPlayer(playerID, &Message{
		Type: MsgTrainComplete,
		Data: TrainCompleteData{
			VillageID: villageID,
			TroopType: troopType,
			NewTotal:  newTotal,
		},
	})

	topic := fmt.Sprintf("village:%d", villageID)
	h.SendToTopic(topic, &Message{
		Type: MsgTrainComplete,
		Data: TrainCompleteData{
			VillageID: villageID,
			TroopType: troopType,
			NewTotal:  newTotal,
		},
	})
}

// BroadcastBuildComplete notifies the owning player that a building upgrade finished.
func (h *Hub) BroadcastBuildComplete(playerID, villageID int64, buildingType string, newLevel int) {
	h.SendToPlayer(playerID, &Message{
		Type: MsgBuildComplete,
		Data: BuildCompleteData{
			VillageID:    villageID,
			BuildingType: buildingType,
			NewLevel:     newLevel,
		},
	})

	// Also send to anyone subscribed to this village topic.
	topic := fmt.Sprintf("village:%d", villageID)
	h.SendToTopic(topic, &Message{
		Type: MsgBuildComplete,
		Data: BuildCompleteData{
			VillageID:    villageID,
			BuildingType: buildingType,
			NewLevel:     newLevel,
		},
	})
}

package websocket

// Message is the standard WebSocket message envelope.
// All messages (client→server and server→client) follow this format.
type Message struct {
	Type string      `json:"type"`
	Data interface{} `json:"data,omitempty"`
}

// --- Client → Server message types ---
const (
	// MsgSubscribe requests subscription to topics (e.g. "village:123", "map:5,10").
	MsgSubscribe = "subscribe"
	// MsgUnsubscribe removes subscriptions.
	MsgUnsubscribe = "unsubscribe"
	// MsgPing is a client heartbeat.
	MsgPing = "ping"
)

// SubscribeData is the payload for subscribe/unsubscribe messages.
type SubscribeData struct {
	Topics []string `json:"topics"`
}

// --- Server → Client message types ---
const (
	// MsgConnectionReady is sent once after successful auth.
	MsgConnectionReady = "connection_ready"
	// MsgSubscriptionConfirmed acknowledges a subscription.
	MsgSubscriptionConfirmed = "subscription_confirmed"
	// MsgPong is the server heartbeat response.
	MsgPong = "pong"
	// MsgError carries error information.
	MsgError = "error"

	// Game event messages
	MsgBuildComplete  = "build_complete"
	MsgResourceUpdate = "resource_update"
	MsgTrainComplete  = "train_complete"
	MsgAttackIncoming = "attack_incoming"
	MsgCombatResult   = "combat_result"
	MsgWorldEvent     = "world_event"
	MsgAnnouncement   = "announcement"
)

// ConnectionReadyData is sent on successful WebSocket connection.
type ConnectionReadyData struct {
	PlayerID   int64  `json:"player_id"`
	ServerTime string `json:"server_time"`
}

// ErrorData carries error information to the client.
type ErrorData struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

// BuildCompleteData is sent when a building upgrade finishes.
type BuildCompleteData struct {
	VillageID    int64  `json:"village_id"`
	BuildingType string `json:"building_type"`
	NewLevel     int    `json:"new_level"`
}

// TrainCompleteData is sent when a troop finishes training.
type TrainCompleteData struct {
	VillageID int64  `json:"village_id"`
	TroopType string `json:"troop_type"`
	NewTotal  int    `json:"new_total"`
}

// ResourceUpdateData pushes current resource values.
type ResourceUpdateData struct {
	VillageID int64   `json:"village_id"`
	Food      float64 `json:"food"`
	Water     float64 `json:"water"`
	Lumber    float64 `json:"lumber"`
	Stone     float64 `json:"stone"`
}

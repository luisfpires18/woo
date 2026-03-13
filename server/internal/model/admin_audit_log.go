package model

// AdminAuditLog records admin configuration changes.
type AdminAuditLog struct {
	ID            int64  `json:"id"`
	AdminPlayerID int64  `json:"admin_player_id"`
	Action        string `json:"action"` // create, update, delete
	EntityType    string `json:"entity_type"`
	EntityID      *int64 `json:"entity_id,omitempty"`
	OldValueJSON  string `json:"old_value_json,omitempty"`
	NewValueJSON  string `json:"new_value_json,omitempty"`
	CreatedAt     string `json:"created_at"`
}

package model

import "time"

// ResourceBuildingConfig holds per-kingdom cosmetic configuration for a resource building slot.
// There are 96 rows total: 8 kingdoms × 4 resource types × 3 slots.
type ResourceBuildingConfig struct {
	ID           int64     `json:"id"`
	ResourceType string    `json:"resource_type"` // food | water | lumber | stone
	Slot         int       `json:"slot"`          // 1, 2, or 3
	Kingdom      string    `json:"kingdom"`
	DisplayName  string    `json:"display_name"`
	Description  string    `json:"description"`
	DefaultIcon  string    `json:"default_icon"`
	SpritePath   *string   `json:"sprite_path,omitempty"`
	UpdatedAt    time.Time `json:"updated_at"`
}

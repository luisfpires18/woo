package model

import "time"

// BuildingDisplayConfig holds per-kingdom cosmetic configuration for a village building.
// There are 48 rows total: 6 building types × 8 kingdoms.
type BuildingDisplayConfig struct {
	ID           int64     `json:"id"`
	BuildingType string    `json:"building_type"` // town_hall | barracks | stable | archery | workshop | special
	Kingdom      string    `json:"kingdom"`
	DisplayName  string    `json:"display_name"`
	Description  string    `json:"description"`
	DefaultIcon  string    `json:"default_icon"`
	UpdatedAt    time.Time `json:"updated_at"`
}

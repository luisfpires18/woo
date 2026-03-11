package model

import "time"

// TroopDisplayConfig holds per-kingdom cosmetic configuration for a troop type.
type TroopDisplayConfig struct {
	ID               int64     `json:"id"`
	TroopType        string    `json:"troop_type"`
	Kingdom          string    `json:"kingdom"`
	TrainingBuilding string    `json:"training_building"` // barracks | stable | archery | workshop | special
	DisplayName      string    `json:"display_name"`
	Description      string    `json:"description"`
	DefaultIcon      string    `json:"default_icon"`
	UpdatedAt        time.Time `json:"updated_at"`
}

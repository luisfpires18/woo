package model

import "time"

// BuildingQueue represents a building upgrade currently under construction.
type BuildingQueue struct {
	ID           int64     `json:"id"`
	VillageID    int64     `json:"village_id"`
	BuildingType string    `json:"building_type"`
	TargetLevel  int       `json:"target_level"`
	StartedAt    time.Time `json:"started_at"`
	CompletesAt  time.Time `json:"completes_at"`
}

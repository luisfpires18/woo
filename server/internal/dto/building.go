package dto

import "time"

// StartUpgradeRequest is the request body for starting a building upgrade.
type StartUpgradeRequest struct {
	BuildingType string `json:"building_type"`
}

// BuildingQueueResponse represents a queued building upgrade in API responses.
type BuildingQueueResponse struct {
	ID           int64     `json:"id"`
	BuildingType string    `json:"building_type"`
	TargetLevel  int       `json:"target_level"`
	StartedAt    time.Time `json:"started_at"`
	CompletesAt  time.Time `json:"completes_at"`
}

// BuildingCostResponse contains the computed cost for upgrading a building to the next level.
type BuildingCostResponse struct {
	BuildingType string  `json:"building_type"`
	CurrentLevel int     `json:"current_level"`
	TargetLevel  int     `json:"target_level"`
	Iron         float64 `json:"iron"`
	Wood         float64 `json:"wood"`
	Stone        float64 `json:"stone"`
	Food         float64 `json:"food"`
	TimeSec      int     `json:"time_seconds"`
}

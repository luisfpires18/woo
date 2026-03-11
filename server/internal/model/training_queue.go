package model

import "time"

// TrainingQueue represents a single training queue entry.
// One unit completes at CompletesAt; Quantity is decremented each tick until 0.
type TrainingQueue struct {
	ID               int64     `json:"id"`
	VillageID        int64     `json:"village_id"`
	TroopType        string    `json:"troop_type"`
	Quantity         int       `json:"quantity"`
	OriginalQuantity int       `json:"original_quantity"`
	EachDurationSec  int       `json:"each_duration_sec"`
	StartedAt        time.Time `json:"started_at"`
	CompletesAt      time.Time `json:"completes_at"`
}

package dto

import "time"

// --- Training request/response types ---

// StartTrainingRequest is the payload for POST /api/villages/{id}/train.
type StartTrainingRequest struct {
	TroopType string `json:"troop_type"`
	Quantity  int    `json:"quantity"`
}

// TrainingQueueResponse represents a training queue item in API responses.
type TrainingQueueResponse struct {
	ID               int64     `json:"id"`
	TroopType        string    `json:"troop_type"`
	Quantity         int       `json:"quantity"`
	OriginalQuantity int       `json:"original_quantity"`
	EachDurationSec  int       `json:"each_duration_sec"`
	StartedAt        time.Time `json:"started_at"`
	CompletesAt      time.Time `json:"completes_at"`
}

// TrainingCostResponse is returned by the training cost preview endpoint.
type TrainingCostResponse struct {
	TroopType    string  `json:"troop_type"`
	Quantity     int     `json:"quantity"`
	TotalFood    float64 `json:"total_food"`
	TotalWater   float64 `json:"total_water"`
	TotalLumber  float64 `json:"total_lumber"`
	TotalStone   float64 `json:"total_stone"`
	TotalGold    float64 `json:"total_gold"`
	EachTimeSec  int     `json:"each_time_sec"`
	TotalTimeSec int     `json:"total_time_sec"`
}

// TroopInfo represents a troop group in API responses.
type TroopInfo struct {
	Type     string `json:"type"`
	Quantity int    `json:"quantity"`
	Status   string `json:"status"`
}

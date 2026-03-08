package model

// Troop represents a group of units of the same type stationed in a village.
// Uses a stacked model: one row per (village, type) with a quantity.
type Troop struct {
	ID        int64  `json:"id"`
	VillageID int64  `json:"village_id"`
	Type      string `json:"type"`
	Quantity  int    `json:"quantity"`
	Status    string `json:"status"` // stationed, marching, defending, returning
}

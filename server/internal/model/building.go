package model

// Building represents a building within a village.
type Building struct {
	ID           int64  `json:"id"`
	VillageID    int64  `json:"village_id"`
	BuildingType string `json:"building_type"`
	Level        int    `json:"level"`
}

package model

import "time"

// Resources represents the resource snapshot for a village.
type Resources struct {
	VillageID       int64     `json:"village_id"`
	Iron            float64   `json:"iron"`
	Wood            float64   `json:"wood"`
	Stone           float64   `json:"stone"`
	Food            float64   `json:"food"`
	IronRate        float64   `json:"iron_rate"`
	WoodRate        float64   `json:"wood_rate"`
	StoneRate       float64   `json:"stone_rate"`
	FoodRate        float64   `json:"food_rate"`
	FoodConsumption float64   `json:"food_consumption"`
	MaxStorage      float64   `json:"max_storage"`
	LastUpdated     time.Time `json:"last_updated"`
}

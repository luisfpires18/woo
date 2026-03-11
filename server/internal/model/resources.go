package model

import "time"

// Resources represents the resource snapshot for a village.
// The four base resources are: food, water, lumber, stone.
type Resources struct {
	VillageID       int64     `json:"village_id"`
	Food            float64   `json:"food"`
	Water           float64   `json:"water"`
	Lumber          float64   `json:"lumber"`
	Stone           float64   `json:"stone"`
	FoodRate        float64   `json:"food_rate"`
	WaterRate       float64   `json:"water_rate"`
	LumberRate      float64   `json:"lumber_rate"`
	StoneRate       float64   `json:"stone_rate"`
	FoodConsumption float64   `json:"food_consumption"`
	PopUsed         int       `json:"pop_used"`
	MaxFood         float64   `json:"max_food"`
	MaxWater        float64   `json:"max_water"`
	MaxLumber       float64   `json:"max_lumber"`
	MaxStone        float64   `json:"max_stone"`
	LastUpdated     time.Time `json:"last_updated"`
}

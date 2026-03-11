package dto

// VillageResponse is returned for village detail endpoints.
type VillageResponse struct {
	ID            int64                   `json:"id"`
	PlayerID      int64                   `json:"player_id"`
	Name          string                  `json:"name"`
	X             int                     `json:"x"`
	Y             int                     `json:"y"`
	IsCapital     bool                    `json:"is_capital"`
	Gold          float64                 `json:"gold"`
	Buildings     []BuildingInfo          `json:"buildings"`
	Resources     *ResourcesResponse      `json:"resources"`
	BuildQueue    []BuildingQueueResponse `json:"build_queue"`
	Troops        []TroopInfo             `json:"troops"`
	TrainingQueue []TrainingQueueResponse `json:"training_queue"`
}

// BuildingInfo represents a building in API responses.
type BuildingInfo struct {
	ID           int64  `json:"id"`
	BuildingType string `json:"building_type"`
	Level        int    `json:"level"`
}

// ResourcesResponse represents current resources (after lazy calculation).
// Base resources: food, water, lumber, stone.
type ResourcesResponse struct {
	Food            float64 `json:"food"`
	Water           float64 `json:"water"`
	Lumber          float64 `json:"lumber"`
	Stone           float64 `json:"stone"`
	FoodRate        float64 `json:"food_rate"`
	WaterRate       float64 `json:"water_rate"`
	LumberRate      float64 `json:"lumber_rate"`
	StoneRate       float64 `json:"stone_rate"`
	FoodConsumption float64 `json:"food_consumption"`
	MaxFood         float64 `json:"max_food"`
	MaxWater        float64 `json:"max_water"`
	MaxLumber       float64 `json:"max_lumber"`
	MaxStone        float64 `json:"max_stone"`
	PopCap          int     `json:"pop_cap"`
	PopUsed         int     `json:"pop_used"`
}

// VillageListItem is a summary for village list endpoints.
type VillageListItem struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	X         int    `json:"x"`
	Y         int    `json:"y"`
	IsCapital bool   `json:"is_capital"`
}

// RenameVillageRequest is the payload for renaming a village.
type RenameVillageRequest struct {
	Name string `json:"name"`
}

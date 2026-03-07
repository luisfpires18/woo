package dto

// VillageResponse is returned for village detail endpoints.
type VillageResponse struct {
	ID         int64                   `json:"id"`
	PlayerID   int64                   `json:"player_id"`
	Name       string                  `json:"name"`
	X          int                     `json:"x"`
	Y          int                     `json:"y"`
	IsCapital  bool                    `json:"is_capital"`
	Buildings  []BuildingInfo          `json:"buildings"`
	Resources  *ResourcesResponse      `json:"resources"`
	BuildQueue []BuildingQueueResponse `json:"build_queue"`
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
	MaxStorage      float64 `json:"max_storage"`
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

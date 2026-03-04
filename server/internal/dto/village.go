package dto

// VillageResponse is returned for village detail endpoints.
type VillageResponse struct {
	ID        int64              `json:"id"`
	PlayerID  int64              `json:"player_id"`
	Name      string             `json:"name"`
	X         int                `json:"x"`
	Y         int                `json:"y"`
	IsCapital bool               `json:"is_capital"`
	Buildings []BuildingInfo     `json:"buildings"`
	Resources *ResourcesResponse `json:"resources"`
}

// BuildingInfo represents a building in API responses.
type BuildingInfo struct {
	ID           int64  `json:"id"`
	BuildingType string `json:"building_type"`
	Level        int    `json:"level"`
}

// ResourcesResponse represents current resources (after lazy calculation).
type ResourcesResponse struct {
	Iron            float64 `json:"iron"`
	Wood            float64 `json:"wood"`
	Stone           float64 `json:"stone"`
	Food            float64 `json:"food"`
	IronRate        float64 `json:"iron_rate"`
	WoodRate        float64 `json:"wood_rate"`
	StoneRate       float64 `json:"stone_rate"`
	FoodRate        float64 `json:"food_rate"`
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

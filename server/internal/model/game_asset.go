package model

import "time"

// GameAsset represents a game entity (building, resource, unit) with its
// display metadata — emoji fallback icon. Sprites are resolved by filesystem convention.
type GameAsset struct {
	ID          string    `json:"id"`
	Category    string    `json:"category"` // "building", "resource", "unit"
	DisplayName string    `json:"display_name"`
	DefaultIcon string    `json:"default_icon"` // emoji fallback
	UpdatedAt   time.Time `json:"updated_at"`
}

// Category constants for game assets.
const (
	AssetCategoryBuilding      = "building"
	AssetCategoryResource      = "resource"
	AssetCategoryUnit          = "unit"
	AssetCategoryKingdomFlag   = "kingdom_flag"
	AssetCategoryVillageMarker = "village_marker"
	AssetCategoryZoneTile      = "zone_tile"
	AssetCategoryTerrainTile   = "terrain_tile"
)

// ValidAssetCategories is the set of accepted category values.
var ValidAssetCategories = map[string]bool{
	AssetCategoryBuilding:      true,
	AssetCategoryResource:      true,
	AssetCategoryUnit:          true,
	AssetCategoryKingdomFlag:   true,
	AssetCategoryVillageMarker: true,
	AssetCategoryZoneTile:      true,
	AssetCategoryTerrainTile:   true,
}

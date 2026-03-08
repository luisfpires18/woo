package model

import "time"

// GameAsset represents a game entity (building, resource, unit) with its
// display metadata — emoji fallback icon plus optional uploaded sprite.
type GameAsset struct {
	ID           string    `json:"id"`
	Category     string    `json:"category"` // "building", "resource", "unit"
	DisplayName  string    `json:"display_name"`
	DefaultIcon  string    `json:"default_icon"` // emoji fallback
	SpritePath   *string   `json:"sprite_path"`  // relative path under uploads/, nullable
	SpriteWidth  int       `json:"sprite_width"`
	SpriteHeight int       `json:"sprite_height"`
	UpdatedAt    time.Time `json:"updated_at"`
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

// Expected sprite dimensions per category.
var AssetSpriteDimensions = map[string][2]int{
	AssetCategoryBuilding:      {96, 96},
	AssetCategoryResource:      {32, 32},
	AssetCategoryUnit:          {64, 64},
	AssetCategoryKingdomFlag:   {256, 256},
	AssetCategoryVillageMarker: {256, 256},
	AssetCategoryZoneTile:      {256, 256},
	AssetCategoryTerrainTile:   {256, 256},
}

// MaxSpriteBytes per category.
var AssetMaxSpriteBytes = map[string]int64{
	AssetCategoryBuilding:      512 * 1024,  // 512 KB
	AssetCategoryResource:      128 * 1024,  // 128 KB
	AssetCategoryUnit:          256 * 1024,  // 256 KB
	AssetCategoryKingdomFlag:   1024 * 1024, // 1 MB
	AssetCategoryVillageMarker: 512 * 1024,  // 512 KB
	AssetCategoryZoneTile:      512 * 1024,  // 512 KB
	AssetCategoryTerrainTile:   512 * 1024,  // 512 KB
}

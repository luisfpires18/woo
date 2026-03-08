package model

import "time"

// MapTemplate represents a saved map layout that can be applied to the live world map.
// Stored as JSON files on disk in server/data/templates/.
type MapTemplate struct {
	Name        string         `json:"name"`
	Description string         `json:"description"`
	MapSize     int            `json:"map_size"` // e.g. 51 means -25 to +25
	Tiles       []TemplateTile `json:"tiles"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}

// TemplateTile represents a single tile entry in a map template.
type TemplateTile struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	TerrainType string `json:"terrain_type"`
	KingdomZone string `json:"kingdom_zone"`
}

// TileZoneUpdate describes a zone change for a single tile (used for batch zone painting).
type TileZoneUpdate struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	KingdomZone string `json:"kingdom_zone"`
}

// ValidKingdomZones is the set of valid kingdom zone values for templates.
var ValidKingdomZones = map[string]bool{
	ZoneMoraphys:   true,
	ZoneVeridor:    true,
	ZoneSylvara:    true,
	ZoneArkazia:    true,
	ZoneDraxys:     true,
	ZoneZandres:    true,
	ZoneLumus:      true,
	ZoneNordalh:    true,
	ZoneDrakanith:  true,
	ZoneDarkReach:  true,
	ZoneWilderness: true,
}

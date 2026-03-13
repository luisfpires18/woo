package model

// MapTile represents a single tile on the world map grid.
type MapTile struct {
	X             int    `json:"x"`
	Y             int    `json:"y"`
	TerrainType   string `json:"terrain_type"`
	KingdomZone   string `json:"kingdom_zone"`
	OwnerPlayerID *int64 `json:"owner_player_id,omitempty"`
	VillageID     *int64 `json:"village_id,omitempty"`
	VillageName   string `json:"village_name,omitempty"`
	OwnerName     string `json:"owner_name,omitempty"`
	CampID        *int64 `json:"camp_id,omitempty"`
}

// Terrain type constants.
const (
	TerrainPlains   = "plains"
	TerrainForest   = "forest"
	TerrainMountain = "mountain"
	TerrainWater    = "water"
	TerrainDesert   = "desert"
	TerrainSwamp    = "swamp"
	TerrainChasm    = "chasm"
	TerrainBridge   = "bridge"
)

// Kingdom zone constants.
const (
	ZoneMoraphys   = "moraphys"
	ZoneVeridor    = "veridor"
	ZoneSylvara    = "sylvara"
	ZoneArkazia    = "arkazia"
	ZoneDraxys     = "draxys"
	ZoneZandres    = "zandres"
	ZoneLumus      = "lumus"
	ZoneNordalh    = "nordalh"
	ZoneDrakanith  = "drakanith"
	ZoneDarkReach  = "dark_reach"
	ZoneWilderness = "wilderness"
)

// Map dimension constants.
const (
	MapHalf = 25            // map ranges from -25 to +25
	MapSize = MapHalf*2 + 1 // 51
)

// TileTerrainUpdate describes a terrain change for a single tile.
type TileTerrainUpdate struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	TerrainType string `json:"terrain_type"`
}

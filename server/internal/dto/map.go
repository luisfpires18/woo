package dto

import "github.com/luisfpires18/woo/internal/model"

// MapChunkRequest represents the query parameters for a map chunk request.
type MapChunkRequest struct {
	X     int `json:"x"`
	Y     int `json:"y"`
	Range int `json:"range"`
}

// MapChunkResponse is the server response containing a grid of map tiles.
type MapChunkResponse struct {
	CenterX int           `json:"center_x"`
	CenterY int           `json:"center_y"`
	Range   int           `json:"range"`
	Tiles   []MapTileInfo `json:"tiles"`
}

// MapTileInfo is the wire-format for a single map tile.
type MapTileInfo struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	Terrain     string `json:"terrain"`
	Zone        string `json:"zone"`
	VillageID   *int64 `json:"village_id,omitempty"`
	VillageName string `json:"village_name,omitempty"`
	OwnerName   string `json:"owner_name,omitempty"`
	CampID      *int64 `json:"camp_id,omitempty"`
}

// MapTileFromModel converts a model.MapTile to a MapTileInfo DTO.
func MapTileFromModel(t *model.MapTile) MapTileInfo {
	return MapTileInfo{
		X:           t.X,
		Y:           t.Y,
		Terrain:     t.TerrainType,
		Zone:        t.KingdomZone,
		VillageID:   t.VillageID,
		VillageName: t.VillageName,
		OwnerName:   t.OwnerName,
		CampID:      t.CampID,
	}
}

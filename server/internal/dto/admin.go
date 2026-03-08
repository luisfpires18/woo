package dto

import "time"

// --- Player management ---

// PlayerListItem is a single player in an admin player listing.
type PlayerListItem struct {
	ID          int64      `json:"id"`
	Username    string     `json:"username"`
	Email       string     `json:"email"`
	Kingdom     string     `json:"kingdom"`
	Role        string     `json:"role"`
	CreatedAt   time.Time  `json:"created_at"`
	LastLoginAt *time.Time `json:"last_login_at,omitempty"`
}

// PlayerListResponse is the paginated response for GET /api/admin/players.
type PlayerListResponse struct {
	Players []*PlayerListItem `json:"players"`
	Total   int64             `json:"total"`
	Offset  int               `json:"offset"`
	Limit   int               `json:"limit"`
}

// UpdateRoleRequest is the payload for PATCH /api/admin/players/{id}/role.
type UpdateRoleRequest struct {
	Role string `json:"role"`
}

// --- World config ---

// WorldConfigResponse is the response for GET /api/admin/config.
type WorldConfigResponse struct {
	Configs []*WorldConfigEntry `json:"configs"`
}

// WorldConfigEntry is a single config key-value pair.
type WorldConfigEntry struct {
	Key         string    `json:"key"`
	Value       string    `json:"value"`
	Description string    `json:"description,omitempty"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// SetConfigRequest is the payload for PUT /api/admin/config/{key}.
type SetConfigRequest struct {
	Value string `json:"value"`
}

// --- Server stats ---

// StatsResponse is the response for GET /api/admin/stats.
type StatsResponse struct {
	TotalPlayers  int64 `json:"total_players"`
	TotalVillages int64 `json:"total_villages"`
}

// --- Announcements ---

// CreateAnnouncementRequest is the payload for POST /api/admin/announcements.
type CreateAnnouncementRequest struct {
	Title     string  `json:"title"`
	Content   string  `json:"content"`
	ExpiresAt *string `json:"expires_at,omitempty"` // ISO 8601 datetime string, optional
}

// AnnouncementResponse is a single announcement in the listing.
type AnnouncementResponse struct {
	ID        int64      `json:"id"`
	Title     string     `json:"title"`
	Content   string     `json:"content"`
	AuthorID  int64      `json:"author_id"`
	CreatedAt time.Time  `json:"created_at"`
	ExpiresAt *time.Time `json:"expires_at,omitempty"`
}

// --- Game assets ---

// GameAssetDTO is a single game asset in the listing.
type GameAssetDTO struct {
	ID           string    `json:"id"`
	Category     string    `json:"category"`
	DisplayName  string    `json:"display_name"`
	DefaultIcon  string    `json:"default_icon"`
	SpriteURL    *string   `json:"sprite_url"`
	SpriteWidth  int       `json:"sprite_width"`
	SpriteHeight int       `json:"sprite_height"`
	UpdatedAt    time.Time `json:"updated_at"`
}

// GameAssetListResponse is the response for GET /api/admin/assets.
type GameAssetListResponse struct {
	Assets []*GameAssetDTO `json:"assets"`
}

// CreateGameAssetRequest is the payload for POST /api/admin/assets.
type CreateGameAssetRequest struct {
	ID          string `json:"id"`
	Category    string `json:"category"`
	DisplayName string `json:"display_name"`
	DefaultIcon string `json:"default_icon,omitempty"`
}

// UpdateTerrainRequest is the payload for PUT /api/admin/map/terrain.
type UpdateTerrainRequest struct {
	Tiles []TileTerrainPaint `json:"tiles"`
}

// TileTerrainPaint describes a single tile terrain change.
type TileTerrainPaint struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	TerrainType string `json:"terrain_type"`
}

// --- Building display configs ---

// BuildingDisplayConfigDTO is a single building display config in API responses.
type BuildingDisplayConfigDTO struct {
	ID           int64   `json:"id"`
	BuildingType string  `json:"building_type"`
	Kingdom      string  `json:"kingdom"`
	DisplayName  string  `json:"display_name"`
	Description  string  `json:"description"`
	DefaultIcon  string  `json:"default_icon"`
	SpriteURL    *string `json:"sprite_url"`
	UpdatedAt    string  `json:"updated_at"`
}

// BuildingDisplayConfigListResponse is the response for GET /api/admin/building-displays.
type BuildingDisplayConfigListResponse struct {
	Configs []*BuildingDisplayConfigDTO `json:"configs"`
}

// UpdateBuildingDisplayConfigRequest is the payload for PUT /api/admin/building-displays/{id}.
type UpdateBuildingDisplayConfigRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	DefaultIcon string `json:"default_icon"`
}

// --- Resource building configs ---

// ResourceBuildingConfigDTO is a single resource building config in API responses.
type ResourceBuildingConfigDTO struct {
	ID           int64   `json:"id"`
	ResourceType string  `json:"resource_type"`
	Slot         int     `json:"slot"`
	Kingdom      string  `json:"kingdom"`
	DisplayName  string  `json:"display_name"`
	Description  string  `json:"description"`
	DefaultIcon  string  `json:"default_icon"`
	SpriteURL    *string `json:"sprite_url"`
	UpdatedAt    string  `json:"updated_at"`
}

// ResourceBuildingConfigListResponse is the response for GET /api/admin/resource-buildings.
type ResourceBuildingConfigListResponse struct {
	Configs []*ResourceBuildingConfigDTO `json:"configs"`
}

// UpdateResourceBuildingConfigRequest is the payload for PUT /api/admin/resource-buildings/{id}.
type UpdateResourceBuildingConfigRequest struct {
	DisplayName string `json:"display_name"`
	Description string `json:"description"`
	DefaultIcon string `json:"default_icon"`
}

// --- Map templates ---

// CreateTemplateRequest is the payload for POST /api/admin/templates.
type CreateTemplateRequest struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	MapSize     int    `json:"map_size"` // odd number, e.g. 51 means -25..+25. 0 = default (51).
}

// ResizeTemplateRequest is the payload for PUT /api/admin/templates/{name}/resize.
type ResizeTemplateRequest struct {
	MapSize int `json:"map_size"` // new odd size, e.g. 21, 51, 101
}

// UpdateTemplateTerrainRequest is the payload for PUT /api/admin/templates/{name}/terrain.
type UpdateTemplateTerrainRequest struct {
	Tiles []TileTerrainPaint `json:"tiles"`
}

// UpdateTemplateZonesRequest is the payload for PUT /api/admin/templates/{name}/zones.
type UpdateTemplateZonesRequest struct {
	Tiles []TileZonePaint `json:"tiles"`
}

// TileZonePaint describes a single tile zone change.
type TileZonePaint struct {
	X           int    `json:"x"`
	Y           int    `json:"y"`
	KingdomZone string `json:"kingdom_zone"`
}

// ApplyTemplateRequest is the payload for POST /api/admin/templates/{name}/apply.
type ApplyTemplateRequest struct {
	Confirm bool `json:"confirm"`
}

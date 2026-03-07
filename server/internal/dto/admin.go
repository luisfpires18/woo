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

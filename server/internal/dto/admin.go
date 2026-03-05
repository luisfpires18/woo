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

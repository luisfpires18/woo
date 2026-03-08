package dto

import "time"

// ── Season responses ─────────────────────────────────────────────────────────

// SeasonResponse is returned when listing seasons.
type SeasonResponse struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	Status               string     `json:"status"`
	StartDate            *string    `json:"start_date,omitempty"`
	StartedAt            *time.Time `json:"started_at,omitempty"`
	EndedAt              *time.Time `json:"ended_at,omitempty"`
	PlayerCount          int        `json:"player_count"`
	MapTemplateName      string     `json:"map_template_name"`
	GameSpeed            float64    `json:"game_speed"`
	ResourceMultiplier   float64    `json:"resource_multiplier"`
	MaxVillagesPerPlayer int        `json:"max_villages_per_player"`
	WeaponsOfChaosCount  int        `json:"weapons_of_chaos_count"`
	MapWidth             int        `json:"map_width"`
	MapHeight            int        `json:"map_height"`
	CreatedAt            time.Time  `json:"created_at"`
	UpdatedAt            time.Time  `json:"updated_at"`
}

// SeasonDetailResponse is returned for a single season, includes whether
// the requesting player has joined.
type SeasonDetailResponse struct {
	SeasonResponse
	Joined  bool   `json:"joined"`
	Kingdom string `json:"kingdom,omitempty"` // player's kingdom in this season, empty if not joined
}

// SeasonPlayerResponse is a player's participation info.
type SeasonPlayerResponse struct {
	PlayerID int64     `json:"player_id"`
	Username string    `json:"username"`
	Kingdom  string    `json:"kingdom"`
	JoinedAt time.Time `json:"joined_at"`
}

// ── Season requests ──────────────────────────────────────────────────────────

// CreateSeasonRequest is the payload for POST /api/admin/seasons.
type CreateSeasonRequest struct {
	Name                 string  `json:"name"`
	Description          string  `json:"description"`
	StartDate            *string `json:"start_date,omitempty"`
	MapTemplateName      string  `json:"map_template_name"`
	GameSpeed            float64 `json:"game_speed"`
	ResourceMultiplier   float64 `json:"resource_multiplier"`
	MaxVillagesPerPlayer int     `json:"max_villages_per_player"`
	WeaponsOfChaosCount  int     `json:"weapons_of_chaos_count"`
	MapWidth             int     `json:"map_width"`
	MapHeight            int     `json:"map_height"`
}

// UpdateSeasonRequest is the payload for PUT /api/admin/seasons/{id}.
type UpdateSeasonRequest struct {
	Name                 *string  `json:"name,omitempty"`
	Description          *string  `json:"description,omitempty"`
	StartDate            *string  `json:"start_date,omitempty"`
	MapTemplateName      *string  `json:"map_template_name,omitempty"`
	GameSpeed            *float64 `json:"game_speed,omitempty"`
	ResourceMultiplier   *float64 `json:"resource_multiplier,omitempty"`
	MaxVillagesPerPlayer *int     `json:"max_villages_per_player,omitempty"`
	WeaponsOfChaosCount  *int     `json:"weapons_of_chaos_count,omitempty"`
	MapWidth             *int     `json:"map_width,omitempty"`
	MapHeight            *int     `json:"map_height,omitempty"`
}

// JoinSeasonRequest is the payload for POST /api/seasons/{id}/join.
type JoinSeasonRequest struct {
	Kingdom string `json:"kingdom"`
}

// ── Player profile ───────────────────────────────────────────────────────────

// PlayerProfileResponse is returned for GET /api/player/profile.
type PlayerProfileResponse struct {
	ID            int64                `json:"id"`
	Username      string               `json:"username"`
	Email         string               `json:"email"`
	Role          string               `json:"role"`
	CreatedAt     time.Time            `json:"created_at"`
	TotalSeasons  int                  `json:"total_seasons"`
	SeasonHistory []SeasonHistoryEntry `json:"season_history"`
}

// SeasonHistoryEntry is a single season in the player's history.
type SeasonHistoryEntry struct {
	SeasonID     int64     `json:"season_id"`
	SeasonName   string    `json:"season_name"`
	SeasonStatus string    `json:"season_status"`
	Kingdom      string    `json:"kingdom"`
	JoinedAt     time.Time `json:"joined_at"`
	VillageCount int       `json:"village_count"`
}

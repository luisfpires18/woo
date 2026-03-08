package model

import "time"

// Season status constants.
const (
	SeasonStatusUpcoming = "upcoming"
	SeasonStatusActive   = "active"
	SeasonStatusEnded    = "ended"
	SeasonStatusArchived = "archived"
)

// Season represents a game world/server with a timed lifecycle.
type Season struct {
	ID                   int64      `json:"id"`
	Name                 string     `json:"name"`
	Description          string     `json:"description"`
	Status               string     `json:"status"`
	StartDate            *string    `json:"start_date,omitempty"`
	StartedAt            *time.Time `json:"started_at,omitempty"`
	EndedAt              *time.Time `json:"ended_at,omitempty"`
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

// SeasonPlayer represents a player's participation in a specific season.
type SeasonPlayer struct {
	ID       int64     `json:"id"`
	SeasonID int64     `json:"season_id"`
	PlayerID int64     `json:"player_id"`
	Kingdom  string    `json:"kingdom"`
	JoinedAt time.Time `json:"joined_at"`
}

// SeasonHistoryRow is a joined result for player profile season history.
type SeasonHistoryRow struct {
	SeasonID     int64
	SeasonName   string
	SeasonStatus string
	Kingdom      string
	JoinedAt     string // raw datetime string from SQLite
	VillageCount int
}

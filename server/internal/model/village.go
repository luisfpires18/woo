package model

import "time"

// Village represents a player's village on the world map.
type Village struct {
	ID        int64     `json:"id"`
	PlayerID  int64     `json:"player_id"`
	Name      string    `json:"name"`
	X         int       `json:"x"`
	Y         int       `json:"y"`
	IsCapital bool      `json:"is_capital"`
	SeasonID  *int64    `json:"season_id,omitempty"`
	CreatedAt time.Time `json:"created_at"`
}

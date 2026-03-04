package model

import "time"

// RefreshToken represents a stored refresh token for session management.
type RefreshToken struct {
	ID        int64     `json:"id"`
	PlayerID  int64     `json:"player_id"`
	TokenHash string    `json:"-"`
	ExpiresAt time.Time `json:"expires_at"`
	CreatedAt time.Time `json:"created_at"`
}

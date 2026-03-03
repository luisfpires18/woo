package model

import "time"

// Player represents a registered game player.
type Player struct {
	ID            int64      `json:"id"`
	Username      string     `json:"username"`
	Email         string     `json:"email"`
	PasswordHash  string     `json:"-"`
	Kingdom       string     `json:"kingdom"`
	OAuthProvider string     `json:"oauth_provider,omitempty"`
	OAuthID       string     `json:"oauth_id,omitempty"`
	CreatedAt     time.Time  `json:"created_at"`
	LastLoginAt   *time.Time `json:"last_login_at,omitempty"`
}

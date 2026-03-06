package dto

// RegisterRequest is the payload for POST /api/auth/register.
type RegisterRequest struct {
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

// ChooseKingdomRequest is the payload for PUT /api/player/kingdom.
type ChooseKingdomRequest struct {
	Kingdom string `json:"kingdom"`
}

// LoginRequest is the payload for POST /api/auth/login.
type LoginRequest struct {
	Login    string `json:"login"`
	Password string `json:"password"`
}

// RefreshRequest is the payload for POST /api/auth/refresh.
type RefreshRequest struct {
	RefreshToken string `json:"refresh_token"`
}

// AuthResponse is returned on successful login/register.
type AuthResponse struct {
	AccessToken  string      `json:"access_token"`
	RefreshToken string      `json:"refresh_token"`
	Player       *PlayerInfo `json:"player"`
}

// PlayerInfo is the public player data returned in auth responses.
type PlayerInfo struct {
	ID       int64  `json:"id"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Kingdom  string `json:"kingdom"`
	Role     string `json:"role"`
}

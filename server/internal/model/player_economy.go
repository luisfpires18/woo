package model

// PlayerEconomy holds per-player economy state (gold, etc.).
// Gold is shared across all of a player's villages.
type PlayerEconomy struct {
	PlayerID int64   `json:"player_id"`
	Gold     float64 `json:"gold"`
}

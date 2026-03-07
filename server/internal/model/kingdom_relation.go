package model

import "time"

// KingdomRelation represents the diplomatic standing between two kingdoms.
type KingdomRelation struct {
	KingdomA  string    `json:"kingdom_a"`
	KingdomB  string    `json:"kingdom_b"`
	Standing  int       `json:"standing"` // -1000 to +1000
	Status    string    `json:"status"`   // allied, friendly, neutral, hostile, war
	UpdatedAt time.Time `json:"updated_at"`
}

// Diplomacy status constants.
const (
	DiplomacyAllied   = "allied"
	DiplomacyFriendly = "friendly"
	DiplomacyNeutral  = "neutral"
	DiplomacyHostile  = "hostile"
	DiplomacyWar      = "war"
)

// StandingToStatus converts a numeric standing to a diplomacy status.
func StandingToStatus(standing int) string {
	switch {
	case standing >= 500:
		return DiplomacyAllied
	case standing >= 200:
		return DiplomacyFriendly
	case standing > -200:
		return DiplomacyNeutral
	case standing > -500:
		return DiplomacyHostile
	default:
		return DiplomacyWar
	}
}

package service

// validKingdoms lists the playable kingdoms.
// NPC-only kingdoms (e.g. zandres, lumus, drakanith) are intentionally excluded.
var validKingdoms = map[string]bool{
	"veridor": true,
	"sylvara": true,
	"arkazia": true,
	"draxys":  true,
	"nordalh": true,
}

// IsValidKingdom checks if the given kingdom string is a valid playable kingdom.
func IsValidKingdom(kingdom string) bool {
	return validKingdoms[kingdom]
}

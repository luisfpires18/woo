package model

// Expedition represents a player's troop movement to attack a camp.
type Expedition struct {
	ID         int64  `json:"id"`
	PlayerID   int64  `json:"player_id"`
	VillageID  int64  `json:"village_id"`
	CampID     int64  `json:"camp_id"`
	TroopsJSON string `json:"troops_json"`
	DepartedAt string `json:"departed_at"`
	ArrivesAt  string `json:"arrives_at"`
	ReturnsAt  string `json:"returns_at,omitempty"`
	Status     string `json:"status"` // marching, battling, returning, completed
	SeasonID   *int64 `json:"season_id,omitempty"`
}

// Expedition status constants.
const (
	ExpeditionMarching  = "marching"
	ExpeditionBattling  = "battling"
	ExpeditionReturning = "returning"
	ExpeditionCompleted = "completed"
)

// ExpeditionTroop represents a single troop group snapshot on an expedition.
type ExpeditionTroop struct {
	TroopType         string  `json:"troop_type"`
	Quantity          int     `json:"quantity"`
	OriginalQuantity  int     `json:"original_quantity"`
	HP                int     `json:"hp"`
	AttackPower       int     `json:"attack_power"`
	AttackInterval    int     `json:"attack_interval"`
	DefensePercent    float64 `json:"defense_percent"`
	CritChancePercent float64 `json:"crit_chance_percent"`
	Speed             int     `json:"speed"`
	Carry             int     `json:"carry"`
}

package model

// Camp represents a runtime camp instance on the world map.
type Camp struct {
	ID             int64  `json:"id"`
	CampTemplateID int64  `json:"camp_template_id"`
	TileX          int    `json:"tile_x"`
	TileY          int    `json:"tile_y"`
	BeastsJSON     string `json:"beasts_json"`
	SpawnedAt      string `json:"spawned_at"`
	Status         string `json:"status"` // active, under_attack, cleared
	SeasonID       *int64 `json:"season_id,omitempty"`
	SpawnRuleID    *int64 `json:"spawn_rule_id,omitempty"`
}

// Camp status constants.
const (
	CampStatusActive      = "active"
	CampStatusUnderAttack = "under_attack"
	CampStatusCleared     = "cleared"
)

// CampBeast represents a single beast instance inside a camp (stored in beasts_json).
type CampBeast struct {
	BeastTemplateID   int64   `json:"beast_template_id"`
	Name              string  `json:"name"`
	SpriteKey         string  `json:"sprite_key"`
	HP                int     `json:"hp"`
	MaxHP             int     `json:"max_hp"`
	AttackPower       int     `json:"attack_power"`
	AttackInterval    int     `json:"attack_interval"`
	DefensePercent    float64 `json:"defense_percent"`
	CritChancePercent float64 `json:"crit_chance_percent"`
}

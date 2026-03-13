package model

// BeastTemplate defines an admin-configurable beast type.
type BeastTemplate struct {
	ID                int64   `json:"id"`
	Name              string  `json:"name"`
	SpriteKey         string  `json:"sprite_key"`
	HP                int     `json:"hp"`
	AttackPower       int     `json:"attack_power"`
	AttackInterval    int     `json:"attack_interval"`
	DefensePercent    float64 `json:"defense_percent"`
	CritChancePercent float64 `json:"crit_chance_percent"`
	Description       string  `json:"description"`
	CreatedAt         string  `json:"created_at"`
	UpdatedAt         string  `json:"updated_at"`
	UpdatedBy         *int64  `json:"updated_by,omitempty"`
}

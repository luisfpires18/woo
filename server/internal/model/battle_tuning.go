package model

// BattleTuning holds global battle configuration (singleton, id=1).
type BattleTuning struct {
	TickDurationMs        int     `json:"tick_duration_ms"`
	CritDamageMultiplier  float64 `json:"crit_damage_multiplier"`
	MaxDefensePercent     float64 `json:"max_defense_percent"`
	MaxCritChancePercent  float64 `json:"max_crit_chance_percent"`
	MinAttackInterval     int     `json:"min_attack_interval"`
	MarchSpeedTilesPerMin float64 `json:"march_speed_tiles_per_min"`
	MaxTicks              int     `json:"max_ticks"`
	UpdatedAt             string  `json:"updated_at"`
	UpdatedBy             *int64  `json:"updated_by,omitempty"`
}

package model

// RewardTable defines a named collection of possible loot drops.
type RewardTable struct {
	ID        int64  `json:"id"`
	Name      string `json:"name"`
	CreatedAt string `json:"created_at"`
	UpdatedAt string `json:"updated_at"`
	UpdatedBy *int64 `json:"updated_by,omitempty"`
}

// RewardTableEntry defines a single reward entry within a reward table.
type RewardTableEntry struct {
	ID            int64   `json:"id"`
	RewardTableID int64   `json:"reward_table_id"`
	RewardType    string  `json:"reward_type"` // food, water, lumber, stone, gold, rune_fragment
	MinAmount     int     `json:"min_amount"`
	MaxAmount     int     `json:"max_amount"`
	DropChancePct float64 `json:"drop_chance_pct"`
	CreatedAt     string  `json:"created_at"`
}

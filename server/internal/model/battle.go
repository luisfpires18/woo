package model

// Battle stores the result of a resolved combat encounter.
type Battle struct {
	ID                   int64  `json:"id"`
	ExpeditionID         int64  `json:"expedition_id"`
	AttackerSnapshotJSON string `json:"attacker_snapshot_json"`
	DefenderSnapshotJSON string `json:"defender_snapshot_json"`
	Result               string `json:"result"` // attacker_won, defender_won, draw
	AttackerLossesJSON   string `json:"attacker_losses_json"`
	DefenderLossesJSON   string `json:"defender_losses_json"`
	RewardsJSON          string `json:"rewards_json"`
	ReplayData           []byte `json:"replay_data,omitempty"`
	Seed                 int64  `json:"seed"`
	ResolvedAt           string `json:"resolved_at"`
	DurationTicks        int    `json:"duration_ticks"`
}

// Battle result constants.
const (
	BattleResultAttackerWon = "attacker_won"
	BattleResultDefenderWon = "defender_won"
	BattleResultDraw        = "draw"
)

// BattleReward represents a single reward item from a battle.
type BattleReward struct {
	RewardType string `json:"reward_type"` // food, water, lumber, stone, gold, rune_fragment
	Amount     int    `json:"amount"`
}

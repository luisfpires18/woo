package dto

import "time"

// ── Camp responses ───────────────────────────────────────────────────────────

// CampBeastResponse represents a beast inside a camp.
type CampBeastResponse struct {
	BeastTemplateID   int64   `json:"beast_template_id"`
	Name              string  `json:"name"`
	SpriteKey         string  `json:"sprite_key"`
	HP                int     `json:"hp"`
	MaxHP             int     `json:"max_hp"`
	AttackPower       int     `json:"attack_power"`
	AttackInterval    int     `json:"attack_interval"`
	DefensePercent    float64 `json:"defense_percent"`
	CritChancePercent float64 `json:"crit_chance_percent"`
	Count             int     `json:"count"`
}

// CampResponse is returned when listing or viewing a camp on the map.
type CampResponse struct {
	ID           int64               `json:"id"`
	TemplateName string              `json:"template_name"`
	Tier         int                 `json:"tier"`
	SpriteKey    string              `json:"sprite_key"`
	TileX        int                 `json:"tile_x"`
	TileY        int                 `json:"tile_y"`
	Status       string              `json:"status"`
	SpawnedAt    time.Time           `json:"spawned_at"`
	Beasts       []CampBeastResponse `json:"beasts"`
}

// ── Expedition requests ──────────────────────────────────────────────────────

// DispatchExpeditionRequest is the body for sending troops to a camp.
type DispatchExpeditionRequest struct {
	CampID int64           `json:"camp_id"`
	Troops []TroopDispatch `json:"troops"`
}

// TroopDispatch specifies a troop type and quantity to send.
type TroopDispatch struct {
	TroopType string `json:"troop_type"`
	Quantity  int    `json:"quantity"`
}

// ── Expedition responses ─────────────────────────────────────────────────────

// ExpeditionTroopResponse represents a troop group on an expedition.
type ExpeditionTroopResponse struct {
	TroopType        string `json:"troop_type"`
	QuantitySent     int    `json:"quantity_sent"`
	QuantitySurvived int    `json:"quantity_survived"`
}

// ExpeditionResponse is returned when listing or viewing an expedition.
type ExpeditionResponse struct {
	ID           int64                     `json:"id"`
	VillageID    int64                     `json:"village_id"`
	CampID       int64                     `json:"camp_id"`
	Troops       []ExpeditionTroopResponse `json:"troops"`
	Status       string                    `json:"status"`
	DispatchedAt time.Time                 `json:"dispatched_at"`
	ArrivesAt    time.Time                 `json:"arrives_at"`
	ReturnAt     *time.Time                `json:"return_at,omitempty"`
	CompletedAt  *time.Time                `json:"completed_at,omitempty"`
	BattleID     *int64                    `json:"battle_id,omitempty"`
}

// ── Battle report responses ──────────────────────────────────────────────────

// BattleLosses shows aggregate troop losses per side.
type BattleLosses struct {
	TotalSent     int `json:"total_sent"`
	TotalLost     int `json:"total_lost"`
	TotalSurvived int `json:"total_survived"`
}

// BattleRewardResponse shows a single reward item.
type BattleRewardResponse struct {
	ResourceType string `json:"resource_type"`
	Amount       int    `json:"amount"`
}

// BattleReportResponse is the summary returned for a completed battle.
type BattleReportResponse struct {
	ID             int64                  `json:"id"`
	ExpeditionID   int64                  `json:"expedition_id"`
	CampID         int64                  `json:"camp_id"`
	Result         string                 `json:"result"`
	AttackerLosses BattleLosses           `json:"attacker_losses"`
	DefenderLosses BattleLosses           `json:"defender_losses"`
	Rewards        []BattleRewardResponse `json:"rewards"`
	FoughtAt       time.Time              `json:"fought_at"`
}

// ── Admin: Beast Template ────────────────────────────────────────────────────

// CreateBeastTemplateRequest is the body for creating a beast template.
type CreateBeastTemplateRequest struct {
	Name              string  `json:"name"`
	SpriteKey         string  `json:"sprite_key"`
	HP                int     `json:"hp"`
	AttackPower       int     `json:"attack_power"`
	AttackInterval    int     `json:"attack_interval"`
	DefensePercent    float64 `json:"defense_percent"`
	CritChancePercent float64 `json:"crit_chance_percent"`
	Description       string  `json:"description"`
}

// UpdateBeastTemplateRequest is the body for updating a beast template.
type UpdateBeastTemplateRequest struct {
	Name              *string  `json:"name,omitempty"`
	SpriteKey         *string  `json:"sprite_key,omitempty"`
	HP                *int     `json:"hp,omitempty"`
	AttackPower       *int     `json:"attack_power,omitempty"`
	AttackInterval    *int     `json:"attack_interval,omitempty"`
	DefensePercent    *float64 `json:"defense_percent,omitempty"`
	CritChancePercent *float64 `json:"crit_chance_percent,omitempty"`
	Description       *string  `json:"description,omitempty"`
}

// BeastTemplateResponse is returned when listing/viewing beast templates.
type BeastTemplateResponse struct {
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
}

// ── Admin: Camp Template ─────────────────────────────────────────────────────

// CampBeastSlotRequest is one beast slot in a camp template.
type CampBeastSlotRequest struct {
	BeastTemplateID int64 `json:"beast_template_id"`
	MinCount        int   `json:"min_count"`
	MaxCount        int   `json:"max_count"`
}

// CreateCampTemplateRequest is the body for creating a camp template.
type CreateCampTemplateRequest struct {
	Name          string                 `json:"name"`
	Tier          int                    `json:"tier"`
	SpriteKey     string                 `json:"sprite_key"`
	Description   string                 `json:"description"`
	RewardTableID *int64                 `json:"reward_table_id,omitempty"`
	BeastSlots    []CampBeastSlotRequest `json:"beast_slots"`
}

// UpdateCampTemplateRequest is the body for updating a camp template.
type UpdateCampTemplateRequest struct {
	Name          *string                 `json:"name,omitempty"`
	Tier          *int                    `json:"tier,omitempty"`
	SpriteKey     *string                 `json:"sprite_key,omitempty"`
	Description   *string                 `json:"description,omitempty"`
	RewardTableID *int64                  `json:"reward_table_id,omitempty"`
	BeastSlots    *[]CampBeastSlotRequest `json:"beast_slots,omitempty"`
}

// CampTemplateResponse is returned when listing/viewing camp templates.
type CampTemplateResponse struct {
	ID            int64                   `json:"id"`
	Name          string                  `json:"name"`
	Tier          int                     `json:"tier"`
	SpriteKey     string                  `json:"sprite_key"`
	Description   string                  `json:"description"`
	RewardTableID *int64                  `json:"reward_table_id,omitempty"`
	BeastSlots    []CampBeastSlotResponse `json:"beast_slots"`
	CreatedAt     string                  `json:"created_at"`
	UpdatedAt     string                  `json:"updated_at"`
}

// CampBeastSlotResponse is a beast slot in a camp template response.
type CampBeastSlotResponse struct {
	ID              int64  `json:"id"`
	BeastTemplateID int64  `json:"beast_template_id"`
	BeastName       string `json:"beast_name"`
	MinCount        int    `json:"min_count"`
	MaxCount        int    `json:"max_count"`
}

// ── Admin: Spawn Rule ────────────────────────────────────────────────────────

// CampTemplatePoolEntry specifies weights for camp template selection.
type CampTemplatePoolEntry struct {
	CampTemplateID int64 `json:"camp_template_id"`
	Weight         int   `json:"weight"`
}

// CreateSpawnRuleRequest is the body for creating a spawn rule.
type CreateSpawnRuleRequest struct {
	Name               string                  `json:"name"`
	TerrainTypes       []string                `json:"terrain_types"`
	ZoneTypes          []string                `json:"zone_types"`
	CampTemplatePool   []CampTemplatePoolEntry `json:"camp_template_pool"`
	MaxCamps           int                     `json:"max_camps"`
	SpawnIntervalSec   int                     `json:"spawn_interval_sec"`
	DespawnAfterSec    int                     `json:"despawn_after_sec"`
	MinCampDistance    int                     `json:"min_camp_distance"`
	MinVillageDistance int                     `json:"min_village_distance"`
	Enabled            bool                    `json:"enabled"`
}

// UpdateSpawnRuleRequest is the body for updating a spawn rule.
type UpdateSpawnRuleRequest struct {
	Name               *string                  `json:"name,omitempty"`
	TerrainTypes       *[]string                `json:"terrain_types,omitempty"`
	ZoneTypes          *[]string                `json:"zone_types,omitempty"`
	CampTemplatePool   *[]CampTemplatePoolEntry `json:"camp_template_pool,omitempty"`
	MaxCamps           *int                     `json:"max_camps,omitempty"`
	SpawnIntervalSec   *int                     `json:"spawn_interval_sec,omitempty"`
	DespawnAfterSec    *int                     `json:"despawn_after_sec,omitempty"`
	MinCampDistance    *int                     `json:"min_camp_distance,omitempty"`
	MinVillageDistance *int                     `json:"min_village_distance,omitempty"`
	Enabled            *bool                    `json:"enabled,omitempty"`
}

// SpawnRuleResponse is returned when listing/viewing spawn rules.
type SpawnRuleResponse struct {
	ID                 int64                   `json:"id"`
	Name               string                  `json:"name"`
	TerrainTypes       []string                `json:"terrain_types"`
	ZoneTypes          []string                `json:"zone_types"`
	CampTemplatePool   []CampTemplatePoolEntry `json:"camp_template_pool"`
	MaxCamps           int                     `json:"max_camps"`
	SpawnIntervalSec   int                     `json:"spawn_interval_sec"`
	DespawnAfterSec    int                     `json:"despawn_after_sec"`
	MinCampDistance    int                     `json:"min_camp_distance"`
	MinVillageDistance int                     `json:"min_village_distance"`
	Enabled            bool                    `json:"enabled"`
	CreatedAt          string                  `json:"created_at"`
	UpdatedAt          string                  `json:"updated_at"`
}

// ── Admin: Reward Table ──────────────────────────────────────────────────────

// RewardEntryRequest is one entry in a reward table.
type RewardEntryRequest struct {
	RewardType string `json:"reward_type"`
	MinAmount  int    `json:"min_amount"`
	MaxAmount  int    `json:"max_amount"`
	DropChance int    `json:"drop_chance"` // 0-100
}

// CreateRewardTableRequest is the body for creating a reward table.
type CreateRewardTableRequest struct {
	Name    string               `json:"name"`
	Entries []RewardEntryRequest `json:"entries"`
}

// UpdateRewardTableRequest is the body for updating a reward table.
type UpdateRewardTableRequest struct {
	Name    *string               `json:"name,omitempty"`
	Entries *[]RewardEntryRequest `json:"entries,omitempty"`
}

// RewardTableResponse is returned when listing/viewing reward tables.
type RewardTableResponse struct {
	ID        int64                 `json:"id"`
	Name      string                `json:"name"`
	Entries   []RewardEntryResponse `json:"entries"`
	CreatedAt string                `json:"created_at"`
	UpdatedAt string                `json:"updated_at"`
}

// RewardEntryResponse is one entry in a reward table response.
type RewardEntryResponse struct {
	ID         int64  `json:"id"`
	RewardType string `json:"reward_type"`
	MinAmount  int    `json:"min_amount"`
	MaxAmount  int    `json:"max_amount"`
	DropChance int    `json:"drop_chance"`
}

// ── Admin: Battle Tuning ─────────────────────────────────────────────────────

// UpdateBattleTuningRequest is the body for updating battle tuning.
type UpdateBattleTuningRequest struct {
	TickDurationMs        *int     `json:"tick_duration_ms,omitempty"`
	CritDamageMultiplier  *float64 `json:"crit_damage_multiplier,omitempty"`
	MaxDefensePercent     *float64 `json:"max_defense_percent,omitempty"`
	MaxCritChancePercent  *float64 `json:"max_crit_chance_percent,omitempty"`
	MinAttackInterval     *int     `json:"min_attack_interval,omitempty"`
	MarchSpeedTilesPerMin *float64 `json:"march_speed_tiles_per_min,omitempty"`
	MaxTicks              *int     `json:"max_ticks,omitempty"`
}

// BattleTuningResponse is returned when viewing battle tuning.
type BattleTuningResponse struct {
	TickDurationMs        int     `json:"tick_duration_ms"`
	CritDamageMultiplier  float64 `json:"crit_damage_multiplier"`
	MaxDefensePercent     float64 `json:"max_defense_percent"`
	MaxCritChancePercent  float64 `json:"max_crit_chance_percent"`
	MinAttackInterval     int     `json:"min_attack_interval"`
	MarchSpeedTilesPerMin float64 `json:"march_speed_tiles_per_min"`
	MaxTicks              int     `json:"max_ticks"`
	UpdatedAt             string  `json:"updated_at"`
}

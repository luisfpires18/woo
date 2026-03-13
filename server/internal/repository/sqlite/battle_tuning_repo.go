package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type battleTuningRepo struct {
	db *sql.DB
}

// NewBattleTuningRepo creates a new SQLite-backed BattleTuningRepository.
func NewBattleTuningRepo(db *sql.DB) *battleTuningRepo {
	return &battleTuningRepo{db: db}
}

func (r *battleTuningRepo) Get(ctx context.Context) (*model.BattleTuning, error) {
	var bt model.BattleTuning
	err := r.db.QueryRowContext(ctx,
		`SELECT tick_duration_ms, crit_damage_multiplier, max_defense_percent, max_crit_chance_percent,
		 min_attack_interval, march_speed_tiles_per_min, max_ticks, updated_at, updated_by
		 FROM battle_tuning WHERE id = 1`,
	).Scan(&bt.TickDurationMs, &bt.CritDamageMultiplier, &bt.MaxDefensePercent, &bt.MaxCritChancePercent,
		&bt.MinAttackInterval, &bt.MarchSpeedTilesPerMin, &bt.MaxTicks, &bt.UpdatedAt, &bt.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get battle tuning: %w", err)
	}
	return &bt, nil
}

func (r *battleTuningRepo) Update(ctx context.Context, bt *model.BattleTuning) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE battle_tuning SET tick_duration_ms = ?, crit_damage_multiplier = ?, max_defense_percent = ?,
		 max_crit_chance_percent = ?, min_attack_interval = ?, march_speed_tiles_per_min = ?, max_ticks = ?,
		 updated_at = datetime('now'), updated_by = ?
		 WHERE id = 1`,
		bt.TickDurationMs, bt.CritDamageMultiplier, bt.MaxDefensePercent, bt.MaxCritChancePercent,
		bt.MinAttackInterval, bt.MarchSpeedTilesPerMin, bt.MaxTicks, bt.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("update battle tuning: %w", err)
	}
	return nil
}

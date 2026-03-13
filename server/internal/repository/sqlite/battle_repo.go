package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type battleRepo struct {
	db *sql.DB
}

// NewBattleRepo creates a new SQLite-backed BattleRepository.
func NewBattleRepo(db *sql.DB) *battleRepo {
	return &battleRepo{db: db}
}

func (r *battleRepo) Create(ctx context.Context, battle *model.Battle) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO battles (expedition_id, attacker_snapshot_json, defender_snapshot_json, result,
		 attacker_losses_json, defender_losses_json, rewards_json, replay_data, seed, duration_ticks)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		battle.ExpeditionID, battle.AttackerSnapshotJSON, battle.DefenderSnapshotJSON, battle.Result,
		battle.AttackerLossesJSON, battle.DefenderLossesJSON, battle.RewardsJSON,
		battle.ReplayData, battle.Seed, battle.DurationTicks,
	)
	if err != nil {
		return fmt.Errorf("create battle: %w", err)
	}
	id, _ := result.LastInsertId()
	battle.ID = id
	return nil
}

func (r *battleRepo) GetByID(ctx context.Context, id int64) (*model.Battle, error) {
	var battle model.Battle
	err := r.db.QueryRowContext(ctx,
		`SELECT id, expedition_id, attacker_snapshot_json, defender_snapshot_json, result,
		 attacker_losses_json, defender_losses_json, rewards_json, seed, resolved_at, duration_ticks
		 FROM battles WHERE id = ?`, id,
	).Scan(&battle.ID, &battle.ExpeditionID, &battle.AttackerSnapshotJSON, &battle.DefenderSnapshotJSON,
		&battle.Result, &battle.AttackerLossesJSON, &battle.DefenderLossesJSON, &battle.RewardsJSON,
		&battle.Seed, &battle.ResolvedAt, &battle.DurationTicks)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get battle %d: %w", id, err)
	}
	return &battle, nil
}

func (r *battleRepo) GetByExpeditionID(ctx context.Context, expeditionID int64) (*model.Battle, error) {
	var battle model.Battle
	err := r.db.QueryRowContext(ctx,
		`SELECT id, expedition_id, attacker_snapshot_json, defender_snapshot_json, result,
		 attacker_losses_json, defender_losses_json, rewards_json, seed, resolved_at, duration_ticks
		 FROM battles WHERE expedition_id = ?`, expeditionID,
	).Scan(&battle.ID, &battle.ExpeditionID, &battle.AttackerSnapshotJSON, &battle.DefenderSnapshotJSON,
		&battle.Result, &battle.AttackerLossesJSON, &battle.DefenderLossesJSON, &battle.RewardsJSON,
		&battle.Seed, &battle.ResolvedAt, &battle.DurationTicks)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get battle by expedition %d: %w", expeditionID, err)
	}
	return &battle, nil
}

func (r *battleRepo) GetReplayData(ctx context.Context, id int64) ([]byte, error) {
	var data []byte
	err := r.db.QueryRowContext(ctx,
		`SELECT replay_data FROM battles WHERE id = ?`, id,
	).Scan(&data)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get replay data for battle %d: %w", id, err)
	}
	return data, nil
}

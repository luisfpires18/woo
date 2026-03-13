package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/luisfpires18/woo/internal/model"
)

type campRepo struct {
	db *sql.DB
}

// NewCampRepo creates a new SQLite-backed CampRepository.
func NewCampRepo(db *sql.DB) *campRepo {
	return &campRepo{db: db}
}

func (r *campRepo) Create(ctx context.Context, camp *model.Camp) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO camps (camp_template_id, tile_x, tile_y, beasts_json, status, season_id, spawn_rule_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		camp.CampTemplateID, camp.TileX, camp.TileY, camp.BeastsJSON, camp.Status, camp.SeasonID, camp.SpawnRuleID,
	)
	if err != nil {
		return fmt.Errorf("create camp: %w", err)
	}
	id, _ := result.LastInsertId()
	camp.ID = id

	// Update world_map tile to reference this camp
	_, err = r.db.ExecContext(ctx,
		`UPDATE world_map SET camp_id = ? WHERE x = ? AND y = ?`,
		camp.ID, camp.TileX, camp.TileY,
	)
	if err != nil {
		return fmt.Errorf("link camp %d to tile (%d,%d): %w", camp.ID, camp.TileX, camp.TileY, err)
	}
	return nil
}

func (r *campRepo) GetByID(ctx context.Context, id int64) (*model.Camp, error) {
	var camp model.Camp
	err := r.db.QueryRowContext(ctx,
		`SELECT id, camp_template_id, tile_x, tile_y, beasts_json, spawned_at, status, season_id, spawn_rule_id
		 FROM camps WHERE id = ?`, id,
	).Scan(&camp.ID, &camp.CampTemplateID, &camp.TileX, &camp.TileY, &camp.BeastsJSON,
		&camp.SpawnedAt, &camp.Status, &camp.SeasonID, &camp.SpawnRuleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get camp %d: %w", id, err)
	}
	return &camp, nil
}

func (r *campRepo) GetByTile(ctx context.Context, x, y int) (*model.Camp, error) {
	var camp model.Camp
	err := r.db.QueryRowContext(ctx,
		`SELECT id, camp_template_id, tile_x, tile_y, beasts_json, spawned_at, status, season_id, spawn_rule_id
		 FROM camps WHERE tile_x = ? AND tile_y = ? AND status != 'cleared'`, x, y,
	).Scan(&camp.ID, &camp.CampTemplateID, &camp.TileX, &camp.TileY, &camp.BeastsJSON,
		&camp.SpawnedAt, &camp.Status, &camp.SeasonID, &camp.SpawnRuleID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get camp at (%d,%d): %w", x, y, err)
	}
	return &camp, nil
}

func (r *campRepo) GetByStatus(ctx context.Context, status string) ([]*model.Camp, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, camp_template_id, tile_x, tile_y, beasts_json, spawned_at, status, season_id, spawn_rule_id
		 FROM camps WHERE status = ? ORDER BY id ASC`, status,
	)
	if err != nil {
		return nil, fmt.Errorf("list camps by status %s: %w", status, err)
	}
	defer rows.Close()
	return scanCamps(rows)
}

func (r *campRepo) ListActive(ctx context.Context) ([]*model.Camp, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, camp_template_id, tile_x, tile_y, beasts_json, spawned_at, status, season_id, spawn_rule_id
		 FROM camps WHERE status IN ('active', 'under_attack') ORDER BY id ASC`,
	)
	if err != nil {
		return nil, fmt.Errorf("list active camps: %w", err)
	}
	defer rows.Close()
	return scanCamps(rows)
}

func (r *campRepo) CountBySpawnRule(ctx context.Context, spawnRuleID int64) (int, error) {
	var count int
	err := r.db.QueryRowContext(ctx,
		`SELECT COUNT(*) FROM camps WHERE spawn_rule_id = ? AND status IN ('active', 'under_attack')`,
		spawnRuleID,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count camps by spawn rule %d: %w", spawnRuleID, err)
	}
	return count, nil
}

func (r *campRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE camps SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update camp %d status: %w", id, err)
	}

	// If cleared, unlink from world_map
	if status == model.CampStatusCleared {
		_, err = r.db.ExecContext(ctx,
			`UPDATE world_map SET camp_id = NULL WHERE camp_id = ?`, id)
		if err != nil {
			return fmt.Errorf("unlink cleared camp %d from world_map: %w", id, err)
		}
	}
	return nil
}

func (r *campRepo) Delete(ctx context.Context, id int64) error {
	// Unlink from world_map first
	_, _ = r.db.ExecContext(ctx, `UPDATE world_map SET camp_id = NULL WHERE camp_id = ?`, id)

	_, err := r.db.ExecContext(ctx, `DELETE FROM camps WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete camp %d: %w", id, err)
	}
	return nil
}

func (r *campRepo) GetExpiredCamps(ctx context.Context, now time.Time) ([]*model.Camp, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT c.id, c.camp_template_id, c.tile_x, c.tile_y, c.beasts_json, c.spawned_at, c.status, c.season_id, c.spawn_rule_id
		 FROM camps c
		 JOIN spawn_rules sr ON c.spawn_rule_id = sr.id
		 WHERE c.status = 'active'
		   AND sr.despawn_after_sec > 0
		   AND datetime(c.spawned_at, '+' || sr.despawn_after_sec || ' seconds') <= ?`,
		now.UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, fmt.Errorf("get expired camps: %w", err)
	}
	defer rows.Close()
	return scanCamps(rows)
}

func scanCamps(rows *sql.Rows) ([]*model.Camp, error) {
	var camps []*model.Camp
	for rows.Next() {
		var camp model.Camp
		if err := rows.Scan(&camp.ID, &camp.CampTemplateID, &camp.TileX, &camp.TileY, &camp.BeastsJSON,
			&camp.SpawnedAt, &camp.Status, &camp.SeasonID, &camp.SpawnRuleID); err != nil {
			return nil, fmt.Errorf("scan camp: %w", err)
		}
		camps = append(camps, &camp)
	}
	return camps, rows.Err()
}

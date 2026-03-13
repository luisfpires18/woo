package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type spawnRuleRepo struct {
	db *sql.DB
}

// NewSpawnRuleRepo creates a new SQLite-backed SpawnRuleRepository.
func NewSpawnRuleRepo(db *sql.DB) *spawnRuleRepo {
	return &spawnRuleRepo{db: db}
}

func (r *spawnRuleRepo) Create(ctx context.Context, rule *model.SpawnRule) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO spawn_rules (name, terrain_types_json, zone_types_json, camp_template_pool_json,
		 max_camps, spawn_interval_sec, despawn_after_sec, min_camp_distance, min_village_distance, enabled, updated_by)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		rule.Name, rule.TerrainTypesJSON, rule.ZoneTypesJSON, rule.CampTemplatePoolJSON,
		rule.MaxCamps, rule.SpawnIntervalSec, rule.DespawnAfterSec,
		rule.MinCampDistance, rule.MinVillageDistance, rule.Enabled, rule.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("create spawn rule: %w", err)
	}
	id, _ := result.LastInsertId()
	rule.ID = id
	return nil
}

func (r *spawnRuleRepo) GetByID(ctx context.Context, id int64) (*model.SpawnRule, error) {
	var rule model.SpawnRule
	var enabled int
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, terrain_types_json, zone_types_json, camp_template_pool_json,
		 max_camps, spawn_interval_sec, despawn_after_sec, min_camp_distance, min_village_distance,
		 enabled, created_at, updated_at, updated_by
		 FROM spawn_rules WHERE id = ?`, id,
	).Scan(&rule.ID, &rule.Name, &rule.TerrainTypesJSON, &rule.ZoneTypesJSON, &rule.CampTemplatePoolJSON,
		&rule.MaxCamps, &rule.SpawnIntervalSec, &rule.DespawnAfterSec,
		&rule.MinCampDistance, &rule.MinVillageDistance,
		&enabled, &rule.CreatedAt, &rule.UpdatedAt, &rule.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get spawn rule %d: %w", id, err)
	}
	rule.Enabled = enabled == 1
	return &rule, nil
}

func (r *spawnRuleRepo) GetAll(ctx context.Context) ([]*model.SpawnRule, error) {
	return r.queryRules(ctx, `SELECT id, name, terrain_types_json, zone_types_json, camp_template_pool_json,
		 max_camps, spawn_interval_sec, despawn_after_sec, min_camp_distance, min_village_distance,
		 enabled, created_at, updated_at, updated_by
		 FROM spawn_rules ORDER BY id ASC`)
}

func (r *spawnRuleRepo) GetEnabled(ctx context.Context) ([]*model.SpawnRule, error) {
	return r.queryRules(ctx, `SELECT id, name, terrain_types_json, zone_types_json, camp_template_pool_json,
		 max_camps, spawn_interval_sec, despawn_after_sec, min_camp_distance, min_village_distance,
		 enabled, created_at, updated_at, updated_by
		 FROM spawn_rules WHERE enabled = 1 ORDER BY id ASC`)
}

func (r *spawnRuleRepo) queryRules(ctx context.Context, query string) ([]*model.SpawnRule, error) {
	rows, err := r.db.QueryContext(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("query spawn rules: %w", err)
	}
	defer rows.Close()

	var rules []*model.SpawnRule
	for rows.Next() {
		var rule model.SpawnRule
		var enabled int
		if err := rows.Scan(&rule.ID, &rule.Name, &rule.TerrainTypesJSON, &rule.ZoneTypesJSON, &rule.CampTemplatePoolJSON,
			&rule.MaxCamps, &rule.SpawnIntervalSec, &rule.DespawnAfterSec,
			&rule.MinCampDistance, &rule.MinVillageDistance,
			&enabled, &rule.CreatedAt, &rule.UpdatedAt, &rule.UpdatedBy); err != nil {
			return nil, fmt.Errorf("scan spawn rule: %w", err)
		}
		rule.Enabled = enabled == 1
		rules = append(rules, &rule)
	}
	return rules, rows.Err()
}

func (r *spawnRuleRepo) Update(ctx context.Context, rule *model.SpawnRule) error {
	enabledInt := 0
	if rule.Enabled {
		enabledInt = 1
	}
	_, err := r.db.ExecContext(ctx,
		`UPDATE spawn_rules SET name = ?, terrain_types_json = ?, zone_types_json = ?, camp_template_pool_json = ?,
		 max_camps = ?, spawn_interval_sec = ?, despawn_after_sec = ?, min_camp_distance = ?, min_village_distance = ?,
		 enabled = ?, updated_at = datetime('now'), updated_by = ?
		 WHERE id = ?`,
		rule.Name, rule.TerrainTypesJSON, rule.ZoneTypesJSON, rule.CampTemplatePoolJSON,
		rule.MaxCamps, rule.SpawnIntervalSec, rule.DespawnAfterSec,
		rule.MinCampDistance, rule.MinVillageDistance,
		enabledInt, rule.UpdatedBy, rule.ID,
	)
	if err != nil {
		return fmt.Errorf("update spawn rule %d: %w", rule.ID, err)
	}
	return nil
}

func (r *spawnRuleRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM spawn_rules WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete spawn rule %d: %w", id, err)
	}
	return nil
}

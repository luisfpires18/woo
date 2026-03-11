package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

// TroopDisplayConfigRepo implements repository.TroopDisplayConfigRepository using SQLite.
type TroopDisplayConfigRepo struct {
	db *sql.DB
}

// NewTroopDisplayConfigRepo creates a new TroopDisplayConfigRepo.
func NewTroopDisplayConfigRepo(db *sql.DB) *TroopDisplayConfigRepo {
	return &TroopDisplayConfigRepo{db: db}
}

func scanTroopDisplayConfig(row interface{ Scan(dest ...any) error }) (*model.TroopDisplayConfig, error) {
	var c model.TroopDisplayConfig
	var updatedAtStr string
	err := row.Scan(
		&c.ID, &c.TroopType, &c.Kingdom, &c.TrainingBuilding,
		&c.DisplayName, &c.Description, &c.DefaultIcon,
		&updatedAtStr,
	)
	if err != nil {
		return nil, err
	}
	c.UpdatedAt, _ = parseTime(updatedAtStr)
	return &c, nil
}

// GetAll returns every troop display config ordered by training_building, troop_type.
func (r *TroopDisplayConfigRepo) GetAll(ctx context.Context) ([]*model.TroopDisplayConfig, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, troop_type, kingdom, training_building, display_name, description, default_icon, updated_at
		 FROM troop_display_configs ORDER BY training_building, troop_type`)
	if err != nil {
		return nil, fmt.Errorf("query troop_display_configs: %w", err)
	}
	defer rows.Close()

	var configs []*model.TroopDisplayConfig
	for rows.Next() {
		c, err := scanTroopDisplayConfig(rows)
		if err != nil {
			return nil, fmt.Errorf("scan troop_display_config: %w", err)
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}

// GetByKingdom returns all troop display configs for a specific kingdom.
func (r *TroopDisplayConfigRepo) GetByKingdom(ctx context.Context, kingdom string) ([]*model.TroopDisplayConfig, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, troop_type, kingdom, training_building, display_name, description, default_icon, updated_at
		 FROM troop_display_configs WHERE kingdom = ? ORDER BY training_building, troop_type`, kingdom)
	if err != nil {
		return nil, fmt.Errorf("query troop_display_configs by kingdom: %w", err)
	}
	defer rows.Close()

	var configs []*model.TroopDisplayConfig
	for rows.Next() {
		c, err := scanTroopDisplayConfig(rows)
		if err != nil {
			return nil, fmt.Errorf("scan troop_display_config: %w", err)
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}

// GetByID returns a single troop display config by its primary key.
func (r *TroopDisplayConfigRepo) GetByID(ctx context.Context, id int64) (*model.TroopDisplayConfig, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, troop_type, kingdom, training_building, display_name, description, default_icon, updated_at
		 FROM troop_display_configs WHERE id = ?`, id)
	c, err := scanTroopDisplayConfig(row)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return c, err
}

// Update persists changes to display_name, description, and default_icon for a config row.
func (r *TroopDisplayConfigRepo) Update(ctx context.Context, cfg *model.TroopDisplayConfig) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE troop_display_configs SET display_name = ?, description = ?, default_icon = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		cfg.DisplayName, cfg.Description, cfg.DefaultIcon, cfg.ID)
	if err != nil {
		return fmt.Errorf("update troop_display_config %d: %w", cfg.ID, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

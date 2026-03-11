package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

// ResourceBuildingConfigRepo implements repository.ResourceBuildingConfigRepository using SQLite.
type ResourceBuildingConfigRepo struct {
	db *sql.DB
}

// NewResourceBuildingConfigRepo creates a new ResourceBuildingConfigRepo.
func NewResourceBuildingConfigRepo(db *sql.DB) *ResourceBuildingConfigRepo {
	return &ResourceBuildingConfigRepo{db: db}
}

func scanResourceBuildingConfig(row interface{ Scan(dest ...any) error }) (*model.ResourceBuildingConfig, error) {
	var c model.ResourceBuildingConfig
	var updatedAtStr string
	err := row.Scan(
		&c.ID, &c.ResourceType, &c.Slot, &c.Kingdom,
		&c.DisplayName, &c.Description, &c.DefaultIcon,
		&updatedAtStr,
	)
	if err != nil {
		return nil, err
	}
	c.UpdatedAt, _ = parseTime(updatedAtStr)
	return &c, nil
}

// GetAll returns every resource building config ordered by resource_type, slot, kingdom.
func (r *ResourceBuildingConfigRepo) GetAll(ctx context.Context) ([]*model.ResourceBuildingConfig, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, resource_type, slot, kingdom, display_name, description, default_icon, updated_at
		 FROM resource_building_configs ORDER BY resource_type, slot, kingdom`)
	if err != nil {
		return nil, fmt.Errorf("query resource_building_configs: %w", err)
	}
	defer rows.Close()

	var configs []*model.ResourceBuildingConfig
	for rows.Next() {
		c, err := scanResourceBuildingConfig(rows)
		if err != nil {
			return nil, fmt.Errorf("scan resource_building_config: %w", err)
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}

// GetByKingdom returns all resource building configs for a specific kingdom.
func (r *ResourceBuildingConfigRepo) GetByKingdom(ctx context.Context, kingdom string) ([]*model.ResourceBuildingConfig, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, resource_type, slot, kingdom, display_name, description, default_icon, updated_at
		 FROM resource_building_configs WHERE kingdom = ? ORDER BY resource_type, slot`, kingdom)
	if err != nil {
		return nil, fmt.Errorf("query resource_building_configs by kingdom: %w", err)
	}
	defer rows.Close()

	var configs []*model.ResourceBuildingConfig
	for rows.Next() {
		c, err := scanResourceBuildingConfig(rows)
		if err != nil {
			return nil, fmt.Errorf("scan resource_building_config: %w", err)
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}

// GetByID returns a single resource building config by its primary key.
func (r *ResourceBuildingConfigRepo) GetByID(ctx context.Context, id int64) (*model.ResourceBuildingConfig, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, resource_type, slot, kingdom, display_name, description, default_icon, updated_at
		 FROM resource_building_configs WHERE id = ?`, id)
	c, err := scanResourceBuildingConfig(row)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return c, err
}

// Update persists changes to display_name, description, and default_icon for a config row.
func (r *ResourceBuildingConfigRepo) Update(ctx context.Context, cfg *model.ResourceBuildingConfig) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE resource_building_configs SET display_name = ?, description = ?, default_icon = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		cfg.DisplayName, cfg.Description, cfg.DefaultIcon, cfg.ID)
	if err != nil {
		return fmt.Errorf("update resource_building_config %d: %w", cfg.ID, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

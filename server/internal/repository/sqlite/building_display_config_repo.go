package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

// BuildingDisplayConfigRepo implements repository.BuildingDisplayConfigRepository using SQLite.
type BuildingDisplayConfigRepo struct {
	db *sql.DB
}

// NewBuildingDisplayConfigRepo creates a new BuildingDisplayConfigRepo.
func NewBuildingDisplayConfigRepo(db *sql.DB) *BuildingDisplayConfigRepo {
	return &BuildingDisplayConfigRepo{db: db}
}

func scanBuildingDisplayConfig(row interface{ Scan(dest ...any) error }) (*model.BuildingDisplayConfig, error) {
	var c model.BuildingDisplayConfig
	var spritePath sql.NullString
	var updatedAtStr string
	err := row.Scan(
		&c.ID, &c.BuildingType, &c.Kingdom,
		&c.DisplayName, &c.Description, &c.DefaultIcon,
		&spritePath, &updatedAtStr,
	)
	if err != nil {
		return nil, err
	}
	if spritePath.Valid {
		c.SpritePath = &spritePath.String
	}
	c.UpdatedAt, _ = parseTime(updatedAtStr)
	return &c, nil
}

// GetAll returns every building display config ordered by building_type, kingdom.
func (r *BuildingDisplayConfigRepo) GetAll(ctx context.Context) ([]*model.BuildingDisplayConfig, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, building_type, kingdom, display_name, description, default_icon, sprite_path, updated_at
		 FROM building_display_configs ORDER BY building_type, kingdom`)
	if err != nil {
		return nil, fmt.Errorf("query building_display_configs: %w", err)
	}
	defer rows.Close()

	var configs []*model.BuildingDisplayConfig
	for rows.Next() {
		c, err := scanBuildingDisplayConfig(rows)
		if err != nil {
			return nil, fmt.Errorf("scan building_display_config: %w", err)
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}

// GetByKingdom returns all building display configs for a specific kingdom.
func (r *BuildingDisplayConfigRepo) GetByKingdom(ctx context.Context, kingdom string) ([]*model.BuildingDisplayConfig, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, building_type, kingdom, display_name, description, default_icon, sprite_path, updated_at
		 FROM building_display_configs WHERE kingdom = ? ORDER BY building_type`, kingdom)
	if err != nil {
		return nil, fmt.Errorf("query building_display_configs by kingdom: %w", err)
	}
	defer rows.Close()

	var configs []*model.BuildingDisplayConfig
	for rows.Next() {
		c, err := scanBuildingDisplayConfig(rows)
		if err != nil {
			return nil, fmt.Errorf("scan building_display_config: %w", err)
		}
		configs = append(configs, c)
	}
	return configs, rows.Err()
}

// GetByID returns a single building display config by its primary key.
func (r *BuildingDisplayConfigRepo) GetByID(ctx context.Context, id int64) (*model.BuildingDisplayConfig, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, building_type, kingdom, display_name, description, default_icon, sprite_path, updated_at
		 FROM building_display_configs WHERE id = ?`, id)
	c, err := scanBuildingDisplayConfig(row)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return c, err
}

// Update persists changes to display_name, description, and default_icon for a config row.
func (r *BuildingDisplayConfigRepo) Update(ctx context.Context, cfg *model.BuildingDisplayConfig) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE building_display_configs SET display_name = ?, description = ?, default_icon = ?, updated_at = CURRENT_TIMESTAMP
		 WHERE id = ?`,
		cfg.DisplayName, cfg.Description, cfg.DefaultIcon, cfg.ID)
	if err != nil {
		return fmt.Errorf("update building_display_config %d: %w", cfg.ID, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

// UpdateSprite sets the sprite_path for the given config row.
func (r *BuildingDisplayConfigRepo) UpdateSprite(ctx context.Context, id int64, spritePath *string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE building_display_configs SET sprite_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		spritePath, id)
	if err != nil {
		return fmt.Errorf("update sprite for building_display_config %d: %w", id, err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

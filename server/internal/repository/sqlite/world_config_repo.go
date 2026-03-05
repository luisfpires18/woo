package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type worldConfigRepo struct {
	db *sql.DB
}

// NewWorldConfigRepo creates a new SQLite-backed WorldConfigRepository.
func NewWorldConfigRepo(db *sql.DB) *worldConfigRepo {
	return &worldConfigRepo{db: db}
}

func (r *worldConfigRepo) Get(ctx context.Context, key string) (*model.WorldConfig, error) {
	var cfg model.WorldConfig
	var updatedAtStr string
	err := r.db.QueryRowContext(ctx,
		`SELECT key, value, description, updated_at FROM world_config WHERE key = ?`, key,
	).Scan(&cfg.Key, &cfg.Value, &cfg.Description, &updatedAtStr)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get world config %q: %w", key, err)
	}
	cfg.UpdatedAt, _ = parseTime(updatedAtStr)
	return &cfg, nil
}

func (r *worldConfigRepo) GetAll(ctx context.Context) ([]*model.WorldConfig, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT key, value, description, updated_at FROM world_config ORDER BY key`,
	)
	if err != nil {
		return nil, fmt.Errorf("list world config: %w", err)
	}
	defer rows.Close()

	var configs []*model.WorldConfig
	for rows.Next() {
		var cfg model.WorldConfig
		var updatedAtStr string
		if err := rows.Scan(&cfg.Key, &cfg.Value, &cfg.Description, &updatedAtStr); err != nil {
			return nil, fmt.Errorf("scan world config: %w", err)
		}
		cfg.UpdatedAt, _ = parseTime(updatedAtStr)
		configs = append(configs, &cfg)
	}
	return configs, rows.Err()
}

func (r *worldConfigRepo) Set(ctx context.Context, key, value string) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE world_config SET value = ?, updated_at = datetime('now') WHERE key = ?`,
		value, key,
	)
	if err != nil {
		return fmt.Errorf("set world config %q: %w", key, err)
	}
	rows, _ := result.RowsAffected()
	if rows == 0 {
		return model.ErrNotFound
	}
	return nil
}

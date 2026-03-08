package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

// GameAssetRepo implements repository.GameAssetRepository using SQLite.
type GameAssetRepo struct {
	db *sql.DB
}

// NewGameAssetRepo creates a new GameAssetRepo.
func NewGameAssetRepo(db *sql.DB) *GameAssetRepo {
	return &GameAssetRepo{db: db}
}

func scanGameAsset(row interface{ Scan(dest ...any) error }) (*model.GameAsset, error) {
	var a model.GameAsset
	var spritePath sql.NullString
	var updatedAtStr string
	err := row.Scan(
		&a.ID, &a.Category, &a.DisplayName, &a.DefaultIcon,
		&spritePath, &a.SpriteWidth, &a.SpriteHeight, &updatedAtStr,
	)
	if err != nil {
		return nil, err
	}
	if spritePath.Valid {
		a.SpritePath = &spritePath.String
	}
	a.UpdatedAt, _ = parseTime(updatedAtStr)
	return &a, nil
}

// GetAll returns every game asset ordered by category then ID.
func (r *GameAssetRepo) GetAll(ctx context.Context) ([]*model.GameAsset, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, category, display_name, default_icon, sprite_path, sprite_width, sprite_height, updated_at
		 FROM game_assets ORDER BY category, id`)
	if err != nil {
		return nil, fmt.Errorf("query game_assets: %w", err)
	}
	defer rows.Close()

	var assets []*model.GameAsset
	for rows.Next() {
		a, err := scanGameAsset(rows)
		if err != nil {
			return nil, fmt.Errorf("scan game_asset: %w", err)
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

// GetByID returns a single game asset by its canonical ID.
func (r *GameAssetRepo) GetByID(ctx context.Context, id string) (*model.GameAsset, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, category, display_name, default_icon, sprite_path, sprite_width, sprite_height, updated_at
		 FROM game_assets WHERE id = ?`, id)
	a, err := scanGameAsset(row)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	return a, err
}

// GetByCategory returns all game assets of the given category.
func (r *GameAssetRepo) GetByCategory(ctx context.Context, category string) ([]*model.GameAsset, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, category, display_name, default_icon, sprite_path, sprite_width, sprite_height, updated_at
		 FROM game_assets WHERE category = ? ORDER BY id`, category)
	if err != nil {
		return nil, fmt.Errorf("query game_assets by category: %w", err)
	}
	defer rows.Close()

	var assets []*model.GameAsset
	for rows.Next() {
		a, err := scanGameAsset(rows)
		if err != nil {
			return nil, fmt.Errorf("scan game_asset: %w", err)
		}
		assets = append(assets, a)
	}
	return assets, rows.Err()
}

// UpdateSprite sets the sprite_path for the given asset.
func (r *GameAssetRepo) UpdateSprite(ctx context.Context, id string, spritePath *string) error {
	res, err := r.db.ExecContext(ctx,
		`UPDATE game_assets SET sprite_path = ?, updated_at = CURRENT_TIMESTAMP WHERE id = ?`,
		spritePath, id)
	if err != nil {
		return fmt.Errorf("update sprite: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

// Create inserts a new game asset (e.g., when adding units later).
func (r *GameAssetRepo) Create(ctx context.Context, a *model.GameAsset) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		a.ID, a.Category, a.DisplayName, a.DefaultIcon, a.SpriteWidth, a.SpriteHeight)
	if err != nil {
		return fmt.Errorf("create game_asset: %w", err)
	}
	return nil
}

// Delete removes a game asset row by ID.
func (r *GameAssetRepo) Delete(ctx context.Context, id string) error {
	res, err := r.db.ExecContext(ctx, `DELETE FROM game_assets WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete game_asset: %w", err)
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		return model.ErrNotFound
	}
	return nil
}

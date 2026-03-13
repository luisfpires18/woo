package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type campTemplateRepo struct {
	db *sql.DB
}

// NewCampTemplateRepo creates a new SQLite-backed CampTemplateRepository.
func NewCampTemplateRepo(db *sql.DB) *campTemplateRepo {
	return &campTemplateRepo{db: db}
}

func (r *campTemplateRepo) Create(ctx context.Context, ct *model.CampTemplate) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO camp_templates (name, tier, sprite_key, description, reward_table_id, updated_by)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		ct.Name, ct.Tier, ct.SpriteKey, ct.Description, ct.RewardTableID, ct.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("create camp template: %w", err)
	}
	id, _ := result.LastInsertId()
	ct.ID = id
	return nil
}

func (r *campTemplateRepo) GetByID(ctx context.Context, id int64) (*model.CampTemplate, error) {
	var ct model.CampTemplate
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, tier, sprite_key, description, reward_table_id, created_at, updated_at, updated_by
		 FROM camp_templates WHERE id = ?`, id,
	).Scan(&ct.ID, &ct.Name, &ct.Tier, &ct.SpriteKey, &ct.Description, &ct.RewardTableID,
		&ct.CreatedAt, &ct.UpdatedAt, &ct.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get camp template %d: %w", id, err)
	}
	return &ct, nil
}

func (r *campTemplateRepo) GetAll(ctx context.Context) ([]*model.CampTemplate, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, tier, sprite_key, description, reward_table_id, created_at, updated_at, updated_by
		 FROM camp_templates ORDER BY tier ASC, id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list camp templates: %w", err)
	}
	defer rows.Close()

	var templates []*model.CampTemplate
	for rows.Next() {
		var ct model.CampTemplate
		if err := rows.Scan(&ct.ID, &ct.Name, &ct.Tier, &ct.SpriteKey, &ct.Description, &ct.RewardTableID,
			&ct.CreatedAt, &ct.UpdatedAt, &ct.UpdatedBy); err != nil {
			return nil, fmt.Errorf("scan camp template: %w", err)
		}
		templates = append(templates, &ct)
	}
	return templates, rows.Err()
}

func (r *campTemplateRepo) Update(ctx context.Context, ct *model.CampTemplate) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE camp_templates SET name = ?, tier = ?, sprite_key = ?, description = ?,
		 reward_table_id = ?, updated_at = datetime('now'), updated_by = ?
		 WHERE id = ?`,
		ct.Name, ct.Tier, ct.SpriteKey, ct.Description, ct.RewardTableID, ct.UpdatedBy, ct.ID,
	)
	if err != nil {
		return fmt.Errorf("update camp template %d: %w", ct.ID, err)
	}
	return nil
}

func (r *campTemplateRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM camp_templates WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete camp template %d: %w", id, err)
	}
	return nil
}

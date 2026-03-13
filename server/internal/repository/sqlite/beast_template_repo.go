package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type beastTemplateRepo struct {
	db *sql.DB
}

// NewBeastTemplateRepo creates a new SQLite-backed BeastTemplateRepository.
func NewBeastTemplateRepo(db *sql.DB) *beastTemplateRepo {
	return &beastTemplateRepo{db: db}
}

func (r *beastTemplateRepo) Create(ctx context.Context, bt *model.BeastTemplate) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO beast_templates (name, sprite_key, hp, attack_power, attack_interval, defense_percent, crit_chance_percent, description, updated_by)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		bt.Name, bt.SpriteKey, bt.HP, bt.AttackPower, bt.AttackInterval,
		bt.DefensePercent, bt.CritChancePercent, bt.Description, bt.UpdatedBy,
	)
	if err != nil {
		return fmt.Errorf("create beast template: %w", err)
	}
	id, _ := result.LastInsertId()
	bt.ID = id
	return nil
}

func (r *beastTemplateRepo) GetByID(ctx context.Context, id int64) (*model.BeastTemplate, error) {
	var bt model.BeastTemplate
	err := r.db.QueryRowContext(ctx,
		`SELECT id, name, sprite_key, hp, attack_power, attack_interval, defense_percent, crit_chance_percent, description, created_at, updated_at, updated_by
		 FROM beast_templates WHERE id = ?`, id,
	).Scan(&bt.ID, &bt.Name, &bt.SpriteKey, &bt.HP, &bt.AttackPower, &bt.AttackInterval,
		&bt.DefensePercent, &bt.CritChancePercent, &bt.Description, &bt.CreatedAt, &bt.UpdatedAt, &bt.UpdatedBy)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get beast template %d: %w", id, err)
	}
	return &bt, nil
}

func (r *beastTemplateRepo) GetAll(ctx context.Context) ([]*model.BeastTemplate, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, name, sprite_key, hp, attack_power, attack_interval, defense_percent, crit_chance_percent, description, created_at, updated_at, updated_by
		 FROM beast_templates ORDER BY id ASC`)
	if err != nil {
		return nil, fmt.Errorf("list beast templates: %w", err)
	}
	defer rows.Close()

	var templates []*model.BeastTemplate
	for rows.Next() {
		var bt model.BeastTemplate
		if err := rows.Scan(&bt.ID, &bt.Name, &bt.SpriteKey, &bt.HP, &bt.AttackPower, &bt.AttackInterval,
			&bt.DefensePercent, &bt.CritChancePercent, &bt.Description, &bt.CreatedAt, &bt.UpdatedAt, &bt.UpdatedBy); err != nil {
			return nil, fmt.Errorf("scan beast template: %w", err)
		}
		templates = append(templates, &bt)
	}
	return templates, rows.Err()
}

func (r *beastTemplateRepo) Update(ctx context.Context, bt *model.BeastTemplate) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE beast_templates SET name = ?, sprite_key = ?, hp = ?, attack_power = ?, attack_interval = ?,
		 defense_percent = ?, crit_chance_percent = ?, description = ?, updated_at = datetime('now'), updated_by = ?
		 WHERE id = ?`,
		bt.Name, bt.SpriteKey, bt.HP, bt.AttackPower, bt.AttackInterval,
		bt.DefensePercent, bt.CritChancePercent, bt.Description, bt.UpdatedBy, bt.ID,
	)
	if err != nil {
		return fmt.Errorf("update beast template %d: %w", bt.ID, err)
	}
	return nil
}

func (r *beastTemplateRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM beast_templates WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete beast template %d: %w", id, err)
	}
	return nil
}

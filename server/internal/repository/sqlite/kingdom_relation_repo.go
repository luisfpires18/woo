package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type kingdomRelationRepo struct {
	db *sql.DB
}

// NewKingdomRelationRepo creates a new SQLite-backed KingdomRelationRepository.
func NewKingdomRelationRepo(db *sql.DB) *kingdomRelationRepo {
	return &kingdomRelationRepo{db: db}
}

func (r *kingdomRelationRepo) GetAll(ctx context.Context) ([]*model.KingdomRelation, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT kingdom_a, kingdom_b, standing, status, updated_at FROM kingdom_relations ORDER BY kingdom_a, kingdom_b`,
	)
	if err != nil {
		return nil, fmt.Errorf("query kingdom relations: %w", err)
	}
	defer rows.Close()

	var rels []*model.KingdomRelation
	for rows.Next() {
		rel := &model.KingdomRelation{}
		var updatedAtStr string
		if err := rows.Scan(&rel.KingdomA, &rel.KingdomB, &rel.Standing, &rel.Status, &updatedAtStr); err != nil {
			return nil, fmt.Errorf("scan kingdom relation: %w", err)
		}
		rel.UpdatedAt, _ = parseTime(updatedAtStr)
		rels = append(rels, rel)
	}
	return rels, rows.Err()
}

func (r *kingdomRelationRepo) Get(ctx context.Context, kingdomA, kingdomB string) (*model.KingdomRelation, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT kingdom_a, kingdom_b, standing, status, updated_at
		 FROM kingdom_relations WHERE kingdom_a = ? AND kingdom_b = ?`,
		kingdomA, kingdomB,
	)

	rel := &model.KingdomRelation{}
	var updatedAtStr string
	err := row.Scan(&rel.KingdomA, &rel.KingdomB, &rel.Standing, &rel.Status, &updatedAtStr)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan kingdom relation: %w", err)
	}
	rel.UpdatedAt, _ = parseTime(updatedAtStr)
	return rel, nil
}

func (r *kingdomRelationRepo) Upsert(ctx context.Context, rel *model.KingdomRelation) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO kingdom_relations (kingdom_a, kingdom_b, standing, status, updated_at)
		 VALUES (?, ?, ?, ?, datetime('now'))
		 ON CONFLICT (kingdom_a, kingdom_b)
		 DO UPDATE SET standing = ?, status = ?, updated_at = datetime('now')`,
		rel.KingdomA, rel.KingdomB, rel.Standing, rel.Status,
		rel.Standing, rel.Status,
	)
	if err != nil {
		return fmt.Errorf("upsert kingdom relation %s↔%s: %w", rel.KingdomA, rel.KingdomB, err)
	}
	return nil
}

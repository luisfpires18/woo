package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type campBeastSlotRepo struct {
	db *sql.DB
}

// NewCampBeastSlotRepo creates a new SQLite-backed CampBeastSlotRepository.
func NewCampBeastSlotRepo(db *sql.DB) *campBeastSlotRepo {
	return &campBeastSlotRepo{db: db}
}

func (r *campBeastSlotRepo) Create(ctx context.Context, slot *model.CampBeastSlot) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO camp_beast_slots (camp_template_id, beast_template_id, min_count, max_count)
		 VALUES (?, ?, ?, ?)`,
		slot.CampTemplateID, slot.BeastTemplateID, slot.MinCount, slot.MaxCount,
	)
	if err != nil {
		return fmt.Errorf("create camp beast slot: %w", err)
	}
	id, _ := result.LastInsertId()
	slot.ID = id
	return nil
}

func (r *campBeastSlotRepo) GetByCampTemplateID(ctx context.Context, campTemplateID int64) ([]*model.CampBeastSlot, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, camp_template_id, beast_template_id, min_count, max_count
		 FROM camp_beast_slots WHERE camp_template_id = ? ORDER BY id ASC`, campTemplateID,
	)
	if err != nil {
		return nil, fmt.Errorf("list camp beast slots for template %d: %w", campTemplateID, err)
	}
	defer rows.Close()

	var slots []*model.CampBeastSlot
	for rows.Next() {
		var s model.CampBeastSlot
		if err := rows.Scan(&s.ID, &s.CampTemplateID, &s.BeastTemplateID, &s.MinCount, &s.MaxCount); err != nil {
			return nil, fmt.Errorf("scan camp beast slot: %w", err)
		}
		slots = append(slots, &s)
	}
	return slots, rows.Err()
}

func (r *campBeastSlotRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM camp_beast_slots WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete camp beast slot %d: %w", id, err)
	}
	return nil
}

func (r *campBeastSlotRepo) DeleteByCampTemplateID(ctx context.Context, campTemplateID int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM camp_beast_slots WHERE camp_template_id = ?`, campTemplateID)
	if err != nil {
		return fmt.Errorf("delete camp beast slots for template %d: %w", campTemplateID, err)
	}
	return nil
}

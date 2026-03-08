package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type troopRepo struct {
	db *sql.DB
}

// NewTroopRepo creates a new SQLite-backed TroopRepository.
func NewTroopRepo(db *sql.DB) *troopRepo {
	return &troopRepo{db: db}
}

func (r *troopRepo) GetByVillageID(ctx context.Context, villageID int64) ([]*model.Troop, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, village_id, type, quantity, status
		 FROM troops WHERE village_id = ? ORDER BY type ASC`, villageID,
	)
	if err != nil {
		return nil, fmt.Errorf("list troops for village %d: %w", villageID, err)
	}
	defer rows.Close()

	var troops []*model.Troop
	for rows.Next() {
		var t model.Troop
		if err := rows.Scan(&t.ID, &t.VillageID, &t.Type, &t.Quantity, &t.Status); err != nil {
			return nil, fmt.Errorf("scan troop row: %w", err)
		}
		troops = append(troops, &t)
	}
	return troops, rows.Err()
}

func (r *troopRepo) GetByVillageAndType(ctx context.Context, villageID int64, troopType string) (*model.Troop, error) {
	var t model.Troop
	err := r.db.QueryRowContext(ctx,
		`SELECT id, village_id, type, quantity, status
		 FROM troops WHERE village_id = ? AND type = ?`, villageID, troopType,
	).Scan(&t.ID, &t.VillageID, &t.Type, &t.Quantity, &t.Status)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get troop by village and type: %w", err)
	}
	return &t, nil
}

func (r *troopRepo) Upsert(ctx context.Context, troop *model.Troop) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO troops (village_id, type, quantity, status)
		 VALUES (?, ?, ?, ?)
		 ON CONFLICT(village_id, type) DO UPDATE SET quantity = quantity + excluded.quantity`,
		troop.VillageID, troop.Type, troop.Quantity, troop.Status,
	)
	if err != nil {
		return fmt.Errorf("upsert troop: %w", err)
	}
	return nil
}

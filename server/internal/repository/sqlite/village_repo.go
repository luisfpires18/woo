package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type villageRepo struct {
	db *sql.DB
}

// NewVillageRepo creates a new SQLite-backed VillageRepository.
func NewVillageRepo(db *sql.DB) *villageRepo {
	return &villageRepo{db: db}
}

func (r *villageRepo) Create(ctx context.Context, village *model.Village) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO villages (player_id, name, x, y, is_capital, created_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		village.PlayerID, village.Name, village.X, village.Y, village.IsCapital, village.CreatedAt,
	)
	if err != nil {
		return fmt.Errorf("insert village: %w", err)
	}
	id, _ := result.LastInsertId()
	village.ID = id
	return nil
}

func (r *villageRepo) GetByID(ctx context.Context, id int64) (*model.Village, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, player_id, name, x, y, is_capital, created_at
		 FROM villages WHERE id = ?`, id,
	)
	return scanVillage(row)
}

func (r *villageRepo) ListByPlayerID(ctx context.Context, playerID int64) ([]*model.Village, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, player_id, name, x, y, is_capital, created_at
		 FROM villages WHERE player_id = ? ORDER BY created_at ASC`, playerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list villages for player %d: %w", playerID, err)
	}
	defer rows.Close()

	var villages []*model.Village
	for rows.Next() {
		var v model.Village
		var createdAtStr string
		if err := rows.Scan(&v.ID, &v.PlayerID, &v.Name, &v.X, &v.Y, &v.IsCapital, &createdAtStr); err != nil {
			return nil, fmt.Errorf("scan village row: %w", err)
		}
		v.CreatedAt, _ = parseTime(createdAtStr)
		villages = append(villages, &v)
	}
	return villages, rows.Err()
}

func (r *villageRepo) Update(ctx context.Context, village *model.Village) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE villages SET name = ?, is_capital = ? WHERE id = ?`,
		village.Name, village.IsCapital, village.ID,
	)
	if err != nil {
		return fmt.Errorf("update village %d: %w", village.ID, err)
	}
	return nil
}

func (r *villageRepo) GetByCoordinates(ctx context.Context, x, y int) (*model.Village, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, player_id, name, x, y, is_capital, created_at
		 FROM villages WHERE x = ? AND y = ?`, x, y,
	)
	return scanVillage(row)
}

func scanVillage(row *sql.Row) (*model.Village, error) {
	var v model.Village
	var createdAtStr string
	err := row.Scan(&v.ID, &v.PlayerID, &v.Name, &v.X, &v.Y, &v.IsCapital, &createdAtStr)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan village: %w", err)
	}
	v.CreatedAt, _ = parseTime(createdAtStr)
	return &v, nil
}

func (r *villageRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM villages`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count villages: %w", err)
	}
	return count, nil
}

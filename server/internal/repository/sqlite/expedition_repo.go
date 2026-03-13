package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/luisfpires18/woo/internal/model"
)

type expeditionRepo struct {
	db *sql.DB
}

// NewExpeditionRepo creates a new SQLite-backed ExpeditionRepository.
func NewExpeditionRepo(db *sql.DB) *expeditionRepo {
	return &expeditionRepo{db: db}
}

func (r *expeditionRepo) Create(ctx context.Context, exp *model.Expedition) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO expeditions (player_id, village_id, camp_id, troops_json, departed_at, arrives_at, returns_at, status, season_id)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		exp.PlayerID, exp.VillageID, exp.CampID, exp.TroopsJSON,
		exp.DepartedAt, exp.ArrivesAt, exp.ReturnsAt, exp.Status, exp.SeasonID,
	)
	if err != nil {
		return fmt.Errorf("create expedition: %w", err)
	}
	id, _ := result.LastInsertId()
	exp.ID = id
	return nil
}

func (r *expeditionRepo) GetByID(ctx context.Context, id int64) (*model.Expedition, error) {
	var exp model.Expedition
	var returnsAt sql.NullString
	err := r.db.QueryRowContext(ctx,
		`SELECT id, player_id, village_id, camp_id, troops_json, departed_at, arrives_at, returns_at, status, season_id
		 FROM expeditions WHERE id = ?`, id,
	).Scan(&exp.ID, &exp.PlayerID, &exp.VillageID, &exp.CampID, &exp.TroopsJSON,
		&exp.DepartedAt, &exp.ArrivesAt, &returnsAt, &exp.Status, &exp.SeasonID)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get expedition %d: %w", id, err)
	}
	if returnsAt.Valid {
		exp.ReturnsAt = returnsAt.String
	}
	return &exp, nil
}

func (r *expeditionRepo) GetByPlayerID(ctx context.Context, playerID int64) ([]*model.Expedition, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, player_id, village_id, camp_id, troops_json, departed_at, arrives_at, returns_at, status, season_id
		 FROM expeditions WHERE player_id = ? ORDER BY departed_at DESC`, playerID,
	)
	if err != nil {
		return nil, fmt.Errorf("list expeditions for player %d: %w", playerID, err)
	}
	defer rows.Close()
	return scanExpeditions(rows)
}

func (r *expeditionRepo) GetArrivedExpeditions(ctx context.Context, now time.Time) ([]*model.Expedition, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, player_id, village_id, camp_id, troops_json, departed_at, arrives_at, returns_at, status, season_id
		 FROM expeditions WHERE status = 'marching' AND arrives_at <= ?
		 ORDER BY arrives_at ASC`,
		now.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("get arrived expeditions: %w", err)
	}
	defer rows.Close()
	return scanExpeditions(rows)
}

func (r *expeditionRepo) GetReturningExpeditions(ctx context.Context, now time.Time) ([]*model.Expedition, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, player_id, village_id, camp_id, troops_json, departed_at, arrives_at, returns_at, status, season_id
		 FROM expeditions WHERE status = 'returning' AND returns_at <= ?
		 ORDER BY returns_at ASC`,
		now.UTC().Format(time.RFC3339),
	)
	if err != nil {
		return nil, fmt.Errorf("get returning expeditions: %w", err)
	}
	defer rows.Close()
	return scanExpeditions(rows)
}

func (r *expeditionRepo) UpdateStatus(ctx context.Context, id int64, status string) error {
	_, err := r.db.ExecContext(ctx, `UPDATE expeditions SET status = ? WHERE id = ?`, status, id)
	if err != nil {
		return fmt.Errorf("update expedition %d status: %w", id, err)
	}
	return nil
}

func (r *expeditionRepo) UpdateReturnsAt(ctx context.Context, id int64, returnsAt time.Time) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE expeditions SET returns_at = ? WHERE id = ?`,
		returnsAt.UTC().Format(time.RFC3339), id,
	)
	if err != nil {
		return fmt.Errorf("update expedition %d returns_at: %w", id, err)
	}
	return nil
}

func (r *expeditionRepo) Update(ctx context.Context, exp *model.Expedition) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE expeditions SET troops_json = ?, status = ?, returns_at = ? WHERE id = ?`,
		exp.TroopsJSON, exp.Status, exp.ReturnsAt, exp.ID,
	)
	if err != nil {
		return fmt.Errorf("update expedition %d: %w", exp.ID, err)
	}
	return nil
}

func scanExpeditions(rows *sql.Rows) ([]*model.Expedition, error) {
	var expeditions []*model.Expedition
	for rows.Next() {
		var exp model.Expedition
		var returnsAt sql.NullString
		if err := rows.Scan(&exp.ID, &exp.PlayerID, &exp.VillageID, &exp.CampID, &exp.TroopsJSON,
			&exp.DepartedAt, &exp.ArrivesAt, &returnsAt, &exp.Status, &exp.SeasonID); err != nil {
			return nil, fmt.Errorf("scan expedition: %w", err)
		}
		if returnsAt.Valid {
			exp.ReturnsAt = returnsAt.String
		}
		expeditions = append(expeditions, &exp)
	}
	return expeditions, rows.Err()
}

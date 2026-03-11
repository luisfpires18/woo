package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type playerEconomyRepo struct {
	db *sql.DB
}

// NewPlayerEconomyRepo creates a new SQLite-backed PlayerEconomyRepository.
func NewPlayerEconomyRepo(db *sql.DB) *playerEconomyRepo {
	return &playerEconomyRepo{db: db}
}

func (r *playerEconomyRepo) Create(ctx context.Context, playerID int64, gold float64) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO player_economy (player_id, gold) VALUES (?, ?)`,
		playerID, gold,
	)
	if err != nil {
		return fmt.Errorf("create player economy for player %d: %w", playerID, err)
	}
	return nil
}

func (r *playerEconomyRepo) GetByPlayerID(ctx context.Context, playerID int64) (*model.PlayerEconomy, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT player_id, gold FROM player_economy WHERE player_id = ?`, playerID,
	)
	var pe model.PlayerEconomy
	err := row.Scan(&pe.PlayerID, &pe.Gold)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get player economy for player %d: %w", playerID, err)
	}
	return &pe, nil
}

func (r *playerEconomyRepo) UpdateGold(ctx context.Context, playerID int64, newGold float64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE player_economy SET gold = ? WHERE player_id = ?`,
		newGold, playerID,
	)
	if err != nil {
		return fmt.Errorf("update gold for player %d: %w", playerID, err)
	}
	return nil
}

// DeductGold atomically checks and deducts gold. Returns ErrInsufficientGold if balance is too low.
func (r *playerEconomyRepo) DeductGold(ctx context.Context, playerID int64, amount float64) error {
	result, err := r.db.ExecContext(ctx,
		`UPDATE player_economy SET gold = gold - ? WHERE player_id = ? AND gold >= ?`,
		amount, playerID, amount,
	)
	if err != nil {
		return fmt.Errorf("deduct gold for player %d: %w", playerID, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check gold deduction result for player %d: %w", playerID, err)
	}
	if rows == 0 {
		return model.ErrInsufficientGold
	}
	return nil
}

// DeductGoldTx atomically checks and deducts gold within an existing sql.Tx.
func (r *playerEconomyRepo) DeductGoldTx(ctx context.Context, txIface interface{}, playerID int64, amount float64) error {
	tx, ok := txIface.(*sql.Tx)
	if !ok {
		return fmt.Errorf("DeductGoldTx: expected *sql.Tx, got %T", txIface)
	}
	result, err := tx.ExecContext(ctx,
		`UPDATE player_economy SET gold = gold - ? WHERE player_id = ? AND gold >= ?`,
		amount, playerID, amount,
	)
	if err != nil {
		return fmt.Errorf("deduct gold (tx) for player %d: %w", playerID, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check gold deduction result (tx) for player %d: %w", playerID, err)
	}
	if rows == 0 {
		return model.ErrInsufficientGold
	}
	return nil
}

package sqlite

import (
	"context"
	"database/sql"
	"fmt"
)

// WithTx executes fn within a database transaction. If fn returns an error the
// transaction is rolled back; otherwise it is committed.
func WithTx(ctx context.Context, db *sql.DB, fn func(tx *sql.Tx) error) error {
	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	if err := fn(tx); err != nil {
		return err
	}
	return tx.Commit()
}

// UpdateResourcesTx updates the resources row for a village within the given transaction.
func UpdateResourcesTx(ctx context.Context, tx *sql.Tx, villageID int64, food, water, lumber, stone, foodRate, waterRate, lumberRate, stoneRate, foodConsumption, maxStorage float64, lastUpdated string) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE resources SET food = ?, water = ?, lumber = ?, stone = ?, food_rate = ?, water_rate = ?, lumber_rate = ?, stone_rate = ?, food_consumption = ?, max_storage = ?, last_updated = ?
		 WHERE village_id = ?`,
		food, water, lumber, stone, foodRate, waterRate, lumberRate, stoneRate,
		foodConsumption, maxStorage, lastUpdated, villageID,
	)
	if err != nil {
		return fmt.Errorf("update resources (tx) for village %d: %w", villageID, err)
	}
	return nil
}

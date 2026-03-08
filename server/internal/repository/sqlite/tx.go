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

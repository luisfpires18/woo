package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/luisfpires18/woo/internal/model"
)

type trainingQueueRepo struct {
	db *sql.DB
}

// NewTrainingQueueRepo creates a new SQLite-backed TrainingQueueRepository.
func NewTrainingQueueRepo(db *sql.DB) *trainingQueueRepo {
	return &trainingQueueRepo{db: db}
}

func (r *trainingQueueRepo) Insert(ctx context.Context, item *model.TrainingQueue) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO training_queue (village_id, troop_type, quantity, each_duration_sec, started_at, completes_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		item.VillageID, item.TroopType, item.Quantity, item.EachDurationSec,
		item.StartedAt.UTC().Format("2006-01-02 15:04:05"),
		item.CompletesAt.UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("insert training queue item: %w", err)
	}
	id, _ := result.LastInsertId()
	item.ID = id
	return nil
}

func (r *trainingQueueRepo) GetByID(ctx context.Context, id int64) (*model.TrainingQueue, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, village_id, troop_type, quantity, each_duration_sec, started_at, completes_at
		 FROM training_queue WHERE id = ?`, id,
	)
	var item model.TrainingQueue
	var startedStr, completesStr string
	if err := row.Scan(
		&item.ID, &item.VillageID, &item.TroopType,
		&item.Quantity, &item.EachDurationSec,
		&startedStr, &completesStr,
	); err != nil {
		if err == sql.ErrNoRows {
			return nil, model.ErrNotFound
		}
		return nil, fmt.Errorf("get training queue item %d: %w", id, err)
	}
	item.StartedAt, _ = parseTime(startedStr)
	item.CompletesAt, _ = parseTime(completesStr)
	return &item, nil
}

func (r *trainingQueueRepo) GetByVillageID(ctx context.Context, villageID int64) ([]*model.TrainingQueue, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, village_id, troop_type, quantity, each_duration_sec, started_at, completes_at
		 FROM training_queue WHERE village_id = ? ORDER BY completes_at ASC`, villageID,
	)
	if err != nil {
		return nil, fmt.Errorf("list training queue for village %d: %w", villageID, err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *trainingQueueRepo) GetNextCompleted(ctx context.Context, now time.Time) ([]*model.TrainingQueue, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, village_id, troop_type, quantity, each_duration_sec, started_at, completes_at
		 FROM training_queue WHERE completes_at <= ? ORDER BY completes_at ASC`,
		now.UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, fmt.Errorf("get completed training queue items: %w", err)
	}
	defer rows.Close()

	return r.scanRows(rows)
}

func (r *trainingQueueRepo) Update(ctx context.Context, item *model.TrainingQueue) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE training_queue SET quantity = ?, completes_at = ? WHERE id = ?`,
		item.Quantity,
		item.CompletesAt.UTC().Format("2006-01-02 15:04:05"),
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("update training queue item %d: %w", item.ID, err)
	}
	return nil
}

func (r *trainingQueueRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM training_queue WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete training queue item %d: %w", id, err)
	}
	return nil
}

func (r *trainingQueueRepo) scanRows(rows *sql.Rows) ([]*model.TrainingQueue, error) {
	var items []*model.TrainingQueue
	for rows.Next() {
		var item model.TrainingQueue
		var startedStr, completesStr string
		if err := rows.Scan(
			&item.ID, &item.VillageID, &item.TroopType,
			&item.Quantity, &item.EachDurationSec,
			&startedStr, &completesStr,
		); err != nil {
			return nil, fmt.Errorf("scan training queue row: %w", err)
		}
		item.StartedAt, _ = parseTime(startedStr)
		item.CompletesAt, _ = parseTime(completesStr)
		items = append(items, &item)
	}
	return items, rows.Err()
}

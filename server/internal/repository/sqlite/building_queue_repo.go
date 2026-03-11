package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/luisfpires18/woo/internal/model"
)

type buildingQueueRepo struct {
	db *sql.DB
}

// NewBuildingQueueRepo creates a new SQLite-backed BuildingQueueRepository.
func NewBuildingQueueRepo(db *sql.DB) *buildingQueueRepo {
	return &buildingQueueRepo{db: db}
}

func (r *buildingQueueRepo) Insert(ctx context.Context, item *model.BuildingQueue) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO building_queue (village_id, building_type, target_level, started_at, completes_at)
		 VALUES (?, ?, ?, ?, ?)`,
		item.VillageID, item.BuildingType, item.TargetLevel,
		item.StartedAt.UTC().Format("2006-01-02 15:04:05"),
		item.CompletesAt.UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("insert building queue item: %w", err)
	}
	id, _ := result.LastInsertId()
	item.ID = id
	return nil
}

func (r *buildingQueueRepo) GetByID(ctx context.Context, id int64) (*model.BuildingQueue, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT id, village_id, building_type, target_level, started_at, completes_at
		 FROM building_queue WHERE id = ?`, id,
	)
	var item model.BuildingQueue
	var startedStr, completesStr string
	if err := row.Scan(&item.ID, &item.VillageID, &item.BuildingType, &item.TargetLevel, &startedStr, &completesStr); err != nil {
		return nil, fmt.Errorf("get building queue item %d: %w", id, err)
	}
	item.StartedAt, _ = parseTime(startedStr)
	item.CompletesAt, _ = parseTime(completesStr)
	return &item, nil
}

func (r *buildingQueueRepo) Update(ctx context.Context, item *model.BuildingQueue) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE building_queue SET completes_at = ? WHERE id = ?`,
		item.CompletesAt.UTC().Format("2006-01-02 15:04:05"), item.ID,
	)
	if err != nil {
		return fmt.Errorf("update building queue item %d: %w", item.ID, err)
	}
	return nil
}

func (r *buildingQueueRepo) GetByVillageID(ctx context.Context, villageID int64) ([]*model.BuildingQueue, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, village_id, building_type, target_level, started_at, completes_at
		 FROM building_queue WHERE village_id = ? ORDER BY completes_at ASC`, villageID,
	)
	if err != nil {
		return nil, fmt.Errorf("list building queue for village %d: %w", villageID, err)
	}
	defer rows.Close()

	var items []*model.BuildingQueue
	for rows.Next() {
		var item model.BuildingQueue
		var startedStr, completesStr string
		if err := rows.Scan(&item.ID, &item.VillageID, &item.BuildingType, &item.TargetLevel, &startedStr, &completesStr); err != nil {
			return nil, fmt.Errorf("scan building queue row: %w", err)
		}
		item.StartedAt, _ = parseTime(startedStr)
		item.CompletesAt, _ = parseTime(completesStr)
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (r *buildingQueueRepo) GetCompleted(ctx context.Context, now time.Time) ([]*model.BuildingQueue, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, village_id, building_type, target_level, started_at, completes_at
		 FROM building_queue WHERE completes_at <= ? ORDER BY completes_at ASC`,
		now.UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return nil, fmt.Errorf("get completed building queue items: %w", err)
	}
	defer rows.Close()

	var items []*model.BuildingQueue
	for rows.Next() {
		var item model.BuildingQueue
		var startedStr, completesStr string
		if err := rows.Scan(&item.ID, &item.VillageID, &item.BuildingType, &item.TargetLevel, &startedStr, &completesStr); err != nil {
			return nil, fmt.Errorf("scan building queue row: %w", err)
		}
		item.StartedAt, _ = parseTime(startedStr)
		item.CompletesAt, _ = parseTime(completesStr)
		items = append(items, &item)
	}
	return items, rows.Err()
}

func (r *buildingQueueRepo) Delete(ctx context.Context, id int64) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM building_queue WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete building queue item %d: %w", id, err)
	}
	return nil
}

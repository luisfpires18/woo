package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

// unitOfWork implements repository.UnitOfWork using SQLite transactions.
type unitOfWork struct {
	db *sql.DB
}

// NewUnitOfWork creates a new SQLite-backed UnitOfWork.
func NewUnitOfWork(db *sql.DB) *unitOfWork {
	return &unitOfWork{db: db}
}

// DeductResourcesAndInsertBuildQueue atomically deducts resources and inserts a build queue item.
func (u *unitOfWork) DeductResourcesAndInsertBuildQueue(ctx context.Context, villageID int64, res *model.Resources, item *model.BuildingQueue) error {
	return WithTx(ctx, u.db, func(tx *sql.Tx) error {
		if err := updateResourcesTx(ctx, tx, villageID, res); err != nil {
			return err
		}
		return insertBuildQueueTx(ctx, tx, item)
	})
}

// DeductResourcesAndInsertTrainQueue atomically deducts resources and inserts a training queue item.
func (u *unitOfWork) DeductResourcesAndInsertTrainQueue(ctx context.Context, villageID int64, res *model.Resources, item *model.TrainingQueue) error {
	return WithTx(ctx, u.db, func(tx *sql.Tx) error {
		if err := updateResourcesTx(ctx, tx, villageID, res); err != nil {
			return err
		}
		return insertTrainQueueTx(ctx, tx, item)
	})
}

// CompleteTrainingUnit atomically adds troops, updates resources, and advances/deletes the queue item.
func (u *unitOfWork) CompleteTrainingUnit(ctx context.Context, villageID int64, troopType string, addQty int, res *model.Resources, queueItem *model.TrainingQueue, deleteQueue bool) error {
	return WithTx(ctx, u.db, func(tx *sql.Tx) error {
		// 1. Add troops
		if err := upsertTroopTx(ctx, tx, villageID, troopType, addQty); err != nil {
			return err
		}
		// 2. Update resources (food consumption updated by caller)
		if err := updateResourcesTx(ctx, tx, villageID, res); err != nil {
			return err
		}
		// 3. Advance or delete queue item
		if deleteQueue {
			return deleteTrainQueueTx(ctx, tx, queueItem.ID)
		}
		return updateTrainQueueTx(ctx, tx, queueItem)
	})
}

// CompleteBuildingUpgrade atomically updates building level, refreshes resource rates, and deletes the queue item.
func (u *unitOfWork) CompleteBuildingUpgrade(ctx context.Context, villageID int64, building *model.Building, resources *model.Resources, queueID int64) error {
	return WithTx(ctx, u.db, func(tx *sql.Tx) error {
		// 1. Level up the building
		if err := updateBuildingLevelTx(ctx, tx, building); err != nil {
			return err
		}
		// 2. Update resources (rates recalculated by caller)
		if err := updateResourcesTx(ctx, tx, villageID, resources); err != nil {
			return err
		}
		// 3. Delete completed queue entry
		return deleteBuildQueueTx(ctx, tx, queueID)
	})
}

// ── Internal transactional helpers (unexported, only used within this package) ─

func updateResourcesTx(ctx context.Context, tx *sql.Tx, villageID int64, res *model.Resources) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE resources SET food = ?, water = ?, lumber = ?, stone = ?,
		 food_rate = ?, water_rate = ?, lumber_rate = ?, stone_rate = ?,
		 food_consumption = ?, max_storage = ?, last_updated = ?
		 WHERE village_id = ?`,
		res.Food, res.Water, res.Lumber, res.Stone,
		res.FoodRate, res.WaterRate, res.LumberRate, res.StoneRate,
		res.FoodConsumption, res.MaxStorage,
		res.LastUpdated.UTC().Format("2006-01-02 15:04:05"),
		villageID,
	)
	if err != nil {
		return fmt.Errorf("update resources (tx) for village %d: %w", villageID, err)
	}
	return nil
}

func insertBuildQueueTx(ctx context.Context, tx *sql.Tx, item *model.BuildingQueue) error {
	result, err := tx.ExecContext(ctx,
		`INSERT INTO building_queue (village_id, building_type, target_level, started_at, completes_at)
		 VALUES (?, ?, ?, ?, ?)`,
		item.VillageID, item.BuildingType, item.TargetLevel,
		item.StartedAt.UTC().Format("2006-01-02 15:04:05"),
		item.CompletesAt.UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("insert building queue item (tx): %w", err)
	}
	id, _ := result.LastInsertId()
	item.ID = id
	return nil
}

func insertTrainQueueTx(ctx context.Context, tx *sql.Tx, item *model.TrainingQueue) error {
	result, err := tx.ExecContext(ctx,
		`INSERT INTO training_queue (village_id, troop_type, quantity, each_duration_sec, started_at, completes_at)
		 VALUES (?, ?, ?, ?, ?, ?)`,
		item.VillageID, item.TroopType, item.Quantity, item.EachDurationSec,
		item.StartedAt.UTC().Format("2006-01-02 15:04:05"),
		item.CompletesAt.UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("insert training queue item (tx): %w", err)
	}
	id, _ := result.LastInsertId()
	item.ID = id
	return nil
}

func upsertTroopTx(ctx context.Context, tx *sql.Tx, villageID int64, troopType string, addQuantity int) error {
	_, err := tx.ExecContext(ctx,
		`INSERT INTO troops (village_id, type, quantity, status)
		 VALUES (?, ?, ?, 'stationed')
		 ON CONFLICT(village_id, type) DO UPDATE SET quantity = quantity + ?`,
		villageID, troopType, addQuantity, addQuantity,
	)
	if err != nil {
		return fmt.Errorf("upsert troop (tx): %w", err)
	}
	return nil
}

func updateTrainQueueTx(ctx context.Context, tx *sql.Tx, item *model.TrainingQueue) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE training_queue SET quantity = ?, completes_at = ? WHERE id = ?`,
		item.Quantity,
		item.CompletesAt.UTC().Format("2006-01-02 15:04:05"),
		item.ID,
	)
	if err != nil {
		return fmt.Errorf("update training queue item (tx) %d: %w", item.ID, err)
	}
	return nil
}

func deleteTrainQueueTx(ctx context.Context, tx *sql.Tx, id int64) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM training_queue WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete training queue item (tx) %d: %w", id, err)
	}
	return nil
}

func updateBuildingLevelTx(ctx context.Context, tx *sql.Tx, building *model.Building) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE buildings SET level = ? WHERE id = ?`,
		building.Level, building.ID,
	)
	if err != nil {
		return fmt.Errorf("update building level (tx) for building %d: %w", building.ID, err)
	}
	return nil
}

func deleteBuildQueueTx(ctx context.Context, tx *sql.Tx, id int64) error {
	_, err := tx.ExecContext(ctx, `DELETE FROM building_queue WHERE id = ?`, id)
	if err != nil {
		return fmt.Errorf("delete building queue item (tx) %d: %w", id, err)
	}
	return nil
}

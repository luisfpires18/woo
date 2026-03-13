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

// DeductResourcesGoldAndInsertBuildQueue atomically deducts village resources + player gold + inserts a build queue item.
func (u *unitOfWork) DeductResourcesGoldAndInsertBuildQueue(ctx context.Context, villageID int64, res *model.Resources, playerID int64, goldCost float64, item *model.BuildingQueue) error {
	return WithTx(ctx, u.db, func(tx *sql.Tx) error {
		if err := updateResourcesTx(ctx, tx, villageID, res); err != nil {
			return err
		}
		if goldCost > 0 {
			if err := deductGoldTx(ctx, tx, playerID, goldCost); err != nil {
				return err
			}
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

// DeductResourcesGoldAndInsertTrainQueue atomically deducts village resources + player gold + inserts a training queue item.
func (u *unitOfWork) DeductResourcesGoldAndInsertTrainQueue(ctx context.Context, villageID int64, res *model.Resources, playerID int64, goldCost float64, item *model.TrainingQueue) error {
	return WithTx(ctx, u.db, func(tx *sql.Tx) error {
		if err := updateResourcesTx(ctx, tx, villageID, res); err != nil {
			return err
		}
		if goldCost > 0 {
			if err := deductGoldTx(ctx, tx, playerID, goldCost); err != nil {
				return err
			}
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

// DeductTroopsAndCreateExpedition atomically deducts troops from a village and creates an expedition.
func (u *unitOfWork) DeductTroopsAndCreateExpedition(ctx context.Context, villageID int64, troopDeductions map[string]int, exp *model.Expedition) error {
	return WithTx(ctx, u.db, func(tx *sql.Tx) error {
		for troopType, qty := range troopDeductions {
			result, err := tx.ExecContext(ctx,
				`UPDATE troops SET quantity = quantity - ? WHERE village_id = ? AND type = ? AND quantity >= ?`,
				qty, villageID, troopType, qty,
			)
			if err != nil {
				return fmt.Errorf("deduct troop %s (tx): %w", troopType, err)
			}
			rows, err := result.RowsAffected()
			if err != nil {
				return fmt.Errorf("check troop deduction %s (tx): %w", troopType, err)
			}
			if rows == 0 {
				return fmt.Errorf("insufficient %s troops for expedition", troopType)
			}
		}

		res, err := tx.ExecContext(ctx,
			`INSERT INTO expeditions (player_id, village_id, camp_id, troops_json, departed_at, arrives_at, status, season_id)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			exp.PlayerID, exp.VillageID, exp.CampID, exp.TroopsJSON,
			exp.DepartedAt, exp.ArrivesAt, exp.Status, exp.SeasonID,
		)
		if err != nil {
			return fmt.Errorf("insert expedition (tx): %w", err)
		}
		id, err := res.LastInsertId()
		if err != nil {
			return fmt.Errorf("get expedition last insert id: %w", err)
		}
		exp.ID = id
		return nil
	})
}

// ReturnExpeditionTroops atomically adds surviving troops back and marks the expedition completed.
func (u *unitOfWork) ReturnExpeditionTroops(ctx context.Context, villageID int64, troopAdditions map[string]int, expeditionID int64) error {
	return WithTx(ctx, u.db, func(tx *sql.Tx) error {
		for troopType, qty := range troopAdditions {
			if qty <= 0 {
				continue
			}
			if err := upsertTroopTx(ctx, tx, villageID, troopType, qty); err != nil {
				return err
			}
		}

		_, err := tx.ExecContext(ctx,
			`UPDATE expeditions SET status = 'completed' WHERE id = ?`, expeditionID,
		)
		if err != nil {
			return fmt.Errorf("complete expedition (tx) %d: %w", expeditionID, err)
		}
		return nil
	})
}

// CreateVillageWithSetup atomically creates a village, links the map tile, creates
// starter buildings, resources, and (if needed) the player economy row.
func (u *unitOfWork) CreateVillageWithSetup(
	ctx context.Context,
	village *model.Village,
	tileX, tileY int,
	buildings []*model.Building,
	resources *model.Resources,
	playerID int64,
	startingGold float64,
) error {
	return WithTx(ctx, u.db, func(tx *sql.Tx) error {
		// 1. Insert village
		result, err := tx.ExecContext(ctx,
			`INSERT INTO villages (player_id, name, x, y, is_capital, season_id, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?)`,
			village.PlayerID, village.Name, village.X, village.Y,
			village.IsCapital, village.SeasonID, village.CreatedAt,
		)
		if err != nil {
			return fmt.Errorf("insert village (tx): %w", err)
		}
		vid, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("get village last insert id (tx): %w", err)
		}
		village.ID = vid

		// 2. Link map tile to village
		_, err = tx.ExecContext(ctx,
			`UPDATE world_map SET owner_player_id = ?, village_id = ? WHERE x = ? AND y = ?`,
			playerID, village.ID, tileX, tileY,
		)
		if err != nil {
			return fmt.Errorf("link tile to village (tx): %w", err)
		}

		// 3. Batch-insert starter buildings
		stmt, err := tx.PrepareContext(ctx,
			`INSERT INTO buildings (village_id, building_type, level) VALUES (?, ?, ?)`,
		)
		if err != nil {
			return fmt.Errorf("prepare insert building (tx): %w", err)
		}
		defer stmt.Close()

		for _, b := range buildings {
			b.VillageID = village.ID
			res, err := stmt.ExecContext(ctx, b.VillageID, b.BuildingType, b.Level)
			if err != nil {
				return fmt.Errorf("insert building %s (tx): %w", b.BuildingType, err)
			}
			bid, err := res.LastInsertId()
			if err != nil {
				return fmt.Errorf("get building last insert id (tx): %w", err)
			}
			b.ID = bid
		}

		// 4. Insert starter resources
		resources.VillageID = village.ID
		_, err = tx.ExecContext(ctx,
			`INSERT INTO resources (village_id, food, water, lumber, stone, food_rate, water_rate, lumber_rate, stone_rate, food_consumption, pop_used, max_food, max_water, max_lumber, max_stone, last_updated)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
			resources.VillageID, resources.Food, resources.Water, resources.Lumber, resources.Stone,
			resources.FoodRate, resources.WaterRate, resources.LumberRate, resources.StoneRate,
			resources.FoodConsumption, resources.PopUsed, resources.MaxFood, resources.MaxWater, resources.MaxLumber, resources.MaxStone,
			resources.LastUpdated.UTC().Format("2006-01-02 15:04:05"),
		)
		if err != nil {
			return fmt.Errorf("insert resources (tx): %w", err)
		}

		// 5. Create player economy if not exists (idempotent)
		if startingGold > 0 {
			_, err = tx.ExecContext(ctx,
				`INSERT OR IGNORE INTO player_economy (player_id, gold) VALUES (?, ?)`,
				playerID, startingGold,
			)
			if err != nil {
				return fmt.Errorf("create player economy (tx): %w", err)
			}
		}

		return nil
	})
}

// ── Internal transactional helpers (unexported, only used within this package) ─

func updateResourcesTx(ctx context.Context, tx *sql.Tx, villageID int64, res *model.Resources) error {
	_, err := tx.ExecContext(ctx,
		`UPDATE resources SET food = ?, water = ?, lumber = ?, stone = ?,
		 food_rate = ?, water_rate = ?, lumber_rate = ?, stone_rate = ?,
		 food_consumption = ?, pop_used = ?, max_food = ?, max_water = ?, max_lumber = ?, max_stone = ?, last_updated = ?
		 WHERE village_id = ?`,
		res.Food, res.Water, res.Lumber, res.Stone,
		res.FoodRate, res.WaterRate, res.LumberRate, res.StoneRate,
		res.FoodConsumption, res.PopUsed, res.MaxFood, res.MaxWater, res.MaxLumber, res.MaxStone,
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
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get build queue last insert id: %w", err)
	}
	item.ID = id
	return nil
}

func insertTrainQueueTx(ctx context.Context, tx *sql.Tx, item *model.TrainingQueue) error {
	result, err := tx.ExecContext(ctx,
		`INSERT INTO training_queue (village_id, troop_type, quantity, original_quantity, each_duration_sec, started_at, completes_at)
		 VALUES (?, ?, ?, ?, ?, ?, ?)`,
		item.VillageID, item.TroopType, item.Quantity, item.OriginalQuantity, item.EachDurationSec,
		item.StartedAt.UTC().Format("2006-01-02 15:04:05"),
		item.CompletesAt.UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("insert training queue item (tx): %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get train queue last insert id: %w", err)
	}
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

func deductGoldTx(ctx context.Context, tx *sql.Tx, playerID int64, amount float64) error {
	result, err := tx.ExecContext(ctx,
		`UPDATE player_economy SET gold = gold - ? WHERE player_id = ? AND gold >= ?`,
		amount, playerID, amount,
	)
	if err != nil {
		return fmt.Errorf("deduct gold (tx) for player %d: %w", playerID, err)
	}
	rows, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("check gold deduction (tx) for player %d: %w", playerID, err)
	}
	if rows == 0 {
		return model.ErrInsufficientGold
	}
	return nil
}

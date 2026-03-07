package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type resourceRepo struct {
	db *sql.DB
}

// NewResourceRepo creates a new SQLite-backed ResourceRepository.
func NewResourceRepo(db *sql.DB) *resourceRepo {
	return &resourceRepo{db: db}
}

func (r *resourceRepo) Create(ctx context.Context, res *model.Resources) error {
	_, err := r.db.ExecContext(ctx,
		`INSERT INTO resources (village_id, food, water, lumber, stone, food_rate, water_rate, lumber_rate, stone_rate, food_consumption, max_storage, last_updated)
		 VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?, ?)`,
		res.VillageID, res.Food, res.Water, res.Lumber, res.Stone,
		res.FoodRate, res.WaterRate, res.LumberRate, res.StoneRate,
		res.FoodConsumption, res.MaxStorage, res.LastUpdated.UTC().Format("2006-01-02 15:04:05"),
	)
	if err != nil {
		return fmt.Errorf("insert resources for village %d: %w", res.VillageID, err)
	}
	return nil
}

func (r *resourceRepo) Get(ctx context.Context, villageID int64) (*model.Resources, error) {
	row := r.db.QueryRowContext(ctx,
		`SELECT village_id, food, water, lumber, stone, food_rate, water_rate, lumber_rate, stone_rate, food_consumption, max_storage, last_updated
		 FROM resources WHERE village_id = ?`, villageID,
	)
	var res model.Resources
	var lastUpdatedStr string
	err := row.Scan(
		&res.VillageID, &res.Food, &res.Water, &res.Lumber, &res.Stone,
		&res.FoodRate, &res.WaterRate, &res.LumberRate, &res.StoneRate,
		&res.FoodConsumption, &res.MaxStorage, &lastUpdatedStr,
	)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("get resources for village %d: %w", villageID, err)
	}
	res.LastUpdated, _ = parseTime(lastUpdatedStr)
	return &res, nil
}

func (r *resourceRepo) Update(ctx context.Context, villageID int64, res *model.Resources) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE resources SET food = ?, water = ?, lumber = ?, stone = ?, food_rate = ?, water_rate = ?, lumber_rate = ?, stone_rate = ?, food_consumption = ?, max_storage = ?, last_updated = ?
		 WHERE village_id = ?`,
		res.Food, res.Water, res.Lumber, res.Stone,
		res.FoodRate, res.WaterRate, res.LumberRate, res.StoneRate,
		res.FoodConsumption, res.MaxStorage, res.LastUpdated.UTC().Format("2006-01-02 15:04:05"),
		villageID,
	)
	if err != nil {
		return fmt.Errorf("update resources for village %d: %w", villageID, err)
	}
	return nil
}

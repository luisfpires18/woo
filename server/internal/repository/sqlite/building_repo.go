package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/luisfpires18/woo/internal/model"
)

type buildingRepo struct {
	db *sql.DB
}

// NewBuildingRepo creates a new SQLite-backed BuildingRepository.
func NewBuildingRepo(db *sql.DB) *buildingRepo {
	return &buildingRepo{db: db}
}

func (r *buildingRepo) Create(ctx context.Context, building *model.Building) error {
	result, err := r.db.ExecContext(ctx,
		`INSERT INTO buildings (village_id, building_type, level) VALUES (?, ?, ?)`,
		building.VillageID, building.BuildingType, building.Level,
	)
	if err != nil {
		return fmt.Errorf("insert building: %w", err)
	}
	id, err := result.LastInsertId()
	if err != nil {
		return fmt.Errorf("get last insert id: %w", err)
	}
	building.ID = id
	return nil
}

func (r *buildingRepo) CreateBatch(ctx context.Context, buildings []*model.Building) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	stmt, err := tx.PrepareContext(ctx,
		`INSERT INTO buildings (village_id, building_type, level) VALUES (?, ?, ?)`,
	)
	if err != nil {
		return fmt.Errorf("prepare insert building: %w", err)
	}
	defer stmt.Close()

	for _, b := range buildings {
		result, err := stmt.ExecContext(ctx, b.VillageID, b.BuildingType, b.Level)
		if err != nil {
			return fmt.Errorf("insert building %s: %w", b.BuildingType, err)
		}
		id, err := result.LastInsertId()
		if err != nil {
			return fmt.Errorf("get last insert id for %s: %w", b.BuildingType, err)
		}
		b.ID = id
	}

	return tx.Commit()
}

func (r *buildingRepo) GetByVillageID(ctx context.Context, villageID int64) ([]*model.Building, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT id, village_id, building_type, level
		 FROM buildings WHERE village_id = ? ORDER BY building_type ASC`, villageID,
	)
	if err != nil {
		return nil, fmt.Errorf("list buildings for village %d: %w", villageID, err)
	}
	defer rows.Close()

	var buildings []*model.Building
	for rows.Next() {
		var b model.Building
		if err := rows.Scan(&b.ID, &b.VillageID, &b.BuildingType, &b.Level); err != nil {
			return nil, fmt.Errorf("scan building row: %w", err)
		}
		buildings = append(buildings, &b)
	}
	return buildings, rows.Err()
}

func (r *buildingRepo) Update(ctx context.Context, building *model.Building) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE buildings SET level = ? WHERE id = ?`,
		building.Level, building.ID,
	)
	if err != nil {
		return fmt.Errorf("update building %d: %w", building.ID, err)
	}
	return nil
}

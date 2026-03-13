package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/luisfpires18/woo/internal/model"
)

type worldMapRepo struct {
	db *sql.DB
}

// NewWorldMapRepo creates a new SQLite-backed WorldMapRepository.
func NewWorldMapRepo(db *sql.DB) *worldMapRepo {
	return &worldMapRepo{db: db}
}

// InsertBatch inserts multiple map tiles using batched INSERT statements within a transaction.
func (r *worldMapRepo) InsertBatch(ctx context.Context, tiles []*model.MapTile) error {
	return WithTx(ctx, r.db, func(tx *sql.Tx) error {
		const batchSize = 500
		for i := 0; i < len(tiles); i += batchSize {
			end := i + batchSize
			if end > len(tiles) {
				end = len(tiles)
			}
			batch := tiles[i:end]

			var sb strings.Builder
			sb.WriteString("INSERT OR IGNORE INTO world_map (x, y, terrain_type, kingdom_zone) VALUES ")

			args := make([]any, 0, len(batch)*4)
			for j, t := range batch {
				if j > 0 {
					sb.WriteString(",")
				}
				sb.WriteString("(?,?,?,?)")
				args = append(args, t.X, t.Y, t.TerrainType, t.KingdomZone)
			}

			if _, err := tx.ExecContext(ctx, sb.String(), args...); err != nil {
				return fmt.Errorf("insert map tile batch at offset %d: %w", i, err)
			}
		}
		return nil
	})
}

// DeleteAll removes all tiles from the world map.
func (r *worldMapRepo) DeleteAll(ctx context.Context) error {
	_, err := r.db.ExecContext(ctx, `DELETE FROM world_map`)
	if err != nil {
		return fmt.Errorf("delete all map tiles: %w", err)
	}
	return nil
}

// GetChunk retrieves map tiles in a square region centered on (cx, cy) with the given radius.
// Returns a (2*radius+1)^2 grid of tiles. Includes village name and owner name via JOINs.
func (r *worldMapRepo) GetChunk(ctx context.Context, cx, cy, radius int) ([]*model.MapTile, error) {
	minX, maxX := cx-radius, cx+radius
	minY, maxY := cy-radius, cy+radius

	rows, err := r.db.QueryContext(ctx, `
		SELECT wm.x, wm.y, wm.terrain_type, wm.kingdom_zone,
		       wm.owner_player_id, wm.village_id, wm.camp_id,
		       COALESCE(v.name, ''), COALESCE(p.username, '')
		FROM world_map wm
		LEFT JOIN villages v ON wm.village_id = v.id
		LEFT JOIN players  p ON wm.owner_player_id = p.id
		WHERE wm.x BETWEEN ? AND ? AND wm.y BETWEEN ? AND ?
		ORDER BY wm.y, wm.x`,
		minX, maxX, minY, maxY,
	)
	if err != nil {
		return nil, fmt.Errorf("query map chunk (%d,%d r=%d): %w", cx, cy, radius, err)
	}
	defer rows.Close()

	var tiles []*model.MapTile
	for rows.Next() {
		t := &model.MapTile{}
		var ownerID, villageID, campID sql.NullInt64
		if err := rows.Scan(
			&t.X, &t.Y, &t.TerrainType, &t.KingdomZone,
			&ownerID, &villageID, &campID,
			&t.VillageName, &t.OwnerName,
		); err != nil {
			return nil, fmt.Errorf("scan map tile: %w", err)
		}
		if ownerID.Valid {
			t.OwnerPlayerID = &ownerID.Int64
		}
		if villageID.Valid {
			t.VillageID = &villageID.Int64
		}
		if campID.Valid {
			t.CampID = &campID.Int64
		}
		tiles = append(tiles, t)
	}
	return tiles, rows.Err()
}

// GetTile returns a single tile at the given coordinates.
func (r *worldMapRepo) GetTile(ctx context.Context, x, y int) (*model.MapTile, error) {
	row := r.db.QueryRowContext(ctx, `
		SELECT x, y, terrain_type, kingdom_zone, owner_player_id, village_id, camp_id
		FROM world_map WHERE x = ? AND y = ?`, x, y,
	)

	t := &model.MapTile{}
	var ownerID, villageID, campID sql.NullInt64
	err := row.Scan(&t.X, &t.Y, &t.TerrainType, &t.KingdomZone, &ownerID, &villageID, &campID)
	if err == sql.ErrNoRows {
		return nil, model.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("scan map tile (%d,%d): %w", x, y, err)
	}
	if ownerID.Valid {
		t.OwnerPlayerID = &ownerID.Int64
	}
	if villageID.Valid {
		t.VillageID = &villageID.Int64
	}
	if campID.Valid {
		t.CampID = &campID.Int64
	}
	return t, nil
}

// UpdateTileOwner updates the owner and village association for a tile.
func (r *worldMapRepo) UpdateTileOwner(ctx context.Context, x, y int, playerID, villageID *int64) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE world_map SET owner_player_id = ?, village_id = ? WHERE x = ? AND y = ?`,
		playerID, villageID, x, y,
	)
	if err != nil {
		return fmt.Errorf("update tile owner (%d,%d): %w", x, y, err)
	}
	return nil
}

// Count returns the total number of tiles in the world_map table.
func (r *worldMapRepo) Count(ctx context.Context) (int64, error) {
	var count int64
	err := r.db.QueryRowContext(ctx, `SELECT COUNT(*) FROM world_map`).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("count world map tiles: %w", err)
	}
	return count, nil
}

// GetByZone returns all tiles belonging to a specific kingdom zone.
func (r *worldMapRepo) GetByZone(ctx context.Context, zone string) ([]*model.MapTile, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT x, y, terrain_type, kingdom_zone, owner_player_id, village_id
		 FROM world_map WHERE kingdom_zone = ?`, zone,
	)
	if err != nil {
		return nil, fmt.Errorf("query map zone %s: %w", zone, err)
	}
	defer rows.Close()

	var tiles []*model.MapTile
	for rows.Next() {
		t := &model.MapTile{}
		var ownerID, villageID sql.NullInt64
		if err := rows.Scan(&t.X, &t.Y, &t.TerrainType, &t.KingdomZone, &ownerID, &villageID); err != nil {
			return nil, fmt.Errorf("scan zone tile: %w", err)
		}
		if ownerID.Valid {
			t.OwnerPlayerID = &ownerID.Int64
		}
		if villageID.Valid {
			t.VillageID = &villageID.Int64
		}
		tiles = append(tiles, t)
	}
	return tiles, rows.Err()
}

// GetDistinctZones returns all distinct non-empty, non-wilderness kingdom zones currently placed.
func (r *worldMapRepo) GetDistinctZones(ctx context.Context) ([]string, error) {
	rows, err := r.db.QueryContext(ctx,
		`SELECT DISTINCT kingdom_zone FROM world_map WHERE kingdom_zone NOT IN ('', 'wilderness')`,
	)
	if err != nil {
		return nil, fmt.Errorf("query distinct zones: %w", err)
	}
	defer rows.Close()

	var zones []string
	for rows.Next() {
		var z string
		if err := rows.Scan(&z); err != nil {
			return nil, fmt.Errorf("scan zone: %w", err)
		}
		zones = append(zones, z)
	}
	return zones, rows.Err()
}

// UpdateTilesZone sets the kingdom_zone for all tiles within a circular radius of (cx, cy).
func (r *worldMapRepo) UpdateTilesZone(ctx context.Context, cx, cy, radius int, zone string) error {
	_, err := r.db.ExecContext(ctx,
		`UPDATE world_map SET kingdom_zone = ?
		 WHERE (x - ?) * (x - ?) + (y - ?) * (y - ?) <= ? * ?`,
		zone, cx, cx, cy, cy, radius, radius,
	)
	if err != nil {
		return fmt.Errorf("update tiles zone (%d,%d r=%d): %w", cx, cy, radius, err)
	}
	return nil
}

// UpdateTerrain updates the terrain_type for a batch of tiles in a single transaction.
func (r *worldMapRepo) UpdateTerrain(ctx context.Context, tiles []model.TileTerrainUpdate) error {
	return WithTx(ctx, r.db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx,
			`UPDATE world_map SET terrain_type = ? WHERE x = ? AND y = ?`)
		if err != nil {
			return fmt.Errorf("prepare terrain update: %w", err)
		}
		defer stmt.Close()

		for _, t := range tiles {
			if _, err := stmt.ExecContext(ctx, t.TerrainType, t.X, t.Y); err != nil {
				return fmt.Errorf("update terrain (%d,%d): %w", t.X, t.Y, err)
			}
		}
		return nil
	})
}

// GetSpawnCandidates returns plains tiles with no village, optionally filtering by zone.
func (r *worldMapRepo) GetSpawnCandidates(ctx context.Context, zone string) ([]*model.MapTile, error) {
	var query string
	var args []any
	if zone != "" {
		query = `SELECT x, y, terrain_type, kingdom_zone, owner_player_id, village_id
		         FROM world_map WHERE terrain_type = 'plains' AND village_id IS NULL AND kingdom_zone = ?`
		args = []any{zone}
	} else {
		query = `SELECT x, y, terrain_type, kingdom_zone, owner_player_id, village_id
		         FROM world_map WHERE terrain_type = 'plains' AND village_id IS NULL`
	}

	rows, err := r.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query spawn candidates: %w", err)
	}
	defer rows.Close()

	var tiles []*model.MapTile
	for rows.Next() {
		t := &model.MapTile{}
		var ownerID, villageID sql.NullInt64
		if err := rows.Scan(&t.X, &t.Y, &t.TerrainType, &t.KingdomZone, &ownerID, &villageID); err != nil {
			return nil, fmt.Errorf("scan spawn candidate: %w", err)
		}
		if ownerID.Valid {
			t.OwnerPlayerID = &ownerID.Int64
		}
		if villageID.Valid {
			t.VillageID = &villageID.Int64
		}
		tiles = append(tiles, t)
	}
	return tiles, rows.Err()
}

// UpdateTilesZoneBatch updates the kingdom_zone for each tile individually in a single transaction.
func (r *worldMapRepo) UpdateTilesZoneBatch(ctx context.Context, tiles []model.TemplateTile) error {
	return WithTx(ctx, r.db, func(tx *sql.Tx) error {
		stmt, err := tx.PrepareContext(ctx,
			`UPDATE world_map SET kingdom_zone = ? WHERE x = ? AND y = ?`)
		if err != nil {
			return fmt.Errorf("prepare zone batch update: %w", err)
		}
		defer stmt.Close()

		for _, t := range tiles {
			if _, err := stmt.ExecContext(ctx, t.KingdomZone, t.X, t.Y); err != nil {
				return fmt.Errorf("update zone (%d,%d): %w", t.X, t.Y, err)
			}
		}
		return nil
	})
}

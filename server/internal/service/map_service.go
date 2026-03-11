package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"math/rand"

	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Map generation errors.
var (
	ErrMapAlreadyGenerated = fmt.Errorf("world map already generated")
	ErrNoDirectionLeft     = fmt.Errorf("all kingdom placement directions are taken")
)

// MapService handles world map generation, queries, and dynamic kingdom zone placement.
type MapService struct {
	mapRepo     repository.WorldMapRepository
	villageRepo repository.VillageRepository
}

// NewMapService creates a new MapService.
func NewMapService(mapRepo repository.WorldMapRepository, villageRepo repository.VillageRepository) *MapService {
	return &MapService{
		mapRepo:     mapRepo,
		villageRepo: villageRepo,
	}
}

// directionSlot defines a predefined placement position for a kingdom zone.
type directionSlot struct {
	CenterX int
	CenterY int
}

// The 5 placement directions on a 51×51 grid (±25).
// Radius 10 each — no overlaps at these distances.
var directionSlots = []directionSlot{
	{CenterX: 0, CenterY: 15},  // North
	{CenterX: 0, CenterY: -15}, // South
	{CenterX: 15, CenterY: 0},  // East
	{CenterX: -15, CenterY: 0}, // West
	{CenterX: 0, CenterY: 0},   // Center
}

// KingdomZoneRadius is the radius of each kingdom zone in tiles.
const KingdomZoneRadius = 10

// GenerateMap generates a procedural world map using simplex noise.
// Idempotent — skips if already generated.
// All tiles default to plains. Terrain is painted manually via the admin map editor.
// Kingdom zones are placed dynamically when the first player joins.
func (s *MapService) GenerateMap(ctx context.Context) error {
	count, err := s.mapRepo.Count(ctx)
	if err != nil {
		return fmt.Errorf("check map count: %w", err)
	}
	expectedTiles := int64(model.MapSize * model.MapSize)
	if count == expectedTiles {
		slog.Info("world map already generated", "tiles", count)
		return nil
	}

	if count > 0 {
		slog.Info("backfilling partial world map", "existing_tiles", count, "expected_tiles", expectedTiles)
	} else {
		slog.Info("generating world map", "size", model.MapSize)
	}

	tiles := make([]*model.MapTile, 0, model.MapSize*model.MapSize)
	for y := -model.MapHalf; y <= model.MapHalf; y++ {
		for x := -model.MapHalf; x <= model.MapHalf; x++ {
			tiles = append(tiles, &model.MapTile{
				X:           x,
				Y:           y,
				TerrainType: model.TerrainPlains,
				KingdomZone: model.ZoneWilderness,
			})
		}
	}

	if err := s.mapRepo.InsertBatch(ctx, tiles); err != nil {
		return fmt.Errorf("insert map tiles: %w", err)
	}

	finalCount, err := s.mapRepo.Count(ctx)
	if err != nil {
		return fmt.Errorf("count generated map tiles: %w", err)
	}

	slog.Info("world map ready", "tiles", finalCount)
	return nil
}

// PlaceKingdomZone assigns a random available direction to a kingdom and updates all tiles
// within the zone radius. Returns the zone center coordinates.
func (s *MapService) PlaceKingdomZone(ctx context.Context, kingdom string) (int, int, error) {
	// Check which zones are already placed
	placedZones, err := s.mapRepo.GetDistinctZones(ctx)
	if err != nil {
		return 0, 0, fmt.Errorf("get placed zones: %w", err)
	}

	// Figure out which direction slots are already taken
	takenSlots := make(map[int]bool) // index into directionSlots
	for _, zone := range placedZones {
		// Find the slot this zone occupies by checking existing tiles
		zoneTiles, err := s.mapRepo.GetByZone(ctx, zone)
		if err != nil || len(zoneTiles) == 0 {
			continue
		}
		// Compute center of mass to identify which slot
		var sumX, sumY int
		for _, t := range zoneTiles {
			sumX += t.X
			sumY += t.Y
		}
		avgX := sumX / len(zoneTiles)
		avgY := sumY / len(zoneTiles)

		// Match to closest slot
		bestIdx := -1
		bestDist := 9999
		for i, slot := range directionSlots {
			dx := avgX - slot.CenterX
			dy := avgY - slot.CenterY
			dist := dx*dx + dy*dy
			if dist < bestDist {
				bestDist = dist
				bestIdx = i
			}
		}
		if bestIdx >= 0 {
			takenSlots[bestIdx] = true
		}
	}

	// Collect available slots
	var available []int
	for i := range directionSlots {
		if !takenSlots[i] {
			available = append(available, i)
		}
	}

	if len(available) == 0 {
		return 0, 0, ErrNoDirectionLeft
	}

	// Pick a random available slot
	chosen := available[rand.Intn(len(available))]
	slot := directionSlots[chosen]

	// Update all tiles within the radius to this kingdom's zone
	if err := s.mapRepo.UpdateTilesZone(ctx, slot.CenterX, slot.CenterY, KingdomZoneRadius, kingdom); err != nil {
		return 0, 0, fmt.Errorf("place kingdom zone %s: %w", kingdom, err)
	}

	slog.Info("kingdom zone placed",
		"kingdom", kingdom,
		"center_x", slot.CenterX,
		"center_y", slot.CenterY,
		"radius", KingdomZoneRadius,
		"slot", chosen,
	)

	return slot.CenterX, slot.CenterY, nil
}

// GetMapChunk returns a chunk of map tiles centered on (cx, cy) with the given radius.
func (s *MapService) GetMapChunk(ctx context.Context, cx, cy, radius int) ([]*model.MapTile, error) {
	if radius < 1 {
		radius = 10
	}
	if radius > 40 {
		radius = 40
	}

	tiles, err := s.mapRepo.GetChunk(ctx, cx, cy, radius)
	if err != nil {
		return nil, fmt.Errorf("get map chunk: %w", err)
	}
	return tiles, nil
}

// GetTile returns a single map tile.
func (s *MapService) GetTile(ctx context.Context, x, y int) (*model.MapTile, error) {
	return s.mapRepo.GetTile(ctx, x, y)
}

// FindSpawnTile finds a suitable tile for a new village spawn.
// Prefers tiles within the player's kingdom zone. Falls back to any plains tile
// if no zone-specific tiles exist (e.g. zones not yet painted on the template).
func (s *MapService) FindSpawnTile(ctx context.Context, kingdom string) (int, int, error) {
	x, y, found, err := s.findAvailableSpawnInScope(ctx, kingdom)
	if err != nil {
		return 0, 0, err
	}
	if found {
		return x, y, nil
	}

	slog.Warn("no zone tiles for kingdom, falling back to any plains tile", "kingdom", kingdom)

	x, y, found, err = s.findAvailableSpawnInScope(ctx, "")
	if err != nil {
		return 0, 0, err
	}
	if found {
		return x, y, nil
	}

	return 0, 0, ErrNoSpawnTile
}

func (s *MapService) findAvailableSpawnInScope(ctx context.Context, zone string) (int, int, bool, error) {
	candidates, err := s.mapRepo.GetSpawnCandidates(ctx, zone)
	if err != nil {
		if zone != "" {
			return 0, 0, false, fmt.Errorf("get zone spawn candidates for %s: %w", zone, err)
		}
		return 0, 0, false, fmt.Errorf("get fallback spawn candidates: %w", err)
	}

	if len(candidates) == 0 {
		return 0, 0, false, nil
	}

	return s.selectAvailableSpawnCandidate(ctx, candidates)
}

func (s *MapService) selectAvailableSpawnCandidate(ctx context.Context, candidates []*model.MapTile) (int, int, bool, error) {
	rand.Shuffle(len(candidates), func(i, j int) {
		candidates[i], candidates[j] = candidates[j], candidates[i]
	})

	for _, candidate := range candidates {
		village, err := s.villageRepo.GetByCoordinates(ctx, candidate.X, candidate.Y)
		switch {
		case errors.Is(err, model.ErrNotFound):
			return candidate.X, candidate.Y, true, nil
		case err != nil:
			return 0, 0, false, fmt.Errorf("check spawn candidate (%d,%d): %w", candidate.X, candidate.Y, err)
		default:
			if syncErr := s.mapRepo.UpdateTileOwner(ctx, candidate.X, candidate.Y, &village.PlayerID, &village.ID); syncErr != nil {
				slog.Warn("failed to sync stale spawn tile owner",
					"x", candidate.X,
					"y", candidate.Y,
					"village_id", village.ID,
					"error", syncErr,
				)
			}
		}
	}

	return 0, 0, false, nil
}

// UpdateTileOwner links a map tile to a village and player.
func (s *MapService) UpdateTileOwner(ctx context.Context, x, y int, playerID, villageID int64) error {
	return s.mapRepo.UpdateTileOwner(ctx, x, y, &playerID, &villageID)
}

// ValidTerrainTypes is the set of terrain types allowed for painting.
var ValidTerrainTypes = map[string]bool{
	model.TerrainPlains:   true,
	model.TerrainForest:   true,
	model.TerrainMountain: true,
	model.TerrainWater:    true,
	model.TerrainDesert:   true,
	model.TerrainSwamp:    true,
}

// UpdateTerrain validates and applies terrain changes for the admin paint tool.
func (s *MapService) UpdateTerrain(ctx context.Context, tiles []model.TileTerrainUpdate) error {
	if len(tiles) == 0 {
		return fmt.Errorf("no tiles provided")
	}
	if len(tiles) > 500 {
		return fmt.Errorf("too many tiles in one request (max 500)")
	}
	for _, t := range tiles {
		if !ValidTerrainTypes[t.TerrainType] {
			return fmt.Errorf("invalid terrain type: %s", t.TerrainType)
		}
		if t.X < -model.MapHalf || t.X > model.MapHalf || t.Y < -model.MapHalf || t.Y > model.MapHalf {
			return fmt.Errorf("tile (%d,%d) out of map bounds", t.X, t.Y)
		}
	}
	return s.mapRepo.UpdateTerrain(ctx, tiles)
}

package service

import (
	"context"
	"testing"
	"time"

	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/testutil"
)

// newTestMapService creates a MapService backed by an in-memory DB with migrations applied.
func newTestMapService(t *testing.T) *MapService {
	t.Helper()
	db := testutil.NewTestDB(t)
	if _, err := db.Exec(`
		DELETE FROM building_queue;
		DELETE FROM training_queue;
		DELETE FROM troops;
		DELETE FROM resources;
		DELETE FROM buildings;
		DELETE FROM villages;
		DELETE FROM season_players;
		DELETE FROM world_map;
	`); err != nil {
		t.Fatalf("reset seeded map test data: %v", err)
	}
	worldMapRepo := sqlite.NewWorldMapRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	return NewMapService(worldMapRepo, villageRepo)
}

// --- GenerateMap tests ---

func TestGenerateMap_CreatesExpectedTileCount(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	err := svc.GenerateMap(ctx)
	if err != nil {
		t.Fatalf("GenerateMap failed: %v", err)
	}

	// 51×51 = 2601 tiles
	tiles, err := svc.GetMapChunk(ctx, 0, 0, 25)
	if err != nil {
		t.Fatalf("GetMapChunk failed: %v", err)
	}

	expected := 51 * 51
	if len(tiles) != expected {
		t.Errorf("tile count: got %d, want %d", len(tiles), expected)
	}
}

func TestGenerateMap_AllTilesArePlains(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	tiles, err := svc.GetMapChunk(ctx, 0, 0, 25)
	if err != nil {
		t.Fatalf("GetMapChunk: %v", err)
	}

	for _, tile := range tiles {
		if tile.TerrainType != model.TerrainPlains {
			t.Errorf("tile (%d,%d): got terrain %q, want %q", tile.X, tile.Y, tile.TerrainType, model.TerrainPlains)
		}
	}
}

func TestGenerateMap_AllTilesStartAsWilderness(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	tiles, err := svc.GetMapChunk(ctx, 0, 0, 25)
	if err != nil {
		t.Fatalf("GetMapChunk: %v", err)
	}

	for _, tile := range tiles {
		if tile.KingdomZone != model.ZoneWilderness {
			t.Errorf("tile (%d,%d): got zone %q, want %q", tile.X, tile.Y, tile.KingdomZone, model.ZoneWilderness)
		}
	}
}

func TestGenerateMap_Idempotent(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("first GenerateMap: %v", err)
	}
	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("second GenerateMap should be idempotent, got: %v", err)
	}
}

func TestGenerateMap_BackfillsPartialSeededMap(t *testing.T) {
	ctx := context.Background()
	db := testutil.NewTestDB(t)
	worldMapRepo := sqlite.NewWorldMapRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	svc := NewMapService(worldMapRepo, villageRepo)

	countBefore, err := worldMapRepo.Count(ctx)
	if err != nil {
		t.Fatalf("count before: %v", err)
	}
	if countBefore >= int64(model.MapSize*model.MapSize) {
		t.Fatalf("expected partial seeded map, got %d tiles", countBefore)
	}

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	countAfter, err := worldMapRepo.Count(ctx)
	if err != nil {
		t.Fatalf("count after: %v", err)
	}
	if countAfter != int64(model.MapSize*model.MapSize) {
		t.Fatalf("expected %d tiles after backfill, got %d", model.MapSize*model.MapSize, countAfter)
	}

	tile, err := svc.GetTile(ctx, 0, 0)
	if err != nil {
		t.Fatalf("GetTile(0,0): %v", err)
	}
	if tile.VillageID == nil {
		t.Fatal("expected seeded village tile ownership to be preserved")
	}
}

// --- PlaceKingdomZone tests ---

func TestPlaceKingdomZone_AssignsDirection(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	cx, cy, err := svc.PlaceKingdomZone(ctx, model.ZoneVeridor)
	if err != nil {
		t.Fatalf("PlaceKingdomZone: %v", err)
	}

	// Center should be one of the 5 predefined slots
	validCenters := [][2]int{{0, 15}, {0, -15}, {15, 0}, {-15, 0}, {0, 0}}
	found := false
	for _, vc := range validCenters {
		if cx == vc[0] && cy == vc[1] {
			found = true
			break
		}
	}
	if !found {
		t.Errorf("unexpected zone center (%d,%d)", cx, cy)
	}

	// Tiles in the zone should now be veridor
	tiles, err := svc.GetMapChunk(ctx, cx, cy, 10)
	if err != nil {
		t.Fatalf("GetMapChunk: %v", err)
	}

	veridorCount := 0
	for _, tile := range tiles {
		if tile.KingdomZone == model.ZoneVeridor {
			veridorCount++
		}
	}

	if veridorCount == 0 {
		t.Error("no veridor tiles found after PlaceKingdomZone")
	}
}

func TestPlaceKingdomZone_SecondKingdomGetsDifferentSlot(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	cx1, cy1, err := svc.PlaceKingdomZone(ctx, model.ZoneVeridor)
	if err != nil {
		t.Fatalf("PlaceKingdomZone veridor: %v", err)
	}

	cx2, cy2, err := svc.PlaceKingdomZone(ctx, model.ZoneSylvara)
	if err != nil {
		t.Fatalf("PlaceKingdomZone sylvara: %v", err)
	}

	if cx1 == cx2 && cy1 == cy2 {
		t.Errorf("both kingdoms placed at same center (%d,%d)", cx1, cy1)
	}
}

func TestPlaceKingdomZone_AllFiveKingdoms(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	kingdoms := []string{
		model.ZoneVeridor, model.ZoneSylvara, model.ZoneArkazia,
		model.ZoneDraxys, model.ZoneNordalh,
	}

	centers := make(map[[2]int]string)
	for _, k := range kingdoms {
		cx, cy, err := svc.PlaceKingdomZone(ctx, k)
		if err != nil {
			t.Fatalf("PlaceKingdomZone(%s): %v", k, err)
		}
		center := [2]int{cx, cy}
		if prev, ok := centers[center]; ok {
			t.Errorf("kingdoms %q and %q both placed at (%d,%d)", prev, k, cx, cy)
		}
		centers[center] = k
	}

	// All 5 should be placed
	if len(centers) != 5 {
		t.Errorf("expected 5 unique centers, got %d", len(centers))
	}
}

func TestPlaceKingdomZone_SixthKingdomFails(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	// Place all 5 playable kingdoms
	kingdoms := []string{
		model.ZoneVeridor, model.ZoneSylvara, model.ZoneArkazia,
		model.ZoneDraxys, model.ZoneNordalh,
	}
	for _, k := range kingdoms {
		if _, _, err := svc.PlaceKingdomZone(ctx, k); err != nil {
			t.Fatalf("PlaceKingdomZone(%s): %v", k, err)
		}
	}

	// A 6th should fail (all 5 slots taken)
	_, _, err := svc.PlaceKingdomZone(ctx, "extra_kingdom")
	if err == nil {
		t.Error("expected error placing 6th kingdom, got nil")
	}
}

// --- FindSpawnTile tests ---

func TestFindSpawnTile_WithZone(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	// Pre-place the kingdom zone so FindSpawnTile can find it
	if _, _, err := svc.PlaceKingdomZone(ctx, model.ZoneArkazia); err != nil {
		t.Fatalf("PlaceKingdomZone(arkazia): %v", err)
	}

	x, y, err := svc.FindSpawnTile(ctx, model.ZoneArkazia)
	if err != nil {
		t.Fatalf("FindSpawnTile(arkazia): %v", err)
	}

	// Verify the spawn tile is in the Arkazia zone
	tile, err := svc.GetTile(ctx, x, y)
	if err != nil {
		t.Fatalf("GetTile(%d,%d): %v", x, y, err)
	}

	if tile.KingdomZone != model.ZoneArkazia {
		t.Errorf("spawn tile zone: got %q, want %q", tile.KingdomZone, model.ZoneArkazia)
	}
	if tile.TerrainType != model.TerrainPlains {
		t.Errorf("spawn tile terrain: got %q, want %q", tile.TerrainType, model.TerrainPlains)
	}
}

func TestFindSpawnTile_FallbackNoZone(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	// No zones placed — FindSpawnTile should fall back to any plains tile
	x, y, err := svc.FindSpawnTile(ctx, model.ZoneVeridor)
	if err != nil {
		t.Fatalf("FindSpawnTile(veridor) with no zones should fallback, got: %v", err)
	}

	tile, err := svc.GetTile(ctx, x, y)
	if err != nil {
		t.Fatalf("GetTile(%d,%d): %v", x, y, err)
	}
	if tile.TerrainType != model.TerrainPlains {
		t.Errorf("fallback spawn tile terrain: got %q, want %q", tile.TerrainType, model.TerrainPlains)
	}
}

func TestSelectAvailableSpawnCandidate_SkipsStaleOccupiedTile(t *testing.T) {
	ctx := context.Background()
	db := testutil.NewTestDB(t)
	worldMapRepo := sqlite.NewWorldMapRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	playerRepo := sqlite.NewPlayerRepo(db)
	svc := NewMapService(worldMapRepo, villageRepo)

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	player := &model.Player{
		Username:     "spawnsync",
		Email:        "spawnsync@example.com",
		PasswordHash: "not-a-real-hash",
		Kingdom:      model.ZoneVeridor,
		CreatedAt:    time.Now().UTC(),
	}
	if err := playerRepo.Create(ctx, player); err != nil {
		t.Fatalf("create player: %v", err)
	}

	occupied := &model.Village{
		PlayerID:  player.ID,
		Name:      "Occupied",
		X:         10,
		Y:         10,
		IsCapital: true,
		CreatedAt: time.Now().UTC(),
	}
	if err := villageRepo.Create(ctx, occupied); err != nil {
		t.Fatalf("create occupied village: %v", err)
	}

	x, y, found, err := svc.selectAvailableSpawnCandidate(ctx, []*model.MapTile{
		{X: 10, Y: 10, TerrainType: model.TerrainPlains},
		{X: 11, Y: 10, TerrainType: model.TerrainPlains},
	})
	if err != nil {
		t.Fatalf("selectAvailableSpawnCandidate: %v", err)
	}
	if !found {
		t.Fatal("expected to find a free spawn candidate")
	}
	if x != 11 || y != 10 {
		t.Fatalf("expected free tile (11,10), got (%d,%d)", x, y)
	}

	tile, err := svc.GetTile(ctx, 10, 10)
	if err != nil {
		t.Fatalf("GetTile(10,10): %v", err)
	}
	if tile.VillageID == nil || *tile.VillageID != occupied.ID {
		t.Fatalf("expected stale tile to sync village_id %d, got %+v", occupied.ID, tile.VillageID)
	}
	if tile.OwnerPlayerID == nil || *tile.OwnerPlayerID != player.ID {
		t.Fatalf("expected stale tile to sync owner_player_id %d, got %+v", player.ID, tile.OwnerPlayerID)
	}
}

func TestFindSpawnTile_MultipleKingdoms(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	kingdoms := []string{
		model.ZoneVeridor, model.ZoneSylvara, model.ZoneArkazia,
		model.ZoneDraxys, model.ZoneNordalh,
	}

	// Pre-place all kingdom zones
	for _, k := range kingdoms {
		if _, _, err := svc.PlaceKingdomZone(ctx, k); err != nil {
			t.Fatalf("PlaceKingdomZone(%s): %v", k, err)
		}
	}

	for _, k := range kingdoms {
		t.Run(k, func(t *testing.T) {
			x, y, err := svc.FindSpawnTile(ctx, k)
			if err != nil {
				t.Fatalf("FindSpawnTile(%s): %v", k, err)
			}
			if x < -model.MapHalf || x > model.MapHalf || y < -model.MapHalf || y > model.MapHalf {
				t.Errorf("spawn (%d,%d) outside map bounds", x, y)
			}
		})
	}
}

// --- GetMapChunk tests ---

func TestGetMapChunk_RadiusClamping(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	tests := []struct {
		name     string
		radius   int
		wantSize int
	}{
		{"normal radius 5", 5, 11 * 11},
		{"radius 0 clamps to 10", 0, 21 * 21},
		// radius 40 => side 81, but map is only 51, so repo returns all within range
		// At center (0,0) only -25..+25 exist: min(40,25) means we get 51*51
		{"large radius clamps to 40", 100, 51 * 51},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tiles, err := svc.GetMapChunk(ctx, 0, 0, tt.radius)
			if err != nil {
				t.Fatalf("GetMapChunk(%d): %v", tt.radius, err)
			}
			if len(tiles) != tt.wantSize {
				t.Errorf("tiles for radius=%d: got %d, want %d", tt.radius, len(tiles), tt.wantSize)
			}
		})
	}
}

func TestGetTile_ReturnsOriginTile(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	tile, err := svc.GetTile(ctx, 0, 0)
	if err != nil {
		t.Fatalf("GetTile(0,0): %v", err)
	}

	if tile.X != 0 || tile.Y != 0 {
		t.Errorf("expected (0,0), got (%d,%d)", tile.X, tile.Y)
	}
	if tile.TerrainType != model.TerrainPlains {
		t.Errorf("origin terrain: got %q, want %q", tile.TerrainType, model.TerrainPlains)
	}
	if tile.KingdomZone != model.ZoneWilderness {
		t.Errorf("origin zone: got %q, want %q", tile.KingdomZone, model.ZoneWilderness)
	}
}

func TestGetTile_NotFound(t *testing.T) {
	svc := newTestMapService(t)
	ctx := context.Background()

	if err := svc.GenerateMap(ctx); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	// Outside the map bounds
	_, err := svc.GetTile(ctx, 999, 999)
	if err == nil {
		t.Error("expected error for tile outside map bounds")
	}
}

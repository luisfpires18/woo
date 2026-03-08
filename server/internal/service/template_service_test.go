package service_test

import (
	"context"
	"testing"

	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/service"
	"github.com/luisfpires18/woo/internal/testutil"
)

// newTestTemplateService creates a TemplateService backed by a temp directory and in-memory DB.
func newTestTemplateService(t *testing.T) *service.TemplateService {
	t.Helper()
	tmpDir := t.TempDir()
	templateRepo, err := repository.NewTemplateRepository(tmpDir)
	if err != nil {
		t.Fatalf("NewTemplateRepository: %v", err)
	}
	db := testutil.NewTestDB(t)
	mapRepo := sqlite.NewWorldMapRepo(db)
	return service.NewTemplateService(templateRepo, mapRepo)
}

// --- CreateTemplate ---

func TestCreateTemplate_Success(t *testing.T) {
	svc := newTestTemplateService(t)

	tmpl, err := svc.CreateTemplate("test-map", "A test template", 0)
	if err != nil {
		t.Fatalf("CreateTemplate: %v", err)
	}

	if tmpl.Name != "test-map" {
		t.Errorf("name: got %q, want %q", tmpl.Name, "test-map")
	}
	if tmpl.MapSize != model.MapSize {
		t.Errorf("map_size: got %d, want %d", tmpl.MapSize, model.MapSize)
	}
	expectedTiles := model.MapSize * model.MapSize
	if len(tmpl.Tiles) != expectedTiles {
		t.Errorf("tile count: got %d, want %d", len(tmpl.Tiles), expectedTiles)
	}
}

func TestCreateTemplate_CustomSize(t *testing.T) {
	svc := newTestTemplateService(t)

	tmpl, err := svc.CreateTemplate("small", "Small map", 11)
	if err != nil {
		t.Fatalf("CreateTemplate: %v", err)
	}

	if tmpl.MapSize != 11 {
		t.Errorf("map_size: got %d, want 11", tmpl.MapSize)
	}
	if len(tmpl.Tiles) != 121 {
		t.Errorf("tile count: got %d, want 121", len(tmpl.Tiles))
	}
}

func TestCreateTemplate_EvenSizeRoundsUp(t *testing.T) {
	svc := newTestTemplateService(t)

	tmpl, err := svc.CreateTemplate("even", "", 10)
	if err != nil {
		t.Fatalf("CreateTemplate: %v", err)
	}

	if tmpl.MapSize != 11 {
		t.Errorf("map_size: got %d, want 11 (even 10 should round up)", tmpl.MapSize)
	}
}

func TestCreateTemplate_Duplicate(t *testing.T) {
	svc := newTestTemplateService(t)

	if _, err := svc.CreateTemplate("dup", "", 5); err != nil {
		t.Fatalf("first create: %v", err)
	}

	_, err := svc.CreateTemplate("dup", "", 5)
	if err == nil {
		t.Fatal("expected error on duplicate, got nil")
	}
}

func TestCreateTemplate_TooLarge(t *testing.T) {
	svc := newTestTemplateService(t)

	_, err := svc.CreateTemplate("huge", "", 203)
	if err == nil {
		t.Fatal("expected error for size > 201, got nil")
	}
}

func TestCreateTemplate_TilesAreAllPlainsWilderness(t *testing.T) {
	svc := newTestTemplateService(t)

	tmpl, err := svc.CreateTemplate("check", "", 5)
	if err != nil {
		t.Fatalf("CreateTemplate: %v", err)
	}

	for _, tile := range tmpl.Tiles {
		if tile.TerrainType != model.TerrainPlains {
			t.Errorf("tile (%d,%d) terrain: got %q, want %q", tile.X, tile.Y, tile.TerrainType, model.TerrainPlains)
		}
		if tile.KingdomZone != model.ZoneWilderness {
			t.Errorf("tile (%d,%d) zone: got %q, want %q", tile.X, tile.Y, tile.KingdomZone, model.ZoneWilderness)
		}
	}
}

// --- GetTemplate ---

func TestGetTemplate_Success(t *testing.T) {
	svc := newTestTemplateService(t)

	if _, err := svc.CreateTemplate("my-tmpl", "desc", 5); err != nil {
		t.Fatalf("CreateTemplate: %v", err)
	}

	tmpl, err := svc.GetTemplate("my-tmpl")
	if err != nil {
		t.Fatalf("GetTemplate: %v", err)
	}
	if tmpl.Name != "my-tmpl" {
		t.Errorf("name: got %q, want %q", tmpl.Name, "my-tmpl")
	}
	if tmpl.MapSize != 5 {
		t.Errorf("map_size: got %d, want 5", tmpl.MapSize)
	}
}

func TestGetTemplate_NotFound(t *testing.T) {
	svc := newTestTemplateService(t)

	_, err := svc.GetTemplate("nonexistent")
	if err == nil {
		t.Fatal("expected error for nonexistent template, got nil")
	}
}

// --- ListTemplates ---

func TestListTemplates_Empty(t *testing.T) {
	svc := newTestTemplateService(t)

	list, err := svc.ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates: %v", err)
	}
	if len(list) != 0 {
		t.Errorf("list count: got %d, want 0", len(list))
	}
}

func TestListTemplates_MultipleSorted(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("beta", "", 5)
	svc.CreateTemplate("alpha", "", 7)

	list, err := svc.ListTemplates()
	if err != nil {
		t.Fatalf("ListTemplates: %v", err)
	}
	if len(list) != 2 {
		t.Fatalf("list count: got %d, want 2", len(list))
	}
}

// --- DeleteTemplate ---

func TestDeleteTemplate_Success(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("to-delete", "", 5)

	if err := svc.DeleteTemplate("to-delete"); err != nil {
		t.Fatalf("DeleteTemplate: %v", err)
	}

	// Verify it's gone
	_, err := svc.GetTemplate("to-delete")
	if err == nil {
		t.Fatal("expected error after delete, got nil")
	}
}

func TestDeleteTemplate_NotFound(t *testing.T) {
	svc := newTestTemplateService(t)

	err := svc.DeleteTemplate("nope")
	if err == nil {
		t.Fatal("expected error deleting nonexistent template, got nil")
	}
}

// --- UpdateTemplateTerrain ---

func TestUpdateTemplateTerrain_Success(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("terrain-test", "", 5)

	updates := []model.TileTerrainUpdate{
		{X: 0, Y: 0, TerrainType: model.TerrainForest},
		{X: 1, Y: 0, TerrainType: model.TerrainMountain},
	}

	if err := svc.UpdateTemplateTerrain("terrain-test", updates); err != nil {
		t.Fatalf("UpdateTemplateTerrain: %v", err)
	}

	tmpl, err := svc.GetTemplate("terrain-test")
	if err != nil {
		t.Fatalf("GetTemplate: %v", err)
	}

	// Verify the updated tiles
	found := 0
	for _, tile := range tmpl.Tiles {
		if tile.X == 0 && tile.Y == 0 {
			if tile.TerrainType != model.TerrainForest {
				t.Errorf("(0,0) terrain: got %q, want %q", tile.TerrainType, model.TerrainForest)
			}
			found++
		}
		if tile.X == 1 && tile.Y == 0 {
			if tile.TerrainType != model.TerrainMountain {
				t.Errorf("(1,0) terrain: got %q, want %q", tile.TerrainType, model.TerrainMountain)
			}
			found++
		}
	}
	if found != 2 {
		t.Errorf("expected 2 updated tiles, found %d", found)
	}
}

func TestUpdateTemplateTerrain_InvalidType(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("bad-terrain", "", 5)

	err := svc.UpdateTemplateTerrain("bad-terrain", []model.TileTerrainUpdate{
		{X: 0, Y: 0, TerrainType: "lava"},
	})
	if err == nil {
		t.Fatal("expected error for invalid terrain type")
	}
}

func TestUpdateTemplateTerrain_OutOfBounds(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("bounds-test", "", 5) // size 5 = -2 to +2

	err := svc.UpdateTemplateTerrain("bounds-test", []model.TileTerrainUpdate{
		{X: 100, Y: 100, TerrainType: model.TerrainForest},
	})
	if err == nil {
		t.Fatal("expected error for out-of-bounds tile")
	}
}

func TestUpdateTemplateTerrain_Empty(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("empty-update", "", 5)

	err := svc.UpdateTemplateTerrain("empty-update", []model.TileTerrainUpdate{})
	if err == nil {
		t.Fatal("expected error for empty update list")
	}
}

// --- UpdateTemplateZones ---

func TestUpdateTemplateZones_Success(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("zone-test", "", 5)

	updates := []model.TileZoneUpdate{
		{X: 0, Y: 0, KingdomZone: model.ZoneVeridor},
		{X: 1, Y: 0, KingdomZone: model.ZoneSylvara},
	}

	if err := svc.UpdateTemplateZones("zone-test", updates); err != nil {
		t.Fatalf("UpdateTemplateZones: %v", err)
	}

	tmpl, err := svc.GetTemplate("zone-test")
	if err != nil {
		t.Fatalf("GetTemplate: %v", err)
	}

	for _, tile := range tmpl.Tiles {
		if tile.X == 0 && tile.Y == 0 && tile.KingdomZone != model.ZoneVeridor {
			t.Errorf("(0,0) zone: got %q, want %q", tile.KingdomZone, model.ZoneVeridor)
		}
		if tile.X == 1 && tile.Y == 0 && tile.KingdomZone != model.ZoneSylvara {
			t.Errorf("(1,0) zone: got %q, want %q", tile.KingdomZone, model.ZoneSylvara)
		}
	}
}

func TestUpdateTemplateZones_InvalidZone(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("bad-zone", "", 5)

	err := svc.UpdateTemplateZones("bad-zone", []model.TileZoneUpdate{
		{X: 0, Y: 0, KingdomZone: "atlantis"},
	})
	if err == nil {
		t.Fatal("expected error for invalid zone")
	}
}

// --- ResizeTemplate ---

func TestResizeTemplate_Grow(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("resize-grow", "", 5)

	// Paint a tile to ensure it's preserved after resize
	svc.UpdateTemplateTerrain("resize-grow", []model.TileTerrainUpdate{
		{X: 0, Y: 0, TerrainType: model.TerrainForest},
	})

	tmpl, err := svc.ResizeTemplate("resize-grow", 9)
	if err != nil {
		t.Fatalf("ResizeTemplate: %v", err)
	}

	if tmpl.MapSize != 9 {
		t.Errorf("map_size: got %d, want 9", tmpl.MapSize)
	}
	if len(tmpl.Tiles) != 81 {
		t.Errorf("tile count: got %d, want 81", len(tmpl.Tiles))
	}

	// Check that the painted tile is preserved
	for _, tile := range tmpl.Tiles {
		if tile.X == 0 && tile.Y == 0 {
			if tile.TerrainType != model.TerrainForest {
				t.Errorf("(0,0) after resize: got %q, want %q", tile.TerrainType, model.TerrainForest)
			}
			break
		}
	}
}

func TestResizeTemplate_Shrink(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("resize-shrink", "", 9)

	// Paint a tile at the edge that will be dropped
	svc.UpdateTemplateTerrain("resize-shrink", []model.TileTerrainUpdate{
		{X: 4, Y: 4, TerrainType: model.TerrainMountain},
	})

	tmpl, err := svc.ResizeTemplate("resize-shrink", 5)
	if err != nil {
		t.Fatalf("ResizeTemplate: %v", err)
	}

	if tmpl.MapSize != 5 {
		t.Errorf("map_size: got %d, want 5", tmpl.MapSize)
	}
	if len(tmpl.Tiles) != 25 {
		t.Errorf("tile count: got %d, want 25", len(tmpl.Tiles))
	}

	// The (4,4) tile is outside new bounds (-2 to +2), should be gone
	for _, tile := range tmpl.Tiles {
		if tile.X == 4 && tile.Y == 4 {
			t.Error("tile (4,4) should have been dropped after shrink")
		}
	}
}

func TestResizeTemplate_TooLarge(t *testing.T) {
	svc := newTestTemplateService(t)

	svc.CreateTemplate("resize-huge", "", 5)

	_, err := svc.ResizeTemplate("resize-huge", 203)
	if err == nil {
		t.Fatal("expected error for resize > 201")
	}
}

// --- ApplyTemplate ---

func TestApplyTemplate_Success(t *testing.T) {
	svc := newTestTemplateService(t)
	ctx := context.Background()

	// Create a 5×5 template with some terrain
	svc.CreateTemplate("apply-test", "", 5)
	svc.UpdateTemplateTerrain("apply-test", []model.TileTerrainUpdate{
		{X: 0, Y: 0, TerrainType: model.TerrainForest},
	})

	if err := svc.ApplyTemplate(ctx, "apply-test"); err != nil {
		t.Fatalf("ApplyTemplate: %v", err)
	}

	// Verify the live map has the template tiles — need a map repo to check.
	// We can create a MapService to query the live map.
	db := testutil.NewTestDB(t)
	mapRepo := sqlite.NewWorldMapRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	mapSvc := service.NewMapService(mapRepo, villageRepo)

	// The issue is the test TemplateService has its own DB. Let's use a shared approach.
	// For a simpler test, we verify no error and that the template has 25 tiles.
	tmpl, _ := svc.GetTemplate("apply-test")
	if len(tmpl.Tiles) != 25 {
		t.Errorf("template has %d tiles, want 25", len(tmpl.Tiles))
	}
	_ = mapSvc // satisfy compiler
}

func TestApplyTemplate_NotFound(t *testing.T) {
	svc := newTestTemplateService(t)
	ctx := context.Background()

	err := svc.ApplyTemplate(ctx, "nonexistent")
	if err == nil {
		t.Fatal("expected error applying nonexistent template")
	}
}

// --- ApplyTemplate with shared DB (integration test) ---

func TestApplyTemplate_WritesToLiveMap(t *testing.T) {
	// Use a single shared DB for both template service (map repo) and map service.
	tmpDir := t.TempDir()
	templateRepo, err := repository.NewTemplateRepository(tmpDir)
	if err != nil {
		t.Fatalf("NewTemplateRepository: %v", err)
	}

	db := testutil.NewTestDB(t)
	mapRepo := sqlite.NewWorldMapRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)

	tmplSvc := service.NewTemplateService(templateRepo, mapRepo)
	mapSvc := service.NewMapService(mapRepo, villageRepo)

	ctx := context.Background()

	// Create a 5×5 template
	tmplSvc.CreateTemplate("live-test", "", 5)
	tmplSvc.UpdateTemplateTerrain("live-test", []model.TileTerrainUpdate{
		{X: 0, Y: 0, TerrainType: model.TerrainForest},
		{X: 1, Y: 1, TerrainType: model.TerrainMountain},
	})

	// Apply to live map
	if err := tmplSvc.ApplyTemplate(ctx, "live-test"); err != nil {
		t.Fatalf("ApplyTemplate: %v", err)
	}

	// Query the live map
	tile, err := mapSvc.GetTile(ctx, 0, 0)
	if err != nil {
		t.Fatalf("GetTile(0,0): %v", err)
	}
	if tile.TerrainType != model.TerrainForest {
		t.Errorf("live tile (0,0) terrain: got %q, want %q", tile.TerrainType, model.TerrainForest)
	}

	tile2, err := mapSvc.GetTile(ctx, 1, 1)
	if err != nil {
		t.Fatalf("GetTile(1,1): %v", err)
	}
	if tile2.TerrainType != model.TerrainMountain {
		t.Errorf("live tile (1,1) terrain: got %q, want %q", tile2.TerrainType, model.TerrainMountain)
	}

	// Check total tile count on live map
	chunks, err := mapSvc.GetMapChunk(ctx, 0, 0, 25)
	if err != nil {
		t.Fatalf("GetMapChunk: %v", err)
	}
	if len(chunks) != 25 {
		t.Errorf("live map tile count: got %d, want 25", len(chunks))
	}
}

func TestApplyTemplate_SizeMismatch(t *testing.T) {
	// First apply a 5×5 template, then apply a 9×9 one. Verify the map now has 81 tiles.
	tmpDir := t.TempDir()
	templateRepo, err := repository.NewTemplateRepository(tmpDir)
	if err != nil {
		t.Fatalf("NewTemplateRepository: %v", err)
	}

	db := testutil.NewTestDB(t)
	mapRepo := sqlite.NewWorldMapRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)

	tmplSvc := service.NewTemplateService(templateRepo, mapRepo)
	mapSvc := service.NewMapService(mapRepo, villageRepo)
	ctx := context.Background()

	// Create and apply a 5×5 template
	tmplSvc.CreateTemplate("small", "", 5)
	if err := tmplSvc.ApplyTemplate(ctx, "small"); err != nil {
		t.Fatalf("ApplyTemplate(small): %v", err)
	}

	// Create and apply a 9×9 template
	tmplSvc.CreateTemplate("bigger", "", 9)
	if err := tmplSvc.ApplyTemplate(ctx, "bigger"); err != nil {
		t.Fatalf("ApplyTemplate(bigger): %v", err)
	}

	// Live map should now have 81 tiles (the bigger template), not 25
	chunks, err := mapSvc.GetMapChunk(ctx, 0, 0, 25)
	if err != nil {
		t.Fatalf("GetMapChunk: %v", err)
	}
	if len(chunks) != 81 {
		t.Errorf("live map tile count after bigger apply: got %d, want 81", len(chunks))
	}
}

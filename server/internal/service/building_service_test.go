package service

import (
	"context"
	"testing"
	"time"

	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/testutil"
)

// helper: create a test building service with a player and village already set up.
func setupBuildingTest(t *testing.T) (*BuildingService, int64, int64) {
	t.Helper()

	db := testutil.NewTestDB(t)

	playerRepo := sqlite.NewPlayerRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	buildingRepo := sqlite.NewBuildingRepo(db)
	resourceRepo := sqlite.NewResourceRepo(db)
	queueRepo := sqlite.NewBuildingQueueRepo(db)

	svc := NewBuildingService(db, villageRepo, buildingRepo, resourceRepo, queueRepo, playerRepo)

	// Create a test player
	player := &model.Player{
		Username:     "testplayer",
		Email:        "test@test.com",
		PasswordHash: "$2a$12$dummy",
		Kingdom:      "veridor",
		Role:         model.RolePlayer,
		CreatedAt:    time.Now().UTC(),
	}
	if err := playerRepo.Create(context.Background(), player); err != nil {
		t.Fatalf("create player: %v", err)
	}

	// Create a village with starter buildings
	villageSvc := NewVillageService(villageRepo, buildingRepo, resourceRepo)
	village, err := villageSvc.CreateFirstVillage(context.Background(), player.ID, "veridor", "testplayer")
	if err != nil {
		t.Fatalf("create village: %v", err)
	}

	return svc, player.ID, village.ID
}

func TestStartUpgrade_Success(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	// Upgrade iron_mine from level 0 to level 1 (no prerequisites, affordable)
	resp, err := svc.StartUpgrade(ctx, playerID, villageID, "iron_mine")
	if err != nil {
		t.Fatalf("expected success, got error: %v", err)
	}

	if resp.BuildingType != "iron_mine" {
		t.Errorf("expected building_type iron_mine, got %s", resp.BuildingType)
	}
	if resp.TargetLevel != 1 {
		t.Errorf("expected target_level 1, got %d", resp.TargetLevel)
	}
	if resp.CompletesAt.Before(resp.StartedAt) {
		t.Error("completes_at should be after started_at")
	}
}

func TestStartUpgrade_InsufficientResources(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	// Upgrade town_hall repeatedly until resources are depleted
	// First, drain resources by upgrading iron_mine (costs 100/80/50/30 = total 260)
	// Starting resources are 500 each.
	// Town Hall costs 200/200/200/100 at level 1.
	_, err := svc.StartUpgrade(ctx, playerID, villageID, "town_hall")
	if err != nil {
		t.Fatalf("first upgrade should succeed: %v", err)
	}

	// Complete the first build manually so we can attempt another
	if err := svc.CompleteBuilds(ctx); err != nil {
		// won't complete yet (time hasn't passed), cancel instead
	}
	// Cancel the upgrade so we can start another
	queue, _ := svc.queueRepo.GetByVillageID(ctx, villageID)
	if len(queue) > 0 {
		_ = svc.CancelUpgrade(ctx, playerID, villageID, queue[0].ID)
	}

	// After spending 200/200/200/100, resources are now 300/300/300/400.
	// Try town_hall again (200/200/200/100) — should still work
	_, err = svc.StartUpgrade(ctx, playerID, villageID, "town_hall")
	if err != nil {
		t.Fatalf("second upgrade should succeed: %v", err)
	}
	queue, _ = svc.queueRepo.GetByVillageID(ctx, villageID)
	if len(queue) > 0 {
		_ = svc.CancelUpgrade(ctx, playerID, villageID, queue[0].ID)
	}

	// Now resources are 100/100/100/300 — town_hall costs 200/200/200/100 → should fail
	_, err = svc.StartUpgrade(ctx, playerID, villageID, "town_hall")
	if err == nil {
		t.Fatal("expected insufficient resources error")
	}
	if err != model.ErrInsufficientResources {
		t.Errorf("expected ErrInsufficientResources, got: %v", err)
	}
}

func TestStartUpgrade_BuildingInProgress(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	// Start first upgrade
	_, err := svc.StartUpgrade(ctx, playerID, villageID, "iron_mine")
	if err != nil {
		t.Fatalf("first upgrade should succeed: %v", err)
	}

	// Try to start second upgrade — should fail
	_, err = svc.StartUpgrade(ctx, playerID, villageID, "lumber_mill")
	if err == nil {
		t.Fatal("expected building in progress error")
	}
	if err != model.ErrBuildingInProgress {
		t.Errorf("expected ErrBuildingInProgress, got: %v", err)
	}
}

func TestStartUpgrade_MaxLevelReached(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	// Manually set iron_mine to max level (20)
	buildings, _ := svc.buildingRepo.GetByVillageID(ctx, villageID)
	for _, b := range buildings {
		if b.BuildingType == "iron_mine" {
			b.Level = 20
			_ = svc.buildingRepo.Update(ctx, b)
			break
		}
	}

	_, err := svc.StartUpgrade(ctx, playerID, villageID, "iron_mine")
	if err == nil {
		t.Fatal("expected max level error")
	}
	if err != model.ErrMaxLevelReached {
		t.Errorf("expected ErrMaxLevelReached, got: %v", err)
	}
}

func TestStartUpgrade_PrerequisitesNotMet(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	// Barracks requires Town Hall level 3. Town Hall is at level 0.
	_, err := svc.StartUpgrade(ctx, playerID, villageID, "barracks")
	if err == nil {
		t.Fatal("expected prerequisites not met error")
	}
	if !isPrereqError(err) {
		t.Errorf("expected ErrPrerequisitesNotMet, got: %v", err)
	}
}

func TestStartUpgrade_PrerequisitesMet(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	// Set Town Hall to level 3 to satisfy barracks prerequisite
	buildings, _ := svc.buildingRepo.GetByVillageID(ctx, villageID)
	for _, b := range buildings {
		if b.BuildingType == "town_hall" {
			b.Level = 3
			_ = svc.buildingRepo.Update(ctx, b)
			break
		}
	}

	// Barracks should now be upgradeable
	resp, err := svc.StartUpgrade(ctx, playerID, villageID, "barracks")
	if err != nil {
		t.Fatalf("expected success with prerequisites met, got: %v", err)
	}
	if resp.BuildingType != "barracks" {
		t.Errorf("expected barracks, got %s", resp.BuildingType)
	}
}

func TestStartUpgrade_NotOwner(t *testing.T) {
	svc, _, villageID := setupBuildingTest(t)
	ctx := context.Background()

	// Use a non-existent player ID
	_, err := svc.StartUpgrade(ctx, 99999, villageID, "iron_mine")
	if err == nil {
		t.Fatal("expected not owner error")
	}
	if err != ErrNotOwner {
		t.Errorf("expected ErrNotOwner, got: %v", err)
	}
}

func TestStartUpgrade_UnknownBuilding(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	_, err := svc.StartUpgrade(ctx, playerID, villageID, "nonexistent_building")
	if err == nil {
		t.Fatal("expected unknown building error")
	}
	if err != ErrUnknownBuilding {
		t.Errorf("expected ErrUnknownBuilding, got: %v", err)
	}
}

func TestCompleteBuilds_PromotesLevel(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	// Start an upgrade
	_, err := svc.StartUpgrade(ctx, playerID, villageID, "iron_mine")
	if err != nil {
		t.Fatalf("upgrade should succeed: %v", err)
	}

	// Manually set the queue item's completes_at to the past
	queue, _ := svc.queueRepo.GetByVillageID(ctx, villageID)
	if len(queue) == 0 {
		t.Fatal("expected a queue item")
	}
	// Delete and re-insert with past completion time
	_ = svc.queueRepo.Delete(ctx, queue[0].ID)
	pastItem := &model.BuildingQueue{
		VillageID:    villageID,
		BuildingType: "iron_mine",
		TargetLevel:  1,
		StartedAt:    time.Now().UTC().Add(-5 * time.Minute),
		CompletesAt:  time.Now().UTC().Add(-1 * time.Minute),
	}
	_ = svc.queueRepo.Insert(ctx, pastItem)

	// Run complete builds
	err = svc.CompleteBuilds(ctx)
	if err != nil {
		t.Fatalf("CompleteBuilds failed: %v", err)
	}

	// Verify building level is now 1
	buildings, _ := svc.buildingRepo.GetByVillageID(ctx, villageID)
	for _, b := range buildings {
		if b.BuildingType == "iron_mine" {
			if b.Level != 1 {
				t.Errorf("expected iron_mine level 1, got %d", b.Level)
			}
			break
		}
	}

	// Verify queue is empty
	queue, _ = svc.queueRepo.GetByVillageID(ctx, villageID)
	if len(queue) != 0 {
		t.Errorf("expected empty queue, got %d items", len(queue))
	}
}

func TestCompleteBuilds_UpdatesResourceRates(t *testing.T) {
	svc, _, villageID := setupBuildingTest(t)
	ctx := context.Background()

	// Insert a completed queue item for iron_mine → level 1
	pastItem := &model.BuildingQueue{
		VillageID:    villageID,
		BuildingType: "iron_mine",
		TargetLevel:  1,
		StartedAt:    time.Now().UTC().Add(-5 * time.Minute),
		CompletesAt:  time.Now().UTC().Add(-1 * time.Minute),
	}
	_ = svc.queueRepo.Insert(ctx, pastItem)

	err := svc.CompleteBuilds(ctx)
	if err != nil {
		t.Fatalf("CompleteBuilds failed: %v", err)
	}

	// Verify iron rate was updated: base(1) + level(1) * rate_per_level(2) = 3
	res, _ := svc.resourceRepo.Get(ctx, villageID)
	expectedRate := 1.0 + 2.0*1 // BaseResourceRate + RatePerLevel * level
	if res.IronRate != expectedRate {
		t.Errorf("expected iron_rate %.1f, got %.1f", expectedRate, res.IronRate)
	}
}

func TestGetUpgradeCost(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	cost, err := svc.GetUpgradeCost(ctx, playerID, villageID, "iron_mine")
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	if cost.BuildingType != "iron_mine" {
		t.Errorf("expected iron_mine, got %s", cost.BuildingType)
	}
	if cost.CurrentLevel != 0 {
		t.Errorf("expected current level 0, got %d", cost.CurrentLevel)
	}
	if cost.TargetLevel != 1 {
		t.Errorf("expected target level 1, got %d", cost.TargetLevel)
	}
	// Iron mine at level 1: base costs 100/80/50/30
	if cost.Iron != 100 || cost.Wood != 80 || cost.Stone != 50 || cost.Food != 30 {
		t.Errorf("unexpected costs: iron=%.0f wood=%.0f stone=%.0f food=%.0f", cost.Iron, cost.Wood, cost.Stone, cost.Food)
	}
	if cost.TimeSec <= 0 {
		t.Error("expected positive time_seconds")
	}
}

func TestCancelUpgrade(t *testing.T) {
	svc, playerID, villageID := setupBuildingTest(t)
	ctx := context.Background()

	resp, err := svc.StartUpgrade(ctx, playerID, villageID, "iron_mine")
	if err != nil {
		t.Fatalf("upgrade should succeed: %v", err)
	}

	err = svc.CancelUpgrade(ctx, playerID, villageID, resp.ID)
	if err != nil {
		t.Fatalf("cancel should succeed: %v", err)
	}

	// Verify queue is empty
	queue, _ := svc.queueRepo.GetByVillageID(ctx, villageID)
	if len(queue) != 0 {
		t.Errorf("expected empty queue after cancel, got %d items", len(queue))
	}
}

func isPrereqError(err error) bool {
	// errors.Is doesn't work with wrapped errors that add context via %w + extra text
	// so we check the string
	return err != nil && (err == model.ErrPrerequisitesNotMet || containsError(err, model.ErrPrerequisitesNotMet))
}

func containsError(err, target error) bool {
	for err != nil {
		if err == target {
			return true
		}
		unwrapped, ok := err.(interface{ Unwrap() error })
		if !ok {
			break
		}
		err = unwrapped.Unwrap()
	}
	return false
}

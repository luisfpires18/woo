package service

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/testutil"
)

// setupTrainingTest creates a TrainingService with an Arkazia player, village, and barracks at level 1.
func setupTrainingTest(t *testing.T) (*TrainingService, int64, int64, repository.BuildingRepository) {
	t.Helper()

	db := testutil.NewTestDB(t)

	playerRepo := sqlite.NewPlayerRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	buildingRepo := sqlite.NewBuildingRepo(db)
	resourceRepo := sqlite.NewResourceRepo(db)
	troopRepo := sqlite.NewTroopRepo(db)
	queueRepo := sqlite.NewTrainingQueueRepo(db)

	svc := NewTrainingService(sqlite.NewUnitOfWork(db), villageRepo, buildingRepo, resourceRepo, troopRepo, queueRepo, playerRepo)

	// Create an Arkazia player
	player := &model.Player{
		Username:     "gladiator",
		Email:        "glad@test.com",
		PasswordHash: "$2a$12$dummy",
		Kingdom:      "arkazia",
		Role:         model.RolePlayer,
		CreatedAt:    time.Now().UTC(),
	}
	if err := playerRepo.Create(context.Background(), player); err != nil {
		t.Fatalf("create player: %v", err)
	}

	// Create a village with starter buildings (includes barracks at level 0)
	villageSvc := NewVillageService(villageRepo, buildingRepo, resourceRepo, nil)
	village, err := villageSvc.CreateFirstVillage(context.Background(), player.ID, "arkazia", "gladiator")
	if err != nil {
		t.Fatalf("create village: %v", err)
	}

	// Level up barracks to 1 so training is possible
	buildings, err := buildingRepo.GetByVillageID(context.Background(), village.ID)
	if err != nil {
		t.Fatalf("get buildings: %v", err)
	}
	for _, b := range buildings {
		if b.BuildingType == "barracks" {
			b.Level = 1
			if err := buildingRepo.Update(context.Background(), b); err != nil {
				t.Fatalf("update barracks: %v", err)
			}
			break
		}
	}

	return svc, player.ID, village.ID, buildingRepo
}

func TestStartTraining_Success(t *testing.T) {
	svc, playerID, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	resp, err := svc.StartTraining(ctx, playerID, villageID, "iron_legionary", 2)
	if err != nil {
		t.Fatalf("expected success, got: %v", err)
	}

	if resp.TroopType != "iron_legionary" {
		t.Errorf("troop_type: got %q, want %q", resp.TroopType, "iron_legionary")
	}
	if resp.Quantity != 2 {
		t.Errorf("quantity: got %d, want 2", resp.Quantity)
	}
	if resp.EachDurationSec <= 0 {
		t.Errorf("each_duration_sec should be > 0, got %d", resp.EachDurationSec)
	}
}

func TestStartTraining_InsufficientResources(t *testing.T) {
	svc, playerID, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	// Iron legionary costs Food:100, Water:50, Lumber:60, Stone:40 per unit.
	// Starting resources are 500 each. Training 100 units should exceed resources.
	_, err := svc.StartTraining(ctx, playerID, villageID, "iron_legionary", 100)
	if err == nil {
		t.Fatal("expected insufficient resources error")
	}
	if !errors.Is(err, model.ErrInsufficientResources) {
		t.Errorf("error: got %v, want ErrInsufficientResources", err)
	}
}

func TestStartTraining_UnknownTroop(t *testing.T) {
	svc, playerID, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	_, err := svc.StartTraining(ctx, playerID, villageID, "dragon_rider", 1)
	if err == nil {
		t.Fatal("expected unknown troop error")
	}
	if !errors.Is(err, ErrUnknownTroop) {
		t.Errorf("error: got %v, want ErrUnknownTroop", err)
	}
}

func TestStartTraining_InvalidQuantity(t *testing.T) {
	svc, playerID, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	_, err := svc.StartTraining(ctx, playerID, villageID, "iron_legionary", 0)
	if err == nil {
		t.Fatal("expected invalid quantity error")
	}
	if !errors.Is(err, ErrInvalidQuantity) {
		t.Errorf("error: got %v, want ErrInvalidQuantity", err)
	}
}

func TestStartTraining_BuildingRequirement(t *testing.T) {
	svc, playerID, villageID, buildingRepo := setupTrainingTest(t)
	ctx := context.Background()

	// Reset barracks to level 0 — training should fail
	buildings, _ := buildingRepo.GetByVillageID(ctx, villageID)
	for _, b := range buildings {
		if b.BuildingType == "barracks" {
			b.Level = 0
			if err := buildingRepo.Update(ctx, b); err != nil {
				t.Fatalf("reset barracks: %v", err)
			}
			break
		}
	}

	_, err := svc.StartTraining(ctx, playerID, villageID, "iron_legionary", 1)
	if err == nil {
		t.Fatal("expected building requirement error")
	}
	if !errors.Is(err, ErrTrainingBuildingReq) {
		t.Errorf("error: got %v, want ErrTrainingBuildingReq", err)
	}
}

func TestStartTraining_WrongKingdom(t *testing.T) {
	// Create a Veridor player and try to train Arkazia troops.
	db := testutil.NewTestDB(t)

	playerRepo := sqlite.NewPlayerRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	buildingRepo := sqlite.NewBuildingRepo(db)
	resourceRepo := sqlite.NewResourceRepo(db)
	troopRepo := sqlite.NewTroopRepo(db)
	queueRepo := sqlite.NewTrainingQueueRepo(db)

	svc := NewTrainingService(sqlite.NewUnitOfWork(db), villageRepo, buildingRepo, resourceRepo, troopRepo, queueRepo, playerRepo)

	player := &model.Player{
		Username:     "sailor",
		Email:        "sail@test.com",
		PasswordHash: "$2a$12$dummy",
		Kingdom:      "veridor",
		Role:         model.RolePlayer,
		CreatedAt:    time.Now().UTC(),
	}
	if err := playerRepo.Create(context.Background(), player); err != nil {
		t.Fatalf("create player: %v", err)
	}

	villageSvc := NewVillageService(villageRepo, buildingRepo, resourceRepo, nil)
	village, err := villageSvc.CreateFirstVillage(context.Background(), player.ID, "veridor", "sailor")
	if err != nil {
		t.Fatalf("create village: %v", err)
	}

	// Level up barracks
	buildings, _ := buildingRepo.GetByVillageID(context.Background(), village.ID)
	for _, b := range buildings {
		if b.BuildingType == "barracks" {
			b.Level = 1
			buildingRepo.Update(context.Background(), b)
			break
		}
	}

	_, err = svc.StartTraining(context.Background(), player.ID, village.ID, "iron_legionary", 1)
	if err == nil {
		t.Fatal("expected kingdom restriction error")
	}
}

func TestStartTraining_NotOwner(t *testing.T) {
	svc, _, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	// Use a non-existent player ID
	_, err := svc.StartTraining(ctx, 9999, villageID, "iron_legionary", 1)
	if err == nil {
		t.Fatal("expected not-owner error")
	}
	if !errors.Is(err, ErrNotOwner) {
		t.Errorf("error: got %v, want ErrNotOwner", err)
	}
}

func TestCompleteTraining_Success(t *testing.T) {
	svc, playerID, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	// Start training 2 iron legionaries
	_, err := svc.StartTraining(ctx, playerID, villageID, "iron_legionary", 2)
	if err != nil {
		t.Fatalf("start training: %v", err)
	}

	// Fast-forward: directly update the queue so completes_at is in the past
	queue, err := svc.queueRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		t.Fatalf("get queue: %v", err)
	}
	if len(queue) != 1 {
		t.Fatalf("expected 1 queue item, got %d", len(queue))
	}
	queue[0].CompletesAt = time.Now().UTC().Add(-1 * time.Second)
	if err := svc.queueRepo.Update(ctx, queue[0]); err != nil {
		t.Fatalf("update queue: %v", err)
	}

	// Complete training — should complete 1 unit
	events, err := svc.CompleteTraining(ctx)
	if err != nil {
		t.Fatalf("complete training: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}
	if events[0].TroopType != "iron_legionary" {
		t.Errorf("troop_type: got %q, want %q", events[0].TroopType, "iron_legionary")
	}
	if events[0].NewTotal != 1 {
		t.Errorf("new_total: got %d, want 1", events[0].NewTotal)
	}

	// Remaining queue should have quantity 1
	queue2, err := svc.queueRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		t.Fatalf("get queue after complete: %v", err)
	}
	if len(queue2) != 1 {
		t.Fatalf("expected 1 remaining queue item, got %d", len(queue2))
	}
	if queue2[0].Quantity != 1 {
		t.Errorf("remaining quantity: got %d, want 1", queue2[0].Quantity)
	}
}

func TestCompleteTraining_LastUnit(t *testing.T) {
	svc, playerID, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	// Start training 1 iron legionary
	_, err := svc.StartTraining(ctx, playerID, villageID, "iron_legionary", 1)
	if err != nil {
		t.Fatalf("start training: %v", err)
	}

	// Fast-forward
	queue, _ := svc.queueRepo.GetByVillageID(ctx, villageID)
	queue[0].CompletesAt = time.Now().UTC().Add(-1 * time.Second)
	svc.queueRepo.Update(ctx, queue[0])

	events, err := svc.CompleteTraining(ctx)
	if err != nil {
		t.Fatalf("complete training: %v", err)
	}
	if len(events) != 1 {
		t.Fatalf("expected 1 event, got %d", len(events))
	}

	// Queue should now be empty
	queue2, _ := svc.queueRepo.GetByVillageID(ctx, villageID)
	if len(queue2) != 0 {
		t.Errorf("queue should be empty, got %d items", len(queue2))
	}

	// Verify troops
	troops, err := svc.GetTroops(ctx, villageID)
	if err != nil {
		t.Fatalf("get troops: %v", err)
	}
	if len(troops) != 1 {
		t.Fatalf("expected 1 troop type, got %d", len(troops))
	}
	if troops[0].Type != "iron_legionary" || troops[0].Quantity != 1 {
		t.Errorf("troop: got type=%q qty=%d, want iron_legionary qty=1", troops[0].Type, troops[0].Quantity)
	}
}

func TestCancelTraining(t *testing.T) {
	svc, playerID, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	resp, err := svc.StartTraining(ctx, playerID, villageID, "iron_legionary", 3)
	if err != nil {
		t.Fatalf("start training: %v", err)
	}

	err = svc.CancelTraining(ctx, playerID, villageID, resp.ID)
	if err != nil {
		t.Fatalf("cancel training: %v", err)
	}

	// Queue should be empty
	queue, _ := svc.queueRepo.GetByVillageID(ctx, villageID)
	if len(queue) != 0 {
		t.Errorf("queue should be empty after cancel, got %d", len(queue))
	}
}

func TestGetTrainingCost(t *testing.T) {
	svc, playerID, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	costResp, err := svc.GetTrainingCost(ctx, playerID, villageID, "iron_legionary", 5)
	if err != nil {
		t.Fatalf("get training cost: %v", err)
	}

	// Iron legionary: Food:100, Water:50, Lumber:60, Stone:40 per unit.
	if costResp.TotalFood != 500 {
		t.Errorf("total_food: got %.0f, want 500", costResp.TotalFood)
	}
	if costResp.TotalWater != 250 {
		t.Errorf("total_water: got %.0f, want 250", costResp.TotalWater)
	}
	if costResp.TotalLumber != 300 {
		t.Errorf("total_lumber: got %.0f, want 300", costResp.TotalLumber)
	}
	if costResp.TotalStone != 200 {
		t.Errorf("total_stone: got %.0f, want 200", costResp.TotalStone)
	}
	if costResp.Quantity != 5 {
		t.Errorf("quantity: got %d, want 5", costResp.Quantity)
	}
}

func TestGetTrainingQueue(t *testing.T) {
	svc, playerID, villageID, _ := setupTrainingTest(t)
	ctx := context.Background()

	// Initially empty
	queue, err := svc.GetTrainingQueue(ctx, villageID)
	if err != nil {
		t.Fatalf("get queue: %v", err)
	}
	if len(queue) != 0 {
		t.Errorf("expected empty queue, got %d", len(queue))
	}

	// Train some units
	svc.StartTraining(ctx, playerID, villageID, "iron_legionary", 1)

	queue, err = svc.GetTrainingQueue(ctx, villageID)
	if err != nil {
		t.Fatalf("get queue: %v", err)
	}
	if len(queue) != 1 {
		t.Errorf("expected 1 queue item, got %d", len(queue))
	}
}

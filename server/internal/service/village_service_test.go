package service_test

import (
	"context"
	"testing"
	"time"

	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/service"
	"github.com/luisfpires18/woo/internal/testutil"
)

func newTestVillageService(t *testing.T) (*service.VillageService, *model.Player) {
	t.Helper()
	db := testutil.NewTestDB(t)

	playerRepo := sqlite.NewPlayerRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	buildingRepo := sqlite.NewBuildingRepo(db)
	resourceRepo := sqlite.NewResourceRepo(db)

	// Create a test player
	player := &model.Player{
		Username:     "testplayer",
		Email:        "test@example.com",
		PasswordHash: "not-a-real-hash",
		Kingdom:      "veridor",
		CreatedAt:    time.Now().UTC(),
	}
	if err := playerRepo.Create(context.Background(), player); err != nil {
		t.Fatalf("create test player: %v", err)
	}

	svc := service.NewVillageService(villageRepo, buildingRepo, resourceRepo, nil)
	return svc, player
}

func TestCreateFirstVillage_Success(t *testing.T) {
	svc, player := newTestVillageService(t)
	ctx := context.Background()

	village, err := svc.CreateFirstVillage(ctx, player.ID, player.Kingdom, player.Username)
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}

	if village.PlayerID != player.ID {
		t.Errorf("expected player ID %d, got %d", player.ID, village.PlayerID)
	}
	if !village.IsCapital {
		t.Error("first village should be capital")
	}
	if village.Name != "testplayer's Village" {
		t.Errorf("expected name 'testplayer's Village', got %q", village.Name)
	}
	// Coordinates should be within map bounds
	if village.X < -200 || village.X > 200 || village.Y < -200 || village.Y > 200 {
		t.Errorf("village coords (%d, %d) out of bounds", village.X, village.Y)
	}
}

func TestGetVillage_Success(t *testing.T) {
	svc, player := newTestVillageService(t)
	ctx := context.Background()

	village, err := svc.CreateFirstVillage(ctx, player.ID, player.Kingdom, player.Username)
	if err != nil {
		t.Fatalf("create village failed: %v", err)
	}

	resp, err := svc.GetVillage(ctx, village.ID, player.ID)
	if err != nil {
		t.Fatalf("get village failed: %v", err)
	}

	if resp.ID != village.ID {
		t.Errorf("expected village ID %d, got %d", village.ID, resp.ID)
	}
	if len(resp.Buildings) == 0 {
		t.Error("expected buildings to be created")
	}

	// Veridor should have 21 buildings (town_hall + barracks + stable + archery + workshop + special + storage + provisions + reservoir + 12 resource fields)
	if len(resp.Buildings) != 21 {
		t.Errorf("expected 21 buildings for veridor, got %d", len(resp.Buildings))
	}

	// Check that archery, workshop, special exist (universal military buildings)
	wantTypes := map[string]bool{"archery": false, "workshop": false, "special": false}
	for _, b := range resp.Buildings {
		if _, ok := wantTypes[b.BuildingType]; ok {
			wantTypes[b.BuildingType] = true
		}
	}
	for bt, found := range wantTypes {
		if !found {
			t.Errorf("village should have %s building", bt)
		}
	}

	// Check resources (allow small delta due to lazy calculation — rates are per-second now)
	if resp.Resources == nil {
		t.Fatal("expected resources")
	}
	if resp.Resources.Food < 499 || resp.Resources.Food > 510 {
		t.Errorf("expected ~500 food, got %f", resp.Resources.Food)
	}
	if resp.Resources.MaxFood != 1200 {
		t.Errorf("expected 1200 max food storage, got %f", resp.Resources.MaxFood)
	}
	if resp.Resources.MaxWater != 1200 {
		t.Errorf("expected 1200 max water storage, got %f", resp.Resources.MaxWater)
	}
	if resp.Resources.MaxLumber != 1200 {
		t.Errorf("expected 1200 max lumber storage, got %f", resp.Resources.MaxLumber)
	}
	if resp.Resources.MaxStone != 1200 {
		t.Errorf("expected 1200 max stone storage, got %f", resp.Resources.MaxStone)
	}
}

func TestGetVillage_NotOwner(t *testing.T) {
	svc, player := newTestVillageService(t)
	ctx := context.Background()

	village, err := svc.CreateFirstVillage(ctx, player.ID, player.Kingdom, player.Username)
	if err != nil {
		t.Fatalf("create village failed: %v", err)
	}

	// Try to access with a different player ID
	_, err = svc.GetVillage(ctx, village.ID, player.ID+999)
	if err != service.ErrNotOwner {
		t.Errorf("expected ErrNotOwner, got: %v", err)
	}
}

func TestGetVillage_NotFound(t *testing.T) {
	svc, _ := newTestVillageService(t)
	ctx := context.Background()

	_, err := svc.GetVillage(ctx, 99999, 1)
	if err != service.ErrVillageNotFound {
		t.Errorf("expected ErrVillageNotFound, got: %v", err)
	}
}

func TestListVillages(t *testing.T) {
	svc, player := newTestVillageService(t)
	ctx := context.Background()

	_, err := svc.CreateFirstVillage(ctx, player.ID, player.Kingdom, player.Username)
	if err != nil {
		t.Fatalf("create village failed: %v", err)
	}

	villages, err := svc.ListVillages(ctx, player.ID)
	if err != nil {
		t.Fatalf("list villages failed: %v", err)
	}
	if len(villages) != 1 {
		t.Errorf("expected 1 village, got %d", len(villages))
	}
	if !villages[0].IsCapital {
		t.Error("first village should be capital")
	}
}

func TestCreateFirstVillage_Sylvara(t *testing.T) {
	db := testutil.NewTestDB(t)
	playerRepo := sqlite.NewPlayerRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	buildingRepo := sqlite.NewBuildingRepo(db)
	resourceRepo := sqlite.NewResourceRepo(db)
	ctx := context.Background()

	player := &model.Player{
		Username:     "sylvara_player",
		Email:        "sylvara@example.com",
		PasswordHash: "not-a-real-hash",
		Kingdom:      "sylvara",
		CreatedAt:    time.Now().UTC(),
	}
	if err := playerRepo.Create(ctx, player); err != nil {
		t.Fatalf("create player: %v", err)
	}

	svc := service.NewVillageService(villageRepo, buildingRepo, resourceRepo, nil)
	village, err := svc.CreateFirstVillage(ctx, player.ID, player.Kingdom, player.Username)
	if err != nil {
		t.Fatalf("create village: %v", err)
	}

	resp, err := svc.GetVillage(ctx, village.ID, player.ID)
	if err != nil {
		t.Fatalf("get village: %v", err)
	}

	// Sylvara gets same 21 buildings as every kingdom (no kingdom-specific buildings)
	if len(resp.Buildings) != 21 {
		t.Errorf("expected 21 buildings for sylvara, got %d", len(resp.Buildings))
	}
}

func TestCreateFirstVillage_Arkazia(t *testing.T) {
	db := testutil.NewTestDB(t)
	playerRepo := sqlite.NewPlayerRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	buildingRepo := sqlite.NewBuildingRepo(db)
	resourceRepo := sqlite.NewResourceRepo(db)
	ctx := context.Background()

	player := &model.Player{
		Username:     "arkazia_player",
		Email:        "arkazia@example.com",
		PasswordHash: "not-a-real-hash",
		Kingdom:      "arkazia",
		CreatedAt:    time.Now().UTC(),
	}
	if err := playerRepo.Create(ctx, player); err != nil {
		t.Fatalf("create player: %v", err)
	}

	svc := service.NewVillageService(villageRepo, buildingRepo, resourceRepo, nil)
	village, err := svc.CreateFirstVillage(ctx, player.ID, player.Kingdom, player.Username)
	if err != nil {
		t.Fatalf("create village: %v", err)
	}

	resp, err := svc.GetVillage(ctx, village.ID, player.ID)
	if err != nil {
		t.Fatalf("get village: %v", err)
	}

	// Arkazia gets same 21 buildings as every kingdom (no kingdom-specific buildings)
	if len(resp.Buildings) != 21 {
		t.Errorf("expected 21 buildings for arkazia, got %d", len(resp.Buildings))
	}
}

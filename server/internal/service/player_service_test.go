package service_test

import (
	"context"
	"errors"
	"testing"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/service"
	"github.com/luisfpires18/woo/internal/testutil"
)

func newPlayerTestEnv(t *testing.T) (*service.PlayerService, int64) {
	t.Helper()
	db := testutil.NewTestDB(t)

	playerRepo := sqlite.NewPlayerRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	buildingRepo := sqlite.NewBuildingRepo(db)
	resourceRepo := sqlite.NewResourceRepo(db)
	worldMapRepo := sqlite.NewWorldMapRepo(db)
	refreshTokenRepo := sqlite.NewRefreshTokenRepo(db)
	playerEconRepo := sqlite.NewPlayerEconomyRepo(db)

	authService := service.NewAuthService(playerRepo, refreshTokenRepo, "test-secret", "woo-test")
	mapService := service.NewMapService(worldMapRepo, villageRepo)
	villageService := service.NewVillageService(villageRepo, buildingRepo, resourceRepo, playerEconRepo, mapService)
	playerService := service.NewPlayerService(playerRepo, villageService)

	// Generate map so village creation works
	if err := mapService.GenerateMap(context.Background()); err != nil {
		t.Fatalf("GenerateMap: %v", err)
	}

	// Register a player
	resp, err := authService.Register(context.Background(), &dto.RegisterRequest{
		Username: "testplayer",
		Email:    "test@test.com",
		Password: "Strong@123",
	})
	if err != nil {
		t.Fatalf("Register: %v", err)
	}

	return playerService, resp.Player.ID
}

func TestPlayerService_GetMe_Success(t *testing.T) {
	svc, playerID := newPlayerTestEnv(t)

	info, err := svc.GetMe(context.Background(), playerID)
	if err != nil {
		t.Fatalf("GetMe: %v", err)
	}

	if info.Username != "testplayer" {
		t.Errorf("username: got %q, want %q", info.Username, "testplayer")
	}
	if info.ID != playerID {
		t.Errorf("id: got %d, want %d", info.ID, playerID)
	}
}

func TestPlayerService_GetMe_NotFound(t *testing.T) {
	svc, _ := newPlayerTestEnv(t)

	_, err := svc.GetMe(context.Background(), 9999)
	if err == nil {
		t.Fatal("expected error for nonexistent player")
	}
	if !errors.Is(err, model.ErrNotFound) {
		t.Errorf("error: got %v, want ErrNotFound", err)
	}
}

func TestPlayerService_ChooseKingdom_Success(t *testing.T) {
	svc, playerID := newPlayerTestEnv(t)

	info, villageID, err := svc.ChooseKingdom(context.Background(), playerID, "veridor")
	if err != nil {
		t.Fatalf("ChooseKingdom: %v", err)
	}

	if info.Kingdom != "veridor" {
		t.Errorf("kingdom: got %q, want %q", info.Kingdom, "veridor")
	}
	if villageID == 0 {
		t.Error("expected non-zero village ID")
	}
}

func TestPlayerService_ChooseKingdom_InvalidKingdom(t *testing.T) {
	svc, playerID := newPlayerTestEnv(t)

	_, _, err := svc.ChooseKingdom(context.Background(), playerID, "atlantis")
	if !errors.Is(err, service.ErrInvalidKingdom) {
		t.Errorf("error: got %v, want ErrInvalidKingdom", err)
	}
}

func TestPlayerService_ChooseKingdom_AlreadyChosen(t *testing.T) {
	svc, playerID := newPlayerTestEnv(t)

	_, _, err := svc.ChooseKingdom(context.Background(), playerID, "veridor")
	if err != nil {
		t.Fatalf("first ChooseKingdom: %v", err)
	}

	_, _, err = svc.ChooseKingdom(context.Background(), playerID, "sylvara")
	if !errors.Is(err, service.ErrKingdomAlreadyChosen) {
		t.Errorf("error: got %v, want ErrKingdomAlreadyChosen", err)
	}
}

func TestPlayerService_ChooseKingdom_PlayerNotFound(t *testing.T) {
	svc, _ := newPlayerTestEnv(t)

	_, _, err := svc.ChooseKingdom(context.Background(), 9999, "veridor")
	if err == nil {
		t.Fatal("expected error for nonexistent player")
	}
}

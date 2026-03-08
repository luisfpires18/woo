package service_test

import (
	"context"
	"testing"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/service"
	"github.com/luisfpires18/woo/internal/testutil"
)

func newTestAdminService(t *testing.T) (*service.AdminService, *service.AuthService) {
	t.Helper()
	db := testutil.NewTestDB(t)
	playerRepo := sqlite.NewPlayerRepo(db)
	villageRepo := sqlite.NewVillageRepo(db)
	worldConfigRepo := sqlite.NewWorldConfigRepo(db)
	announcementRepo := sqlite.NewAnnouncementRepo(db)
	refreshTokenRepo := sqlite.NewRefreshTokenRepo(db)
	gameAssetRepo := sqlite.NewGameAssetRepo(db)
	resBuildingConfigRepo := sqlite.NewResourceBuildingConfigRepo(db)
	buildingDisplayConfigRepo := sqlite.NewBuildingDisplayConfigRepo(db)

	adminSvc := service.NewAdminService(playerRepo, villageRepo, worldConfigRepo, announcementRepo, gameAssetRepo, resBuildingConfigRepo, buildingDisplayConfigRepo)
	authSvc := service.NewAuthService(playerRepo, refreshTokenRepo, "test-secret", "woo-test")
	return adminSvc, authSvc
}

func TestListPlayers(t *testing.T) {
	adminSvc, authSvc := newTestAdminService(t)
	ctx := context.Background()

	// Register two players (seed admin + 2 more via register = 3 total)
	authSvc.Register(ctx, &dto.RegisterRequest{
		Username: "player1", Email: "p1@test.com", Password: "password123",
	})
	authSvc.Register(ctx, &dto.RegisterRequest{
		Username: "player2", Email: "p2@test.com", Password: "password123",
	})

	resp, err := adminSvc.ListPlayers(ctx, 0, 20)
	if err != nil {
		t.Fatalf("list players: %v", err)
	}

	// Seed admin (from migration 019) + 2 registered = 3
	if resp.Total < 2 {
		t.Errorf("expected at least 2 players, got %d", resp.Total)
	}
	if len(resp.Players) < 2 {
		t.Errorf("expected at least 2 players in list, got %d", len(resp.Players))
	}
}

func TestListPlayers_Pagination(t *testing.T) {
	adminSvc, authSvc := newTestAdminService(t)
	ctx := context.Background()

	for i := 0; i < 5; i++ {
		authSvc.Register(ctx, &dto.RegisterRequest{
			Username: "pageplayer" + string(rune('a'+i)),
			Email:    "page" + string(rune('a'+i)) + "@test.com",
			Password: "password123",
		})
	}

	resp, err := adminSvc.ListPlayers(ctx, 0, 2)
	if err != nil {
		t.Fatalf("list players page 1: %v", err)
	}
	if len(resp.Players) != 2 {
		t.Errorf("expected 2 players in page, got %d", len(resp.Players))
	}
	if resp.Limit != 2 {
		t.Errorf("expected limit 2, got %d", resp.Limit)
	}
}

func TestUpdatePlayerRole(t *testing.T) {
	adminSvc, authSvc := newTestAdminService(t)
	ctx := context.Background()

	resp, err := authSvc.Register(ctx, &dto.RegisterRequest{
		Username: "upgradeuser", Email: "upgrade@test.com", Password: "password123",
	})
	if err != nil {
		t.Fatalf("register: %v", err)
	}

	// Promote to admin
	if err := adminSvc.UpdatePlayerRole(ctx, resp.Player.ID, model.RoleAdmin); err != nil {
		t.Fatalf("update role to admin: %v", err)
	}

	// Demote back to player
	if err := adminSvc.UpdatePlayerRole(ctx, resp.Player.ID, model.RolePlayer); err != nil {
		t.Fatalf("update role to player: %v", err)
	}
}

func TestUpdatePlayerRole_InvalidRole(t *testing.T) {
	adminSvc, _ := newTestAdminService(t)
	ctx := context.Background()

	err := adminSvc.UpdatePlayerRole(ctx, 1, "superadmin")
	if err != service.ErrInvalidRole {
		t.Errorf("expected ErrInvalidRole, got: %v", err)
	}
}

func TestGetWorldConfig(t *testing.T) {
	adminSvc, _ := newTestAdminService(t)
	ctx := context.Background()

	resp, err := adminSvc.GetWorldConfig(ctx)
	if err != nil {
		t.Fatalf("get world config: %v", err)
	}

	if len(resp.Configs) < 6 {
		t.Errorf("expected at least 6 config entries, got %d", len(resp.Configs))
	}

	// Verify one known config
	found := false
	for _, c := range resp.Configs {
		if c.Key == "game_speed" && c.Value == "1.0" {
			found = true
			break
		}
	}
	if !found {
		t.Error("expected game_speed config with value 1.0")
	}
}

func TestSetWorldConfig(t *testing.T) {
	adminSvc, _ := newTestAdminService(t)
	ctx := context.Background()

	if err := adminSvc.SetWorldConfig(ctx, "game_speed", "2.0"); err != nil {
		t.Fatalf("set config: %v", err)
	}

	// Verify it changed
	resp, _ := adminSvc.GetWorldConfig(ctx)
	for _, c := range resp.Configs {
		if c.Key == "game_speed" {
			if c.Value != "2.0" {
				t.Errorf("expected game_speed 2.0, got %q", c.Value)
			}
			return
		}
	}
	t.Error("game_speed config not found after update")
}

func TestSetWorldConfig_NotFound(t *testing.T) {
	adminSvc, _ := newTestAdminService(t)
	ctx := context.Background()

	err := adminSvc.SetWorldConfig(ctx, "nonexistent_key", "value")
	if err == nil {
		t.Error("expected error for nonexistent config key")
	}
}

func TestGetStats(t *testing.T) {
	adminSvc, authSvc := newTestAdminService(t)
	ctx := context.Background()

	authSvc.Register(ctx, &dto.RegisterRequest{
		Username: "statsplayer", Email: "stats@test.com", Password: "password123",
	})

	stats, err := adminSvc.GetStats(ctx)
	if err != nil {
		t.Fatalf("get stats: %v", err)
	}
	if stats.TotalPlayers < 1 {
		t.Errorf("expected at least 1 player, got %d", stats.TotalPlayers)
	}
}

func TestAnnouncements_CRUD(t *testing.T) {
	adminSvc, authSvc := newTestAdminService(t)
	ctx := context.Background()

	// Need a player ID for author
	resp, _ := authSvc.Register(ctx, &dto.RegisterRequest{
		Username: "announcer", Email: "ann@test.com", Password: "password123",
	})
	authorID := resp.Player.ID

	// Create
	ann, err := adminSvc.CreateAnnouncement(ctx, &dto.CreateAnnouncementRequest{
		Title:   "Server Maintenance",
		Content: "The server will be down for maintenance at midnight.",
	}, authorID)
	if err != nil {
		t.Fatalf("create announcement: %v", err)
	}
	if ann.ID == 0 {
		t.Error("expected non-zero announcement ID")
	}
	if ann.Title != "Server Maintenance" {
		t.Errorf("expected title 'Server Maintenance', got %q", ann.Title)
	}

	// List
	list, err := adminSvc.ListAnnouncements(ctx)
	if err != nil {
		t.Fatalf("list announcements: %v", err)
	}
	if len(list) < 1 {
		t.Fatal("expected at least 1 announcement")
	}

	// Delete
	if err := adminSvc.DeleteAnnouncement(ctx, ann.ID); err != nil {
		t.Fatalf("delete announcement: %v", err)
	}

	// Verify deleted
	list, _ = adminSvc.ListAnnouncements(ctx)
	for _, a := range list {
		if a.ID == ann.ID {
			t.Error("announcement should have been deleted")
		}
	}
}

func TestCreateAnnouncement_Empty(t *testing.T) {
	adminSvc, _ := newTestAdminService(t)
	ctx := context.Background()

	_, err := adminSvc.CreateAnnouncement(ctx, &dto.CreateAnnouncementRequest{
		Title:   "",
		Content: "some content",
	}, 1)
	if err == nil {
		t.Error("expected error for empty title")
	}
}

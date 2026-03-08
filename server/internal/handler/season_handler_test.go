package handler_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/luisfpires18/woo/internal/dto"
)

// ── Admin: Create Season ─────────────────────────────────────────────────────

func TestAdminCreateSeason_Success(t *testing.T) {
	env := newTestEnv(t)

	body, _ := json.Marshal(dto.CreateSeasonRequest{
		Name:        "Season 1",
		Description: "First test season",
		GameSpeed:   2.0,
		MapWidth:    51,
		MapHeight:   51,
	})

	req := httptest.NewRequest(http.MethodPost, "/seasons", bytes.NewReader(body))
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()

	env.SeasonHandler.AdminCreateSeason(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	var result struct {
		Season dto.SeasonResponse `json:"season"`
	}
	if err := json.Unmarshal(data, &result); err != nil {
		t.Fatalf("unmarshal: %v", err)
	}
	if result.Season.Name != "Season 1" {
		t.Errorf("expected name 'Season 1', got %q", result.Season.Name)
	}
	if result.Season.Status != "upcoming" {
		t.Errorf("expected status 'upcoming', got %q", result.Season.Status)
	}
	if result.Season.GameSpeed != 2.0 {
		t.Errorf("expected game_speed 2.0, got %f", result.Season.GameSpeed)
	}
}

func TestAdminCreateSeason_EmptyName(t *testing.T) {
	env := newTestEnv(t)

	body, _ := json.Marshal(dto.CreateSeasonRequest{Name: ""})
	req := httptest.NewRequest(http.MethodPost, "/seasons", bytes.NewReader(body))
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()

	env.SeasonHandler.AdminCreateSeason(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", rec.Code)
	}
}

// ── Admin: Launch / End / Archive lifecycle ───────────────────────────────────

func TestSeasonLifecycle(t *testing.T) {
	env := newTestEnv(t)

	// Create
	createBody, _ := json.Marshal(dto.CreateSeasonRequest{
		Name:     "Lifecycle Test",
		MapWidth: 51, MapHeight: 51,
	})
	req := httptest.NewRequest(http.MethodPost, "/seasons", bytes.NewReader(createBody))
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()
	env.SeasonHandler.AdminCreateSeason(rec, req)
	if rec.Code != http.StatusCreated {
		t.Fatalf("create: expected 201, got %d: %s", rec.Code, rec.Body.String())
	}

	data, _ := decodeEnvelope(t, rec)
	var created struct {
		Season dto.SeasonResponse `json:"season"`
	}
	json.Unmarshal(data, &created)
	seasonID := created.Season.ID
	idStr := fmt.Sprintf("%d", seasonID)

	// Launch
	req = httptest.NewRequest(http.MethodPost, "/seasons/"+idStr+"/launch", nil)
	req.SetPathValue("id", idStr)
	req = req.WithContext(authCtx(1, "admin"))
	rec = httptest.NewRecorder()
	env.SeasonHandler.AdminLaunchSeason(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("launch: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	data, _ = decodeEnvelope(t, rec)
	var launched struct {
		Season dto.SeasonResponse `json:"season"`
	}
	json.Unmarshal(data, &launched)
	if launched.Season.Status != "active" {
		t.Errorf("expected active, got %q", launched.Season.Status)
	}

	// End
	req = httptest.NewRequest(http.MethodPost, "/seasons/"+idStr+"/end", nil)
	req.SetPathValue("id", idStr)
	req = req.WithContext(authCtx(1, "admin"))
	rec = httptest.NewRecorder()
	env.SeasonHandler.AdminEndSeason(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("end: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	data, _ = decodeEnvelope(t, rec)
	var ended struct {
		Season dto.SeasonResponse `json:"season"`
	}
	json.Unmarshal(data, &ended)
	if ended.Season.Status != "ended" {
		t.Errorf("expected ended, got %q", ended.Season.Status)
	}

	// Archive
	req = httptest.NewRequest(http.MethodPost, "/seasons/"+idStr+"/archive", nil)
	req.SetPathValue("id", idStr)
	req = req.WithContext(authCtx(1, "admin"))
	rec = httptest.NewRecorder()
	env.SeasonHandler.AdminArchiveSeason(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("archive: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}
	data, _ = decodeEnvelope(t, rec)
	var archived struct {
		Season dto.SeasonResponse `json:"season"`
	}
	json.Unmarshal(data, &archived)
	if archived.Season.Status != "archived" {
		t.Errorf("expected archived, got %q", archived.Season.Status)
	}
	_ = seasonID
}

// ── Player: Join Season ──────────────────────────────────────────────────────

func TestJoinSeason_Success(t *testing.T) {
	env := newTestEnv(t)

	// Register a player
	resp, err := env.AuthService.Register(t.Context(), &dto.RegisterRequest{Username: "testplayer", Email: "test@example.com", Password: "Password1!"})
	if err != nil {
		t.Fatalf("register: %v", err)
	}
	playerID := resp.Player.ID

	// Admin creates + launches a season
	ctx := authCtx(1, "admin")
	season, err := env.SeasonService.CreateSeason(ctx, dto.CreateSeasonRequest{
		Name: "Joinable Season", MapWidth: 51, MapHeight: 51,
	})
	if err != nil {
		t.Fatalf("create season: %v", err)
	}
	_, err = env.SeasonService.LaunchSeason(ctx, season.ID)
	if err != nil {
		t.Fatalf("launch season: %v", err)
	}

	// Player joins
	body, _ := json.Marshal(dto.JoinSeasonRequest{Kingdom: "veridor"})
	req := httptest.NewRequest(http.MethodPost, "/1/join", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()

	env.SeasonHandler.JoinSeason(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("join: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	var result struct {
		Season    dto.SeasonDetailResponse `json:"season"`
		VillageID int64                    `json:"village_id"`
	}
	json.Unmarshal(data, &result)
	if !result.Season.Joined {
		t.Error("expected joined to be true")
	}
	if result.Season.Kingdom != "veridor" {
		t.Errorf("expected kingdom 'veridor', got %q", result.Season.Kingdom)
	}
	if result.VillageID == 0 {
		t.Error("expected a village ID")
	}
}

func TestJoinSeason_AlreadyJoined(t *testing.T) {
	env := newTestEnv(t)

	resp, _ := env.AuthService.Register(t.Context(), &dto.RegisterRequest{Username: "testplayer2", Email: "test2@example.com", Password: "Password1!"})
	playerID := resp.Player.ID

	ctx := authCtx(1, "admin")
	season, _ := env.SeasonService.CreateSeason(ctx, dto.CreateSeasonRequest{
		Name: "Double Join Test", MapWidth: 51, MapHeight: 51,
	})
	env.SeasonService.LaunchSeason(ctx, season.ID)

	// First join
	body, _ := json.Marshal(dto.JoinSeasonRequest{Kingdom: "sylvara"})
	req := httptest.NewRequest(http.MethodPost, "/1/join", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.SeasonHandler.JoinSeason(rec, req)
	if rec.Code != http.StatusOK {
		t.Fatalf("first join: expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	// Second join — should fail
	body, _ = json.Marshal(dto.JoinSeasonRequest{Kingdom: "arkazia"})
	req = httptest.NewRequest(http.MethodPost, "/1/join", bytes.NewReader(body))
	req.SetPathValue("id", "1")
	req = req.WithContext(authCtx(playerID, "player"))
	rec = httptest.NewRecorder()
	env.SeasonHandler.JoinSeason(rec, req)
	if rec.Code != http.StatusConflict {
		t.Fatalf("second join: expected 409, got %d: %s", rec.Code, rec.Body.String())
	}
}

// ── Player: List Seasons ─────────────────────────────────────────────────────

func TestListSeasons(t *testing.T) {
	env := newTestEnv(t)

	// Count pre-seeded seasons
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(authCtx(1, "player"))
	rec := httptest.NewRecorder()
	env.SeasonHandler.ListSeasons(rec, req)
	var baseline struct {
		Seasons []*dto.SeasonResponse `json:"seasons"`
	}
	data, _ := decodeEnvelope(t, rec)
	json.Unmarshal(data, &baseline)
	baseCount := len(baseline.Seasons)

	ctx := authCtx(1, "admin")
	env.SeasonService.CreateSeason(ctx, dto.CreateSeasonRequest{Name: "S1", MapWidth: 51, MapHeight: 51})
	env.SeasonService.CreateSeason(ctx, dto.CreateSeasonRequest{Name: "S2", MapWidth: 51, MapHeight: 51})

	req = httptest.NewRequest(http.MethodGet, "/", nil)
	req = req.WithContext(authCtx(1, "player"))
	rec = httptest.NewRecorder()

	env.SeasonHandler.ListSeasons(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", rec.Code)
	}

	data, _ = decodeEnvelope(t, rec)
	var result struct {
		Seasons []*dto.SeasonResponse `json:"seasons"`
	}
	json.Unmarshal(data, &result)
	if len(result.Seasons) != baseCount+2 {
		t.Errorf("expected %d seasons, got %d", baseCount+2, len(result.Seasons))
	}
}

// ── Player: Profile ──────────────────────────────────────────────────────────

func TestGetProfile(t *testing.T) {
	env := newTestEnv(t)

	// Register player + join a season
	resp, _ := env.AuthService.Register(t.Context(), &dto.RegisterRequest{Username: "profiletest", Email: "profile@example.com", Password: "Password1!"})
	playerID := resp.Player.ID

	ctx := authCtx(1, "admin")
	season, _ := env.SeasonService.CreateSeason(ctx, dto.CreateSeasonRequest{
		Name: "Profile Season", MapWidth: 51, MapHeight: 51,
	})
	env.SeasonService.LaunchSeason(ctx, season.ID)
	env.SeasonService.JoinSeason(authCtx(playerID, "player"), season.ID, playerID, "arkazia")

	// Get profile
	req := httptest.NewRequest(http.MethodGet, "/api/player/profile", nil)
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()

	env.PlayerHandler.GetProfile(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d: %s", rec.Code, rec.Body.String())
	}

	data, _ := decodeEnvelope(t, rec)
	var result struct {
		Profile dto.PlayerProfileResponse `json:"profile"`
	}
	json.Unmarshal(data, &result)
	if result.Profile.Username != "profiletest" {
		t.Errorf("expected username 'profiletest', got %q", result.Profile.Username)
	}
	if result.Profile.TotalSeasons != 1 {
		t.Errorf("expected 1 season, got %d", result.Profile.TotalSeasons)
	}
	if len(result.Profile.SeasonHistory) != 1 {
		t.Errorf("expected 1 history entry, got %d", len(result.Profile.SeasonHistory))
	}
	if result.Profile.SeasonHistory[0].Kingdom != "arkazia" {
		t.Errorf("expected kingdom 'arkazia', got %q", result.Profile.SeasonHistory[0].Kingdom)
	}
}

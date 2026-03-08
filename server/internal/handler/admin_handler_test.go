package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
)

// --- ListPlayers ---

func TestAdminHandler_ListPlayers(t *testing.T) {
	env := newTestEnv(t)

	// Register a couple of players (admin seed creates wright but test DB is fresh)
	registerAndLogin(t, env, "player1", "p1@test.com", "Strong@123")
	registerAndLogin(t, env, "player2", "p2@test.com", "Strong@123")

	req := httptest.NewRequest("GET", "/api/admin/players?offset=0&limit=10", nil)
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()
	env.AdminHandler.ListPlayers(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

// --- GetStats ---

func TestAdminHandler_GetStats(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest("GET", "/api/admin/stats", nil)
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()
	env.AdminHandler.GetStats(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

// --- GetWorldConfig ---

func TestAdminHandler_GetWorldConfig(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest("GET", "/api/admin/config", nil)
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()
	env.AdminHandler.GetWorldConfig(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	// Should have at least the default seed values
	var configResp struct {
		Configs []json.RawMessage `json:"configs"`
	}
	json.Unmarshal(data, &configResp)
	if len(configResp.Configs) < 3 {
		t.Errorf("expected at least 3 config entries, got %d", len(configResp.Configs))
	}
}

// --- SetWorldConfig ---

func TestAdminHandler_SetWorldConfig(t *testing.T) {
	env := newTestEnv(t)

	body := `{"value":"2.0"}`
	req := httptest.NewRequest("PUT", "/api/admin/config/game_speed", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("key", "game_speed")
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()
	env.AdminHandler.SetWorldConfig(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

// --- Announcements CRUD ---

func TestAdminHandler_Announcements_CRUD(t *testing.T) {
	env := newTestEnv(t)

	// Must register an admin-level player first (announcements FK to players.id)
	adminID, _ := registerAndLogin(t, env, "admin1", "admin1@test.com", "Strong@123")

	// Create announcement
	body := `{"title":"Test Announcement","content":"This is a test."}`
	createReq := httptest.NewRequest("POST", "/api/admin/announcements", strings.NewReader(body))
	createReq.Header.Set("Content-Type", "application/json")
	createReq = createReq.WithContext(authCtx(adminID, "admin"))
	createRec := httptest.NewRecorder()
	env.AdminHandler.CreateAnnouncement(createRec, createReq)

	if createRec.Code != http.StatusCreated {
		t.Fatalf("create status: got %d, want %d. Body: %s", createRec.Code, http.StatusCreated, createRec.Body.String())
	}

	// List announcements
	listReq := httptest.NewRequest("GET", "/api/admin/announcements", nil)
	listReq = listReq.WithContext(authCtx(adminID, "admin"))
	listRec := httptest.NewRecorder()
	env.AdminHandler.ListAnnouncements(listRec, listReq)

	if listRec.Code != http.StatusOK {
		t.Fatalf("list status: got %d, want %d", listRec.Code, http.StatusOK)
	}

	listData, _ := decodeEnvelope(t, listRec)
	var announcements []struct {
		ID int64 `json:"id"`
	}
	json.Unmarshal(listData, &announcements)
	if len(announcements) != 1 {
		t.Fatalf("announcement count: got %d, want 1", len(announcements))
	}

	// Delete announcement
	delReq := httptest.NewRequest("DELETE", "/api/admin/announcements/"+strconv.Itoa(int(announcements[0].ID)), nil)
	delReq.SetPathValue("id", strconv.Itoa(int(announcements[0].ID)))
	delReq = delReq.WithContext(authCtx(adminID, "admin"))
	delRec := httptest.NewRecorder()
	env.AdminHandler.DeleteAnnouncement(delRec, delReq)

	if delRec.Code != http.StatusOK {
		t.Fatalf("delete status: got %d, want %d", delRec.Code, http.StatusOK)
	}
}

func TestAdminHandler_UpdatePlayerRole(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "roleuser", "role@test.com", "Strong@123")

	body := `{"role":"admin"}`
	req := httptest.NewRequest("PATCH", "/api/admin/players/"+strconv.Itoa(int(playerID))+"/role", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", strconv.Itoa(int(playerID)))
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()
	env.AdminHandler.UpdatePlayerRole(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestAdminHandler_UpdatePlayerRole_InvalidRole(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "badrole", "badrole@test.com", "Strong@123")

	body := `{"role":"superadmin"}`
	req := httptest.NewRequest("PATCH", "/api/admin/players/"+strconv.Itoa(int(playerID))+"/role", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", strconv.Itoa(int(playerID)))
	req = req.WithContext(authCtx(1, "admin"))
	rec := httptest.NewRecorder()
	env.AdminHandler.UpdatePlayerRole(rec, req)

	if rec.Code == http.StatusOK {
		t.Error("expected error for invalid role, got 200")
	}
}

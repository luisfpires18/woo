package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/luisfpires18/woo/internal/dto"
)

// --- GetMe ---

func TestPlayerHandler_GetMe_Success(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "meuser", "me@test.com", "Strong@123")

	req := httptest.NewRequest("GET", "/api/player/me", nil)
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.PlayerHandler.GetMe(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	var resp struct {
		Player dto.PlayerInfo `json:"player"`
	}
	json.Unmarshal(data, &resp)
	if resp.Player.Username != "meuser" {
		t.Errorf("username: got %q, want %q", resp.Player.Username, "meuser")
	}
	if resp.Player.ID != playerID {
		t.Errorf("id: got %d, want %d", resp.Player.ID, playerID)
	}
}

func TestPlayerHandler_GetMe_Unauthenticated(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest("GET", "/api/player/me", nil)
	rec := httptest.NewRecorder()
	env.PlayerHandler.GetMe(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- ChooseKingdom ---

func TestPlayerHandler_ChooseKingdom_Success(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "kingdomuser", "kingdom@test.com", "Strong@123")

	body := `{"kingdom":"veridor"}`
	req := httptest.NewRequest("PUT", "/api/player/kingdom", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.PlayerHandler.ChooseKingdom(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	var resp struct {
		Player    dto.PlayerInfo `json:"player"`
		VillageID float64        `json:"village_id"`
	}
	json.Unmarshal(data, &resp)
	if resp.Player.Kingdom != "veridor" {
		t.Errorf("kingdom: got %q, want %q", resp.Player.Kingdom, "veridor")
	}
	if resp.VillageID == 0 {
		t.Error("expected non-zero village_id")
	}
}

func TestPlayerHandler_ChooseKingdom_InvalidKingdom(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "badkingdom", "badk@test.com", "Strong@123")

	body := `{"kingdom":"atlantis"}`
	req := httptest.NewRequest("PUT", "/api/player/kingdom", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.PlayerHandler.ChooseKingdom(rec, req)

	if rec.Code == http.StatusOK {
		t.Error("expected error for invalid kingdom, got 200")
	}
}

func TestPlayerHandler_ChooseKingdom_AlreadyChosen(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "doublechoose", "double@test.com", "Strong@123")

	// First choice
	body := `{"kingdom":"veridor"}`
	req := httptest.NewRequest("PUT", "/api/player/kingdom", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.PlayerHandler.ChooseKingdom(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("first choose failed: %d %s", rec.Code, rec.Body.String())
	}

	// Second choice — should be rejected
	body2 := `{"kingdom":"sylvara"}`
	req2 := httptest.NewRequest("PUT", "/api/player/kingdom", strings.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	req2 = req2.WithContext(authCtx(playerID, "player"))
	rec2 := httptest.NewRecorder()
	env.PlayerHandler.ChooseKingdom(rec2, req2)

	if rec2.Code == http.StatusOK {
		t.Error("expected error for already-chosen kingdom, got 200")
	}
}

func TestPlayerHandler_ChooseKingdom_Unauthenticated(t *testing.T) {
	env := newTestEnv(t)

	body := `{"kingdom":"veridor"}`
	req := httptest.NewRequest("PUT", "/api/player/kingdom", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.PlayerHandler.ChooseKingdom(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

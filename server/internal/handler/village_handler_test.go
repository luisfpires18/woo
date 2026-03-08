package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/luisfpires18/woo/internal/dto"
)

// registerAndLogin registers a player and returns the player ID and access token.
func registerAndLogin(t *testing.T, env *testEnv, username, email, password string) (int64, string) {
	t.Helper()
	body := `{"username":"` + username + `","email":"` + email + `","password":"` + password + `"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.AuthHandler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("register %s failed: %d %s", username, rec.Code, rec.Body.String())
	}

	data, _ := decodeEnvelope(t, rec)
	var resp dto.AuthResponse
	json.Unmarshal(data, &resp)
	return resp.Player.ID, resp.AccessToken
}

// chooseKingdomForPlayer calls ChooseKingdom handler for a given player.
func chooseKingdomForPlayer(t *testing.T, env *testEnv, playerID int64, kingdom string) {
	t.Helper()
	body := `{"kingdom":"` + kingdom + `"}`
	req := httptest.NewRequest("PUT", "/api/player/kingdom", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.PlayerHandler.ChooseKingdom(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("choose kingdom for player %d failed: %d %s", playerID, rec.Code, rec.Body.String())
	}
}

// --- ListVillages ---

func TestVillageHandler_ListVillages_Success(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "viluser", "vil@test.com", "Strong@123")
	chooseKingdomForPlayer(t, env, playerID, "veridor")

	req := httptest.NewRequest("GET", "/api/villages", nil)
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.VillageHandler.ListVillages(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	var villages []json.RawMessage
	json.Unmarshal(data, &villages)
	if len(villages) != 1 {
		t.Errorf("village count: got %d, want 1", len(villages))
	}
}

func TestVillageHandler_ListVillages_Unauthenticated(t *testing.T) {
	env := newTestEnv(t)

	req := httptest.NewRequest("GET", "/api/villages", nil)
	rec := httptest.NewRecorder()
	env.VillageHandler.ListVillages(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- GetVillage ---

func TestVillageHandler_GetVillage_Success(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "getviluser", "getvil@test.com", "Strong@123")
	chooseKingdomForPlayer(t, env, playerID, "sylvara")

	// First list villages to get the ID
	listReq := httptest.NewRequest("GET", "/api/villages", nil)
	listReq = listReq.WithContext(authCtx(playerID, "player"))
	listRec := httptest.NewRecorder()
	env.VillageHandler.ListVillages(listRec, listReq)

	listData, _ := decodeEnvelope(t, listRec)
	var villages []struct {
		ID int64 `json:"id"`
	}
	json.Unmarshal(listData, &villages)
	if len(villages) == 0 {
		t.Fatal("no villages found")
	}

	// Now get that specific village — need to use path value
	req := httptest.NewRequest("GET", "/api/villages/"+strconv.Itoa(int(villages[0].ID)), nil)
	req.SetPathValue("id", strconv.Itoa(int(villages[0].ID)))
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.VillageHandler.GetVillage(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestVillageHandler_GetVillage_NotOwner(t *testing.T) {
	env := newTestEnv(t)

	player1, _ := registerAndLogin(t, env, "owner1", "owner1@test.com", "Strong@123")
	player2, _ := registerAndLogin(t, env, "owner2", "owner2@test.com", "Strong@123")
	chooseKingdomForPlayer(t, env, player1, "veridor")

	// List player1's villages
	listReq := httptest.NewRequest("GET", "/api/villages", nil)
	listReq = listReq.WithContext(authCtx(player1, "player"))
	listRec := httptest.NewRecorder()
	env.VillageHandler.ListVillages(listRec, listReq)

	listData, _ := decodeEnvelope(t, listRec)
	var villages []struct {
		ID int64 `json:"id"`
	}
	json.Unmarshal(listData, &villages)

	// Player2 tries to access player1's village
	req := httptest.NewRequest("GET", "/api/villages/"+strconv.Itoa(int(villages[0].ID)), nil)
	req.SetPathValue("id", strconv.Itoa(int(villages[0].ID)))
	req = req.WithContext(authCtx(player2, "player"))
	rec := httptest.NewRecorder()
	env.VillageHandler.GetVillage(rec, req)

	if rec.Code != http.StatusForbidden {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusForbidden)
	}
}

func TestVillageHandler_GetVillage_NotFound(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "notfoundvil", "notfound@test.com", "Strong@123")

	req := httptest.NewRequest("GET", "/api/villages/9999", nil)
	req.SetPathValue("id", "9999")
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.VillageHandler.GetVillage(rec, req)

	if rec.Code != http.StatusNotFound {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusNotFound)
	}
}

// --- RenameVillage ---

func TestVillageHandler_RenameVillage_Success(t *testing.T) {
	env := newTestEnv(t)

	playerID, _ := registerAndLogin(t, env, "renameuser", "rename@test.com", "Strong@123")
	chooseKingdomForPlayer(t, env, playerID, "arkazia")

	// Get village ID
	listReq := httptest.NewRequest("GET", "/api/villages", nil)
	listReq = listReq.WithContext(authCtx(playerID, "player"))
	listRec := httptest.NewRecorder()
	env.VillageHandler.ListVillages(listRec, listReq)

	listData, _ := decodeEnvelope(t, listRec)
	var villages []struct {
		ID int64 `json:"id"`
	}
	json.Unmarshal(listData, &villages)

	// Rename
	body := `{"name":"New Village Name"}`
	req := httptest.NewRequest("PUT", "/api/villages/"+strconv.Itoa(int(villages[0].ID))+"/name", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", strconv.Itoa(int(villages[0].ID)))
	req = req.WithContext(authCtx(playerID, "player"))
	rec := httptest.NewRecorder()
	env.VillageHandler.RenameVillage(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

package handler_test

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"

	"github.com/luisfpires18/woo/internal/dto"
)

// setupArkaziaPlayer registers an Arkazia player, chooses kingdom, and levels barracks to lv1.
// Returns playerID, villageID.
func setupArkaziaPlayer(t *testing.T, env *testEnv) (int64, int64) {
	t.Helper()

	playerID, _ := registerAndLogin(t, env, "gladtest", "glad@test.com", "Strong@123")
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
	if len(villages) == 0 {
		t.Fatal("no villages")
	}
	villageID := villages[0].ID

	// Level up barracks to 1 via direct DB update
	_, err := env.DB.ExecContext(context.Background(),
		`UPDATE buildings SET level = 1 WHERE village_id = ? AND building_type = 'barracks'`, villageID)
	if err != nil {
		t.Fatalf("level up barracks: %v", err)
	}

	return playerID, villageID
}

// --- StartTraining ---

func TestTrainingHandler_StartTraining_Success(t *testing.T) {
	env := newTestEnv(t)
	playerID, villageID := setupArkaziaPlayer(t, env)

	body := `{"troop_type":"iron_legionary","quantity":2}`
	req := httptest.NewRequest("POST", "/api/villages/"+strconv.FormatInt(villageID, 10)+"/train", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authCtx(playerID, "player"))
	req.SetPathValue("id", strconv.FormatInt(villageID, 10))
	rec := httptest.NewRecorder()

	env.TrainingHandler.StartTraining(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d, body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	data, _ := decodeEnvelope(t, rec)
	var resp dto.TrainingQueueResponse
	json.Unmarshal(data, &resp)

	if resp.TroopType != "iron_legionary" {
		t.Errorf("troop_type: got %q, want %q", resp.TroopType, "iron_legionary")
	}
	if resp.Quantity != 2 {
		t.Errorf("quantity: got %d, want 2", resp.Quantity)
	}
}

func TestTrainingHandler_StartTraining_UnknownTroop(t *testing.T) {
	env := newTestEnv(t)
	playerID, villageID := setupArkaziaPlayer(t, env)

	body := `{"troop_type":"dragon_rider","quantity":1}`
	req := httptest.NewRequest("POST", "/api/villages/"+strconv.FormatInt(villageID, 10)+"/train", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req = req.WithContext(authCtx(playerID, "player"))
	req.SetPathValue("id", strconv.FormatInt(villageID, 10))
	rec := httptest.NewRecorder()

	env.TrainingHandler.StartTraining(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

func TestTrainingHandler_StartTraining_Unauthenticated(t *testing.T) {
	env := newTestEnv(t)

	body := `{"troop_type":"iron_legionary","quantity":1}`
	req := httptest.NewRequest("POST", "/api/villages/1/train", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	req.SetPathValue("id", "1")
	rec := httptest.NewRecorder()

	env.TrainingHandler.StartTraining(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- GetTrainingCost ---

func TestTrainingHandler_GetTrainingCost_Success(t *testing.T) {
	env := newTestEnv(t)
	playerID, villageID := setupArkaziaPlayer(t, env)

	req := httptest.NewRequest("GET", "/api/villages/"+strconv.FormatInt(villageID, 10)+"/train/cost?troop_type=iron_legionary&quantity=3", nil)
	req = req.WithContext(authCtx(playerID, "player"))
	req.SetPathValue("id", strconv.FormatInt(villageID, 10))
	rec := httptest.NewRecorder()

	env.TrainingHandler.GetTrainingCost(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d, body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	data, _ := decodeEnvelope(t, rec)
	var resp dto.TrainingCostResponse
	json.Unmarshal(data, &resp)

	// Iron legionary: Food:100, Water:50, Lumber:60, Stone:40 × 3
	if resp.TotalFood != 300 {
		t.Errorf("total_food: got %.0f, want 300", resp.TotalFood)
	}
	if resp.Quantity != 3 {
		t.Errorf("quantity: got %d, want 3", resp.Quantity)
	}
}

// --- CancelTraining ---

func TestTrainingHandler_CancelTraining_Success(t *testing.T) {
	env := newTestEnv(t)
	playerID, villageID := setupArkaziaPlayer(t, env)

	// Start training first
	body := `{"troop_type":"iron_legionary","quantity":1}`
	startReq := httptest.NewRequest("POST", "/api/villages/"+strconv.FormatInt(villageID, 10)+"/train", strings.NewReader(body))
	startReq.Header.Set("Content-Type", "application/json")
	startReq = startReq.WithContext(authCtx(playerID, "player"))
	startReq.SetPathValue("id", strconv.FormatInt(villageID, 10))
	startRec := httptest.NewRecorder()
	env.TrainingHandler.StartTraining(startRec, startReq)

	data, _ := decodeEnvelope(t, startRec)
	var resp dto.TrainingQueueResponse
	json.Unmarshal(data, &resp)

	// Now cancel
	queueIDStr := strconv.FormatInt(resp.ID, 10)
	cancelReq := httptest.NewRequest("DELETE", "/api/villages/"+strconv.FormatInt(villageID, 10)+"/train/"+queueIDStr, nil)
	cancelReq = cancelReq.WithContext(authCtx(playerID, "player"))
	cancelReq.SetPathValue("id", strconv.FormatInt(villageID, 10))
	cancelReq.SetPathValue("queueId", queueIDStr)
	cancelRec := httptest.NewRecorder()

	env.TrainingHandler.CancelTraining(cancelRec, cancelReq)

	if cancelRec.Code != http.StatusNoContent {
		t.Errorf("status: got %d, want %d, body: %s", cancelRec.Code, http.StatusNoContent, cancelRec.Body.String())
	}
}

// --- ListTroops ---

func TestTrainingHandler_ListTroops_Empty(t *testing.T) {
	env := newTestEnv(t)
	playerID, villageID := setupArkaziaPlayer(t, env)

	req := httptest.NewRequest("GET", "/api/villages/"+strconv.FormatInt(villageID, 10)+"/troops", nil)
	req = req.WithContext(authCtx(playerID, "player"))
	req.SetPathValue("id", strconv.FormatInt(villageID, 10))
	rec := httptest.NewRecorder()

	env.TrainingHandler.ListTroops(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}

	data, _ := decodeEnvelope(t, rec)
	var result struct {
		Troops []dto.TroopInfo `json:"troops"`
	}
	json.Unmarshal(data, &result)

	if len(result.Troops) != 0 {
		t.Errorf("expected 0 troops, got %d", len(result.Troops))
	}
}

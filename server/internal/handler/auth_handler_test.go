package handler_test

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/luisfpires18/woo/internal/dto"
)

// --- Register ---

func TestAuthHandler_Register_Success(t *testing.T) {
	env := newTestEnv(t)

	body := `{"username":"newplayer","email":"new@test.com","password":"Strong@123"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()

	env.AuthHandler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusCreated, rec.Body.String())
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	var resp dto.AuthResponse
	if err := json.Unmarshal(data, &resp); err != nil {
		t.Fatalf("decode AuthResponse: %v", err)
	}

	if resp.AccessToken == "" {
		t.Error("expected non-empty access_token")
	}
	if resp.RefreshToken == "" {
		t.Error("expected non-empty refresh_token")
	}
	if resp.Player.Username != "newplayer" {
		t.Errorf("username: got %q, want %q", resp.Player.Username, "newplayer")
	}
}

func TestAuthHandler_Register_DuplicateUsername(t *testing.T) {
	env := newTestEnv(t)

	body := `{"username":"dup","email":"dup@test.com","password":"Strong@123"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.AuthHandler.Register(rec, req)

	if rec.Code != http.StatusCreated {
		t.Fatalf("first register failed: %d %s", rec.Code, rec.Body.String())
	}

	// Second registration with same username
	body2 := `{"username":"dup","email":"dup2@test.com","password":"Strong@123"}`
	req2 := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body2))
	req2.Header.Set("Content-Type", "application/json")
	rec2 := httptest.NewRecorder()
	env.AuthHandler.Register(rec2, req2)

	if rec2.Code != http.StatusConflict {
		t.Errorf("status: got %d, want %d", rec2.Code, http.StatusConflict)
	}
}

func TestAuthHandler_Register_WeakPassword(t *testing.T) {
	env := newTestEnv(t)

	body := `{"username":"weakpw","email":"weak@test.com","password":"123"}`
	req := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.AuthHandler.Register(rec, req)

	if rec.Code != http.StatusBadRequest {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusBadRequest)
	}
}

// --- Login ---

func TestAuthHandler_Login_Success(t *testing.T) {
	env := newTestEnv(t)

	// Register first
	regBody := `{"username":"logintest","email":"login@test.com","password":"Strong@123"}`
	regReq := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	env.AuthHandler.Register(regRec, regReq)

	// Login with email
	body := `{"login":"login@test.com","password":"Strong@123"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.AuthHandler.Login(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}

	data, errMsg := decodeEnvelope(t, rec)
	if errMsg != "" {
		t.Fatalf("unexpected error: %s", errMsg)
	}

	var resp dto.AuthResponse
	json.Unmarshal(data, &resp)
	if resp.AccessToken == "" {
		t.Error("expected non-empty access_token")
	}
}

func TestAuthHandler_Login_WrongPassword(t *testing.T) {
	env := newTestEnv(t)

	// Register
	regBody := `{"username":"wrongpw","email":"wrongpw@test.com","password":"Strong@123"}`
	regReq := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	env.AuthHandler.Register(regRec, regReq)

	// Login with wrong password
	body := `{"login":"wrongpw@test.com","password":"WrongPass@1"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.AuthHandler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

func TestAuthHandler_Login_NonexistentUser(t *testing.T) {
	env := newTestEnv(t)

	body := `{"login":"nobody@test.com","password":"Strong@123"}`
	req := httptest.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.AuthHandler.Login(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- Refresh ---

func TestAuthHandler_Refresh_Success(t *testing.T) {
	env := newTestEnv(t)

	// Register to get a refresh token
	regBody := `{"username":"refreshtest","email":"refresh@test.com","password":"Strong@123"}`
	regReq := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	env.AuthHandler.Register(regRec, regReq)

	data, _ := decodeEnvelope(t, regRec)
	var authResp dto.AuthResponse
	json.Unmarshal(data, &authResp)

	// Use refresh token
	body := `{"refresh_token":"` + authResp.RefreshToken + `"}`
	req := httptest.NewRequest("POST", "/api/auth/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.AuthHandler.Refresh(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d. Body: %s", rec.Code, http.StatusOK, rec.Body.String())
	}
}

func TestAuthHandler_Refresh_InvalidToken(t *testing.T) {
	env := newTestEnv(t)

	body := `{"refresh_token":"invalid-token-here"}`
	req := httptest.NewRequest("POST", "/api/auth/refresh", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.AuthHandler.Refresh(rec, req)

	if rec.Code != http.StatusUnauthorized {
		t.Errorf("status: got %d, want %d", rec.Code, http.StatusUnauthorized)
	}
}

// --- Logout ---

func TestAuthHandler_Logout_Success(t *testing.T) {
	env := newTestEnv(t)

	// Register to get tokens
	regBody := `{"username":"logouttest","email":"logout@test.com","password":"Strong@123"}`
	regReq := httptest.NewRequest("POST", "/api/auth/register", strings.NewReader(regBody))
	regReq.Header.Set("Content-Type", "application/json")
	regRec := httptest.NewRecorder()
	env.AuthHandler.Register(regRec, regReq)

	data, _ := decodeEnvelope(t, regRec)
	var authResp dto.AuthResponse
	json.Unmarshal(data, &authResp)

	// Logout
	body := `{"refresh_token":"` + authResp.RefreshToken + `"}`
	req := httptest.NewRequest("POST", "/api/auth/logout", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	rec := httptest.NewRecorder()
	env.AuthHandler.Logout(rec, req)

	if rec.Code != http.StatusOK {
		t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
	}
}

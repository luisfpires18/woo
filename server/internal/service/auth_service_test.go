package service_test

import (
	"context"
	"testing"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/repository/sqlite"
	"github.com/luisfpires18/woo/internal/service"
	"github.com/luisfpires18/woo/internal/testutil"
)

func newTestAuthService(t *testing.T) *service.AuthService {
	t.Helper()
	db := testutil.NewTestDB(t)
	playerRepo := sqlite.NewPlayerRepo(db)
	refreshTokenRepo := sqlite.NewRefreshTokenRepo(db)
	return service.NewAuthService(playerRepo, refreshTokenRepo, "test-secret-key", "woo-test")
}

func TestRegister_Success(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	resp, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "testplayer",
		Email:    "test@example.com",
		Password: "securepass123",
		Kingdom:  "veridor",
	})
	if err != nil {
		t.Fatalf("expected no error, got: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("expected access token")
	}
	if resp.RefreshToken == "" {
		t.Error("expected refresh token")
	}
	if resp.Player.Username != "testplayer" {
		t.Errorf("expected username 'testplayer', got %q", resp.Player.Username)
	}
	if resp.Player.Kingdom != "veridor" {
		t.Errorf("expected kingdom 'veridor', got %q", resp.Player.Kingdom)
	}
	if resp.Player.Role != "player" {
		t.Errorf("expected role 'player', got %q", resp.Player.Role)
	}
}

func TestRegister_DuplicateEmail(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	req := &dto.RegisterRequest{
		Username: "player1",
		Email:    "dupe@example.com",
		Password: "securepass123",
		Kingdom:  "sylvara",
	}
	if _, err := svc.Register(ctx, req); err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	req.Username = "player2"
	_, err := svc.Register(ctx, req)
	if err != service.ErrEmailTaken {
		t.Errorf("expected ErrEmailTaken, got: %v", err)
	}
}

func TestRegister_DuplicateUsername(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	req := &dto.RegisterRequest{
		Username: "samename",
		Email:    "first@example.com",
		Password: "securepass123",
		Kingdom:  "arkazia",
	}
	if _, err := svc.Register(ctx, req); err != nil {
		t.Fatalf("first register failed: %v", err)
	}

	req.Email = "second@example.com"
	_, err := svc.Register(ctx, req)
	if err != service.ErrUsernameTaken {
		t.Errorf("expected ErrUsernameTaken, got: %v", err)
	}
}

func TestRegister_InvalidKingdom(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "player",
		Email:    "test@example.com",
		Password: "securepass123",
		Kingdom:  "mordor",
	})
	if err != service.ErrInvalidKingdom {
		t.Errorf("expected ErrInvalidKingdom, got: %v", err)
	}
}

func TestRegister_WeakPassword(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "player",
		Email:    "test@example.com",
		Password: "short",
		Kingdom:  "veridor",
	})
	if err != service.ErrWeakPassword {
		t.Errorf("expected ErrWeakPassword, got: %v", err)
	}
}

func TestRegister_InvalidEmail(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "player",
		Email:    "notanemail",
		Password: "securepass123",
		Kingdom:  "veridor",
	})
	if err != service.ErrInvalidEmail {
		t.Errorf("expected ErrInvalidEmail, got: %v", err)
	}
}

func TestRegister_InvalidUsername(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "ab",
		Email:    "test@example.com",
		Password: "securepass123",
		Kingdom:  "veridor",
	})
	if err != service.ErrInvalidUsername {
		t.Errorf("expected ErrInvalidUsername, got: %v", err)
	}
}

func TestLogin_Success(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	// Register first
	_, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "logintest",
		Email:    "login@example.com",
		Password: "securepass123",
		Kingdom:  "sylvara",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Login with email
	resp, err := svc.Login(ctx, &dto.LoginRequest{
		Login:    "login@example.com",
		Password: "securepass123",
	})
	if err != nil {
		t.Fatalf("login with email failed: %v", err)
	}
	if resp.AccessToken == "" {
		t.Error("expected access token")
	}
	if resp.Player.Username != "logintest" {
		t.Errorf("expected username 'logintest', got %q", resp.Player.Username)
	}

	// Login with username
	resp2, err := svc.Login(ctx, &dto.LoginRequest{
		Login:    "logintest",
		Password: "securepass123",
	})
	if err != nil {
		t.Fatalf("login with username failed: %v", err)
	}
	if resp2.Player.Username != "logintest" {
		t.Errorf("expected username 'logintest', got %q", resp2.Player.Username)
	}
}

func TestLogin_WrongPassword(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "wrongpw",
		Email:    "wrong@example.com",
		Password: "correctpass1",
		Kingdom:  "arkazia",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	_, err = svc.Login(ctx, &dto.LoginRequest{
		Login:    "wrong@example.com",
		Password: "wrongpass00",
	})
	if err != service.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestLogin_NonexistentEmail(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.Login(ctx, &dto.LoginRequest{
		Login:    "noexist@example.com",
		Password: "securepass123",
	})
	if err != service.ErrInvalidCredentials {
		t.Errorf("expected ErrInvalidCredentials, got: %v", err)
	}
}

func TestRefreshToken_Success(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	// Register to get tokens
	regResp, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "refreshtest",
		Email:    "refresh@example.com",
		Password: "securepass123",
		Kingdom:  "veridor",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Use refresh token to get new tokens
	newResp, err := svc.RefreshToken(ctx, regResp.RefreshToken)
	if err != nil {
		t.Fatalf("refresh failed: %v", err)
	}
	if newResp.AccessToken == "" {
		t.Error("expected new access token")
	}
	if newResp.RefreshToken == regResp.RefreshToken {
		t.Error("expected rotated refresh token")
	}
}

func TestRefreshToken_Invalid(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	_, err := svc.RefreshToken(ctx, "invalid-token-string")
	if err != service.ErrInvalidRefreshToken {
		t.Errorf("expected ErrInvalidRefreshToken, got: %v", err)
	}
}

func TestRefreshToken_SingleUse(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	regResp, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "singleuse",
		Email:    "single@example.com",
		Password: "securepass123",
		Kingdom:  "sylvara",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Use refresh token once
	if _, err := svc.RefreshToken(ctx, regResp.RefreshToken); err != nil {
		t.Fatalf("first refresh failed: %v", err)
	}

	// Reuse same refresh token — should fail (single-use rotation)
	_, err = svc.RefreshToken(ctx, regResp.RefreshToken)
	if err != service.ErrInvalidRefreshToken {
		t.Errorf("expected ErrInvalidRefreshToken on reuse, got: %v", err)
	}
}

func TestValidateAccessToken(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	resp, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "validatetest",
		Email:    "validate@example.com",
		Password: "securepass123",
		Kingdom:  "arkazia",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	playerID, role, err := svc.ValidateAccessToken(resp.AccessToken)
	if err != nil {
		t.Fatalf("validate failed: %v", err)
	}
	if playerID != resp.Player.ID {
		t.Errorf("expected player ID %d, got %d", resp.Player.ID, playerID)
	}
	if role != "player" {
		t.Errorf("expected role 'player', got %q", role)
	}
}

func TestValidateAccessToken_Invalid(t *testing.T) {
	svc := newTestAuthService(t)

	_, _, err := svc.ValidateAccessToken("not.a.valid.jwt")
	if err == nil {
		t.Error("expected error for invalid token")
	}
}

func TestLogout(t *testing.T) {
	svc := newTestAuthService(t)
	ctx := context.Background()

	resp, err := svc.Register(ctx, &dto.RegisterRequest{
		Username: "logouttest",
		Email:    "logout@example.com",
		Password: "securepass123",
		Kingdom:  "veridor",
	})
	if err != nil {
		t.Fatalf("register failed: %v", err)
	}

	// Logout should invalidate the refresh token
	if err := svc.Logout(ctx, resp.RefreshToken); err != nil {
		t.Fatalf("logout failed: %v", err)
	}

	// Refresh with the same token should now fail
	_, err = svc.RefreshToken(ctx, resp.RefreshToken)
	if err != service.ErrInvalidRefreshToken {
		t.Errorf("expected ErrInvalidRefreshToken after logout, got: %v", err)
	}
}

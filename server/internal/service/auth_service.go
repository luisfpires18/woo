package service

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Auth errors.
var (
	ErrInvalidCredentials  = errors.New("invalid username/email or password")
	ErrEmailTaken          = errors.New("email already registered")
	ErrUsernameTaken       = errors.New("username already taken")
	ErrInvalidKingdom      = errors.New("invalid kingdom")
	ErrInvalidRefreshToken = errors.New("invalid or expired refresh token")
	ErrWeakPassword        = errors.New("password must be at least 8 characters")
	ErrInvalidEmail        = errors.New("invalid email address")
	ErrInvalidUsername     = errors.New("username must be 3-20 characters, alphanumeric and underscores only")
)

var validKingdoms = map[string]bool{
	"veridor":   true,
	"sylvara":   true,
	"arkazia":   true,
	"draxys":    true,
	"zandres":   true,
	"lumus":     true,
	"nordalh":   true,
	"drakanith": true,
}

// IsValidKingdom checks if the given kingdom string is a valid playable kingdom.
func IsValidKingdom(kingdom string) bool {
	return validKingdoms[kingdom]
}

const (
	accessTokenDuration  = 15 * time.Minute
	refreshTokenDuration = 7 * 24 * time.Hour
	bcryptCost           = 12
)

// AuthService handles authentication business logic.
type AuthService struct {
	playerRepo       repository.PlayerRepository
	refreshTokenRepo repository.RefreshTokenRepository
	jwtSecret        []byte
	jwtIssuer        string
}

// NewAuthService creates a new AuthService.
func NewAuthService(
	playerRepo repository.PlayerRepository,
	refreshTokenRepo repository.RefreshTokenRepository,
	jwtSecret string,
	jwtIssuer string,
) *AuthService {
	return &AuthService{
		playerRepo:       playerRepo,
		refreshTokenRepo: refreshTokenRepo,
		jwtSecret:        []byte(jwtSecret),
		jwtIssuer:        jwtIssuer,
	}
}

// Register creates a new player account and returns auth tokens.
func (s *AuthService) Register(ctx context.Context, req *dto.RegisterRequest) (*dto.AuthResponse, error) {
	// Validate input
	if err := s.validateRegisterInput(req); err != nil {
		return nil, err
	}

	// Check if email is already taken
	existing, err := s.playerRepo.GetByEmail(ctx, req.Email)
	if err != nil && !errors.Is(err, model.ErrNotFound) {
		return nil, fmt.Errorf("check email: %w", err)
	}
	if existing != nil {
		return nil, ErrEmailTaken
	}

	// Hash password
	hash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcryptCost)
	if err != nil {
		return nil, fmt.Errorf("hash password: %w", err)
	}

	// Create player — kingdom is empty until chosen post-login
	player := &model.Player{
		Username:     req.Username,
		Email:        strings.ToLower(strings.TrimSpace(req.Email)),
		PasswordHash: string(hash),
		Kingdom:      "",
		CreatedAt:    time.Now().UTC(),
	}

	if err := s.playerRepo.Create(ctx, player); err != nil {
		if strings.Contains(err.Error(), "UNIQUE constraint failed") {
			if strings.Contains(err.Error(), "username") {
				return nil, ErrUsernameTaken
			}
			return nil, ErrEmailTaken
		}
		return nil, fmt.Errorf("create player: %w", err)
	}

	// Generate tokens
	return s.generateAuthResponse(ctx, player)
}

// Login authenticates a player with email or username and password.
func (s *AuthService) Login(ctx context.Context, req *dto.LoginRequest) (*dto.AuthResponse, error) {
	login := strings.TrimSpace(req.Login)

	var player *model.Player
	var err error

	// If the login contains '@', treat it as an email; otherwise as a username.
	if strings.Contains(login, "@") {
		player, err = s.playerRepo.GetByEmail(ctx, strings.ToLower(login))
	} else {
		player, err = s.playerRepo.GetByUsername(ctx, login)
	}

	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrInvalidCredentials
		}
		return nil, fmt.Errorf("get player: %w", err)
	}

	if player.PasswordHash == "" {
		return nil, ErrInvalidCredentials // OAuth-only account
	}

	if err := bcrypt.CompareHashAndPassword([]byte(player.PasswordHash), []byte(req.Password)); err != nil {
		return nil, ErrInvalidCredentials
	}

	// Update last login
	_ = s.playerRepo.UpdateLastLogin(ctx, player.ID)

	return s.generateAuthResponse(ctx, player)
}

// RefreshToken validates a refresh token and issues new tokens.
func (s *AuthService) RefreshToken(ctx context.Context, rawToken string) (*dto.AuthResponse, error) {
	tokenHash := hashToken(rawToken)

	stored, err := s.refreshTokenRepo.GetByTokenHash(ctx, tokenHash)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrInvalidRefreshToken
		}
		return nil, fmt.Errorf("get refresh token: %w", err)
	}

	if time.Now().UTC().After(stored.ExpiresAt) {
		// Token expired — clean it up
		_ = s.refreshTokenRepo.DeleteByTokenHash(ctx, tokenHash)
		return nil, ErrInvalidRefreshToken
	}

	// Delete old token (single use)
	_ = s.refreshTokenRepo.DeleteByTokenHash(ctx, tokenHash)

	// Get player
	player, err := s.playerRepo.GetByID(ctx, stored.PlayerID)
	if err != nil {
		return nil, fmt.Errorf("get player for refresh: %w", err)
	}

	return s.generateAuthResponse(ctx, player)
}

// Logout invalidates a refresh token.
func (s *AuthService) Logout(ctx context.Context, rawToken string) error {
	tokenHash := hashToken(rawToken)
	return s.refreshTokenRepo.DeleteByTokenHash(ctx, tokenHash)
}

// ValidateAccessToken parses and validates a JWT access token, returning the player ID and role.
func (s *AuthService) ValidateAccessToken(tokenString string) (int64, string, error) {
	token, err := jwt.Parse(tokenString, func(t *jwt.Token) (any, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
		}
		return s.jwtSecret, nil
	})
	if err != nil {
		return 0, "", fmt.Errorf("parse token: %w", err)
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok || !token.Valid {
		return 0, "", errors.New("invalid token claims")
	}

	playerIDFloat, ok := claims["sub"].(float64)
	if !ok {
		return 0, "", errors.New("invalid subject claim")
	}

	role, _ := claims["role"].(string)
	if role == "" {
		role = "player"
	}

	return int64(playerIDFloat), role, nil
}

// generateAuthResponse creates JWT access token + refresh token for a player.
func (s *AuthService) generateAuthResponse(ctx context.Context, player *model.Player) (*dto.AuthResponse, error) {
	// Generate JWT access token
	now := time.Now().UTC()
	accessClaims := jwt.MapClaims{
		"sub":      player.ID,
		"username": player.Username,
		"kingdom":  player.Kingdom,
		"role":     player.Role,
		"iss":      s.jwtIssuer,
		"iat":      now.Unix(),
		"exp":      now.Add(accessTokenDuration).Unix(),
	}
	accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessClaims)
	accessTokenString, err := accessToken.SignedString(s.jwtSecret)
	if err != nil {
		return nil, fmt.Errorf("sign access token: %w", err)
	}

	// Generate opaque refresh token
	rawRefreshToken, err := generateRandomToken(32)
	if err != nil {
		return nil, fmt.Errorf("generate refresh token: %w", err)
	}

	// Store hashed refresh token
	refreshToken := &model.RefreshToken{
		PlayerID:  player.ID,
		TokenHash: hashToken(rawRefreshToken),
		ExpiresAt: now.Add(refreshTokenDuration),
		CreatedAt: now,
	}
	if err := s.refreshTokenRepo.Create(ctx, refreshToken); err != nil {
		return nil, fmt.Errorf("store refresh token: %w", err)
	}

	return &dto.AuthResponse{
		AccessToken:  accessTokenString,
		RefreshToken: rawRefreshToken,
		Player: &dto.PlayerInfo{
			ID:       player.ID,
			Username: player.Username,
			Email:    player.Email,
			Kingdom:  player.Kingdom,
			Role:     player.Role,
		},
	}, nil
}

func (s *AuthService) validateRegisterInput(req *dto.RegisterRequest) error {
	// Username: 3-20 chars, alphanumeric + underscore
	username := strings.TrimSpace(req.Username)
	if len(username) < 3 || len(username) > 20 {
		return ErrInvalidUsername
	}
	for _, c := range username {
		if !((c >= 'a' && c <= 'z') || (c >= 'A' && c <= 'Z') || (c >= '0' && c <= '9') || c == '_') {
			return ErrInvalidUsername
		}
	}

	// Email: basic check
	email := strings.TrimSpace(req.Email)
	if !strings.Contains(email, "@") || !strings.Contains(email, ".") || len(email) < 5 {
		return ErrInvalidEmail
	}

	// Password: min 8 chars
	if len(req.Password) < 8 {
		return ErrWeakPassword
	}

	// Kingdom is now chosen post-registration; skip validation here
	return nil
}

// generateRandomToken generates a cryptographically random hex token.
func generateRandomToken(bytes int) (string, error) {
	b := make([]byte, bytes)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	return hex.EncodeToString(b), nil
}

// hashToken returns the SHA-256 hash of a token string.
func hashToken(token string) string {
	h := sha256.Sum256([]byte(token))
	return hex.EncodeToString(h[:])
}

package service

import (
	"context"
	"fmt"
	"log/slog"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/repository"
)

var (
	// ErrKingdomAlreadyChosen is returned when a player tries to change their kingdom.
	ErrKingdomAlreadyChosen = fmt.Errorf("kingdom already chosen")
)

// PlayerService handles player business logic.
type PlayerService struct {
	playerRepo     repository.PlayerRepository
	villageService *VillageService
}

// NewPlayerService creates a new PlayerService.
func NewPlayerService(
	playerRepo repository.PlayerRepository,
	villageService *VillageService,
) *PlayerService {
	return &PlayerService{
		playerRepo:     playerRepo,
		villageService: villageService,
	}
}

// GetMe returns the authenticated player's info.
func (s *PlayerService) GetMe(ctx context.Context, playerID int64) (*dto.PlayerInfo, error) {
	player, err := s.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("get player: %w", err)
	}

	return &dto.PlayerInfo{
		ID:       player.ID,
		Username: player.Username,
		Email:    player.Email,
		Kingdom:  player.Kingdom,
		Role:     player.Role,
	}, nil
}

// ChooseKingdom sets the player's kingdom (one-time) and creates their first village.
// Returns the player info and the new village ID.
func (s *PlayerService) ChooseKingdom(ctx context.Context, playerID int64, kingdom string) (*dto.PlayerInfo, int64, error) {
	if !IsValidKingdom(kingdom) {
		return nil, 0, ErrInvalidKingdom
	}

	player, err := s.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, 0, fmt.Errorf("get player: %w", err)
	}

	if player.Kingdom != "" {
		return nil, 0, ErrKingdomAlreadyChosen
	}

	// Create first village BEFORE committing the kingdom, so we can abort cleanly.
	village, err := s.villageService.CreateFirstVillage(ctx, playerID, kingdom, player.Username)
	if err != nil {
		slog.Error("failed to create first village", "player_id", playerID, "error", err)
		return nil, 0, fmt.Errorf("create first village: %w", err)
	}

	// Village created successfully — now persist the kingdom choice.
	if err := s.playerRepo.UpdateKingdom(ctx, playerID, kingdom); err != nil {
		return nil, 0, fmt.Errorf("update kingdom: %w", err)
	}

	info := &dto.PlayerInfo{
		ID:       player.ID,
		Username: player.Username,
		Email:    player.Email,
		Kingdom:  kingdom,
		Role:     player.Role,
	}

	return info, village.ID, nil
}

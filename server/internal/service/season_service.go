package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Season service errors.
var (
	ErrSeasonNotFound     = errors.New("season not found")
	ErrSeasonNotUpcoming  = errors.New("season is not in upcoming status")
	ErrSeasonNotActive    = errors.New("season is not active")
	ErrSeasonNotEnded     = errors.New("season is not in ended status")
	ErrAlreadyJoined      = errors.New("you have already joined this season")
	ErrSeasonNameRequired = errors.New("season name is required")
	ErrSeasonNameTaken    = errors.New("season name already exists")
)

// SeasonService handles season/server lifecycle and player enrollment.
type SeasonService struct {
	seasonRepo     repository.SeasonRepository
	playerRepo     repository.PlayerRepository
	villageService *VillageService
}

// NewSeasonService creates a new SeasonService.
func NewSeasonService(
	seasonRepo repository.SeasonRepository,
	playerRepo repository.PlayerRepository,
	villageService *VillageService,
) *SeasonService {
	return &SeasonService{
		seasonRepo:     seasonRepo,
		playerRepo:     playerRepo,
		villageService: villageService,
	}
}

// ── Admin operations ─────────────────────────────────────────────────────────

// CreateSeason creates a new season in "upcoming" status.
func (s *SeasonService) CreateSeason(ctx context.Context, req dto.CreateSeasonRequest) (*dto.SeasonResponse, error) {
	if req.Name == "" {
		return nil, ErrSeasonNameRequired
	}
	// Apply defaults
	if req.GameSpeed <= 0 {
		req.GameSpeed = 1.0
	}
	if req.ResourceMultiplier <= 0 {
		req.ResourceMultiplier = 1.0
	}
	if req.MaxVillagesPerPlayer <= 0 {
		req.MaxVillagesPerPlayer = 5
	}
	if req.WeaponsOfChaosCount <= 0 {
		req.WeaponsOfChaosCount = 7
	}
	if req.MapWidth <= 0 {
		req.MapWidth = 51
	}
	if req.MapHeight <= 0 {
		req.MapHeight = 51
	}

	season := &model.Season{
		Name:                 req.Name,
		Description:          req.Description,
		Status:               model.SeasonStatusUpcoming,
		StartDate:            req.StartDate,
		MapTemplateName:      req.MapTemplateName,
		GameSpeed:            req.GameSpeed,
		ResourceMultiplier:   req.ResourceMultiplier,
		MaxVillagesPerPlayer: req.MaxVillagesPerPlayer,
		WeaponsOfChaosCount:  req.WeaponsOfChaosCount,
		MapWidth:             req.MapWidth,
		MapHeight:            req.MapHeight,
	}

	if err := s.seasonRepo.Create(ctx, season); err != nil {
		return nil, fmt.Errorf("create season: %w", err)
	}

	slog.Info("season created", "id", season.ID, "name", season.Name)
	return s.toResponse(season, 0), nil
}

// UpdateSeason updates a season (only allowed while upcoming).
func (s *SeasonService) UpdateSeason(ctx context.Context, id int64, req dto.UpdateSeasonRequest) (*dto.SeasonResponse, error) {
	season, err := s.seasonRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrSeasonNotFound
		}
		return nil, fmt.Errorf("get season: %w", err)
	}

	if season.Status != model.SeasonStatusUpcoming {
		return nil, ErrSeasonNotUpcoming
	}

	// Apply partial updates
	if req.Name != nil {
		season.Name = *req.Name
	}
	if req.Description != nil {
		season.Description = *req.Description
	}
	if req.StartDate != nil {
		season.StartDate = req.StartDate
	}
	if req.MapTemplateName != nil {
		season.MapTemplateName = *req.MapTemplateName
	}
	if req.GameSpeed != nil {
		season.GameSpeed = *req.GameSpeed
	}
	if req.ResourceMultiplier != nil {
		season.ResourceMultiplier = *req.ResourceMultiplier
	}
	if req.MaxVillagesPerPlayer != nil {
		season.MaxVillagesPerPlayer = *req.MaxVillagesPerPlayer
	}
	if req.WeaponsOfChaosCount != nil {
		season.WeaponsOfChaosCount = *req.WeaponsOfChaosCount
	}
	if req.MapWidth != nil {
		season.MapWidth = *req.MapWidth
	}
	if req.MapHeight != nil {
		season.MapHeight = *req.MapHeight
	}

	if err := s.seasonRepo.Update(ctx, season); err != nil {
		return nil, fmt.Errorf("update season: %w", err)
	}

	count, _ := s.seasonRepo.GetSeasonPlayerCount(ctx, id)
	return s.toResponse(season, count), nil
}

// LaunchSeason transitions an "upcoming" season to "active".
func (s *SeasonService) LaunchSeason(ctx context.Context, id int64) (*dto.SeasonResponse, error) {
	season, err := s.seasonRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrSeasonNotFound
		}
		return nil, fmt.Errorf("get season: %w", err)
	}

	if season.Status != model.SeasonStatusUpcoming {
		return nil, ErrSeasonNotUpcoming
	}

	if err := s.seasonRepo.UpdateStatus(ctx, id, model.SeasonStatusActive); err != nil {
		return nil, fmt.Errorf("launch season: %w", err)
	}

	slog.Info("season launched", "id", id, "name", season.Name)

	// Re-fetch to get updated timestamps
	season, _ = s.seasonRepo.GetByID(ctx, id)
	count, _ := s.seasonRepo.GetSeasonPlayerCount(ctx, id)
	return s.toResponse(season, count), nil
}

// EndSeason transitions an "active" season to "ended".
func (s *SeasonService) EndSeason(ctx context.Context, id int64) (*dto.SeasonResponse, error) {
	season, err := s.seasonRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrSeasonNotFound
		}
		return nil, fmt.Errorf("get season: %w", err)
	}

	if season.Status != model.SeasonStatusActive {
		return nil, ErrSeasonNotActive
	}

	if err := s.seasonRepo.UpdateStatus(ctx, id, model.SeasonStatusEnded); err != nil {
		return nil, fmt.Errorf("end season: %w", err)
	}

	slog.Info("season ended", "id", id, "name", season.Name)
	season, _ = s.seasonRepo.GetByID(ctx, id)
	count, _ := s.seasonRepo.GetSeasonPlayerCount(ctx, id)
	return s.toResponse(season, count), nil
}

// ArchiveSeason transitions an "ended" season to "archived".
func (s *SeasonService) ArchiveSeason(ctx context.Context, id int64) (*dto.SeasonResponse, error) {
	season, err := s.seasonRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrSeasonNotFound
		}
		return nil, fmt.Errorf("get season: %w", err)
	}

	if season.Status != model.SeasonStatusEnded {
		return nil, ErrSeasonNotEnded
	}

	if err := s.seasonRepo.UpdateStatus(ctx, id, model.SeasonStatusArchived); err != nil {
		return nil, fmt.Errorf("archive season: %w", err)
	}

	slog.Info("season archived", "id", id, "name", season.Name)
	season, _ = s.seasonRepo.GetByID(ctx, id)
	count, _ := s.seasonRepo.GetSeasonPlayerCount(ctx, id)
	return s.toResponse(season, count), nil
}

// DeleteSeason deletes an "upcoming" season.
func (s *SeasonService) DeleteSeason(ctx context.Context, id int64) error {
	season, err := s.seasonRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return ErrSeasonNotFound
		}
		return fmt.Errorf("get season: %w", err)
	}

	if season.Status != model.SeasonStatusUpcoming {
		return ErrSeasonNotUpcoming
	}

	return s.seasonRepo.Delete(ctx, id)
}

// ── Player operations ────────────────────────────────────────────────────────

// ListSeasons returns all seasons visible to players (upcoming, active, ended).
func (s *SeasonService) ListSeasons(ctx context.Context, statusFilter string) ([]*dto.SeasonResponse, error) {
	seasons, err := s.seasonRepo.List(ctx, statusFilter)
	if err != nil {
		return nil, fmt.Errorf("list seasons: %w", err)
	}

	results := make([]*dto.SeasonResponse, 0, len(seasons))
	for _, season := range seasons {
		count, _ := s.seasonRepo.GetSeasonPlayerCount(ctx, season.ID)
		results = append(results, s.toResponse(season, count))
	}
	return results, nil
}

// GetSeason returns season detail with the player's join status.
func (s *SeasonService) GetSeason(ctx context.Context, id, playerID int64) (*dto.SeasonDetailResponse, error) {
	season, err := s.seasonRepo.GetByID(ctx, id)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrSeasonNotFound
		}
		return nil, fmt.Errorf("get season: %w", err)
	}

	count, _ := s.seasonRepo.GetSeasonPlayerCount(ctx, id)
	resp := s.toResponse(season, count)

	detail := &dto.SeasonDetailResponse{
		SeasonResponse: *resp,
	}

	// Check if player has joined
	sp, err := s.seasonRepo.GetSeasonPlayer(ctx, id, playerID)
	if err == nil && sp != nil {
		detail.Joined = true
		detail.Kingdom = sp.Kingdom
	}

	return detail, nil
}

// JoinSeason adds a player to an active season with their chosen kingdom,
// and creates their first village in that season.
func (s *SeasonService) JoinSeason(ctx context.Context, seasonID, playerID int64, kingdom string) (*dto.SeasonDetailResponse, int64, error) {
	if !IsValidKingdom(kingdom) {
		return nil, 0, ErrInvalidKingdom
	}

	season, err := s.seasonRepo.GetByID(ctx, seasonID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, 0, ErrSeasonNotFound
		}
		return nil, 0, fmt.Errorf("get season: %w", err)
	}

	if season.Status != model.SeasonStatusActive {
		return nil, 0, ErrSeasonNotActive
	}

	// Check if already joined
	_, err = s.seasonRepo.GetSeasonPlayer(ctx, seasonID, playerID)
	if err == nil {
		return nil, 0, ErrAlreadyJoined
	}
	if !errors.Is(err, model.ErrNotFound) {
		return nil, 0, fmt.Errorf("check season player: %w", err)
	}

	// Get player info for village naming
	player, err := s.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, 0, fmt.Errorf("get player: %w", err)
	}

	// Add to season
	sp := &model.SeasonPlayer{
		SeasonID: seasonID,
		PlayerID: playerID,
		Kingdom:  kingdom,
	}
	if err := s.seasonRepo.AddPlayer(ctx, sp); err != nil {
		return nil, 0, fmt.Errorf("add season player: %w", err)
	}

	// Update player's kingdom on the global players table if not set
	// (first season they join sets their "default" kingdom)
	if player.Kingdom == "" {
		_ = s.playerRepo.UpdateKingdom(ctx, playerID, kingdom)
	}

	// Create first village for this season
	village, err := s.villageService.CreateFirstVillageForSeason(ctx, playerID, kingdom, player.Username, seasonID)
	if err != nil {
		slog.Error("failed to create season village", "player_id", playerID, "season_id", seasonID, "error", err)
		return nil, 0, fmt.Errorf("create season village: %w", err)
	}

	count, _ := s.seasonRepo.GetSeasonPlayerCount(ctx, seasonID)
	resp := s.toResponse(season, count)
	detail := &dto.SeasonDetailResponse{
		SeasonResponse: *resp,
		Joined:         true,
		Kingdom:        kingdom,
	}

	return detail, village.ID, nil
}

// GetMySeasons returns all seasons the player is participating in.
func (s *SeasonService) GetMySeasons(ctx context.Context, playerID int64) ([]*dto.SeasonDetailResponse, error) {
	seasonPlayers, err := s.seasonRepo.ListPlayerSeasons(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("list player seasons: %w", err)
	}

	results := make([]*dto.SeasonDetailResponse, 0, len(seasonPlayers))
	for _, sp := range seasonPlayers {
		season, err := s.seasonRepo.GetByID(ctx, sp.SeasonID)
		if err != nil {
			continue
		}
		count, _ := s.seasonRepo.GetSeasonPlayerCount(ctx, season.ID)
		resp := s.toResponse(season, count)
		results = append(results, &dto.SeasonDetailResponse{
			SeasonResponse: *resp,
			Joined:         true,
			Kingdom:        sp.Kingdom,
		})
	}
	return results, nil
}

// ── Profile ──────────────────────────────────────────────────────────────────

// GetPlayerProfile returns cross-season profile data.
func (s *SeasonService) GetPlayerProfile(ctx context.Context, playerID int64) (*dto.PlayerProfileResponse, error) {
	player, err := s.playerRepo.GetByID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("get player: %w", err)
	}

	history, err := s.seasonRepo.GetPlayerSeasonHistory(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("get season history: %w", err)
	}

	seasonEntries := make([]dto.SeasonHistoryEntry, 0, len(history))
	for _, h := range history {
		joinedAt, _ := time.Parse("2006-01-02 15:04:05", h.JoinedAt)
		seasonEntries = append(seasonEntries, dto.SeasonHistoryEntry{
			SeasonID:     h.SeasonID,
			SeasonName:   h.SeasonName,
			SeasonStatus: h.SeasonStatus,
			Kingdom:      h.Kingdom,
			JoinedAt:     joinedAt,
			VillageCount: h.VillageCount,
		})
	}

	return &dto.PlayerProfileResponse{
		ID:            player.ID,
		Username:      player.Username,
		Email:         player.Email,
		Role:          player.Role,
		CreatedAt:     player.CreatedAt,
		TotalSeasons:  len(seasonEntries),
		SeasonHistory: seasonEntries,
	}, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func (s *SeasonService) toResponse(season *model.Season, playerCount int) *dto.SeasonResponse {
	return &dto.SeasonResponse{
		ID:                   season.ID,
		Name:                 season.Name,
		Description:          season.Description,
		Status:               season.Status,
		StartDate:            season.StartDate,
		StartedAt:            season.StartedAt,
		EndedAt:              season.EndedAt,
		PlayerCount:          playerCount,
		MapTemplateName:      season.MapTemplateName,
		GameSpeed:            season.GameSpeed,
		ResourceMultiplier:   season.ResourceMultiplier,
		MaxVillagesPerPlayer: season.MaxVillagesPerPlayer,
		WeaponsOfChaosCount:  season.WeaponsOfChaosCount,
		MapWidth:             season.MapWidth,
		MapHeight:            season.MapHeight,
		CreatedAt:            season.CreatedAt,
		UpdatedAt:            season.UpdatedAt,
	}
}

package service

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Admin errors.
var (
	ErrInvalidRole = errors.New("role must be 'player' or 'admin'")
)

// AdminService handles admin-only business logic.
type AdminService struct {
	playerRepo       repository.PlayerRepository
	villageRepo      repository.VillageRepository
	worldConfigRepo  repository.WorldConfigRepository
	announcementRepo repository.AnnouncementRepository
}

// NewAdminService creates a new AdminService.
func NewAdminService(
	playerRepo repository.PlayerRepository,
	villageRepo repository.VillageRepository,
	worldConfigRepo repository.WorldConfigRepository,
	announcementRepo repository.AnnouncementRepository,
) *AdminService {
	return &AdminService{
		playerRepo:       playerRepo,
		villageRepo:      villageRepo,
		worldConfigRepo:  worldConfigRepo,
		announcementRepo: announcementRepo,
	}
}

// --- Player management ---

// ListPlayers returns a paginated list of all players.
func (s *AdminService) ListPlayers(ctx context.Context, offset, limit int) (*dto.PlayerListResponse, error) {
	if limit <= 0 {
		limit = 20
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	players, err := s.playerRepo.ListAll(ctx, offset, limit)
	if err != nil {
		return nil, fmt.Errorf("list players: %w", err)
	}

	total, err := s.playerRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("count players: %w", err)
	}

	items := make([]*dto.PlayerListItem, 0, len(players))
	for _, p := range players {
		items = append(items, &dto.PlayerListItem{
			ID:          p.ID,
			Username:    p.Username,
			Email:       p.Email,
			Kingdom:     p.Kingdom,
			Role:        p.Role,
			CreatedAt:   p.CreatedAt,
			LastLoginAt: p.LastLoginAt,
		})
	}

	return &dto.PlayerListResponse{
		Players: items,
		Total:   total,
		Offset:  offset,
		Limit:   limit,
	}, nil
}

// UpdatePlayerRole changes a player's role.
func (s *AdminService) UpdatePlayerRole(ctx context.Context, playerID int64, role string) error {
	if role != model.RolePlayer && role != model.RoleAdmin {
		return ErrInvalidRole
	}
	return s.playerRepo.UpdateRole(ctx, playerID, role)
}

// --- World config ---

// GetWorldConfig returns all configuration entries.
func (s *AdminService) GetWorldConfig(ctx context.Context) (*dto.WorldConfigResponse, error) {
	configs, err := s.worldConfigRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("get world config: %w", err)
	}

	entries := make([]*dto.WorldConfigEntry, 0, len(configs))
	for _, c := range configs {
		entries = append(entries, &dto.WorldConfigEntry{
			Key:         c.Key,
			Value:       c.Value,
			Description: c.Description,
			UpdatedAt:   c.UpdatedAt,
		})
	}

	return &dto.WorldConfigResponse{Configs: entries}, nil
}

// SetWorldConfig updates a single configuration value.
func (s *AdminService) SetWorldConfig(ctx context.Context, key, value string) error {
	if key == "" || value == "" {
		return errors.New("key and value are required")
	}
	return s.worldConfigRepo.Set(ctx, key, value)
}

// --- Server stats ---

// GetStats returns aggregate server statistics.
func (s *AdminService) GetStats(ctx context.Context) (*dto.StatsResponse, error) {
	playerCount, err := s.playerRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("count players: %w", err)
	}

	villageCount, err := s.villageRepo.Count(ctx)
	if err != nil {
		return nil, fmt.Errorf("count villages: %w", err)
	}

	return &dto.StatsResponse{
		TotalPlayers:  playerCount,
		TotalVillages: villageCount,
	}, nil
}

// --- Announcements ---

// CreateAnnouncement publishes a new server-wide announcement.
func (s *AdminService) CreateAnnouncement(ctx context.Context, req *dto.CreateAnnouncementRequest, authorID int64) (*dto.AnnouncementResponse, error) {
	if req.Title == "" || req.Content == "" {
		return nil, errors.New("title and content are required")
	}

	a := &model.Announcement{
		Title:    req.Title,
		Content:  req.Content,
		AuthorID: authorID,
	}

	if req.ExpiresAt != nil && *req.ExpiresAt != "" {
		t, err := time.Parse(time.RFC3339, *req.ExpiresAt)
		if err != nil {
			return nil, fmt.Errorf("invalid expires_at format (use RFC3339): %w", err)
		}
		utc := t.UTC()
		a.ExpiresAt = &utc
	}

	if err := s.announcementRepo.Create(ctx, a); err != nil {
		return nil, fmt.Errorf("create announcement: %w", err)
	}

	return &dto.AnnouncementResponse{
		ID:        a.ID,
		Title:     a.Title,
		Content:   a.Content,
		AuthorID:  a.AuthorID,
		CreatedAt: a.CreatedAt,
		ExpiresAt: a.ExpiresAt,
	}, nil
}

// ListAnnouncements returns all active (non-expired) announcements.
func (s *AdminService) ListAnnouncements(ctx context.Context) ([]*dto.AnnouncementResponse, error) {
	list, err := s.announcementRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list announcements: %w", err)
	}

	result := make([]*dto.AnnouncementResponse, 0, len(list))
	for _, a := range list {
		result = append(result, &dto.AnnouncementResponse{
			ID:        a.ID,
			Title:     a.Title,
			Content:   a.Content,
			AuthorID:  a.AuthorID,
			CreatedAt: a.CreatedAt,
			ExpiresAt: a.ExpiresAt,
		})
	}
	return result, nil
}

// DeleteAnnouncement removes an announcement by ID.
func (s *AdminService) DeleteAnnouncement(ctx context.Context, id int64) error {
	return s.announcementRepo.Delete(ctx, id)
}

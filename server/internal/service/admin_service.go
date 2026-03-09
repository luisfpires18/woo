package service

import (
	"context"
	"errors"
	"fmt"
	"strconv"
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
	playerRepo                repository.PlayerRepository
	villageRepo               repository.VillageRepository
	announcementRepo          repository.AnnouncementRepository
	gameAssetRepo             repository.GameAssetRepository
	resBuildingConfigRepo     repository.ResourceBuildingConfigRepository
	buildingDisplayConfigRepo repository.BuildingDisplayConfigRepository
}

// NewAdminService creates a new AdminService.
func NewAdminService(
	playerRepo repository.PlayerRepository,
	villageRepo repository.VillageRepository,
	announcementRepo repository.AnnouncementRepository,
	gameAssetRepo repository.GameAssetRepository,
	resBuildingConfigRepo repository.ResourceBuildingConfigRepository,
	buildingDisplayConfigRepo repository.BuildingDisplayConfigRepository,
) *AdminService {
	return &AdminService{
		playerRepo:                playerRepo,
		villageRepo:               villageRepo,
		announcementRepo:          announcementRepo,
		gameAssetRepo:             gameAssetRepo,
		resBuildingConfigRepo:     resBuildingConfigRepo,
		buildingDisplayConfigRepo: buildingDisplayConfigRepo,
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

// --- Game assets ---

// ListGameAssets returns all game assets.
func (s *AdminService) ListGameAssets(ctx context.Context) (*dto.GameAssetListResponse, error) {
	assets, err := s.gameAssetRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list game assets: %w", err)
	}

	items := make([]*dto.GameAssetDTO, 0, len(assets))
	for _, a := range assets {
		var spriteURL *string
		if a.SpritePath != nil {
			u := "/uploads/" + *a.SpritePath + "?v=" + strconv.FormatInt(a.UpdatedAt.Unix(), 10)
			spriteURL = &u
		}
		items = append(items, &dto.GameAssetDTO{
			ID:           a.ID,
			Category:     a.Category,
			DisplayName:  a.DisplayName,
			DefaultIcon:  a.DefaultIcon,
			SpriteURL:    spriteURL,
			SpriteWidth:  a.SpriteWidth,
			SpriteHeight: a.SpriteHeight,
			UpdatedAt:    a.UpdatedAt,
		})
	}

	return &dto.GameAssetListResponse{Assets: items}, nil
}

// GetGameAsset returns a single game asset by ID.
func (s *AdminService) GetGameAsset(ctx context.Context, id string) (*model.GameAsset, error) {
	return s.gameAssetRepo.GetByID(ctx, id)
}

// UpdateGameAssetSprite sets the sprite_path for a game asset.
func (s *AdminService) UpdateGameAssetSprite(ctx context.Context, id string, spritePath *string) error {
	return s.gameAssetRepo.UpdateSprite(ctx, id, spritePath)
}

// CreateGameAsset inserts a new game asset row (used for adding variants of zone/terrain tiles).
func (s *AdminService) CreateGameAsset(ctx context.Context, asset *model.GameAsset) error {
	if asset.ID == "" {
		return errors.New("id is required")
	}
	if asset.Category == "" {
		return errors.New("category is required")
	}
	if asset.DisplayName == "" {
		return errors.New("display_name is required")
	}
	// Validate category is known
	if _, ok := model.AssetSpriteDimensions[asset.Category]; !ok {
		return fmt.Errorf("unknown category: %s", asset.Category)
	}
	// Set default dimensions from category
	dims := model.AssetSpriteDimensions[asset.Category]
	asset.SpriteWidth = dims[0]
	asset.SpriteHeight = dims[1]
	if asset.DefaultIcon == "" {
		asset.DefaultIcon = "🖼️"
	}
	return s.gameAssetRepo.Create(ctx, asset)
}

// DeleteGameAsset removes a game asset row by ID (used for removing variants).
func (s *AdminService) DeleteGameAsset(ctx context.Context, id string) error {
	return s.gameAssetRepo.Delete(ctx, id)
}

// --- Resource building configs ---

// ListResourceBuildingConfigs returns all resource building configs, optionally filtered by kingdom.
func (s *AdminService) ListResourceBuildingConfigs(ctx context.Context, kingdom string) (*dto.ResourceBuildingConfigListResponse, error) {
	var configs []*model.ResourceBuildingConfig
	var err error
	if kingdom != "" {
		configs, err = s.resBuildingConfigRepo.GetByKingdom(ctx, kingdom)
	} else {
		configs, err = s.resBuildingConfigRepo.GetAll(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("list resource building configs: %w", err)
	}

	items := make([]*dto.ResourceBuildingConfigDTO, 0, len(configs))
	for _, c := range configs {
		var spriteURL *string
		if c.SpritePath != nil {
			u := "/uploads/" + *c.SpritePath + "?v=" + strconv.FormatInt(c.UpdatedAt.Unix(), 10)
			spriteURL = &u
		}
		items = append(items, &dto.ResourceBuildingConfigDTO{
			ID:           c.ID,
			ResourceType: c.ResourceType,
			Slot:         c.Slot,
			Kingdom:      c.Kingdom,
			DisplayName:  c.DisplayName,
			Description:  c.Description,
			DefaultIcon:  c.DefaultIcon,
			SpriteURL:    spriteURL,
			UpdatedAt:    c.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &dto.ResourceBuildingConfigListResponse{Configs: items}, nil
}

// GetResourceBuildingConfig returns a single resource building config by ID.
func (s *AdminService) GetResourceBuildingConfig(ctx context.Context, id int64) (*model.ResourceBuildingConfig, error) {
	return s.resBuildingConfigRepo.GetByID(ctx, id)
}

// UpdateResourceBuildingConfig updates display_name, description, and default_icon for a config.
func (s *AdminService) UpdateResourceBuildingConfig(ctx context.Context, id int64, req *dto.UpdateResourceBuildingConfigRequest) error {
	if req.DisplayName == "" {
		return errors.New("display_name is required")
	}
	cfg, err := s.resBuildingConfigRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get resource building config: %w", err)
	}
	cfg.DisplayName = req.DisplayName
	cfg.Description = req.Description
	if req.DefaultIcon != "" {
		cfg.DefaultIcon = req.DefaultIcon
	}
	return s.resBuildingConfigRepo.Update(ctx, cfg)
}

// UpdateResourceBuildingConfigSprite sets the sprite_path for a resource building config.
func (s *AdminService) UpdateResourceBuildingConfigSprite(ctx context.Context, id int64, spritePath *string) error {
	return s.resBuildingConfigRepo.UpdateSprite(ctx, id, spritePath)
}

// --- Building display configs ---

// ListBuildingDisplayConfigs returns all building display configs, optionally filtered by kingdom.
func (s *AdminService) ListBuildingDisplayConfigs(ctx context.Context, kingdom string) (*dto.BuildingDisplayConfigListResponse, error) {
	var configs []*model.BuildingDisplayConfig
	var err error
	if kingdom != "" {
		configs, err = s.buildingDisplayConfigRepo.GetByKingdom(ctx, kingdom)
	} else {
		configs, err = s.buildingDisplayConfigRepo.GetAll(ctx)
	}
	if err != nil {
		return nil, fmt.Errorf("list building display configs: %w", err)
	}

	items := make([]*dto.BuildingDisplayConfigDTO, 0, len(configs))
	for _, c := range configs {
		var spriteURL *string
		if c.SpritePath != nil {
			u := "/uploads/" + *c.SpritePath + "?v=" + strconv.FormatInt(c.UpdatedAt.Unix(), 10)
			spriteURL = &u
		}
		items = append(items, &dto.BuildingDisplayConfigDTO{
			ID:           c.ID,
			BuildingType: c.BuildingType,
			Kingdom:      c.Kingdom,
			DisplayName:  c.DisplayName,
			Description:  c.Description,
			DefaultIcon:  c.DefaultIcon,
			SpriteURL:    spriteURL,
			UpdatedAt:    c.UpdatedAt.Format(time.RFC3339),
		})
	}

	return &dto.BuildingDisplayConfigListResponse{Configs: items}, nil
}

// GetBuildingDisplayConfig returns a single building display config by ID.
func (s *AdminService) GetBuildingDisplayConfig(ctx context.Context, id int64) (*model.BuildingDisplayConfig, error) {
	return s.buildingDisplayConfigRepo.GetByID(ctx, id)
}

// UpdateBuildingDisplayConfig updates display_name, description, and default_icon for a config.
func (s *AdminService) UpdateBuildingDisplayConfig(ctx context.Context, id int64, req *dto.UpdateBuildingDisplayConfigRequest) error {
	if req.DisplayName == "" {
		return errors.New("display_name is required")
	}
	cfg, err := s.buildingDisplayConfigRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get building display config: %w", err)
	}
	cfg.DisplayName = req.DisplayName
	cfg.Description = req.Description
	if req.DefaultIcon != "" {
		cfg.DefaultIcon = req.DefaultIcon
	}
	return s.buildingDisplayConfigRepo.Update(ctx, cfg)
}

// UpdateBuildingDisplayConfigSprite sets the sprite_path for a building display config.
func (s *AdminService) UpdateBuildingDisplayConfigSprite(ctx context.Context, id int64, spritePath *string) error {
	return s.buildingDisplayConfigRepo.UpdateSprite(ctx, id, spritePath)
}

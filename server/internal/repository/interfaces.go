package repository

import (
	"context"
	"time"

	"github.com/luisfpires18/woo/internal/model"
)

// PlayerRepository defines data access operations for players.
type PlayerRepository interface {
	Create(ctx context.Context, player *model.Player) error
	GetByID(ctx context.Context, id int64) (*model.Player, error)
	GetByEmail(ctx context.Context, email string) (*model.Player, error)
	GetByUsername(ctx context.Context, username string) (*model.Player, error)
	GetByOAuth(ctx context.Context, provider, oauthID string) (*model.Player, error)
	UpdateLastLogin(ctx context.Context, id int64) error
	UpdateRole(ctx context.Context, id int64, role string) error
	ListAll(ctx context.Context, offset, limit int) ([]*model.Player, error)
	Count(ctx context.Context) (int64, error)
}

// VillageRepository defines data access operations for villages.
type VillageRepository interface {
	Create(ctx context.Context, village *model.Village) error
	GetByID(ctx context.Context, id int64) (*model.Village, error)
	ListByPlayerID(ctx context.Context, playerID int64) ([]*model.Village, error)
	Update(ctx context.Context, village *model.Village) error
	GetByCoordinates(ctx context.Context, x, y int) (*model.Village, error)
	Count(ctx context.Context) (int64, error)
}

// BuildingRepository defines data access operations for buildings.
type BuildingRepository interface {
	Create(ctx context.Context, building *model.Building) error
	CreateBatch(ctx context.Context, buildings []*model.Building) error
	GetByVillageID(ctx context.Context, villageID int64) ([]*model.Building, error)
	Update(ctx context.Context, building *model.Building) error
}

// ResourceRepository defines data access operations for village resources.
type ResourceRepository interface {
	Get(ctx context.Context, villageID int64) (*model.Resources, error)
	Update(ctx context.Context, villageID int64, resources *model.Resources) error
	Create(ctx context.Context, resources *model.Resources) error
}

// RefreshTokenRepository defines data access operations for refresh tokens.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteAllByPlayerID(ctx context.Context, playerID int64) error
}

// WorldConfigRepository defines data access operations for world configuration.
type WorldConfigRepository interface {
	Get(ctx context.Context, key string) (*model.WorldConfig, error)
	GetAll(ctx context.Context) ([]*model.WorldConfig, error)
	Set(ctx context.Context, key, value string) error
}

// AnnouncementRepository defines data access operations for announcements.
type AnnouncementRepository interface {
	Create(ctx context.Context, announcement *model.Announcement) error
	ListActive(ctx context.Context) ([]*model.Announcement, error)
	Delete(ctx context.Context, id int64) error
}

// BuildingQueueRepository defines data access operations for the building construction queue.
type BuildingQueueRepository interface {
	Insert(ctx context.Context, item *model.BuildingQueue) error
	GetByVillageID(ctx context.Context, villageID int64) ([]*model.BuildingQueue, error)
	GetCompleted(ctx context.Context, now time.Time) ([]*model.BuildingQueue, error)
	Delete(ctx context.Context, id int64) error
}

// GameAssetRepository defines data access operations for game assets (buildings, resources, units).
type GameAssetRepository interface {
	GetAll(ctx context.Context) ([]*model.GameAsset, error)
	GetByID(ctx context.Context, id string) (*model.GameAsset, error)
	GetByCategory(ctx context.Context, category string) ([]*model.GameAsset, error)
	UpdateSprite(ctx context.Context, id string, spritePath *string) error
	Create(ctx context.Context, asset *model.GameAsset) error
}

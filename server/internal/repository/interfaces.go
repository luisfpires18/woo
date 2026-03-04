package repository

import (
	"context"

	"github.com/luisfpires18/woo/internal/model"
)

// PlayerRepository defines data access operations for players.
type PlayerRepository interface {
	Create(ctx context.Context, player *model.Player) error
	GetByID(ctx context.Context, id int64) (*model.Player, error)
	GetByEmail(ctx context.Context, email string) (*model.Player, error)
	GetByOAuth(ctx context.Context, provider, oauthID string) (*model.Player, error)
	UpdateLastLogin(ctx context.Context, id int64) error
}

// VillageRepository defines data access operations for villages.
type VillageRepository interface {
	Create(ctx context.Context, village *model.Village) error
	GetByID(ctx context.Context, id int64) (*model.Village, error)
	ListByPlayerID(ctx context.Context, playerID int64) ([]*model.Village, error)
	Update(ctx context.Context, village *model.Village) error
	GetByCoordinates(ctx context.Context, x, y int) (*model.Village, error)
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

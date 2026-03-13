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
	UpdateKingdom(ctx context.Context, id int64, kingdom string) error
	ListAll(ctx context.Context, offset, limit int) ([]*model.Player, error)
	Count(ctx context.Context) (int64, error)
}

// VillageRepository defines data access operations for villages.
type VillageRepository interface {
	Create(ctx context.Context, village *model.Village) error
	GetByID(ctx context.Context, id int64) (*model.Village, error)
	ListByPlayerID(ctx context.Context, playerID int64) ([]*model.Village, error)
	ListByPlayerAndSeason(ctx context.Context, playerID int64, seasonID int64) ([]*model.Village, error)
	Update(ctx context.Context, village *model.Village) error
	GetByCoordinates(ctx context.Context, x, y int) (*model.Village, error)
	Count(ctx context.Context) (int64, error)
	CountByPlayerAndSeason(ctx context.Context, playerID int64, seasonID int64) (int64, error)
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

// PlayerEconomyRepository defines data access operations for per-player economy (gold).
type PlayerEconomyRepository interface {
	Create(ctx context.Context, playerID int64, gold float64) error
	GetByPlayerID(ctx context.Context, playerID int64) (*model.PlayerEconomy, error)
	UpdateGold(ctx context.Context, playerID int64, newGold float64) error
	// DeductGold atomically deducts gold if sufficient balance exists. Returns ErrInsufficientGold otherwise.
	DeductGold(ctx context.Context, playerID int64, amount float64) error
	// DeductGoldTx atomically deducts gold within an existing transaction.
	DeductGoldTx(ctx context.Context, tx interface{}, playerID int64, amount float64) error
}

// RefreshTokenRepository defines data access operations for refresh tokens.
type RefreshTokenRepository interface {
	Create(ctx context.Context, token *model.RefreshToken) error
	GetByTokenHash(ctx context.Context, tokenHash string) (*model.RefreshToken, error)
	DeleteByTokenHash(ctx context.Context, tokenHash string) error
	DeleteAllByPlayerID(ctx context.Context, playerID int64) error
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
	GetByID(ctx context.Context, id int64) (*model.BuildingQueue, error)
	GetByVillageID(ctx context.Context, villageID int64) ([]*model.BuildingQueue, error)
	GetCompleted(ctx context.Context, now time.Time) ([]*model.BuildingQueue, error)
	Update(ctx context.Context, item *model.BuildingQueue) error
	Delete(ctx context.Context, id int64) error
}

// GameAssetRepository defines data access operations for game assets (buildings, resources, units).
type GameAssetRepository interface {
	GetAll(ctx context.Context) ([]*model.GameAsset, error)
	GetByID(ctx context.Context, id string) (*model.GameAsset, error)
	GetByCategory(ctx context.Context, category string) ([]*model.GameAsset, error)
	Create(ctx context.Context, asset *model.GameAsset) error
	Delete(ctx context.Context, id string) error
}

// ResourceBuildingConfigRepository defines data access operations for per-kingdom resource building cosmetics.
type ResourceBuildingConfigRepository interface {
	GetAll(ctx context.Context) ([]*model.ResourceBuildingConfig, error)
	GetByKingdom(ctx context.Context, kingdom string) ([]*model.ResourceBuildingConfig, error)
	GetByID(ctx context.Context, id int64) (*model.ResourceBuildingConfig, error)
	Update(ctx context.Context, cfg *model.ResourceBuildingConfig) error
}

// BuildingDisplayConfigRepository defines data access operations for per-kingdom village building cosmetics.
type BuildingDisplayConfigRepository interface {
	GetAll(ctx context.Context) ([]*model.BuildingDisplayConfig, error)
	GetByKingdom(ctx context.Context, kingdom string) ([]*model.BuildingDisplayConfig, error)
	GetByID(ctx context.Context, id int64) (*model.BuildingDisplayConfig, error)
	Update(ctx context.Context, cfg *model.BuildingDisplayConfig) error
}

// TroopDisplayConfigRepository defines data access operations for per-kingdom troop cosmetics.
type TroopDisplayConfigRepository interface {
	GetAll(ctx context.Context) ([]*model.TroopDisplayConfig, error)
	GetByKingdom(ctx context.Context, kingdom string) ([]*model.TroopDisplayConfig, error)
	GetByID(ctx context.Context, id int64) (*model.TroopDisplayConfig, error)
	Update(ctx context.Context, cfg *model.TroopDisplayConfig) error
}

// WorldMapRepository defines data access operations for the world map tile grid.
type WorldMapRepository interface {
	// InsertBatch inserts multiple map tiles in a single transaction.
	InsertBatch(ctx context.Context, tiles []*model.MapTile) error
	// DeleteAll removes all tiles from the world map.
	DeleteAll(ctx context.Context) error
	// GetChunk returns tiles within a rectangular region centered on (cx, cy) with the given radius.
	GetChunk(ctx context.Context, cx, cy, radius int) ([]*model.MapTile, error)
	// GetTile returns a single tile at the given coordinates.
	GetTile(ctx context.Context, x, y int) (*model.MapTile, error)
	// UpdateTileOwner sets the owner and village for a tile.
	UpdateTileOwner(ctx context.Context, x, y int, playerID, villageID *int64) error
	// Count returns the total number of tiles in the world map.
	Count(ctx context.Context) (int64, error)
	// GetByZone returns all tiles belonging to a specific kingdom zone.
	GetByZone(ctx context.Context, zone string) ([]*model.MapTile, error)
	// GetDistinctZones returns all distinct non-empty, non-wilderness kingdom zones currently placed.
	GetDistinctZones(ctx context.Context) ([]string, error)
	// UpdateTilesZone sets the kingdom_zone for all tiles within a circular radius of (cx, cy).
	UpdateTilesZone(ctx context.Context, cx, cy, radius int, zone string) error
	// UpdateTerrain updates the terrain_type for a batch of tiles.
	UpdateTerrain(ctx context.Context, tiles []model.TileTerrainUpdate) error
	// UpdateTilesZoneBatch updates the kingdom_zone for each tile individually (for template apply).
	UpdateTilesZoneBatch(ctx context.Context, tiles []model.TemplateTile) error
	// GetSpawnCandidates returns plains tiles with no village, optionally filtered by zone.
	// If zone is empty, returns candidates from the entire map.
	GetSpawnCandidates(ctx context.Context, zone string) ([]*model.MapTile, error)
}

// KingdomRelationRepository defines data access operations for kingdom diplomacy.
type KingdomRelationRepository interface {
	GetAll(ctx context.Context) ([]*model.KingdomRelation, error)
	Get(ctx context.Context, kingdomA, kingdomB string) (*model.KingdomRelation, error)
	Upsert(ctx context.Context, rel *model.KingdomRelation) error
}

// TroopRepository defines data access operations for troops stationed in villages.
type TroopRepository interface {
	GetByVillageID(ctx context.Context, villageID int64) ([]*model.Troop, error)
	GetByVillageAndType(ctx context.Context, villageID int64, troopType string) (*model.Troop, error)
	Upsert(ctx context.Context, troop *model.Troop) error
}

// TrainingQueueRepository defines data access operations for the troop training queue.
type TrainingQueueRepository interface {
	Insert(ctx context.Context, item *model.TrainingQueue) error
	GetByID(ctx context.Context, id int64) (*model.TrainingQueue, error)
	GetByVillageID(ctx context.Context, villageID int64) ([]*model.TrainingQueue, error)
	GetNextCompleted(ctx context.Context, now time.Time) ([]*model.TrainingQueue, error)
	Update(ctx context.Context, item *model.TrainingQueue) error
	Delete(ctx context.Context, id int64) error
}

// SeasonRepository defines data access operations for seasons.
type SeasonRepository interface {
	Create(ctx context.Context, season *model.Season) error
	GetByID(ctx context.Context, id int64) (*model.Season, error)
	List(ctx context.Context, statusFilter string) ([]*model.Season, error)
	Update(ctx context.Context, season *model.Season) error
	UpdateStatus(ctx context.Context, id int64, status string) error
	Delete(ctx context.Context, id int64) error

	// Season players
	AddPlayer(ctx context.Context, sp *model.SeasonPlayer) error
	RemovePlayer(ctx context.Context, seasonID, playerID int64) error
	GetSeasonPlayer(ctx context.Context, seasonID, playerID int64) (*model.SeasonPlayer, error)
	ListSeasonPlayers(ctx context.Context, seasonID int64) ([]*model.SeasonPlayer, error)
	ListPlayerSeasons(ctx context.Context, playerID int64) ([]*model.SeasonPlayer, error)
	GetSeasonPlayerCount(ctx context.Context, seasonID int64) (int, error)

	// Season player history with season info (for profile)
	GetPlayerSeasonHistory(ctx context.Context, playerID int64) ([]model.SeasonHistoryRow, error)
}

// BeastTemplateRepository defines data access operations for beast templates.
type BeastTemplateRepository interface {
	Create(ctx context.Context, bt *model.BeastTemplate) error
	GetByID(ctx context.Context, id int64) (*model.BeastTemplate, error)
	GetAll(ctx context.Context) ([]*model.BeastTemplate, error)
	Update(ctx context.Context, bt *model.BeastTemplate) error
	Delete(ctx context.Context, id int64) error
}

// CampTemplateRepository defines data access operations for camp templates.
type CampTemplateRepository interface {
	Create(ctx context.Context, ct *model.CampTemplate) error
	GetByID(ctx context.Context, id int64) (*model.CampTemplate, error)
	GetAll(ctx context.Context) ([]*model.CampTemplate, error)
	Update(ctx context.Context, ct *model.CampTemplate) error
	Delete(ctx context.Context, id int64) error
}

// CampBeastSlotRepository defines data access operations for camp beast slots.
type CampBeastSlotRepository interface {
	Create(ctx context.Context, slot *model.CampBeastSlot) error
	GetByCampTemplateID(ctx context.Context, campTemplateID int64) ([]*model.CampBeastSlot, error)
	Delete(ctx context.Context, id int64) error
	DeleteByCampTemplateID(ctx context.Context, campTemplateID int64) error
}

// SpawnRuleRepository defines data access operations for spawn rules.
type SpawnRuleRepository interface {
	Create(ctx context.Context, rule *model.SpawnRule) error
	GetByID(ctx context.Context, id int64) (*model.SpawnRule, error)
	GetAll(ctx context.Context) ([]*model.SpawnRule, error)
	GetEnabled(ctx context.Context) ([]*model.SpawnRule, error)
	Update(ctx context.Context, rule *model.SpawnRule) error
	Delete(ctx context.Context, id int64) error
}

// RewardTableRepository defines data access operations for reward tables.
type RewardTableRepository interface {
	Create(ctx context.Context, rt *model.RewardTable) error
	GetByID(ctx context.Context, id int64) (*model.RewardTable, error)
	GetAll(ctx context.Context) ([]*model.RewardTable, error)
	Delete(ctx context.Context, id int64) error
	Update(ctx context.Context, rt *model.RewardTable) error
}

// RewardTableEntryRepository defines data access operations for reward table entries.
type RewardTableEntryRepository interface {
	Create(ctx context.Context, entry *model.RewardTableEntry) error
	GetByRewardTableID(ctx context.Context, rewardTableID int64) ([]*model.RewardTableEntry, error)
	Delete(ctx context.Context, id int64) error
	DeleteByRewardTableID(ctx context.Context, rewardTableID int64) error
}

// CampRepository defines data access operations for runtime camp instances.
type CampRepository interface {
	Create(ctx context.Context, camp *model.Camp) error
	GetByID(ctx context.Context, id int64) (*model.Camp, error)
	GetByTile(ctx context.Context, x, y int) (*model.Camp, error)
	GetByStatus(ctx context.Context, status string) ([]*model.Camp, error)
	CountBySpawnRule(ctx context.Context, spawnRuleID int64) (int, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	Delete(ctx context.Context, id int64) error
	GetExpiredCamps(ctx context.Context, now time.Time) ([]*model.Camp, error)
	ListActive(ctx context.Context) ([]*model.Camp, error)
}

// ExpeditionRepository defines data access operations for expeditions.
type ExpeditionRepository interface {
	Create(ctx context.Context, exp *model.Expedition) error
	GetByID(ctx context.Context, id int64) (*model.Expedition, error)
	GetByPlayerID(ctx context.Context, playerID int64) ([]*model.Expedition, error)
	GetArrivedExpeditions(ctx context.Context, now time.Time) ([]*model.Expedition, error)
	GetReturningExpeditions(ctx context.Context, now time.Time) ([]*model.Expedition, error)
	UpdateStatus(ctx context.Context, id int64, status string) error
	UpdateReturnsAt(ctx context.Context, id int64, returnsAt time.Time) error
	Update(ctx context.Context, exp *model.Expedition) error
}

// BattleRepository defines data access operations for battle results.
type BattleRepository interface {
	Create(ctx context.Context, battle *model.Battle) error
	GetByID(ctx context.Context, id int64) (*model.Battle, error)
	GetByExpeditionID(ctx context.Context, expeditionID int64) (*model.Battle, error)
	GetReplayData(ctx context.Context, id int64) ([]byte, error)
}

// BattleTuningRepository defines data access operations for the battle tuning singleton.
type BattleTuningRepository interface {
	Get(ctx context.Context) (*model.BattleTuning, error)
	Update(ctx context.Context, tuning *model.BattleTuning) error
}

// AdminAuditLogRepository defines data access operations for the admin audit log.
type AdminAuditLogRepository interface {
	Create(ctx context.Context, entry *model.AdminAuditLog) error
	List(ctx context.Context, entityType string, limit, offset int) ([]*model.AdminAuditLog, error)
}

// UnitOfWork encapsulates multi-table transactional operations.
// Keeps database transaction details behind the repository abstraction boundary.
type UnitOfWork interface {
	// DeductResourcesAndInsertBuildQueue atomically deducts resources and inserts a build queue item.
	DeductResourcesAndInsertBuildQueue(ctx context.Context, villageID int64, res *model.Resources, item *model.BuildingQueue) error

	// DeductResourcesGoldAndInsertBuildQueue atomically deducts village resources + player gold + inserts a build queue item.
	DeductResourcesGoldAndInsertBuildQueue(ctx context.Context, villageID int64, res *model.Resources, playerID int64, goldCost float64, item *model.BuildingQueue) error

	// DeductResourcesAndInsertTrainQueue atomically deducts resources and inserts a training queue item.
	DeductResourcesAndInsertTrainQueue(ctx context.Context, villageID int64, res *model.Resources, item *model.TrainingQueue) error

	// DeductResourcesGoldAndInsertTrainQueue atomically deducts village resources + player gold + inserts a training queue item.
	DeductResourcesGoldAndInsertTrainQueue(ctx context.Context, villageID int64, res *model.Resources, playerID int64, goldCost float64, item *model.TrainingQueue) error

	// CompleteTrainingUnit atomically adds troops, updates resources (food consumption), and advances/deletes the queue item.
	CompleteTrainingUnit(ctx context.Context, villageID int64, troopType string, addQty int, res *model.Resources, queueItem *model.TrainingQueue, deleteQueue bool) error

	// CompleteBuildingUpgrade atomically updates building level, refreshes resource rates, and deletes the queue item.
	CompleteBuildingUpgrade(ctx context.Context, villageID int64, building *model.Building, resources *model.Resources, queueID int64) error

	// DeductTroopsAndCreateExpedition atomically deducts troops from a village and creates an expedition.
	DeductTroopsAndCreateExpedition(ctx context.Context, villageID int64, troopDeductions map[string]int, exp *model.Expedition) error

	// ReturnExpeditionTroops atomically adds surviving troops back and marks the expedition completed.
	ReturnExpeditionTroops(ctx context.Context, villageID int64, troopAdditions map[string]int, expeditionID int64) error

	// CreateVillageWithSetup atomically creates a village, links the map tile, creates starter buildings, resources, and player economy.
	CreateVillageWithSetup(ctx context.Context, village *model.Village, tileX, tileY int, buildings []*model.Building, resources *model.Resources, playerID int64, startingGold float64) error
}

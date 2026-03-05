package service

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/luisfpires18/woo/internal/config"
	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
	sqlt "github.com/luisfpires18/woo/internal/repository/sqlite"
)

// Building service errors.
var (
	ErrUnknownBuilding = errors.New("unknown building type")
)

// QueueTxInserter is an optional interface for queue repos that support transactional inserts.
type QueueTxInserter interface {
	InsertTx(ctx context.Context, tx *sql.Tx, item *model.BuildingQueue) error
}

// BuildingService handles building construction business logic.
type BuildingService struct {
	db            *sql.DB
	villageRepo   repository.VillageRepository
	buildingRepo  repository.BuildingRepository
	resourceRepo  repository.ResourceRepository
	queueRepo     repository.BuildingQueueRepository
	queueTx       QueueTxInserter
	playerRepo    repository.PlayerRepository
}

// NewBuildingService creates a new BuildingService.
func NewBuildingService(
	db *sql.DB,
	villageRepo repository.VillageRepository,
	buildingRepo repository.BuildingRepository,
	resourceRepo repository.ResourceRepository,
	queueRepo repository.BuildingQueueRepository,
	playerRepo repository.PlayerRepository,
) *BuildingService {
	svc := &BuildingService{
		db:           db,
		villageRepo:  villageRepo,
		buildingRepo: buildingRepo,
		resourceRepo: resourceRepo,
		queueRepo:    queueRepo,
		playerRepo:   playerRepo,
	}
	// If the queue repo supports transactional inserts, save a reference.
	if txRepo, ok := queueRepo.(QueueTxInserter); ok {
		svc.queueTx = txRepo
	}
	return svc
}

// StartUpgrade begins a building upgrade for the given village.
func (s *BuildingService) StartUpgrade(ctx context.Context, playerID, villageID int64, buildingType string) (*dto.BuildingQueueResponse, error) {
	// 1. Validate building type exists in config
	bldCfg, ok := config.BuildingConfigs[buildingType]
	if !ok {
		return nil, ErrUnknownBuilding
	}

	// 2. Validate village ownership
	village, err := s.villageRepo.GetByID(ctx, villageID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrVillageNotFound
		}
		return nil, fmt.Errorf("get village: %w", err)
	}
	if village.PlayerID != playerID {
		return nil, ErrNotOwner
	}

	// 3. Check kingdom restriction
	if bldCfg.KingdomOnly != "" {
		player, err := s.playerRepo.GetByID(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("get player: %w", err)
		}
		if player.Kingdom != bldCfg.KingdomOnly {
			return nil, fmt.Errorf("building %s is only available to %s kingdom", buildingType, bldCfg.KingdomOnly)
		}
	}

	// 4. Check no build already in progress
	queue, err := s.queueRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get build queue: %w", err)
	}
	if len(queue) > 0 {
		return nil, model.ErrBuildingInProgress
	}

	// 5. Get current buildings
	buildings, err := s.buildingRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get buildings: %w", err)
	}

	buildingMap := make(map[string]*model.Building, len(buildings))
	for _, b := range buildings {
		buildingMap[b.BuildingType] = b
	}

	// 6. Find the target building and check max level
	building, exists := buildingMap[buildingType]
	if !exists {
		return nil, fmt.Errorf("building type %s not found in this village", buildingType)
	}
	targetLevel := building.Level + 1
	if targetLevel > bldCfg.MaxLevel {
		return nil, model.ErrMaxLevelReached
	}

	// 7. Check prerequisites
	for _, prereq := range bldCfg.Prerequisites {
		prereqBuilding, ok := buildingMap[prereq.BuildingType]
		if !ok || prereqBuilding.Level < prereq.MinLevel {
			prereqCfg := config.BuildingConfigs[prereq.BuildingType]
			return nil, fmt.Errorf("%w: requires %s level %d", model.ErrPrerequisitesNotMet, prereqCfg.DisplayName, prereq.MinLevel)
		}
	}

	// 8. Calculate cost
	cost, err := config.CostAtLevel(buildingType, targetLevel)
	if err != nil {
		return nil, fmt.Errorf("calculate cost: %w", err)
	}
	timeSec, err := config.TimeAtLevel(buildingType, targetLevel)
	if err != nil {
		return nil, fmt.Errorf("calculate time: %w", err)
	}

	// 9. Flush lazy resources (calculate current amounts)
	res, err := s.resourceRepo.Get(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get resources: %w", err)
	}
	now := time.Now().UTC()
	elapsed := now.Sub(res.LastUpdated).Hours()
	if elapsed > 0 {
		res.Iron = clampStorage(res.Iron+res.IronRate*elapsed, res.MaxStorage)
		res.Wood = clampStorage(res.Wood+res.WoodRate*elapsed, res.MaxStorage)
		res.Stone = clampStorage(res.Stone+res.StoneRate*elapsed, res.MaxStorage)
		res.Food = clampStorage(res.Food+(res.FoodRate-res.FoodConsumption)*elapsed, res.MaxStorage)
		if res.Food < 0 {
			res.Food = 0
		}
	}

	// 10. Check sufficient resources
	if res.Iron < cost.Iron || res.Wood < cost.Wood || res.Stone < cost.Stone || res.Food < cost.Food {
		return nil, model.ErrInsufficientResources
	}

	// 11. Deduct resources + insert queue atomically
	res.Iron -= cost.Iron
	res.Wood -= cost.Wood
	res.Stone -= cost.Stone
	res.Food -= cost.Food
	res.LastUpdated = now

	completesAt := now.Add(time.Duration(timeSec) * time.Second)
	queueItem := &model.BuildingQueue{
		VillageID:    villageID,
		BuildingType: buildingType,
		TargetLevel:  targetLevel,
		StartedAt:    now,
		CompletesAt:  completesAt,
	}

	err = sqlt.WithTx(ctx, s.db, func(tx *sql.Tx) error {
		if err := sqlt.UpdateResourcesTx(ctx, tx, villageID,
			res.Iron, res.Wood, res.Stone, res.Food,
			res.IronRate, res.WoodRate, res.StoneRate, res.FoodRate,
			res.FoodConsumption, res.MaxStorage,
			res.LastUpdated.UTC().Format("2006-01-02 15:04:05"),
		); err != nil {
			return err
		}
		return s.queueTx.InsertTx(ctx, tx, queueItem)
	})
	if err != nil {
		return nil, fmt.Errorf("execute upgrade transaction: %w", err)
	}

	return &dto.BuildingQueueResponse{
		ID:           queueItem.ID,
		BuildingType: queueItem.BuildingType,
		TargetLevel:  queueItem.TargetLevel,
		StartedAt:    queueItem.StartedAt,
		CompletesAt:  queueItem.CompletesAt,
	}, nil
}

// GetUpgradeCost returns the cost for upgrading a building in the given village.
func (s *BuildingService) GetUpgradeCost(ctx context.Context, playerID, villageID int64, buildingType string) (*dto.BuildingCostResponse, error) {
	if _, ok := config.BuildingConfigs[buildingType]; !ok {
		return nil, ErrUnknownBuilding
	}

	village, err := s.villageRepo.GetByID(ctx, villageID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return nil, ErrVillageNotFound
		}
		return nil, fmt.Errorf("get village: %w", err)
	}
	if village.PlayerID != playerID {
		return nil, ErrNotOwner
	}

	buildings, err := s.buildingRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get buildings: %w", err)
	}

	var currentLevel int
	for _, b := range buildings {
		if b.BuildingType == buildingType {
			currentLevel = b.Level
			break
		}
	}

	targetLevel := currentLevel + 1
	cfg := config.BuildingConfigs[buildingType]
	if targetLevel > cfg.MaxLevel {
		return nil, model.ErrMaxLevelReached
	}

	cost, err := config.CostAtLevel(buildingType, targetLevel)
	if err != nil {
		return nil, fmt.Errorf("cost: %w", err)
	}
	timeSec, err := config.TimeAtLevel(buildingType, targetLevel)
	if err != nil {
		return nil, fmt.Errorf("time: %w", err)
	}

	return &dto.BuildingCostResponse{
		BuildingType: buildingType,
		CurrentLevel: currentLevel,
		TargetLevel:  targetLevel,
		Iron:         cost.Iron,
		Wood:         cost.Wood,
		Stone:        cost.Stone,
		Food:         cost.Food,
		TimeSec:      timeSec,
	}, nil
}

// CancelUpgrade removes a queued upgrade. No resources are refunded.
func (s *BuildingService) CancelUpgrade(ctx context.Context, playerID, villageID, queueID int64) error {
	village, err := s.villageRepo.GetByID(ctx, villageID)
	if err != nil {
		if errors.Is(err, model.ErrNotFound) {
			return ErrVillageNotFound
		}
		return fmt.Errorf("get village: %w", err)
	}
	if village.PlayerID != playerID {
		return ErrNotOwner
	}

	// Verify the queue item belongs to this village
	queue, err := s.queueRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return fmt.Errorf("get queue: %w", err)
	}
	found := false
	for _, q := range queue {
		if q.ID == queueID {
			found = true
			break
		}
	}
	if !found {
		return model.ErrNotFound
	}

	return s.queueRepo.Delete(ctx, queueID)
}

// CompleteBuilds processes all building queue items whose completes_at has passed.
// Called by the game tick loop.
func (s *BuildingService) CompleteBuilds(ctx context.Context) error {
	now := time.Now().UTC()
	completed, err := s.queueRepo.GetCompleted(ctx, now)
	if err != nil {
		return fmt.Errorf("get completed builds: %w", err)
	}

	for _, item := range completed {
		if err := s.completeBuild(ctx, item); err != nil {
			slog.Error("failed to complete build",
				"queue_id", item.ID,
				"village_id", item.VillageID,
				"building_type", item.BuildingType,
				"error", err,
			)
			// Continue processing other items even if one fails
			continue
		}
		slog.Info("building upgrade completed",
			"village_id", item.VillageID,
			"building_type", item.BuildingType,
			"level", item.TargetLevel,
		)
	}
	return nil
}

// completeBuild applies a single completed build: levels up the building,
// updates resource rates if applicable, and removes the queue entry.
func (s *BuildingService) completeBuild(ctx context.Context, item *model.BuildingQueue) error {
	// Get current buildings to find the one to upgrade
	buildings, err := s.buildingRepo.GetByVillageID(ctx, item.VillageID)
	if err != nil {
		return fmt.Errorf("get buildings: %w", err)
	}

	var targetBuilding *model.Building
	for _, b := range buildings {
		if b.BuildingType == item.BuildingType {
			targetBuilding = b
			break
		}
	}
	if targetBuilding == nil {
		return fmt.Errorf("building %s not found in village %d", item.BuildingType, item.VillageID)
	}

	// Level up
	targetBuilding.Level = item.TargetLevel
	if err := s.buildingRepo.Update(ctx, targetBuilding); err != nil {
		return fmt.Errorf("update building level: %w", err)
	}

	// Update resource rates if it's a resource-producing building
	if err := s.updateResourceRates(ctx, item.VillageID, item.BuildingType, item.TargetLevel); err != nil {
		return fmt.Errorf("update resource rates: %w", err)
	}

	// Remove completed queue entry
	if err := s.queueRepo.Delete(ctx, item.ID); err != nil {
		return fmt.Errorf("delete queue item: %w", err)
	}

	return nil
}

// updateResourceRates updates the village resource rates/storage when a resource building is upgraded.
func (s *BuildingService) updateResourceRates(ctx context.Context, villageID int64, buildingType string, newLevel int) error {
	res, err := s.resourceRepo.Get(ctx, villageID)
	if err != nil {
		return fmt.Errorf("get resources: %w", err)
	}

	changed := false
	newRate := config.BaseResourceRate + config.RatePerLevel*float64(newLevel)

	switch buildingType {
	case "iron_mine":
		res.IronRate = newRate
		changed = true
	case "lumber_mill":
		res.WoodRate = newRate
		changed = true
	case "quarry":
		res.StoneRate = newRate
		changed = true
	case "farm":
		res.FoodRate = newRate
		changed = true
	case "warehouse":
		res.MaxStorage = config.BaseStorage + config.StoragePerLevel*float64(newLevel)
		changed = true
	}

	if changed {
		if err := s.resourceRepo.Update(ctx, villageID, res); err != nil {
			return fmt.Errorf("update resource rates: %w", err)
		}
	}
	return nil
}

// GetBuildQueue returns the current build queue for a village.
func (s *BuildingService) GetBuildQueue(ctx context.Context, villageID int64) ([]dto.BuildingQueueResponse, error) {
	items, err := s.queueRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get build queue: %w", err)
	}

	result := make([]dto.BuildingQueueResponse, 0, len(items))
	for _, item := range items {
		result = append(result, dto.BuildingQueueResponse{
			ID:           item.ID,
			BuildingType: item.BuildingType,
			TargetLevel:  item.TargetLevel,
			StartedAt:    item.StartedAt,
			CompletesAt:  item.CompletesAt,
		})
	}
	return result, nil
}

func clampStorage(val, max float64) float64 {
	if val > max {
		return max
	}
	return val
}

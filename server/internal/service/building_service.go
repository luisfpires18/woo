package service

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"time"

	"github.com/luisfpires18/woo/internal/config"
	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Building service errors.
var (
	ErrUnknownBuilding = errors.New("unknown building type")
)

// BuildingService handles building construction business logic.
type BuildingService struct {
	uow          repository.UnitOfWork
	villageRepo  repository.VillageRepository
	buildingRepo repository.BuildingRepository
	resourceRepo repository.ResourceRepository
	queueRepo    repository.BuildingQueueRepository
	playerRepo   repository.PlayerRepository
}

// NewBuildingService creates a new BuildingService.
func NewBuildingService(
	uow repository.UnitOfWork,
	villageRepo repository.VillageRepository,
	buildingRepo repository.BuildingRepository,
	resourceRepo repository.ResourceRepository,
	queueRepo repository.BuildingQueueRepository,
	playerRepo repository.PlayerRepository,
) *BuildingService {
	return &BuildingService{
		uow:          uow,
		villageRepo:  villageRepo,
		buildingRepo: buildingRepo,
		resourceRepo: resourceRepo,
		queueRepo:    queueRepo,
		playerRepo:   playerRepo,
	}
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

	// 3. Check no build already in progress
	queue, err := s.queueRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get build queue: %w", err)
	}
	if len(queue) > 0 {
		return nil, model.ErrBuildingInProgress
	}

	// 4. Get current buildings
	buildings, err := s.buildingRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get buildings: %w", err)
	}

	buildingMap := make(map[string]*model.Building, len(buildings))
	for _, b := range buildings {
		buildingMap[b.BuildingType] = b
	}

	// 5. Find the target building and check max level
	building, exists := buildingMap[buildingType]
	if !exists {
		return nil, fmt.Errorf("building type %s not found in this village", buildingType)
	}
	targetLevel := building.Level + 1
	if targetLevel > bldCfg.MaxLevel {
		return nil, model.ErrMaxLevelReached
	}

	// 6. Check prerequisites
	for _, prereq := range bldCfg.Prerequisites {
		prereqBuilding, ok := buildingMap[prereq.BuildingType]
		if !ok || prereqBuilding.Level < prereq.MinLevel {
			prereqCfg := config.BuildingConfigs[prereq.BuildingType]
			return nil, fmt.Errorf("%w: requires %s level %d", model.ErrPrerequisitesNotMet, prereqCfg.DisplayName, prereq.MinLevel)
		}
	}

	// 7. Calculate cost
	cost, err := config.CostAtLevel(buildingType, targetLevel)
	if err != nil {
		return nil, fmt.Errorf("calculate cost: %w", err)
	}
	timeSec, err := config.TimeAtLevel(buildingType, targetLevel)
	if err != nil {
		return nil, fmt.Errorf("calculate time: %w", err)
	}

	// 8. Flush lazy resources (calculate current amounts)
	res, err := s.resourceRepo.Get(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get resources: %w", err)
	}
	now := time.Now().UTC()
	FlushResources(res, now)

	// 9. Check sufficient resources
	if res.Food < cost.Food || res.Water < cost.Water || res.Lumber < cost.Lumber || res.Stone < cost.Stone {
		return nil, model.ErrInsufficientResources
	}

	// 10. Deduct resources + insert queue atomically
	res.Food -= cost.Food
	res.Water -= cost.Water
	res.Lumber -= cost.Lumber
	res.Stone -= cost.Stone
	res.LastUpdated = now

	completesAt := now.Add(time.Duration(timeSec) * time.Second)
	queueItem := &model.BuildingQueue{
		VillageID:    villageID,
		BuildingType: buildingType,
		TargetLevel:  targetLevel,
		StartedAt:    now,
		CompletesAt:  completesAt,
	}

	err = s.uow.DeductResourcesAndInsertBuildQueue(ctx, villageID, res, queueItem)
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
		Food:         cost.Food,
		Water:        cost.Water,
		Lumber:       cost.Lumber,
		Stone:        cost.Stone,
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

// BuildCompletionEvent describes a single completed building upgrade.
type BuildCompletionEvent struct {
	PlayerID     int64
	VillageID    int64
	BuildingType string
	NewLevel     int
}

// CompleteBuilds processes all building queue items whose completes_at has passed.
// Returns the list of successfully completed events for notification purposes.
// Called by the game tick loop.
func (s *BuildingService) CompleteBuilds(ctx context.Context) ([]BuildCompletionEvent, error) {
	now := time.Now().UTC()
	completed, err := s.queueRepo.GetCompleted(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("get completed builds: %w", err)
	}

	var events []BuildCompletionEvent
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

		// Look up player owning this village for WS notification.
		village, err := s.villageRepo.GetByID(ctx, item.VillageID)
		if err != nil {
			slog.Warn("could not look up village owner for notification",
				"village_id", item.VillageID, "error", err)
		} else {
			events = append(events, BuildCompletionEvent{
				PlayerID:     village.PlayerID,
				VillageID:    item.VillageID,
				BuildingType: item.BuildingType,
				NewLevel:     item.TargetLevel,
			})
		}

		slog.Info("building upgrade completed",
			"village_id", item.VillageID,
			"building_type", item.BuildingType,
			"level", item.TargetLevel,
		)
	}
	return events, nil
}

// completeBuild applies a single completed build: atomically levels up the building,
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

	// Prepare building with new level
	targetBuilding.Level = item.TargetLevel

	// Get resources and recalculate rates based on the new building state
	res, err := s.resourceRepo.Get(ctx, item.VillageID)
	if err != nil {
		return fmt.Errorf("get resources: %w", err)
	}
	now := time.Now().UTC()
	FlushResources(res, now)
	res.LastUpdated = now

	// Update resource rates if it's a resource-producing building (inline calculation)
	resType := config.ResourceTypeForBuilding(item.BuildingType)
	if resType != "" {
		// Sum levels of all buildings (including the newly leveled one) that produce this resource type
		totalLevel := 0
		for _, b := range buildings {
			if config.ResourceTypeForBuilding(b.BuildingType) == resType {
				if b.BuildingType == item.BuildingType {
					totalLevel += targetBuilding.Level // Use the new level
				} else {
					totalLevel += b.Level
				}
			}
		}
		newRate := config.BaseResourceRate + config.RatePerLevel*float64(totalLevel)
		switch resType {
		case "food":
			res.FoodRate = newRate
		case "water":
			res.WaterRate = newRate
		case "lumber":
			res.LumberRate = newRate
		case "stone":
			res.StoneRate = newRate
		}
	}

	// Atomically: level up building, update resources, and remove queue entry
	if err := s.uow.CompleteBuildingUpgrade(ctx, item.VillageID, targetBuilding, res, item.ID); err != nil {
		return fmt.Errorf("complete building upgrade transaction: %w", err)
	}

	return nil
}

// updateResourceRates updates the village resource rates/storage when a building is upgraded.
// For resource buildings, it sums the levels of all 3 buildings producing the same resource type.
func (s *BuildingService) updateResourceRates(ctx context.Context, villageID int64, buildingType string, buildings []*model.Building) error {
	res, err := s.resourceRepo.Get(ctx, villageID)
	if err != nil {
		return fmt.Errorf("get resources: %w", err)
	}

	changed := false

	// Check if it's a resource-producing building
	resType := config.ResourceTypeForBuilding(buildingType)
	if resType != "" {
		// Sum levels of all buildings that produce this resource type
		totalLevel := 0
		for _, b := range buildings {
			if config.ResourceTypeForBuilding(b.BuildingType) == resType {
				totalLevel += b.Level
			}
		}
		newRate := config.BaseResourceRate + config.RatePerLevel*float64(totalLevel)
		switch resType {
		case "food":
			res.FoodRate = newRate
		case "water":
			res.WaterRate = newRate
		case "lumber":
			res.LumberRate = newRate
		case "stone":
			res.StoneRate = newRate
		}
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

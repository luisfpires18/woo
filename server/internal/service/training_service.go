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

// Training service errors.
var (
	ErrUnknownTroop        = errors.New("unknown troop type")
	ErrTrainingBuildingReq = errors.New("training building requirement not met")
	ErrInvalidQuantity     = errors.New("quantity must be at least 1")
	ErrInsufficientPop     = errors.New("not enough population capacity")
)

// TrainCompletionEvent describes a single unit completing training.
type TrainCompletionEvent struct {
	PlayerID  int64
	VillageID int64
	TroopType string
	NewTotal  int
}

// TrainingService handles troop training business logic.
type TrainingService struct {
	uow            repository.UnitOfWork
	villageRepo    repository.VillageRepository
	buildingRepo   repository.BuildingRepository
	resourceRepo   repository.ResourceRepository
	troopRepo      repository.TroopRepository
	queueRepo      repository.TrainingQueueRepository
	playerRepo     repository.PlayerRepository
	playerEconRepo repository.PlayerEconomyRepository
}

// NewTrainingService creates a new TrainingService.
func NewTrainingService(
	uow repository.UnitOfWork,
	villageRepo repository.VillageRepository,
	buildingRepo repository.BuildingRepository,
	resourceRepo repository.ResourceRepository,
	troopRepo repository.TroopRepository,
	queueRepo repository.TrainingQueueRepository,
	playerRepo repository.PlayerRepository,
	playerEconRepo repository.PlayerEconomyRepository,
) *TrainingService {
	return &TrainingService{
		uow:            uow,
		villageRepo:    villageRepo,
		buildingRepo:   buildingRepo,
		resourceRepo:   resourceRepo,
		troopRepo:      troopRepo,
		queueRepo:      queueRepo,
		playerRepo:     playerRepo,
		playerEconRepo: playerEconRepo,
	}
}

// StartTraining begins training a batch of troops.
func (s *TrainingService) StartTraining(ctx context.Context, playerID, villageID int64, troopType string, quantity int) (*dto.TrainingQueueResponse, error) {
	// 1. Validate troop type
	troopCfg, ok := config.TroopConfigs[troopType]
	if !ok {
		return nil, ErrUnknownTroop
	}

	// 2. Validate quantity
	if quantity < 1 {
		return nil, ErrInvalidQuantity
	}

	// 3. Validate village ownership
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

	// 4. Check kingdom restriction
	if troopCfg.Kingdom != "" {
		player, err := s.playerRepo.GetByID(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("get player: %w", err)
		}
		if player.Kingdom != troopCfg.Kingdom {
			return nil, fmt.Errorf("troop %s is only available to %s kingdom", troopType, troopCfg.Kingdom)
		}
	}

	// 5. Check training building exists and meets level requirement
	buildings, err := s.buildingRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get buildings: %w", err)
	}

	var trainingBuildingLevel int
	for _, b := range buildings {
		if b.BuildingType == troopCfg.TrainingBuilding {
			trainingBuildingLevel = b.Level
			break
		}
	}
	if trainingBuildingLevel < troopCfg.BuildingLevelReq {
		return nil, fmt.Errorf("%w: requires %s level %d (current: %d)",
			ErrTrainingBuildingReq, troopCfg.TrainingBuilding, troopCfg.BuildingLevelReq, trainingBuildingLevel)
	}

	// 6. Calculate per-unit training time with speed multiplier
	eachTimeSec, err := config.TrainingTime(troopType, trainingBuildingLevel)
	if err != nil {
		if errors.Is(err, config.ErrUnknownTroop) {
			return nil, ErrUnknownTroop
		}
		return nil, fmt.Errorf("calculate training time: %w", err)
	}

	// 7. Calculate total cost
	totalCost, err := config.TrainingCost(troopType, quantity)
	if err != nil {
		if errors.Is(err, config.ErrUnknownTroop) {
			return nil, ErrUnknownTroop
		}
		return nil, fmt.Errorf("calculate training cost: %w", err)
	}

	// 8. Flush lazy resources
	res, err := s.resourceRepo.Get(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get resources: %w", err)
	}
	now := time.Now().UTC()
	FlushResources(res, now)

	// 9. Check sufficient resources
	if res.Food < totalCost.Food || res.Water < totalCost.Water ||
		res.Lumber < totalCost.Lumber || res.Stone < totalCost.Stone {
		return nil, model.ErrInsufficientResources
	}

	// 9a. Check sufficient gold (per-player currency)
	if totalCost.Gold > 0 {
		econ, err := s.playerEconRepo.GetByPlayerID(ctx, playerID)
		if err != nil {
			return nil, fmt.Errorf("get player economy: %w", err)
		}
		if econ.Gold < totalCost.Gold {
			return nil, model.ErrInsufficientGold
		}
	}

	// 9b. Check sufficient population capacity
	popCost := config.TroopPopCost(troopType)
	popCap := config.CalculatePopCap(buildings)
	if res.PopUsed+popCost*quantity > popCap {
		return nil, ErrInsufficientPop
	}

	// 10. Deduct resources + insert queue atomically
	res.Food -= totalCost.Food
	res.Water -= totalCost.Water
	res.Lumber -= totalCost.Lumber
	res.Stone -= totalCost.Stone
	res.PopUsed += popCost * quantity
	res.LastUpdated = now

	completesAt := now.Add(time.Duration(eachTimeSec) * time.Second)
	queueItem := &model.TrainingQueue{
		VillageID:        villageID,
		TroopType:        troopType,
		Quantity:         quantity,
		OriginalQuantity: quantity,
		EachDurationSec:  eachTimeSec,
		StartedAt:        now,
		CompletesAt:      completesAt,
	}

	err = s.uow.DeductResourcesGoldAndInsertTrainQueue(ctx, villageID, res, playerID, totalCost.Gold, queueItem)
	if err != nil {
		return nil, fmt.Errorf("execute training transaction: %w", err)
	}

	return &dto.TrainingQueueResponse{
		ID:               queueItem.ID,
		TroopType:        queueItem.TroopType,
		Quantity:         queueItem.Quantity,
		OriginalQuantity: queueItem.OriginalQuantity,
		EachDurationSec:  queueItem.EachDurationSec,
		StartedAt:        queueItem.StartedAt,
		CompletesAt:      queueItem.CompletesAt,
	}, nil
}

// CompleteTraining processes all training queue items where the next unit is done.
// For each completed unit: add 1 to troop count, increment food consumption,
// decrement queue quantity, advance completes_at (or delete if quantity → 0).
// Returns events for WebSocket notification.
func (s *TrainingService) CompleteTraining(ctx context.Context) ([]TrainCompletionEvent, error) {
	now := time.Now().UTC()
	completed, err := s.queueRepo.GetNextCompleted(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("get completed training: %w", err)
	}

	var events []TrainCompletionEvent
	for _, item := range completed {
		// Process each completed unit
		troopCfg, ok := config.TroopConfigs[item.TroopType]
		if !ok {
			slog.Error("unknown troop type in queue", "troop_type", item.TroopType, "queue_id", item.ID)
			continue
		}

		// Fetch resources so we can update food consumption
		res, err := s.resourceRepo.Get(ctx, item.VillageID)
		if err != nil {
			slog.Error("failed to get resources for training completion", "error", err, "queue_id", item.ID)
			continue
		}
		res.FoodConsumption += troopCfg.FoodUpkeep

		// Prepare queue update: decrement quantity or mark for deletion
		item.Quantity--
		deleteQueue := item.Quantity <= 0
		if !deleteQueue {
			item.CompletesAt = now.Add(time.Duration(item.EachDurationSec) * time.Second)
		}

		// Atomically: add troop + update resources + advance/delete queue
		if err := s.uow.CompleteTrainingUnit(ctx, item.VillageID, item.TroopType, 1, res, item, deleteQueue); err != nil {
			slog.Error("failed to complete training unit", "error", err, "queue_id", item.ID)
			continue
		}

		// Look up the new total for the event
		var newTotal int
		troop, err := s.troopRepo.GetByVillageAndType(ctx, item.VillageID, item.TroopType)
		if err == nil {
			newTotal = troop.Quantity
		}

		// Look up player owning this village
		village, err := s.villageRepo.GetByID(ctx, item.VillageID)
		if err != nil {
			slog.Warn("could not look up village owner for notification", "village_id", item.VillageID, "error", err)
			continue
		}

		events = append(events, TrainCompletionEvent{
			PlayerID:  village.PlayerID,
			VillageID: item.VillageID,
			TroopType: item.TroopType,
			NewTotal:  newTotal,
		})

		slog.Info("training unit completed",
			"village_id", item.VillageID,
			"troop_type", item.TroopType,
			"remaining", item.Quantity,
		)
	}
	return events, nil
}

// CancelTraining cancels a training queue entry. No refund (matches building cancel pattern).
func (s *TrainingService) CancelTraining(ctx context.Context, playerID, villageID, queueID int64) error {
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

	return s.queueRepo.Delete(ctx, queueID)
}

// GetTrainingQueue returns all training queue items for a village.
func (s *TrainingService) GetTrainingQueue(ctx context.Context, villageID int64) ([]dto.TrainingQueueResponse, error) {
	items, err := s.queueRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get training queue: %w", err)
	}

	result := make([]dto.TrainingQueueResponse, len(items))
	for i, item := range items {
		result[i] = dto.TrainingQueueResponse{
			ID:               item.ID,
			TroopType:        item.TroopType,
			Quantity:         item.Quantity,
			OriginalQuantity: item.OriginalQuantity,
			EachDurationSec:  item.EachDurationSec,
			StartedAt:        item.StartedAt,
			CompletesAt:      item.CompletesAt,
		}
	}
	return result, nil
}

// GetTrainingCost returns the cost preview for training a batch of troops.
func (s *TrainingService) GetTrainingCost(ctx context.Context, playerID, villageID int64, troopType string, quantity int) (*dto.TrainingCostResponse, error) {
	troopCfg, ok := config.TroopConfigs[troopType]
	if !ok {
		return nil, ErrUnknownTroop
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

	// Get training building level for time calculation
	buildings, err := s.buildingRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get buildings: %w", err)
	}

	var buildingLevel int
	for _, b := range buildings {
		if b.BuildingType == troopCfg.TrainingBuilding {
			buildingLevel = b.Level
			break
		}
	}

	eachTimeSec, err := config.TrainingTime(troopType, buildingLevel)
	if err != nil {
		if errors.Is(err, config.ErrUnknownTroop) {
			return nil, ErrUnknownTroop
		}
		return nil, fmt.Errorf("calculate training time: %w", err)
	}

	totalCost, err := config.TrainingCost(troopType, quantity)
	if err != nil {
		if errors.Is(err, config.ErrUnknownTroop) {
			return nil, ErrUnknownTroop
		}
		return nil, fmt.Errorf("calculate cost: %w", err)
	}

	return &dto.TrainingCostResponse{
		TroopType:    troopType,
		Quantity:     quantity,
		TotalFood:    totalCost.Food,
		TotalWater:   totalCost.Water,
		TotalLumber:  totalCost.Lumber,
		TotalStone:   totalCost.Stone,
		TotalGold:    totalCost.Gold,
		EachTimeSec:  eachTimeSec,
		TotalTimeSec: eachTimeSec * quantity,
	}, nil
}

// InstantCompleteTraining sets a training queue item's completes_at to now so the
// game loop picks it up on the next tick. Admin-only — no ownership check.
func (s *TrainingService) InstantCompleteTraining(ctx context.Context, queueID int64) error {
	item, err := s.queueRepo.GetByID(ctx, queueID)
	if err != nil {
		return fmt.Errorf("get training queue item: %w", err)
	}
	now := time.Now().UTC()
	item.CompletesAt = now
	if err := s.queueRepo.Update(ctx, item); err != nil {
		return fmt.Errorf("update training queue item: %w", err)
	}
	return nil
}

// GetTroops returns all troops stationed in a village. Requires ownership of the village.
func (s *TrainingService) GetTroops(ctx context.Context, playerID, villageID int64) ([]dto.TroopInfo, error) {
	// Verify village ownership
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

	troops, err := s.troopRepo.GetByVillageID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get troops: %w", err)
	}

	result := make([]dto.TroopInfo, len(troops))
	for i, t := range troops {
		result[i] = dto.TroopInfo{
			Type:     t.Type,
			Quantity: t.Quantity,
			Status:   t.Status,
		}
	}
	return result, nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"

	"github.com/luisfpires18/woo/internal/config"
	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Village errors.
var (
	ErrVillageNotFound = errors.New("village not found")
	ErrNotOwner        = errors.New("you do not own this village")
	ErrNoSpawnTile     = errors.New("no available spawn tile found")
	ErrInvalidName     = errors.New("village name must be between 2 and 30 characters")
)

// Starting building types for every village.
var commonBuildings = []string{
	"town_hall",
	"food_1", "food_2", "food_3",
	"water_1", "water_2", "water_3",
	"lumber_1", "lumber_2", "lumber_3",
	"stone_1", "stone_2", "stone_3",
	"barracks", "stable", "archery", "workshop", "special",
	"storage", "provisions", "reservoir",
}

// Starting resource config.
const (
	maxSpawnAttempts = 100
)

// VillageService handles village business logic.
type VillageService struct {
	villageRepo     repository.VillageRepository
	buildingRepo    repository.BuildingRepository
	resourceRepo    repository.ResourceRepository
	playerEconRepo  repository.PlayerEconomyRepository
	mapService      *MapService
}

// NewVillageService creates a new VillageService.
func NewVillageService(
	villageRepo repository.VillageRepository,
	buildingRepo repository.BuildingRepository,
	resourceRepo repository.ResourceRepository,
	playerEconRepo repository.PlayerEconomyRepository,
	mapService *MapService,
) *VillageService {
	return &VillageService{
		villageRepo:    villageRepo,
		buildingRepo:   buildingRepo,
		resourceRepo:   resourceRepo,
		playerEconRepo: playerEconRepo,
		mapService:     mapService,
	}
}

// CreateFirstVillage creates a player's first (capital) village with starter buildings and resources.
func (s *VillageService) CreateFirstVillage(ctx context.Context, playerID int64, kingdom, username string) (*model.Village, error) {
	return s.createVillageCore(ctx, playerID, kingdom, username, nil)
}

// CreateFirstVillageForSeason creates a player's first village in a specific season.
func (s *VillageService) CreateFirstVillageForSeason(ctx context.Context, playerID int64, kingdom, username string, seasonID int64) (*model.Village, error) {
	return s.createVillageCore(ctx, playerID, kingdom, username, &seasonID)
}

// createVillageCore is the shared implementation for village creation. If seasonID is non-nil
// the village is scoped to that season.
func (s *VillageService) createVillageCore(ctx context.Context, playerID int64, kingdom, username string, seasonID *int64) (*model.Village, error) {
	x, y, err := s.findSpawnLocation(ctx, kingdom)
	if err != nil {
		return nil, fmt.Errorf("find spawn location: %w", err)
	}

	now := time.Now().UTC()
	village := &model.Village{
		PlayerID:  playerID,
		Name:      username + "'s Village",
		X:         x,
		Y:         y,
		IsCapital: true,
		SeasonID:  seasonID,
		CreatedAt: now,
	}
	if err := s.villageRepo.Create(ctx, village); err != nil {
		return nil, fmt.Errorf("create village: %w", err)
	}

	// Link the map tile to this village and player
	if s.mapService != nil {
		if err := s.mapService.UpdateTileOwner(ctx, x, y, playerID, village.ID); err != nil {
			return nil, fmt.Errorf("link tile to village: %w", err)
		}
	}

	// Create starter buildings (all at level 0 = slot exists but not built)
	buildings := make([]*model.Building, 0, len(commonBuildings)+1)
	for _, bt := range commonBuildings {
		buildings = append(buildings, &model.Building{
			VillageID:    village.ID,
			BuildingType: bt,
			Level:        0,
		})
	}
	if err := s.buildingRepo.CreateBatch(ctx, buildings); err != nil {
		return nil, fmt.Errorf("create starter buildings: %w", err)
	}

	// Create initial resources
	resources := &model.Resources{
		VillageID:       village.ID,
		Food:            config.StartingResources,
		Water:           config.StartingResources,
		Lumber:          config.StartingResources,
		Stone:           config.StartingResources,
		FoodRate:        config.StartingRate,
		WaterRate:       config.StartingRate,
		LumberRate:      config.StartingRate,
		StoneRate:       config.StartingRate,
		FoodConsumption: 0,
		MaxFood:         config.StartingStorage,
		MaxWater:        config.StartingStorage,
		MaxLumber:       config.StartingStorage,
		MaxStone:        config.StartingStorage,
		LastUpdated:     now,
	}
	if err := s.resourceRepo.Create(ctx, resources); err != nil {
		return nil, fmt.Errorf("create starter resources: %w", err)
	}

	// Create player economy (gold) — idempotent, only for first village.
	if s.playerEconRepo != nil {
		_, econErr := s.playerEconRepo.GetByPlayerID(ctx, playerID)
		if econErr != nil {
			if err := s.playerEconRepo.Create(ctx, playerID, config.StartingGold); err != nil {
				return nil, fmt.Errorf("create player economy: %w", err)
			}
		}
	}

	return village, nil
}

// GetVillage retrieves a village by ID and verifies ownership. Resources are lazily calculated.
func (s *VillageService) GetVillage(ctx context.Context, villageID, playerID int64) (*dto.VillageResponse, error) {
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

	resources, err := s.getCalculatedResources(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get resources: %w", err)
	}

	// Fetch player gold
	var gold float64
	if s.playerEconRepo != nil {
		econ, err := s.playerEconRepo.GetByPlayerID(ctx, playerID)
		if err == nil {
			gold = econ.Gold
		}
	}

	resp := s.buildVillageResponse(village, buildings, resources)
	resp.Gold = gold
	return resp, nil
}

// ListVillages lists all villages for a player.
func (s *VillageService) ListVillages(ctx context.Context, playerID int64) ([]dto.VillageListItem, error) {
	villages, err := s.villageRepo.ListByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("list villages: %w", err)
	}

	items := make([]dto.VillageListItem, 0, len(villages))
	for _, v := range villages {
		items = append(items, dto.VillageListItem{
			ID:        v.ID,
			Name:      v.Name,
			X:         v.X,
			Y:         v.Y,
			IsCapital: v.IsCapital,
		})
	}
	return items, nil
}

// getCalculatedResources performs lazy resource calculation: stored + rate × elapsed, capped at max.
// If no resources row exists (e.g. admin-seeded village), it self-heals by creating default resources.
func (s *VillageService) getCalculatedResources(ctx context.Context, villageID int64) (*model.Resources, error) {
	res, err := s.resourceRepo.Get(ctx, villageID)
	if errors.Is(err, model.ErrNotFound) {
		// Self-heal: create default resources for villages missing a resources row
		now := time.Now().UTC()
		res = &model.Resources{
			VillageID:       villageID,
			Food:            config.StartingResources,
			Water:           config.StartingResources,
			Lumber:          config.StartingResources,
			Stone:           config.StartingResources,
			FoodRate:        config.StartingRate,
			WaterRate:       config.StartingRate,
			LumberRate:      config.StartingRate,
			StoneRate:       config.StartingRate,
			FoodConsumption: 0,
			MaxFood:         config.StartingStorage,
			MaxWater:        config.StartingStorage,
			MaxLumber:       config.StartingStorage,
			MaxStone:        config.StartingStorage,
			LastUpdated:     now,
		}
		if createErr := s.resourceRepo.Create(ctx, res); createErr != nil {
			return nil, fmt.Errorf("create missing resources for village %d: %w", villageID, createErr)
		}
		return res, nil
	}
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	if !FlushResources(res, now) {
		return res, nil
	}

	// Persist the recalculated snapshot
	if err := s.resourceRepo.Update(ctx, villageID, res); err != nil {
		return nil, fmt.Errorf("update calculated resources: %w", err)
	}

	return res, nil
}

// findSpawnLocation finds a suitable tile for a new village, preferring the player's kingdom zone.
func (s *VillageService) findSpawnLocation(ctx context.Context, kingdom string) (int, int, error) {
	// If map service is available, use zone-aware spawning (auto-places zone if needed)
	if s.mapService != nil {
		return s.mapService.FindSpawnTile(ctx, kingdom)
	}

	// Fallback: random coordinate (pre-map-generation path)
	for i := 0; i < maxSpawnAttempts; i++ {
		x := rand.Intn(model.MapHalf*2+1) - model.MapHalf
		y := rand.Intn(model.MapHalf*2+1) - model.MapHalf

		_, err := s.villageRepo.GetByCoordinates(ctx, x, y)
		if errors.Is(err, model.ErrNotFound) {
			return x, y, nil // tile is free
		}
		if err != nil {
			return 0, 0, fmt.Errorf("check coordinates (%d,%d): %w", x, y, err)
		}
		// Tile occupied, try again
	}
	return 0, 0, ErrNoSpawnTile
}

func (s *VillageService) buildVillageResponse(village *model.Village, buildings []*model.Building, resources *model.Resources) *dto.VillageResponse {
	buildingInfos := make([]dto.BuildingInfo, 0, len(buildings))
	for _, b := range buildings {
		buildingInfos = append(buildingInfos, dto.BuildingInfo{
			ID:           b.ID,
			BuildingType: b.BuildingType,
			Level:        b.Level,
		})
	}

	popCap := config.CalculatePopCap(buildings)

	return &dto.VillageResponse{
		ID:        village.ID,
		Name:      village.Name,
		X:         village.X,
		Y:         village.Y,
		IsCapital: village.IsCapital,
		Buildings: buildingInfos,
		Resources: &dto.ResourcesResponse{
			Food:            resources.Food,
			Water:           resources.Water,
			Lumber:          resources.Lumber,
			Stone:           resources.Stone,
			FoodRate:        resources.FoodRate,
			WaterRate:       resources.WaterRate,
			LumberRate:      resources.LumberRate,
			StoneRate:       resources.StoneRate,
			FoodConsumption: resources.FoodConsumption,
			MaxFood:         resources.MaxFood,
			MaxWater:        resources.MaxWater,
			MaxLumber:       resources.MaxLumber,
			MaxStone:        resources.MaxStone,
			PopCap:          popCap,
			PopUsed:         resources.PopUsed,
		},
	}
}

// RenameVillage changes the name of a village owned by the player.
func (s *VillageService) RenameVillage(ctx context.Context, villageID, playerID int64, newName string) (*dto.VillageListItem, error) {
	// Validate name length
	if len(newName) < 2 || len(newName) > 30 {
		return nil, ErrInvalidName
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

	village.Name = newName
	if err := s.villageRepo.Update(ctx, village); err != nil {
		return nil, fmt.Errorf("update village name: %w", err)
	}

	return &dto.VillageListItem{
		ID:        village.ID,
		Name:      village.Name,
		X:         village.X,
		Y:         village.Y,
		IsCapital: village.IsCapital,
	}, nil
}

package service

import (
	"context"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Village errors.
var (
	ErrVillageNotFound = errors.New("village not found")
	ErrNotOwner        = errors.New("you do not own this village")
	ErrNoSpawnTile     = errors.New("no available spawn tile found")
)

// Starting building types for every village. Kingdom-specific building is appended.
var commonBuildings = []string{
	"town_hall", "iron_mine", "lumber_mill", "quarry", "farm", "warehouse",
	"barracks", "stable", "forge", "rune_altar", "walls", "marketplace",
	"embassy", "watchtower",
}

// Kingdom-specific bonus building.
var kingdomBuilding = map[string]string{
	"veridor": "dock",
	"sylvara": "grove_sanctum",
	"arkazia": "colosseum",
}

// Starting resource config.
const (
	startingResources   = 500.0
	startingRate        = 30.0
	startingStorage     = 1000.0
	mapHalf             = 200 // map goes from -200 to +200
	spawnMinDist        = 10  // don't spawn within 10 tiles of center (Moraphys)
	maxSpawnAttempts    = 100
)

// VillageService handles village business logic.
type VillageService struct {
	villageRepo  repository.VillageRepository
	buildingRepo repository.BuildingRepository
	resourceRepo repository.ResourceRepository
}

// NewVillageService creates a new VillageService.
func NewVillageService(
	villageRepo repository.VillageRepository,
	buildingRepo repository.BuildingRepository,
	resourceRepo repository.ResourceRepository,
) *VillageService {
	return &VillageService{
		villageRepo:  villageRepo,
		buildingRepo: buildingRepo,
		resourceRepo: resourceRepo,
	}
}

// CreateFirstVillage creates a player's first (capital) village with starter buildings and resources.
func (s *VillageService) CreateFirstVillage(ctx context.Context, playerID int64, kingdom, username string) (*model.Village, error) {
	x, y, err := s.findSpawnLocation(ctx)
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
		CreatedAt: now,
	}
	if err := s.villageRepo.Create(ctx, village); err != nil {
		return nil, fmt.Errorf("create village: %w", err)
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
	// Add kingdom-specific building
	if kb, ok := kingdomBuilding[kingdom]; ok {
		buildings = append(buildings, &model.Building{
			VillageID:    village.ID,
			BuildingType: kb,
			Level:        0,
		})
	}
	if err := s.buildingRepo.CreateBatch(ctx, buildings); err != nil {
		return nil, fmt.Errorf("create starter buildings: %w", err)
	}

	// Create initial resources
	resources := &model.Resources{
		VillageID:       village.ID,
		Iron:            startingResources,
		Wood:            startingResources,
		Stone:           startingResources,
		Food:            startingResources,
		IronRate:        startingRate,
		WoodRate:        startingRate,
		StoneRate:       startingRate,
		FoodRate:        startingRate,
		FoodConsumption: 0,
		MaxStorage:      startingStorage,
		LastUpdated:     now,
	}
	if err := s.resourceRepo.Create(ctx, resources); err != nil {
		return nil, fmt.Errorf("create starter resources: %w", err)
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

	return s.buildVillageResponse(village, buildings, resources), nil
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
func (s *VillageService) getCalculatedResources(ctx context.Context, villageID int64) (*model.Resources, error) {
	res, err := s.resourceRepo.Get(ctx, villageID)
	if err != nil {
		return nil, err
	}

	now := time.Now().UTC()
	elapsed := now.Sub(res.LastUpdated).Hours()
	if elapsed <= 0 {
		return res, nil
	}

	// Calculate current resources: stored + rate*hours, capped at max_storage
	res.Iron = math.Min(res.Iron+res.IronRate*elapsed, res.MaxStorage)
	res.Wood = math.Min(res.Wood+res.WoodRate*elapsed, res.MaxStorage)
	res.Stone = math.Min(res.Stone+res.StoneRate*elapsed, res.MaxStorage)
	res.Food = math.Min(res.Food+(res.FoodRate-res.FoodConsumption)*elapsed, res.MaxStorage)
	if res.Food < 0 {
		res.Food = 0
	}
	res.LastUpdated = now

	// Persist the recalculated snapshot
	if err := s.resourceRepo.Update(ctx, villageID, res); err != nil {
		return nil, fmt.Errorf("update calculated resources: %w", err)
	}

	return res, nil
}

// findSpawnLocation finds a random unoccupied tile on the world map for a new village.
func (s *VillageService) findSpawnLocation(ctx context.Context) (int, int, error) {
	for i := 0; i < maxSpawnAttempts; i++ {
		x := rand.Intn(mapHalf*2+1) - mapHalf // -200 to 200
		y := rand.Intn(mapHalf*2+1) - mapHalf

		// Skip center zone (Moraphys stronghold area)
		if abs(x) < spawnMinDist && abs(y) < spawnMinDist {
			continue
		}

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

	return &dto.VillageResponse{
		ID:        village.ID,
		Name:      village.Name,
		X:         village.X,
		Y:         village.Y,
		IsCapital: village.IsCapital,
		Buildings: buildingInfos,
		Resources: &dto.ResourcesResponse{
			Iron:            resources.Iron,
			Wood:            resources.Wood,
			Stone:           resources.Stone,
			Food:            resources.Food,
			IronRate:        resources.IronRate,
			WoodRate:        resources.WoodRate,
			StoneRate:       resources.StoneRate,
			FoodRate:        resources.FoodRate,
			FoodConsumption: resources.FoodConsumption,
			MaxStorage:      resources.MaxStorage,
		},
	}
}

func abs(x int) int {
	if x < 0 {
		return -x
	}
	return x
}

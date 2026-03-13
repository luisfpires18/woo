package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"log/slog"
	"math"
	"math/rand"
	"time"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Camp spawning errors.
var (
	ErrNoSpawnCandidates = errors.New("no valid spawn candidate tiles")
	ErrMaxCampsReached   = errors.New("maximum camps reached for this spawn rule")
)

// CampService handles camp spawning, despawning, and querying.
type CampService struct {
	campRepo          repository.CampRepository
	campTemplateRepo  repository.CampTemplateRepository
	beastSlotRepo     repository.CampBeastSlotRepository
	beastTemplateRepo repository.BeastTemplateRepository
	spawnRuleRepo     repository.SpawnRuleRepository
	worldMapRepo      repository.WorldMapRepository
	villageRepo       repository.VillageRepository
}

// NewCampService creates a new CampService.
func NewCampService(
	campRepo repository.CampRepository,
	campTemplateRepo repository.CampTemplateRepository,
	beastSlotRepo repository.CampBeastSlotRepository,
	beastTemplateRepo repository.BeastTemplateRepository,
	spawnRuleRepo repository.SpawnRuleRepository,
	worldMapRepo repository.WorldMapRepository,
	villageRepo repository.VillageRepository,
) *CampService {
	return &CampService{
		campRepo:          campRepo,
		campTemplateRepo:  campTemplateRepo,
		beastSlotRepo:     beastSlotRepo,
		beastTemplateRepo: beastTemplateRepo,
		spawnRuleRepo:     spawnRuleRepo,
		worldMapRepo:      worldMapRepo,
		villageRepo:       villageRepo,
	}
}

// SpawnCamps processes all enabled spawn rules and spawns camps where needed.
// Called periodically by the game loop.
func (s *CampService) SpawnCamps(ctx context.Context) (int, error) {
	rules, err := s.spawnRuleRepo.GetEnabled(ctx)
	if err != nil {
		return 0, fmt.Errorf("get enabled spawn rules: %w", err)
	}

	totalSpawned := 0
	rng := rand.New(rand.NewSource(time.Now().UnixNano()))

	for _, rule := range rules {
		count, err := s.spawnForRule(ctx, rule, rng)
		if err != nil {
			slog.Warn("camp spawn failed for rule", "rule_id", rule.ID, "error", err)
			continue
		}
		totalSpawned += count
	}
	return totalSpawned, nil
}

// DespawnExpiredCamps removes camps that have exceeded their lifetime.
// Called periodically by the game loop.
func (s *CampService) DespawnExpiredCamps(ctx context.Context) (int, error) {
	now := time.Now().UTC()
	expired, err := s.campRepo.GetExpiredCamps(ctx, now)
	if err != nil {
		return 0, fmt.Errorf("get expired camps: %w", err)
	}

	despawned := 0
	for _, camp := range expired {
		if err := s.worldMapRepo.UpdateTileOwner(ctx, camp.TileX, camp.TileY, nil, nil); err != nil {
			slog.Warn("failed to clear tile for despawned camp", "camp_id", camp.ID, "error", err)
			continue
		}
		if err := s.campRepo.Delete(ctx, camp.ID); err != nil {
			slog.Warn("failed to delete despawned camp", "camp_id", camp.ID, "error", err)
			continue
		}
		despawned++
	}
	return despawned, nil
}

// ListActiveCamps returns all active camp instances as DTOs enriched with template data.
func (s *CampService) ListActiveCamps(ctx context.Context) ([]dto.CampResponse, error) {
	camps, err := s.campRepo.ListActive(ctx)
	if err != nil {
		return nil, err
	}

	result := make([]dto.CampResponse, 0, len(camps))
	for _, c := range camps {
		tmpl, err := s.campTemplateRepo.GetByID(ctx, c.CampTemplateID)
		if err != nil {
			slog.Warn("skip camp with missing template", "camp_id", c.ID, "template_id", c.CampTemplateID)
			continue
		}

		cr := dto.CampResponse{
			ID:           c.ID,
			TemplateName: tmpl.Name,
			Tier:         tmpl.Tier,
			SpriteKey:    tmpl.SpriteKey,
			TileX:        c.TileX,
			TileY:        c.TileY,
			Status:       c.Status,
			SpawnedAt:    time.Now(), // default
		}
		if t, err2 := time.Parse(time.RFC3339, c.SpawnedAt); err2 == nil {
			cr.SpawnedAt = t
		}
		result = append(result, cr)
	}
	return result, nil
}

// GetCampWithBeasts returns a camp with its beast details and camp template.
func (s *CampService) GetCampWithBeasts(ctx context.Context, campID int64) (*model.Camp, []model.CampBeast, *model.CampTemplate, error) {
	camp, err := s.campRepo.GetByID(ctx, campID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get camp: %w", err)
	}

	var beasts []model.CampBeast
	if err := json.Unmarshal([]byte(camp.BeastsJSON), &beasts); err != nil {
		return nil, nil, nil, fmt.Errorf("parse camp beasts: %w", err)
	}

	tmpl, err := s.campTemplateRepo.GetByID(ctx, camp.CampTemplateID)
	if err != nil {
		return nil, nil, nil, fmt.Errorf("get camp template: %w", err)
	}

	return camp, beasts, tmpl, nil
}

// ── Internal spawn logic ─────────────────────────────────────────────────────

func (s *CampService) spawnForRule(ctx context.Context, rule *model.SpawnRule, rng *rand.Rand) (int, error) {
	// Check current count
	currentCount, err := s.campRepo.CountBySpawnRule(ctx, rule.ID)
	if err != nil {
		return 0, fmt.Errorf("count camps for rule: %w", err)
	}
	if currentCount >= rule.MaxCamps {
		return 0, nil // at capacity
	}

	// Parse terrain and zone filters
	var terrainTypes []string
	if err := json.Unmarshal([]byte(rule.TerrainTypesJSON), &terrainTypes); err != nil {
		return 0, fmt.Errorf("parse terrain types: %w", err)
	}

	var zoneTypes []string
	if err := json.Unmarshal([]byte(rule.ZoneTypesJSON), &zoneTypes); err != nil {
		return 0, fmt.Errorf("parse zone types: %w", err)
	}

	// Parse camp template pool
	var pool []model.CampTemplatePoolEntry
	if err := json.Unmarshal([]byte(rule.CampTemplatePoolJSON), &pool); err != nil {
		return 0, fmt.Errorf("parse camp template pool: %w", err)
	}
	if len(pool) == 0 {
		return 0, nil
	}

	// Get candidate tiles (wilderness tiles with no village and no camp)
	candidates, err := s.getSpawnCandidates(ctx, terrainTypes, zoneTypes, rule)
	if err != nil {
		return 0, err
	}
	if len(candidates) == 0 {
		return 0, nil
	}

	// Spawn one camp per tick (rate-limited)
	toSpawn := rule.MaxCamps - currentCount
	if toSpawn > 1 {
		toSpawn = 1 // one per tick per rule
	}

	spawned := 0
	for i := 0; i < toSpawn && len(candidates) > 0; i++ {
		// Pick random candidate tile
		idx := rng.Intn(len(candidates))
		tile := candidates[idx]
		candidates = append(candidates[:idx], candidates[idx+1:]...)

		// Pick weighted random camp template
		templateID := pickWeightedTemplate(pool, rng)
		if templateID == 0 {
			continue
		}

		// Load template and generate beasts
		beasts, err := s.generateBeasts(ctx, templateID, rng)
		if err != nil {
			slog.Warn("failed to generate beasts for camp", "template_id", templateID, "error", err)
			continue
		}

		beastsJSON, err := json.Marshal(beasts)
		if err != nil {
			continue
		}

		now := time.Now().UTC().Format(time.RFC3339)
		camp := &model.Camp{
			CampTemplateID: templateID,
			TileX:          tile.X,
			TileY:          tile.Y,
			BeastsJSON:     string(beastsJSON),
			SpawnedAt:      now,
			Status:         model.CampStatusActive,
			SpawnRuleID:    &rule.ID,
		}

		if err := s.campRepo.Create(ctx, camp); err != nil {
			slog.Warn("failed to create camp", "error", err)
			continue
		}

		spawned++
	}

	return spawned, nil
}

func (s *CampService) getSpawnCandidates(ctx context.Context, terrainTypes, zoneTypes []string, rule *model.SpawnRule) ([]*model.MapTile, error) {
	// Get all potential tiles matching terrain/zone criteria
	var allTiles []*model.MapTile

	// Get tiles per zone
	for _, zone := range zoneTypes {
		tiles, err := s.worldMapRepo.GetByZone(ctx, zone)
		if err != nil {
			return nil, fmt.Errorf("get tiles by zone %s: %w", zone, err)
		}
		allTiles = append(allTiles, tiles...)
	}

	// Filter by terrain type and availability
	terrainSet := make(map[string]bool, len(terrainTypes))
	for _, t := range terrainTypes {
		terrainSet[t] = true
	}

	// Get active camps for distance checks
	activeCamps, err := s.campRepo.ListActive(ctx)
	if err != nil {
		return nil, fmt.Errorf("list active camps: %w", err)
	}

	var candidates []*model.MapTile
	for _, tile := range allTiles {
		// Must match terrain type
		if !terrainSet[tile.TerrainType] {
			continue
		}
		// Must not have a village
		if tile.VillageID != nil {
			continue
		}
		// Must not have a camp already
		if tile.CampID != nil {
			continue
		}
		// Check min distance from other camps
		if !checkMinDistance(tile, activeCamps, rule.MinCampDistance) {
			continue
		}
		candidates = append(candidates, tile)
	}

	return candidates, nil
}

func checkMinDistance(tile *model.MapTile, camps []*model.Camp, minDist int) bool {
	if minDist <= 0 {
		return true
	}
	for _, camp := range camps {
		dx := float64(tile.X - camp.TileX)
		dy := float64(tile.Y - camp.TileY)
		dist := math.Sqrt(dx*dx + dy*dy)
		if dist < float64(minDist) {
			return false
		}
	}
	return true
}

func pickWeightedTemplate(pool []model.CampTemplatePoolEntry, rng *rand.Rand) int64 {
	totalWeight := 0
	for _, entry := range pool {
		totalWeight += entry.Weight
	}
	if totalWeight == 0 {
		return 0
	}

	roll := rng.Intn(totalWeight)
	cumulative := 0
	for _, entry := range pool {
		cumulative += entry.Weight
		if roll < cumulative {
			return entry.CampTemplateID
		}
	}
	return pool[len(pool)-1].CampTemplateID
}

func (s *CampService) generateBeasts(ctx context.Context, campTemplateID int64, rng *rand.Rand) ([]model.CampBeast, error) {
	slots, err := s.beastSlotRepo.GetByCampTemplateID(ctx, campTemplateID)
	if err != nil {
		return nil, fmt.Errorf("get beast slots: %w", err)
	}

	var beasts []model.CampBeast
	for _, slot := range slots {
		template, err := s.beastTemplateRepo.GetByID(ctx, slot.BeastTemplateID)
		if err != nil {
			return nil, fmt.Errorf("get beast template %d: %w", slot.BeastTemplateID, err)
		}

		// Random count between min and max
		count := slot.MinCount
		if slot.MaxCount > slot.MinCount {
			count = slot.MinCount + rng.Intn(slot.MaxCount-slot.MinCount+1)
		}

		for i := 0; i < count; i++ {
			beasts = append(beasts, model.CampBeast{
				BeastTemplateID:   template.ID,
				Name:              template.Name,
				SpriteKey:         template.SpriteKey,
				HP:                template.HP,
				MaxHP:             template.HP,
				AttackPower:       template.AttackPower,
				AttackInterval:    template.AttackInterval,
				DefensePercent:    template.DefensePercent,
				CritChancePercent: template.CritChancePercent,
			})
		}
	}

	return beasts, nil
}

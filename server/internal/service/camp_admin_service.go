package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Camp admin errors.
var (
	ErrBeastTemplateNotFound = errors.New("beast template not found")
	ErrCampTemplateNotFound  = errors.New("camp template not found")
	ErrSpawnRuleNotFound     = errors.New("spawn rule not found")
	ErrRewardTableNotFound   = errors.New("reward table not found")
)

// CampAdminService handles admin CRUD for camp-related configuration.
type CampAdminService struct {
	beastTemplateRepo repository.BeastTemplateRepository
	campTemplateRepo  repository.CampTemplateRepository
	beastSlotRepo     repository.CampBeastSlotRepository
	spawnRuleRepo     repository.SpawnRuleRepository
	rewardTableRepo   repository.RewardTableRepository
	rewardEntryRepo   repository.RewardTableEntryRepository
	battleTuningRepo  repository.BattleTuningRepository
	auditLogRepo      repository.AdminAuditLogRepository
}

// NewCampAdminService creates a new CampAdminService.
func NewCampAdminService(
	beastTemplateRepo repository.BeastTemplateRepository,
	campTemplateRepo repository.CampTemplateRepository,
	beastSlotRepo repository.CampBeastSlotRepository,
	spawnRuleRepo repository.SpawnRuleRepository,
	rewardTableRepo repository.RewardTableRepository,
	rewardEntryRepo repository.RewardTableEntryRepository,
	battleTuningRepo repository.BattleTuningRepository,
	auditLogRepo repository.AdminAuditLogRepository,
) *CampAdminService {
	return &CampAdminService{
		beastTemplateRepo: beastTemplateRepo,
		campTemplateRepo:  campTemplateRepo,
		beastSlotRepo:     beastSlotRepo,
		spawnRuleRepo:     spawnRuleRepo,
		rewardTableRepo:   rewardTableRepo,
		rewardEntryRepo:   rewardEntryRepo,
		battleTuningRepo:  battleTuningRepo,
		auditLogRepo:      auditLogRepo,
	}
}

// ── Beast Templates ──────────────────────────────────────────────────────────

// ListBeastTemplates returns all beast templates.
func (s *CampAdminService) ListBeastTemplates(ctx context.Context) ([]dto.BeastTemplateResponse, error) {
	templates, err := s.beastTemplateRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list beast templates: %w", err)
	}
	result := make([]dto.BeastTemplateResponse, len(templates))
	for i, t := range templates {
		result[i] = toBeastTemplateDTO(t)
	}
	return result, nil
}

// CreateBeastTemplate creates a new beast template.
func (s *CampAdminService) CreateBeastTemplate(ctx context.Context, adminID int64, req dto.CreateBeastTemplateRequest) (*dto.BeastTemplateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	bt := &model.BeastTemplate{
		Name:              req.Name,
		SpriteKey:         req.SpriteKey,
		HP:                req.HP,
		AttackPower:       req.AttackPower,
		AttackInterval:    req.AttackInterval,
		DefensePercent:    req.DefensePercent,
		CritChancePercent: req.CritChancePercent,
		Description:       req.Description,
		CreatedAt:         now,
		UpdatedAt:         now,
		UpdatedBy:         &adminID,
	}

	if err := s.beastTemplateRepo.Create(ctx, bt); err != nil {
		return nil, fmt.Errorf("create beast template: %w", err)
	}

	s.audit(ctx, adminID, "create", "beast_template", &bt.ID, nil, bt)

	resp := toBeastTemplateDTO(bt)
	return &resp, nil
}

// UpdateBeastTemplate updates an existing beast template.
func (s *CampAdminService) UpdateBeastTemplate(ctx context.Context, adminID, id int64, req dto.UpdateBeastTemplateRequest) (*dto.BeastTemplateResponse, error) {
	bt, err := s.beastTemplateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get beast template: %w", err)
	}
	if bt == nil {
		return nil, ErrBeastTemplateNotFound
	}

	old := *bt
	if req.Name != nil {
		bt.Name = *req.Name
	}
	if req.SpriteKey != nil {
		bt.SpriteKey = *req.SpriteKey
	}
	if req.HP != nil {
		bt.HP = *req.HP
	}
	if req.AttackPower != nil {
		bt.AttackPower = *req.AttackPower
	}
	if req.AttackInterval != nil {
		bt.AttackInterval = *req.AttackInterval
	}
	if req.DefensePercent != nil {
		bt.DefensePercent = *req.DefensePercent
	}
	if req.CritChancePercent != nil {
		bt.CritChancePercent = *req.CritChancePercent
	}
	if req.Description != nil {
		bt.Description = *req.Description
	}
	bt.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	bt.UpdatedBy = &adminID

	if err := s.beastTemplateRepo.Update(ctx, bt); err != nil {
		return nil, fmt.Errorf("update beast template: %w", err)
	}

	s.audit(ctx, adminID, "update", "beast_template", &id, &old, bt)

	resp := toBeastTemplateDTO(bt)
	return &resp, nil
}

// DeleteBeastTemplate deletes a beast template.
func (s *CampAdminService) DeleteBeastTemplate(ctx context.Context, adminID, id int64) error {
	bt, err := s.beastTemplateRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get beast template: %w", err)
	}
	if bt == nil {
		return ErrBeastTemplateNotFound
	}
	if err := s.beastTemplateRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete beast template: %w", err)
	}
	s.audit(ctx, adminID, "delete", "beast_template", &id, bt, nil)
	return nil
}

// ── Camp Templates ───────────────────────────────────────────────────────────

// ListCampTemplates returns all camp templates with their beast slots.
func (s *CampAdminService) ListCampTemplates(ctx context.Context) ([]dto.CampTemplateResponse, error) {
	templates, err := s.campTemplateRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list camp templates: %w", err)
	}

	result := make([]dto.CampTemplateResponse, len(templates))
	for i, ct := range templates {
		slots, err := s.beastSlotRepo.GetByCampTemplateID(ctx, ct.ID)
		if err != nil {
			return nil, fmt.Errorf("get beast slots for template %d: %w", ct.ID, err)
		}

		slotDTOs := make([]dto.CampBeastSlotResponse, len(slots))
		for j, slot := range slots {
			beastName := ""
			if bt, _ := s.beastTemplateRepo.GetByID(ctx, slot.BeastTemplateID); bt != nil {
				beastName = bt.Name
			}
			slotDTOs[j] = dto.CampBeastSlotResponse{
				ID:              slot.ID,
				BeastTemplateID: slot.BeastTemplateID,
				BeastName:       beastName,
				MinCount:        slot.MinCount,
				MaxCount:        slot.MaxCount,
			}
		}

		result[i] = dto.CampTemplateResponse{
			ID:            ct.ID,
			Name:          ct.Name,
			Tier:          ct.Tier,
			SpriteKey:     ct.SpriteKey,
			Description:   ct.Description,
			RewardTableID: ct.RewardTableID,
			BeastSlots:    slotDTOs,
			CreatedAt:     ct.CreatedAt,
			UpdatedAt:     ct.UpdatedAt,
		}
	}
	return result, nil
}

// CreateCampTemplate creates a new camp template with beast slots.
func (s *CampAdminService) CreateCampTemplate(ctx context.Context, adminID int64, req dto.CreateCampTemplateRequest) (*dto.CampTemplateResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	ct := &model.CampTemplate{
		Name:          req.Name,
		Tier:          req.Tier,
		SpriteKey:     req.SpriteKey,
		Description:   req.Description,
		RewardTableID: req.RewardTableID,
		CreatedAt:     now,
		UpdatedAt:     now,
		UpdatedBy:     &adminID,
	}

	if err := s.campTemplateRepo.Create(ctx, ct); err != nil {
		return nil, fmt.Errorf("create camp template: %w", err)
	}

	// Create beast slots
	var slotDTOs []dto.CampBeastSlotResponse
	for _, slotReq := range req.BeastSlots {
		slot := &model.CampBeastSlot{
			CampTemplateID:  ct.ID,
			BeastTemplateID: slotReq.BeastTemplateID,
			MinCount:        slotReq.MinCount,
			MaxCount:        slotReq.MaxCount,
		}
		if err := s.beastSlotRepo.Create(ctx, slot); err != nil {
			return nil, fmt.Errorf("create beast slot: %w", err)
		}
		beastName := ""
		if bt, _ := s.beastTemplateRepo.GetByID(ctx, slotReq.BeastTemplateID); bt != nil {
			beastName = bt.Name
		}
		slotDTOs = append(slotDTOs, dto.CampBeastSlotResponse{
			ID:              slot.ID,
			BeastTemplateID: slotReq.BeastTemplateID,
			BeastName:       beastName,
			MinCount:        slotReq.MinCount,
			MaxCount:        slotReq.MaxCount,
		})
	}

	s.audit(ctx, adminID, "create", "camp_template", &ct.ID, nil, ct)

	return &dto.CampTemplateResponse{
		ID:            ct.ID,
		Name:          ct.Name,
		Tier:          ct.Tier,
		SpriteKey:     ct.SpriteKey,
		Description:   ct.Description,
		RewardTableID: ct.RewardTableID,
		BeastSlots:    slotDTOs,
		CreatedAt:     ct.CreatedAt,
		UpdatedAt:     ct.UpdatedAt,
	}, nil
}

// UpdateCampTemplate updates an existing camp template.
func (s *CampAdminService) UpdateCampTemplate(ctx context.Context, adminID, id int64, req dto.UpdateCampTemplateRequest) (*dto.CampTemplateResponse, error) {
	ct, err := s.campTemplateRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get camp template: %w", err)
	}
	if ct == nil {
		return nil, ErrCampTemplateNotFound
	}

	if req.Name != nil {
		ct.Name = *req.Name
	}
	if req.Tier != nil {
		ct.Tier = *req.Tier
	}
	if req.SpriteKey != nil {
		ct.SpriteKey = *req.SpriteKey
	}
	if req.Description != nil {
		ct.Description = *req.Description
	}
	if req.RewardTableID != nil {
		ct.RewardTableID = req.RewardTableID
	}
	ct.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	ct.UpdatedBy = &adminID

	if err := s.campTemplateRepo.Update(ctx, ct); err != nil {
		return nil, fmt.Errorf("update camp template: %w", err)
	}

	// Replace beast slots if provided
	if req.BeastSlots != nil {
		if err := s.beastSlotRepo.DeleteByCampTemplateID(ctx, id); err != nil {
			return nil, fmt.Errorf("delete old beast slots: %w", err)
		}
		for _, slotReq := range *req.BeastSlots {
			slot := &model.CampBeastSlot{
				CampTemplateID:  id,
				BeastTemplateID: slotReq.BeastTemplateID,
				MinCount:        slotReq.MinCount,
				MaxCount:        slotReq.MaxCount,
			}
			if err := s.beastSlotRepo.Create(ctx, slot); err != nil {
				return nil, fmt.Errorf("create beast slot: %w", err)
			}
		}
	}

	s.audit(ctx, adminID, "update", "camp_template", &id, nil, ct)

	// Re-fetch with slots for response
	templates, _ := s.ListCampTemplates(ctx)
	for _, t := range templates {
		if t.ID == id {
			return &t, nil
		}
	}
	return nil, nil
}

// DeleteCampTemplate deletes a camp template and its beast slots.
func (s *CampAdminService) DeleteCampTemplate(ctx context.Context, adminID, id int64) error {
	ct, err := s.campTemplateRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get camp template: %w", err)
	}
	if ct == nil {
		return ErrCampTemplateNotFound
	}
	if err := s.beastSlotRepo.DeleteByCampTemplateID(ctx, id); err != nil {
		return fmt.Errorf("delete beast slots: %w", err)
	}
	if err := s.campTemplateRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete camp template: %w", err)
	}
	s.audit(ctx, adminID, "delete", "camp_template", &id, ct, nil)
	return nil
}

// ── Spawn Rules ──────────────────────────────────────────────────────────────

// ListSpawnRules returns all spawn rules.
func (s *CampAdminService) ListSpawnRules(ctx context.Context) ([]dto.SpawnRuleResponse, error) {
	rules, err := s.spawnRuleRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list spawn rules: %w", err)
	}
	result := make([]dto.SpawnRuleResponse, len(rules))
	for i, r := range rules {
		result[i] = toSpawnRuleDTO(r)
	}
	return result, nil
}

// CreateSpawnRule creates a new spawn rule.
func (s *CampAdminService) CreateSpawnRule(ctx context.Context, adminID int64, req dto.CreateSpawnRuleRequest) (*dto.SpawnRuleResponse, error) {
	terrainJSON, _ := json.Marshal(req.TerrainTypes)
	zoneJSON, _ := json.Marshal(req.ZoneTypes)
	poolJSON, _ := json.Marshal(req.CampTemplatePool)

	now := time.Now().UTC().Format(time.RFC3339)
	rule := &model.SpawnRule{
		Name:                 req.Name,
		TerrainTypesJSON:     string(terrainJSON),
		ZoneTypesJSON:        string(zoneJSON),
		CampTemplatePoolJSON: string(poolJSON),
		MaxCamps:             req.MaxCamps,
		SpawnIntervalSec:     req.SpawnIntervalSec,
		DespawnAfterSec:      req.DespawnAfterSec,
		MinCampDistance:      req.MinCampDistance,
		MinVillageDistance:   req.MinVillageDistance,
		Enabled:              req.Enabled,
		CreatedAt:            now,
		UpdatedAt:            now,
		UpdatedBy:            &adminID,
	}

	if err := s.spawnRuleRepo.Create(ctx, rule); err != nil {
		return nil, fmt.Errorf("create spawn rule: %w", err)
	}

	s.audit(ctx, adminID, "create", "spawn_rule", &rule.ID, nil, rule)

	resp := toSpawnRuleDTO(rule)
	return &resp, nil
}

// UpdateSpawnRule updates an existing spawn rule.
func (s *CampAdminService) UpdateSpawnRule(ctx context.Context, adminID, id int64, req dto.UpdateSpawnRuleRequest) (*dto.SpawnRuleResponse, error) {
	rule, err := s.spawnRuleRepo.GetByID(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("get spawn rule: %w", err)
	}
	if rule == nil {
		return nil, ErrSpawnRuleNotFound
	}

	if req.Name != nil {
		rule.Name = *req.Name
	}
	if req.TerrainTypes != nil {
		j, _ := json.Marshal(*req.TerrainTypes)
		rule.TerrainTypesJSON = string(j)
	}
	if req.ZoneTypes != nil {
		j, _ := json.Marshal(*req.ZoneTypes)
		rule.ZoneTypesJSON = string(j)
	}
	if req.CampTemplatePool != nil {
		j, _ := json.Marshal(*req.CampTemplatePool)
		rule.CampTemplatePoolJSON = string(j)
	}
	if req.MaxCamps != nil {
		rule.MaxCamps = *req.MaxCamps
	}
	if req.SpawnIntervalSec != nil {
		rule.SpawnIntervalSec = *req.SpawnIntervalSec
	}
	if req.DespawnAfterSec != nil {
		rule.DespawnAfterSec = *req.DespawnAfterSec
	}
	if req.MinCampDistance != nil {
		rule.MinCampDistance = *req.MinCampDistance
	}
	if req.MinVillageDistance != nil {
		rule.MinVillageDistance = *req.MinVillageDistance
	}
	if req.Enabled != nil {
		rule.Enabled = *req.Enabled
	}
	rule.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	rule.UpdatedBy = &adminID

	if err := s.spawnRuleRepo.Update(ctx, rule); err != nil {
		return nil, fmt.Errorf("update spawn rule: %w", err)
	}

	s.audit(ctx, adminID, "update", "spawn_rule", &id, nil, rule)

	resp := toSpawnRuleDTO(rule)
	return &resp, nil
}

// DeleteSpawnRule deletes a spawn rule.
func (s *CampAdminService) DeleteSpawnRule(ctx context.Context, adminID, id int64) error {
	rule, err := s.spawnRuleRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get spawn rule: %w", err)
	}
	if rule == nil {
		return ErrSpawnRuleNotFound
	}
	if err := s.spawnRuleRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete spawn rule: %w", err)
	}
	s.audit(ctx, adminID, "delete", "spawn_rule", &id, rule, nil)
	return nil
}

// ── Reward Tables ────────────────────────────────────────────────────────────

// ListRewardTables returns all reward tables with their entries.
func (s *CampAdminService) ListRewardTables(ctx context.Context) ([]dto.RewardTableResponse, error) {
	tables, err := s.rewardTableRepo.GetAll(ctx)
	if err != nil {
		return nil, fmt.Errorf("list reward tables: %w", err)
	}

	result := make([]dto.RewardTableResponse, len(tables))
	for i, rt := range tables {
		entries, err := s.rewardEntryRepo.GetByRewardTableID(ctx, rt.ID)
		if err != nil {
			return nil, fmt.Errorf("get reward entries: %w", err)
		}
		entryDTOs := make([]dto.RewardEntryResponse, len(entries))
		for j, e := range entries {
			entryDTOs[j] = dto.RewardEntryResponse{
				ID:         e.ID,
				RewardType: e.RewardType,
				MinAmount:  e.MinAmount,
				MaxAmount:  e.MaxAmount,
				DropChance: int(e.DropChancePct),
			}
		}
		result[i] = dto.RewardTableResponse{
			ID:        rt.ID,
			Name:      rt.Name,
			Entries:   entryDTOs,
			CreatedAt: rt.CreatedAt,
			UpdatedAt: rt.UpdatedAt,
		}
	}
	return result, nil
}

// CreateRewardTable creates a new reward table with entries.
func (s *CampAdminService) CreateRewardTable(ctx context.Context, adminID int64, req dto.CreateRewardTableRequest) (*dto.RewardTableResponse, error) {
	now := time.Now().UTC().Format(time.RFC3339)
	rt := &model.RewardTable{
		Name:      req.Name,
		CreatedAt: now,
		UpdatedAt: now,
		UpdatedBy: &adminID,
	}

	if err := s.rewardTableRepo.Create(ctx, rt); err != nil {
		return nil, fmt.Errorf("create reward table: %w", err)
	}

	var entryDTOs []dto.RewardEntryResponse
	for _, entryReq := range req.Entries {
		entry := &model.RewardTableEntry{
			RewardTableID: rt.ID,
			RewardType:    entryReq.RewardType,
			MinAmount:     entryReq.MinAmount,
			MaxAmount:     entryReq.MaxAmount,
			DropChancePct: float64(entryReq.DropChance),
			CreatedAt:     now,
		}
		if err := s.rewardEntryRepo.Create(ctx, entry); err != nil {
			return nil, fmt.Errorf("create reward entry: %w", err)
		}
		entryDTOs = append(entryDTOs, dto.RewardEntryResponse{
			ID:         entry.ID,
			RewardType: entry.RewardType,
			MinAmount:  entry.MinAmount,
			MaxAmount:  entry.MaxAmount,
			DropChance: entryReq.DropChance,
		})
	}

	s.audit(ctx, adminID, "create", "reward_table", &rt.ID, nil, rt)

	return &dto.RewardTableResponse{
		ID:        rt.ID,
		Name:      rt.Name,
		Entries:   entryDTOs,
		CreatedAt: rt.CreatedAt,
		UpdatedAt: rt.UpdatedAt,
	}, nil
}

// DeleteRewardTable deletes a reward table and its entries.
func (s *CampAdminService) DeleteRewardTable(ctx context.Context, adminID, id int64) error {
	rt, err := s.rewardTableRepo.GetByID(ctx, id)
	if err != nil {
		return fmt.Errorf("get reward table: %w", err)
	}
	if rt == nil {
		return ErrRewardTableNotFound
	}
	if err := s.rewardEntryRepo.DeleteByRewardTableID(ctx, id); err != nil {
		return fmt.Errorf("delete reward entries: %w", err)
	}
	if err := s.rewardTableRepo.Delete(ctx, id); err != nil {
		return fmt.Errorf("delete reward table: %w", err)
	}
	s.audit(ctx, adminID, "delete", "reward_table", &id, rt, nil)
	return nil
}

// ── Battle Tuning ────────────────────────────────────────────────────────────

// GetBattleTuning returns the current battle tuning configuration.
func (s *CampAdminService) GetBattleTuning(ctx context.Context) (*dto.BattleTuningResponse, error) {
	tuning, err := s.battleTuningRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get battle tuning: %w", err)
	}
	return &dto.BattleTuningResponse{
		TickDurationMs:        tuning.TickDurationMs,
		CritDamageMultiplier:  tuning.CritDamageMultiplier,
		MaxDefensePercent:     tuning.MaxDefensePercent,
		MaxCritChancePercent:  tuning.MaxCritChancePercent,
		MinAttackInterval:     tuning.MinAttackInterval,
		MarchSpeedTilesPerMin: tuning.MarchSpeedTilesPerMin,
		MaxTicks:              tuning.MaxTicks,
		UpdatedAt:             tuning.UpdatedAt,
	}, nil
}

// UpdateBattleTuning updates the battle tuning configuration.
func (s *CampAdminService) UpdateBattleTuning(ctx context.Context, adminID int64, req dto.UpdateBattleTuningRequest) (*dto.BattleTuningResponse, error) {
	tuning, err := s.battleTuningRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get battle tuning: %w", err)
	}

	if req.TickDurationMs != nil {
		tuning.TickDurationMs = *req.TickDurationMs
	}
	if req.CritDamageMultiplier != nil {
		tuning.CritDamageMultiplier = *req.CritDamageMultiplier
	}
	if req.MaxDefensePercent != nil {
		tuning.MaxDefensePercent = *req.MaxDefensePercent
	}
	if req.MaxCritChancePercent != nil {
		tuning.MaxCritChancePercent = *req.MaxCritChancePercent
	}
	if req.MinAttackInterval != nil {
		tuning.MinAttackInterval = *req.MinAttackInterval
	}
	if req.MarchSpeedTilesPerMin != nil {
		tuning.MarchSpeedTilesPerMin = *req.MarchSpeedTilesPerMin
	}
	if req.MaxTicks != nil {
		tuning.MaxTicks = *req.MaxTicks
	}
	tuning.UpdatedAt = time.Now().UTC().Format(time.RFC3339)
	tuning.UpdatedBy = &adminID

	if err := s.battleTuningRepo.Update(ctx, tuning); err != nil {
		return nil, fmt.Errorf("update battle tuning: %w", err)
	}

	s.audit(ctx, adminID, "update", "battle_tuning", nil, nil, tuning)

	return &dto.BattleTuningResponse{
		TickDurationMs:        tuning.TickDurationMs,
		CritDamageMultiplier:  tuning.CritDamageMultiplier,
		MaxDefensePercent:     tuning.MaxDefensePercent,
		MaxCritChancePercent:  tuning.MaxCritChancePercent,
		MinAttackInterval:     tuning.MinAttackInterval,
		MarchSpeedTilesPerMin: tuning.MarchSpeedTilesPerMin,
		MaxTicks:              tuning.MaxTicks,
		UpdatedAt:             tuning.UpdatedAt,
	}, nil
}

// ── Helpers ──────────────────────────────────────────────────────────────────

func (s *CampAdminService) audit(ctx context.Context, adminID int64, action, entityType string, entityID *int64, oldVal, newVal interface{}) {
	oldJSON, _ := json.Marshal(oldVal)
	newJSON, _ := json.Marshal(newVal)

	entry := &model.AdminAuditLog{
		AdminPlayerID: adminID,
		Action:        action,
		EntityType:    entityType,
		EntityID:      entityID,
		OldValueJSON:  string(oldJSON),
		NewValueJSON:  string(newJSON),
		CreatedAt:     time.Now().UTC().Format(time.RFC3339),
	}
	s.auditLogRepo.Create(ctx, entry)
}

func toBeastTemplateDTO(bt *model.BeastTemplate) dto.BeastTemplateResponse {
	return dto.BeastTemplateResponse{
		ID:                bt.ID,
		Name:              bt.Name,
		SpriteKey:         bt.SpriteKey,
		HP:                bt.HP,
		AttackPower:       bt.AttackPower,
		AttackInterval:    bt.AttackInterval,
		DefensePercent:    bt.DefensePercent,
		CritChancePercent: bt.CritChancePercent,
		Description:       bt.Description,
		CreatedAt:         bt.CreatedAt,
		UpdatedAt:         bt.UpdatedAt,
	}
}

func toSpawnRuleDTO(r *model.SpawnRule) dto.SpawnRuleResponse {
	var terrainTypes []string
	json.Unmarshal([]byte(r.TerrainTypesJSON), &terrainTypes)
	var zoneTypes []string
	json.Unmarshal([]byte(r.ZoneTypesJSON), &zoneTypes)
	var pool []dto.CampTemplatePoolEntry
	json.Unmarshal([]byte(r.CampTemplatePoolJSON), &pool)

	return dto.SpawnRuleResponse{
		ID:                 r.ID,
		Name:               r.Name,
		TerrainTypes:       terrainTypes,
		ZoneTypes:          zoneTypes,
		CampTemplatePool:   pool,
		MaxCamps:           r.MaxCamps,
		SpawnIntervalSec:   r.SpawnIntervalSec,
		DespawnAfterSec:    r.DespawnAfterSec,
		MinCampDistance:    r.MinCampDistance,
		MinVillageDistance: r.MinVillageDistance,
		Enabled:            r.Enabled,
		CreatedAt:          r.CreatedAt,
		UpdatedAt:          r.UpdatedAt,
	}
}

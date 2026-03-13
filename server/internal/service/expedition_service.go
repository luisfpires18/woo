package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/luisfpires18/woo/internal/battle"
	"github.com/luisfpires18/woo/internal/config"
	"github.com/luisfpires18/woo/internal/dto"
	"github.com/luisfpires18/woo/internal/model"
	"github.com/luisfpires18/woo/internal/repository"
)

// Expedition service errors.
var (
	ErrCampNotFound       = errors.New("camp not found")
	ErrCampNotActive      = errors.New("camp is not active")
	ErrNoTroopsSent       = errors.New("must send at least one troop")
	ErrInsufficientTroops = errors.New("not enough troops in village")
	ErrExpeditionNotFound = errors.New("expedition not found")
	ErrBattleNotFound     = errors.New("battle not found")
)

// ExpeditionCompletionEvent describes an expedition that has been resolved.
type ExpeditionCompletionEvent struct {
	PlayerID     int64
	VillageID    int64
	ExpeditionID int64
	CampID       int64
	Result       string // attacker_won, defender_won, draw
}

// ExpeditionReturnEvent describes troops returning home.
type ExpeditionReturnEvent struct {
	PlayerID     int64
	VillageID    int64
	ExpeditionID int64
}

// ExpeditionService handles dispatching troops, resolving battles, and returning troops.
type ExpeditionService struct {
	uow              repository.UnitOfWork
	expeditionRepo   repository.ExpeditionRepository
	campRepo         repository.CampRepository
	battleRepo       repository.BattleRepository
	battleTuningRepo repository.BattleTuningRepository
	troopRepo        repository.TroopRepository
	villageRepo      repository.VillageRepository
	rewardTableRepo  repository.RewardTableRepository
	rewardEntryRepo  repository.RewardTableEntryRepository
	campTemplateRepo repository.CampTemplateRepository
	resourceRepo     repository.ResourceRepository
}

// NewExpeditionService creates a new ExpeditionService.
func NewExpeditionService(
	uow repository.UnitOfWork,
	expeditionRepo repository.ExpeditionRepository,
	campRepo repository.CampRepository,
	battleRepo repository.BattleRepository,
	battleTuningRepo repository.BattleTuningRepository,
	troopRepo repository.TroopRepository,
	villageRepo repository.VillageRepository,
	rewardTableRepo repository.RewardTableRepository,
	rewardEntryRepo repository.RewardTableEntryRepository,
	campTemplateRepo repository.CampTemplateRepository,
	resourceRepo repository.ResourceRepository,
) *ExpeditionService {
	return &ExpeditionService{
		uow:              uow,
		expeditionRepo:   expeditionRepo,
		campRepo:         campRepo,
		battleRepo:       battleRepo,
		battleTuningRepo: battleTuningRepo,
		troopRepo:        troopRepo,
		villageRepo:      villageRepo,
		rewardTableRepo:  rewardTableRepo,
		rewardEntryRepo:  rewardEntryRepo,
		campTemplateRepo: campTemplateRepo,
		resourceRepo:     resourceRepo,
	}
}

// DispatchExpedition sends troops from a village to attack a camp.
func (s *ExpeditionService) DispatchExpedition(ctx context.Context, playerID, villageID int64, req dto.DispatchExpeditionRequest) (*dto.ExpeditionResponse, error) {
	// 1. Validate village ownership
	village, err := s.villageRepo.GetByID(ctx, villageID)
	if err != nil {
		return nil, fmt.Errorf("get village: %w", err)
	}
	if village == nil {
		return nil, ErrVillageNotFound
	}
	if village.PlayerID != playerID {
		return nil, ErrNotOwner
	}

	// 2. Validate camp exists and is active
	camp, err := s.campRepo.GetByID(ctx, req.CampID)
	if err != nil {
		return nil, fmt.Errorf("get camp: %w", err)
	}
	if camp == nil {
		return nil, ErrCampNotFound
	}
	if camp.Status != model.CampStatusActive {
		return nil, ErrCampNotActive
	}

	// 3. Validate troops
	if len(req.Troops) == 0 {
		return nil, ErrNoTroopsSent
	}

	// Build troop deduction map and snapshot
	troopDeductions := make(map[string]int)
	var troopSnapshot []model.ExpeditionTroop
	var slowestSpeed int = math.MaxInt32

	for _, td := range req.Troops {
		if td.Quantity <= 0 {
			continue
		}

		troopCfg, ok := config.TroopConfigs[td.TroopType]
		if !ok {
			return nil, fmt.Errorf("unknown troop type: %s", td.TroopType)
		}

		// Check available troops
		troop, err := s.troopRepo.GetByVillageAndType(ctx, villageID, td.TroopType)
		if err != nil {
			return nil, fmt.Errorf("get troop: %w", err)
		}
		if troop == nil || troop.Quantity < td.Quantity {
			return nil, ErrInsufficientTroops
		}

		troopDeductions[td.TroopType] = td.Quantity

		troopSnapshot = append(troopSnapshot, model.ExpeditionTroop{
			TroopType:         td.TroopType,
			Quantity:          td.Quantity,
			OriginalQuantity:  td.Quantity,
			HP:                troopCfg.Attack, // use attack as base HP for troops
			AttackPower:       troopCfg.Attack,
			AttackInterval:    2, // default tick interval for troops
			DefensePercent:    float64(troopCfg.DefInfantry),
			CritChancePercent: 0,
			Speed:             troopCfg.Speed,
			Carry:             troopCfg.Carry,
		})

		if troopCfg.Speed < slowestSpeed && troopCfg.Speed > 0 {
			slowestSpeed = troopCfg.Speed
		}
	}

	if len(troopDeductions) == 0 {
		return nil, ErrNoTroopsSent
	}

	// 4. Calculate travel time based on distance and speed
	tuning, err := s.battleTuningRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get battle tuning: %w", err)
	}

	distance := tileDistance(village.X, village.Y, camp.TileX, camp.TileY)
	marchSpeed := tuning.MarchSpeedTilesPerMin
	if marchSpeed <= 0 {
		marchSpeed = 1.0
	}
	travelMinutes := distance / marchSpeed
	travelDuration := time.Duration(travelMinutes * float64(time.Minute))
	if travelDuration < 10*time.Second {
		travelDuration = 10 * time.Second // minimum travel time
	}

	now := time.Now().UTC()
	arrivesAt := now.Add(travelDuration)

	troopsJSON, err := json.Marshal(troopSnapshot)
	if err != nil {
		return nil, fmt.Errorf("marshal troops: %w", err)
	}

	expedition := &model.Expedition{
		PlayerID:   playerID,
		VillageID:  villageID,
		CampID:     req.CampID,
		TroopsJSON: string(troopsJSON),
		DepartedAt: now.Format(time.RFC3339),
		ArrivesAt:  arrivesAt.Format(time.RFC3339),
		Status:     model.ExpeditionMarching,
	}

	// 5. Atomically deduct troops and create expedition
	if err := s.uow.DeductTroopsAndCreateExpedition(ctx, villageID, troopDeductions, expedition); err != nil {
		return nil, fmt.Errorf("dispatch expedition: %w", err)
	}

	// 6. Mark camp as under attack
	if err := s.campRepo.UpdateStatus(ctx, req.CampID, model.CampStatusUnderAttack); err != nil {
		return nil, fmt.Errorf("update camp status: %w", err)
	}

	return &dto.ExpeditionResponse{
		ID:           expedition.ID,
		VillageID:    villageID,
		CampID:       req.CampID,
		Troops:       toExpeditionTroopDTOs(troopSnapshot),
		Status:       model.ExpeditionMarching,
		DispatchedAt: now,
		ArrivesAt:    arrivesAt,
	}, nil
}

// ResolveArrivedExpeditions processes expeditions that have arrived at their destination.
// Called periodically by the game loop.
func (s *ExpeditionService) ResolveArrivedExpeditions(ctx context.Context) ([]ExpeditionCompletionEvent, error) {
	now := time.Now().UTC()
	arrived, err := s.expeditionRepo.GetArrivedExpeditions(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("get arrived expeditions: %w", err)
	}

	var events []ExpeditionCompletionEvent
	for _, exp := range arrived {
		event, err := s.resolveExpedition(ctx, exp)
		if err != nil {
			continue
		}
		events = append(events, *event)
	}
	return events, nil
}

// ReturnCompletedExpeditions processes expeditions that have returned home.
// Called periodically by the game loop.
func (s *ExpeditionService) ReturnCompletedExpeditions(ctx context.Context) ([]ExpeditionReturnEvent, error) {
	now := time.Now().UTC()
	returning, err := s.expeditionRepo.GetReturningExpeditions(ctx, now)
	if err != nil {
		return nil, fmt.Errorf("get returning expeditions: %w", err)
	}

	var events []ExpeditionReturnEvent
	for _, exp := range returning {
		event, err := s.returnExpedition(ctx, exp)
		if err != nil {
			continue
		}
		events = append(events, *event)
	}
	return events, nil
}

// GetExpeditionsByPlayer returns all expeditions for a player.
func (s *ExpeditionService) GetExpeditionsByPlayer(ctx context.Context, playerID int64) ([]dto.ExpeditionResponse, error) {
	exps, err := s.expeditionRepo.GetByPlayerID(ctx, playerID)
	if err != nil {
		return nil, fmt.Errorf("get expeditions: %w", err)
	}

	var result []dto.ExpeditionResponse
	for _, exp := range exps {
		var troops []model.ExpeditionTroop
		if err := json.Unmarshal([]byte(exp.TroopsJSON), &troops); err != nil {
			continue
		}

		resp := dto.ExpeditionResponse{
			ID:        exp.ID,
			VillageID: exp.VillageID,
			CampID:    exp.CampID,
			Troops:    toExpeditionTroopDTOs(troops),
			Status:    exp.Status,
		}

		if t, err := time.Parse(time.RFC3339, exp.DepartedAt); err == nil {
			resp.DispatchedAt = t
		}
		if t, err := time.Parse(time.RFC3339, exp.ArrivesAt); err == nil {
			resp.ArrivesAt = t
		}
		if exp.ReturnsAt != "" {
			if t, err := time.Parse(time.RFC3339, exp.ReturnsAt); err == nil {
				resp.ReturnAt = &t
			}
		}

		// Look up battle for non-marching expeditions
		if exp.Status != model.ExpeditionMarching {
			if b, err := s.battleRepo.GetByExpeditionID(ctx, exp.ID); err == nil && b != nil {
				resp.BattleID = &b.ID
				if exp.Status == model.ExpeditionCompleted {
					if t, err := time.Parse(time.RFC3339, b.ResolvedAt); err == nil {
						resp.CompletedAt = &t
					}
				}
			}
		}

		result = append(result, resp)
	}
	return result, nil
}

// GetBattleReport returns the battle report for an expedition.
func (s *ExpeditionService) GetBattleReport(ctx context.Context, playerID int64, battleID int64) (*dto.BattleReportResponse, error) {
	b, err := s.battleRepo.GetByID(ctx, battleID)
	if err != nil {
		return nil, fmt.Errorf("get battle: %w", err)
	}
	if b == nil {
		return nil, ErrBattleNotFound
	}

	// Verify the player owns this expedition
	exp, err := s.expeditionRepo.GetByID(ctx, b.ExpeditionID)
	if err != nil || exp == nil {
		return nil, ErrExpeditionNotFound
	}
	if exp.PlayerID != playerID {
		return nil, ErrNotOwner
	}

	// Unmarshal per-unit losses from DB and aggregate
	var rawAttackerLosses []storedUnitLoss
	if b.AttackerLossesJSON != "" {
		json.Unmarshal([]byte(b.AttackerLossesJSON), &rawAttackerLosses)
	}
	var rawDefenderLosses []storedUnitLoss
	if b.DefenderLossesJSON != "" {
		json.Unmarshal([]byte(b.DefenderLossesJSON), &rawDefenderLosses)
	}

	attackerLosses := aggregateLosses(rawAttackerLosses)
	defenderLosses := aggregateLosses(rawDefenderLosses)

	var rawRewards []model.BattleReward
	if b.RewardsJSON != "" {
		json.Unmarshal([]byte(b.RewardsJSON), &rawRewards)
	}
	rewards := make([]dto.BattleRewardResponse, len(rawRewards))
	for i, r := range rawRewards {
		rewards[i] = dto.BattleRewardResponse{
			ResourceType: r.RewardType,
			Amount:       r.Amount,
		}
	}

	foughtAt, _ := time.Parse(time.RFC3339, b.ResolvedAt)

	return &dto.BattleReportResponse{
		ID:             b.ID,
		ExpeditionID:   b.ExpeditionID,
		CampID:         exp.CampID,
		Result:         b.Result,
		AttackerLosses: attackerLosses,
		DefenderLosses: defenderLosses,
		Rewards:        rewards,
		FoughtAt:       foughtAt,
	}, nil
}

// GetBattleReplay returns the raw replay data blob for a battle.
func (s *ExpeditionService) GetBattleReplay(ctx context.Context, playerID int64, battleID int64) ([]byte, error) {
	b, err := s.battleRepo.GetByID(ctx, battleID)
	if err != nil {
		return nil, fmt.Errorf("get battle: %w", err)
	}
	if b == nil {
		return nil, ErrBattleNotFound
	}

	exp, err := s.expeditionRepo.GetByID(ctx, b.ExpeditionID)
	if err != nil || exp == nil {
		return nil, ErrExpeditionNotFound
	}
	if exp.PlayerID != playerID {
		return nil, ErrNotOwner
	}

	return s.battleRepo.GetReplayData(ctx, battleID)
}

// ── Internal battle resolution ───────────────────────────────────────────────

func (s *ExpeditionService) resolveExpedition(ctx context.Context, exp *model.Expedition) (*ExpeditionCompletionEvent, error) {
	// Mark as battling
	if err := s.expeditionRepo.UpdateStatus(ctx, exp.ID, model.ExpeditionBattling); err != nil {
		return nil, fmt.Errorf("update expedition status: %w", err)
	}

	// Load camp and beasts
	camp, err := s.campRepo.GetByID(ctx, exp.CampID)
	if err != nil || camp == nil {
		return nil, fmt.Errorf("get camp for battle: %w", err)
	}

	var campBeasts []model.CampBeast
	if err := json.Unmarshal([]byte(camp.BeastsJSON), &campBeasts); err != nil {
		return nil, fmt.Errorf("parse camp beasts: %w", err)
	}

	// Load troop snapshot
	var troops []model.ExpeditionTroop
	if err := json.Unmarshal([]byte(exp.TroopsJSON), &troops); err != nil {
		return nil, fmt.Errorf("parse expedition troops: %w", err)
	}

	// Load tuning
	tuning, err := s.battleTuningRepo.Get(ctx)
	if err != nil {
		return nil, fmt.Errorf("get battle tuning: %w", err)
	}

	// Convert to battle units
	attackers := troopsToUnits(troops)
	defenders := beastsToUnits(campBeasts)

	engineTuning := battle.Tuning{
		CritDamageMultiplier: tuning.CritDamageMultiplier,
		MaxDefensePercent:    tuning.MaxDefensePercent,
		MaxCritChancePercent: tuning.MaxCritChancePercent,
		MinAttackInterval:    tuning.MinAttackInterval,
		MaxTicks:             tuning.MaxTicks,
	}

	// Run deterministic battle
	seed := time.Now().UnixNano() ^ exp.ID
	attackerCopy := cloneBattleUnits(attackers)
	defenderCopy := cloneBattleUnits(defenders)
	result := battle.Simulate(attackerCopy, defenderCopy, engineTuning, seed)

	// Build replay data
	replayData, _ := battle.BuildReplayJSON(attackers, defenders, result, tuning.TickDurationMs)

	// Calculate losses (stored as JSON in DB with per-unit detail)
	attackerLosses := calculateTroopLossesForStorage(troops, result.AttackerSurvivors)
	defenderLosses := calculateBeastLossesForStorage(campBeasts, result.DefenderSurvivors)

	// Roll rewards if attacker won
	var rewards []model.BattleReward
	if result.Outcome == "attacker_won" {
		rewards = s.rollRewards(ctx, camp.CampTemplateID)
	}

	attackerLossesJSON, _ := json.Marshal(attackerLosses)
	defenderLossesJSON, _ := json.Marshal(defenderLosses)
	rewardsJSON, _ := json.Marshal(rewards)
	attackerSnapJSON, _ := json.Marshal(attackers)
	defenderSnapJSON, _ := json.Marshal(defenders)

	now := time.Now().UTC()
	battleModel := &model.Battle{
		ExpeditionID:         exp.ID,
		AttackerSnapshotJSON: string(attackerSnapJSON),
		DefenderSnapshotJSON: string(defenderSnapJSON),
		Result:               result.Outcome,
		AttackerLossesJSON:   string(attackerLossesJSON),
		DefenderLossesJSON:   string(defenderLossesJSON),
		RewardsJSON:          string(rewardsJSON),
		ReplayData:           replayData,
		Seed:                 seed,
		ResolvedAt:           now.Format(time.RFC3339),
		DurationTicks:        result.DurationTicks,
	}

	if err := s.battleRepo.Create(ctx, battleModel); err != nil {
		return nil, fmt.Errorf("create battle record: %w", err)
	}

	// Calculate return time (same as travel time)
	village, err := s.villageRepo.GetByID(ctx, exp.VillageID)
	if err != nil {
		return nil, fmt.Errorf("get village for return: %w", err)
	}

	distance := tileDistance(village.X, village.Y, camp.TileX, camp.TileY)
	marchSpeed := tuning.MarchSpeedTilesPerMin
	if marchSpeed <= 0 {
		marchSpeed = 1.0
	}
	returnMinutes := distance / marchSpeed
	returnDuration := time.Duration(returnMinutes * float64(time.Minute))
	if returnDuration < 10*time.Second {
		returnDuration = 10 * time.Second
	}
	returnsAt := now.Add(returnDuration)

	// Update survivor troops in expedition JSON for return
	survivorTroops := survivorTroopsFromResult(troops, result.AttackerSurvivors)
	survivorJSON, _ := json.Marshal(survivorTroops)

	exp.TroopsJSON = string(survivorJSON)
	exp.Status = model.ExpeditionReturning
	exp.ReturnsAt = returnsAt.Format(time.RFC3339)
	if err := s.expeditionRepo.Update(ctx, exp); err != nil {
		return nil, fmt.Errorf("update expedition for return: %w", err)
	}

	// Update camp status
	if result.Outcome == "attacker_won" {
		s.campRepo.UpdateStatus(ctx, camp.ID, model.CampStatusCleared)
	} else {
		s.campRepo.UpdateStatus(ctx, camp.ID, model.CampStatusActive)
	}

	// Apply resource rewards to village
	if result.Outcome == "attacker_won" && len(rewards) > 0 {
		s.applyRewards(ctx, exp.VillageID, rewards)
	}

	return &ExpeditionCompletionEvent{
		PlayerID:     exp.PlayerID,
		VillageID:    exp.VillageID,
		ExpeditionID: exp.ID,
		CampID:       exp.CampID,
		Result:       result.Outcome,
	}, nil
}

func (s *ExpeditionService) returnExpedition(ctx context.Context, exp *model.Expedition) (*ExpeditionReturnEvent, error) {
	// Parse surviving troops
	var troops []model.ExpeditionTroop
	if err := json.Unmarshal([]byte(exp.TroopsJSON), &troops); err != nil {
		return nil, fmt.Errorf("parse returning troops: %w", err)
	}

	// Build troop additions map
	troopAdditions := make(map[string]int)
	for _, t := range troops {
		if t.Quantity > 0 {
			troopAdditions[t.TroopType] += t.Quantity
		}
	}

	// Atomically return troops and mark expedition completed
	if err := s.uow.ReturnExpeditionTroops(ctx, exp.VillageID, troopAdditions, exp.ID); err != nil {
		return nil, fmt.Errorf("return expedition troops: %w", err)
	}

	return &ExpeditionReturnEvent{
		PlayerID:     exp.PlayerID,
		VillageID:    exp.VillageID,
		ExpeditionID: exp.ID,
	}, nil
}

func (s *ExpeditionService) rollRewards(ctx context.Context, campTemplateID int64) []model.BattleReward {
	template, err := s.campTemplateRepo.GetByID(ctx, campTemplateID)
	if err != nil || template == nil || template.RewardTableID == nil {
		return nil
	}

	entries, err := s.rewardEntryRepo.GetByRewardTableID(ctx, *template.RewardTableID)
	if err != nil {
		return nil
	}

	rng := rand.New(rand.NewSource(time.Now().UnixNano()))
	var rewards []model.BattleReward

	for _, entry := range entries {
		// Roll drop chance
		if rng.Float64()*100 > entry.DropChancePct {
			continue
		}
		amount := entry.MinAmount
		if entry.MaxAmount > entry.MinAmount {
			amount = entry.MinAmount + rng.Intn(entry.MaxAmount-entry.MinAmount+1)
		}
		if amount > 0 {
			rewards = append(rewards, model.BattleReward{
				RewardType: entry.RewardType,
				Amount:     amount,
			})
		}
	}

	return rewards
}

func (s *ExpeditionService) applyRewards(ctx context.Context, villageID int64, rewards []model.BattleReward) {
	res, err := s.resourceRepo.Get(ctx, villageID)
	if err != nil {
		return
	}

	for _, reward := range rewards {
		switch reward.RewardType {
		case "food":
			res.Food += float64(reward.Amount)
		case "water":
			res.Water += float64(reward.Amount)
		case "lumber":
			res.Lumber += float64(reward.Amount)
		case "stone":
			res.Stone += float64(reward.Amount)
		}
	}

	s.resourceRepo.Update(ctx, villageID, res)
}

// ── Conversion helpers ───────────────────────────────────────────────────────

func tileDistance(x1, y1, x2, y2 int) float64 {
	dx := float64(x1 - x2)
	dy := float64(y1 - y2)
	return math.Sqrt(dx*dx + dy*dy)
}

func troopsToUnits(troops []model.ExpeditionTroop) []battle.Unit {
	var units []battle.Unit
	id := 1
	for _, t := range troops {
		for i := 0; i < t.Quantity; i++ {
			units = append(units, battle.Unit{
				ID:                id,
				Side:              "attacker",
				Name:              t.TroopType,
				HP:                t.HP,
				MaxHP:             t.HP,
				AttackPower:       t.AttackPower,
				AttackInterval:    t.AttackInterval,
				DefensePercent:    t.DefensePercent,
				CritChancePercent: t.CritChancePercent,
			})
			id++
		}
	}
	return units
}

func beastsToUnits(beasts []model.CampBeast) []battle.Unit {
	var units []battle.Unit
	id := 10000 // offset to avoid ID conflicts with attackers
	for _, b := range beasts {
		units = append(units, battle.Unit{
			ID:                id,
			Side:              "defender",
			Name:              b.Name,
			SpriteKey:         b.SpriteKey,
			HP:                b.HP,
			MaxHP:             b.MaxHP,
			AttackPower:       b.AttackPower,
			AttackInterval:    b.AttackInterval,
			DefensePercent:    b.DefensePercent,
			CritChancePercent: b.CritChancePercent,
		})
		id++
	}
	return units
}

func cloneBattleUnits(units []battle.Unit) []battle.Unit {
	c := make([]battle.Unit, len(units))
	copy(c, units)
	return c
}

func calculateTroopLossesForStorage(original []model.ExpeditionTroop, survivors []battle.Unit) []storedUnitLoss {
	// Count surviving units by troop type
	survByCType := make(map[string]int)
	for _, u := range survivors {
		survByCType[u.Name]++
	}

	var losses []storedUnitLoss
	for _, t := range original {
		survived := survByCType[t.TroopType]
		sent := t.OriginalQuantity
		if sent == 0 {
			sent = t.Quantity // backward compat
		}
		lost := sent - survived
		losses = append(losses, storedUnitLoss{
			TroopType: t.TroopType,
			Sent:      sent,
			Lost:      lost,
		})
	}
	return losses
}

func calculateBeastLossesForStorage(original []model.CampBeast, survivors []battle.Unit) []storedUnitLoss {
	survByName := make(map[string]int)
	for _, u := range survivors {
		survByName[u.Name]++
	}

	// Count original by name
	origByName := make(map[string]int)
	for _, b := range original {
		origByName[b.Name]++
	}

	var losses []storedUnitLoss
	for name, count := range origByName {
		survived := survByName[name]
		lost := count - survived
		losses = append(losses, storedUnitLoss{
			TroopType: name,
			Sent:      count,
			Lost:      lost,
		})
	}
	return losses
}

// storedUnitLoss is the per-unit-type loss format stored as JSON in the DB.
type storedUnitLoss struct {
	TroopType string `json:"troop_type"`
	Sent      int    `json:"sent"`
	Lost      int    `json:"lost"`
}

func aggregateLosses(losses []storedUnitLoss) dto.BattleLosses {
	var totalSent, totalLost int
	for _, l := range losses {
		totalSent += l.Sent
		totalLost += l.Lost
	}
	return dto.BattleLosses{
		TotalSent:     totalSent,
		TotalLost:     totalLost,
		TotalSurvived: totalSent - totalLost,
	}
}

func survivorTroopsFromResult(original []model.ExpeditionTroop, survivors []battle.Unit) []model.ExpeditionTroop {
	// Count surviving by troop type
	survCount := make(map[string]int)
	for _, u := range survivors {
		survCount[u.Name]++
	}

	var result []model.ExpeditionTroop
	for _, t := range original {
		survived := survCount[t.TroopType]
		origQty := t.OriginalQuantity
		if origQty == 0 {
			origQty = t.Quantity // backward compat
		}
		// Include all troop types (even those with 0 survivors) so we can report them
		t.Quantity = survived
		t.OriginalQuantity = origQty
		result = append(result, t)
	}
	return result
}

func toExpeditionTroopDTOs(troops []model.ExpeditionTroop) []dto.ExpeditionTroopResponse {
	result := make([]dto.ExpeditionTroopResponse, len(troops))
	for i, t := range troops {
		sent := t.OriginalQuantity
		if sent == 0 {
			sent = t.Quantity // backward compat for pre-existing expeditions
		}
		result[i] = dto.ExpeditionTroopResponse{
			TroopType:        t.TroopType,
			QuantitySent:     sent,
			QuantitySurvived: t.Quantity,
		}
	}
	return result
}

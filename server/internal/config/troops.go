package config

import "math"

// TroopConfig holds base stats and training data for a troop type.
type TroopConfig struct {
	DisplayName      string
	Kingdom          string       // kingdom slug — empty means all kingdoms
	TrainingBuilding string       // building type required to train (e.g. "barracks", "stable")
	BuildingLevelReq int          // minimum level of the training building
	BaseCost         ResourceCost // cost per single unit
	BaseTimeSec      int          // training time per unit in seconds (at barracks level 1)
	FoodUpkeep       float64      // food consumption per unit per hour
	Attack           int
	DefInfantry      int // defense against infantry
	DefCavalry       int // defense against cavalry
	Speed            int // tiles per hour
	Carry            int // resource carrying capacity per unit
}

// TroopConfigs is the authoritative registry of all troop types and their stats.
// Currently contains Arkazia troops only. Veridor and Sylvara to be added later.
var TroopConfigs = map[string]TroopConfig{
	// --- Arkazia (Mountain / Gladiator / Iron) ---
	"iron_legionary": {
		DisplayName:      "Iron Legionary",
		Kingdom:          "arkazia",
		TrainingBuilding: "barracks",
		BuildingLevelReq: 1,
		BaseCost:         ResourceCost{Food: 100, Water: 50, Lumber: 60, Stone: 40},
		BaseTimeSec:      120, // 2 min
		FoodUpkeep:       1,
		Attack:           50,
		DefInfantry:      55,
		DefCavalry:       40,
		Speed:            5,
		Carry:            55,
	},
	"crossbowman": {
		DisplayName:      "Crossbowman",
		Kingdom:          "arkazia",
		TrainingBuilding: "barracks",
		BuildingLevelReq: 3,
		BaseCost:         ResourceCost{Food: 80, Water: 60, Lumber: 100, Stone: 30},
		BaseTimeSec:      150, // 2.5 min
		FoodUpkeep:       1,
		Attack:           50,
		DefInfantry:      30,
		DefCavalry:       25,
		Speed:            5,
		Carry:            35,
	},
	"mountain_knight": {
		DisplayName:      "Mountain Knight",
		Kingdom:          "arkazia",
		TrainingBuilding: "stable",
		BuildingLevelReq: 1,
		BaseCost:         ResourceCost{Food: 200, Water: 100, Lumber: 150, Stone: 80},
		BaseTimeSec:      300, // 5 min
		FoodUpkeep:       3,
		Attack:           80,
		DefInfantry:      40,
		DefCavalry:       50,
		Speed:            10,
		Carry:            90,
	},
	"shieldbearer": {
		DisplayName:      "Shieldbearer",
		Kingdom:          "arkazia",
		TrainingBuilding: "barracks",
		BuildingLevelReq: 5,
		BaseCost:         ResourceCost{Food: 120, Water: 80, Lumber: 60, Stone: 120},
		BaseTimeSec:      180, // 3 min
		FoodUpkeep:       2,
		Attack:           20,
		DefInfantry:      80,
		DefCavalry:       70,
		Speed:            3,
		Carry:            10,
	},
	"gladiator": {
		DisplayName:      "Gladiator",
		Kingdom:          "arkazia",
		TrainingBuilding: "special",
		BuildingLevelReq: 1,
		BaseCost:         ResourceCost{Food: 300, Water: 150, Lumber: 200, Stone: 200},
		BaseTimeSec:      600, // 10 min
		FoodUpkeep:       3,
		Attack:           95,
		DefInfantry:      50,
		DefCavalry:       45,
		Speed:            7,
		Carry:            60,
	},
	"battering_ram": {
		DisplayName:      "Battering Ram",
		Kingdom:          "arkazia",
		TrainingBuilding: "barracks",
		BuildingLevelReq: 10,
		BaseCost:         ResourceCost{Food: 250, Water: 80, Lumber: 300, Stone: 200},
		BaseTimeSec:      480, // 8 min
		FoodUpkeep:       4,
		Attack:           70,
		DefInfantry:      25,
		DefCavalry:       10,
		Speed:            3,
		Carry:            0,
	},
	"mountain_scout": {
		DisplayName:      "Mountain Scout",
		Kingdom:          "arkazia",
		TrainingBuilding: "barracks",
		BuildingLevelReq: 1,
		BaseCost:         ResourceCost{Food: 60, Water: 40, Lumber: 30, Stone: 20},
		BaseTimeSec:      60, // 1 min
		FoodUpkeep:       1,
		Attack:           12,
		DefInfantry:      10,
		DefCavalry:       10,
		Speed:            12,
		Carry:            0,
	},
}

// TrainingSpeedMultiplier returns the training speed multiplier for a given
// training building level. Higher levels = faster training.
// Based on progression-and-scaling.md: lv1=1.0, lv5=1.25, lv10=1.6, lv15=2.0, lv20=2.5
func TrainingSpeedMultiplier(buildingLevel int) float64 {
	if buildingLevel <= 0 {
		return 0 // building not built — can't train
	}
	if buildingLevel >= 20 {
		return 2.5
	}
	// Linear interpolation within the breakpoints.
	breakpoints := []struct {
		level      int
		multiplier float64
	}{
		{1, 1.0},
		{5, 1.25},
		{10, 1.6},
		{15, 2.0},
		{20, 2.5},
	}
	for i := 1; i < len(breakpoints); i++ {
		if buildingLevel <= breakpoints[i].level {
			lo := breakpoints[i-1]
			hi := breakpoints[i]
			frac := float64(buildingLevel-lo.level) / float64(hi.level-lo.level)
			return lo.multiplier + frac*(hi.multiplier-lo.multiplier)
		}
	}
	return 1.0
}

// TrainingTime calculates the training time per unit in seconds, accounting for building level.
// Formula: base_time / speed_multiplier (rounded up).
func TrainingTime(troopType string, buildingLevel int) (int, error) {
	cfg, ok := TroopConfigs[troopType]
	if !ok {
		return 0, ErrUnknownTroop
	}
	mult := TrainingSpeedMultiplier(buildingLevel)
	if mult <= 0 {
		return 0, ErrBuildingNotBuilt
	}
	seconds := float64(cfg.BaseTimeSec) / mult
	return int(math.Ceil(seconds)), nil
}

// TrainingCost returns the total resource cost for training a given quantity of troops.
func TrainingCost(troopType string, quantity int) (ResourceCost, error) {
	cfg, ok := TroopConfigs[troopType]
	if !ok {
		return ResourceCost{}, ErrUnknownTroop
	}
	qty := float64(quantity)
	return ResourceCost{
		Food:   cfg.BaseCost.Food * qty,
		Water:  cfg.BaseCost.Water * qty,
		Lumber: cfg.BaseCost.Lumber * qty,
		Stone:  cfg.BaseCost.Stone * qty,
	}, nil
}

// Troop config errors.
var (
	ErrUnknownTroop     = errStr("unknown troop type")
	ErrBuildingNotBuilt = errStr("training building not built")
)

// errStr is a simple error type for config-level errors.
type errStr string

func (e errStr) Error() string { return string(e) }

package config

import (
	"fmt"
	"math"
	"strings"

	"github.com/luisfpires18/woo/internal/model"
)

// ResourceCost represents the resource cost for an action.
// The four base resources are: Food, Water, Lumber, Stone.
type ResourceCost struct {
	Food   float64
	Water  float64
	Lumber float64
	Stone  float64
}

// BuildingPrerequisite describes a required building type and minimum level.
type BuildingPrerequisite struct {
	BuildingType string
	MinLevel     int
}

// BuildingConfig holds base stats and scaling data for a building type.
type BuildingConfig struct {
	DisplayName     string
	BaseCost        ResourceCost
	BaseTimeSec     int                    // build time at level 1 in seconds
	ScalingFactor   float64                // cost multiplier per level
	TimeFactor      float64                // time multiplier per level
	MaxLevel        int                    // maximum building level
	Prerequisites   []BuildingPrerequisite // required buildings
	PopCapPerLevel  int                    // population capacity added per level of this building
}

// BasePopulation is the base population capacity every village starts with.
const BasePopulation = 20

// resourceBuildingCost is the shared base cost for all 12 resource field buildings.
var resourceBuildingCost = ResourceCost{Food: 0, Water: 0, Lumber: 90, Stone: 50}

// newResourceBuilding returns a BuildingConfig for a resource-producing building slot.
func newResourceBuilding(displayName string) BuildingConfig {
	return BuildingConfig{
		DisplayName:    displayName,
		BaseCost:       resourceBuildingCost,
		BaseTimeSec:    120, // 2 min
		ScalingFactor:  1.5,
		TimeFactor:     1.5,
		MaxLevel:       20,
		PopCapPerLevel: 3,
	}
}

// BuildingConfigs is the authoritative registry of all building types and their stats.
var BuildingConfigs = map[string]BuildingConfig{
	// --- Village buildings ---
	"town_hall": {
		DisplayName:    "Town Hall",
		BaseCost:       ResourceCost{Food: 0, Water: 0, Lumber: 250, Stone: 250},
		BaseTimeSec:    300, // 5 min
		ScalingFactor:  1.7,
		TimeFactor:     1.7,
		MaxLevel:       20,
		PopCapPerLevel: 15,
	},
	"barracks": {
		DisplayName:    "Barracks",
		BaseCost:       ResourceCost{Food: 50, Water: 0, Lumber: 200, Stone: 150},
		BaseTimeSec:    300,
		ScalingFactor:  1.8,
		TimeFactor:     1.8,
		MaxLevel:       20,
		PopCapPerLevel: 5,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 3},
		},
	},
	"stable": {
		DisplayName:    "Stable",
		BaseCost:       ResourceCost{Food: 40, Water: 60, Lumber: 250, Stone: 200},
		BaseTimeSec:    480, // 8 min
		ScalingFactor:  1.8,
		TimeFactor:     1.8,
		MaxLevel:       15,
		PopCapPerLevel: 5,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 5},
			{BuildingType: "barracks", MinLevel: 5},
		},
	},
	"archery": {
		DisplayName:    "Archery",
		BaseCost:       ResourceCost{Food: 0, Water: 0, Lumber: 200, Stone: 150},
		BaseTimeSec:    300, // 5 min
		ScalingFactor:  1.8,
		TimeFactor:     1.8,
		MaxLevel:       15,
		PopCapPerLevel: 5,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 3},
		},
	},
	"workshop": {
		DisplayName:    "Workshop",
		BaseCost:       ResourceCost{Food: 0, Water: 0, Lumber: 350, Stone: 300},
		BaseTimeSec:    600, // 10 min
		ScalingFactor:  1.8,
		TimeFactor:     1.8,
		MaxLevel:       15,
		PopCapPerLevel: 5,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 7},
			{BuildingType: "barracks", MinLevel: 5},
		},
	},
	"special": {
		DisplayName:    "Special",
		BaseCost:       ResourceCost{Food: 40, Water: 40, Lumber: 300, Stone: 350},
		BaseTimeSec:    900, // 15 min
		ScalingFactor:  1.8,
		TimeFactor:     1.8,
		MaxLevel:       15,
		PopCapPerLevel: 5,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 10},
			{BuildingType: "barracks", MinLevel: 7},
			{BuildingType: "stable", MinLevel: 5},
		},
	},

	// --- Storage buildings (increase resource capacity) ---
	"storage": {
		DisplayName:    "Storage",
		BaseCost:       ResourceCost{Food: 0, Water: 0, Lumber: 150, Stone: 120},
		BaseTimeSec:    180, // 3 min
		ScalingFactor:  1.6,
		TimeFactor:     1.6,
		MaxLevel:       20,
		PopCapPerLevel: 5,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 1},
		},
	},
	"provisions": {
		DisplayName:    "Provisions",
		BaseCost:       ResourceCost{Food: 40, Water: 0, Lumber: 120, Stone: 100},
		BaseTimeSec:    180, // 3 min
		ScalingFactor:  1.6,
		TimeFactor:     1.6,
		MaxLevel:       20,
		PopCapPerLevel: 5,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 1},
		},
	},
	"reservoir": {
		DisplayName:    "Reservoir",
		BaseCost:       ResourceCost{Food: 0, Water: 0, Lumber: 100, Stone: 150},
		BaseTimeSec:    180, // 3 min
		ScalingFactor:  1.6,
		TimeFactor:     1.6,
		MaxLevel:       20,
		PopCapPerLevel: 5,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 1},
		},
	},

	// --- Resource field buildings (3 per resource type) ---
	"food_1":   newResourceBuilding("Food Field I"),
	"food_2":   newResourceBuilding("Food Field II"),
	"food_3":   newResourceBuilding("Food Field III"),
	"water_1":  newResourceBuilding("Water Field I"),
	"water_2":  newResourceBuilding("Water Field II"),
	"water_3":  newResourceBuilding("Water Field III"),
	"lumber_1": newResourceBuilding("Lumber Field I"),
	"lumber_2": newResourceBuilding("Lumber Field II"),
	"lumber_3": newResourceBuilding("Lumber Field III"),
	"stone_1":  newResourceBuilding("Stone Field I"),
	"stone_2":  newResourceBuilding("Stone Field II"),
	"stone_3":  newResourceBuilding("Stone Field III"),
}

// ResourceTypeForBuilding returns the resource type a building produces, or "" if not a resource building.
// Mapping: food_1/2/3 → "food", water_1/2/3 → "water", lumber_1/2/3 → "lumber", stone_1/2/3 → "stone".
func ResourceTypeForBuilding(buildingType string) string {
	for _, res := range []string{"food", "water", "lumber", "stone"} {
		if strings.HasPrefix(buildingType, res+"_") {
			return res
		}
	}
	return ""
}

// IsResourceBuilding returns true if the building type produces resources.
func IsResourceBuilding(buildingType string) bool {
	return ResourceTypeForBuilding(buildingType) != ""
}

// ResourceBuildingTypes returns all 12 resource building type IDs.
func ResourceBuildingTypes() []string {
	return []string{
		"food_1", "food_2", "food_3",
		"water_1", "water_2", "water_3",
		"lumber_1", "lumber_2", "lumber_3",
		"stone_1", "stone_2", "stone_3",
	}
}

// IsStorageBuilding returns true if the building type increases resource storage capacity.
func IsStorageBuilding(buildingType string) bool {
	_, ok := StorageBuildingTypes[buildingType]
	return ok
}

// StorageCaps holds per-resource storage capacity for a village.
type StorageCaps struct {
	MaxFood   float64
	MaxWater  float64
	MaxLumber float64
	MaxStone  float64
}

// CalculateStorageCaps computes per-resource storage caps for a village based on its storage buildings.
// Each storage building type only increases the cap for its corresponding resource(s):
//   - Storage → lumber, stone
//   - Provisions → food
//   - Reservoir → water
func CalculateStorageCaps(buildings []*model.Building) StorageCaps {
	caps := StorageCaps{
		MaxFood:   BaseStorage,
		MaxWater:  BaseStorage,
		MaxLumber: BaseStorage,
		MaxStone:  BaseStorage,
	}
	for _, b := range buildings {
		if resourceTypes, ok := StorageBuildingTypes[b.BuildingType]; ok {
			bonus := float64(b.Level) * StoragePerLevel
			for _, rt := range resourceTypes {
				switch rt {
				case "food":
					caps.MaxFood += bonus
				case "water":
					caps.MaxWater += bonus
				case "lumber":
					caps.MaxLumber += bonus
				case "stone":
					caps.MaxStone += bonus
				}
			}
		}
	}
	return caps
}

// CalculatePopCap computes the total population capacity for a village based on its buildings.
// Pop cap = BasePopulation + sum(building.Level * building.PopCapPerLevel) for all built buildings.
func CalculatePopCap(buildings []*model.Building) int {
	cap := BasePopulation
	for _, b := range buildings {
		if b.Level <= 0 {
			continue
		}
		cfg, ok := BuildingConfigs[b.BuildingType]
		if !ok {
			continue
		}
		cap += b.Level * cfg.PopCapPerLevel
	}
	return cap
}

// CostAtLevel calculates the resource cost for upgrading a building to the given level.
// Formula: base_cost × scaling_factor^(level-1)
func CostAtLevel(buildingType string, level int) (ResourceCost, error) {
	cfg, ok := BuildingConfigs[buildingType]
	if !ok {
		return ResourceCost{}, fmt.Errorf("unknown building type: %s", buildingType)
	}
	if level < 1 || level > cfg.MaxLevel {
		return ResourceCost{}, fmt.Errorf("level %d out of range [1, %d]", level, cfg.MaxLevel)
	}

	mult := math.Pow(cfg.ScalingFactor, float64(level-1))
	return ResourceCost{
		Food:   math.Round(cfg.BaseCost.Food * mult),
		Water:  math.Round(cfg.BaseCost.Water * mult),
		Lumber: math.Round(cfg.BaseCost.Lumber * mult),
		Stone:  math.Round(cfg.BaseCost.Stone * mult),
	}, nil
}

// TimeAtLevel calculates the build time in seconds for upgrading a building to the given level.
// Formula: base_time × time_factor^(level-1)
func TimeAtLevel(buildingType string, level int) (int, error) {
	cfg, ok := BuildingConfigs[buildingType]
	if !ok {
		return 0, fmt.Errorf("unknown building type: %s", buildingType)
	}
	if level < 1 || level > cfg.MaxLevel {
		return 0, fmt.Errorf("level %d out of range [1, %d]", level, cfg.MaxLevel)
	}

	seconds := float64(cfg.BaseTimeSec) * math.Pow(cfg.TimeFactor, float64(level-1))
	return int(math.Round(seconds)), nil
}

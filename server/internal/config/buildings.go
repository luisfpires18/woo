package config

import (
	"fmt"
	"math"
)

// ResourceCost represents the resource cost for an action.
type ResourceCost struct {
	Iron  float64
	Wood  float64
	Stone float64
	Food  float64
}

// BuildingPrerequisite describes a required building type and minimum level.
type BuildingPrerequisite struct {
	BuildingType string
	MinLevel     int
}

// BuildingConfig holds base stats and scaling data for a building type.
type BuildingConfig struct {
	DisplayName   string
	BaseCost      ResourceCost
	BaseTimeSec   int                    // build time at level 1 in seconds
	ScalingFactor float64                // cost multiplier per level
	TimeFactor    float64                // time multiplier per level
	MaxLevel      int                    // maximum building level
	Prerequisites []BuildingPrerequisite // required buildings
	KingdomOnly   string                 // empty = all kingdoms, otherwise kingdom slug
}

// BuildingConfigs is the authoritative registry of all building types and their stats.
// Values are draft from docs/01-game-design/progression-and-scaling.md until
// game-template.md is finalized.
var BuildingConfigs = map[string]BuildingConfig{
	"town_hall": {
		DisplayName:   "Town Hall",
		BaseCost:      ResourceCost{Iron: 200, Wood: 200, Stone: 200, Food: 100},
		BaseTimeSec:   300, // 5 min
		ScalingFactor: 1.7,
		TimeFactor:    1.7,
		MaxLevel:      20,
	},
	"iron_mine": {
		DisplayName:   "Iron Mine",
		BaseCost:      ResourceCost{Iron: 100, Wood: 80, Stone: 50, Food: 30},
		BaseTimeSec:   120, // 2 min
		ScalingFactor: 1.5,
		TimeFactor:    1.5,
		MaxLevel:      20,
	},
	"lumber_mill": {
		DisplayName:   "Lumber Mill",
		BaseCost:      ResourceCost{Iron: 80, Wood: 100, Stone: 50, Food: 30},
		BaseTimeSec:   120,
		ScalingFactor: 1.5,
		TimeFactor:    1.5,
		MaxLevel:      20,
	},
	"quarry": {
		DisplayName:   "Quarry",
		BaseCost:      ResourceCost{Iron: 80, Wood: 50, Stone: 100, Food: 30},
		BaseTimeSec:   120,
		ScalingFactor: 1.5,
		TimeFactor:    1.5,
		MaxLevel:      20,
	},
	"farm": {
		DisplayName:   "Farm",
		BaseCost:      ResourceCost{Iron: 50, Wood: 80, Stone: 50, Food: 20},
		BaseTimeSec:   120,
		ScalingFactor: 1.5,
		TimeFactor:    1.5,
		MaxLevel:      20,
	},
	"warehouse": {
		DisplayName:   "Warehouse",
		BaseCost:      ResourceCost{Iron: 120, Wood: 120, Stone: 100, Food: 50},
		BaseTimeSec:   180, // 3 min
		ScalingFactor: 1.6,
		TimeFactor:    1.6,
		MaxLevel:      20,
	},
	"barracks": {
		DisplayName:   "Barracks",
		BaseCost:      ResourceCost{Iron: 200, Wood: 150, Stone: 100, Food: 80},
		BaseTimeSec:   300,
		ScalingFactor: 1.8,
		TimeFactor:    1.8,
		MaxLevel:      20,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 3},
		},
	},
	"stable": {
		DisplayName:   "Stable",
		BaseCost:      ResourceCost{Iron: 300, Wood: 200, Stone: 150, Food: 120},
		BaseTimeSec:   480, // 8 min
		ScalingFactor: 1.8,
		TimeFactor:    1.8,
		MaxLevel:      15,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 5},
			{BuildingType: "barracks", MinLevel: 5},
		},
	},
	"forge": {
		DisplayName:   "Forge",
		BaseCost:      ResourceCost{Iron: 250, Wood: 180, Stone: 200, Food: 100},
		BaseTimeSec:   480,
		ScalingFactor: 1.8,
		TimeFactor:    1.8,
		MaxLevel:      10,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 5},
			{BuildingType: "barracks", MinLevel: 3},
		},
	},
	"rune_altar": {
		DisplayName:   "Rune Altar",
		BaseCost:      ResourceCost{Iron: 300, Wood: 250, Stone: 250, Food: 150},
		BaseTimeSec:   600, // 10 min
		ScalingFactor: 1.9,
		TimeFactor:    1.9,
		MaxLevel:      10,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 7},
			{BuildingType: "forge", MinLevel: 3},
		},
	},
	"walls": {
		DisplayName:   "Walls",
		BaseCost:      ResourceCost{Iron: 150, Wood: 100, Stone: 200, Food: 50},
		BaseTimeSec:   240, // 4 min
		ScalingFactor: 1.6,
		TimeFactor:    1.6,
		MaxLevel:      20,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 2},
		},
	},
	"marketplace": {
		DisplayName:   "Marketplace",
		BaseCost:      ResourceCost{Iron: 180, Wood: 180, Stone: 120, Food: 80},
		BaseTimeSec:   300,
		ScalingFactor: 1.6,
		TimeFactor:    1.6,
		MaxLevel:      15,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 5},
			{BuildingType: "warehouse", MinLevel: 3},
		},
	},
	"embassy": {
		DisplayName:   "Embassy",
		BaseCost:      ResourceCost{Iron: 200, Wood: 200, Stone: 200, Food: 100},
		BaseTimeSec:   480,
		ScalingFactor: 1.7,
		TimeFactor:    1.7,
		MaxLevel:      10,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 8},
		},
	},
	"watchtower": {
		DisplayName:   "Watchtower",
		BaseCost:      ResourceCost{Iron: 150, Wood: 100, Stone: 150, Food: 60},
		BaseTimeSec:   240,
		ScalingFactor: 1.6,
		TimeFactor:    1.6,
		MaxLevel:      10,
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 3},
			{BuildingType: "walls", MinLevel: 1},
		},
	},
	"dock": {
		DisplayName:   "Dock",
		BaseCost:      ResourceCost{Iron: 250, Wood: 300, Stone: 150, Food: 100},
		BaseTimeSec:   480,
		ScalingFactor: 1.8,
		TimeFactor:    1.8,
		MaxLevel:      15,
		KingdomOnly:   "veridor",
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 6},
		},
	},
	"grove_sanctum": {
		DisplayName:   "Grove Sanctum",
		BaseCost:      ResourceCost{Iron: 200, Wood: 300, Stone: 200, Food: 150},
		BaseTimeSec:   480,
		ScalingFactor: 1.8,
		TimeFactor:    1.8,
		MaxLevel:      15,
		KingdomOnly:   "sylvara",
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 6},
		},
	},
	"colosseum": {
		DisplayName:   "Colosseum",
		BaseCost:      ResourceCost{Iron: 300, Wood: 200, Stone: 300, Food: 100},
		BaseTimeSec:   480,
		ScalingFactor: 1.8,
		TimeFactor:    1.8,
		MaxLevel:      15,
		KingdomOnly:   "arkazia",
		Prerequisites: []BuildingPrerequisite{
			{BuildingType: "town_hall", MinLevel: 6},
		},
	},
}

// ResourceRatePerLevel defines how much each resource building produces per hour per level.
// Rate = base_rate + (level * rate_per_level)
const BaseResourceRate = 10.0  // rate at level 0 (idle)
const RatePerLevel = 20.0     // additional rate per level

// StoragePerLevel defines warehouse capacity scaling.
const BaseStorage = 1000.0    // storage at warehouse level 0
const StoragePerLevel = 500.0 // additional storage per warehouse level

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
		Iron:  math.Round(cfg.BaseCost.Iron * mult),
		Wood:  math.Round(cfg.BaseCost.Wood * mult),
		Stone: math.Round(cfg.BaseCost.Stone * mult),
		Food:  math.Round(cfg.BaseCost.Food * mult),
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

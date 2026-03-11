package config

import "sort"

// ---------------------------------------------------------------------------
// Generated DTO types — shared contract between genconfig and parity_test.
// These structs define the JSON schema for client/src/config/generated/*.json.
// ---------------------------------------------------------------------------

// GeneratedResourceCost is the JSON-serialisable form of ResourceCost.
type GeneratedResourceCost struct {
	Food   float64 `json:"food"`
	Water  float64 `json:"water"`
	Lumber float64 `json:"lumber"`
	Stone  float64 `json:"stone"`
}

// GeneratedPrerequisite is the JSON-serialisable form of BuildingPrerequisite.
type GeneratedPrerequisite struct {
	BuildingType string `json:"buildingType"`
	MinLevel     int    `json:"minLevel"`
}

// GeneratedBuildingConfig is the JSON-serialisable form of BuildingConfig.
type GeneratedBuildingConfig struct {
	DisplayName    string                  `json:"displayName"`
	BaseCost       GeneratedResourceCost   `json:"baseCost"`
	BaseTimeSec    int                     `json:"baseTimeSec"`
	ScalingFactor  float64                 `json:"scalingFactor"`
	TimeFactor     float64                 `json:"timeFactor"`
	MaxLevel       int                     `json:"maxLevel"`
	Prerequisites  []GeneratedPrerequisite `json:"prerequisites"`
	PopCapPerLevel int                     `json:"popCapPerLevel"`
}

// GeneratedTroopConfig is the JSON-serialisable form of TroopConfig.
type GeneratedTroopConfig struct {
	DisplayName      string                `json:"displayName"`
	Kingdom          string                `json:"kingdom"`
	TrainingBuilding string                `json:"trainingBuilding"`
	BuildingLevelReq int                   `json:"buildingLevelReq"`
	BaseCost         GeneratedResourceCost `json:"baseCost"`
	BaseTimeSec      int                   `json:"baseTimeSec"`
	FoodUpkeep       float64               `json:"foodUpkeep"`
	Attack           int                   `json:"attack"`
	DefInfantry      int                   `json:"defInfantry"`
	DefCavalry       int                   `json:"defCavalry"`
	Speed            int                   `json:"speed"`
	Carry            int                   `json:"carry"`
	PopCost          int                   `json:"popCost"`
}

// GeneratedResourceEconomy holds all resource-economy constants.
type GeneratedResourceEconomy struct {
	StartingResources float64             `json:"startingResources"`
	StartingRate      float64             `json:"startingRate"`
	BaseResourceRate  float64             `json:"baseResourceRate"`
	RatePerLevel      float64             `json:"ratePerLevel"`
	BaseStorage       float64             `json:"baseStorage"`
	StoragePerLevel   float64             `json:"storagePerLevel"`
	StorageBuildings  map[string][]string `json:"storageBuildings"`
	BasePopulation    int                 `json:"basePopulation"`
}

// ToGeneratedResourceEconomy reads the current constants and returns the DTO.
func ToGeneratedResourceEconomy() GeneratedResourceEconomy {
	// Deep-copy the map so callers can't mutate the original.
	sb := make(map[string][]string, len(StorageBuildingTypes))
	for k, v := range StorageBuildingTypes {
		cp := make([]string, len(v))
		copy(cp, v)
		sb[k] = cp
	}
	return GeneratedResourceEconomy{
		StartingResources: StartingResources,
		StartingRate:      StartingRate,
		BaseResourceRate:  BaseResourceRate,
		RatePerLevel:      RatePerLevel,
		BaseStorage:       BaseStorage,
		StoragePerLevel:   StoragePerLevel,
		StorageBuildings:  sb,
		BasePopulation:    BasePopulation,
	}
}

// ---------------------------------------------------------------------------
// Conversion helpers
// ---------------------------------------------------------------------------

// ToGeneratedBuilding converts a BuildingConfig to its generated DTO.
func ToGeneratedBuilding(cfg BuildingConfig) GeneratedBuildingConfig {
	prereqs := make([]GeneratedPrerequisite, len(cfg.Prerequisites))
	for i, p := range cfg.Prerequisites {
		prereqs[i] = GeneratedPrerequisite{
			BuildingType: p.BuildingType,
			MinLevel:     p.MinLevel,
		}
	}
	return GeneratedBuildingConfig{
		DisplayName:    cfg.DisplayName,
		BaseCost:       toGeneratedCost(cfg.BaseCost),
		BaseTimeSec:    cfg.BaseTimeSec,
		ScalingFactor:  cfg.ScalingFactor,
		TimeFactor:     cfg.TimeFactor,
		MaxLevel:       cfg.MaxLevel,
		Prerequisites:  prereqs,
		PopCapPerLevel: cfg.PopCapPerLevel,
	}
}

// ToGeneratedTroop converts a TroopConfig to its generated DTO.
func ToGeneratedTroop(cfg TroopConfig) GeneratedTroopConfig {
	return GeneratedTroopConfig{
		DisplayName:      cfg.DisplayName,
		Kingdom:          cfg.Kingdom,
		TrainingBuilding: cfg.TrainingBuilding,
		BuildingLevelReq: cfg.BuildingLevelReq,
		BaseCost:         toGeneratedCost(cfg.BaseCost),
		BaseTimeSec:      cfg.BaseTimeSec,
		FoodUpkeep:       cfg.FoodUpkeep,
		Attack:           cfg.Attack,
		DefInfantry:      cfg.DefInfantry,
		DefCavalry:       cfg.DefCavalry,
		Speed:            cfg.Speed,
		Carry:            cfg.Carry,
		PopCost:          effectiveTroopPopCost(cfg),
	}
}

// SortedBuildingKeys returns BuildingConfigs keys in deterministic sorted order.
func SortedBuildingKeys() []string {
	keys := make([]string, 0, len(BuildingConfigs))
	for k := range BuildingConfigs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// SortedTroopKeys returns TroopConfigs keys in deterministic sorted order.
func SortedTroopKeys() []string {
	keys := make([]string, 0, len(TroopConfigs))
	for k := range TroopConfigs {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func toGeneratedCost(c ResourceCost) GeneratedResourceCost {
	return GeneratedResourceCost{
		Food:   c.Food,
		Water:  c.Water,
		Lumber: c.Lumber,
		Stone:  c.Stone,
	}
}

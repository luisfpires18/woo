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
	PopCost          int // population consumed when this troop is trained (0 = use building-type default)
}

// defaultPopCostByBuilding maps training building types to their default population cost per troop.
var defaultPopCostByBuilding = map[string]int{
	"barracks": 1,
	"stable":   2,
	"archery":  2,
	"workshop": 3,
	"special":  3,
}

// effectiveTroopPopCost returns the pop cost for a TroopConfig struct without requiring a map lookup.
// Used internally by ToGeneratedTroop.
func effectiveTroopPopCost(cfg TroopConfig) int {
	if cfg.PopCost > 0 {
		return cfg.PopCost
	}
	if def, ok := defaultPopCostByBuilding[cfg.TrainingBuilding]; ok {
		return def
	}
	return 1
}

// TroopPopCost returns the population cost for training one unit of the given troop type.
// Uses the configured PopCost if set (> 0), otherwise falls back to the building-type default.
func TroopPopCost(troopType string) int {
	cfg, ok := TroopConfigs[troopType]
	if !ok {
		return 1 // safe default
	}
	if cfg.PopCost > 0 {
		return cfg.PopCost
	}
	if def, ok := defaultPopCostByBuilding[cfg.TrainingBuilding]; ok {
		return def
	}
	return 1
}

// TroopConfigs is the authoritative registry of all troop types and their stats.
// Each kingdom has ~19-21 units across 5 buildings: barracks (4), stable (3-4),
// archery (4), workshop (4), special (4-5).
// Stats are balanced per building-tier across kingdoms for fairness.
var TroopConfigs = map[string]TroopConfig{

	// ═══════════════════════════════════════════════════════════════════════
	// SYLVARA — Forest / Druid / Nature
	// ═══════════════════════════════════════════════════════════════════════

	// Barracks
	"sylvara_rootguard_spearmen": {
		DisplayName: "Rootguard Spearmen", Kingdom: "sylvara", TrainingBuilding: "barracks", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 100, Water: 50, Lumber: 60, Stone: 40}, BaseTimeSec: 120,
		FoodUpkeep: 1, Attack: 45, DefInfantry: 55, DefCavalry: 40, Speed: 5, Carry: 55,
	},
	"sylvara_thornblade_wardens": {
		DisplayName: "Thornblade Wardens", Kingdom: "sylvara", TrainingBuilding: "barracks", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 120, Water: 60, Lumber: 80, Stone: 50}, BaseTimeSec: 150,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 50, DefCavalry: 35, Speed: 5, Carry: 50,
	},
	"sylvara_boarhide_axemen": {
		DisplayName: "Boarhide Axemen", Kingdom: "sylvara", TrainingBuilding: "barracks", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 140, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 180,
		FoodUpkeep: 1, Attack: 65, DefInfantry: 60, DefCavalry: 45, Speed: 4, Carry: 45,
	},
	"sylvara_leafknife_stalkers": {
		DisplayName: "Leafknife Stalkers", Kingdom: "sylvara", TrainingBuilding: "barracks", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 160, Water: 100, Lumber: 120, Stone: 70}, BaseTimeSec: 210,
		FoodUpkeep: 2, Attack: 75, DefInfantry: 45, DefCavalry: 30, Speed: 7, Carry: 40,
	},
	// Stable
	"sylvara_hart_lancers": {
		DisplayName: "Hart Lancers", Kingdom: "sylvara", TrainingBuilding: "stable", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 180, Water: 100, Lumber: 120, Stone: 80}, BaseTimeSec: 270,
		FoodUpkeep: 2, Attack: 75, DefInfantry: 35, DefCavalry: 55, Speed: 10, Carry: 80,
	},
	"sylvara_boar_riders": {
		DisplayName: "Boar Riders", Kingdom: "sylvara", TrainingBuilding: "stable", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 220, Water: 120, Lumber: 150, Stone: 100}, BaseTimeSec: 330,
		FoodUpkeep: 3, Attack: 90, DefInfantry: 45, DefCavalry: 60, Speed: 8, Carry: 70,
	},
	"sylvara_wolf_outriders": {
		DisplayName: "Wolf Outriders", Kingdom: "sylvara", TrainingBuilding: "stable", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 150, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 240,
		FoodUpkeep: 2, Attack: 60, DefInfantry: 25, DefCavalry: 40, Speed: 14, Carry: 50,
	},
	// Archery
	"sylvara_longbow_wardens": {
		DisplayName: "Longbow Wardens", Kingdom: "sylvara", TrainingBuilding: "archery", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 80, Water: 60, Lumber: 100, Stone: 30}, BaseTimeSec: 130,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 35,
	},
	"sylvara_owlshot_hunters": {
		DisplayName: "Owlshot Hunters", Kingdom: "sylvara", TrainingBuilding: "archery", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 100, Water: 70, Lumber: 120, Stone: 40}, BaseTimeSec: 160,
		FoodUpkeep: 1, Attack: 65, DefInfantry: 20, DefCavalry: 15, Speed: 6, Carry: 30,
	},
	"sylvara_thorn_javeliners": {
		DisplayName: "Thorn Javeliners", Kingdom: "sylvara", TrainingBuilding: "archery", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 90, Water: 60, Lumber: 80, Stone: 35}, BaseTimeSec: 140,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 30, DefCavalry: 25, Speed: 6, Carry: 40,
	},
	"sylvara_stonesling_foresters": {
		DisplayName: "Stone-sling Foresters", Kingdom: "sylvara", TrainingBuilding: "archery", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 60, Water: 40, Lumber: 50, Stone: 25}, BaseTimeSec: 100,
		FoodUpkeep: 1, Attack: 40, DefInfantry: 20, DefCavalry: 15, Speed: 5, Carry: 30,
	},
	// Workshop
	"sylvara_vine_ballista": {
		DisplayName: "Vine Ballista", Kingdom: "sylvara", TrainingBuilding: "workshop", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 20, DefCavalry: 15, Speed: 3, Carry: 0,
	},
	"sylvara_log_trebuchet": {
		DisplayName: "Log Trebuchet", Kingdom: "sylvara", TrainingBuilding: "workshop", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 400, Stone: 200}, BaseTimeSec: 500,
		FoodUpkeep: 4, Attack: 100, DefInfantry: 15, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"sylvara_root_ram": {
		DisplayName: "Root Ram", Kingdom: "sylvara", TrainingBuilding: "workshop", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 220, Water: 80, Lumber: 350, Stone: 180}, BaseTimeSec: 450,
		FoodUpkeep: 4, Attack: 70, DefInfantry: 30, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"sylvara_sporepot_thrower": {
		DisplayName: "Spore-Pot Thrower", Kingdom: "sylvara", TrainingBuilding: "workshop", BuildingLevelReq: 12,
		BaseCost: ResourceCost{Food: 280, Water: 150, Lumber: 350, Stone: 200}, BaseTimeSec: 550,
		FoodUpkeep: 5, Attack: 90, DefInfantry: 10, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	// Special
	"sylvara_beastmasters": {
		DisplayName: "Beastmasters", Kingdom: "sylvara", TrainingBuilding: "special", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 180, Stone: 120}, BaseTimeSec: 500,
		FoodUpkeep: 3, Attack: 85, DefInfantry: 55, DefCavalry: 50, Speed: 6, Carry: 40,
	},
	"sylvara_shapeshifter_scouts": {
		DisplayName: "Shapeshifter Scouts", Kingdom: "sylvara", TrainingBuilding: "special", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 150, Stone: 100}, BaseTimeSec: 450,
		FoodUpkeep: 2, Attack: 70, DefInfantry: 40, DefCavalry: 35, Speed: 12, Carry: 20,
	},
	"sylvara_life_healers": {
		DisplayName: "Life Healers", Kingdom: "sylvara", TrainingBuilding: "special", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 300, Water: 200, Lumber: 200, Stone: 150}, BaseTimeSec: 550,
		FoodUpkeep: 3, Attack: 30, DefInfantry: 70, DefCavalry: 65, Speed: 5, Carry: 20,
	},
	"sylvara_antler_seers": {
		DisplayName: "Antler Seers", Kingdom: "sylvara", TrainingBuilding: "special", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 280, Water: 180, Lumber: 180, Stone: 130}, BaseTimeSec: 520,
		FoodUpkeep: 3, Attack: 40, DefInfantry: 60, DefCavalry: 55, Speed: 5, Carry: 15,
	},

	// ═══════════════════════════════════════════════════════════════════════
	// ARKAZIA — Mountain / Fortress / Iron
	// ═══════════════════════════════════════════════════════════════════════

	// Barracks
	"arkazia_shield_guard": {
		DisplayName: "Shield Guard", Kingdom: "arkazia", TrainingBuilding: "barracks", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 100, Water: 50, Lumber: 60, Stone: 40}, BaseTimeSec: 120,
		FoodUpkeep: 1, Attack: 45, DefInfantry: 55, DefCavalry: 40, Speed: 5, Carry: 55,
	},
	"arkazia_pike_guard": {
		DisplayName: "Pike Guard", Kingdom: "arkazia", TrainingBuilding: "barracks", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 120, Water: 60, Lumber: 80, Stone: 50}, BaseTimeSec: 150,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 60, DefCavalry: 50, Speed: 4, Carry: 45,
	},
	"arkazia_stonebreaker": {
		DisplayName: "Stonebreaker", Kingdom: "arkazia", TrainingBuilding: "barracks", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 140, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 180,
		FoodUpkeep: 1, Attack: 70, DefInfantry: 55, DefCavalry: 40, Speed: 4, Carry: 45,
	},
	"arkazia_rampart_halberdiers": {
		DisplayName: "Rampart Halberdiers", Kingdom: "arkazia", TrainingBuilding: "barracks", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 160, Water: 100, Lumber: 120, Stone: 70}, BaseTimeSec: 210,
		FoodUpkeep: 2, Attack: 65, DefInfantry: 65, DefCavalry: 55, Speed: 4, Carry: 40,
	},
	// Stable
	"arkazia_arknight_lancers": {
		DisplayName: "Arknight Lancers", Kingdom: "arkazia", TrainingBuilding: "stable", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 180, Water: 100, Lumber: 120, Stone: 80}, BaseTimeSec: 270,
		FoodUpkeep: 2, Attack: 80, DefInfantry: 40, DefCavalry: 55, Speed: 10, Carry: 80,
	},
	"arkazia_redcrest_cavaliers": {
		DisplayName: "Redcrest Cavaliers", Kingdom: "arkazia", TrainingBuilding: "stable", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 220, Water: 120, Lumber: 150, Stone: 100}, BaseTimeSec: 330,
		FoodUpkeep: 3, Attack: 85, DefInfantry: 50, DefCavalry: 65, Speed: 8, Carry: 70,
	},
	"arkazia_banner_riders": {
		DisplayName: "Banner Riders", Kingdom: "arkazia", TrainingBuilding: "stable", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 200, Water: 110, Lumber: 140, Stone: 90}, BaseTimeSec: 300,
		FoodUpkeep: 2, Attack: 60, DefInfantry: 45, DefCavalry: 60, Speed: 9, Carry: 60,
	},
	"arkazia_ravine_pursuers": {
		DisplayName: "Ravine Pursuers", Kingdom: "arkazia", TrainingBuilding: "stable", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 150, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 240,
		FoodUpkeep: 2, Attack: 55, DefInfantry: 25, DefCavalry: 40, Speed: 14, Carry: 50,
	},
	// Archery
	"arkazia_hill_slingers": {
		DisplayName: "Hill Slingers", Kingdom: "arkazia", TrainingBuilding: "archery", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 60, Water: 40, Lumber: 50, Stone: 25}, BaseTimeSec: 100,
		FoodUpkeep: 1, Attack: 40, DefInfantry: 20, DefCavalry: 15, Speed: 5, Carry: 30,
	},
	"arkazia_ridge_crossbowmen": {
		DisplayName: "Ridge Crossbowmen", Kingdom: "arkazia", TrainingBuilding: "archery", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 80, Water: 60, Lumber: 100, Stone: 30}, BaseTimeSec: 130,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 35,
	},
	"arkazia_javelin_climbers": {
		DisplayName: "Javelin Climbers", Kingdom: "arkazia", TrainingBuilding: "archery", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 90, Water: 60, Lumber: 80, Stone: 35}, BaseTimeSec: 140,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 30, DefCavalry: 25, Speed: 6, Carry: 40,
	},
	"arkazia_crag_bowmen": {
		DisplayName: "Crag Bowmen", Kingdom: "arkazia", TrainingBuilding: "archery", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 100, Water: 70, Lumber: 120, Stone: 40}, BaseTimeSec: 160,
		FoodUpkeep: 1, Attack: 65, DefInfantry: 20, DefCavalry: 15, Speed: 5, Carry: 30,
	},
	// Workshop
	"arkazia_mountain_trebuchet": {
		DisplayName: "Mountain Trebuchet", Kingdom: "arkazia", TrainingBuilding: "workshop", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 400, Stone: 200}, BaseTimeSec: 500,
		FoodUpkeep: 4, Attack: 100, DefInfantry: 15, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"arkazia_ram_engineers": {
		DisplayName: "Ram Engineers", Kingdom: "arkazia", TrainingBuilding: "workshop", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 220, Water: 80, Lumber: 350, Stone: 180}, BaseTimeSec: 450,
		FoodUpkeep: 4, Attack: 70, DefInfantry: 30, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"arkazia_mantlet_pushers": {
		DisplayName: "Mantlet Pushers", Kingdom: "arkazia", TrainingBuilding: "workshop", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 30, DefInfantry: 50, DefCavalry: 40, Speed: 3, Carry: 0,
	},
	"arkazia_bridgewright_crew": {
		DisplayName: "Bridgewright Crew", Kingdom: "arkazia", TrainingBuilding: "workshop", BuildingLevelReq: 12,
		BaseCost: ResourceCost{Food: 280, Water: 150, Lumber: 350, Stone: 200}, BaseTimeSec: 550,
		FoodUpkeep: 5, Attack: 40, DefInfantry: 40, DefCavalry: 30, Speed: 3, Carry: 0,
	},
	// Special
	"arkazia_banner_knights": {
		DisplayName: "Banner Knights", Kingdom: "arkazia", TrainingBuilding: "special", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 180, Stone: 120}, BaseTimeSec: 500,
		FoodUpkeep: 3, Attack: 85, DefInfantry: 55, DefCavalry: 50, Speed: 8, Carry: 40,
	},
	"arkazia_oathsworn_champions": {
		DisplayName: "Oathsworn Champions", Kingdom: "arkazia", TrainingBuilding: "special", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 300, Water: 150, Lumber: 200, Stone: 150}, BaseTimeSec: 600,
		FoodUpkeep: 3, Attack: 95, DefInfantry: 60, DefCavalry: 50, Speed: 5, Carry: 30,
	},
	"arkazia_bastion_marshals": {
		DisplayName: "Bastion Marshals", Kingdom: "arkazia", TrainingBuilding: "special", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 350, Water: 200, Lumber: 250, Stone: 200}, BaseTimeSec: 650,
		FoodUpkeep: 4, Attack: 60, DefInfantry: 80, DefCavalry: 70, Speed: 4, Carry: 20,
	},
	"arkazia_arknight_captains": {
		DisplayName: "Arknight Captains", Kingdom: "arkazia", TrainingBuilding: "special", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 400, Water: 200, Lumber: 300, Stone: 250}, BaseTimeSec: 700,
		FoodUpkeep: 4, Attack: 90, DefInfantry: 65, DefCavalry: 60, Speed: 9, Carry: 50,
	},

	// ═══════════════════════════════════════════════════════════════════════
	// VERIDOR — Trade / Naval / Roads
	// ═══════════════════════════════════════════════════════════════════════

	// Barracks
	"veridor_road_legionaries": {
		DisplayName: "Road Legionaries", Kingdom: "veridor", TrainingBuilding: "barracks", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 100, Water: 50, Lumber: 60, Stone: 40}, BaseTimeSec: 120,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 50, DefCavalry: 40, Speed: 5, Carry: 55,
	},
	"veridor_harbor_pike": {
		DisplayName: "Harbor Pike", Kingdom: "veridor", TrainingBuilding: "barracks", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 120, Water: 60, Lumber: 80, Stone: 50}, BaseTimeSec: 150,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 60, DefCavalry: 50, Speed: 4, Carry: 45,
	},
	"veridor_cutlass_marines": {
		DisplayName: "Cutlass Marines", Kingdom: "veridor", TrainingBuilding: "barracks", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 140, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 180,
		FoodUpkeep: 1, Attack: 70, DefInfantry: 50, DefCavalry: 35, Speed: 5, Carry: 45,
	},
	"veridor_wharf_axemen": {
		DisplayName: "Wharf Axemen", Kingdom: "veridor", TrainingBuilding: "barracks", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 160, Water: 100, Lumber: 120, Stone: 70}, BaseTimeSec: 210,
		FoodUpkeep: 2, Attack: 75, DefInfantry: 55, DefCavalry: 40, Speed: 4, Carry: 40,
	},
	// Stable
	"veridor_river_lancers": {
		DisplayName: "River Lancers", Kingdom: "veridor", TrainingBuilding: "stable", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 180, Water: 100, Lumber: 120, Stone: 80}, BaseTimeSec: 270,
		FoodUpkeep: 2, Attack: 75, DefInfantry: 35, DefCavalry: 55, Speed: 10, Carry: 80,
	},
	"veridor_courier_riders": {
		DisplayName: "Courier Riders", Kingdom: "veridor", TrainingBuilding: "stable", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 130, Water: 70, Lumber: 80, Stone: 50}, BaseTimeSec: 200,
		FoodUpkeep: 1, Attack: 40, DefInfantry: 20, DefCavalry: 30, Speed: 16, Carry: 40,
	},
	"veridor_marsh_scouts": {
		DisplayName: "Marsh Scouts", Kingdom: "veridor", TrainingBuilding: "stable", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 150, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 240,
		FoodUpkeep: 2, Attack: 55, DefInfantry: 25, DefCavalry: 40, Speed: 12, Carry: 50,
	},
	"veridor_road_wardens": {
		DisplayName: "Road Wardens", Kingdom: "veridor", TrainingBuilding: "stable", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 220, Water: 120, Lumber: 150, Stone: 100}, BaseTimeSec: 330,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 50, DefCavalry: 65, Speed: 8, Carry: 70,
	},
	// Archery
	"veridor_deck_arbalesters": {
		DisplayName: "Deck Arbalesters", Kingdom: "veridor", TrainingBuilding: "archery", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 80, Water: 60, Lumber: 100, Stone: 30}, BaseTimeSec: 130,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 35,
	},
	"veridor_highland_longbowmen": {
		DisplayName: "Highland Longbowmen", Kingdom: "veridor", TrainingBuilding: "archery", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 100, Water: 70, Lumber: 120, Stone: 40}, BaseTimeSec: 160,
		FoodUpkeep: 1, Attack: 65, DefInfantry: 20, DefCavalry: 15, Speed: 5, Carry: 30,
	},
	"veridor_harpoon_casters": {
		DisplayName: "Harpoon Casters", Kingdom: "veridor", TrainingBuilding: "archery", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 90, Water: 60, Lumber: 80, Stone: 35}, BaseTimeSec: 140,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 30, DefCavalry: 25, Speed: 5, Carry: 35,
	},
	"veridor_pavise_marksmen": {
		DisplayName: "Pavise Marksmen", Kingdom: "veridor", TrainingBuilding: "archery", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 120, Water: 80, Lumber: 140, Stone: 50}, BaseTimeSec: 180,
		FoodUpkeep: 1, Attack: 70, DefInfantry: 35, DefCavalry: 25, Speed: 4, Carry: 25,
	},
	// Workshop
	"veridor_harbor_ballista": {
		DisplayName: "Harbor Ballista", Kingdom: "veridor", TrainingBuilding: "workshop", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 20, DefCavalry: 15, Speed: 3, Carry: 0,
	},
	"veridor_mangonel": {
		DisplayName: "Mangonel", Kingdom: "veridor", TrainingBuilding: "workshop", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 400, Stone: 200}, BaseTimeSec: 500,
		FoodUpkeep: 4, Attack: 100, DefInfantry: 15, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"veridor_firepot_crane": {
		DisplayName: "Firepot Crane", Kingdom: "veridor", TrainingBuilding: "workshop", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 280, Water: 150, Lumber: 350, Stone: 200}, BaseTimeSec: 550,
		FoodUpkeep: 5, Attack: 90, DefInfantry: 10, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"veridor_pavise_wagon": {
		DisplayName: "Pavise Wagon", Kingdom: "veridor", TrainingBuilding: "workshop", BuildingLevelReq: 12,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 30, DefInfantry: 50, DefCavalry: 40, Speed: 3, Carry: 0,
	},
	// Special
	"veridor_hydra_hunters": {
		DisplayName: "Hydra Hunters", Kingdom: "veridor", TrainingBuilding: "special", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 180, Stone: 120}, BaseTimeSec: 500,
		FoodUpkeep: 3, Attack: 85, DefInfantry: 55, DefCavalry: 50, Speed: 5, Carry: 40,
	},
	"veridor_tidemark_duelists": {
		DisplayName: "Tidemark Duelists", Kingdom: "veridor", TrainingBuilding: "special", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 280, Water: 140, Lumber: 200, Stone: 140}, BaseTimeSec: 550,
		FoodUpkeep: 3, Attack: 90, DefInfantry: 50, DefCavalry: 45, Speed: 7, Carry: 30,
	},
	"veridor_bluecoat_captains": {
		DisplayName: "Bluecoat Captains", Kingdom: "veridor", TrainingBuilding: "special", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 350, Water: 200, Lumber: 250, Stone: 200}, BaseTimeSec: 650,
		FoodUpkeep: 4, Attack: 70, DefInfantry: 70, DefCavalry: 65, Speed: 5, Carry: 20,
	},
	"veridor_skiff_raiders": {
		DisplayName: "Skiff Raiders", Kingdom: "veridor", TrainingBuilding: "special", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 300, Water: 150, Lumber: 200, Stone: 150}, BaseTimeSec: 600,
		FoodUpkeep: 3, Attack: 95, DefInfantry: 45, DefCavalry: 40, Speed: 8, Carry: 50,
	},

	// ═══════════════════════════════════════════════════════════════════════
	// DRAXYS — Desert / Arena / Scorpion
	// ═══════════════════════════════════════════════════════════════════════

	// Barracks
	"draxys_sandshield_infantry": {
		DisplayName: "Sandshield Infantry", Kingdom: "draxys", TrainingBuilding: "barracks", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 100, Water: 50, Lumber: 60, Stone: 40}, BaseTimeSec: 120,
		FoodUpkeep: 1, Attack: 45, DefInfantry: 55, DefCavalry: 40, Speed: 5, Carry: 55,
	},
	"draxys_khopesh_guard": {
		DisplayName: "Khopesh Guard", Kingdom: "draxys", TrainingBuilding: "barracks", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 120, Water: 60, Lumber: 80, Stone: 50}, BaseTimeSec: 150,
		FoodUpkeep: 1, Attack: 60, DefInfantry: 45, DefCavalry: 35, Speed: 5, Carry: 50,
	},
	"draxys_dune_axemen": {
		DisplayName: "Dune Axemen", Kingdom: "draxys", TrainingBuilding: "barracks", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 140, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 180,
		FoodUpkeep: 1, Attack: 70, DefInfantry: 50, DefCavalry: 35, Speed: 5, Carry: 45,
	},
	"draxys_wadi_lashers": {
		DisplayName: "Wadi Lashers", Kingdom: "draxys", TrainingBuilding: "barracks", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 160, Water: 100, Lumber: 120, Stone: 70}, BaseTimeSec: 210,
		FoodUpkeep: 2, Attack: 75, DefInfantry: 40, DefCavalry: 30, Speed: 6, Carry: 35,
	},
	// Stable
	"draxys_scorpion_riders": {
		DisplayName: "Scorpion Riders", Kingdom: "draxys", TrainingBuilding: "stable", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 200, Water: 120, Lumber: 130, Stone: 90}, BaseTimeSec: 300,
		FoodUpkeep: 3, Attack: 85, DefInfantry: 45, DefCavalry: 55, Speed: 9, Carry: 70,
	},
	"draxys_dune_lancers": {
		DisplayName: "Dune Lancers", Kingdom: "draxys", TrainingBuilding: "stable", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 180, Water: 100, Lumber: 120, Stone: 80}, BaseTimeSec: 270,
		FoodUpkeep: 2, Attack: 75, DefInfantry: 35, DefCavalry: 55, Speed: 10, Carry: 80,
	},
	"draxys_camel_skirmishers": {
		DisplayName: "Camel Skirmishers", Kingdom: "draxys", TrainingBuilding: "stable", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 150, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 240,
		FoodUpkeep: 2, Attack: 55, DefInfantry: 30, DefCavalry: 45, Speed: 12, Carry: 60,
	},
	"draxys_dust_chasers": {
		DisplayName: "Dust Chasers", Kingdom: "draxys", TrainingBuilding: "stable", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 130, Water: 70, Lumber: 80, Stone: 50}, BaseTimeSec: 200,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 20, DefCavalry: 35, Speed: 16, Carry: 40,
	},
	// Archery
	"draxys_oasis_rangers": {
		DisplayName: "Oasis Rangers", Kingdom: "draxys", TrainingBuilding: "archery", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 80, Water: 60, Lumber: 100, Stone: 30}, BaseTimeSec: 130,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 35,
	},
	"draxys_sun_slingers": {
		DisplayName: "Sun Slingers", Kingdom: "draxys", TrainingBuilding: "archery", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 60, Water: 40, Lumber: 50, Stone: 25}, BaseTimeSec: 100,
		FoodUpkeep: 1, Attack: 40, DefInfantry: 20, DefCavalry: 15, Speed: 6, Carry: 30,
	},
	"draxys_chakram_dancers": {
		DisplayName: "Chakram Dancers", Kingdom: "draxys", TrainingBuilding: "archery", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 100, Water: 70, Lumber: 120, Stone: 40}, BaseTimeSec: 160,
		FoodUpkeep: 1, Attack: 60, DefInfantry: 20, DefCavalry: 15, Speed: 6, Carry: 30,
	},
	"draxys_javelin_skirmishers": {
		DisplayName: "Javelin Skirmishers", Kingdom: "draxys", TrainingBuilding: "archery", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 90, Water: 60, Lumber: 80, Stone: 35}, BaseTimeSec: 140,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 25, DefCavalry: 20, Speed: 7, Carry: 35,
	},
	// Workshop
	"draxys_bolt_thrower": {
		DisplayName: "Bolt Thrower", Kingdom: "draxys", TrainingBuilding: "workshop", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 20, DefCavalry: 15, Speed: 3, Carry: 0,
	},
	"draxys_firepot_mangonel": {
		DisplayName: "Firepot Mangonel", Kingdom: "draxys", TrainingBuilding: "workshop", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 280, Water: 150, Lumber: 350, Stone: 200}, BaseTimeSec: 550,
		FoodUpkeep: 5, Attack: 95, DefInfantry: 10, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"draxys_siege_tower": {
		DisplayName: "Siege Tower", Kingdom: "draxys", TrainingBuilding: "workshop", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 300, Water: 150, Lumber: 450, Stone: 250}, BaseTimeSec: 600,
		FoodUpkeep: 5, Attack: 50, DefInfantry: 60, DefCavalry: 50, Speed: 2, Carry: 0,
	},
	"draxys_scorpion_cage_wagon": {
		DisplayName: "Scorpion Cage Wagon", Kingdom: "draxys", TrainingBuilding: "workshop", BuildingLevelReq: 12,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 300, Stone: 180}, BaseTimeSec: 500,
		FoodUpkeep: 4, Attack: 85, DefInfantry: 15, DefCavalry: 10, Speed: 3, Carry: 0,
	},
	// Special
	"draxys_gladiators": {
		DisplayName: "Gladiators", Kingdom: "draxys", TrainingBuilding: "special", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 180, Stone: 120}, BaseTimeSec: 500,
		FoodUpkeep: 3, Attack: 95, DefInfantry: 50, DefCavalry: 45, Speed: 6, Carry: 30,
	},
	"draxys_netfighters": {
		DisplayName: "Netfighters", Kingdom: "draxys", TrainingBuilding: "special", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 220, Water: 110, Lumber: 160, Stone: 110}, BaseTimeSec: 470,
		FoodUpkeep: 2, Attack: 70, DefInfantry: 45, DefCavalry: 40, Speed: 6, Carry: 25,
	},
	"draxys_arena_spearmen": {
		DisplayName: "Arena Spearmen", Kingdom: "draxys", TrainingBuilding: "special", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 280, Water: 140, Lumber: 200, Stone: 140}, BaseTimeSec: 550,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 55, DefCavalry: 50, Speed: 5, Carry: 25,
	},
	"draxys_pit_brutes": {
		DisplayName: "Pit Brutes", Kingdom: "draxys", TrainingBuilding: "special", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 350, Water: 180, Lumber: 250, Stone: 180}, BaseTimeSec: 650,
		FoodUpkeep: 4, Attack: 100, DefInfantry: 60, DefCavalry: 40, Speed: 4, Carry: 20,
	},
	"draxys_beast_tamers": {
		DisplayName: "Beast Tamers", Kingdom: "draxys", TrainingBuilding: "special", BuildingLevelReq: 10,
		BaseCost: ResourceCost{Food: 300, Water: 160, Lumber: 200, Stone: 150}, BaseTimeSec: 580,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 50, DefCavalry: 45, Speed: 7, Carry: 30,
	},

	// ═══════════════════════════════════════════════════════════════════════
	// NORDALH — Northern / Forge / Viking
	// ═══════════════════════════════════════════════════════════════════════

	// Barracks
	"nordalh_hearth_guards": {
		DisplayName: "Hearth Guards", Kingdom: "nordalh", TrainingBuilding: "barracks", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 100, Water: 50, Lumber: 60, Stone: 40}, BaseTimeSec: 120,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 55, DefCavalry: 40, Speed: 5, Carry: 55,
	},
	"nordalh_fjord_spearmen": {
		DisplayName: "Fjord Spearmen", Kingdom: "nordalh", TrainingBuilding: "barracks", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 120, Water: 60, Lumber: 80, Stone: 50}, BaseTimeSec: 150,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 60, DefCavalry: 45, Speed: 5, Carry: 50,
	},
	"nordalh_iceshore_raiders": {
		DisplayName: "Ice-Shore Raiders", Kingdom: "nordalh", TrainingBuilding: "barracks", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 140, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 180,
		FoodUpkeep: 1, Attack: 70, DefInfantry: 45, DefCavalry: 35, Speed: 6, Carry: 50,
	},
	"nordalh_chain_wardens": {
		DisplayName: "Chain Wardens", Kingdom: "nordalh", TrainingBuilding: "barracks", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 160, Water: 100, Lumber: 120, Stone: 70}, BaseTimeSec: 210,
		FoodUpkeep: 2, Attack: 65, DefInfantry: 60, DefCavalry: 50, Speed: 4, Carry: 40,
	},
	// Stable
	"nordalh_direwolf_riders": {
		DisplayName: "Direwolf Riders", Kingdom: "nordalh", TrainingBuilding: "stable", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 180, Water: 100, Lumber: 120, Stone: 80}, BaseTimeSec: 270,
		FoodUpkeep: 2, Attack: 80, DefInfantry: 35, DefCavalry: 50, Speed: 12, Carry: 70,
	},
	"nordalh_elk_lancers": {
		DisplayName: "Elk Lancers", Kingdom: "nordalh", TrainingBuilding: "stable", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 220, Water: 120, Lumber: 150, Stone: 100}, BaseTimeSec: 330,
		FoodUpkeep: 3, Attack: 85, DefInfantry: 45, DefCavalry: 60, Speed: 9, Carry: 75,
	},
	"nordalh_snow_riders": {
		DisplayName: "Snow Riders", Kingdom: "nordalh", TrainingBuilding: "stable", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 150, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 240,
		FoodUpkeep: 2, Attack: 55, DefInfantry: 25, DefCavalry: 40, Speed: 14, Carry: 50,
	},
	"nordalh_fang_cavaliers": {
		DisplayName: "Fang Cavaliers", Kingdom: "nordalh", TrainingBuilding: "stable", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 240, Water: 130, Lumber: 160, Stone: 110}, BaseTimeSec: 350,
		FoodUpkeep: 3, Attack: 90, DefInfantry: 40, DefCavalry: 55, Speed: 8, Carry: 65,
	},
	// Archery
	"nordalh_frostbow_hunters": {
		DisplayName: "Frost Bow Hunters", Kingdom: "nordalh", TrainingBuilding: "archery", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 80, Water: 60, Lumber: 100, Stone: 30}, BaseTimeSec: 130,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 35,
	},
	"nordalh_harpoon_throwers": {
		DisplayName: "Harpoon Throwers", Kingdom: "nordalh", TrainingBuilding: "archery", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 100, Water: 70, Lumber: 120, Stone: 40}, BaseTimeSec: 160,
		FoodUpkeep: 1, Attack: 60, DefInfantry: 30, DefCavalry: 25, Speed: 5, Carry: 35,
	},
	"nordalh_cliff_crossbowmen": {
		DisplayName: "Cliff Crossbowmen", Kingdom: "nordalh", TrainingBuilding: "archery", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 90, Water: 60, Lumber: 110, Stone: 35}, BaseTimeSec: 140,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 30,
	},
	"nordalh_storm_slingers": {
		DisplayName: "Storm Slingers", Kingdom: "nordalh", TrainingBuilding: "archery", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 60, Water: 40, Lumber: 50, Stone: 25}, BaseTimeSec: 100,
		FoodUpkeep: 1, Attack: 40, DefInfantry: 20, DefCavalry: 15, Speed: 5, Carry: 30,
	},
	// Workshop
	"nordalh_cliff_ballista": {
		DisplayName: "Cliff Ballista", Kingdom: "nordalh", TrainingBuilding: "workshop", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 20, DefCavalry: 15, Speed: 3, Carry: 0,
	},
	"nordalh_stone_trebuchet": {
		DisplayName: "Stone Trebuchet", Kingdom: "nordalh", TrainingBuilding: "workshop", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 400, Stone: 200}, BaseTimeSec: 500,
		FoodUpkeep: 4, Attack: 100, DefInfantry: 15, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"nordalh_ram_sled": {
		DisplayName: "Ram Sled", Kingdom: "nordalh", TrainingBuilding: "workshop", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 220, Water: 80, Lumber: 350, Stone: 180}, BaseTimeSec: 450,
		FoodUpkeep: 4, Attack: 70, DefInfantry: 30, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"nordalh_boiling_pitch_crew": {
		DisplayName: "Boiling Pitch Crew", Kingdom: "nordalh", TrainingBuilding: "workshop", BuildingLevelReq: 12,
		BaseCost: ResourceCost{Food: 280, Water: 150, Lumber: 300, Stone: 200}, BaseTimeSec: 520,
		FoodUpkeep: 4, Attack: 85, DefInfantry: 15, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	// Special
	"nordalh_smith_retinues": {
		DisplayName: "Smith Retinues", Kingdom: "nordalh", TrainingBuilding: "special", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 180, Stone: 120}, BaseTimeSec: 500,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 65, DefCavalry: 55, Speed: 4, Carry: 30,
	},
	"nordalh_coyote_blademasters": {
		DisplayName: "Coyote Blademasters", Kingdom: "nordalh", TrainingBuilding: "special", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 280, Water: 140, Lumber: 200, Stone: 140}, BaseTimeSec: 550,
		FoodUpkeep: 3, Attack: 90, DefInfantry: 50, DefCavalry: 45, Speed: 7, Carry: 25,
	},
	"nordalh_runeforged_forgers": {
		DisplayName: "Runeforged Forgers", Kingdom: "nordalh", TrainingBuilding: "special", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 400, Water: 200, Lumber: 300, Stone: 250}, BaseTimeSec: 700,
		FoodUpkeep: 4, Attack: 85, DefInfantry: 70, DefCavalry: 60, Speed: 4, Carry: 20,
	},
	"nordalh_ulfhednar_champions": {
		DisplayName: "Ulfhednar Champions", Kingdom: "nordalh", TrainingBuilding: "special", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 350, Water: 180, Lumber: 250, Stone: 180}, BaseTimeSec: 650,
		FoodUpkeep: 4, Attack: 100, DefInfantry: 40, DefCavalry: 35, Speed: 6, Carry: 20,
	},

	// ═══════════════════════════════════════════════════════════════════════
	// ZANDRES — Underground / Crystal / Tech
	// ═══════════════════════════════════════════════════════════════════════

	// Barracks
	"zandres_door_wardens": {
		DisplayName: "Door Wardens", Kingdom: "zandres", TrainingBuilding: "barracks", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 100, Water: 50, Lumber: 60, Stone: 40}, BaseTimeSec: 120,
		FoodUpkeep: 1, Attack: 45, DefInfantry: 60, DefCavalry: 45, Speed: 4, Carry: 50,
	},
	"zandres_karst_pikemen": {
		DisplayName: "Karst Pikemen", Kingdom: "zandres", TrainingBuilding: "barracks", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 120, Water: 60, Lumber: 80, Stone: 50}, BaseTimeSec: 150,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 55, DefCavalry: 50, Speed: 4, Carry: 45,
	},
	"zandres_lattice_halberdiers": {
		DisplayName: "Lattice Halberdiers", Kingdom: "zandres", TrainingBuilding: "barracks", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 140, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 180,
		FoodUpkeep: 1, Attack: 65, DefInfantry: 60, DefCavalry: 50, Speed: 4, Carry: 40,
	},
	"zandres_survey_suppressors": {
		DisplayName: "Survey Suppressors", Kingdom: "zandres", TrainingBuilding: "barracks", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 160, Water: 100, Lumber: 120, Stone: 70}, BaseTimeSec: 210,
		FoodUpkeep: 2, Attack: 60, DefInfantry: 65, DefCavalry: 55, Speed: 4, Carry: 35,
	},
	// Stable
	"zandres_cave_strider_riders": {
		DisplayName: "Cave Strider Riders", Kingdom: "zandres", TrainingBuilding: "stable", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 160, Water: 90, Lumber: 100, Stone: 70}, BaseTimeSec: 240,
		FoodUpkeep: 2, Attack: 60, DefInfantry: 25, DefCavalry: 40, Speed: 12, Carry: 60,
	},
	"zandres_beetle_lancers": {
		DisplayName: "Beetle Lancers", Kingdom: "zandres", TrainingBuilding: "stable", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 240, Water: 130, Lumber: 160, Stone: 120}, BaseTimeSec: 360,
		FoodUpkeep: 3, Attack: 75, DefInfantry: 55, DefCavalry: 65, Speed: 6, Carry: 60,
	},
	"zandres_survey_couriers": {
		DisplayName: "Survey Couriers", Kingdom: "zandres", TrainingBuilding: "stable", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 130, Water: 70, Lumber: 80, Stone: 50}, BaseTimeSec: 200,
		FoodUpkeep: 1, Attack: 40, DefInfantry: 20, DefCavalry: 30, Speed: 16, Carry: 40,
	},
	"zandres_burrow_guards": {
		DisplayName: "Burrow Guards", Kingdom: "zandres", TrainingBuilding: "stable", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 220, Water: 120, Lumber: 150, Stone: 100}, BaseTimeSec: 330,
		FoodUpkeep: 3, Attack: 70, DefInfantry: 50, DefCavalry: 60, Speed: 7, Carry: 60,
	},
	// Archery
	"zandres_crystal_boltcasters": {
		DisplayName: "Crystal Boltcasters", Kingdom: "zandres", TrainingBuilding: "archery", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 80, Water: 60, Lumber: 100, Stone: 30}, BaseTimeSec: 130,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 35,
	},
	"zandres_resonance_slingers": {
		DisplayName: "Resonance Slingers", Kingdom: "zandres", TrainingBuilding: "archery", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 70, Water: 50, Lumber: 60, Stone: 30}, BaseTimeSec: 110,
		FoodUpkeep: 1, Attack: 45, DefInfantry: 20, DefCavalry: 15, Speed: 5, Carry: 30,
	},
	"zandres_survey_needlers": {
		DisplayName: "Survey Needlers", Kingdom: "zandres", TrainingBuilding: "archery", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 90, Water: 60, Lumber: 80, Stone: 35}, BaseTimeSec: 140,
		FoodUpkeep: 1, Attack: 50, DefInfantry: 20, DefCavalry: 15, Speed: 6, Carry: 30,
	},
	"zandres_prism_markers": {
		DisplayName: "Prism Markers", Kingdom: "zandres", TrainingBuilding: "archery", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 100, Water: 70, Lumber: 120, Stone: 40}, BaseTimeSec: 160,
		FoodUpkeep: 1, Attack: 45, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 25,
	},
	// Workshop
	"zandres_resonance_ballista": {
		DisplayName: "Resonance Ballista", Kingdom: "zandres", TrainingBuilding: "workshop", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 85, DefInfantry: 20, DefCavalry: 15, Speed: 3, Carry: 0,
	},
	"zandres_drill_ram": {
		DisplayName: "Drill Ram", Kingdom: "zandres", TrainingBuilding: "workshop", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 220, Water: 80, Lumber: 350, Stone: 180}, BaseTimeSec: 450,
		FoodUpkeep: 4, Attack: 75, DefInfantry: 30, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"zandres_stonerail_thrower": {
		DisplayName: "Stone-Rail Thrower", Kingdom: "zandres", TrainingBuilding: "workshop", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 280, Water: 150, Lumber: 400, Stone: 220}, BaseTimeSec: 560,
		FoodUpkeep: 5, Attack: 100, DefInfantry: 15, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"zandres_barrier_cart": {
		DisplayName: "Barrier Cart", Kingdom: "zandres", TrainingBuilding: "workshop", BuildingLevelReq: 12,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 30, DefInfantry: 55, DefCavalry: 45, Speed: 3, Carry: 0,
	},
	// Special
	"zandres_powertech_adepts": {
		DisplayName: "Power-Tech Adepts", Kingdom: "zandres", TrainingBuilding: "special", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 180, Stone: 120}, BaseTimeSec: 500,
		FoodUpkeep: 3, Attack: 85, DefInfantry: 50, DefCavalry: 45, Speed: 5, Carry: 30,
	},
	"zandres_beacon_surveyors": {
		DisplayName: "Beacon Surveyors", Kingdom: "zandres", TrainingBuilding: "special", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 220, Water: 110, Lumber: 160, Stone: 110}, BaseTimeSec: 470,
		FoodUpkeep: 2, Attack: 40, DefInfantry: 55, DefCavalry: 50, Speed: 6, Carry: 20,
	},
	"zandres_capacitor_sentries": {
		DisplayName: "Capacitor Sentries", Kingdom: "zandres", TrainingBuilding: "special", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 350, Water: 200, Lumber: 250, Stone: 200}, BaseTimeSec: 650,
		FoodUpkeep: 4, Attack: 70, DefInfantry: 75, DefCavalry: 65, Speed: 4, Carry: 15,
	},
	"zandres_magnet_lashers": {
		DisplayName: "Magnet Lashers", Kingdom: "zandres", TrainingBuilding: "special", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 300, Water: 160, Lumber: 200, Stone: 150}, BaseTimeSec: 580,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 55, DefCavalry: 50, Speed: 5, Carry: 25,
	},

	// ═══════════════════════════════════════════════════════════════════════
	// LUMUS — Sacred / Radiant / Temple
	// ═══════════════════════════════════════════════════════════════════════

	// Barracks
	"lumus_ringwall_wardens": {
		DisplayName: "Ring-Wall Wardens", Kingdom: "lumus", TrainingBuilding: "barracks", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 100, Water: 50, Lumber: 60, Stone: 40}, BaseTimeSec: 120,
		FoodUpkeep: 1, Attack: 40, DefInfantry: 60, DefCavalry: 45, Speed: 4, Carry: 50,
	},
	"lumus_sun_monks": {
		DisplayName: "Sun Monks", Kingdom: "lumus", TrainingBuilding: "barracks", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 120, Water: 60, Lumber: 80, Stone: 50}, BaseTimeSec: 150,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 55, DefCavalry: 40, Speed: 5, Carry: 45,
	},
	"lumus_prism_guards": {
		DisplayName: "Prism Guards", Kingdom: "lumus", TrainingBuilding: "barracks", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 140, Water: 80, Lumber: 100, Stone: 60}, BaseTimeSec: 180,
		FoodUpkeep: 1, Attack: 60, DefInfantry: 65, DefCavalry: 50, Speed: 4, Carry: 40,
	},
	"lumus_eclipse_wardens": {
		DisplayName: "Eclipse Wardens", Kingdom: "lumus", TrainingBuilding: "barracks", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 160, Water: 100, Lumber: 120, Stone: 70}, BaseTimeSec: 210,
		FoodUpkeep: 2, Attack: 55, DefInfantry: 70, DefCavalry: 60, Speed: 4, Carry: 35,
	},
	// Stable
	"lumus_sunrider_lancers": {
		DisplayName: "Sunrider Lancers", Kingdom: "lumus", TrainingBuilding: "stable", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 180, Water: 100, Lumber: 120, Stone: 80}, BaseTimeSec: 270,
		FoodUpkeep: 2, Attack: 75, DefInfantry: 35, DefCavalry: 55, Speed: 10, Carry: 80,
	},
	"lumus_dawn_couriers": {
		DisplayName: "Dawn Couriers", Kingdom: "lumus", TrainingBuilding: "stable", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 130, Water: 70, Lumber: 80, Stone: 50}, BaseTimeSec: 200,
		FoodUpkeep: 1, Attack: 40, DefInfantry: 20, DefCavalry: 30, Speed: 16, Carry: 40,
	},
	"lumus_halo_riders": {
		DisplayName: "Halo Riders", Kingdom: "lumus", TrainingBuilding: "stable", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 200, Water: 110, Lumber: 140, Stone: 90}, BaseTimeSec: 300,
		FoodUpkeep: 2, Attack: 70, DefInfantry: 35, DefCavalry: 50, Speed: 10, Carry: 60,
	},
	"lumus_whitecloak_escorts": {
		DisplayName: "Whitecloak Escorts", Kingdom: "lumus", TrainingBuilding: "stable", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 220, Water: 120, Lumber: 150, Stone: 100}, BaseTimeSec: 330,
		FoodUpkeep: 3, Attack: 60, DefInfantry: 50, DefCavalry: 65, Speed: 8, Carry: 60,
	},
	// Archery
	"lumus_sunshot_archers": {
		DisplayName: "Sunshot Archers", Kingdom: "lumus", TrainingBuilding: "archery", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 80, Water: 60, Lumber: 100, Stone: 30}, BaseTimeSec: 130,
		FoodUpkeep: 1, Attack: 55, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 35,
	},
	"lumus_halo_chakramists": {
		DisplayName: "Halo Chakramists", Kingdom: "lumus", TrainingBuilding: "archery", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 100, Water: 70, Lumber: 120, Stone: 40}, BaseTimeSec: 160,
		FoodUpkeep: 1, Attack: 60, DefInfantry: 20, DefCavalry: 15, Speed: 6, Carry: 30,
	},
	"lumus_prism_sling_monks": {
		DisplayName: "Prism Sling Monks", Kingdom: "lumus", TrainingBuilding: "archery", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 60, Water: 40, Lumber: 50, Stone: 25}, BaseTimeSec: 100,
		FoodUpkeep: 1, Attack: 40, DefInfantry: 20, DefCavalry: 15, Speed: 5, Carry: 30,
	},
	"lumus_glare_casters": {
		DisplayName: "Glare Casters", Kingdom: "lumus", TrainingBuilding: "archery", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 90, Water: 60, Lumber: 80, Stone: 35}, BaseTimeSec: 140,
		FoodUpkeep: 1, Attack: 45, DefInfantry: 25, DefCavalry: 20, Speed: 5, Carry: 25,
	},
	// Workshop
	"lumus_mirror_ballista": {
		DisplayName: "Mirror Ballista", Kingdom: "lumus", TrainingBuilding: "workshop", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 80, DefInfantry: 20, DefCavalry: 15, Speed: 3, Carry: 0,
	},
	"lumus_sunfire_trebuchet": {
		DisplayName: "Sunfire Trebuchet", Kingdom: "lumus", TrainingBuilding: "workshop", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 400, Stone: 200}, BaseTimeSec: 500,
		FoodUpkeep: 4, Attack: 100, DefInfantry: 15, DefCavalry: 10, Speed: 2, Carry: 0,
	},
	"lumus_glare_tower": {
		DisplayName: "Glare Tower", Kingdom: "lumus", TrainingBuilding: "workshop", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 280, Water: 150, Lumber: 350, Stone: 200}, BaseTimeSec: 550,
		FoodUpkeep: 5, Attack: 70, DefInfantry: 40, DefCavalry: 35, Speed: 2, Carry: 0,
	},
	"lumus_array_cart": {
		DisplayName: "Array Cart", Kingdom: "lumus", TrainingBuilding: "workshop", BuildingLevelReq: 12,
		BaseCost: ResourceCost{Food: 200, Water: 100, Lumber: 300, Stone: 150}, BaseTimeSec: 400,
		FoodUpkeep: 3, Attack: 35, DefInfantry: 50, DefCavalry: 40, Speed: 3, Carry: 0,
	},
	// Special
	"lumus_sunchorus_masters": {
		DisplayName: "Sun-Chorus Masters", Kingdom: "lumus", TrainingBuilding: "special", BuildingLevelReq: 1,
		BaseCost: ResourceCost{Food: 250, Water: 120, Lumber: 180, Stone: 120}, BaseTimeSec: 500,
		FoodUpkeep: 3, Attack: 75, DefInfantry: 60, DefCavalry: 55, Speed: 5, Carry: 25,
	},
	"lumus_radiant_duelists": {
		DisplayName: "Radiant Duelists", Kingdom: "lumus", TrainingBuilding: "special", BuildingLevelReq: 3,
		BaseCost: ResourceCost{Food: 280, Water: 140, Lumber: 200, Stone: 140}, BaseTimeSec: 550,
		FoodUpkeep: 3, Attack: 90, DefInfantry: 50, DefCavalry: 45, Speed: 7, Carry: 25,
	},
	"lumus_eclipse_watch": {
		DisplayName: "Eclipse Watch", Kingdom: "lumus", TrainingBuilding: "special", BuildingLevelReq: 5,
		BaseCost: ResourceCost{Food: 350, Water: 200, Lumber: 250, Stone: 200}, BaseTimeSec: 650,
		FoodUpkeep: 4, Attack: 55, DefInfantry: 75, DefCavalry: 70, Speed: 4, Carry: 15,
	},
	"lumus_prism_adepts": {
		DisplayName: "Prism Adepts", Kingdom: "lumus", TrainingBuilding: "special", BuildingLevelReq: 8,
		BaseCost: ResourceCost{Food: 300, Water: 180, Lumber: 200, Stone: 150}, BaseTimeSec: 580,
		FoodUpkeep: 3, Attack: 50, DefInfantry: 65, DefCavalry: 60, Speed: 5, Carry: 20,
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

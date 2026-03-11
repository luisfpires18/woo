package config

// ---------------------------------------------------------------------------
// Resource economy constants — single source of truth.
// These are exported to the frontend via genconfig → resources.json.
// ---------------------------------------------------------------------------

// StartingResources is the initial amount of each resource in a new village.
const StartingResources = 500.0

// StartingRate is the initial production rate (per second) for each resource.
const StartingRate = 3.0

// StartingStorage is the initial max storage cap for a new village.
// Must equal BaseStorage so the first-village creation is consistent.
const StartingStorage = BaseStorage

// BaseResourceRate is the passive production rate per second at level 0.
const BaseResourceRate = 1.0

// RatePerLevel is the additional production rate per second per building level.
// Total rate = BaseResourceRate + RatePerLevel × sum(all 3 building levels for that resource).
const RatePerLevel = 2.0

// BaseStorage is the base storage capacity per village before any storage buildings.
const BaseStorage = 1200.0

// StoragePerLevel is the additional capacity per level of a storage building.
const StoragePerLevel = 400.0

// StorageBuildingTypes maps each storage building to the resource types it increases capacity for.
var StorageBuildingTypes = map[string][]string{
	"storage":    {"lumber", "stone"},
	"provisions": {"food"},
	"reservoir":  {"water"},
}

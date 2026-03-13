package model

// SpawnRule defines when and where camps appear on the world map.
type SpawnRule struct {
	ID                   int64  `json:"id"`
	Name                 string `json:"name"`
	TerrainTypesJSON     string `json:"terrain_types_json"`
	ZoneTypesJSON        string `json:"zone_types_json"`
	CampTemplatePoolJSON string `json:"camp_template_pool_json"`
	MaxCamps             int    `json:"max_camps"`
	SpawnIntervalSec     int    `json:"spawn_interval_sec"`
	DespawnAfterSec      int    `json:"despawn_after_sec"`
	MinCampDistance      int    `json:"min_camp_distance"`
	MinVillageDistance   int    `json:"min_village_distance"`
	Enabled              bool   `json:"enabled"`
	CreatedAt            string `json:"created_at"`
	UpdatedAt            string `json:"updated_at"`
	UpdatedBy            *int64 `json:"updated_by,omitempty"`
}

// CampTemplatePoolEntry is a weighted entry in a spawn rule's camp template pool.
type CampTemplatePoolEntry struct {
	CampTemplateID int64 `json:"camp_template_id"`
	Weight         int   `json:"weight"`
}

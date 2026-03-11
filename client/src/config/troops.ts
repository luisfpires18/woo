// Frontend troop config — auto-generated from server/internal/config/troops.go
// Re-run `npm run gen-config` (from repo root) when backend config changes.

import type { BuildingType, Kingdom } from '../types/game';
import type { ResourceCost } from './buildings';
import generatedTroops from './generated/troops.json';

export type { ResourceCost };

export interface TroopConfig {
  displayName: string;
  kingdom: Kingdom;
  trainingBuilding: BuildingType;
  buildingLevelReq: number;
  baseCost: ResourceCost;
  baseTimeSec: number;
  foodUpkeep: number;
  attack: number;
  defInfantry: number;
  defCavalry: number;
  speed: number;
  carry: number;
}

export type TroopType =
  // Sylvara
  | 'sylvara_rootguard_spearmen'
  | 'sylvara_thornblade_wardens'
  | 'sylvara_boarhide_axemen'
  | 'sylvara_leafknife_stalkers'
  | 'sylvara_hart_lancers'
  | 'sylvara_boar_riders'
  | 'sylvara_wolf_outriders'
  | 'sylvara_longbow_wardens'
  | 'sylvara_owlshot_hunters'
  | 'sylvara_thorn_javeliners'
  | 'sylvara_stonesling_foresters'
  | 'sylvara_vine_ballista'
  | 'sylvara_log_trebuchet'
  | 'sylvara_root_ram'
  | 'sylvara_sporepot_thrower'
  | 'sylvara_beastmasters'
  | 'sylvara_shapeshifter_scouts'
  | 'sylvara_life_healers'
  | 'sylvara_antler_seers'
  // Arkazia
  | 'arkazia_shield_guard'
  | 'arkazia_pike_guard'
  | 'arkazia_stonebreaker'
  | 'arkazia_rampart_halberdiers'
  | 'arkazia_arknight_lancers'
  | 'arkazia_redcrest_cavaliers'
  | 'arkazia_banner_riders'
  | 'arkazia_ravine_pursuers'
  | 'arkazia_hill_slingers'
  | 'arkazia_ridge_crossbowmen'
  | 'arkazia_javelin_climbers'
  | 'arkazia_crag_bowmen'
  | 'arkazia_mountain_trebuchet'
  | 'arkazia_ram_engineers'
  | 'arkazia_mantlet_pushers'
  | 'arkazia_bridgewright_crew'
  | 'arkazia_banner_knights'
  | 'arkazia_oathsworn_champions'
  | 'arkazia_bastion_marshals'
  | 'arkazia_arknight_captains'
  // Veridor
  | 'veridor_road_legionaries'
  | 'veridor_harbor_pike'
  | 'veridor_cutlass_marines'
  | 'veridor_wharf_axemen'
  | 'veridor_river_lancers'
  | 'veridor_courier_riders'
  | 'veridor_marsh_scouts'
  | 'veridor_road_wardens'
  | 'veridor_deck_arbalesters'
  | 'veridor_highland_longbowmen'
  | 'veridor_harpoon_casters'
  | 'veridor_pavise_marksmen'
  | 'veridor_harbor_ballista'
  | 'veridor_mangonel'
  | 'veridor_firepot_crane'
  | 'veridor_pavise_wagon'
  | 'veridor_hydra_hunters'
  | 'veridor_tidemark_duelists'
  | 'veridor_bluecoat_captains'
  | 'veridor_skiff_raiders'
  // Draxys
  | 'draxys_sandshield_infantry'
  | 'draxys_khopesh_guard'
  | 'draxys_dune_axemen'
  | 'draxys_wadi_lashers'
  | 'draxys_scorpion_riders'
  | 'draxys_dune_lancers'
  | 'draxys_camel_skirmishers'
  | 'draxys_dust_chasers'
  | 'draxys_oasis_rangers'
  | 'draxys_sun_slingers'
  | 'draxys_chakram_dancers'
  | 'draxys_javelin_skirmishers'
  | 'draxys_bolt_thrower'
  | 'draxys_firepot_mangonel'
  | 'draxys_siege_tower'
  | 'draxys_scorpion_cage_wagon'
  | 'draxys_gladiators'
  | 'draxys_netfighters'
  | 'draxys_arena_spearmen'
  | 'draxys_pit_brutes'
  | 'draxys_beast_tamers'
  // Nordalh
  | 'nordalh_hearth_guards'
  | 'nordalh_fjord_spearmen'
  | 'nordalh_iceshore_raiders'
  | 'nordalh_chain_wardens'
  | 'nordalh_direwolf_riders'
  | 'nordalh_elk_lancers'
  | 'nordalh_snow_riders'
  | 'nordalh_fang_cavaliers'
  | 'nordalh_frostbow_hunters'
  | 'nordalh_harpoon_throwers'
  | 'nordalh_cliff_crossbowmen'
  | 'nordalh_storm_slingers'
  | 'nordalh_cliff_ballista'
  | 'nordalh_stone_trebuchet'
  | 'nordalh_ram_sled'
  | 'nordalh_boiling_pitch_crew'
  | 'nordalh_smith_retinues'
  | 'nordalh_coyote_blademasters'
  | 'nordalh_runeforged_forgers'
  | 'nordalh_ulfhednar_champions'
  // Zandres
  | 'zandres_door_wardens'
  | 'zandres_karst_pikemen'
  | 'zandres_lattice_halberdiers'
  | 'zandres_survey_suppressors'
  | 'zandres_cave_strider_riders'
  | 'zandres_beetle_lancers'
  | 'zandres_survey_couriers'
  | 'zandres_burrow_guards'
  | 'zandres_crystal_boltcasters'
  | 'zandres_resonance_slingers'
  | 'zandres_survey_needlers'
  | 'zandres_prism_markers'
  | 'zandres_resonance_ballista'
  | 'zandres_drill_ram'
  | 'zandres_stonerail_thrower'
  | 'zandres_barrier_cart'
  | 'zandres_powertech_adepts'
  | 'zandres_beacon_surveyors'
  | 'zandres_capacitor_sentries'
  | 'zandres_magnet_lashers'
  // Lumus
  | 'lumus_ringwall_wardens'
  | 'lumus_sun_monks'
  | 'lumus_prism_guards'
  | 'lumus_eclipse_wardens'
  | 'lumus_sunrider_lancers'
  | 'lumus_dawn_couriers'
  | 'lumus_halo_riders'
  | 'lumus_whitecloak_escorts'
  | 'lumus_sunshot_archers'
  | 'lumus_halo_chakramists'
  | 'lumus_prism_sling_monks'
  | 'lumus_glare_casters'
  | 'lumus_mirror_ballista'
  | 'lumus_sunfire_trebuchet'
  | 'lumus_glare_tower'
  | 'lumus_array_cart'
  | 'lumus_sunchorus_masters'
  | 'lumus_radiant_duelists'
  | 'lumus_eclipse_watch'
  | 'lumus_prism_adepts';

/**
 * Authoritative troop config imported from generated JSON.
 * Keys are validated at build time via the TroopType cast.
 */
export const TROOP_CONFIGS = generatedTroops as unknown as Record<TroopType, TroopConfig>;

/** Get troops available for a given kingdom. */
export function getTroopsForKingdom(kingdom: Kingdom): [TroopType, TroopConfig][] {
  return (Object.entries(TROOP_CONFIGS) as [TroopType, TroopConfig][]).filter(
    ([, cfg]) => cfg.kingdom === kingdom,
  );
}

/** Get troops trainable at a specific building for a given kingdom. */
export function getTroopsForBuilding(buildingType: string, kingdom: Kingdom): [TroopType, TroopConfig][] {
  return (Object.entries(TROOP_CONFIGS) as [TroopType, TroopConfig][]).filter(
    ([, cfg]) => cfg.kingdom === kingdom && cfg.trainingBuilding === buildingType,
  );
}

/** Check if a building type is a military building (has troops mapped to it). */
export function isMilitaryBuilding(buildingType: string): boolean {
  return Object.values(TROOP_CONFIGS).some((cfg) => cfg.trainingBuilding === buildingType);
}

/** Format seconds into a human-readable duration (e.g., "2m 30s"). */
export function formatDuration(seconds: number): string {
  if (seconds < 60) return `${seconds}s`;
  const mins = Math.floor(seconds / 60);
  const secs = seconds % 60;
  if (mins < 60) return secs > 0 ? `${mins}m ${secs}s` : `${mins}m`;
  const hrs = Math.floor(mins / 60);
  const remMins = mins % 60;
  return remMins > 0 ? `${hrs}h ${remMins}m` : `${hrs}h`;
}

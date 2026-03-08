// Frontend troop config — mirrors server/internal/config/troops.go
// Must be kept in sync manually when backend config changes.

import type { BuildingType, Kingdom } from '../types/game';

export interface ResourceCost {
  food: number;
  water: number;
  lumber: number;
  stone: number;
}

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
  | 'iron_legionary'
  | 'crossbowman'
  | 'mountain_knight'
  | 'shieldbearer'
  | 'gladiator'
  | 'battering_ram'
  | 'mountain_scout';

export const TROOP_CONFIGS: Record<TroopType, TroopConfig> = {
  iron_legionary: {
    displayName: 'Iron Legionary',
    kingdom: 'arkazia',
    trainingBuilding: 'barracks',
    buildingLevelReq: 1,
    baseCost: { food: 100, water: 50, lumber: 60, stone: 40 },
    baseTimeSec: 120,
    foodUpkeep: 1,
    attack: 50,
    defInfantry: 55,
    defCavalry: 40,
    speed: 5,
    carry: 55,
  },
  crossbowman: {
    displayName: 'Crossbowman',
    kingdom: 'arkazia',
    trainingBuilding: 'barracks',
    buildingLevelReq: 3,
    baseCost: { food: 80, water: 60, lumber: 100, stone: 30 },
    baseTimeSec: 150,
    foodUpkeep: 1,
    attack: 50,
    defInfantry: 30,
    defCavalry: 25,
    speed: 5,
    carry: 35,
  },
  mountain_knight: {
    displayName: 'Mountain Knight',
    kingdom: 'arkazia',
    trainingBuilding: 'stable',
    buildingLevelReq: 1,
    baseCost: { food: 150, water: 100, lumber: 50, stone: 120 },
    baseTimeSec: 300,
    foodUpkeep: 2,
    attack: 80,
    defInfantry: 50,
    defCavalry: 60,
    speed: 8,
    carry: 80,
  },
  shieldbearer: {
    displayName: 'Shieldbearer',
    kingdom: 'arkazia',
    trainingBuilding: 'barracks',
    buildingLevelReq: 5,
    baseCost: { food: 120, water: 80, lumber: 50, stone: 150 },
    baseTimeSec: 180,
    foodUpkeep: 1,
    attack: 30,
    defInfantry: 80,
    defCavalry: 70,
    speed: 4,
    carry: 40,
  },
  gladiator: {
    displayName: 'Gladiator',
    kingdom: 'arkazia',
    trainingBuilding: 'special',
    buildingLevelReq: 1,
    baseCost: { food: 200, water: 100, lumber: 80, stone: 200 },
    baseTimeSec: 600,
    foodUpkeep: 3,
    attack: 95,
    defInfantry: 60,
    defCavalry: 50,
    speed: 5,
    carry: 30,
  },
  battering_ram: {
    displayName: 'Battering Ram',
    kingdom: 'arkazia',
    trainingBuilding: 'barracks',
    buildingLevelReq: 10,
    baseCost: { food: 300, water: 200, lumber: 400, stone: 200 },
    baseTimeSec: 480,
    foodUpkeep: 4,
    attack: 70,
    defInfantry: 20,
    defCavalry: 30,
    speed: 2,
    carry: 0,
  },
  mountain_scout: {
    displayName: 'Mountain Scout',
    kingdom: 'arkazia',
    trainingBuilding: 'barracks',
    buildingLevelReq: 1,
    baseCost: { food: 50, water: 40, lumber: 30, stone: 20 },
    baseTimeSec: 60,
    foodUpkeep: 1,
    attack: 20,
    defInfantry: 10,
    defCavalry: 10,
    speed: 12,
    carry: 20,
  },
};

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

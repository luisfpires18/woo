// Frontend building config — mirrors server/internal/config/buildings.go
// Must be kept in sync manually when backend config changes.

import type { BuildingType } from '../types/game';
import type { BuildingInfo } from '../types/api';

export interface ResourceCost {
  food: number;
  water: number;
  lumber: number;
  stone: number;
}

export interface BuildingPrerequisite {
  buildingType: BuildingType;
  minLevel: number;
}

export interface BuildingConfig {
  displayName: string;
  baseCost: ResourceCost;
  baseTimeSec: number;
  scalingFactor: number;
  timeFactor: number;
  maxLevel: number;
  prerequisites: BuildingPrerequisite[];
}

/** Shared base cost for all 12 resource field buildings. */
const RESOURCE_FIELD_COST: ResourceCost = { food: 60, water: 40, lumber: 80, stone: 50 };

function resourceBuilding(displayName: string): BuildingConfig {
  return {
    displayName,
    baseCost: RESOURCE_FIELD_COST,
    baseTimeSec: 120,
    scalingFactor: 1.5,
    timeFactor: 1.5,
    maxLevel: 20,
    prerequisites: [],
  };
}

export const BUILDING_CONFIGS: Record<BuildingType, BuildingConfig> = {
  // --- Village buildings ---
  town_hall: {
    displayName: 'Town Hall',
    baseCost: { food: 100, water: 200, lumber: 200, stone: 200 },
    baseTimeSec: 300,
    scalingFactor: 1.7,
    timeFactor: 1.7,
    maxLevel: 20,
    prerequisites: [],
  },
  barracks: {
    displayName: 'Barracks',
    baseCost: { food: 80, water: 200, lumber: 150, stone: 100 },
    baseTimeSec: 300,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 20,
    prerequisites: [{ buildingType: 'town_hall', minLevel: 3 }],
  },
  stable: {
    displayName: 'Stable',
    baseCost: { food: 120, water: 300, lumber: 200, stone: 150 },
    baseTimeSec: 480,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    prerequisites: [
      { buildingType: 'town_hall', minLevel: 5 },
      { buildingType: 'barracks', minLevel: 5 },
    ],
  },
  archery: {
    displayName: 'Archery',
    baseCost: { food: 80, water: 150, lumber: 200, stone: 80 },
    baseTimeSec: 300,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    prerequisites: [{ buildingType: 'town_hall', minLevel: 3 }],
  },
  workshop: {
    displayName: 'Workshop',
    baseCost: { food: 100, water: 200, lumber: 300, stone: 250 },
    baseTimeSec: 600,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    prerequisites: [
      { buildingType: 'town_hall', minLevel: 7 },
      { buildingType: 'barracks', minLevel: 5 },
    ],
  },
  special: {
    displayName: 'Special',
    baseCost: { food: 200, water: 300, lumber: 250, stone: 300 },
    baseTimeSec: 900,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    prerequisites: [
      { buildingType: 'town_hall', minLevel: 10 },
      { buildingType: 'barracks', minLevel: 7 },
      { buildingType: 'stable', minLevel: 5 },
    ],
  },

  // --- Resource field buildings (3 per resource type) ---
  food_1: resourceBuilding('Food Field I'),
  food_2: resourceBuilding('Food Field II'),
  food_3: resourceBuilding('Food Field III'),
  water_1: resourceBuilding('Water Field I'),
  water_2: resourceBuilding('Water Field II'),
  water_3: resourceBuilding('Water Field III'),
  lumber_1: resourceBuilding('Lumber Field I'),
  lumber_2: resourceBuilding('Lumber Field II'),
  lumber_3: resourceBuilding('Lumber Field III'),
  stone_1: resourceBuilding('Stone Field I'),
  stone_2: resourceBuilding('Stone Field II'),
  stone_3: resourceBuilding('Stone Field III'),
};

/** All 12 resource building type IDs. */
export const RESOURCE_BUILDING_TYPES: ReadonlySet<string> = new Set([
  'food_1', 'food_2', 'food_3',
  'water_1', 'water_2', 'water_3',
  'lumber_1', 'lumber_2', 'lumber_3',
  'stone_1', 'stone_2', 'stone_3',
]);

/** Resource building IDs grouped by resource type, in display order. */
export const RESOURCE_BUILDING_GROUPS: { resource: string; label: string; emoji: string; types: BuildingType[] }[] = [
  { resource: 'food', label: 'Food', emoji: '🌾', types: ['food_1', 'food_2', 'food_3'] },
  { resource: 'water', label: 'Water', emoji: '💧', types: ['water_1', 'water_2', 'water_3'] },
  { resource: 'lumber', label: 'Lumber', emoji: '🪵', types: ['lumber_1', 'lumber_2', 'lumber_3'] },
  { resource: 'stone', label: 'Stone', emoji: '🪨', types: ['stone_1', 'stone_2', 'stone_3'] },
];

/** Building types shown in the main "Village Buildings" grid (everything except resource fields). */
export const VILLAGE_BUILDING_TYPES: ReadonlySet<string> = new Set(
  (Object.keys(BUILDING_CONFIGS) as BuildingType[]).filter(
    (t) => !RESOURCE_BUILDING_TYPES.has(t),
  ),
);

/**
 * Calculate cost to upgrade a building to the given level.
 * Mirrors server/internal/config/buildings.go CostAtLevel.
 */
export function costAtLevel(buildingType: BuildingType, level: number): ResourceCost {
  const cfg = BUILDING_CONFIGS[buildingType];
  const mult = Math.pow(cfg.scalingFactor, level - 1);
  return {
    food: Math.round(cfg.baseCost.food * mult),
    water: Math.round(cfg.baseCost.water * mult),
    lumber: Math.round(cfg.baseCost.lumber * mult),
    stone: Math.round(cfg.baseCost.stone * mult),
  };
}

/**
 * Calculate build time in seconds for upgrading a building to the given level.
 * Mirrors server/internal/config/buildings.go TimeAtLevel.
 */
export function timeAtLevel(buildingType: BuildingType, level: number): number {
  const cfg = BUILDING_CONFIGS[buildingType];
  return Math.round(cfg.baseTimeSec * Math.pow(cfg.timeFactor, level - 1));
}

export interface PrerequisiteCheck {
  buildingType: BuildingType;
  displayName: string;
  minLevel: number;
  currentLevel: number;
  met: boolean;
}

/**
 * Check whether all prerequisites for a building are satisfied.
 * Returns per-prerequisite details plus an overall `allMet` flag.
 * Pass an optional `displayNameFn` to resolve admin-configured kingdom display names;
 * otherwise falls back to the hardcoded BUILDING_CONFIGS names.
 */
export function checkPrerequisites(
  buildingType: BuildingType,
  buildings: BuildingInfo[],
  displayNameFn?: (type: string) => string,
): { allMet: boolean; checks: PrerequisiteCheck[] } {
  const cfg = BUILDING_CONFIGS[buildingType];
  const levelMap = new Map<string, number>();
  for (const b of buildings) {
    levelMap.set(b.building_type, b.level);
  }

  const resolve = displayNameFn ?? ((t: string) => BUILDING_CONFIGS[t as BuildingType]?.displayName ?? t);

  const checks: PrerequisiteCheck[] = cfg.prerequisites.map((prereq) => {
    const current = levelMap.get(prereq.buildingType) ?? 0;
    return {
      buildingType: prereq.buildingType,
      displayName: resolve(prereq.buildingType),
      minLevel: prereq.minLevel,
      currentLevel: current,
      met: current >= prereq.minLevel,
    };
  });

  return {
    allMet: checks.every((c) => c.met),
    checks,
  };
}

/**
 * Format seconds into a human-readable duration string (e.g. "5m 30s", "1h 15m 30s").
 */
export function formatDuration(totalSeconds: number): string {
  if (totalSeconds <= 0) return '0s';
  const h = Math.floor(totalSeconds / 3600);
  const m = Math.floor((totalSeconds % 3600) / 60);
  const s = totalSeconds % 60;
  const parts: string[] = [];
  if (h > 0) parts.push(`${h}h`);
  if (m > 0) parts.push(`${m}m`);
  if (s > 0 || parts.length === 0) parts.push(`${s}s`);
  return parts.join(' ');
}

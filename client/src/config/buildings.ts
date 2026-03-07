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
  kingdomOnly?: string; // empty = all kingdoms
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
  warehouse: {
    displayName: 'Warehouse',
    baseCost: { food: 50, water: 120, lumber: 120, stone: 100 },
    baseTimeSec: 180,
    scalingFactor: 1.6,
    timeFactor: 1.6,
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
  forge: {
    displayName: 'Forge',
    baseCost: { food: 100, water: 250, lumber: 180, stone: 200 },
    baseTimeSec: 480,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 10,
    prerequisites: [
      { buildingType: 'town_hall', minLevel: 5 },
      { buildingType: 'barracks', minLevel: 3 },
    ],
  },
  rune_altar: {
    displayName: 'Rune Altar',
    baseCost: { food: 150, water: 300, lumber: 250, stone: 250 },
    baseTimeSec: 600,
    scalingFactor: 1.9,
    timeFactor: 1.9,
    maxLevel: 10,
    prerequisites: [
      { buildingType: 'town_hall', minLevel: 7 },
      { buildingType: 'forge', minLevel: 3 },
    ],
  },
  walls: {
    displayName: 'Walls',
    baseCost: { food: 50, water: 150, lumber: 100, stone: 200 },
    baseTimeSec: 240,
    scalingFactor: 1.6,
    timeFactor: 1.6,
    maxLevel: 20,
    prerequisites: [{ buildingType: 'town_hall', minLevel: 2 }],
  },
  marketplace: {
    displayName: 'Marketplace',
    baseCost: { food: 80, water: 180, lumber: 180, stone: 120 },
    baseTimeSec: 300,
    scalingFactor: 1.6,
    timeFactor: 1.6,
    maxLevel: 15,
    prerequisites: [
      { buildingType: 'town_hall', minLevel: 5 },
      { buildingType: 'warehouse', minLevel: 3 },
    ],
  },
  embassy: {
    displayName: 'Embassy',
    baseCost: { food: 100, water: 200, lumber: 200, stone: 200 },
    baseTimeSec: 480,
    scalingFactor: 1.7,
    timeFactor: 1.7,
    maxLevel: 10,
    prerequisites: [{ buildingType: 'town_hall', minLevel: 8 }],
  },
  watchtower: {
    displayName: 'Watchtower',
    baseCost: { food: 60, water: 150, lumber: 100, stone: 150 },
    baseTimeSec: 240,
    scalingFactor: 1.6,
    timeFactor: 1.6,
    maxLevel: 10,
    prerequisites: [
      { buildingType: 'town_hall', minLevel: 3 },
      { buildingType: 'walls', minLevel: 1 },
    ],
  },
  dock: {
    displayName: 'Dock',
    baseCost: { food: 100, water: 250, lumber: 300, stone: 150 },
    baseTimeSec: 480,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    kingdomOnly: 'veridor',
    prerequisites: [{ buildingType: 'town_hall', minLevel: 6 }],
  },
  grove_sanctum: {
    displayName: 'Grove Sanctum',
    baseCost: { food: 150, water: 200, lumber: 300, stone: 200 },
    baseTimeSec: 480,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    kingdomOnly: 'sylvara',
    prerequisites: [{ buildingType: 'town_hall', minLevel: 6 }],
  },
  colosseum: {
    displayName: 'Colosseum',
    baseCost: { food: 100, water: 300, lumber: 200, stone: 300 },
    baseTimeSec: 480,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    kingdomOnly: 'arkazia',
    prerequisites: [{ buildingType: 'town_hall', minLevel: 6 }],
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
 */
export function checkPrerequisites(
  buildingType: BuildingType,
  buildings: BuildingInfo[],
): { allMet: boolean; checks: PrerequisiteCheck[] } {
  const cfg = BUILDING_CONFIGS[buildingType];
  const levelMap = new Map<string, number>();
  for (const b of buildings) {
    levelMap.set(b.building_type, b.level);
  }

  const checks: PrerequisiteCheck[] = cfg.prerequisites.map((prereq) => {
    const current = levelMap.get(prereq.buildingType) ?? 0;
    return {
      buildingType: prereq.buildingType,
      displayName: BUILDING_CONFIGS[prereq.buildingType].displayName,
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

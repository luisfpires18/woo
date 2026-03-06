// Frontend building config — mirrors server/internal/config/buildings.go
// Must be kept in sync manually when backend config changes.

import type { BuildingType } from '../types/game';
import type { BuildingInfo } from '../types/api';

export interface ResourceCost {
  iron: number;
  wood: number;
  stone: number;
  food: number;
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

export const BUILDING_CONFIGS: Record<BuildingType, BuildingConfig> = {
  town_hall: {
    displayName: 'Town Hall',
    baseCost: { iron: 200, wood: 200, stone: 200, food: 100 },
    baseTimeSec: 300,
    scalingFactor: 1.7,
    timeFactor: 1.7,
    maxLevel: 20,
    prerequisites: [],
  },
  iron_mine: {
    displayName: 'Iron Mine',
    baseCost: { iron: 100, wood: 80, stone: 50, food: 30 },
    baseTimeSec: 120,
    scalingFactor: 1.5,
    timeFactor: 1.5,
    maxLevel: 20,
    prerequisites: [],
  },
  lumber_mill: {
    displayName: 'Lumber Mill',
    baseCost: { iron: 80, wood: 100, stone: 50, food: 30 },
    baseTimeSec: 120,
    scalingFactor: 1.5,
    timeFactor: 1.5,
    maxLevel: 20,
    prerequisites: [],
  },
  quarry: {
    displayName: 'Quarry',
    baseCost: { iron: 80, wood: 50, stone: 100, food: 30 },
    baseTimeSec: 120,
    scalingFactor: 1.5,
    timeFactor: 1.5,
    maxLevel: 20,
    prerequisites: [],
  },
  farm: {
    displayName: 'Farm',
    baseCost: { iron: 50, wood: 80, stone: 50, food: 20 },
    baseTimeSec: 120,
    scalingFactor: 1.5,
    timeFactor: 1.5,
    maxLevel: 20,
    prerequisites: [],
  },
  warehouse: {
    displayName: 'Warehouse',
    baseCost: { iron: 120, wood: 120, stone: 100, food: 50 },
    baseTimeSec: 180,
    scalingFactor: 1.6,
    timeFactor: 1.6,
    maxLevel: 20,
    prerequisites: [],
  },
  barracks: {
    displayName: 'Barracks',
    baseCost: { iron: 200, wood: 150, stone: 100, food: 80 },
    baseTimeSec: 300,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 20,
    prerequisites: [{ buildingType: 'town_hall', minLevel: 3 }],
  },
  stable: {
    displayName: 'Stable',
    baseCost: { iron: 300, wood: 200, stone: 150, food: 120 },
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
    baseCost: { iron: 250, wood: 180, stone: 200, food: 100 },
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
    baseCost: { iron: 300, wood: 250, stone: 250, food: 150 },
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
    baseCost: { iron: 150, wood: 100, stone: 200, food: 50 },
    baseTimeSec: 240,
    scalingFactor: 1.6,
    timeFactor: 1.6,
    maxLevel: 20,
    prerequisites: [{ buildingType: 'town_hall', minLevel: 2 }],
  },
  marketplace: {
    displayName: 'Marketplace',
    baseCost: { iron: 180, wood: 180, stone: 120, food: 80 },
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
    baseCost: { iron: 200, wood: 200, stone: 200, food: 100 },
    baseTimeSec: 480,
    scalingFactor: 1.7,
    timeFactor: 1.7,
    maxLevel: 10,
    prerequisites: [{ buildingType: 'town_hall', minLevel: 8 }],
  },
  watchtower: {
    displayName: 'Watchtower',
    baseCost: { iron: 150, wood: 100, stone: 150, food: 60 },
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
    baseCost: { iron: 250, wood: 300, stone: 150, food: 100 },
    baseTimeSec: 480,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    kingdomOnly: 'veridor',
    prerequisites: [{ buildingType: 'town_hall', minLevel: 6 }],
  },
  grove_sanctum: {
    displayName: 'Grove Sanctum',
    baseCost: { iron: 200, wood: 300, stone: 200, food: 150 },
    baseTimeSec: 480,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    kingdomOnly: 'sylvara',
    prerequisites: [{ buildingType: 'town_hall', minLevel: 6 }],
  },
  colosseum: {
    displayName: 'Colosseum',
    baseCost: { iron: 300, wood: 200, stone: 300, food: 100 },
    baseTimeSec: 480,
    scalingFactor: 1.8,
    timeFactor: 1.8,
    maxLevel: 15,
    kingdomOnly: 'arkazia',
    prerequisites: [{ buildingType: 'town_hall', minLevel: 6 }],
  },
};

/** Building types that produce resources (shown in "Resource Fields" section). */
export const RESOURCE_BUILDING_TYPES: ReadonlySet<string> = new Set([
  'iron_mine',
  'lumber_mill',
  'quarry',
  'farm',
]);

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
    iron: Math.round(cfg.baseCost.iron * mult),
    wood: Math.round(cfg.baseCost.wood * mult),
    stone: Math.round(cfg.baseCost.stone * mult),
    food: Math.round(cfg.baseCost.food * mult),
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

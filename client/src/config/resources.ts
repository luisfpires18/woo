// Frontend resource economy config — auto-generated from server/internal/config/resources.go
// Re-run `npm run gen-config` (from repo root) when backend config changes.

import generatedResources from './generated/resources.json';

/** Resource economy constants shared with the server. */
export interface ResourceEconomy {
  /** Initial amount of each resource in a new village. */
  startingResources: number;
  /** Initial production rate (per second) for each resource. */
  startingRate: number;
  /** Initial gold amount for a new player. */
  startingGold: number;
  /** Passive production rate per second at building level 0. */
  baseResourceRate: number;
  /** Additional production rate per second per building level. */
  ratePerLevel: number;
  /** Base storage capacity before any storage buildings. */
  baseStorage: number;
  /** Additional capacity per level of a storage building. */
  storagePerLevel: number;
  /** Maps storage buildings to the resource types they increase. */
  storageBuildings: Record<string, string[]>;
}

export const RESOURCE_ECONOMY = generatedResources as unknown as ResourceEconomy;

// Re-export individual constants for convenience.
export const STARTING_RESOURCES = RESOURCE_ECONOMY.startingResources;
export const STARTING_RATE = RESOURCE_ECONOMY.startingRate;
export const STARTING_GOLD = RESOURCE_ECONOMY.startingGold;
export const BASE_RESOURCE_RATE = RESOURCE_ECONOMY.baseResourceRate;
export const RATE_PER_LEVEL = RESOURCE_ECONOMY.ratePerLevel;
export const BASE_STORAGE = RESOURCE_ECONOMY.baseStorage;
export const STORAGE_PER_LEVEL = RESOURCE_ECONOMY.storagePerLevel;
export const STORAGE_BUILDINGS = RESOURCE_ECONOMY.storageBuildings;

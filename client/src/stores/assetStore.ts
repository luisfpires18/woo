import { create } from 'zustand';
import type { GameAsset, TroopDisplayConfig, BuildingDisplayConfig } from '../types/api';
import { fetchGameAssets, fetchTroopDisplayConfigs, fetchBuildingDisplayConfigs } from '../services/admin';
import { getSpriteUrl } from '../utils/spriteUrl';

/** Convert a troop display config into a synthetic GameAsset so GameIcon works. */
export function troopConfigToAsset(c: TroopDisplayConfig): GameAsset {
  return {
    id: c.troop_type,
    category: 'unit',
    display_name: c.display_name,
    default_icon: c.default_icon,
    sprite_url: getSpriteUrl({ kind: 'unit', id: c.troop_type, kingdom: c.kingdom }),
    updated_at: c.updated_at,
  };
}

/** Convert a building display config into a synthetic GameAsset so GameIcon works. */
export function buildingConfigToAsset(c: BuildingDisplayConfig): GameAsset {
  return {
    id: `${c.building_type}_${c.kingdom}`,
    category: 'building',
    display_name: c.display_name,
    default_icon: c.default_icon,
    sprite_url: getSpriteUrl({ kind: 'building', id: c.building_type, kingdom: c.kingdom }),
    updated_at: c.updated_at,
  };
}

interface AssetState {
  assets: GameAsset[];
  loaded: boolean;
  loading: boolean;

  /** Fetch all game assets (called once on init). */
  load: () => Promise<void>;

  /** Find an asset by its id (e.g. "iron_mine"). */
  getById: (id: string) => GameAsset | undefined;

  /** Replace a single asset in the cache (after update). */
  upsert: (asset: GameAsset) => void;

  /** Add a new asset to the cache (after creating a variant). */
  addAsset: (asset: GameAsset) => void;

  /** Remove an asset from the cache by id. */
  removeAsset: (id: string) => void;
}

export const useAssetStore = create<AssetState>((set, get) => ({
  assets: [],
  loaded: false,
  loading: false,

  load: async () => {
    if (get().loaded || get().loading) return;
    set({ loading: true });
    try {
      const [assetResp, troopResp, buildingResp] = await Promise.all([
        fetchGameAssets(),
        fetchTroopDisplayConfigs().catch(() => ({ configs: [] as TroopDisplayConfig[] })),
        fetchBuildingDisplayConfigs().catch(() => ({ configs: [] as BuildingDisplayConfig[] })),
      ]);
      // Assign convention sprite URLs to game assets
      const gameAssets = assetResp.assets.map((a) => ({
        ...a,
        sprite_url: getSpriteUrl({ kind: a.category as any, id: a.id }),
      }));
      // Merge game assets + troop display configs + building display configs
      const troopAssets = troopResp.configs.map(troopConfigToAsset);
      const buildingAssets = buildingResp.configs.map(buildingConfigToAsset);
      set({ assets: [...gameAssets, ...troopAssets, ...buildingAssets], loaded: true });
    } catch {
      // Silently fail — components fall back to emoji
    } finally {
      set({ loading: false });
    }
  },

  getById: (id: string) => get().assets.find((a) => a.id === id),

  upsert: (asset: GameAsset) =>
    set((s) => ({
      assets: s.assets.map((a) => (a.id === asset.id ? asset : a)),
    })),

  addAsset: (asset: GameAsset) =>
    set((s) => ({
      assets: [...s.assets, asset],
    })),

  removeAsset: (id: string) =>
    set((s) => ({
      assets: s.assets.filter((a) => a.id !== id),
    })),
}));

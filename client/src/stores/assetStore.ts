import { create } from 'zustand';
import type { GameAsset } from '../types/api';
import { fetchGameAssets } from '../services/admin';

interface AssetState {
  assets: GameAsset[];
  loaded: boolean;
  loading: boolean;

  /** Fetch all game assets (called once on init). */
  load: () => Promise<void>;

  /** Find an asset by its id (e.g. "iron_mine"). */
  getById: (id: string) => GameAsset | undefined;

  /** Replace a single asset in the cache (after upload / delete). */
  upsert: (asset: GameAsset) => void;

  /** Add a new asset to the cache (after creating a variant). */
  addAsset: (asset: GameAsset) => void;

  /** Remove an asset from the cache by id. */
  removeAsset: (id: string) => void;

  /** Clear sprite_url for a given asset id. */
  clearSprite: (id: string) => void;
}

export const useAssetStore = create<AssetState>((set, get) => ({
  assets: [],
  loaded: false,
  loading: false,

  load: async () => {
    if (get().loaded || get().loading) return;
    set({ loading: true });
    try {
      const resp = await fetchGameAssets();
      set({ assets: resp.assets, loaded: true });
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

  clearSprite: (id: string) =>
    set((s) => ({
      assets: s.assets.map((a) =>
        a.id === id ? { ...a, sprite_url: null } : a,
      ),
    })),
}));

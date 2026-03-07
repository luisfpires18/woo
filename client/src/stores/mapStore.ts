// Map store — manages loaded map chunks, viewport position, and selected tile

import { create } from 'zustand';
import type { MapTile } from '../types/map';
import { fetchMapChunk } from '../services/map';

interface MapState {
  /** All loaded tiles indexed by "x,y" key */
  tiles: Map<string, MapTile>;
  /** Current viewport center */
  centerX: number;
  centerY: number;
  /** Currently selected tile */
  selectedTile: MapTile | null;
  /** Loading state */
  loading: boolean;
  /** Error message */
  error: string | null;

  /** Load a map chunk from the server and merge into tile cache */
  loadChunk: (x: number, y: number, range?: number) => Promise<void>;
  /** Set viewport center and auto-load surrounding tiles */
  setCenter: (x: number, y: number) => void;
  /** Select a tile for the info panel */
  selectTile: (tile: MapTile | null) => void;
  /** Clear all cached tiles */
  clearTiles: () => void;
}

function tileKey(x: number, y: number): string {
  return `${x},${y}`;
}

export const useMapStore = create<MapState>((set, get) => ({
  tiles: new Map(),
  centerX: 0,
  centerY: 0,
  selectedTile: null,
  loading: false,
  error: null,

  loadChunk: async (x: number, y: number, range: number = 15) => {
    set({ loading: true, error: null });
    try {
      const response = await fetchMapChunk(x, y, range);
      const tiles = new Map(get().tiles);
      for (const tile of response.tiles) {
        tiles.set(tileKey(tile.x, tile.y), tile);
      }
      set({ tiles, loading: false });
    } catch (err) {
      set({
        loading: false,
        error: err instanceof Error ? err.message : 'Failed to load map',
      });
    }
  },

  setCenter: (x: number, y: number) => {
    set({ centerX: x, centerY: y });
    // Auto-load chunk at the new center
    get().loadChunk(x, y, 15);
  },

  selectTile: (tile: MapTile | null) => {
    set({ selectedTile: tile });
  },

  clearTiles: () => {
    set({ tiles: new Map(), selectedTile: null });
  },
}));

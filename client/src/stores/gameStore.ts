import { create } from 'zustand';
import type { VillageResponse, VillageListItem } from '../types/api';

interface GameState {
  /** Full village detail (buildings + resources) for the currently viewed village */
  currentVillage: VillageResponse | null;
  /** List of the player's villages (summary data) */
  villages: VillageListItem[];
  /** Whether the village list has been fetched at least once */
  villagesLoaded: boolean;
  /** Player-level gold balance (shared across all villages) */
  playerGold: number;

  setCurrentVillage: (village: VillageResponse | null) => void;
  setVillages: (villages: VillageListItem[]) => void;
  setPlayerGold: (gold: number) => void;
}

export const useGameStore = create<GameState>((set) => ({
  currentVillage: null,
  villages: [],
  villagesLoaded: false,
  playerGold: 0,

  setCurrentVillage: (village) => {
    // Sync player gold from village response (gold is per-player, returned with every village fetch)
    const gold = village?.gold ?? 0;
    set({ currentVillage: village, playerGold: gold });
  },
  setVillages: (villages) => set({ villages, villagesLoaded: true }),
  setPlayerGold: (gold) => set({ playerGold: gold }),
}));

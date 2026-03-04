import { create } from 'zustand';
import type { VillageResponse, VillageListItem } from '../types/api';

interface GameState {
  /** Full village detail (buildings + resources) for the currently viewed village */
  currentVillage: VillageResponse | null;
  /** List of the player's villages (summary data) */
  villages: VillageListItem[];

  setCurrentVillage: (village: VillageResponse | null) => void;
  setVillages: (villages: VillageListItem[]) => void;
}

export const useGameStore = create<GameState>((set) => ({
  currentVillage: null,
  villages: [],

  setCurrentVillage: (village) => set({ currentVillage: village }),
  setVillages: (villages) => set({ villages }),
}));

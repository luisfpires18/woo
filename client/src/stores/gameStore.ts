import { create } from 'zustand';
import type { Resources } from '../types/game';

interface GameState {
  resources: Resources | null;
  setResources: (resources: Resources) => void;
}

export const useGameStore = create<GameState>((set) => ({
  resources: null,
  setResources: (resources) => set({ resources }),
}));

import { create } from 'zustand';
import type { SeasonResponse, SeasonDetailResponse } from '../types/api';
import { fetchSeasons, fetchMySeasons } from '../services/season';

interface SeasonState {
  /** All seasons (any status) */
  seasons: SeasonResponse[];
  /** Seasons the current player has joined */
  mySeasons: SeasonDetailResponse[];
  /** Whether seasons have been fetched at least once */
  loaded: boolean;
  loading: boolean;

  /** Fetch all seasons, optionally filtered by status */
  loadSeasons: (status?: string) => Promise<void>;
  /** Fetch seasons the current player is part of */
  loadMySeasons: () => Promise<void>;
  /** Update a single season in the local list (after join/admin action) */
  updateSeason: (season: SeasonResponse) => void;
  /** Remove a season from the local list */
  removeSeason: (seasonId: number) => void;
  /** Reset all season state */
  reset: () => void;
}

export const useSeasonStore = create<SeasonState>((set) => ({
  seasons: [],
  mySeasons: [],
  loaded: false,
  loading: false,

  loadSeasons: async (status?: string) => {
    set({ loading: true });
    try {
      const seasons = await fetchSeasons(status);
      set({ seasons, loaded: true });
    } finally {
      set({ loading: false });
    }
  },

  loadMySeasons: async () => {
    try {
      const mySeasons = await fetchMySeasons();
      set({ mySeasons });
    } catch {
      // Silently fail — non-critical
    }
  },

  updateSeason: (season) => {
    set((state) => ({
      seasons: state.seasons.map((s) => (s.id === season.id ? season : s)),
    }));
  },

  removeSeason: (seasonId) => {
    set((state) => ({
      seasons: state.seasons.filter((s) => s.id !== seasonId),
    }));
  },

  reset: () => set({ seasons: [], mySeasons: [], loaded: false, loading: false }),
}));

import { create } from 'zustand';
import type { Kingdom } from '../types/game';
import { useAuthStore } from './authStore';

const VALID_KINGDOMS: readonly Kingdom[] = [
  'veridor', 'sylvara', 'arkazia', 'draxys',
  'zandres', 'lumus', 'nordalh', 'drakanith',
];

function isValidKingdom(v: string | undefined | null): v is Kingdom {
  return !!v && (VALID_KINGDOMS as readonly string[]).includes(v);
}

function applyKingdom(kingdom: Kingdom | null) {
  if (kingdom) {
    document.documentElement.setAttribute('data-kingdom', kingdom);
  } else {
    document.documentElement.removeAttribute('data-kingdom');
  }
}

interface ThemeState {
  kingdom: Kingdom | null;
  /** Initialise on app boot — reads kingdom from authStore and applies. */
  init: () => void;
  /** Set the kingdom theme explicitly (e.g. after kingdom selection). */
  setKingdom: (k: Kingdom | null) => void;
}

export const useThemeStore = create<ThemeState>((set) => ({
  kingdom: null,

  init: () => {
    const raw = useAuthStore.getState().player?.kingdom;
    const kingdom = isValidKingdom(raw) ? raw : null;
    applyKingdom(kingdom);
    set({ kingdom });
  },

  setKingdom: (kingdom) => {
    applyKingdom(kingdom);
    set({ kingdom });
  },
}));

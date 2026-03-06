import { create } from 'zustand';

export type Theme = 'dark' | 'light';

const STORAGE_KEY = 'woo-theme';

function getSystemTheme(): Theme {
  if (typeof window === 'undefined') return 'dark';
  return window.matchMedia('(prefers-color-scheme: light)').matches ? 'light' : 'dark';
}

function getInitialTheme(): Theme {
  if (typeof window === 'undefined') return 'dark';
  const stored = localStorage.getItem(STORAGE_KEY);
  if (stored === 'dark' || stored === 'light') return stored;
  return getSystemTheme();
}

function applyTheme(theme: Theme) {
  document.documentElement.setAttribute('data-theme', theme);
}

interface ThemeState {
  theme: Theme;
  /** Initialise on app boot — reads localStorage / OS preference and applies. */
  init: () => void;
  /** Toggle between dark ↔ light. */
  toggle: () => void;
  /** Set a specific theme. */
  setTheme: (t: Theme) => void;
}

export const useThemeStore = create<ThemeState>((set, get) => ({
  theme: 'dark', // will be overwritten by init()

  init: () => {
    const theme = getInitialTheme();
    applyTheme(theme);
    set({ theme });

    // Listen to OS preference changes (only matters when user hasn't set a manual preference)
    window.matchMedia('(prefers-color-scheme: light)').addEventListener('change', (e) => {
      const stored = localStorage.getItem(STORAGE_KEY);
      if (stored) return; // user has a manual preference — ignore OS change
      const next: Theme = e.matches ? 'light' : 'dark';
      applyTheme(next);
      set({ theme: next });
    });
  },

  toggle: () => {
    const next: Theme = get().theme === 'dark' ? 'light' : 'dark';
    applyTheme(next);
    localStorage.setItem(STORAGE_KEY, next);
    set({ theme: next });
  },

  setTheme: (t: Theme) => {
    applyTheme(t);
    localStorage.setItem(STORAGE_KEY, t);
    set({ theme: t });
  },
}));

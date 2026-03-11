import { create } from 'zustand';
import type { PlayerInfo, RegisterRequest, LoginRequest } from '../types/api';
import {
  registerPlayer,
  loginPlayer,
  logoutPlayer,
} from '../services/auth';
import { getMe } from '../services/player';

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  player: PlayerInfo | null;
  isAuthenticated: boolean;
  /** Whether restore() has been called at least once */
  hydrated: boolean;

  /** Call on app init to restore session from localStorage */
  restore: () => void;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
  /** Update the player object (e.g. after kingdom selection) */
  setPlayer: (player: PlayerInfo) => void;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  accessToken: null,
  refreshToken: null,
  player: null,
  isAuthenticated: false,
  hydrated: false,

  restore: () => {
    const accessToken = localStorage.getItem('access_token');
    const refreshToken = localStorage.getItem('refresh_token');
    const playerStr = localStorage.getItem('player');

    if (accessToken && playerStr) {
      try {
        // Validate cached session with the server before marking as authenticated.
        // This avoids firing protected queries with stale tokens after a backend/DB reset.
        set({ accessToken, refreshToken, player: null, isAuthenticated: false, hydrated: false });

        // Refresh player data from server to avoid stale localStorage
        getMe()
          .then((freshPlayer) => {
            localStorage.setItem('player', JSON.stringify(freshPlayer));
            set({
              accessToken,
              refreshToken,
              player: freshPlayer,
              isAuthenticated: true,
              hydrated: true,
            });
          })
          .catch(() => {
            // Token invalid or server unreachable — clear session
            localStorage.removeItem('access_token');
            localStorage.removeItem('refresh_token');
            localStorage.removeItem('player');
            set({
              accessToken: null,
              refreshToken: null,
              player: null,
              isAuthenticated: false,
              hydrated: true,
            });
          });
        return;
      } catch {
        // Corrupted data — clear it
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('player');
      }
    }
    set({ hydrated: true });
  },

  login: async (data) => {
    const resp = await loginPlayer(data);
    localStorage.setItem('access_token', resp.access_token);
    localStorage.setItem('refresh_token', resp.refresh_token);
    localStorage.setItem('player', JSON.stringify(resp.player));
    set({
      accessToken: resp.access_token,
      refreshToken: resp.refresh_token,
      player: resp.player,
      isAuthenticated: true,
    });
  },

  register: async (data) => {
    const resp = await registerPlayer(data);
    localStorage.setItem('access_token', resp.access_token);
    localStorage.setItem('refresh_token', resp.refresh_token);
    localStorage.setItem('player', JSON.stringify(resp.player));
    set({
      accessToken: resp.access_token,
      refreshToken: resp.refresh_token,
      player: resp.player,
      isAuthenticated: true,
    });
  },

  logout: async () => {
    const rt = get().refreshToken;
    if (rt) {
      try {
        await logoutPlayer(rt);
      } catch {
        // Server logout failed — still clear local state
      }
    }
    localStorage.removeItem('access_token');
    localStorage.removeItem('refresh_token');
    localStorage.removeItem('player');
    set({
      accessToken: null,
      refreshToken: null,
      player: null,
      isAuthenticated: false,
    });
  },

  setPlayer: (player) => {
    localStorage.setItem('player', JSON.stringify(player));
    set({ player });
  },
}));

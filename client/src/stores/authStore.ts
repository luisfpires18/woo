import { create } from 'zustand';
import type { PlayerInfo, RegisterRequest, LoginRequest } from '../types/api';
import {
  registerPlayer,
  loginPlayer,
  logoutPlayer,
} from '../services/auth';

interface AuthState {
  accessToken: string | null;
  refreshToken: string | null;
  player: PlayerInfo | null;
  isAuthenticated: boolean;

  /** Call on app init to restore session from localStorage */
  restore: () => void;
  login: (data: LoginRequest) => Promise<void>;
  register: (data: RegisterRequest) => Promise<void>;
  logout: () => Promise<void>;
}

export const useAuthStore = create<AuthState>((set, get) => ({
  accessToken: null,
  refreshToken: null,
  player: null,
  isAuthenticated: false,

  restore: () => {
    const accessToken = localStorage.getItem('access_token');
    const refreshToken = localStorage.getItem('refresh_token');
    const playerStr = localStorage.getItem('player');

    if (accessToken && playerStr) {
      try {
        const player = JSON.parse(playerStr) as PlayerInfo;
        set({ accessToken, refreshToken, player, isAuthenticated: true });
      } catch {
        // Corrupted data — clear it
        localStorage.removeItem('access_token');
        localStorage.removeItem('refresh_token');
        localStorage.removeItem('player');
      }
    }
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
}));

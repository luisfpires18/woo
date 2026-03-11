import { beforeEach, describe, expect, it, vi } from 'vitest';
import type { PlayerInfo } from '../types/api';

const getMeMock = vi.fn<() => Promise<PlayerInfo>>();

vi.mock('../services/player', () => ({
  getMe: () => getMeMock(),
}));

import { useAuthStore } from './authStore';

const cachedPlayer: PlayerInfo = {
  id: 123,
  username: 'wright',
  email: 'wright@woo.local',
  kingdom: 'arkazia',
  role: 'player',
};

describe('authStore.restore', () => {
  beforeEach(() => {
    localStorage.clear();
    getMeMock.mockReset();
    useAuthStore.setState({
      accessToken: null,
      refreshToken: null,
      player: null,
      isAuthenticated: false,
      hydrated: false,
    });
  });

  it('does not mark the user authenticated before server validation succeeds', async () => {
    let resolveGetMe: (player: PlayerInfo) => void = () => {};
    getMeMock.mockReturnValue(new Promise<PlayerInfo>((resolve) => {
      resolveGetMe = resolve;
    }));

    localStorage.setItem('access_token', 'access-token');
    localStorage.setItem('refresh_token', 'refresh-token');
    localStorage.setItem('player', JSON.stringify(cachedPlayer));

    useAuthStore.getState().restore();

    expect(useAuthStore.getState().isAuthenticated).toBe(false);
    expect(useAuthStore.getState().player).toBeNull();
    expect(useAuthStore.getState().hydrated).toBe(false);

    resolveGetMe(cachedPlayer);
    await Promise.resolve();
    await Promise.resolve();

    expect(useAuthStore.getState().isAuthenticated).toBe(true);
    expect(useAuthStore.getState().player).toEqual(cachedPlayer);
    expect(useAuthStore.getState().hydrated).toBe(true);
  });

  it('clears stale local session when server validation fails', async () => {
    getMeMock.mockRejectedValue(new Error('unauthorized'));

    localStorage.setItem('access_token', 'stale-access');
    localStorage.setItem('refresh_token', 'stale-refresh');
    localStorage.setItem('player', JSON.stringify(cachedPlayer));

    useAuthStore.getState().restore();
    await Promise.resolve();
    await Promise.resolve();

    expect(localStorage.getItem('access_token')).toBeNull();
    expect(localStorage.getItem('refresh_token')).toBeNull();
    expect(localStorage.getItem('player')).toBeNull();
    expect(useAuthStore.getState().isAuthenticated).toBe(false);
    expect(useAuthStore.getState().player).toBeNull();
    expect(useAuthStore.getState().hydrated).toBe(true);
  });
});
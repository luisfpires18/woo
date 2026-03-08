import { api } from './api';
import type { PlayerInfo, ChooseKingdomResponse, PlayerProfileResponse } from '../types/api';

export async function getMe(): Promise<PlayerInfo> {
  const res = await api.get<{ player: PlayerInfo }>('/player/me');
  return res.player;
}

export async function chooseKingdom(kingdom: string): Promise<ChooseKingdomResponse> {
  return api.put<ChooseKingdomResponse>('/player/kingdom', { kingdom });
}

export async function getProfile(): Promise<PlayerProfileResponse> {
  const res = await api.get<{ profile: PlayerProfileResponse }>('/player/profile');
  return res.profile;
}

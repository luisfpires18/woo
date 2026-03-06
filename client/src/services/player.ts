import { api } from './api';
import type { ChooseKingdomResponse } from '../types/api';

export async function chooseKingdom(kingdom: string): Promise<ChooseKingdomResponse> {
  return api.put<ChooseKingdomResponse>('/player/kingdom', { kingdom });
}

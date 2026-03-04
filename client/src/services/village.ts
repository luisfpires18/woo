// Village API service — fetch villages and village details

import { api } from './api';
import type { VillageResponse, VillageListItem } from '../types/api';

export async function fetchVillages(): Promise<VillageListItem[]> {
  return api.get<VillageListItem[]>('/villages');
}

export async function fetchVillage(id: number): Promise<VillageResponse> {
  return api.get<VillageResponse>(`/villages/${id}`);
}

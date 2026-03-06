// Village API service — fetch villages and village details

import { api } from './api';
import type {
  VillageResponse,
  VillageListItem,
  BuildingQueueResponse,
  BuildingCostResponse,
} from '../types/api';

export async function fetchVillages(): Promise<VillageListItem[]> {
  return api.get<VillageListItem[]>('/villages');
}

export async function fetchVillage(id: number): Promise<VillageResponse> {
  return api.get<VillageResponse>(`/villages/${id}`);
}

export async function startUpgrade(
  villageId: number,
  buildingType: string,
): Promise<BuildingQueueResponse> {
  return api.post<BuildingQueueResponse>(`/villages/${villageId}/upgrade`, {
    building_type: buildingType,
  });
}

export async function getUpgradeCost(
  villageId: number,
  buildingType: string,
): Promise<BuildingCostResponse> {
  return api.get<BuildingCostResponse>(
    `/villages/${villageId}/upgrade/cost?building_type=${encodeURIComponent(buildingType)}`,
  );
}

export async function cancelUpgrade(
  villageId: number,
  queueId: number,
): Promise<void> {
  await api.delete(`/villages/${villageId}/upgrade/${queueId}`);
}

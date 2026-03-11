// Village API service — fetch villages and village details

import { api } from './api';
import type {
  VillageResponse,
  VillageListItem,
  BuildingQueueResponse,
  BuildingCostResponse,
  BuildingDisplayConfigListResponse,
  ResourceBuildingConfigListResponse,
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

export async function renameVillage(
  villageId: number,
  name: string,
): Promise<VillageListItem> {
  return api.put<VillageListItem>(`/villages/${villageId}/name`, { name });
}

/** Admin: instantly complete a building queue item. */
export async function instantCompleteBuild(queueId: number): Promise<void> {
  await api.post(`/admin/building/${queueId}/complete`, {});
}

/** Fetch building display configs, optionally filtered by kingdom. */
export async function fetchBuildingDisplayConfigs(kingdom?: string): Promise<BuildingDisplayConfigListResponse> {
  const qs = kingdom ? `?kingdom=${encodeURIComponent(kingdom)}` : '';
  return api.get<BuildingDisplayConfigListResponse>(`/building-display-configs${qs}`);
}

/** Fetch resource building configs, optionally filtered by kingdom. */
export async function fetchResourceBuildingConfigs(kingdom?: string): Promise<ResourceBuildingConfigListResponse> {
  const qs = kingdom ? `?kingdom=${encodeURIComponent(kingdom)}` : '';
  return api.get<ResourceBuildingConfigListResponse>(`/resource-building-configs${qs}`);
}

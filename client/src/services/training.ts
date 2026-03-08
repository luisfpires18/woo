// Training API service — train troops, fetch costs, manage training queue

import { api } from './api';
import type {
  TrainingQueueResponse,
  TrainingCostResponse,
  TroopInfo,
} from '../types/api';

export async function startTraining(
  villageId: number,
  troopType: string,
  quantity: number,
): Promise<TrainingQueueResponse> {
  return api.post<TrainingQueueResponse>(`/villages/${villageId}/train`, {
    troop_type: troopType,
    quantity,
  });
}

export async function getTrainingCost(
  villageId: number,
  troopType: string,
  quantity: number,
): Promise<TrainingCostResponse> {
  return api.get<TrainingCostResponse>(
    `/villages/${villageId}/train/cost?troop_type=${encodeURIComponent(troopType)}&quantity=${quantity}`,
  );
}

export async function cancelTraining(
  villageId: number,
  queueId: number,
): Promise<void> {
  await api.delete(`/villages/${villageId}/train/${queueId}`);
}

export async function fetchTroops(villageId: number): Promise<{ troops: TroopInfo[] }> {
  return api.get<{ troops: TroopInfo[] }>(`/villages/${villageId}/troops`);
}

/** Admin-only: instantly complete a training queue item. */
export async function instantCompleteTraining(queueId: number): Promise<void> {
  await api.post(`/admin/training/${queueId}/complete`, {});
}

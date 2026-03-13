// Camp & Expedition API service — fetches camps, dispatches expeditions, retrieves battle reports

import { api } from './api';
import type {
  CampResponse,
  DispatchExpeditionRequest,
  ExpeditionResponse,
  BattleReportResponse,
  BattleReplayResponse,
  // Admin types
  BeastTemplateResponse,
  CreateBeastTemplateRequest,
  UpdateBeastTemplateRequest,
  CampTemplateResponse,
  CreateCampTemplateRequest,
  UpdateCampTemplateRequest,
  SpawnRuleResponse,
  CreateSpawnRuleRequest,
  UpdateSpawnRuleRequest,
  RewardTableResponse,
  CreateRewardTableRequest,
  UpdateRewardTableRequest,
  RewardEntryRequest,
  BattleTuningResponse,
} from '../types/api';

// ── Player endpoints ────────────────────────────────────────────────────────

/** List all active camps on the world map. */
export async function fetchCamps(): Promise<CampResponse[]> {
  return api.get<CampResponse[]>('/camps');
}

/** Get a specific camp with its beasts. */
export async function fetchCamp(campId: number): Promise<CampResponse> {
  return api.get<CampResponse>(`/camps/${campId}`);
}

/** Dispatch troops from a village to attack a camp. */
export async function dispatchExpedition(
  villageId: number,
  request: DispatchExpeditionRequest,
): Promise<ExpeditionResponse> {
  return api.post<ExpeditionResponse>(`/villages/${villageId}/expeditions`, request);
}

/** Get all expeditions for the current player. */
export async function fetchExpeditions(): Promise<ExpeditionResponse[]> {
  return api.get<ExpeditionResponse[]>('/expeditions');
}

/** Get a battle report by battle ID. */
export async function fetchBattleReport(battleId: number): Promise<BattleReportResponse> {
  return api.get<BattleReportResponse>(`/battles/${battleId}`);
}

/** Get the raw replay JSON for a battle. */
export async function fetchBattleReplay(battleId: number): Promise<BattleReplayResponse> {
  return api.get<BattleReplayResponse>(`/battles/${battleId}/replay`);
}

// ── Admin: Beast Templates ──────────────────────────────────────────────────

export async function fetchBeastTemplates(): Promise<BeastTemplateResponse[]> {
  return api.get<BeastTemplateResponse[]>('/admin/beast-templates');
}

export async function fetchBeastTemplate(id: number): Promise<BeastTemplateResponse> {
  return api.get<BeastTemplateResponse>(`/admin/beast-templates/${id}`);
}

export async function createBeastTemplate(req: CreateBeastTemplateRequest): Promise<BeastTemplateResponse> {
  return api.post<BeastTemplateResponse>('/admin/beast-templates', req);
}

export async function updateBeastTemplate(id: number, req: UpdateBeastTemplateRequest): Promise<BeastTemplateResponse> {
  return api.put<BeastTemplateResponse>(`/admin/beast-templates/${id}`, req);
}

export async function deleteBeastTemplate(id: number): Promise<void> {
  await api.delete(`/admin/beast-templates/${id}`);
}

// ── Admin: Camp Templates ───────────────────────────────────────────────────

export async function fetchCampTemplates(): Promise<CampTemplateResponse[]> {
  return api.get<CampTemplateResponse[]>('/admin/camp-templates');
}

export async function fetchCampTemplate(id: number): Promise<CampTemplateResponse> {
  return api.get<CampTemplateResponse>(`/admin/camp-templates/${id}`);
}

export async function createCampTemplate(req: CreateCampTemplateRequest): Promise<CampTemplateResponse> {
  return api.post<CampTemplateResponse>('/admin/camp-templates', req);
}

export async function updateCampTemplate(id: number, req: UpdateCampTemplateRequest): Promise<CampTemplateResponse> {
  return api.put<CampTemplateResponse>(`/admin/camp-templates/${id}`, req);
}

export async function deleteCampTemplate(id: number): Promise<void> {
  await api.delete(`/admin/camp-templates/${id}`);
}

// ── Admin: Spawn Rules ──────────────────────────────────────────────────────

export async function fetchSpawnRules(): Promise<SpawnRuleResponse[]> {
  return api.get<SpawnRuleResponse[]>('/admin/spawn-rules');
}

export async function fetchSpawnRule(id: number): Promise<SpawnRuleResponse> {
  return api.get<SpawnRuleResponse>(`/admin/spawn-rules/${id}`);
}

export async function createSpawnRule(req: CreateSpawnRuleRequest): Promise<SpawnRuleResponse> {
  return api.post<SpawnRuleResponse>('/admin/spawn-rules', req);
}

export async function updateSpawnRule(id: number, req: UpdateSpawnRuleRequest): Promise<SpawnRuleResponse> {
  return api.put<SpawnRuleResponse>(`/admin/spawn-rules/${id}`, req);
}

export async function deleteSpawnRule(id: number): Promise<void> {
  await api.delete(`/admin/spawn-rules/${id}`);
}

// ── Admin: Reward Tables ────────────────────────────────────────────────────

export async function fetchRewardTables(): Promise<RewardTableResponse[]> {
  return api.get<RewardTableResponse[]>('/admin/reward-tables');
}

export async function fetchRewardTable(id: number): Promise<RewardTableResponse> {
  return api.get<RewardTableResponse>(`/admin/reward-tables/${id}`);
}

export async function createRewardTable(req: CreateRewardTableRequest): Promise<RewardTableResponse> {
  return api.post<RewardTableResponse>('/admin/reward-tables', req);
}

export async function updateRewardTable(id: number, req: UpdateRewardTableRequest): Promise<RewardTableResponse> {
  return api.put<RewardTableResponse>(`/admin/reward-tables/${id}`, req);
}

export async function deleteRewardTable(id: number): Promise<void> {
  await api.delete(`/admin/reward-tables/${id}`);
}

export async function addRewardEntry(tableId: number, entry: RewardEntryRequest): Promise<void> {
  await api.post(`/admin/reward-tables/${tableId}/entries`, entry);
}

export async function deleteRewardEntry(tableId: number, entryId: number): Promise<void> {
  await api.delete(`/admin/reward-tables/${tableId}/entries/${entryId}`);
}

// ── Admin: Battle Tuning ────────────────────────────────────────────────────

export async function fetchBattleTuning(): Promise<BattleTuningResponse> {
  return api.get<BattleTuningResponse>('/admin/battle-tuning');
}

export async function updateBattleTuning(tuning: Partial<BattleTuningResponse>): Promise<BattleTuningResponse> {
  return api.put<BattleTuningResponse>('/admin/battle-tuning', tuning);
}

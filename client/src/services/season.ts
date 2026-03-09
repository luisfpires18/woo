import { api } from './api';
import type {
  SeasonResponse,
  SeasonDetailResponse,
  CreateSeasonRequest,
  JoinSeasonResponse,
} from '../types/api';

// ── Public (no auth required) ────────────────────────────────────────────────

export async function fetchPublicSeasons(status?: string): Promise<SeasonResponse[]> {
  const query = status ? `?status=${status}` : '';
  const res = await api.get<{ seasons: SeasonResponse[] }>(`/public/seasons${query}`, true);
  return res.seasons;
}

// ── Player-facing (auth required) ─────────────────────────────────────────────

export async function fetchMySeasons(): Promise<SeasonDetailResponse[]> {
  const res = await api.get<{ seasons: SeasonDetailResponse[] }>('/seasons/my');
  return res.seasons;
}

export async function joinSeason(id: number, kingdom: string): Promise<JoinSeasonResponse> {
  return api.post<JoinSeasonResponse>(`/seasons/${id}/join`, { kingdom });
}

// ── Admin ────────────────────────────────────────────────────────────────────

export async function adminFetchSeasons(status?: string): Promise<SeasonResponse[]> {
  const query = status ? `?status=${status}` : '';
  const res = await api.get<{ seasons: SeasonResponse[] }>(`/admin/seasons${query}`);
  return res.seasons;
}

export async function adminCreateSeason(data: CreateSeasonRequest): Promise<SeasonResponse> {
  const res = await api.post<{ season: SeasonResponse }>('/admin/seasons', data);
  return res.season;
}

export async function adminDeleteSeason(id: number): Promise<void> {
  await api.delete<{ deleted: boolean }>(`/admin/seasons/${id}`);
}

export async function adminLaunchSeason(id: number): Promise<SeasonResponse> {
  const res = await api.post<{ season: SeasonResponse }>(`/admin/seasons/${id}/launch`, {});
  return res.season;
}

export async function adminEndSeason(id: number): Promise<SeasonResponse> {
  const res = await api.post<{ season: SeasonResponse }>(`/admin/seasons/${id}/end`, {});
  return res.season;
}

export async function adminArchiveSeason(id: number): Promise<SeasonResponse> {
  const res = await api.post<{ season: SeasonResponse }>(`/admin/seasons/${id}/archive`, {});
  return res.season;
}

import { api } from './api';
import type {
  PlayerListResponse,
  WorldConfigResponse,
  StatsResponse,
  AnnouncementResponse,
  CreateAnnouncementRequest,
  SetConfigRequest,
  UpdateRoleRequest,
  GameAssetListResponse,
  GameAsset,
} from '../types/api';

// --- Player management ---

export async function fetchPlayers(offset = 0, limit = 20): Promise<PlayerListResponse> {
  return api.get<PlayerListResponse>(`/admin/players?offset=${offset}&limit=${limit}`);
}

export async function updatePlayerRole(id: number, role: 'player' | 'admin'): Promise<void> {
  const body: UpdateRoleRequest = { role };
  await api.patch<{ message: string }>(`/admin/players/${id}/role`, body);
}

// --- World config ---

export async function fetchWorldConfig(): Promise<WorldConfigResponse> {
  return api.get<WorldConfigResponse>('/admin/config');
}

export async function setWorldConfig(key: string, value: string): Promise<void> {
  const body: SetConfigRequest = { value };
  await api.put<{ message: string }>(`/admin/config/${key}`, body);
}

// --- Server stats ---

export async function fetchStats(): Promise<StatsResponse> {
  return api.get<StatsResponse>('/admin/stats');
}

// --- Announcements ---

export async function fetchAnnouncements(): Promise<AnnouncementResponse[]> {
  return api.get<AnnouncementResponse[]>('/admin/announcements');
}

export async function createAnnouncement(req: CreateAnnouncementRequest): Promise<AnnouncementResponse> {
  return api.post<AnnouncementResponse>('/admin/announcements', req);
}

export async function deleteAnnouncement(id: number): Promise<void> {
  await api.delete<{ message: string }>(`/admin/announcements/${id}`);
}

// --- Game assets ---

export async function fetchGameAssets(): Promise<GameAssetListResponse> {
  return api.get<GameAssetListResponse>('/assets');
}

export async function uploadSprite(id: string, file: File): Promise<GameAsset> {
  const formData = new FormData();
  formData.append('file', file);

  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/admin/assets/${id}/sprite`, {
    method: 'POST',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    body: formData,
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: 'Upload failed' }));
    throw new Error(body.error || `HTTP ${response.status}`);
  }

  const result = await response.json();
  return result.data;
}

export async function deleteSprite(id: string): Promise<void> {
  await api.delete<{ message: string }>(`/admin/assets/${id}/sprite`);
}

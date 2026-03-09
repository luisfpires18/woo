import { api } from './api';
import type {
  PlayerListResponse,
  StatsResponse,
  AnnouncementResponse,
  CreateAnnouncementRequest,
  UpdateRoleRequest,
  GameAssetListResponse,
  GameAsset,
  BuildingDisplayConfig,
  BuildingDisplayConfigListResponse,
} from '../types/api';

// --- Player management ---

export async function fetchPlayers(offset = 0, limit = 20): Promise<PlayerListResponse> {
  return api.get<PlayerListResponse>(`/admin/players?offset=${offset}&limit=${limit}`);
}

export async function updatePlayerRole(id: number, role: 'player' | 'admin'): Promise<void> {
  const body: UpdateRoleRequest = { role };
  await api.patch<{ message: string }>(`/admin/players/${id}/role`, body);
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

export async function createGameAsset(data: {
  id: string;
  category: string;
  display_name: string;
  default_icon?: string;
}): Promise<GameAsset> {
  return api.post<GameAsset>('/admin/assets', data);
}

export async function deleteGameAsset(id: string): Promise<void> {
  await api.delete<{ message: string }>(`/admin/assets/${id}`);
}

// --- Building display configs ---

export async function fetchBuildingDisplayConfigs(): Promise<BuildingDisplayConfigListResponse> {
  return api.get<BuildingDisplayConfigListResponse>('/building-display-configs');
}

export async function updateBuildingDisplayConfig(
  id: number,
  data: { display_name: string; description: string; default_icon: string },
): Promise<BuildingDisplayConfig> {
  const res = await api.put<{ data: BuildingDisplayConfig }>(`/admin/building-displays/${id}`, data);
  return res.data;
}

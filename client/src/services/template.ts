import { api } from './api';
import type { TemplateInfo, MapTemplate } from '../types/map';

// --- Map template API ---

/** List all saved templates (metadata only). */
export async function listTemplates(): Promise<TemplateInfo[]> {
  const res = await api.get<{ templates: TemplateInfo[] }>('/admin/templates');
  return res.templates;
}

/** Get a full template including tile data. */
export async function getTemplate(name: string): Promise<MapTemplate> {
  return api.get<MapTemplate>(`/admin/templates/${encodeURIComponent(name)}`);
}

/** Create a new blank template. */
export async function createTemplate(name: string, description: string, mapSize?: number): Promise<void> {
  await api.post<{ message: string }>('/admin/templates', { name, description, map_size: mapSize || 0 });
}

/** Resize a template (preserves existing tiles within new bounds). */
export async function resizeTemplate(name: string, mapSize: number): Promise<void> {
  await api.put<{ message: string }>(`/admin/templates/${encodeURIComponent(name)}/resize`, { map_size: mapSize });
}

/** Delete a template. */
export async function deleteTemplate(name: string): Promise<void> {
  await api.delete<{ message: string }>(`/admin/templates/${encodeURIComponent(name)}`);
}

/** Paint terrain on a template. */
export interface TemplateTileUpdate {
  x: number;
  y: number;
  terrain_type: string;
}

export async function updateTemplateTerrain(name: string, tiles: TemplateTileUpdate[]): Promise<void> {
  await api.put<{ message: string }>(`/admin/templates/${encodeURIComponent(name)}/terrain`, { tiles });
}

/** Paint zones on a template. */
export interface TemplateZoneUpdate {
  x: number;
  y: number;
  kingdom_zone: string;
}

export async function updateTemplateZones(name: string, tiles: TemplateZoneUpdate[]): Promise<void> {
  await api.put<{ message: string }>(`/admin/templates/${encodeURIComponent(name)}/zones`, { tiles });
}

/** Apply a template to the live world map. */
export async function applyTemplate(name: string): Promise<void> {
  await api.post<{ message: string }>(`/admin/templates/${encodeURIComponent(name)}/apply`, { confirm: true });
}

/** Export a template as a downloadable JSON file. */
export async function exportTemplate(name: string): Promise<void> {
  const token = localStorage.getItem('access_token');
  const response = await fetch(`/api/admin/templates/${encodeURIComponent(name)}/export`, {
    headers: token ? { Authorization: `Bearer ${token}` } : {},
  });

  if (!response.ok) {
    throw new Error(`Export failed: HTTP ${response.status}`);
  }

  const blob = await response.blob();
  const url = URL.createObjectURL(blob);
  const a = document.createElement('a');
  a.href = url;
  a.download = `${name}.json`;
  document.body.appendChild(a);
  a.click();
  document.body.removeChild(a);
  URL.revokeObjectURL(url);
}

/** Import a template from a JSON file upload. */
export async function importTemplate(file: File): Promise<void> {
  const formData = new FormData();
  formData.append('file', file);

  const token = localStorage.getItem('access_token');
  const response = await fetch('/api/admin/templates/import', {
    method: 'POST',
    headers: token ? { Authorization: `Bearer ${token}` } : {},
    body: formData,
  });

  if (!response.ok) {
    const body = await response.json().catch(() => ({ error: 'Import failed' }));
    throw new Error(body.error || `HTTP ${response.status}`);
  }
}

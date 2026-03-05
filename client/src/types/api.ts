// API response and request types

export interface ApiResponse<T> {
  data: T;
  error?: string;
}

export interface ApiError {
  error: string;
  details?: string;
}

export interface PaginatedResponse<T> {
  data: T[];
  total: number;
  page: number;
  perPage: number;
}

// Auth types — mirrors server/internal/dto/auth.go

export interface RegisterRequest {
  username: string;
  email: string;
  password: string;
  kingdom: string;
}

export interface LoginRequest {
  email: string;
  password: string;
}

export interface RefreshRequest {
  refresh_token: string;
}

export interface AuthResponse {
  access_token: string;
  refresh_token: string;
  player: PlayerInfo;
}

export interface PlayerInfo {
  id: number;
  username: string;
  email: string;
  kingdom: string;
  role: 'player' | 'admin';
}

// Village types — mirrors server/internal/dto/village.go

export interface VillageResponse {
  id: number;
  player_id: number;
  name: string;
  x: number;
  y: number;
  is_capital: boolean;
  buildings: BuildingInfo[];
  resources: ResourcesResponse;
}

export interface BuildingInfo {
  id: number;
  building_type: string;
  level: number;
}

export interface ResourcesResponse {
  iron: number;
  wood: number;
  stone: number;
  food: number;
  iron_rate: number;
  wood_rate: number;
  stone_rate: number;
  food_rate: number;
  food_consumption: number;
  max_storage: number;
}

export interface VillageListItem {
  id: number;
  name: string;
  x: number;
  y: number;
  is_capital: boolean;
}

// Admin types — mirrors server/internal/dto/admin.go

export interface PlayerListItem {
  id: number;
  username: string;
  email: string;
  kingdom: string;
  role: 'player' | 'admin';
  created_at: string;
  last_login_at?: string;
}

export interface PlayerListResponse {
  players: PlayerListItem[];
  total: number;
  offset: number;
  limit: number;
}

export interface UpdateRoleRequest {
  role: 'player' | 'admin';
}

export interface WorldConfigEntry {
  key: string;
  value: string;
  description?: string;
  updated_at: string;
}

export interface WorldConfigResponse {
  configs: WorldConfigEntry[];
}

export interface SetConfigRequest {
  value: string;
}

export interface StatsResponse {
  total_players: number;
  total_villages: number;
}

export interface CreateAnnouncementRequest {
  title: string;
  content: string;
  expires_at?: string;
}

export interface AnnouncementResponse {
  id: number;
  title: string;
  content: string;
  author_id: number;
  created_at: string;
  expires_at?: string;
}

// Game asset types — mirrors server/internal/dto/admin.go

export type AssetCategory = 'building' | 'resource' | 'unit';

export interface GameAsset {
  id: string;
  category: AssetCategory;
  display_name: string;
  default_icon: string;
  sprite_url: string | null;
  sprite_width: number;
  sprite_height: number;
  updated_at: string;
}

export interface GameAssetListResponse {
  assets: GameAsset[];
}

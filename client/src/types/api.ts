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

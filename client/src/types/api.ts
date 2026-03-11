// API response and request types

import type { BuildingType, Kingdom } from './game';
import type { TroopType } from '../config/troops';

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
}

export interface LoginRequest {
  login: string;
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
  kingdom: Kingdom | ''; // empty string means kingdom not yet chosen
  role: 'player' | 'admin';
}

export interface ChooseKingdomRequest {
  kingdom: string;
}

export interface ChooseKingdomResponse {
  player: PlayerInfo;
  village_id: number;
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
  build_queue: BuildingQueueResponse[];
  troops: TroopInfo[];
  training_queue: TrainingQueueResponse[];
}

export interface BuildingInfo {
  id: number;
  building_type: BuildingType;
  level: number;
}

export interface ResourcesResponse {
  food: number;
  water: number;
  lumber: number;
  stone: number;
  food_rate: number;
  water_rate: number;
  lumber_rate: number;
  stone_rate: number;
  food_consumption: number;
  max_food: number;
  max_water: number;
  max_lumber: number;
  max_stone: number;
  pop_cap: number;
  pop_used: number;
}

// Building upgrade types — mirrors server/internal/dto/building.go

export interface StartUpgradeRequest {
  building_type: string;
}

export interface BuildingQueueResponse {
  id: number;
  building_type: BuildingType;
  target_level: number;
  started_at: string;
  completes_at: string;
}

export interface BuildingCostResponse {
  building_type: string;
  current_level: number;
  target_level: number;
  food: number;
  water: number;
  lumber: number;
  stone: number;
  time_seconds: number;
}

// Training types — mirrors server/internal/dto/training.go

export interface StartTrainingRequest {
  troop_type: string;
  quantity: number;
}

export interface TrainingQueueResponse {
  id: number;
  troop_type: TroopType;
  quantity: number;
  original_quantity: number;
  each_duration_sec: number;
  started_at: string;
  completes_at: string;
}

export interface TrainingCostResponse {
  troop_type: string;
  quantity: number;
  total_food: number;
  total_water: number;
  total_lumber: number;
  total_stone: number;
  each_time_sec: number;
  total_time_sec: number;
}

export interface TroopInfo {
  type: TroopType;
  quantity: number;
  status: string;
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
  kingdom: Kingdom | '';
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

export type AssetCategory = 'building' | 'resource' | 'unit' | 'kingdom_flag' | 'village_marker' | 'zone_tile' | 'terrain_tile';

export interface GameAsset {
  id: string;
  category: AssetCategory;
  display_name: string;
  default_icon: string;
  sprite_url?: string | null;
  updated_at: string;
}

export interface GameAssetListResponse {
  assets: GameAsset[];
}

// Season types — mirrors server/internal/dto/season.go

export interface SeasonResponse {
  id: number;
  name: string;
  description: string;
  status: 'upcoming' | 'active' | 'ended' | 'archived';
  start_date?: string;
  started_at?: string;
  ended_at?: string;
  player_count: number;
  map_template_name: string;
  game_speed: number;
  resource_multiplier: number;
  max_villages_per_player: number;
  weapons_of_chaos_count: number;
  map_width: number;
  map_height: number;
  created_at: string;
  updated_at: string;
}

export interface SeasonDetailResponse extends SeasonResponse {
  joined: boolean;
  kingdom?: string;
}

export interface CreateSeasonRequest {
  name: string;
  description?: string;
  start_date?: string;
  map_template_name?: string;
  game_speed?: number;
  resource_multiplier?: number;
  max_villages_per_player?: number;
  weapons_of_chaos_count?: number;
  map_width: number;
  map_height: number;
}

export interface UpdateSeasonRequest {
  name?: string;
  description?: string;
  start_date?: string;
  map_template_name?: string;
  game_speed?: number;
  resource_multiplier?: number;
  max_villages_per_player?: number;
  weapons_of_chaos_count?: number;
  map_width?: number;
  map_height?: number;
}

export interface JoinSeasonRequest {
  kingdom: string;
}

export interface JoinSeasonResponse {
  season: SeasonDetailResponse;
  village_id: number;
}

export interface PlayerProfileResponse {
  id: number;
  username: string;
  email: string;
  role: 'player' | 'admin';
  created_at: string;
  total_seasons: number;
  season_history: SeasonHistoryEntry[];
}

export interface SeasonHistoryEntry {
  season_id: number;
  season_name: string;
  season_status: string;
  kingdom: string;
  joined_at: string;
  village_count: number;
}

// Building display config types — mirrors server/internal/dto/admin.go

export interface BuildingDisplayConfig {
  id: number;
  building_type: string;
  kingdom: string;
  display_name: string;
  description: string;
  default_icon: string;
  updated_at: string;
}

export interface BuildingDisplayConfigListResponse {
  configs: BuildingDisplayConfig[];
}

// Troop display config types — mirrors server/internal/dto/admin.go

export interface TroopDisplayConfig {
  id: number;
  troop_type: string;
  kingdom: string;
  training_building: string;
  display_name: string;
  description: string;
  default_icon: string;
  updated_at: string;
}

export interface TroopDisplayConfigListResponse {
  configs: TroopDisplayConfig[];
}

// Resource building config types — mirrors server/internal/dto/admin.go

export interface ResourceBuildingConfig {
  id: number;
  resource_type: string;
  slot: number;
  kingdom: string;
  display_name: string;
  description: string;
  default_icon: string;
  updated_at: string;
}

export interface ResourceBuildingConfigListResponse {
  configs: ResourceBuildingConfig[];
}

// Building sprite info — returned by GET /api/admin/sprites/buildings/{kingdom}

export interface BuildingSpriteInfo {
  filename: string;
  resource_type: string;
  slot: number;
  name: string;
  url: string;
}

export interface BuildingSpriteListResponse {
  sprites: BuildingSpriteInfo[];
}

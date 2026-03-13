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
  gold: number;
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
  gold: number;
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
  total_gold: number;
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

// Display building sprite info — returned by GET /api/admin/sprites/display-buildings/{kingdom}

export interface DisplayBuildingSpriteInfo {
  filename: string;
  building_type: string;
  name: string;
  url: string;
}

export interface DisplayBuildingSpriteListResponse {
  sprites: DisplayBuildingSpriteInfo[];
}

// Troop sprite info — returned by GET /api/admin/sprites/troops/{kingdom}

export interface TroopSpriteInfo {
  filename: string;
  troop_type: string;
  name: string;
  url: string;
}

export interface TroopSpriteListResponse {
  sprites: TroopSpriteInfo[];
}

// ── Camp & Expedition types — mirrors server/internal/dto/camp.go ──────────

/** A camp visible on the world map */
export interface CampResponse {
  id: number;
  template_name: string;
  tier: number;
  tile_x: number;
  tile_y: number;
  status: 'active' | 'under_attack' | 'cleared';
  spawned_at: string;
  beasts: CampBeastResponse[];
}

export interface CampBeastResponse {
  name: string;
  sprite_key: string;
  hp: number;
  attack_power: number;
  attack_interval: number;
  defense_percent: number;
  crit_chance_percent: number;
  count: number;
}

/** Dispatch expedition request */
export interface DispatchExpeditionRequest {
  camp_id: number;
  troops: TroopDispatch[];
}

export interface TroopDispatch {
  troop_type: string;
  quantity: number;
}

/** Expedition status response */
export interface ExpeditionResponse {
  id: number;
  village_id: number;
  camp_id: number;
  status: 'marching' | 'battling' | 'returning' | 'completed';
  dispatched_at: string;
  arrives_at: string;
  return_at?: string;
  completed_at?: string;
  battle_id?: number;
  troops: ExpeditionTroopResponse[];
}

export interface ExpeditionTroopResponse {
  troop_type: string;
  quantity_sent: number;
  quantity_survived: number;
}

/** Battle report */
export interface BattleReportResponse {
  id: number;
  expedition_id: number;
  camp_id: number;
  result: 'attacker_won' | 'defender_won' | 'draw';
  attacker_losses: BattleLosses;
  defender_losses: BattleLosses;
  rewards: BattleRewardResponse[];
  fought_at: string;
}

export interface BattleLosses {
  total_sent: number;
  total_lost: number;
  total_survived: number;
}

export interface BattleRewardResponse {
  resource_type: string;
  amount: number;
}

/** Raw replay JSON */
export interface BattleReplayResponse {
  version: number;
  tick_rate_ms: number;
  attackers: ReplayUnit[];
  defenders: ReplayUnit[];
  events: ReplayEvent[];
  result: string;
  total_ticks: number;
}

export interface ReplayUnit {
  id: number;
  side: 'attacker' | 'defender';
  name: string;
  sprite_key: string;
  hp: number;
  max_hp: number;
  attack_power: number;
  attack_interval: number;
  defense_percent: number;
  crit_chance_percent: number;
}

export interface ReplayEvent {
  tick: number;
  type: 'attack' | 'kill';
  source_id: number;
  target_id: number;
  damage: number;
  is_crit: boolean;
  target_hp_after: number;
  is_kill: boolean;
}

// ── Admin Camp types — mirrors server/internal/dto/camp.go admin DTOs ──────

export interface BeastTemplateResponse {
  id: number;
  name: string;
  sprite_key: string;
  hp: number;
  attack_power: number;
  attack_interval: number;
  defense_percent: number;
  crit_chance_percent: number;
  description: string;
  created_at: string;
  updated_at: string;
}

export interface CreateBeastTemplateRequest {
  name: string;
  sprite_key: string;
  hp: number;
  attack_power: number;
  attack_interval?: number;
  defense_percent?: number;
  crit_chance_percent?: number;
  description?: string;
}

export interface UpdateBeastTemplateRequest {
  name?: string;
  sprite_key?: string;
  hp?: number;
  attack_power?: number;
  attack_interval?: number;
  defense_percent?: number;
  crit_chance_percent?: number;
  description?: string;
}

export interface CampTemplateResponse {
  id: number;
  name: string;
  tier: number;
  min_beasts: number;
  max_beasts: number;
  reward_table_id: number;
  description: string;
  created_at: string;
  updated_at: string;
  beast_slots: CampBeastSlotResponse[];
}

export interface CampBeastSlotResponse {
  id: number;
  beast_template_id: number;
  beast_name: string;
  min_count: number;
  max_count: number;
  weight: number;
}

export interface CreateCampTemplateRequest {
  name: string;
  tier: number;
  min_beasts: number;
  max_beasts: number;
  reward_table_id: number;
  description?: string;
  beast_slots: CampBeastSlotRequest[];
}

export interface CampBeastSlotRequest {
  beast_template_id: number;
  min_count: number;
  max_count: number;
  weight: number;
}

export interface UpdateCampTemplateRequest {
  name?: string;
  tier?: number;
  min_beasts?: number;
  max_beasts?: number;
  reward_table_id?: number;
  description?: string;
}

export interface SpawnRuleResponse {
  id: number;
  name: string;
  enabled: boolean;
  terrain_types: string[];
  zone_types: string[];
  camp_template_pool: CampTemplatePoolEntry[];
  max_camps: number;
  spawn_interval_sec: number;
  despawn_after_sec: number;
  min_camp_distance: number;
  min_village_distance: number;
  created_at: string;
  updated_at: string;
}

export interface CampTemplatePoolEntry {
  camp_template_id: number;
  weight: number;
}

export interface CreateSpawnRuleRequest {
  name: string;
  terrain_types: string[];
  zone_types: string[];
  camp_template_pool: CampTemplatePoolEntry[];
  max_camps: number;
  spawn_interval_sec: number;
  despawn_after_sec: number;
  min_camp_distance?: number;
  min_village_distance?: number;
}

export interface UpdateSpawnRuleRequest {
  name?: string;
  enabled?: boolean;
  terrain_types?: string[];
  zone_types?: string[];
  camp_template_pool?: CampTemplatePoolEntry[];
  max_camps?: number;
  spawn_interval_sec?: number;
  despawn_after_sec?: number;
  min_camp_distance?: number;
  min_village_distance?: number;
}

export interface RewardTableResponse {
  id: number;
  name: string;
  created_at: string;
  updated_at: string;
  entries: RewardEntryResponse[];
}

export interface RewardEntryResponse {
  id: number;
  reward_type: string;
  min_amount: number;
  max_amount: number;
  drop_chance_pct: number;
}

export interface CreateRewardTableRequest {
  name: string;
  entries: RewardEntryRequest[];
}

export interface UpdateRewardTableRequest {
  name?: string;
}

export interface RewardEntryRequest {
  reward_type: string;
  min_amount: number;
  max_amount: number;
  drop_chance_pct: number;
}

export interface BattleTuningResponse {
  tick_duration_ms: number;
  crit_damage_multiplier: number;
  max_defense_percent: number;
  max_crit_chance_percent: number;
  min_attack_interval: number;
  march_speed_tiles_per_min: number;
  max_ticks: number;
}

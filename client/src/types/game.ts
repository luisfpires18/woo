// Game entity types — mirrors server domain models

export type Kingdom = 'veridor' | 'sylvara' | 'arkazia' | 'draxys' | 'zandres' | 'lumus' | 'nordalh' | 'drakanith';

/** Kingdoms that players can currently select during registration. */
export type PlayableKingdom = 'veridor' | 'sylvara' | 'arkazia' | 'draxys' | 'nordalh' | 'zandres' | 'lumus';

export interface Player {
  id: number;
  username: string;
  email: string;
  kingdom: Kingdom;
  created_at: string;
  last_login_at?: string;
}

export interface Village {
  id: number;
  player_id: number;
  name: string;
  x: number;
  y: number;
  is_capital: boolean;
  created_at: string;
}

export interface Building {
  id: number;
  village_id: number;
  building_type: BuildingType;
  level: number;
}

export type BuildingType =
  | 'town_hall'
  | 'food_1'
  | 'food_2'
  | 'food_3'
  | 'water_1'
  | 'water_2'
  | 'water_3'
  | 'lumber_1'
  | 'lumber_2'
  | 'lumber_3'
  | 'stone_1'
  | 'stone_2'
  | 'stone_3'
  | 'barracks'
  | 'stable'
  | 'archery'
  | 'workshop'
  | 'special'
  | 'storage'
  | 'provisions'
  | 'reservoir';

export interface Resources {
  village_id: number;
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
  last_updated: string;
}

// Game entity types — mirrors server domain models

export type Kingdom = 'veridor' | 'sylvara' | 'arkazia' | 'draxys' | 'zandres' | 'lumus' | 'nordalh' | 'drakanith';

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
  | 'iron_mine'
  | 'lumber_mill'
  | 'quarry'
  | 'farm'
  | 'warehouse'
  | 'barracks'
  | 'stable'
  | 'forge'
  | 'rune_altar'
  | 'walls'
  | 'marketplace'
  | 'embassy'
  | 'watchtower'
  | 'dock'
  | 'grove_sanctum'
  | 'colosseum';

export interface Resources {
  village_id: number;
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
  last_updated: string;
}

// Map-related type definitions — mirrors server DTOs

/** Terrain types matching the server enum */
export type TerrainType =
  | 'plains'
  | 'forest'
  | 'mountain'
  | 'water'
  | 'desert'
  | 'swamp'
  | 'chasm'
  | 'bridge';

/** Kingdom zone identifiers */
export type KingdomZone =
  | ''
  | 'moraphys'
  | 'veridor'
  | 'sylvara'
  | 'arkazia'
  | 'draxys'
  | 'zandres'
  | 'lumus'
  | 'nordalh'
  | 'drakanith'
  | 'dark_reach'
  | 'wilderness';

/** A single map tile as returned by the API */
export interface MapTile {
  x: number;
  y: number;
  terrain: TerrainType;
  zone: KingdomZone;
  village_id?: number;
  village_name?: string;
  owner_name?: string;
}

/** Response from GET /api/map */
export interface MapChunkResponse {
  center_x: number;
  center_y: number;
  range: number;
  tiles: MapTile[];
}

/** Terrain rendering config — color, movement modifier, label */
export interface TerrainConfig {
  color: number;
  label: string;
  movementMod: number;
  passable: boolean;
}

/** Map of terrain types to their rendering config */
export const TERRAIN_CONFIG: Record<TerrainType, TerrainConfig> = {
  plains:   { color: 0x7ec850, label: 'Plains',   movementMod: 1.0, passable: true },
  forest:   { color: 0x2d7a3a, label: 'Forest',   movementMod: 0.8, passable: true },
  mountain: { color: 0x8b7355, label: 'Mountain', movementMod: 0.6, passable: true },
  water:    { color: 0x3a7ec8, label: 'Water',    movementMod: 0,   passable: false },
  desert:   { color: 0xd4a843, label: 'Desert',   movementMod: 0.7, passable: true },
  swamp:    { color: 0x5a6e3a, label: 'Swamp',    movementMod: 0.5, passable: true },
  chasm:    { color: 0x1a0a2e, label: 'Chasm',    movementMod: 0,   passable: false },
  bridge:   { color: 0x8b6914, label: 'Bridge',   movementMod: 0.9, passable: true },
};

/** Kingdom zone colors for zone overlay tinting */
export const ZONE_COLORS: Record<string, number> = {
  moraphys:   0x330000,
  veridor:    0x004488,
  sylvara:    0x006622,
  arkazia:    0x664400,
  draxys:     0x880000,
  zandres:    0x886600,
  lumus:      0x446688,
  nordalh:    0x225544,
  drakanith:  0x662200,
  dark_reach: 0x110022,
  wilderness: 0x444444,
};

/** All paintable kingdom zones */
export const KINGDOM_ZONES: string[] = [
  'wilderness',
  'moraphys',
  'veridor',
  'sylvara',
  'arkazia',
  'draxys',
  'zandres',
  'lumus',
  'nordalh',
  'drakanith',
  'dark_reach',
];

/** Zone rendering config — color + label */
export const ZONE_CONFIG: Record<string, { color: number; label: string }> = {
  moraphys:   { color: 0x330000, label: 'Moraphys' },
  veridor:    { color: 0x004488, label: 'Veridor' },
  sylvara:    { color: 0x006622, label: 'Sylvara' },
  arkazia:    { color: 0x664400, label: 'Arkazia' },
  draxys:     { color: 0x880000, label: 'Draxys' },
  zandres:    { color: 0x886600, label: 'Zandres' },
  lumus:      { color: 0x446688, label: 'Lumus' },
  nordalh:    { color: 0x225544, label: 'Nordalh' },
  drakanith:  { color: 0x662200, label: 'Drakanith' },
  dark_reach: { color: 0x110022, label: 'Dark Reach' },
  wilderness: { color: 0x444444, label: 'Wilderness' },
};

// --- Map template types ---

/** Metadata for a map template (no tiles) */
export interface TemplateInfo {
  name: string;
  description: string;
  map_size: number;
  tile_count: number;
  created_at: string;
  updated_at: string;
}

/** A single tile in a map template */
export interface TemplateTile {
  x: number;
  y: number;
  terrain_type: string;
  kingdom_zone: string;
}

/** Full map template with tiles */
export interface MapTemplate {
  name: string;
  description: string;
  map_size: number;
  tiles: TemplateTile[];
  created_at: string;
  updated_at: string;
}

-- 002_seed_data.sql — Seed data for development.
-- World config defaults, admin accounts, game assets, resource building configs.

-- ── World Config defaults ────────────────────────────────────────────────────
INSERT INTO world_config (key, value, description) VALUES
    ('game_speed',              '1.0',  'Global game speed multiplier'),
    ('resource_multiplier',     '1.0',  'Resource production multiplier'),
    ('map_width',               '200',  'World map width in tiles'),
    ('map_height',              '200',  'World map height in tiles'),
    ('weapons_of_chaos_count',  '7',    'Number of Weapons of Chaos per world'),
    ('max_villages_per_player', '5',    'Maximum villages a player can own');

-- ── Admin seed accounts ──────────────────────────────────────────────────────
-- Password for all: Woo123!
-- bcrypt hash: $2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC
-- Kingdom left empty ('') → kingdom picker shown on first login.
INSERT OR IGNORE INTO players (username, email, password_hash, kingdom, role, created_at)
VALUES
    ('wright',   'wright@woo.local',   '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('coyote',   'coyote@woo.local',   '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('veridor',  'veridor@woo.local',  '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('sylvaine', 'sylvaine@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('taldros',  'taldros@woo.local',  '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('blake',    'blake@woo.local',    '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('hawkes',   'hawkes@woo.local',   '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('drakt',    'drakt@woo.local',    '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now'));

-- ── Game Assets: Village buildings (13 non-resource) ─────────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('town_hall',     'building', 'Town Hall',     '🏛️',  96, 96),
    ('warehouse',     'building', 'Warehouse',     '📦',  96, 96),
    ('barracks',      'building', 'Barracks',      '⚔️',  96, 96),
    ('stable',        'building', 'Stable',        '🐴',  96, 96),
    ('forge',         'building', 'Forge',         '🔨',  96, 96),
    ('rune_altar',    'building', 'Rune Altar',    '🔮',  96, 96),
    ('walls',         'building', 'Walls',         '🏰',  96, 96),
    ('marketplace',   'building', 'Marketplace',   '🏪',  96, 96),
    ('embassy',       'building', 'Embassy',       '📜',  96, 96),
    ('watchtower',    'building', 'Watchtower',    '👁️',  96, 96),
    ('dock',          'building', 'Dock',          '⚓',  96, 96),
    ('grove_sanctum', 'building', 'Grove Sanctum', '🌿',  96, 96),
    ('colosseum',     'building', 'Colosseum',     '🏟️',  96, 96);

-- ── Game Assets: Resource buildings (12 = 3 per resource) ────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('food_1',   'building', 'Food Field I',    '🌾', 96, 96),
    ('food_2',   'building', 'Food Field II',   '🐟', 96, 96),
    ('food_3',   'building', 'Food Field III',  '🍎', 96, 96),
    ('water_1',  'building', 'Water Field I',   '💧', 96, 96),
    ('water_2',  'building', 'Water Field II',  '🏞️', 96, 96),
    ('water_3',  'building', 'Water Field III', '🚰', 96, 96),
    ('lumber_1', 'building', 'Lumber Field I',  '🪓', 96, 96),
    ('lumber_2', 'building', 'Lumber Field II', '🪵', 96, 96),
    ('lumber_3', 'building', 'Lumber Field III','🌲', 96, 96),
    ('stone_1',  'building', 'Stone Field I',   '⛏️', 96, 96),
    ('stone_2',  'building', 'Stone Field II',  '🪨', 96, 96),
    ('stone_3',  'building', 'Stone Field III', '⛰️', 96, 96);

-- ── Game Assets: Resources (4) ───────────────────────────────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('food',   'resource', 'Food',   '🌾', 32, 32),
    ('water',  'resource', 'Water',  '💧', 32, 32),
    ('lumber', 'resource', 'Lumber', '🪵', 32, 32),
    ('stone',  'resource', 'Stone',  '🪨', 32, 32);

-- ── Game Assets: Kingdom flags (8) ───────────────────────────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('flag_veridor',   'kingdom_flag', 'Veridor Flag',   '🔵', 256, 256),
    ('flag_sylvara',   'kingdom_flag', 'Sylvara Flag',   '🟢', 256, 256),
    ('flag_arkazia',   'kingdom_flag', 'Arkazia Flag',   '🔴', 256, 256),
    ('flag_draxys',    'kingdom_flag', 'Draxys Flag',    '🟡', 256, 256),
    ('flag_nordalh',   'kingdom_flag', 'Nordalh Flag',   '🟣', 256, 256),
    ('flag_zandres',   'kingdom_flag', 'Zandres Flag',   '🟤', 256, 256),
    ('flag_lumus',     'kingdom_flag', 'Lumus Flag',     '⚪', 256, 256),
    ('flag_drakanith', 'kingdom_flag', 'Drakanith Flag', '🟠', 256, 256);

-- ── Game Assets: Village markers (8, one per kingdom) ────────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('marker_veridor',   'village_marker', 'Veridor Village',   '🏘️', 256, 256),
    ('marker_sylvara',   'village_marker', 'Sylvara Village',   '🏘️', 256, 256),
    ('marker_arkazia',   'village_marker', 'Arkazia Village',   '🏘️', 256, 256),
    ('marker_draxys',    'village_marker', 'Draxys Village',    '🏘️', 256, 256),
    ('marker_nordalh',   'village_marker', 'Nordalh Village',   '🏘️', 256, 256),
    ('marker_zandres',   'village_marker', 'Zandres Village',   '🏘️', 256, 256),
    ('marker_lumus',     'village_marker', 'Lumus Village',     '🏘️', 256, 256),
    ('marker_drakanith', 'village_marker', 'Drakanith Village', '🏘️', 256, 256);

-- ── Game Assets: Zone tiles (8, one per kingdom + 1 default) ─────────────────
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('zone_default',   'zone_tile', 'Default Zone',     '🟩', 256, 256),
    ('zone_veridor',   'zone_tile', 'Veridor Zone',     '🔵', 256, 256),
    ('zone_sylvara',   'zone_tile', 'Sylvara Zone',     '🟢', 256, 256),
    ('zone_arkazia',   'zone_tile', 'Arkazia Zone',     '🔴', 256, 256),
    ('zone_draxys',    'zone_tile', 'Draxys Zone',      '🟡', 256, 256),
    ('zone_nordalh',   'zone_tile', 'Nordalh Zone',     '🟣', 256, 256),
    ('zone_zandres',   'zone_tile', 'Zandres Zone',     '🟤', 256, 256),
    ('zone_lumus',     'zone_tile', 'Lumus Zone',       '⚪', 256, 256),
    ('zone_drakanith', 'zone_tile', 'Drakanith Zone',   '🟠', 256, 256);

-- ── Resource Building Configs (96 rows: 8 kingdoms × 4 resources × 3 slots) ─

-- Food
INSERT INTO resource_building_configs (resource_type, slot, kingdom, display_name, description, default_icon) VALUES
    ('food', 1, 'veridor',   'Farm',    'Grows crops for food.',  '🌾'),
    ('food', 2, 'veridor',   'Fishery', 'Catches fish for food.', '🐟'),
    ('food', 3, 'veridor',   'Orchard', 'Harvests fruit.',        '🍎'),
    ('food', 1, 'sylvara',   'Farm',    'Grows crops for food.',  '🌾'),
    ('food', 2, 'sylvara',   'Fishery', 'Catches fish for food.', '🐟'),
    ('food', 3, 'sylvara',   'Orchard', 'Harvests fruit.',        '🍎'),
    ('food', 1, 'arkazia',   'Farm',    'Grows crops for food.',  '🌾'),
    ('food', 2, 'arkazia',   'Fishery', 'Catches fish for food.', '🐟'),
    ('food', 3, 'arkazia',   'Orchard', 'Harvests fruit.',        '🍎'),
    ('food', 1, 'draxys',    'Farm',    'Grows crops for food.',  '🌾'),
    ('food', 2, 'draxys',    'Fishery', 'Catches fish for food.', '🐟'),
    ('food', 3, 'draxys',    'Orchard', 'Harvests fruit.',        '🍎'),
    ('food', 1, 'zandres',   'Farm',    'Grows crops for food.',  '🌾'),
    ('food', 2, 'zandres',   'Fishery', 'Catches fish for food.', '🐟'),
    ('food', 3, 'zandres',   'Orchard', 'Harvests fruit.',        '🍎'),
    ('food', 1, 'lumus',     'Farm',    'Grows crops for food.',  '🌾'),
    ('food', 2, 'lumus',     'Fishery', 'Catches fish for food.', '🐟'),
    ('food', 3, 'lumus',     'Orchard', 'Harvests fruit.',        '🍎'),
    ('food', 1, 'nordalh',   'Farm',    'Grows crops for food.',  '🌾'),
    ('food', 2, 'nordalh',   'Fishery', 'Catches fish for food.', '🐟'),
    ('food', 3, 'nordalh',   'Orchard', 'Harvests fruit.',        '🍎'),
    ('food', 1, 'drakanith', 'Farm',    'Grows crops for food.',  '🌾'),
    ('food', 2, 'drakanith', 'Fishery', 'Catches fish for food.', '🐟'),
    ('food', 3, 'drakanith', 'Orchard', 'Harvests fruit.',        '🍎');

-- Water
INSERT INTO resource_building_configs (resource_type, slot, kingdom, display_name, description, default_icon) VALUES
    ('water', 1, 'veridor',   'Well',     'Draws water from underground.',  '💧'),
    ('water', 2, 'veridor',   'Spring',   'Collects natural spring water.', '🏞️'),
    ('water', 3, 'veridor',   'Aqueduct', 'Channels water from afar.',     '🚰'),
    ('water', 1, 'sylvara',   'Well',     'Draws water from underground.',  '💧'),
    ('water', 2, 'sylvara',   'Spring',   'Collects natural spring water.', '🏞️'),
    ('water', 3, 'sylvara',   'Aqueduct', 'Channels water from afar.',     '🚰'),
    ('water', 1, 'arkazia',   'Well',     'Draws water from underground.',  '💧'),
    ('water', 2, 'arkazia',   'Spring',   'Collects natural spring water.', '🏞️'),
    ('water', 3, 'arkazia',   'Aqueduct', 'Channels water from afar.',     '🚰'),
    ('water', 1, 'draxys',    'Well',     'Draws water from underground.',  '💧'),
    ('water', 2, 'draxys',    'Spring',   'Collects natural spring water.', '🏞️'),
    ('water', 3, 'draxys',    'Aqueduct', 'Channels water from afar.',     '🚰'),
    ('water', 1, 'zandres',   'Well',     'Draws water from underground.',  '💧'),
    ('water', 2, 'zandres',   'Spring',   'Collects natural spring water.', '🏞️'),
    ('water', 3, 'zandres',   'Aqueduct', 'Channels water from afar.',     '🚰'),
    ('water', 1, 'lumus',     'Well',     'Draws water from underground.',  '💧'),
    ('water', 2, 'lumus',     'Spring',   'Collects natural spring water.', '🏞️'),
    ('water', 3, 'lumus',     'Aqueduct', 'Channels water from afar.',     '🚰'),
    ('water', 1, 'nordalh',   'Well',     'Draws water from underground.',  '💧'),
    ('water', 2, 'nordalh',   'Spring',   'Collects natural spring water.', '🏞️'),
    ('water', 3, 'nordalh',   'Aqueduct', 'Channels water from afar.',     '🚰'),
    ('water', 1, 'drakanith', 'Well',     'Draws water from underground.',  '💧'),
    ('water', 2, 'drakanith', 'Spring',   'Collects natural spring water.', '🏞️'),
    ('water', 3, 'drakanith', 'Aqueduct', 'Channels water from afar.',     '🚰');

-- Lumber
INSERT INTO resource_building_configs (resource_type, slot, kingdom, display_name, description, default_icon) VALUES
    ('lumber', 1, 'veridor',   'Sawmill',     'Cuts timber into lumber.',  '🪓'),
    ('lumber', 2, 'veridor',   'Lumber Camp', 'Harvests trees in bulk.',   '🪵'),
    ('lumber', 3, 'veridor',   'Woodcutter',  'Fells trees for wood.',     '🌲'),
    ('lumber', 1, 'sylvara',   'Sawmill',     'Cuts timber into lumber.',  '🪓'),
    ('lumber', 2, 'sylvara',   'Lumber Camp', 'Harvests trees in bulk.',   '🪵'),
    ('lumber', 3, 'sylvara',   'Woodcutter',  'Fells trees for wood.',     '🌲'),
    ('lumber', 1, 'arkazia',   'Sawmill',     'Cuts timber into lumber.',  '🪓'),
    ('lumber', 2, 'arkazia',   'Lumber Camp', 'Harvests trees in bulk.',   '🪵'),
    ('lumber', 3, 'arkazia',   'Woodcutter',  'Fells trees for wood.',     '🌲'),
    ('lumber', 1, 'draxys',    'Sawmill',     'Cuts timber into lumber.',  '🪓'),
    ('lumber', 2, 'draxys',    'Lumber Camp', 'Harvests trees in bulk.',   '🪵'),
    ('lumber', 3, 'draxys',    'Woodcutter',  'Fells trees for wood.',     '🌲'),
    ('lumber', 1, 'zandres',   'Sawmill',     'Cuts timber into lumber.',  '🪓'),
    ('lumber', 2, 'zandres',   'Lumber Camp', 'Harvests trees in bulk.',   '🪵'),
    ('lumber', 3, 'zandres',   'Woodcutter',  'Fells trees for wood.',     '🌲'),
    ('lumber', 1, 'lumus',     'Sawmill',     'Cuts timber into lumber.',  '🪓'),
    ('lumber', 2, 'lumus',     'Lumber Camp', 'Harvests trees in bulk.',   '🪵'),
    ('lumber', 3, 'lumus',     'Woodcutter',  'Fells trees for wood.',     '🌲'),
    ('lumber', 1, 'nordalh',   'Sawmill',     'Cuts timber into lumber.',  '🪓'),
    ('lumber', 2, 'nordalh',   'Lumber Camp', 'Harvests trees in bulk.',   '🪵'),
    ('lumber', 3, 'nordalh',   'Woodcutter',  'Fells trees for wood.',     '🌲'),
    ('lumber', 1, 'drakanith', 'Sawmill',     'Cuts timber into lumber.',  '🪓'),
    ('lumber', 2, 'drakanith', 'Lumber Camp', 'Harvests trees in bulk.',   '🪵'),
    ('lumber', 3, 'drakanith', 'Woodcutter',  'Fells trees for wood.',     '🌲');

-- Stone
INSERT INTO resource_building_configs (resource_type, slot, kingdom, display_name, description, default_icon) VALUES
    ('stone', 1, 'veridor',   'Quarry',    'Extracts stone from rock.',  '⛏️'),
    ('stone', 2, 'veridor',   'Stone Pit', 'Digs stone from the earth.', '🪨'),
    ('stone', 3, 'veridor',   'Mine',      'Mines deep for stone.',      '⛰️'),
    ('stone', 1, 'sylvara',   'Quarry',    'Extracts stone from rock.',  '⛏️'),
    ('stone', 2, 'sylvara',   'Stone Pit', 'Digs stone from the earth.', '🪨'),
    ('stone', 3, 'sylvara',   'Mine',      'Mines deep for stone.',      '⛰️'),
    ('stone', 1, 'arkazia',   'Quarry',    'Extracts stone from rock.',  '⛏️'),
    ('stone', 2, 'arkazia',   'Stone Pit', 'Digs stone from the earth.', '🪨'),
    ('stone', 3, 'arkazia',   'Mine',      'Mines deep for stone.',      '⛰️'),
    ('stone', 1, 'draxys',    'Quarry',    'Extracts stone from rock.',  '⛏️'),
    ('stone', 2, 'draxys',    'Stone Pit', 'Digs stone from the earth.', '🪨'),
    ('stone', 3, 'draxys',    'Mine',      'Mines deep for stone.',      '⛰️'),
    ('stone', 1, 'zandres',   'Quarry',    'Extracts stone from rock.',  '⛏️'),
    ('stone', 2, 'zandres',   'Stone Pit', 'Digs stone from the earth.', '🪨'),
    ('stone', 3, 'zandres',   'Mine',      'Mines deep for stone.',      '⛰️'),
    ('stone', 1, 'lumus',     'Quarry',    'Extracts stone from rock.',  '⛏️'),
    ('stone', 2, 'lumus',     'Stone Pit', 'Digs stone from the earth.', '🪨'),
    ('stone', 3, 'lumus',     'Mine',      'Mines deep for stone.',      '⛰️'),
    ('stone', 1, 'nordalh',   'Quarry',    'Extracts stone from rock.',  '⛏️'),
    ('stone', 2, 'nordalh',   'Stone Pit', 'Digs stone from the earth.', '🪨'),
    ('stone', 3, 'nordalh',   'Mine',      'Mines deep for stone.',      '⛰️'),
    ('stone', 1, 'drakanith', 'Quarry',    'Extracts stone from rock.',  '⛏️'),
    ('stone', 2, 'drakanith', 'Stone Pit', 'Digs stone from the earth.', '🪨'),
    ('stone', 3, 'drakanith', 'Mine',      'Mines deep for stone.',      '⛰️');

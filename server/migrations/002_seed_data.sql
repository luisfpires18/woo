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
-- Exception: wright has kingdom pre-set to 'arkazia' for dev testing.
INSERT OR IGNORE INTO players (username, email, password_hash, kingdom, role, created_at)
VALUES
    ('wright',   'wright@woo.local',   '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', 'arkazia', 'admin', datetime('now')),
    ('coyote',   'coyote@woo.local',   '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('veridor',  'veridor@woo.local',  '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('sylvaine', 'sylvaine@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('taldros',  'taldros@woo.local',  '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('blake',    'blake@woo.local',    '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('hawkes',   'hawkes@woo.local',   '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
    ('drakt',    'drakt@woo.local',    '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now'));

-- ── Game Assets: Village buildings (6 non-resource) ──────────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon) VALUES
    ('town_hall',  'building', 'Town Hall',  '🏛️'),
    ('barracks',   'building', 'Barracks',   '⚔️'),
    ('stable',     'building', 'Stable',     '🐴'),
    ('archery',    'building', 'Archery',    '🏹'),
    ('workshop',   'building', 'Workshop',   '🔨'),
    ('special',    'building', 'Special',    '⭐');

-- ── Game Assets: Storage buildings (3 = capacity boosters) ───────────────────
INSERT INTO game_assets (id, category, display_name, default_icon) VALUES
    ('storage',    'building', 'Storage',    '🏚️'),
    ('provisions', 'building', 'Provisions', '🍞'),
    ('reservoir',  'building', 'Reservoir',  '🏊');

-- ── Game Assets: Resource buildings (12 = 3 per resource) ────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon) VALUES
    ('food_1',   'building', 'Food Field I',    '🌾'),
    ('food_2',   'building', 'Food Field II',   '🐟'),
    ('food_3',   'building', 'Food Field III',  '🍎'),
    ('water_1',  'building', 'Water Field I',   '💧'),
    ('water_2',  'building', 'Water Field II',  '🏞️'),
    ('water_3',  'building', 'Water Field III', '🚰'),
    ('lumber_1', 'building', 'Lumber Field I',  '🪓'),
    ('lumber_2', 'building', 'Lumber Field II', '🪵'),
    ('lumber_3', 'building', 'Lumber Field III','🌲'),
    ('stone_1',  'building', 'Stone Field I',   '⛏️'),
    ('stone_2',  'building', 'Stone Field II',  '🪨'),
    ('stone_3',  'building', 'Stone Field III', '⛰️');

-- ── Game Assets: Resources (4) ───────────────────────────────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon) VALUES
    ('food',   'resource', 'Food',   '🌾'),
    ('water',  'resource', 'Water',  '💧'),
    ('lumber', 'resource', 'Lumber', '🪵'),
    ('stone',  'resource', 'Stone',  '🪨');

-- ── Game Assets: Kingdom flags (8) ───────────────────────────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon) VALUES
    ('flag_veridor',   'kingdom_flag', 'Veridor Flag',   '🔵'),
    ('flag_sylvara',   'kingdom_flag', 'Sylvara Flag',   '🟢'),
    ('flag_arkazia',   'kingdom_flag', 'Arkazia Flag',   '🔴'),
    ('flag_draxys',    'kingdom_flag', 'Draxys Flag',    '🟡'),
    ('flag_nordalh',   'kingdom_flag', 'Nordalh Flag',   '🟣'),
    ('flag_zandres',   'kingdom_flag', 'Zandres Flag',   '🟤'),
    ('flag_lumus',     'kingdom_flag', 'Lumus Flag',     '⚪'),
    ('flag_drakanith', 'kingdom_flag', 'Drakanith Flag', '🟠');

-- ── Game Assets: Village markers (8, one per kingdom) ────────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon) VALUES
    ('marker_veridor',   'village_marker', 'Veridor Village',   '🏘️'),
    ('marker_sylvara',   'village_marker', 'Sylvara Village',   '🏘️'),
    ('marker_arkazia',   'village_marker', 'Arkazia Village',   '🏘️'),
    ('marker_draxys',    'village_marker', 'Draxys Village',    '🏘️'),
    ('marker_nordalh',   'village_marker', 'Nordalh Village',   '🏘️'),
    ('marker_zandres',   'village_marker', 'Zandres Village',   '🏘️'),
    ('marker_lumus',     'village_marker', 'Lumus Village',     '🏘️'),
    ('marker_drakanith', 'village_marker', 'Drakanith Village', '🏘️');

-- ── Game Assets: Zone tiles (8, one per kingdom + 1 default + wilderness/moraphys/dark_reach) ─
INSERT INTO game_assets (id, category, display_name, default_icon) VALUES
    ('zone_default',    'zone_tile', 'Default Zone',    '🟩'),
    ('zone_veridor',    'zone_tile', 'Veridor Zone',    '🔵'),
    ('zone_sylvara',    'zone_tile', 'Sylvara Zone',    '🟢'),
    ('zone_arkazia',    'zone_tile', 'Arkazia Zone',    '🔴'),
    ('zone_draxys',     'zone_tile', 'Draxys Zone',     '🟡'),
    ('zone_nordalh',    'zone_tile', 'Nordalh Zone',    '🟣'),
    ('zone_zandres',    'zone_tile', 'Zandres Zone',    '🟤'),
    ('zone_lumus',      'zone_tile', 'Lumus Zone',      '⚪'),
    ('zone_drakanith',  'zone_tile', 'Drakanith Zone',  '🟠'),
    ('zone_wilderness', 'zone_tile', 'Wilderness Zone', '🌿'),
    ('zone_moraphys',   'zone_tile', 'Moraphys Zone',   '💀'),
    ('zone_dark_reach', 'zone_tile', 'Dark Reach Zone', '🖤');

-- ── Game Assets: Terrain tiles (6, one per terrain type) ─────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon) VALUES
    ('terrain_plains',   'terrain_tile', 'Plains',   '🟩'),
    ('terrain_forest',   'terrain_tile', 'Forest',   '🌲'),
    ('terrain_mountain', 'terrain_tile', 'Mountain', '⛰️'),
    ('terrain_water',    'terrain_tile', 'Water',    '🌊'),
    ('terrain_desert',   'terrain_tile', 'Desert',   '🏜️'),
    ('terrain_swamp',    'terrain_tile', 'Swamp',    '🐊');

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

-- ── Season Seeds ───────────────────────────────────────────────────────────────

-- dev_season: Always active, unlimited time — used for day-to-day development.
INSERT OR IGNORE INTO seasons (name, description, status, game_speed, resource_multiplier, max_villages_per_player, weapons_of_chaos_count, map_width, map_height, started_at)
VALUES (
    'dev_season',
    'Permanent development server. No time limit — play at your own pace.',
    'active',
    1.0,
    1.0,
    5,
    7,
    51,
    51,
    datetime('now')
);

-- alpha_season: Upcoming — TBA.
INSERT OR IGNORE INTO seasons (name, description, status, game_speed, resource_multiplier, max_villages_per_player, weapons_of_chaos_count, map_width, map_height)
VALUES (
    'alpha_season',
    'First competitive season. Details to be announced.',
    'upcoming',
    1.0,
    1.0,
    5,
    7,
    101,
    101
);

-- ── Wright dev village seed ────────────────────────────────────────────────
-- Creates one ready-to-play village for wright with all buildings constructed.
INSERT OR IGNORE INTO villages (player_id, name, x, y, is_capital, season_id, created_at)
SELECT 
    p.id,
    'Wright''s Keep',
    0,
    0,
    1,
    s.id,
    datetime('now')
FROM players p
JOIN seasons s ON s.name = 'dev_season'
WHERE p.username = 'wright';

INSERT OR IGNORE INTO buildings (village_id, building_type, level)
SELECT
    v.id,
    b.building_type,
    CASE WHEN b.building_type = 'town_hall' THEN 3 ELSE 1 END
FROM villages v
JOIN (
    SELECT 'town_hall' AS building_type
    UNION ALL SELECT 'food_1'
    UNION ALL SELECT 'food_2'
    UNION ALL SELECT 'food_3'
    UNION ALL SELECT 'water_1'
    UNION ALL SELECT 'water_2'
    UNION ALL SELECT 'water_3'
    UNION ALL SELECT 'lumber_1'
    UNION ALL SELECT 'lumber_2'
    UNION ALL SELECT 'lumber_3'
    UNION ALL SELECT 'stone_1'
    UNION ALL SELECT 'stone_2'
    UNION ALL SELECT 'stone_3'
    UNION ALL SELECT 'barracks'
    UNION ALL SELECT 'stable'
    UNION ALL SELECT 'archery'
    UNION ALL SELECT 'workshop'
    UNION ALL SELECT 'special'
    UNION ALL SELECT 'storage'
    UNION ALL SELECT 'provisions'
    UNION ALL SELECT 'reservoir'
) b
WHERE v.name = 'Wright''s Keep'
  AND v.player_id = (SELECT id FROM players WHERE username = 'wright');

INSERT OR IGNORE INTO resources (
    village_id,
    food,
    water,
    lumber,
    stone,
    food_rate,
    water_rate,
    lumber_rate,
    stone_rate,
    food_consumption,
    max_storage,
    last_updated
)
SELECT
    v.id,
    500,
    500,
    500,
    500,
    3,
    3,
    3,
    3,
    0,
    1200,
    datetime('now')
FROM villages v
WHERE v.name = 'Wright''s Keep'
  AND v.player_id = (SELECT id FROM players WHERE username = 'wright');

-- Create map tile for Wright's village
INSERT OR IGNORE INTO world_map (x, y, terrain_type, kingdom_zone, owner_player_id, village_id, season_id)
SELECT
    0,
    0,
    'plains',
    'arkazia',
    p.id,
    v.id,
    s.id
FROM players p
JOIN villages v ON v.player_id = p.id AND v.name = 'Wright''s Keep'
JOIN seasons s ON s.name = 'dev_season'
WHERE p.username = 'wright';

-- Enroll wright in dev_season with Arkazia kingdom
INSERT OR IGNORE INTO season_players (season_id, player_id, kingdom)
SELECT
    s.id,
    p.id,
    'arkazia'
FROM seasons s
JOIN players p ON p.username = 'wright'
WHERE s.name = 'dev_season';

-- ── Building Display Configs (48 rows: 6 building types × 8 kingdoms) ────────
-- Kingdom-specific display names sourced from docs/01-game-design/kingdoms_units_buildlings.md

INSERT INTO building_display_configs (building_type, kingdom, display_name, description, default_icon) VALUES
    ('town_hall', 'veridor',   'Town Hall',     'The seat of power in your village.',  '🏛️'),
    ('town_hall', 'sylvara',   'Town Hall',     'The seat of power in your village.',  '🏛️'),
    ('town_hall', 'arkazia',   'Town Hall',     'The seat of power in your village.',  '🏛️'),
    ('town_hall', 'draxys',    'Town Hall',     'The seat of power in your village.',  '🏛️'),
    ('town_hall', 'nordalh',   'Town Hall',     'The seat of power in your village.',  '🏛️'),
    ('town_hall', 'zandres',   'Town Hall',     'The seat of power in your village.',  '🏛️'),
    ('town_hall', 'lumus',     'Town Hall',     'The seat of power in your village.',  '🏛️'),
    ('town_hall', 'drakanith', 'Town Hall',     'The seat of power in your village.',  '🏛️'),

    ('barracks', 'veridor',   'Road Barracks',      'Practical infantry hall mixing inland soldiery with disciplined port and trade-route warfare.',  '⚔️'),
    ('barracks', 'sylvara',   'Roothall',           'Infantry lodge where forest warbands train for close fighting, ambush defense, and shielded woodland warfare.',  '⚔️'),
    ('barracks', 'arkazia',   'Red Bastion',        'Infantry stronghold focused on fortress discipline, shield warfare, and mountain-line defense.',  '⚔️'),
    ('barracks', 'draxys',    'Sandwall Barracks',  'Infantry training ground for hard desert warfare, fast brutality, and shielded survival fighting.',  '⚔️'),
    ('barracks', 'nordalh',   'Hearth Barracks',    'Infantry hall where hard northern warriors train for brutal close combat and shielded endurance.',  '⚔️'),
    ('barracks', 'zandres',   'Doorwarden Hall',    'Infantry center built around tunnel defense, pressure control, and precise close-order combat.',  '⚔️'),
    ('barracks', 'lumus',     'Prism Barracks',     'Infantry hall focused on disciplined sacred warfare, radiant order, and defensive grace.',  '⚔️'),
    ('barracks', 'drakanith', 'Barracks',           'Trains infantry units.',  '⚔️'),

    ('stable', 'veridor',   'River Cavalry Yard',  'Cavalry stables built for mobility across roads, marshlands, and river trade territory.',  '🐴'),
    ('stable', 'sylvara',   'Beast Hall',          'Mount-yard where riders bond with woodland beasts instead of standard cavalry horses.',  '🐴'),
    ('stable', 'arkazia',   'Arknight Stables',    'Cavalry yard where disciplined riders train for fortress support and decisive charges.',  '🐴'),
    ('stable', 'draxys',    'Scorpion Pens',       'Beast and mount yard favoring desert creatures over noble cavalry tradition.',  '🐴'),
    ('stable', 'nordalh',   'Wolf Kennels',        'Cavalry and beast-rider yard centered on savage mobility and northern hunt traditions.',  '🐴'),
    ('stable', 'zandres',   'Crawler Pens',        'Stable-yard for specialized underground and engineered mounts.',  '🐴'),
    ('stable', 'lumus',     'Sun Court Stables',   'Cavalry yard where speed, ceremony, and holy discipline meet.',  '🐴'),
    ('stable', 'drakanith', 'Stable',              'Trains cavalry units.',  '🐴'),

    ('archery', 'veridor',   'Chart Range',       'Ranged corps blending port defense, ship warfare, and inland marksman traditions.',  '🏹'),
    ('archery', 'sylvara',   'Grove Range',       'Long-range training grounds blending hunting skill with battlefield marksmanship.',  '🏹'),
    ('archery', 'arkazia',   'Ridge Range',       'Missile training for mountain and fortress warfare, favoring reliable ranged support.',  '🏹'),
    ('archery', 'draxys',    'Oasis Range',       'Missile grounds focused on mobility, heat tolerance, and arena-style ranged skill.',  '🏹'),
    ('archery', 'nordalh',   'Ice Loom Range',    'Missile hall for harsh-weather hunters and long-distance northern marksmen.',  '🏹'),
    ('archery', 'zandres',   'Crystal Range',     'Ranged corps using precise mechanical launchers and resonance-based projectile methods.',  '🏹'),
    ('archery', 'lumus',     'Sunshot Range',     'Ranged school mixing traditional missile training with ceremonial and radiant weapons.',  '🏹'),
    ('archery', 'drakanith', 'Archery Range',     'Trains ranged units.',  '🏹'),

    ('workshop', 'veridor',   'Shipwright Siegeyard',  'Siege dockyard where artillery is built with strong naval and coastal influence.',  '🔨'),
    ('workshop', 'sylvara',   'Tree-Sapper Yard',     'Forest siege workshop using woodcraft, roots, sap, and natural pressure weapons.',  '🔨'),
    ('workshop', 'arkazia',   'Stonecaller Yard',     'Siege and engineering yard built around stone, ramps, walls, and heavy siege utility.',  '🔨'),
    ('workshop', 'draxys',    'Sandwall Foundry',     'Siege workshop producing heat-tough engines, traps, fire weapons, and desert assault tools.',  '🔨'),
    ('workshop', 'nordalh',   'Long Forge Siegeyard', 'Siege foundry emphasizing rugged artillery, brute-force rams, and forge-crafted war engines.',  '🔨'),
    ('workshop', 'zandres',   'Resonance Works',      'Advanced siege and engineering workshop built on pressure, drilling, and charged mechanisms.',  '🔨'),
    ('workshop', 'lumus',     'Heliostat Works',      'Siege yard using mirrors, radiant arrays, and formalized support engines.',  '🔨'),
    ('workshop', 'drakanith', 'Workshop',             'Trains siege units.',  '🔨'),

    ('special', 'veridor',   'Admiralty Hall',      'Special command building for naval elites, captains, and exotic maritime fighters.',  '⭐'),
    ('special', 'sylvara',   'Spirit Glade',        'Sacred special building for rare bonded, mystical, and druidic elites.',  '⭐'),
    ('special', 'arkazia',   'Chapter Fortress',    'Special order hall where elite sworn units, commanders, and prestige warriors are raised.',  '⭐'),
    ('special', 'draxys',    'Grand Arena',         'Signature special building where spectacle fighters, slave champions, and blood elites are trained.',  '⭐'),
    ('special', 'nordalh',   'Long Forge Hall',     'Special elite building, home to champions, smith-warriors, and rune-forged retinues.',  '⭐'),
    ('special', 'zandres',   'Circuit Archive',     'Elite technical hall for power-tech adepts, charged guards, and battlefield specialists.',  '⭐'),
    ('special', 'lumus',     'Heliostat Sanctum',   'Special hall for sacred elites, radiant adepts, and ceremonial champions.',  '⭐'),
    ('special', 'drakanith', 'Special Hall',        'Trains elite units.',  '⭐');

-- ── Building Display Configs: Storage buildings (24 rows: 3 types × 8 kingdoms) ─

INSERT INTO building_display_configs (building_type, kingdom, display_name, description, default_icon) VALUES
    ('storage', 'veridor',   'Storage',    'Increases lumber and stone capacity.',  '🏚️'),
    ('storage', 'sylvara',   'Storage',    'Increases lumber and stone capacity.',  '🏚️'),
    ('storage', 'arkazia',   'Storage',    'Increases lumber and stone capacity.',  '🏚️'),
    ('storage', 'draxys',    'Storage',    'Increases lumber and stone capacity.',  '🏚️'),
    ('storage', 'nordalh',   'Storage',    'Increases lumber and stone capacity.',  '🏚️'),
    ('storage', 'zandres',   'Storage',    'Increases lumber and stone capacity.',  '🏚️'),
    ('storage', 'lumus',     'Storage',    'Increases lumber and stone capacity.',  '🏚️'),
    ('storage', 'drakanith', 'Storage',    'Increases lumber and stone capacity.',  '🏚️'),

    ('provisions', 'veridor',   'Provisions',  'Increases food capacity.',  '🍞'),
    ('provisions', 'sylvara',   'Provisions',  'Increases food capacity.',  '🍞'),
    ('provisions', 'arkazia',   'Provisions',  'Increases food capacity.',  '🍞'),
    ('provisions', 'draxys',    'Provisions',  'Increases food capacity.',  '🍞'),
    ('provisions', 'nordalh',   'Provisions',  'Increases food capacity.',  '🍞'),
    ('provisions', 'zandres',   'Provisions',  'Increases food capacity.',  '🍞'),
    ('provisions', 'lumus',     'Provisions',  'Increases food capacity.',  '🍞'),
    ('provisions', 'drakanith', 'Provisions',  'Increases food capacity.',  '🍞'),

    ('reservoir', 'veridor',   'Reservoir',  'Increases water capacity.',  '🏊'),
    ('reservoir', 'sylvara',   'Reservoir',  'Increases water capacity.',  '🏊'),
    ('reservoir', 'arkazia',   'Reservoir',  'Increases water capacity.',  '🏊'),
    ('reservoir', 'draxys',    'Reservoir',  'Increases water capacity.',  '🏊'),
    ('reservoir', 'nordalh',   'Reservoir',  'Increases water capacity.',  '🏊'),
    ('reservoir', 'zandres',   'Reservoir',  'Increases water capacity.',  '🏊'),
    ('reservoir', 'lumus',     'Reservoir',  'Increases water capacity.',  '🏊'),
    ('reservoir', 'drakanith', 'Reservoir',  'Increases water capacity.',  '🏊');

-- ══════════════════════════════════════════════════════════════════════════════
-- TROOP DISPLAY CONFIGS — per-kingdom unit display names & descriptions
-- 7 kingdoms × ~20 units each = ~140 rows (Drakanith has no units yet)
-- ══════════════════════════════════════════════════════════════════════════════

-- ── Sylvara ──────────────────────────────────────────────────────────────────
INSERT INTO troop_display_configs (troop_type, kingdom, training_building, display_name, description, default_icon) VALUES
    ('sylvara_rootguard_spearmen',  'sylvara', 'barracks', 'Rootguard Spearmen',  'Disciplined defensive infantry with spears and living-wood shields, built to hold paths and choke points.', '⚔️'),
    ('sylvara_thornblade_wardens',  'sylvara', 'barracks', 'Thornblade Wardens',  'Sword-and-buckler fighters, quick and precise, used as the reliable core of Sylvaran lines.', '⚔️'),
    ('sylvara_boarhide_axemen',     'sylvara', 'barracks', 'Boarhide Axemen',     'Tougher shock infantry with heavier axes and hide protection, made to break enemy guards.', '⚔️'),
    ('sylvara_leafknife_stalkers',  'sylvara', 'barracks', 'Leafknife Stalkers',  'Fast skirmishing infantry with paired blades, specialized in flanks and finishing wounded targets.', '⚔️'),
    ('sylvara_hart_lancers',        'sylvara', 'stable',   'Hart Lancers',        'Swift stag-mounted lancers, elegant and fast, ideal for hit-and-withdraw charges.', '🐴'),
    ('sylvara_boar_riders',         'sylvara', 'stable',   'Boar Riders',         'Brutal close-range cavalry on armored boars, used to smash light infantry and disrupt formations.', '🐴'),
    ('sylvara_wolf_outriders',      'sylvara', 'stable',   'Wolf Outriders',      'Fast pursuit riders with javelins, meant for harassment, scouting, and hunting routed enemies.', '🐴'),
    ('sylvara_longbow_wardens',     'sylvara', 'archery',  'Longbow Wardens',     'Disciplined long-range archers, the backbone of Sylvara''s missile line.', '🏹'),
    ('sylvara_owlshot_hunters',     'sylvara', 'archery',  'Owlshot Hunters',     'Lighter archers built for stealth, precision, and forest skirmishing.', '🏹'),
    ('sylvara_thorn_javeliners',    'sylvara', 'archery',  'Thorn Javeliners',    'Mobile throwers with hardened javelins, best at mid-range pressure.', '🏹'),
    ('sylvara_stonesling_foresters','sylvara', 'archery',  'Stone-sling Foresters','Cheap but useful support slingers, strong against lightly armored troops.', '🏹'),
    ('sylvara_vine_ballista',       'sylvara', 'workshop', 'Vine Ballista',       'A living-wood bolt launcher used to pin monsters, cavalry, or siege crews.', '🔨'),
    ('sylvara_log_trebuchet',       'sylvara', 'workshop', 'Log Trebuchet',       'Sylvara''s heavy stone-thrower, built from enormous trunk frames.', '🔨'),
    ('sylvara_root_ram',            'sylvara', 'workshop', 'Root Ram',            'A reinforced siege ram covered in bark and root plating.', '🔨'),
    ('sylvara_sporepot_thrower',    'sylvara', 'workshop', 'Spore-Pot Thrower',   'A specialized engine that hurls toxic or blinding spore vessels.', '🔨'),
    ('sylvara_beastmasters',        'sylvara', 'special',  'Beastmasters',        'Handlers who fight beside sacred beasts and enhance nearby animal units.', '⭐'),
    ('sylvara_shapeshifter_scouts', 'sylvara', 'special',  'Shapeshifter Scouts', 'Elusive elites who blur the line between warrior and forest predator.', '⭐'),
    ('sylvara_life_healers',        'sylvara', 'special',  'Life Healers',        'Battlefield sustain units using restorative nature power instead of pure offense.', '⭐'),
    ('sylvara_antler_seers',        'sylvara', 'special',  'Antler Seers',        'Support mystics who guide, bless, and disrupt through omen and ritual.', '⭐');

-- ── Arkazia ──────────────────────────────────────────────────────────────────
INSERT INTO troop_display_configs (troop_type, kingdom, training_building, display_name, description, default_icon) VALUES
    ('arkazia_shield_guard',        'arkazia', 'barracks', 'Shield Guard',        'Core sword-and-kite infantry, dependable, defensive, and visually iconic for Arkazia.', '⚔️'),
    ('arkazia_pike_guard',          'arkazia', 'barracks', 'Pike Guard',          'Anti-cavalry formation troops trained to hold approaches and punish charges.', '⚔️'),
    ('arkazia_stonebreaker',        'arkazia', 'barracks', 'Stonebreaker',        'Heavy shock infantry with massive hammers, built to crush armor and break frontlines.', '⚔️'),
    ('arkazia_rampart_halberdiers', 'arkazia', 'barracks', 'Rampart Halberdiers', 'Polearm specialists for anti-armor fighting, line support, and wall defense.', '⚔️'),
    ('arkazia_arknight_lancers',    'arkazia', 'stable',   'Arknight Lancers',    'Elite lance cavalry, the classic heavy striking force of the kingdom.', '🐴'),
    ('arkazia_redcrest_cavaliers',  'arkazia', 'stable',   'Redcrest Cavaliers',  'Mace-and-shield riders for disciplined melee cavalry engagements.', '🐴'),
    ('arkazia_banner_riders',       'arkazia', 'stable',   'Banner Riders',       'Prestige cavalry who carry command banners and improve allied morale.', '🐴'),
    ('arkazia_ravine_pursuers',     'arkazia', 'stable',   'Ravine Pursuers',     'Lighter horsemen used for pursuit, skirmishing, and flanking in rough terrain.', '🐴'),
    ('arkazia_hill_slingers',       'arkazia', 'archery',  'Hill Slingers',       'Practical missile troops trained for steep terrain and defensive bombardment.', '🏹'),
    ('arkazia_ridge_crossbowmen',   'arkazia', 'archery',  'Ridge Crossbowmen',   'Armor-piercing ranged infantry designed for wall and choke-point fighting.', '🏹'),
    ('arkazia_javelin_climbers',    'arkazia', 'archery',  'Javelin Climbers',    'Mobile throwers suited to broken ground and close support.', '🏹'),
    ('arkazia_crag_bowmen',         'arkazia', 'archery',  'Crag Bowmen',         'Disciplined war-bow troops with better range than the lighter missile corps.', '🏹'),
    ('arkazia_mountain_trebuchet',  'arkazia', 'workshop', 'Mountain Trebuchet',  'A fortress-grade heavy artillery engine for long-range bombardment.', '🔨'),
    ('arkazia_ram_engineers',       'arkazia', 'workshop', 'Ram Engineers',       'Armored crews who build and escort gate-breaking siege rams.', '🔨'),
    ('arkazia_mantlet_pushers',     'arkazia', 'workshop', 'Mantlet Pushers',     'Shielded advance crews bringing cover to infantry and ranged troops.', '🔨'),
    ('arkazia_bridgewright_crew',   'arkazia', 'workshop', 'Bridgewright Crew',   'Tactical engineers able to create assault paths, crossings, and siege access.', '🔨'),
    ('arkazia_banner_knights',      'arkazia', 'special',  'Banner Knights',      'Knightly infantry or mounted elites who fight as mobile standards of the realm.', '⭐'),
    ('arkazia_oathsworn_champions', 'arkazia', 'special',  'Oathsworn Champions', 'Elite duel-capable warriors trusted for hard breaches and last stands.', '⭐'),
    ('arkazia_bastion_marshals',    'arkazia', 'special',  'Bastion Marshals',    'Heavily protected command fighters who anchor lines and inspire nearby troops.', '⭐'),
    ('arkazia_arknight_captains',   'arkazia', 'special',  'Arknight Captains',   'Rare cavalry leaders representing Arkazia''s highest mounted tradition.', '⭐');

-- ── Veridor ──────────────────────────────────────────────────────────────────
INSERT INTO troop_display_configs (troop_type, kingdom, training_building, display_name, description, default_icon) VALUES
    ('veridor_road_legionaries',    'veridor', 'barracks', 'Road Legionaries',    'Sword-and-kite infantry, orderly and professional, suited to roads, ports, and line battles.', '⚔️'),
    ('veridor_harbor_pike',         'veridor', 'barracks', 'Harbor Pike',         'Disciplined spear-wall infantry used to defend wharves, docks, and narrow approaches.', '⚔️'),
    ('veridor_cutlass_marines',     'veridor', 'barracks', 'Cutlass Marines',     'Close-range naval fighters trained for boarding and shoreline combat.', '⚔️'),
    ('veridor_wharf_axemen',        'veridor', 'barracks', 'Wharf Axemen',        'Rougher heavy infantry with boarding axes, useful in brutal melee pushes.', '⚔️'),
    ('veridor_river_lancers',       'veridor', 'stable',   'River Lancers',       'Clean medium cavalry for disciplined charges and open-ground engagements.', '🐴'),
    ('veridor_courier_riders',      'veridor', 'stable',   'Courier Riders',      'Light cavalry scouts focused on speed, messages, and harassment.', '🐴'),
    ('veridor_marsh_scouts',        'veridor', 'stable',   'Marsh Scouts',        'Javelin cavalry adapted to wetlands, skirmish routes, and awkward terrain.', '🐴'),
    ('veridor_road_wardens',        'veridor', 'stable',   'Road Wardens',        'Heavier patrol cavalry who protect routes, caravans, and border roads.', '🐴'),
    ('veridor_deck_arbalesters',    'veridor', 'archery',  'Deck Arbalesters',    'Crossbow troops strong in controlled volleys and anti-armor shooting.', '🏹'),
    ('veridor_highland_longbowmen', 'veridor', 'archery',  'Highland Longbowmen', 'Longer-range archers recruited from Veridor''s tougher upland regions.', '🏹'),
    ('veridor_harpoon_casters',     'veridor', 'archery',  'Harpoon Casters',     'Specialized throwers built for anti-beast, anti-marine, and close naval range.', '🏹'),
    ('veridor_pavise_marksmen',     'veridor', 'archery',  'Pavise Marksmen',     'Heavy crossbow teams protected by large shields for disciplined firing lines.', '🏹'),
    ('veridor_harbor_ballista',     'veridor', 'workshop', 'Harbor Ballista',     'Precision anti-ship or anti-siege launcher adapted for land defense.', '🔨'),
    ('veridor_mangonel',            'veridor', 'workshop', 'Mangonel',            'Practical stone-throwing siege engine for field and city assaults.', '🔨'),
    ('veridor_firepot_crane',       'veridor', 'workshop', 'Firepot Crane',       'Incendiary support engine designed to hurl burning cargo-like munitions.', '🔨'),
    ('veridor_pavise_wagon',        'veridor', 'workshop', 'Pavise Wagon',        'Mobile cover cart that supports advancing missile troops and siege crews.', '🔨'),
    ('veridor_hydra_hunters',       'veridor', 'special',  'Hydra Hunters',       'Specialist spear-and-net elites trained to face monstrous sea threats.', '⭐'),
    ('veridor_tidemark_duelists',   'veridor', 'special',  'Tidemark Duelists',   'Fast naval duelists who excel in cramped close combat.', '⭐'),
    ('veridor_bluecoat_captains',   'veridor', 'special',  'Bluecoat Captains',   'Command fighters who embody Veridor''s disciplined officer class.', '⭐'),
    ('veridor_skiff_raiders',       'veridor', 'special',  'Skiff Raiders',       'Aggressive marine shock troops using hooks, shields, and boarding tactics.', '⭐');

-- ── Draxys ───────────────────────────────────────────────────────────────────
INSERT INTO troop_display_configs (troop_type, kingdom, training_building, display_name, description, default_icon) VALUES
    ('draxys_sandshield_infantry',  'draxys', 'barracks', 'Sandshield Infantry',  'Reliable spear-and-shield core troops built for heat, dust, and formation war.', '⚔️'),
    ('draxys_khopesh_guard',        'draxys', 'barracks', 'Khopesh Guard',        'Curved-blade infantry with a sharper offensive style than basic line troops.', '⚔️'),
    ('draxys_dune_axemen',          'draxys', 'barracks', 'Dune Axemen',          'Hard-hitting shock infantry used to crack open shielded enemies.', '⚔️'),
    ('draxys_wadi_lashers',         'draxys', 'barracks', 'Wadi Lashers',         'Cruel close-range specialists with whip-and-knife fighting styles.', '⚔️'),
    ('draxys_scorpion_riders',      'draxys', 'stable',   'Scorpion Riders',      'Exotic strike cavalry with poisonous, fearsome battlefield presence.', '🐴'),
    ('draxys_dune_lancers',         'draxys', 'stable',   'Dune Lancers',         'Desert charge cavalry balancing speed and power.', '🐴'),
    ('draxys_camel_skirmishers',    'draxys', 'stable',   'Camel Skirmishers',    'Resilient ranged cavalry made for endurance and harassment.', '🐴'),
    ('draxys_dust_chasers',         'draxys', 'stable',   'Dust Chasers',         'Light pursuit cavalry with curved blades, ideal for routing broken enemies.', '🐴'),
    ('draxys_oasis_rangers',        'draxys', 'archery',  'Oasis Rangers',        'Dependable desert bowmen with strong battlefield flexibility.', '🏹'),
    ('draxys_sun_slingers',         'draxys', 'archery',  'Sun Slingers',         'Cheap and mobile ranged troops suited to pressure and anti-light infantry use.', '🏹'),
    ('draxys_chakram_dancers',      'draxys', 'archery',  'Chakram Dancers',      'Stylized thrown-weapon specialists with strong arena flavor.', '🏹'),
    ('draxys_javelin_skirmishers',  'draxys', 'archery',  'Javelin Skirmishers',  'Aggressive mid-range harassers who work well with fast desert movement.', '🏹'),
    ('draxys_bolt_thrower',         'draxys', 'workshop', 'Bolt Thrower',         'Compact anti-armor launcher with precise heavy shots.', '🔨'),
    ('draxys_firepot_mangonel',     'draxys', 'workshop', 'Firepot Mangonel',     'Incendiary siege engine built for fear, area denial, and city assaults.', '🔨'),
    ('draxys_siege_tower',          'draxys', 'workshop', 'Siege Tower',          'Assault engine for storming walls and fortifications.', '🔨'),
    ('draxys_scorpion_cage_wagon',  'draxys', 'workshop', 'Scorpion Cage Wagon',  'Specialized battlefield terror platform using beasts or venom tactics.', '🔨'),
    ('draxys_gladiators',           'draxys', 'special',  'Gladiators',           'Classic arena warriors with shield and crushing weapon, built for crowd-favorite brutality.', '⭐'),
    ('draxys_netfighters',          'draxys', 'special',  'Netfighters',          'Control specialists using entangling tools and piercing follow-up attacks.', '⭐'),
    ('draxys_arena_spearmen',       'draxys', 'special',  'Arena Spearmen',       'Precise reach-fighters using agile spear work and arena discipline.', '⭐'),
    ('draxys_pit_brutes',           'draxys', 'special',  'Pit Brutes',           'Overwhelming heavy arena shock troops, less refined but terrifying in melee.', '⭐'),
    ('draxys_beast_tamers',         'draxys', 'special',  'Beast Tamers',         'Handlers who direct trained arena beasts into the fight.', '⭐');

-- ── Nordalh ──────────────────────────────────────────────────────────────────
INSERT INTO troop_display_configs (troop_type, kingdom, training_building, display_name, description, default_icon) VALUES
    ('nordalh_hearth_guards',       'nordalh', 'barracks', 'Hearth Guards',       'Dependable axe-and-shield infantry, proud and stubborn in defensive fights.', '⚔️'),
    ('nordalh_fjord_spearmen',      'nordalh', 'barracks', 'Fjord Spearmen',      'Disciplined spear infantry, simpler but effective in line warfare.', '⚔️'),
    ('nordalh_iceshore_raiders',    'nordalh', 'barracks', 'Ice-Shore Raiders',   'Aggressive coastal fighters using blade and thrown axe for fast violence.', '⚔️'),
    ('nordalh_chain_wardens',       'nordalh', 'barracks', 'Chain Wardens',       'Heavier fighters with flails, built to batter shields and unsettle enemy lines.', '⚔️'),
    ('nordalh_direwolf_riders',     'nordalh', 'stable',   'Direwolf Riders',     'Fast, fearsome strike cavalry with strong pursuit identity.', '🐴'),
    ('nordalh_elk_lancers',         'nordalh', 'stable',   'Elk Lancers',         'Unusual but noble northern charge cavalry with a strong kingdom silhouette.', '🐴'),
    ('nordalh_snow_riders',         'nordalh', 'stable',   'Snow Riders',         'Light javelin cavalry meant for scouting, raids, and harassment.', '🐴'),
    ('nordalh_fang_cavaliers',      'nordalh', 'stable',   'Fang Cavaliers',      'Heavier mounted raiders with axes, built for direct melee aggression.', '🐴'),
    ('nordalh_frostbow_hunters',    'nordalh', 'archery',  'Frost Bow Hunters',   'Reliable long-range bow troops shaped by survival and discipline.', '🏹'),
    ('nordalh_harpoon_throwers',    'nordalh', 'archery',  'Harpoon Throwers',    'Heavy missile infantry built to punch into armor, beasts, and larger targets.', '🏹'),
    ('nordalh_cliff_crossbowmen',   'nordalh', 'archery',  'Cliff Crossbowmen',   'Crossbow troops ideal for defending harsh terrain and elevated ground.', '🏹'),
    ('nordalh_storm_slingers',      'nordalh', 'archery',  'Storm Slingers',      'Light ranged support who pressure enemies with cheap but constant fire.', '🏹'),
    ('nordalh_cliff_ballista',      'nordalh', 'workshop', 'Cliff Ballista',      'Heavy launcher ideal for anti-siege and anti-monster roles.', '🔨'),
    ('nordalh_stone_trebuchet',     'nordalh', 'workshop', 'Stone Trebuchet',     'Straightforward heavy bombardment engine built for real destruction.', '🔨'),
    ('nordalh_ram_sled',            'nordalh', 'workshop', 'Ram Sled',            'Reinforced breach device adapted to icy or rough approaches.', '🔨'),
    ('nordalh_boiling_pitch_crew',  'nordalh', 'workshop', 'Boiling Pitch Crew',  'Close-support siege team using heat and burning liquid defense.', '🔨'),
    ('nordalh_smith_retinues',      'nordalh', 'special',  'Smith Retinues',      'Armored hammer elites tied to the kingdom''s forge prestige.', '⭐'),
    ('nordalh_coyote_blademasters', 'nordalh', 'special',  'Coyote Blademasters', 'Swift elite swordsmen with a colder, more dangerous dueling edge.', '⭐'),
    ('nordalh_runeforged_forgers',  'nordalh', 'special',  'Runeforged Forgers',  'Rare heavy elites wielding empowered forge-crafted weapons.', '⭐'),
    ('nordalh_ulfhednar_champions', 'nordalh', 'special',  'Ulfhednar Champions', 'Terrifying berserker-style elites who trade composure for killing force.', '⭐');

-- ── Zandres ──────────────────────────────────────────────────────────────────
INSERT INTO troop_display_configs (troop_type, kingdom, training_building, display_name, description, default_icon) VALUES
    ('zandres_door_wardens',        'zandres', 'barracks', 'Door Wardens',        'Disciplined spear-and-shield infantry designed to hold gates, corridors, and narrow fronts.', '⚔️'),
    ('zandres_karst_pikemen',       'zandres', 'barracks', 'Karst Pikemen',       'Longer-reach infantry used to block charges and dominate lanes underground or above.', '⚔️'),
    ('zandres_lattice_halberdiers', 'zandres', 'barracks', 'Lattice Halberdiers', 'Versatile polearm troops for anti-armor and formation disruption.', '⚔️'),
    ('zandres_survey_suppressors',  'zandres', 'barracks', 'Survey Suppressors',  'Hard utility infantry mixing mace work with tech-assisted control tools.', '⚔️'),
    ('zandres_cave_strider_riders', 'zandres', 'stable',   'Cave Strider Riders', 'Swift reptilian riders used for scouting and tunnel pursuit.', '🐴'),
    ('zandres_beetle_lancers',      'zandres', 'stable',   'Beetle Lancers',      'Heavily protected cavalry with slow, crushing impact.', '🐴'),
    ('zandres_survey_couriers',     'zandres', 'stable',   'Survey Couriers',     'Fast mounted skirmishers linking outposts and harassing weak targets.', '🐴'),
    ('zandres_burrow_guards',       'zandres', 'stable',   'Burrow Guards',       'Heavier mounted enforcers for secure escort and close defense.', '🐴'),
    ('zandres_crystal_boltcasters', 'zandres', 'archery',  'Crystal Boltcasters', 'Disciplined crossbow troops with strong armor-piercing function.', '🏹'),
    ('zandres_resonance_slingers',  'zandres', 'archery',  'Resonance Slingers',  'Technical slingers using crystal ammunition for sharp impact.', '🏹'),
    ('zandres_survey_needlers',     'zandres', 'archery',  'Survey Needlers',     'Rapid-firing light marksmen built for suppression and exact shots.', '🏹'),
    ('zandres_prism_markers',       'zandres', 'archery',  'Prism Markers',       'Ranged support specialists who designate, weaken, or expose enemy targets.', '🏹'),
    ('zandres_resonance_ballista',  'zandres', 'workshop', 'Resonance Ballista',  'Powerful anti-structure or anti-heavy launcher with high precision.', '🔨'),
    ('zandres_drill_ram',           'zandres', 'workshop', 'Drill Ram',           'Reinforced breach engine designed for walls, gates, and fortified tunnels.', '🔨'),
    ('zandres_stonerail_thrower',   'zandres', 'workshop', 'Stone-Rail Thrower',  'Technical bombardment engine using guided or stabilized heavy shots.', '🔨'),
    ('zandres_barrier_cart',        'zandres', 'workshop', 'Barrier Cart',        'Deployable cover-and-control vehicle supporting pushes and defense.', '🔨'),
    ('zandres_powertech_adepts',    'zandres', 'special',  'Power-Tech Adepts',   'Weaponized engineers using rods, packs, and directed energy-like strikes.', '⭐'),
    ('zandres_beacon_surveyors',    'zandres', 'special',  'Beacon Surveyors',    'Advanced support units who guide, mark, and coordinate battlefield movement.', '⭐'),
    ('zandres_capacitor_sentries',  'zandres', 'special',  'Capacitor Sentries',  'Disciplined elite guards whose gear stores and releases heavy impact force.', '⭐'),
    ('zandres_magnet_lashers',      'zandres', 'special',  'Magnet Lashers',      'Rare specialists using metal-drawing or binding tools to disrupt armored foes.', '⭐');

-- ── Lumus ────────────────────────────────────────────────────────────────────
INSERT INTO troop_display_configs (troop_type, kingdom, training_building, display_name, description, default_icon) VALUES
    ('lumus_ringwall_wardens',      'lumus', 'barracks', 'Ring-Wall Wardens',     'Shielded mace infantry who protect lines and sacred ground.', '⚔️'),
    ('lumus_sun_monks',             'lumus', 'barracks', 'Sun Monks',             'Disciplined close-combat fighters using body training and blunt martial weapons.', '⚔️'),
    ('lumus_prism_guards',          'lumus', 'barracks', 'Prism Guards',          'Elite temple soldiers wielding rods or scepters with ritual precision.', '⚔️'),
    ('lumus_eclipse_wardens',       'lumus', 'barracks', 'Eclipse Wardens',       'Balanced defensive elites using staff forms and control-oriented combat.', '⚔️'),
    ('lumus_sunrider_lancers',      'lumus', 'stable',   'Sunrider Lancers',      'Clean, prestigious lance cavalry with a bright noble battlefield image.', '🐴'),
    ('lumus_dawn_couriers',         'lumus', 'stable',   'Dawn Couriers',         'Fast support cavalry for scouting, message carrying, and harassment.', '🐴'),
    ('lumus_halo_riders',           'lumus', 'stable',   'Halo Riders',           'Elegant saber cavalry suited to sweeping melee passes.', '🐴'),
    ('lumus_whitecloak_escorts',    'lumus', 'stable',   'Whitecloak Escorts',    'Honor cavalry used to defend priests, nobles, and high-value units.', '🐴'),
    ('lumus_sunshot_archers',       'lumus', 'archery',  'Sunshot Archers',       'Disciplined archers representing the reliable ranged core of Lumus.', '🏹'),
    ('lumus_halo_chakramists',      'lumus', 'archery',  'Halo Chakramists',      'Stylized thrown-disc fighters with distinctive sacred-war flair.', '🏹'),
    ('lumus_prism_sling_monks',     'lumus', 'archery',  'Prism Sling Monks',     'Light ranged ascetics using simple tools with surprising battlefield value.', '🏹'),
    ('lumus_glare_casters',         'lumus', 'archery',  'Glare Casters',         'Ranged support specialists using light-focused rods to harass and disrupt.', '🏹'),
    ('lumus_mirror_ballista',       'lumus', 'workshop', 'Mirror Ballista',       'Precision launcher associated with reflective and focused battlefield aesthetics.', '🔨'),
    ('lumus_sunfire_trebuchet',     'lumus', 'workshop', 'Sunfire Trebuchet',     'Heavy bombardment engine for long-range siege and holy spectacle.', '🔨'),
    ('lumus_glare_tower',           'lumus', 'workshop', 'Glare Tower',           'Static or movable support piece used for zone denial and ranged dominance.', '🔨'),
    ('lumus_array_cart',            'lumus', 'workshop', 'Array Cart',            'Support platform carrying reflectors, shields, or energy-enhancing battlefield gear.', '🔨'),
    ('lumus_sunchorus_masters',     'lumus', 'special',  'Sun-Chorus Masters',    'Elite staff fighters who blend martial discipline with ritual authority.', '⭐'),
    ('lumus_radiant_duelists',      'lumus', 'special',  'Radiant Duelists',      'Holy close-combat champions, fast and precise but not lightly trained.', '⭐'),
    ('lumus_eclipse_watch',         'lumus', 'special',  'Eclipse Watch',         'Specialized guardians skilled in vision denial, control, and composed defense.', '⭐'),
    ('lumus_prism_adepts',          'lumus', 'special',  'Prism Adepts',          'Support-heavy elites who project barriers, marks, or light-born attacks.', '⭐');

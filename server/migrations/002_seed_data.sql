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

-- ── Game Assets: Village buildings (6 non-resource) ──────────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('town_hall',  'building', 'Town Hall',  '🏛️',  96, 96),
    ('barracks',   'building', 'Barracks',   '⚔️',  96, 96),
    ('stable',     'building', 'Stable',     '🐴',  96, 96),
    ('archery',    'building', 'Archery',    '🏹',  96, 96),
    ('workshop',   'building', 'Workshop',   '🔨',  96, 96),
    ('special',    'building', 'Special',    '⭐',  96, 96);

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

-- ── Game Assets: Zone tiles (8, one per kingdom + 1 default + wilderness/moraphys/dark_reach) ─
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('zone_default',    'zone_tile', 'Default Zone',    '🟩', 256, 256),
    ('zone_veridor',    'zone_tile', 'Veridor Zone',    '🔵', 256, 256),
    ('zone_sylvara',    'zone_tile', 'Sylvara Zone',    '🟢', 256, 256),
    ('zone_arkazia',    'zone_tile', 'Arkazia Zone',    '🔴', 256, 256),
    ('zone_draxys',     'zone_tile', 'Draxys Zone',     '🟡', 256, 256),
    ('zone_nordalh',    'zone_tile', 'Nordalh Zone',    '🟣', 256, 256),
    ('zone_zandres',    'zone_tile', 'Zandres Zone',    '🟤', 256, 256),
    ('zone_lumus',      'zone_tile', 'Lumus Zone',      '⚪', 256, 256),
    ('zone_drakanith',  'zone_tile', 'Drakanith Zone',  '🟠', 256, 256),
    ('zone_wilderness', 'zone_tile', 'Wilderness Zone', '🌿', 256, 256),
    ('zone_moraphys',   'zone_tile', 'Moraphys Zone',   '💀', 256, 256),
    ('zone_dark_reach', 'zone_tile', 'Dark Reach Zone', '🖤', 256, 256);

-- ── Game Assets: Terrain tiles (6, one per terrain type) ─────────────────────
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('terrain_plains',   'terrain_tile', 'Plains',   '🟩', 256, 256),
    ('terrain_forest',   'terrain_tile', 'Forest',   '🌲', 256, 256),
    ('terrain_mountain', 'terrain_tile', 'Mountain', '⛰️', 256, 256),
    ('terrain_water',    'terrain_tile', 'Water',    '🌊', 256, 256),
    ('terrain_desert',   'terrain_tile', 'Desert',   '🏜️', 256, 256),
    ('terrain_swamp',    'terrain_tile', 'Swamp',    '🐊', 256, 256);

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

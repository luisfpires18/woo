-- Game assets: canonical lookup for buildings, resources, and units with emoji fallback + optional sprite.
CREATE TABLE IF NOT EXISTS game_assets (
    id           TEXT PRIMARY KEY,
    category     TEXT NOT NULL CHECK (category IN ('building', 'resource', 'unit')),
    display_name TEXT NOT NULL,
    default_icon TEXT NOT NULL,
    sprite_path  TEXT,
    sprite_width  INTEGER NOT NULL DEFAULT 0,
    sprite_height INTEGER NOT NULL DEFAULT 0,
    updated_at   DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- Seed buildings (17)
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('town_hall',     'building', 'Town Hall',     '🏛️',  96, 96),
    ('iron_mine',     'building', 'Iron Mine',     '⛏️',  96, 96),
    ('lumber_mill',   'building', 'Lumber Mill',   '🪓',  96, 96),
    ('quarry',        'building', 'Quarry',        '🪨',  96, 96),
    ('farm',          'building', 'Farm',          '🌾',  96, 96),
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

-- Seed resources (4)
INSERT INTO game_assets (id, category, display_name, default_icon, sprite_width, sprite_height) VALUES
    ('iron',  'resource', 'Iron',  '⛏️', 32, 32),
    ('wood',  'resource', 'Wood',  '🪵', 32, 32),
    ('stone', 'resource', 'Stone', '🪨', 32, 32),
    ('food',  'resource', 'Food',  '🌾', 32, 32);

-- Units will be seeded later when troop rosters are defined in game-template.md

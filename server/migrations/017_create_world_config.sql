-- World configuration key-value store.
-- Used by admins to tune game parameters at runtime.
CREATE TABLE IF NOT EXISTS world_config (
    key         TEXT PRIMARY KEY,
    value       TEXT NOT NULL,
    description TEXT,
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

-- Seed default configuration values
INSERT INTO world_config (key, value, description) VALUES
    ('game_speed',              '1.0',  'Global game speed multiplier'),
    ('resource_multiplier',     '1.0',  'Resource production multiplier'),
    ('map_width',               '200',  'World map width in tiles'),
    ('map_height',              '200',  'World map height in tiles'),
    ('weapons_of_chaos_count',  '7',    'Number of Weapons of Chaos per world'),
    ('max_villages_per_player', '5',    'Maximum villages a player can own');

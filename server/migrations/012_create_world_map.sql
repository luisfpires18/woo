-- 012_create_world_map.sql
CREATE TABLE IF NOT EXISTS world_map (
    x INTEGER NOT NULL,
    y INTEGER NOT NULL,
    terrain_type TEXT NOT NULL CHECK (terrain_type IN ('plains', 'forest', 'mountain', 'water', 'desert', 'swamp')),
    owner_player_id INTEGER REFERENCES players(id),
    village_id INTEGER REFERENCES villages(id),
    has_oasis INTEGER NOT NULL DEFAULT 0,
    has_chaos_shrine INTEGER NOT NULL DEFAULT 0,
    PRIMARY KEY (x, y)
);

CREATE INDEX IF NOT EXISTS idx_world_map_owner ON world_map(owner_player_id);
CREATE INDEX IF NOT EXISTS idx_world_map_village ON world_map(village_id);

-- 014_create_weapons_of_chaos.sql
CREATE TABLE IF NOT EXISTS weapons_of_chaos (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    name TEXT NOT NULL UNIQUE,
    weapon_type TEXT NOT NULL,
    attack_bonus INTEGER NOT NULL,
    defense_bonus INTEGER NOT NULL,
    effects_json TEXT NOT NULL,
    location_x INTEGER,
    location_y INTEGER,
    wielder_player_id INTEGER REFERENCES players(id),
    held_by_moraphys INTEGER NOT NULL DEFAULT 0,
    claimed_at TEXT
);

CREATE INDEX IF NOT EXISTS idx_woc_wielder ON weapons_of_chaos(wielder_player_id);

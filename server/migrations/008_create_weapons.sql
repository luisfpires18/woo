-- 008_create_weapons.sql
CREATE TABLE IF NOT EXISTS weapons (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL REFERENCES players(id),
    name TEXT NOT NULL,
    weapon_type TEXT NOT NULL CHECK (weapon_type IN ('sword', 'axe', 'bow', 'spear', 'shield', 'staff')),
    tier TEXT NOT NULL CHECK (tier IN ('common', 'rare', 'epic', 'legendary', 'mythic')),
    attack_bonus INTEGER NOT NULL DEFAULT 0,
    defense_bonus INTEGER NOT NULL DEFAULT 0,
    rune_slots INTEGER NOT NULL DEFAULT 0,
    durability INTEGER NOT NULL,
    max_durability INTEGER NOT NULL,
    equipped_on TEXT,
    stats_json TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_weapons_player_id ON weapons(player_id);
CREATE INDEX IF NOT EXISTS idx_weapons_tier ON weapons(tier);

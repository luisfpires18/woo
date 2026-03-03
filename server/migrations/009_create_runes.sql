-- 009_create_runes.sql
CREATE TABLE IF NOT EXISTS runes (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id INTEGER NOT NULL REFERENCES players(id),
    rune_type TEXT NOT NULL,
    rarity TEXT NOT NULL CHECK (rarity IN ('fragment', 'minor', 'major', 'grand', 'primordial')),
    effect_json TEXT NOT NULL,
    weapon_id INTEGER REFERENCES weapons(id),
    created_at TEXT NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_runes_player_id ON runes(player_id);
CREATE INDEX IF NOT EXISTS idx_runes_weapon_id ON runes(weapon_id);
CREATE INDEX IF NOT EXISTS idx_runes_rarity ON runes(rarity);

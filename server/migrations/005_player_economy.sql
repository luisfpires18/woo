-- 005_player_economy.sql — Per-player gold currency.
-- Gold is shared across all villages and earned through gameplay.

CREATE TABLE IF NOT EXISTS player_economy (
    player_id  INTEGER PRIMARY KEY REFERENCES players(id),
    gold       REAL    NOT NULL DEFAULT 0,
    created_at TEXT    NOT NULL DEFAULT (datetime('now'))
);

-- Seed gold for any existing players.
INSERT OR IGNORE INTO player_economy (player_id, gold)
SELECT id, 200 FROM players;

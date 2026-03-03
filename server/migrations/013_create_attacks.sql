-- 013_create_attacks.sql
CREATE TABLE IF NOT EXISTS attacks (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    attacker_player_id INTEGER NOT NULL REFERENCES players(id),
    attacker_village_id INTEGER NOT NULL REFERENCES villages(id),
    target_x INTEGER NOT NULL,
    target_y INTEGER NOT NULL,
    attack_type TEXT NOT NULL CHECK (attack_type IN ('attack', 'raid', 'scout', 'reinforce')),
    troops_json TEXT NOT NULL,
    weapons_json TEXT,
    departed_at TEXT NOT NULL DEFAULT (datetime('now')),
    arrives_at TEXT NOT NULL,
    status TEXT NOT NULL DEFAULT 'marching' CHECK (status IN ('marching', 'arrived', 'returning', 'completed')),
    result_json TEXT
);

CREATE INDEX IF NOT EXISTS idx_attacks_attacker ON attacks(attacker_player_id);
CREATE INDEX IF NOT EXISTS idx_attacks_arrives_at ON attacks(arrives_at);
CREATE INDEX IF NOT EXISTS idx_attacks_status ON attacks(status);

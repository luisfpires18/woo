-- 011_create_alliance_members.sql
CREATE TABLE IF NOT EXISTS alliance_members (
    alliance_id INTEGER NOT NULL REFERENCES alliances(id),
    player_id INTEGER NOT NULL REFERENCES players(id) UNIQUE,
    role TEXT NOT NULL DEFAULT 'member' CHECK (role IN ('leader', 'officer', 'member')),
    joined_at TEXT NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (alliance_id, player_id)
);

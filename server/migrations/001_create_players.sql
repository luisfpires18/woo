-- 001_create_players.sql
CREATE TABLE IF NOT EXISTS players (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    username TEXT NOT NULL UNIQUE,
    email TEXT NOT NULL UNIQUE,
    password_hash TEXT,
    kingdom TEXT NOT NULL DEFAULT '' CHECK (kingdom IN ('', 'veridor', 'sylvara', 'arkazia', 'draxys', 'zandres', 'lumus', 'nordalh', 'drakanith')),
    oauth_provider TEXT,
    oauth_id TEXT,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    last_login_at TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_players_email ON players(email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_players_username ON players(username);
CREATE INDEX IF NOT EXISTS idx_players_oauth ON players(oauth_provider, oauth_id);

-- 001_schema.sql — Baseline schema for all WOO tables.
-- Consolidates the original 001–018 + 022 migration files.
-- During development, edit this file freely and delete woo.db to rebuild.
-- Once production launches, freeze this file and add new numbered migrations.

PRAGMA foreign_keys = ON;

-- ── Players ──────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS players (
    id             INTEGER PRIMARY KEY AUTOINCREMENT,
    username       TEXT    NOT NULL UNIQUE,
    email          TEXT    NOT NULL UNIQUE,
    password_hash  TEXT,
    kingdom        TEXT    NOT NULL DEFAULT '' CHECK (kingdom IN ('', 'veridor', 'sylvara', 'arkazia', 'draxys', 'zandres', 'lumus', 'nordalh', 'drakanith')),
    role           TEXT    NOT NULL DEFAULT 'player' CHECK (role IN ('player', 'admin')),
    oauth_provider TEXT,
    oauth_id       TEXT,
    created_at     TEXT    NOT NULL DEFAULT (datetime('now')),
    last_login_at  TEXT
);

CREATE UNIQUE INDEX IF NOT EXISTS idx_players_email    ON players(email);
CREATE UNIQUE INDEX IF NOT EXISTS idx_players_username ON players(username);
CREATE        INDEX IF NOT EXISTS idx_players_oauth    ON players(oauth_provider, oauth_id);

-- ── Refresh Tokens ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS refresh_tokens (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id  INTEGER NOT NULL REFERENCES players(id),
    token_hash TEXT    NOT NULL UNIQUE,
    expires_at TEXT    NOT NULL,
    created_at TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_refresh_tokens_player_id ON refresh_tokens(player_id);
CREATE INDEX IF NOT EXISTS idx_refresh_tokens_hash      ON refresh_tokens(token_hash);

-- ── Villages ─────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS villages (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id  INTEGER NOT NULL REFERENCES players(id),
    name       TEXT    NOT NULL,
    x          INTEGER NOT NULL,
    y          INTEGER NOT NULL,
    is_capital INTEGER NOT NULL DEFAULT 0,
    created_at TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE        INDEX IF NOT EXISTS idx_villages_player_id ON villages(player_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_villages_coords    ON villages(x, y);

-- ── Buildings ────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS buildings (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    village_id    INTEGER NOT NULL REFERENCES villages(id),
    building_type TEXT    NOT NULL,
    level         INTEGER NOT NULL DEFAULT 0
);

CREATE        INDEX IF NOT EXISTS idx_buildings_village_id   ON buildings(village_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_buildings_village_type ON buildings(village_id, building_type);

-- ── Building Queue ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS building_queue (
    id            INTEGER PRIMARY KEY AUTOINCREMENT,
    village_id    INTEGER NOT NULL REFERENCES villages(id),
    building_type TEXT    NOT NULL,
    target_level  INTEGER NOT NULL,
    started_at    TEXT    NOT NULL DEFAULT (datetime('now')),
    completes_at  TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_building_queue_village_id    ON building_queue(village_id);
CREATE INDEX IF NOT EXISTS idx_building_queue_completes_at  ON building_queue(completes_at);

-- ── Resources (lazy-calculated per village) ──────────────────────────────────
CREATE TABLE IF NOT EXISTS resources (
    village_id       INTEGER PRIMARY KEY REFERENCES villages(id),
    food             REAL NOT NULL DEFAULT 500,
    water            REAL NOT NULL DEFAULT 500,
    lumber           REAL NOT NULL DEFAULT 500,
    stone            REAL NOT NULL DEFAULT 500,
    food_rate        REAL NOT NULL DEFAULT 3,
    water_rate       REAL NOT NULL DEFAULT 3,
    lumber_rate      REAL NOT NULL DEFAULT 3,
    stone_rate       REAL NOT NULL DEFAULT 3,
    food_consumption REAL NOT NULL DEFAULT 0,
    max_storage      REAL NOT NULL DEFAULT 1200,
    last_updated     TEXT NOT NULL DEFAULT (datetime('now'))
);

-- ── Troops ───────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS troops (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    village_id INTEGER NOT NULL REFERENCES villages(id),
    type       TEXT    NOT NULL,
    quantity   INTEGER NOT NULL DEFAULT 0,
    status     TEXT    NOT NULL DEFAULT 'stationed' CHECK (status IN ('stationed', 'marching', 'defending', 'returning'))
);

CREATE        INDEX IF NOT EXISTS idx_troops_village_id   ON troops(village_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_troops_village_type ON troops(village_id, type);

-- ── Training Queue ───────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS training_queue (
    id           INTEGER PRIMARY KEY AUTOINCREMENT,
    village_id   INTEGER NOT NULL REFERENCES villages(id),
    troop_type   TEXT    NOT NULL,
    quantity     INTEGER NOT NULL,
    started_at   TEXT    NOT NULL DEFAULT (datetime('now')),
    completes_at TEXT    NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_training_queue_village_id   ON training_queue(village_id);
CREATE INDEX IF NOT EXISTS idx_training_queue_completes_at ON training_queue(completes_at);

-- ── Weapons ──────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS weapons (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id       INTEGER NOT NULL REFERENCES players(id),
    name            TEXT    NOT NULL,
    weapon_type     TEXT    NOT NULL CHECK (weapon_type IN ('sword', 'axe', 'bow', 'spear', 'shield', 'staff')),
    tier            TEXT    NOT NULL CHECK (tier IN ('common', 'rare', 'epic', 'legendary', 'mythic')),
    attack_bonus    INTEGER NOT NULL DEFAULT 0,
    defense_bonus   INTEGER NOT NULL DEFAULT 0,
    rune_slots      INTEGER NOT NULL DEFAULT 0,
    durability      INTEGER NOT NULL,
    max_durability  INTEGER NOT NULL,
    equipped_on     TEXT,
    stats_json      TEXT,
    created_at      TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_weapons_player_id ON weapons(player_id);
CREATE INDEX IF NOT EXISTS idx_weapons_tier      ON weapons(tier);

-- ── Runes ────────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS runes (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id   INTEGER NOT NULL REFERENCES players(id),
    rune_type   TEXT    NOT NULL,
    rarity      TEXT    NOT NULL CHECK (rarity IN ('fragment', 'minor', 'major', 'grand', 'primordial')),
    effect_json TEXT    NOT NULL,
    weapon_id   INTEGER REFERENCES weapons(id),
    created_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_runes_player_id ON runes(player_id);
CREATE INDEX IF NOT EXISTS idx_runes_weapon_id ON runes(weapon_id);
CREATE INDEX IF NOT EXISTS idx_runes_rarity    ON runes(rarity);

-- ── Alliances ────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS alliances (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    name        TEXT    NOT NULL UNIQUE,
    tag         TEXT    NOT NULL UNIQUE,
    leader_id   INTEGER NOT NULL REFERENCES players(id),
    max_members INTEGER NOT NULL DEFAULT 10,
    created_at  TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE TABLE IF NOT EXISTS alliance_members (
    alliance_id INTEGER NOT NULL REFERENCES alliances(id),
    player_id   INTEGER NOT NULL REFERENCES players(id) UNIQUE,
    role        TEXT    NOT NULL DEFAULT 'member' CHECK (role IN ('leader', 'officer', 'member')),
    joined_at   TEXT    NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (alliance_id, player_id)
);

-- ── World Map ────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS world_map (
    x                INTEGER NOT NULL,
    y                INTEGER NOT NULL,
    terrain_type     TEXT    NOT NULL CHECK (terrain_type IN ('plains')),
    kingdom_zone     TEXT    NOT NULL DEFAULT '' CHECK (kingdom_zone IN ('', 'veridor', 'sylvara', 'arkazia', 'draxys', 'nordalh', 'wilderness')),
    owner_player_id  INTEGER REFERENCES players(id),
    village_id       INTEGER REFERENCES villages(id),
    PRIMARY KEY (x, y)
);

CREATE INDEX IF NOT EXISTS idx_world_map_owner   ON world_map(owner_player_id);
CREATE INDEX IF NOT EXISTS idx_world_map_village ON world_map(village_id);
CREATE INDEX IF NOT EXISTS idx_world_map_zone    ON world_map(kingdom_zone);

-- ── Kingdom Relations (diplomacy standings) ──────────────────────────────────
CREATE TABLE IF NOT EXISTS kingdom_relations (
    kingdom_a  TEXT    NOT NULL,
    kingdom_b  TEXT    NOT NULL,
    standing   INTEGER NOT NULL DEFAULT 0,
    status     TEXT    NOT NULL DEFAULT 'neutral' CHECK (status IN ('allied', 'friendly', 'neutral', 'hostile', 'war')),
    updated_at TEXT    NOT NULL DEFAULT (datetime('now')),
    PRIMARY KEY (kingdom_a, kingdom_b)
);

-- ── Attacks ──────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS attacks (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    attacker_player_id  INTEGER NOT NULL REFERENCES players(id),
    attacker_village_id INTEGER NOT NULL REFERENCES villages(id),
    target_x            INTEGER NOT NULL,
    target_y            INTEGER NOT NULL,
    attack_type         TEXT    NOT NULL CHECK (attack_type IN ('attack', 'raid', 'scout', 'reinforce')),
    troops_json         TEXT    NOT NULL,
    weapons_json        TEXT,
    departed_at         TEXT    NOT NULL DEFAULT (datetime('now')),
    arrives_at          TEXT    NOT NULL,
    status              TEXT    NOT NULL DEFAULT 'marching' CHECK (status IN ('marching', 'arrived', 'returning', 'completed')),
    result_json         TEXT
);

CREATE INDEX IF NOT EXISTS idx_attacks_attacker   ON attacks(attacker_player_id);
CREATE INDEX IF NOT EXISTS idx_attacks_arrives_at ON attacks(arrives_at);
CREATE INDEX IF NOT EXISTS idx_attacks_status     ON attacks(status);

-- ── Weapons of Chaos ─────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS weapons_of_chaos (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    name              TEXT    NOT NULL UNIQUE,
    weapon_type       TEXT    NOT NULL,
    attack_bonus      INTEGER NOT NULL,
    defense_bonus     INTEGER NOT NULL,
    effects_json      TEXT    NOT NULL,
    location_x        INTEGER,
    location_y        INTEGER,
    wielder_player_id INTEGER REFERENCES players(id),
    held_by_moraphys  INTEGER NOT NULL DEFAULT 0,
    claimed_at        TEXT
);

CREATE INDEX IF NOT EXISTS idx_woc_wielder ON weapons_of_chaos(wielder_player_id);

-- ── World Config (admin key-value store) ─────────────────────────────────────
CREATE TABLE IF NOT EXISTS world_config (
    key         TEXT PRIMARY KEY,
    value       TEXT NOT NULL,
    description TEXT,
    updated_at  TEXT NOT NULL DEFAULT (datetime('now'))
);

-- ── Announcements ────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS announcements (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    title      TEXT    NOT NULL,
    content    TEXT    NOT NULL,
    author_id  INTEGER NOT NULL REFERENCES players(id),
    created_at TEXT    NOT NULL DEFAULT (datetime('now')),
    expires_at TEXT
);

-- ── Game Assets (sprites / icons lookup) ─────────────────────────────────────
CREATE TABLE IF NOT EXISTS game_assets (
    id            TEXT PRIMARY KEY,
    category      TEXT     NOT NULL CHECK (category IN ('building', 'resource', 'unit', 'kingdom_flag', 'village_marker', 'zone_tile')),
    display_name  TEXT     NOT NULL,
    default_icon  TEXT     NOT NULL,
    sprite_path   TEXT,
    sprite_width  INTEGER  NOT NULL DEFAULT 0,
    sprite_height INTEGER  NOT NULL DEFAULT 0,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP
);

-- ── Resource Building Configs (admin-customisable per kingdom) ───────────────
CREATE TABLE IF NOT EXISTS resource_building_configs (
    id            INTEGER  PRIMARY KEY AUTOINCREMENT,
    resource_type TEXT     NOT NULL CHECK (resource_type IN ('food', 'water', 'lumber', 'stone')),
    slot          INTEGER  NOT NULL CHECK (slot BETWEEN 1 AND 3),
    kingdom       TEXT     NOT NULL,
    display_name  TEXT     NOT NULL,
    description   TEXT     NOT NULL DEFAULT '',
    default_icon  TEXT     NOT NULL DEFAULT '🏗️',
    sprite_path   TEXT,
    updated_at    DATETIME NOT NULL DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(resource_type, slot, kingdom)
);

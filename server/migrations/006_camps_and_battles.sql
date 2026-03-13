-- ── Beast Templates (admin-configurable beast types) ─────────────────────────
CREATE TABLE IF NOT EXISTS beast_templates (
    id                  INTEGER PRIMARY KEY AUTOINCREMENT,
    name                TEXT    NOT NULL,
    sprite_key          TEXT    NOT NULL DEFAULT '',
    hp                  INTEGER NOT NULL CHECK (hp >= 1),
    attack_power        INTEGER NOT NULL CHECK (attack_power >= 1),
    attack_interval     INTEGER NOT NULL DEFAULT 5 CHECK (attack_interval >= 1),
    defense_percent     REAL    NOT NULL DEFAULT 0 CHECK (defense_percent >= 0 AND defense_percent <= 100),
    crit_chance_percent REAL    NOT NULL DEFAULT 0 CHECK (crit_chance_percent >= 0 AND crit_chance_percent <= 100),
    description         TEXT    NOT NULL DEFAULT '',
    created_at          TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at          TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_by          INTEGER REFERENCES players(id)
);

-- ── Reward Tables ────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS reward_tables (
    id         INTEGER PRIMARY KEY AUTOINCREMENT,
    name       TEXT NOT NULL,
    created_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_at TEXT NOT NULL DEFAULT (datetime('now')),
    updated_by INTEGER REFERENCES players(id)
);

CREATE TABLE IF NOT EXISTS reward_table_entries (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    reward_table_id  INTEGER NOT NULL REFERENCES reward_tables(id) ON DELETE CASCADE,
    reward_type      TEXT    NOT NULL CHECK (reward_type IN ('food', 'water', 'lumber', 'stone', 'gold', 'rune_fragment')),
    min_amount       INTEGER NOT NULL DEFAULT 1 CHECK (min_amount >= 0),
    max_amount       INTEGER NOT NULL DEFAULT 1 CHECK (max_amount >= 1),
    drop_chance_pct  REAL    NOT NULL DEFAULT 100 CHECK (drop_chance_pct > 0 AND drop_chance_pct <= 100),
    created_at       TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_reward_entries_table ON reward_table_entries(reward_table_id);

-- ── Camp Templates (admin-configurable camp types) ───────────────────────────
CREATE TABLE IF NOT EXISTS camp_templates (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    name            TEXT    NOT NULL,
    tier            INTEGER NOT NULL DEFAULT 1 CHECK (tier >= 1 AND tier <= 10),
    sprite_key      TEXT    NOT NULL DEFAULT '',
    description     TEXT    NOT NULL DEFAULT '',
    reward_table_id INTEGER REFERENCES reward_tables(id),
    created_at      TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at      TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_by      INTEGER REFERENCES players(id)
);

-- ── Camp Beast Slots (which beasts spawn in a camp template) ─────────────────
CREATE TABLE IF NOT EXISTS camp_beast_slots (
    id                INTEGER PRIMARY KEY AUTOINCREMENT,
    camp_template_id  INTEGER NOT NULL REFERENCES camp_templates(id) ON DELETE CASCADE,
    beast_template_id INTEGER NOT NULL REFERENCES beast_templates(id),
    min_count         INTEGER NOT NULL DEFAULT 1 CHECK (min_count >= 1),
    max_count         INTEGER NOT NULL DEFAULT 1 CHECK (max_count >= 1 AND max_count >= min_count)
);

CREATE INDEX IF NOT EXISTS idx_camp_beast_slots_template ON camp_beast_slots(camp_template_id);

-- ── Spawn Rules (admin-configurable camp spawning) ───────────────────────────
CREATE TABLE IF NOT EXISTS spawn_rules (
    id                      INTEGER PRIMARY KEY AUTOINCREMENT,
    name                    TEXT    NOT NULL,
    terrain_types_json      TEXT    NOT NULL DEFAULT '["plains"]',
    zone_types_json         TEXT    NOT NULL DEFAULT '["wilderness"]',
    camp_template_pool_json TEXT    NOT NULL DEFAULT '[]',
    max_camps               INTEGER NOT NULL DEFAULT 10 CHECK (max_camps >= 1),
    spawn_interval_sec      INTEGER NOT NULL DEFAULT 60 CHECK (spawn_interval_sec >= 10),
    despawn_after_sec       INTEGER NOT NULL DEFAULT 0,
    min_camp_distance       INTEGER NOT NULL DEFAULT 2 CHECK (min_camp_distance >= 1),
    min_village_distance    INTEGER NOT NULL DEFAULT 2 CHECK (min_village_distance >= 1),
    enabled                 INTEGER NOT NULL DEFAULT 1,
    created_at              TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_at              TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_by              INTEGER REFERENCES players(id)
);

-- ── Battle Tuning (singleton admin config) ───────────────────────────────────
CREATE TABLE IF NOT EXISTS battle_tuning (
    id                      INTEGER PRIMARY KEY CHECK (id = 1),
    tick_duration_ms        INTEGER NOT NULL DEFAULT 200,
    crit_damage_multiplier  REAL    NOT NULL DEFAULT 2.0,
    max_defense_percent     REAL    NOT NULL DEFAULT 90.0,
    max_crit_chance_percent REAL    NOT NULL DEFAULT 75.0,
    min_attack_interval     INTEGER NOT NULL DEFAULT 1,
    march_speed_tiles_per_min REAL  NOT NULL DEFAULT 2.0,
    max_ticks               INTEGER NOT NULL DEFAULT 10000,
    updated_at              TEXT    NOT NULL DEFAULT (datetime('now')),
    updated_by              INTEGER REFERENCES players(id)
);

-- Insert default battle tuning singleton
INSERT OR IGNORE INTO battle_tuning (id) VALUES (1);

-- ── Camps (runtime instances on the map) ─────────────────────────────────────
CREATE TABLE IF NOT EXISTS camps (
    id               INTEGER PRIMARY KEY AUTOINCREMENT,
    camp_template_id INTEGER NOT NULL REFERENCES camp_templates(id),
    tile_x           INTEGER NOT NULL,
    tile_y           INTEGER NOT NULL,
    beasts_json      TEXT    NOT NULL,
    spawned_at       TEXT    NOT NULL DEFAULT (datetime('now')),
    status           TEXT    NOT NULL DEFAULT 'active' CHECK (status IN ('active', 'under_attack', 'cleared')),
    season_id        INTEGER,
    spawn_rule_id    INTEGER REFERENCES spawn_rules(id)
);

CREATE INDEX IF NOT EXISTS idx_camps_status    ON camps(status);
CREATE INDEX IF NOT EXISTS idx_camps_tile      ON camps(tile_x, tile_y);
CREATE INDEX IF NOT EXISTS idx_camps_season    ON camps(season_id);

-- ── Add camp_id to world_map ─────────────────────────────────────────────────
ALTER TABLE world_map ADD COLUMN camp_id INTEGER REFERENCES camps(id);
CREATE INDEX IF NOT EXISTS idx_world_map_camp ON world_map(camp_id);

-- ── Expeditions ──────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS expeditions (
    id          INTEGER PRIMARY KEY AUTOINCREMENT,
    player_id   INTEGER NOT NULL REFERENCES players(id),
    village_id  INTEGER NOT NULL REFERENCES villages(id),
    camp_id     INTEGER NOT NULL REFERENCES camps(id),
    troops_json TEXT    NOT NULL,
    departed_at TEXT    NOT NULL DEFAULT (datetime('now')),
    arrives_at  TEXT    NOT NULL,
    returns_at  TEXT,
    status      TEXT    NOT NULL DEFAULT 'marching' CHECK (status IN ('marching', 'battling', 'returning', 'completed')),
    season_id   INTEGER
);

CREATE INDEX IF NOT EXISTS idx_expeditions_player  ON expeditions(player_id);
CREATE INDEX IF NOT EXISTS idx_expeditions_status  ON expeditions(status);
CREATE INDEX IF NOT EXISTS idx_expeditions_arrives ON expeditions(arrives_at);
CREATE INDEX IF NOT EXISTS idx_expeditions_returns ON expeditions(returns_at);

-- ── Battles ──────────────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS battles (
    id                      INTEGER PRIMARY KEY AUTOINCREMENT,
    expedition_id           INTEGER NOT NULL UNIQUE REFERENCES expeditions(id),
    attacker_snapshot_json  TEXT    NOT NULL,
    defender_snapshot_json  TEXT    NOT NULL,
    result                  TEXT    NOT NULL CHECK (result IN ('attacker_won', 'defender_won', 'draw')),
    attacker_losses_json    TEXT    NOT NULL DEFAULT '{}',
    defender_losses_json    TEXT    NOT NULL DEFAULT '{}',
    rewards_json            TEXT    NOT NULL DEFAULT '[]',
    replay_data             BLOB,
    seed                    INTEGER NOT NULL,
    resolved_at             TEXT    NOT NULL DEFAULT (datetime('now')),
    duration_ticks          INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_battles_expedition ON battles(expedition_id);

-- ── Admin Audit Log ──────────────────────────────────────────────────────────
CREATE TABLE IF NOT EXISTS admin_audit_log (
    id              INTEGER PRIMARY KEY AUTOINCREMENT,
    admin_player_id INTEGER NOT NULL REFERENCES players(id),
    action          TEXT    NOT NULL CHECK (action IN ('create', 'update', 'delete')),
    entity_type     TEXT    NOT NULL,
    entity_id       INTEGER,
    old_value_json  TEXT,
    new_value_json  TEXT,
    created_at      TEXT    NOT NULL DEFAULT (datetime('now'))
);

CREATE INDEX IF NOT EXISTS idx_audit_log_entity ON admin_audit_log(entity_type, entity_id);
CREATE INDEX IF NOT EXISTS idx_audit_log_admin  ON admin_audit_log(admin_player_id);

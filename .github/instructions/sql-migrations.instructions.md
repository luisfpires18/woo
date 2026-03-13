---
description: "Use when writing or modifying SQL migration files. Covers naming, idempotency, foreign keys, index conventions, and the dev vs production workflow."
applyTo: "server/migrations/**/*.sql"
---

# SQL Migration Conventions

Full reference: `docs/06-database/database-guide.md`

## File Naming

Sequential numbering + snake_case: `001_schema.sql`, `002_seed_data.sql`, `003_per_resource_storage.sql`.
Never skip numbers. Check existing files before picking the next number.

## Idempotency

All `CREATE TABLE` and `CREATE INDEX` use `IF NOT EXISTS`:

```sql
CREATE TABLE IF NOT EXISTS players (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    ...
);
CREATE INDEX IF NOT EXISTS idx_players_email ON players(email);
```

## Foreign Keys

Always enable at the top of migration files:
```sql
PRAGMA foreign_keys = ON;
```

## Index Naming

Follow `idx_<table>_<columns>`:
- `idx_villages_player_id`
- `idx_building_queue_completion_time`

## Constraints

Use `CHECK` constraints for enum-like fields:
```sql
role TEXT NOT NULL DEFAULT 'player' CHECK(role IN ('player', 'admin'))
```

## Development vs Production

- **During development** (now): Edit baseline files (`001_schema.sql`, `002_seed_data.sql`) directly. Delete `woo.db` and rebuild.
- **After production launch**: Freeze baselines. Append-only numbered migrations. No editing existing files.

## After Migration Changes

1. Delete database: `Remove-Item "d:\Workspace\WOO\server\data\woo.db*" -Force`
2. Rebuild and restart server.

## Lazy Resource Calculation

Resources use snapshot + timestamp pattern. **Never** create triggers or cron jobs that update resources periodically:
```sql
-- Stored: last known value + timestamp
food_stored REAL NOT NULL DEFAULT 0,
food_rate REAL NOT NULL DEFAULT 0,
last_resource_update DATETIME NOT NULL DEFAULT (datetime('now'))
-- Current = stored + (rate × seconds_since_update)
```

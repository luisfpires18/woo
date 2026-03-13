---
description: "Create a new SQL migration file with proper numbering, conventions, and safety checks."
agent: "agent"
argument-hint: "Describe the schema change (e.g., 'Add alliances table with name, leader_id, created_at')"
---

# New Migration

Generate a numbered SQL migration following WOO database conventions.

## Steps

1. **Check existing migrations** in `server/migrations/` to determine the next sequential number.
2. **Create file**: `server/migrations/{NNN}_{description}.sql`

## Template

```sql
-- Migration {NNN}: {Description}
-- {Brief explanation of what this migration does}

PRAGMA foreign_keys = ON;

CREATE TABLE IF NOT EXISTS {table_name} (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    -- columns...
    created_at DATETIME NOT NULL DEFAULT (datetime('now')),
    updated_at DATETIME NOT NULL DEFAULT (datetime('now'))
);

-- Indexes: idx_{table}_{columns}
CREATE INDEX IF NOT EXISTS idx_{table}_{column} ON {table}({column});
```

## Rules

- `PRAGMA foreign_keys = ON` at top
- `IF NOT EXISTS` on all CREATE statements
- `CHECK` constraints for enum-like fields: `CHECK(status IN ('active', 'inactive'))`
- Index naming: `idx_{table}_{columns}`
- Foreign keys: `REFERENCES {parent_table}(id)`
- For resource fields: use lazy calculation pattern (stored + rate + last_update timestamp)

## After Creating

Remind user:
1. Delete database: `Remove-Item "d:\Workspace\WOO\server\data\woo.db*" -Force`
2. Rebuild server: `cd d:\Workspace\WOO\server; go build -o server.exe ./cmd/server`
3. If a new table was added: create corresponding model, repository interface, and SQLite implementation.

# Database Guide

> SQLite patterns, migration system, repository conventions, and PostgreSQL migration path. Read before writing any database code.

---

## Current Stack: SQLite

| Tool | Purpose |
|------|---------|
| **modernc.org/sqlite** | Pure Go SQLite driver (no CGO dependency) |
| **database/sql** | Standard Go database interface |
| Manual migration runner | Sequential SQL files applied on startup |

### Why SQLite First?

- Zero external dependencies for development
- Single file database — easy to backup, reset, share
- Fast enough for development and early playtesting
- The repository pattern ensures DB-agnostic business logic

### SQLite Configuration

Apply these PRAGMAs on every connection:

```go
func NewConnection(dbPath string) *sql.DB {
    db, err := sql.Open("sqlite", dbPath)
    if err != nil {
        panic(fmt.Sprintf("failed to open database: %v", err))
    }

    // Enable WAL mode for concurrent reads
    db.Exec("PRAGMA journal_mode=WAL")

    // Enable foreign key enforcement
    db.Exec("PRAGMA foreign_keys=ON")

    // Increase cache size for better performance
    db.Exec("PRAGMA cache_size=-64000") // 64MB

    // Synchronous mode: NORMAL is safe with WAL
    db.Exec("PRAGMA synchronous=NORMAL")

    // Set connection pool limits
    db.SetMaxOpenConns(1)   // SQLite only supports 1 writer
    db.SetMaxIdleConns(1)
    db.SetConnMaxLifetime(0) // Don't close idle connections

    return db
}
```

**Important**: SQLite supports only **one writer at a time**. WAL mode allows concurrent reads during writes, but write operations are serialized. This is fine for development but is why we migrate to PostgreSQL for production.

---

## Migration System

### Convention

Migrations are numbered SQL files in `server/migrations/`:

```
migrations/
├── 001_schema.sql          # All CREATE TABLE + CREATE INDEX statements
└── 002_seed_data.sql        # World config, admin accounts, game assets, resource building configs
```

> **Dev-phase approach**: While in active development, we keep a small set of consolidated baseline migrations. When the schema changes, edit the baseline files directly and delete `woo.db` to rebuild from scratch. Switch to append-only numbered migrations once the game reaches production with real player data.

### Migration Runner

```go
func RunMigrations(db *sql.DB, migrationsDir string) error {
    // Create schema_migrations table if not exists
    _, err := db.Exec(`
        CREATE TABLE IF NOT EXISTS schema_migrations (
            version INTEGER PRIMARY KEY,
            applied_at TEXT NOT NULL
        )
    `)
    if err != nil {
        return fmt.Errorf("create migrations table: %w", err)
    }

    // Get current version
    var currentVersion int
    db.QueryRow("SELECT COALESCE(MAX(version), 0) FROM schema_migrations").Scan(&currentVersion)

    // Read and apply pending migrations
    files, _ := os.ReadDir(migrationsDir)
    for _, f := range files {
        version := parseVersion(f.Name()) // Extract number from filename
        if version <= currentVersion {
            continue
        }

        sql, err := os.ReadFile(filepath.Join(migrationsDir, f.Name()))
        if err != nil {
            return fmt.Errorf("read migration %s: %w", f.Name(), err)
        }

        tx, _ := db.Begin()
        if _, err := tx.Exec(string(sql)); err != nil {
            tx.Rollback()
            return fmt.Errorf("apply migration %s: %w", f.Name(), err)
        }
        tx.Exec("INSERT INTO schema_migrations (version, applied_at) VALUES (?, ?)",
            version, time.Now().UTC().Format(time.RFC3339))
        tx.Commit()

        slog.Info("applied migration", "version", version, "file", f.Name())
    }
    return nil
}
```

### Migration Rules

#### Development Phase (current)

1. **Edit baseline files directly.** The schema is consolidated into `001_schema.sql` and `002_seed_data.sql`. Modify them as needed.
2. **Delete `woo.db` to rebuild.** After changing baseline migrations, delete the database file and restart the server.
3. **Each migration is a single transaction.** If it fails, it rolls back completely.
4. **Test migrations**: Write a test that applies all migrations to an in-memory database to verify they work.

#### Production Phase (future)

1. **Never modify an existing migration.** Always create a new numbered file (e.g., `003_add_column.sql`).
2. **Migrations are forward-only.** No down migrations (too risky for production data).
3. **Each migration is a single transaction.** If it fails, it rolls back completely.
4. **Test migrations**: Write a test that applies all migrations to an in-memory database to verify they work.

---

## Repository Pattern

### Interface Definition

Interfaces are defined in `server/internal/repository/interfaces.go`. Each entity gets its own interface.

```go
package repository

import (
    "context"
    "github.com/luisfpires18/woo/internal/model"
)

type PlayerRepository interface {
    Create(ctx context.Context, player *model.Player) error
    GetByID(ctx context.Context, id int64) (*model.Player, error)
    GetByEmail(ctx context.Context, email string) (*model.Player, error)
    GetByOAuth(ctx context.Context, provider, oauthID string) (*model.Player, error)
    UpdateLastLogin(ctx context.Context, id int64) error
}

type VillageRepository interface {
    Create(ctx context.Context, village *model.Village) error
    GetByID(ctx context.Context, id int64) (*model.Village, error)
    ListByPlayerID(ctx context.Context, playerID int64) ([]*model.Village, error)
    Update(ctx context.Context, village *model.Village) error
    GetByCoordinates(ctx context.Context, x, y int) (*model.Village, error)
}

type BuildingRepository interface {
    Create(ctx context.Context, building *model.Building) error
    GetByVillageID(ctx context.Context, villageID int64) ([]*model.Building, error)
    Update(ctx context.Context, building *model.Building) error
}

type ResourceRepository interface {
    Get(ctx context.Context, villageID int64) (*model.Resources, error)
    Update(ctx context.Context, villageID int64, resources *model.Resources) error
}

// ... additional repository interfaces for troops, weapons, runes, attacks, etc.
```

### SQLite Implementation

```go
// repository/sqlite/player_repo.go
package sqlite

import (
    "context"
    "database/sql"
    "fmt"
    "github.com/luisfpires18/woo/internal/model"
)

type playerRepo struct {
    db *sql.DB
}

func NewPlayerRepo(db *sql.DB) *playerRepo {
    return &playerRepo{db: db}
}

func (r *playerRepo) Create(ctx context.Context, player *model.Player) error {
    result, err := r.db.ExecContext(ctx,
        `INSERT INTO players (username, email, password_hash, kingdom, oauth_provider, oauth_id, created_at)
         VALUES (?, ?, ?, ?, ?, ?, ?)`,
        player.Username, player.Email, player.PasswordHash,
        player.Kingdom, player.OAuthProvider, player.OAuthID, player.CreatedAt,
    )
    if err != nil {
        return fmt.Errorf("insert player: %w", err)
    }
    id, _ := result.LastInsertId()
    player.ID = id
    return nil
}

func (r *playerRepo) GetByID(ctx context.Context, id int64) (*model.Player, error) {
    row := r.db.QueryRowContext(ctx,
        `SELECT id, username, email, password_hash, kingdom, oauth_provider, oauth_id, created_at, last_login_at
         FROM players WHERE id = ?`, id,
    )
    return scanPlayer(row)
}

func (r *playerRepo) GetByEmail(ctx context.Context, email string) (*model.Player, error) {
    row := r.db.QueryRowContext(ctx,
        `SELECT id, username, email, password_hash, kingdom, oauth_provider, oauth_id, created_at, last_login_at
         FROM players WHERE email = ?`, email,
    )
    return scanPlayer(row)
}

func scanPlayer(row *sql.Row) (*model.Player, error) {
    var p model.Player
    err := row.Scan(&p.ID, &p.Username, &p.Email, &p.PasswordHash,
        &p.Kingdom, &p.OAuthProvider, &p.OAuthID, &p.CreatedAt, &p.LastLoginAt)
    if err == sql.ErrNoRows {
        return nil, model.ErrNotFound
    }
    if err != nil {
        return nil, fmt.Errorf("scan player: %w", err)
    }
    return &p, nil
}
```

---

## Query Patterns

### Parameterized Queries Only

```go
// GOOD — parameterized
db.QueryRowContext(ctx, "SELECT * FROM players WHERE email = ?", email)

// BAD — string concatenation (SQL injection vulnerability!)
db.QueryRowContext(ctx, "SELECT * FROM players WHERE email = '" + email + "'")
```

### Batch Operations

For inserting multiple rows (e.g., initializing a village's buildings):

```go
func (r *buildingRepo) CreateBatch(ctx context.Context, buildings []*model.Building) error {
    tx, err := r.db.BeginTx(ctx, nil)
    if err != nil {
        return fmt.Errorf("begin transaction: %w", err)
    }
    defer tx.Rollback()

    stmt, err := tx.PrepareContext(ctx,
        "INSERT INTO buildings (village_id, building_type, level) VALUES (?, ?, ?)")
    if err != nil {
        return fmt.Errorf("prepare statement: %w", err)
    }
    defer stmt.Close()

    for _, b := range buildings {
        if _, err := stmt.ExecContext(ctx, b.VillageID, b.BuildingType, b.Level); err != nil {
            return fmt.Errorf("insert building %s: %w", b.BuildingType, err)
        }
    }

    return tx.Commit()
}
```

### Avoid N+1 Queries

```go
// BAD — N+1 (1 query for villages, then N queries for buildings)
villages, _ := villageRepo.ListByPlayerID(ctx, playerID)
for _, v := range villages {
    buildings, _ := buildingRepo.GetByVillageID(ctx, v.ID)
    v.Buildings = buildings
}

// GOOD — Single query with JOIN, or batch query
buildings, _ := buildingRepo.GetByVillageIDs(ctx, villageIDs)
```

---

## Lazy Resource Calculation

**This is critical for performance.** Do NOT update resources in the database every tick/second.

### How It Works

1. The `resources` table stores a **snapshot**: the last known resource values + production rates + timestamp.
2. When reading resources (for display, for building, for anything):
   ```
   current_food = stored_food + (food_rate_per_second × seconds_since_last_update)
   ```
3. The snapshot is **written** only when something changes the resources:
   - Player starts building (deducts resources)
   - Player trains troops (deducts resources)
   - Player receives a trade
   - Village is attacked/raided (resources stolen)
   - Player logs in (refresh snapshot)

### Why?

- A game with 10,000 villages updating resources every second = 10,000 DB writes/sec = database meltdown.
- Lazy calculation = 0 DB writes/sec for idle villages. Only active villages get writes.

### Food Consumption

Food is special — troops consume food. Net food rate can be negative:

```
food_rate = farm_production - troop_upkeep
```

If `food_rate < 0` and `current_food <= 0`:
- Troops start dying (starvation)
- Handled during resource snapshot writes or periodic starvation checks

---

## PostgreSQL Migration Path

When the game is ready for production, SQLite will be replaced with PostgreSQL.

### What Changes

| Aspect | SQLite | PostgreSQL |
|--------|--------|-----------|
| Driver | `modernc.org/sqlite` | `github.com/jackc/pgx/v5` |
| Auto-increment | `INTEGER PRIMARY KEY AUTOINCREMENT` | `SERIAL` or `BIGSERIAL` |
| Timestamps | `TEXT` (ISO 8601 string) | `TIMESTAMPTZ` |
| JSON columns | `TEXT` (parsed in Go) | `JSONB` (indexable, queryable) |
| Booleans | `INTEGER` (0/1) | `BOOLEAN` |
| Concurrency | Single-writer (WAL helps reads) | Full concurrent read/write |
| Connection pool | `MaxOpenConns(1)` | `MaxOpenConns(25+)` |

### What Stays the Same

- **Repository interfaces** — unchanged
- **Service layer** — unchanged
- **Handler layer** — unchanged
- **Domain models** — unchanged
- **Migration file names** — same convention, new SQL syntax

### Migration Steps (When Ready)

1. Create `repository/postgres/` package implementing same interfaces
2. Rewrite migration SQL files for PostgreSQL syntax
3. Write data export/import script (SQLite → PostgreSQL)
4. Swap driver in `config.go` and `main.go`
5. Adjust connection pool settings
6. Test thoroughly with integration tests

---

## Indexing Strategy

| Table | Column(s) | Type | Why |
|-------|----------|------|-----|
| players | email | UNIQUE | Login lookup |
| players | username | UNIQUE | Display name uniqueness |
| players | (oauth_provider, oauth_id) | INDEX | OAuth login lookup |
| villages | player_id | INDEX | List player's villages |
| villages | (x, y) | UNIQUE | Map coordinate uniqueness |
| buildings | village_id | INDEX | List village buildings |
| building_queue | completes_at | INDEX | Game tick: find completed builds |
| building_queue | village_id | INDEX | Check village's build queue |
| resources | (PK = village_id) | — | Direct lookup |
| troops | village_id | INDEX | List village troops |
| world_map | (PK = x, y) | — | Tile lookup |
| world_map | owner_player_id | INDEX | Find player's territory |
| attacks | arrives_at | INDEX | Game tick: find arriving troops |
| attacks | attacker_player_id | INDEX | Player's outgoing attacks |
| weapons_of_chaos | wielder_player_id | INDEX | Find wielder's weapons |

---

## Testing Database Code

Use **in-memory SQLite** for repository tests:

```go
func setupTestDB(t *testing.T) *sql.DB {
    db, err := sql.Open("sqlite", ":memory:")
    if err != nil {
        t.Fatal(err)
    }
    db.Exec("PRAGMA foreign_keys=ON")

    // Apply all migrations
    RunMigrations(db, "../../../migrations")

    t.Cleanup(func() { db.Close() })
    return db
}

func TestPlayerRepo_Create(t *testing.T) {
    db := setupTestDB(t)
    repo := sqlite.NewPlayerRepo(db)

    player := &model.Player{
        Username: "testuser",
        Email:    "test@example.com",
        Kingdom:  "arkazia",
    }

    err := repo.Create(context.Background(), player)
    if err != nil {
        t.Fatalf("expected no error, got %v", err)
    }
    if player.ID == 0 {
        t.Fatal("expected player ID to be set")
    }
}
```

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of database guide |
| 2026-03-03 | Fixed import paths to github.com/luisfpires18/woo/internal/model, fixed buildings.type→building_type, expanded migration list to 15 files |
| 2026-03-03 | Consolidated 22 migrations into 2 baseline files (001_schema.sql + 002_seed_data.sql). Updated migration rules with dev-phase vs production-phase guidance |
| 2026-03-10 | Removed stale migration files (020–024). All schema consolidated into 001_schema.sql (training_queue.each_duration_sec, seasons tables, season_id FKs). 002_seed_data.sql includes season seeds. UnitOfWork pattern replaces per-repo WithTx variants—services no longer see *sql.Tx. PRAGMA errors now checked and logged. |
| 2026-03-09 | Added development seed for user wright: one capital village with all 21 building slots pre-built at level 1, and town_hall at level 3. |
| 2026-03-10 | Updated wright seed: set kingdom='arkazia', enrolled in dev_season with Arkazia kingdom via season_players table. Wright now fully ready for season-based gameplay testing. |

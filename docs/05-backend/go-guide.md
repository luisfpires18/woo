# Go Guide

> All conventions for the Go backend. Read before writing any server code.

---

## Technology Stack

| Tool | Purpose |
|------|---------|
| **Go 1.22+** | Server language |
| **net/http** | HTTP server (standard library — avoid frameworks) |
| **gorilla/websocket** | WebSocket library |
| **modernc.org/sqlite** | Pure Go SQLite driver (no CGO required) |
| **golang-jwt/jwt/v5** | JWT generation and validation |
| **golang.org/x/crypto/bcrypt** | Password hashing |
| **slog** | Structured logging (standard library, Go 1.21+) |
| **golangci-lint** | Linting and static analysis |
| **godotenv** | Environment variable loading (dev) |

---

## Project Structure

```
server/
├── cmd/
│   └── server/
│       └── main.go                 # Entry point — wires dependencies, starts server
├── internal/
│   ├── config/
│   │   └── config.go               # Load configuration from env vars / config file
│   ├── handler/
│   │   ├── auth_handler.go          # HTTP handlers for auth endpoints
│   │   ├── auth_handler_test.go
│   │   ├── village_handler.go       # HTTP handlers for village endpoints
│   │   ├── village_handler_test.go
│   │   ├── map_handler.go
│   │   └── ...
│   ├── service/
│   │   ├── auth_service.go          # Auth business logic
│   │   ├── auth_service_test.go
│   │   ├── village_service.go       # Village business logic
│   │   ├── village_service_test.go
│   │   ├── resource_service.go      # Resource calculation (lazy evaluation)
│   │   └── ...
│   ├── repository/
│   │   ├── interfaces.go            # Repository interfaces (PlayerRepo, VillageRepo, etc.)
│   │   ├── sqlite/
│   │   │   ├── connection.go        # SQLite database connection setup
│   │   │   ├── player_repo.go       # PlayerRepo SQLite implementation
│   │   │   ├── player_repo_test.go
│   │   │   ├── village_repo.go
│   │   │   └── ...
│   │   └── postgres/                # Future PostgreSQL implementations
│   ├── model/
│   │   ├── player.go                # Domain model structs (not DB models — no SQL tags)
│   │   ├── village.go
│   │   ├── building.go
│   │   ├── resources.go
│   │   ├── troop.go
│   │   ├── weapon.go
│   │   ├── rune.go
│   │   └── errors.go               # Domain-specific error types
│   ├── middleware/
│   │   ├── auth.go                  # JWT validation middleware
│   │   ├── ratelimit.go             # Request rate limiting
│   │   ├── logging.go               # Request/response logging
│   │   └── cors.go                  # CORS configuration
│   ├── websocket/
│   │   ├── hub.go                   # Central WebSocket hub (manages all connections)
│   │   ├── client.go                # Individual WebSocket client
│   │   ├── messages.go              # Message type definitions
│   │   └── handlers.go              # WebSocket message handlers
│   └── gameloop/
│       ├── ticker.go                # Main game tick loop
│       ├── resource_tick.go         # Resource production completion
│       ├── building_tick.go         # Building queue completion
│       └── combat_tick.go           # Troop movement and combat resolution
├── migrations/
│   ├── 001_create_players.sql
│   ├── 002_create_villages.sql
│   └── ...
├── go.mod
├── go.sum
└── Makefile
```

---

## Architecture: Clean Layers

```
HTTP Request / WebSocket Message
         │
    ┌────▼────┐
    │ Handler  │   Parse request, validate input format, call service, return response
    └────┬────┘   NO business logic. NO SQL.
         │
    ┌────▼────┐
    │ Service  │   Business logic, validation rules, game calculations, orchestration
    └────┬────┘   Calls repository. Returns domain models. NO SQL.
         │
    ┌────▼──────┐
    │ Repository │  Data access ONLY. SQL queries, row scanning, marshaling
    └───────────┘  Returns domain models. Defined as INTERFACES.
```

### Rules

1. **Handlers** never import `repository`. They only know about `service`.
2. **Services** never import `database/sql`. They only know about `repository` interfaces.
3. **Repositories** never contain business logic. They are pure CRUD.
4. **Models** are plain Go structs. No SQL tags, no JSON tags for DB operations. JSON tags are only for API responses (in handler layer or separate DTO structs).

---

## Dependency Injection

Use **constructor injection** — no global variables, no init() magic.

```go
// cmd/server/main.go
func main() {
    cfg := config.Load()

    db := sqlite.NewConnection(cfg.DatabasePath)
    defer db.Close()

    // Repositories
    playerRepo := sqlite.NewPlayerRepo(db)
    villageRepo := sqlite.NewVillageRepo(db)

    // Services
    authService := service.NewAuthService(playerRepo, cfg.JWTSecret)
    villageService := service.NewVillageService(villageRepo)

    // Handlers
    authHandler := handler.NewAuthHandler(authService)
    villageHandler := handler.NewVillageHandler(villageService)

    // Router
    mux := http.NewServeMux()
    mux.HandleFunc("POST /api/auth/register", authHandler.Register)
    mux.HandleFunc("POST /api/auth/login", authHandler.Login)
    mux.HandleFunc("GET /api/villages", villageHandler.List)
    // ...

    // Middleware
    stack := middleware.Chain(
        middleware.Logging,
        middleware.CORS(cfg.AllowedOrigins),
        middleware.RateLimit(cfg.RateLimitPerSecond),
    )

    server := &http.Server{
        Addr:    ":" + cfg.Port,
        Handler: stack(mux),
    }

    slog.Info("server starting", "port", cfg.Port)
    log.Fatal(server.ListenAndServe())
}
```

---

## Interface Design

Define interfaces **where they are consumed**, not where they are implemented.

```go
// repository/interfaces.go
package repository

import (
    "context"
    "server/internal/model"
)

type PlayerRepository interface {
    Create(ctx context.Context, player *model.Player) error
    GetByID(ctx context.Context, id int64) (*model.Player, error)
    GetByEmail(ctx context.Context, email string) (*model.Player, error)
    GetByOAuth(ctx context.Context, provider, oauthID string) (*model.Player, error)
}

type VillageRepository interface {
    Create(ctx context.Context, village *model.Village) error
    GetByID(ctx context.Context, id int64) (*model.Village, error)
    ListByPlayerID(ctx context.Context, playerID int64) ([]*model.Village, error)
    Update(ctx context.Context, village *model.Village) error
}

type ResourceRepository interface {
    Get(ctx context.Context, villageID int64) (*model.Resources, error)
    Update(ctx context.Context, villageID int64, resources *model.Resources) error
}

// ... more interfaces as needed
```

---

## Error Handling

### Domain Errors

```go
// model/errors.go
package model

import "errors"

var (
    ErrNotFound          = errors.New("not found")
    ErrAlreadyExists     = errors.New("already exists")
    ErrInvalidInput      = errors.New("invalid input")
    ErrInsufficientRes   = errors.New("insufficient resources")
    ErrUnauthorized      = errors.New("unauthorized")
    ErrForbidden         = errors.New("forbidden")
    ErrBuildingInProgress = errors.New("building already in progress")
)
```

### Error Wrapping (Mandatory)

Always wrap errors with context:

```go
// Good
func (s *VillageService) GetVillage(ctx context.Context, id int64) (*model.Village, error) {
    village, err := s.repo.GetByID(ctx, id)
    if err != nil {
        return nil, fmt.Errorf("get village %d: %w", id, err)
    }
    return village, nil
}

// Bad — no context
func (s *VillageService) GetVillage(ctx context.Context, id int64) (*model.Village, error) {
    return s.repo.GetByID(ctx, id)  // DON'T: loses context on error
}
```

### HTTP Error Responses

```go
// handler/helpers.go
func writeError(w http.ResponseWriter, status int, message string) {
    w.Header().Set("Content-Type", "application/json")
    w.WriteHeader(status)
    json.NewEncoder(w).Encode(map[string]string{"error": message})
}

// Usage in handler
func (h *VillageHandler) Get(w http.ResponseWriter, r *http.Request) {
    village, err := h.service.GetVillage(r.Context(), villageID)
    if err != nil {
        if errors.Is(err, model.ErrNotFound) {
            writeError(w, http.StatusNotFound, "village not found")
            return
        }
        slog.Error("failed to get village", "error", err, "village_id", villageID)
        writeError(w, http.StatusInternalServerError, "internal server error")
        return
    }
    writeJSON(w, http.StatusOK, village)
}
```

---

## Context Propagation

**Every function that performs I/O takes `context.Context` as its first parameter.**

```go
// Good
func (r *sqlitePlayerRepo) GetByID(ctx context.Context, id int64) (*model.Player, error) {
    row := r.db.QueryRowContext(ctx, "SELECT ... WHERE id = ?", id)
    // ...
}

// Bad
func (r *sqlitePlayerRepo) GetByID(id int64) (*model.Player, error) {
    row := r.db.QueryRow("SELECT ... WHERE id = ?", id)  // DON'T: no context
}
```

---

## Logging

Use Go's built-in `slog` (structured logging):

```go
import "log/slog"

// Info level for normal operations
slog.Info("player registered", "player_id", player.ID, "kingdom", player.Kingdom)

// Warn for expected but notable events
slog.Warn("rate limit exceeded", "player_id", playerID, "endpoint", endpoint)

// Error for unexpected failures
slog.Error("failed to create village", "error", err, "player_id", playerID)
```

### Rules

- **No `fmt.Println` in production code.** Use `slog`.
- **No logging of sensitive data**: passwords, tokens, full email addresses.
- **Structured fields**: Always use key-value pairs, not formatted strings.

---

## Configuration

```go
// internal/config/config.go
package config

import "os"

type Config struct {
    Port               string
    DatabasePath       string
    JWTSecret          string
    AllowedOrigins     []string
    RateLimitPerSecond int
    GoogleClientID     string
    GoogleClientSecret string
    DiscordClientID    string
    DiscordClientSecret string
}

func Load() *Config {
    return &Config{
        Port:               getEnv("PORT", "8080"),
        DatabasePath:       getEnv("DB_PATH", "./data/woo.db"),
        JWTSecret:          mustGetEnv("JWT_SECRET"),
        // ...
    }
}

func getEnv(key, fallback string) string {
    if val := os.Getenv(key); val != "" {
        return val
    }
    return fallback
}

func mustGetEnv(key string) string {
    val := os.Getenv(key)
    if val == "" {
        panic("required env var not set: " + key)
    }
    return val
}
```

---

## Concurrency Patterns

### Game Loop

The game loop runs in a dedicated goroutine (see `docs/05-backend/go-multiplayer.md`).

### Village State Locking

Per-village operations (building, resource calculation) use per-village mutexes to avoid global locks:

```go
type VillageLockManager struct {
    mu    sync.Mutex
    locks map[int64]*sync.RWMutex
}

func (m *VillageLockManager) Lock(villageID int64) {
    m.getOrCreate(villageID).Lock()
}

func (m *VillageLockManager) RLock(villageID int64) {
    m.getOrCreate(villageID).RLock()
}
```

### Worker Pools

For batch operations (combat resolution, resource recalculation across many villages):

```go
func processInParallel(ctx context.Context, items []int64, worker func(context.Context, int64) error, maxWorkers int) error {
    sem := make(chan struct{}, maxWorkers)
    g, ctx := errgroup.WithContext(ctx)

    for _, item := range items {
        item := item // capture
        g.Go(func() error {
            sem <- struct{}{}       // acquire
            defer func() { <-sem }() // release
            return worker(ctx, item)
        })
    }

    return g.Wait()
}
```

---

## Code Style

### Naming

| Entity | Convention | Example |
|--------|-----------|---------|
| Packages | lowercase, short | `handler`, `service`, `model` |
| Exported types | PascalCase | `VillageService`, `PlayerRepository` |
| Unexported types | camelCase | `sqlitePlayerRepo` |
| Functions | PascalCase (exported), camelCase (unexported) | `GetByID`, `parseToken` |
| Constants | PascalCase or ALL_CAPS for env keys | `MaxBuildingLevel`, `DefaultPort` |
| Interfaces | PascalCase, describe behavior | `PlayerRepository`, `TokenGenerator` |

### File Naming

- snake_case always: `village_handler.go`, `auth_service.go`
- Test files: `village_handler_test.go`
- One primary type per file (handler, service, or repo)
- Keep files < 300 lines. Split if larger.

### Formatting

- Run `gofmt` before every commit.
- Run `golangci-lint run` before every commit.
- No dead code. No unused imports.

---

## Makefile

```makefile
.PHONY: run build test lint migrate

run:
	go run cmd/server/main.go

build:
	go build -o bin/server cmd/server/main.go

test:
	go test ./... -v -race -cover

lint:
	golangci-lint run ./...

migrate:
	go run cmd/migrate/main.go
```

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of Go guide |

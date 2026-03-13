---
name: scaffold-endpoint
description: "Multi-step workflow to scaffold a complete REST API endpoint: model, DTO, repository interface, SQLite implementation, service, handler, tests, and route registration. Use when adding a new entity or domain to the backend."
argument-hint: "Entity name and operations (e.g., 'Alliance: Create, List, GetByID, AddMember, RemoveMember')"
---

# Scaffold Endpoint

End-to-end workflow for adding a new REST API entity to the WOO backend. Generates all layers of the clean architecture with tests.

## When to Use

- Adding a new game entity (alliance, expedition, trade, report, etc.)
- Adding a new CRUD domain that needs the full handler → service → repository stack

## Procedure

### Step 1: Gather Requirements

Determine from the user:
- Entity name (singular, PascalCase for Go, snake_case for SQL)
- Fields and their types
- Which CRUD operations are needed
- Any relationships to existing entities (player, village, etc.)
- Whether it needs a migration (new table) or uses existing tables

### Step 2: Create Model

File: `server/internal/model/{entity}.go`

Reference pattern: [model conventions](./references/model-pattern.md)

```go
package model

type {Entity} struct {
    ID        int64
    PlayerID  int64     // if player-owned
    // ... fields
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### Step 3: Create DTO

File: `server/internal/dto/{entity}.go`

Separate request and response types. Never expose internal model directly.

### Step 4: Add Repository Interface

Add to: `server/internal/repository/interfaces.go`

```go
type {Entity}Repository interface {
    Create(ctx context.Context, entity *model.{Entity}) error
    GetByID(ctx context.Context, id int64) (*model.{Entity}, error)
    // ... as needed
}
```

### Step 5: Implement SQLite Repository

File: `server/internal/repository/sqlite/{entity}_repo.go`

- Parameterized queries ONLY
- Constructor: `New{Entity}Repo(db *sql.DB) *{Entity}Repo`
- All methods take `context.Context` as first param

### Step 6: Create Service

File: `server/internal/service/{entity}_service.go`

- Constructor accepts repository interfaces
- Uses `repository.UnitOfWork` for transactions
- Wraps all errors: `fmt.Errorf("context: %w", err)`
- Validates against config where applicable
- Checks ownership, prerequisites, resource sufficiency

Reference: `server/internal/service/building_service.go`

### Step 7: Create Handler

File: `server/internal/handler/{entity}_handler.go`

- Constructor accepts service
- `RegisterRoutes(mux *http.ServeMux)` for route registration
- Extract player from `middleware.PlayerIDFromContext()`
- Map service errors to HTTP status codes
- Use `writeJSON` / `writeError` helpers

Reference: `server/internal/handler/village_handler.go`

### Step 8: Create Tests

Files:
- `server/internal/service/{entity}_service_test.go` — table-driven with `setup{Entity}Test(t)`
- `server/internal/handler/{entity}_handler_test.go` — full-stack with real DB

Reference: `server/internal/service/building_service_test.go`

### Step 9: Create Migration (if needed)

File: `server/migrations/{NNN}_{entity}.sql`

Follow the `new-migration` prompt conventions. After creating, delete woo.db and rebuild.

### Step 10: Wire in main.go

Add to `server/cmd/server/main.go`:
1. Create repository: `{entity}Repo := sqlite.New{Entity}Repo(db)`
2. Create service: `{entity}Svc := service.New{Entity}Service(uow, {entity}Repo, ...)`
3. Create handler: `{entity}Handler := handler.New{Entity}Handler({entity}Svc)`
4. Register routes: `{entity}Handler.RegisterRoutes(mux)`

### Step 11: Verify

1. `cd server && go build ./cmd/server` — compiles
2. `cd server && go test ./internal/service/ -run {Entity}` — service tests pass
3. `cd server && go test ./internal/handler/ -run {Entity}` — handler tests pass

---
description: "Scaffold a new REST API endpoint with handler, service, repository interface, SQLite implementation, and test stubs following WOO clean architecture."
agent: "agent"
argument-hint: "Describe the entity and CRUD operations needed (e.g., 'Expedition entity with Create, List, GetByID')"
---

# New Endpoint Scaffold

Generate all layers for a new REST endpoint following the project's clean architecture.

## What to Generate

For the given entity/domain, create these files:

### 1. Repository Interface (`server/internal/repository/interfaces.go`)

Add the new interface to the existing file:
```go
type {Entity}Repository interface {
    Create(ctx context.Context, entity *model.{Entity}) error
    GetByID(ctx context.Context, id int64) (*model.{Entity}, error)
    // ... operations as needed
}
```

### 2. Model (`server/internal/model/{entity}.go`)

```go
type {Entity} struct {
    ID        int64
    PlayerID  int64
    // ... fields
    CreatedAt time.Time
    UpdatedAt time.Time
}
```

### 3. DTO (`server/internal/dto/{entity}.go`)

Request and response types — decouple API contract from internal models.

### 4. SQLite Repository (`server/internal/repository/sqlite/{entity}_repo.go`)

Implement the interface with parameterized queries. Constructor: `NewRepo(db *sql.DB)`.

### 5. Service (`server/internal/service/{entity}_service.go`)

- Constructor accepts repository interfaces via DI
- Business logic, config lookups, error wrapping with `fmt.Errorf("context: %w", err)`
- Uses `repository.UnitOfWork` for multi-step transactions

### 6. Handler (`server/internal/handler/{entity}_handler.go`)

- Constructor accepts service
- `RegisterRoutes(mux *http.ServeMux)` method
- Extract player ID from context via `middleware.PlayerIDFromContext()`
- Delegate to service, map errors to HTTP status codes
- Use `writeJSON` / `writeError` helpers

### 7. Test Files

- `server/internal/service/{entity}_service_test.go` — table-driven tests with `setup{Entity}Test(t)` helper
- `server/internal/handler/{entity}_handler_test.go` — full-stack handler tests with real DB

### 8. Wire in main.go

Add the new handler registration to `server/cmd/server/main.go`.

## Reference Patterns

Study these existing files for patterns:
- Handler: `server/internal/handler/village_handler.go`
- Service: `server/internal/service/building_service.go`
- Repository: `server/internal/repository/interfaces.go`
- Tests: `server/internal/service/building_service_test.go`

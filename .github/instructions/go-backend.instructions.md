---
description: "Use when writing Go backend code: handlers, services, repositories, middleware, models, DTOs. Covers clean architecture, error handling, dependency injection, and server-authoritative patterns."
applyTo: "server/**/*.go"
---

# Go Backend Conventions

Full reference: `docs/05-backend/go-guide.md`

## Clean Architecture (Strictly Enforced)

**Handler → Service → Repository** — no skipping layers.

- **Handlers**: Extract params, validate request shape, delegate to service, write response. No business logic.
- **Services**: All business logic. Depend on repository *interfaces* (not implementations). Wrap errors with context.
- **Repositories**: Pure CRUD. Interface defined in `repository/interfaces.go`, implemented in `repository/sqlite/`.

Handlers never import `repository`. Services never import `net/http`.

## Dependency Injection

- Constructors accept interfaces: `NewBuildingService(uow repository.UnitOfWork, villageRepo repository.VillageRepository, ...)`
- Wired in `cmd/server/main.go`. No global variables.
- Interfaces defined where consumed (`repository/interfaces.go`), not where implemented.

## Error Handling

- **Always wrap**: `fmt.Errorf("start upgrade: %w", err)`
- **Sentinel errors** in `model/errors.go`: `ErrNotFound`, `ErrInsufficientResources`, etc.
- **Handlers switch on sentinels**: `errors.Is(err, model.ErrNotFound)` → 404, `model.ErrNotOwner` → 403.
- **Services return domain errors**, never HTTP status codes.

## Context Propagation

Every function doing I/O takes `context.Context` as first parameter:
```go
func (s *BuildingService) StartUpgrade(ctx context.Context, playerID, villageID int64, slot string) (*dto.BuildingResponse, error)
```

## Transactions

Use `repository.UnitOfWork` for multi-step operations (deduct resources + queue build).
Named transactional methods on UoW: `CompleteBuildingUpgrade()`, `DeductResourcesAndInsertBuildQueue()`.

## Response Envelope

All responses via `writeJSON(w, status, data)` / `writeError(w, status, message)` from `handler/helpers.go`.
Envelope: `{"data": ...}` or `{"error": "message"}`.

## Server-Authoritative

- ALL game logic runs server-side. Client sends intents, server validates.
- Never trust client timestamps, resource values, or damage calculations.
- Validate ownership, prerequisites, and sufficient resources before every action.

## Logging

Use `slog` (structured logging). Never `fmt.Println` or `log.Println`.

## SQL Safety

Parameterized queries ONLY. Never concatenate user input into SQL strings.

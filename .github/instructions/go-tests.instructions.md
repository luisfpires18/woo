---
description: "Use when writing or modifying Go test files. Covers table-driven tests, test helpers, full-stack handler tests, and the WOO testing patterns."
applyTo: "server/**/*_test.go"
---

# Go Test Conventions

Full reference: `docs/07-testing/testing-guide.md`

## Table-Driven Tests (Mandatory)

```go
tests := []struct {
    name     string
    input    string
    expected int
    wantErr  bool
}{
    {"valid input", "food_1", 1, false},
    {"invalid slot", "bad", 0, true},
}
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        // ...
    })
}
```

## Test Helpers

Use `setup*Test(t *testing.T)` pattern — creates fresh DB, repos, service, and fixtures:

```go
func setupBuildingTest(t *testing.T) (*BuildingService, int64, int64) {
    t.Helper()
    db := testutil.NewTestDB(t)
    // Create repos, service, seed player + village
    return svc, playerID, villageID
}
```

- `testutil.NewTestDB(t)` — fresh in-memory SQLite with all migrations applied.
- `authCtx(playerID)` — injects authenticated context for handler tests.

## Full-Stack Handler Tests (Preferred)

Wire real SQLite + all repositories + all services + handler. This catches integration bugs mocks miss:

```go
func TestHandler_ListVillages(t *testing.T) {
    // Real DB, real repos, real services
    req := httptest.NewRequest("GET", "/api/villages", nil)
    req = req.WithContext(authCtx(playerID))
    rec := httptest.NewRecorder()
    handler.ListVillages(rec, req)
    // Assert status + response body
}
```

## Fresh DB Per Test

Every test starts with a clean database. Never rely on state from previous tests.

## Context

All test methods use `context.Background()` for contexts passed to services/repos.

## Assertions

Use explicit assertions — `if err != nil { t.Fatalf(...) }`. No assertion libraries required.

## What to Test

- **Services**: Business logic, validation, error cases, edge cases.
- **Handlers**: HTTP flow (status codes, response shape, auth).
- **Config**: Parity tests verify Go config matches generated JSON.
- Skip integration tests with `if testing.Short() { t.Skip() }`.

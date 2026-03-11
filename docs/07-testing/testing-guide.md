# Testing Guide

> Test strategy for all layers — frontend, backend, and integration. Read before writing any tests.

---

## Philosophy

- **No code without tests.** Every feature PR must include tests.
- **Test behavior, not implementation.** Tests should verify what a function does, not how it does it.
- **Fast tests by default.** Unit tests must run in milliseconds. Integration tests can take seconds.
- **Deterministic.** No flaky tests. No reliance on timing, network, or external services in unit tests.

---

## Backend Testing (Go)

### Tools

| Tool | Purpose |
|------|---------|
| `testing` (stdlib) | Test framework |
| `httptest` (stdlib) | HTTP handler testing |
| `database/sql` + `:memory:` | In-memory SQLite for repo tests |
| `-race` flag | Race condition detection |
| `-cover` flag | Code coverage reporting |

### Running Tests

```bash
# All tests with race detection and coverage
go test ./... -v -race -cover

# Specific package
go test ./internal/service/... -v

# Specific test
go test ./internal/service/ -run TestVillageService_CreateVillage -v
```

### Table-Driven Tests (Mandatory Pattern)

All service tests must use table-driven tests:

```go
func TestResourceService_CalculateCurrent(t *testing.T) {
    tests := []struct {
        name        string
        stored      model.Resources
        elapsed     time.Duration
        wantFood    float64
        wantLumber  float64
    }{
        {
            name: "1 hour elapsed with base rates",
            stored: model.Resources{
                Food: 100, FoodRate: 30,
                Lumber: 200, LumberRate: 25,
            },
            elapsed:    1 * time.Hour,
            wantFood:   130,
            wantLumber: 225,
        },
        {
            name: "caps at max storage",
            stored: model.Resources{
                Food: 900, FoodRate: 200,
                MaxFoodStorage: 1000,
            },
            elapsed:  1 * time.Hour,
            wantFood: 1000, // capped, not 1100
        },
        {
            name: "zero elapsed returns stored values",
            stored: model.Resources{
                Food: 500, FoodRate: 100,
            },
            elapsed:  0,
            wantFood: 500,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            svc := service.NewResourceService(nil) // nil repo for calc-only tests
            result := svc.Calculate(tt.stored, tt.elapsed)

            if result.Food != tt.wantFood {
                t.Errorf("food: got %.2f, want %.2f", result.Food, tt.wantFood)
            }
            if tt.wantLumber > 0 && result.Lumber != tt.wantLumber {
                t.Errorf("lumber: got %.2f, want %.2f", result.Lumber, tt.wantLumber)
            }
        })
    }
}
```

### Handler Tests (Full-Stack with In-Memory DB)

**Do NOT mock services.** Instead, wire a full in-memory test environment with real services and SQLite:

```go
// handler_test_helpers_test.go
func newTestEnv(t *testing.T) *testEnv {
    // Create in-memory SQLite DB and apply migrations
    db := testutil.NewTestDB(t)

    // Wire all repositories
    playerRepo := sqlite.NewPlayerRepository(db)
    villageRepo := sqlite.NewVillageRepository(db)
    // ... wire all other repos

    // Wire all services
    authSvc := service.NewAuthService(playerRepo, tokenRepo)
    villageSvc := service.NewVillageService(villageRepo, resourceRepo)
    // ... wire all other services

    // Wire all handlers
    authHandler := handler.NewAuthHandler(authSvc)
    villageHandler := handler.NewVillageHandler(villageSvc)
    // ... wire all other handlers

    return &testEnv{
        DB: db,
        AuthHandler: authHandler,
        VillageHandler: villageHandler,
    }
}

// Use authCtx() helper to inject authenticated context without running JWT middleware
func authCtx(playerID int64, role string) context.Context {
    ctx := context.Background()
    ctx = context.WithValue(ctx, middleware.ContextKeyPlayerID, playerID)
    ctx = context.WithValue(ctx, middleware.ContextKeyRole, role)
    return ctx
}

// Use setupArkaziaPlayer() helper to quickly set up a test player
func setupArkaziaPlayer(t *testing.T, env *testEnv) (int64, int64) {
    playerID, _ := registerAndLogin(t, env, "testuser", "test@test.com", "Strong@123")
    chooseKingdomForPlayer(t, env, playerID, "arkazia")

    // Level up barracks to 1 for training tests
    _, err := env.DB.ExecContext(context.Background(),
        `UPDATE buildings SET level = 1 WHERE village_id = ? AND building_type = 'barracks'`, villageID)

    return playerID, villageID
}
```

**Test example:**

```go
func TestVillageHandler_ListVillages_Success(t *testing.T) {
    env := newTestEnv(t)

    // Register and setup player
    playerID, _ := registerAndLogin(t, env, "viluser", "vil@test.com", "Strong@123")
    chooseKingdomForPlayer(t, env, playerID, "veridor")

    // Make HTTP request with authenticated context
    req := httptest.NewRequest("GET", "/api/villages", nil)
    req = req.WithContext(authCtx(playerID, "player"))
    rec := httptest.NewRecorder()

    env.VillageHandler.ListVillages(rec, req)

    if rec.Code != http.StatusOK {
        t.Fatalf("status: got %d, want %d", rec.Code, http.StatusOK)
    }

    // Decode response and assert
    var villages []json.RawMessage
    json.Unmarshal(rec.Body.Bytes(), &villages)
    if len(villages) != 1 {
        t.Errorf("village count: got %d, want 1", len(villages))
    }
}
```

**Advantages:**
- Tests the entire stack: HTTP parsing, handler logic, service logic, repository logic, DB state
- No mocking = no mock-reality gap
- Catches integration bugs early (e.g., SQL errors, transaction issues)
- Real DB state between requests allows testing complex sequences

### Repository Tests

Use in-memory SQLite:

```go
func TestPlayerRepo_GetByEmail(t *testing.T) {
    db := setupTestDB(t) // :memory: SQLite with migrations applied
    repo := sqlite.NewPlayerRepo(db)
    ctx := context.Background()

    // Seed data
    player := &model.Player{
        Username: "testuser",
        Email:    "test@example.com",
        Kingdom:  "veridor",
        CreatedAt: time.Now().UTC().Format(time.RFC3339),
    }
    repo.Create(ctx, player)

    // Test
    found, err := repo.GetByEmail(ctx, "test@example.com")
    if err != nil {
        t.Fatalf("unexpected error: %v", err)
    }
    if found.Username != "testuser" {
        t.Errorf("username: got %s, want testuser", found.Username)
    }

    // Test not found
    _, err = repo.GetByEmail(ctx, "nonexistent@example.com")
    if !errors.Is(err, model.ErrNotFound) {
        t.Errorf("expected ErrNotFound, got %v", err)
    }
}
```

### Mock Pattern

Mock interfaces using function-based structs:

```go
type mockPlayerRepo struct {
    createFn     func(ctx context.Context, p *model.Player) error
    getByIDFn    func(ctx context.Context, id int64) (*model.Player, error)
    getByEmailFn func(ctx context.Context, email string) (*model.Player, error)
}

func (m *mockPlayerRepo) Create(ctx context.Context, p *model.Player) error {
    return m.createFn(ctx, p)
}

func (m *mockPlayerRepo) GetByID(ctx context.Context, id int64) (*model.Player, error) {
    return m.getByIDFn(ctx, id)
}

func (m *mockPlayerRepo) GetByEmail(ctx context.Context, email string) (*model.Player, error) {
    return m.getByEmailFn(ctx, email)
}
```

---

## Frontend Testing (React)

### Tools

| Tool | Purpose |
|------|---------|
| **Vitest** | Test runner (Vite-native, fast) |
| **React Testing Library** | Component rendering and assertions |
| **@testing-library/user-event** | Simulate user interactions |
| **MSW (Mock Service Worker)** | Mock API and WebSocket in tests |

### Running Tests

```bash
# All tests
npm run test

# Watch mode
npm run test:watch

# Coverage
npm run test:coverage
```

### Component Tests

Test components through user interactions, not internal state:

```tsx
// Button.test.tsx
import { render, screen } from '@testing-library/react';
import userEvent from '@testing-library/user-event';
import { Button } from './Button';

describe('Button', () => {
  it('renders with label', () => {
    render(<Button label="Build" onClick={() => {}} />);
    expect(screen.getByText('Build')).toBeInTheDocument();
  });

  it('calls onClick when clicked', async () => {
    const onClick = vi.fn();
    render(<Button label="Build" onClick={onClick} />);

    await userEvent.click(screen.getByText('Build'));
    expect(onClick).toHaveBeenCalledOnce();
  });

  it('is disabled when disabled prop is true', () => {
    render(<Button label="Build" onClick={() => {}} disabled />);
    expect(screen.getByText('Build')).toBeDisabled();
  });

  it('applies variant class', () => {
    render(<Button label="Danger" onClick={() => {}} variant="danger" />);
    const btn = screen.getByText('Danger');
    expect(btn.className).toContain('danger');
  });
});
```

### Hook Tests

Test custom hooks with `renderHook`:

```tsx
import { renderHook, act } from '@testing-library/react';
import { useVillageStore } from '../stores/villageStore';

describe('useVillageStore', () => {
  it('sets active village', () => {
    const { result } = renderHook(() => useVillageStore());

    act(() => {
      result.current.setActiveVillage({ id: 1, name: 'Test Village', x: 10, y: 20 });
    });

    expect(result.current.activeVillage?.name).toBe('Test Village');
  });
});
```

### API Mocking with MSW

```tsx
// test/mocks/handlers.ts
import { http, HttpResponse } from 'msw';

export const handlers = [
  http.get('/api/villages', () => {
    return HttpResponse.json([
      { id: 1, name: 'Capital', x: 50, y: 50, isCapital: true },
    ]);
  }),

  http.post('/api/auth/login', async ({ request }) => {
    const body = await request.json();
    if (body.email === 'test@test.com') {
      return HttpResponse.json({ accessToken: 'mock-token' });
    }
    return HttpResponse.json({ error: 'Invalid credentials' }, { status: 401 });
  }),
];
```

---

## Integration Tests

### What to Integration Test

| Scenario | Level | What Gets Tested |
|----------|-------|-----------------|
| Auth flow: register → login → refresh token | Backend | Handler → Service → Repo → DB |
| Build queue: start build → wait → complete | Backend | Service → Repo → GameLoop tick |
| WebSocket: connect → subscribe → receive event | Full stack | Client WS → Server Hub → Game event |
| Attack flow: send troops → arrive → combat result | Backend | Multiple services + game loop |

### Backend Integration Tests

Use a real (in-memory) SQLite database and test across layers:

```go
func TestIntegration_BuildingFlow(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping integration test")
    }

    db := setupTestDB(t)
    // Wire up real services (not mocks)
    playerRepo := sqlite.NewPlayerRepo(db)
    villageRepo := sqlite.NewVillageRepo(db)
    buildingRepo := sqlite.NewBuildingRepo(db)
    resourceRepo := sqlite.NewResourceRepo(db)

    villageSvc := service.NewVillageService(villageRepo, buildingRepo, resourceRepo)

    // 1. Create player and village
    // 2. Start building construction
    // 3. Verify resources deducted
    // 4. Simulate time passing (or call building completion directly)
    // 5. Verify building level increased
}
```

### Run Integration Tests

```bash
# Skip integration tests in quick runs
go test ./... -short

# Run everything including integration
go test ./... -v -race
```

---

## Coverage Targets

| Layer | Target | Notes |
|-------|--------|-------|
| **Services** (Go) | 80%+ | Core business logic — most critical |
| **Handlers** (Go) | 70%+ | Request/response handling |
| **Repositories** (Go) | 70%+ | SQL query correctness |
| **React Components** | 60%+ | Interaction + rendering tests |
| **Hooks / Stores** | 70%+ | State management logic |
| **Utilities** | 90%+ | Pure functions — easy to test |

### Anti-Cheat Tests

Specifically test that the server rejects:

1. **Insufficient resources**: Attempt to build without enough resources → expect rejection
2. **Invalid building prerequisites**: Build something requiring Town Hall level 5 with level 2 → expect rejection
3. **Double-build**: Start two builds simultaneously → expect second rejected
4. **Rapid-fire actions**: Send 100 build commands in 1 second → expect rate limiting
5. **Invalid coordinates**: Attack a tile outside map bounds → expect rejection
6. **Ownership violations**: Try to build in someone else's village → expect rejection
7. **Tampered payloads**: Send negative resource values, SQL injection, XSS → expect sanitization

---

## Test File Organization

### Backend

```
server/internal/
├── handler/
│   ├── auth_handler.go
│   └── auth_handler_test.go      # Tests in same package
├── service/
│   ├── auth_service.go
│   └── auth_service_test.go
├── repository/sqlite/
│   ├── player_repo.go
│   └── player_repo_test.go
└── integration/                   # Integration test package
    └── building_flow_test.go
```

### Frontend

```
client/src/
├── components/Button/
│   ├── Button.tsx
│   └── Button.test.tsx
├── features/auth/
│   ├── components/LoginForm.tsx
│   └── components/LoginForm.test.tsx
├── hooks/
│   ├── useAuth.ts
│   └── useAuth.test.ts
└── test/                          # Test infrastructure
    ├── mocks/
    │   └── handlers.ts            # MSW handlers
    └── setup.ts                   # Test setup (MSW server, etc.)
```

---

## CI Pipeline (Future)

```yaml
# .github/workflows/test.yml
on: [push, pull_request]
jobs:
  backend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-go@v5
        with:
          go-version: '1.22'
      - run: cd server && go test ./... -v -race -cover

  frontend:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4
      - uses: actions/setup-node@v4
        with:
          node-version: '20'
      - run: cd client && npm ci && npm run test:coverage
```

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of testing guide |
| 2026-03-10 | Updated resource test examples: Iron/Wood → Food/Lumber with per-resource storage caps. |

# System Architecture

> Technical architecture for the full stack. Read before implementing any infrastructure, networking, or auth features.

---

## High-Level Architecture

```
┌─────────────────────────────────────────────────────────────┐
│                        BROWSER (Client)                     │
│                                                             │
│  React 18+ (TypeScript) ── Vite ── Zustand ── React Query  │
│                                                             │
│  ┌──────────┐  ┌──────────┐  ┌──────────┐  ┌────────────┐ │
│  │  REST API │  │WebSocket │  │  OAuth   │  │  Static    │ │
│  │  Calls   │  │Connection│  │  Flows   │  │  Assets    │ │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘  └────────────┘ │
└───────┼──────────────┼─────────────┼────────────────────────┘
        │              │             │
   HTTPS/JSON     WSS/JSON      OAuth Redirect
        │              │             │
┌───────┼──────────────┼─────────────┼────────────────────────┐
│       │         GO SERVER          │                        │
│  ┌────▼─────┐  ┌────▼─────┐  ┌────▼─────┐                 │
│  │  REST    │  │WebSocket │  │  OAuth   │                  │
│  │ Handler  │  │   Hub    │  │ Handler  │                  │
│  └────┬─────┘  └────┬─────┘  └────┬─────┘                 │
│       │              │             │                        │
│  ┌────▼──────────────▼─────────────▼─────┐                 │
│  │              MIDDLEWARE                │                 │
│  │  Auth │ Rate Limit │ Logging │ CORS   │                 │
│  └────────────────┬──────────────────────┘                 │
│                   │                                         │
│  ┌────────────────▼──────────────────────┐                 │
│  │           SERVICE LAYER               │                 │
│  │  AuthService │ VillageService │ ...   │                 │
│  └────────────────┬──────────────────────┘                 │
│                   │                                         │
│  ┌────────────────▼──────────────────────┐                 │
│  │         REPOSITORY LAYER              │                 │
│  │  Interface │ SQLite Impl │ (PG Impl)  │                 │
│  └────────────────┬──────────────────────┘                 │
│                   │                                         │
└───────────────────┼─────────────────────────────────────────┘
                    │
              ┌─────▼─────┐
              │  SQLite    │
              │ (→ Postgres│
              │   later)   │
              └────────────┘
```

---

## Client Architecture

### Technology

| Tool | Purpose |
|------|---------|
| **React 18+** | UI framework |
| **TypeScript** | Type safety (strict mode) |
| **Vite** | Build tool + dev server |
| **Zustand** | Global game state (resources, village, troops) |
| **React Query (TanStack Query)** | Server data fetching, caching, synchronization |
| **React Router** | Client-side routing |
| **CSS Modules** | Scoped component styling |

### Client Folder Structure

```
client/
├── public/                     # Static assets (favicon, fonts, images)
├── src/
│   ├── components/             # Reusable UI components
│   │   ├── Button/
│   │   │   ├── Button.tsx
│   │   │   └── Button.module.css
│   │   ├── Modal/
│   │   ├── Card/
│   │   ├── ResourceBar/
│   │   └── ... 
│   ├── features/               # Feature-specific modules
│   │   ├── auth/               # Login, register, OAuth
│   │   ├── village/            # Village view, buildings, resources
│   │   ├── map/                # World map rendering
│   │   ├── combat/             # Troop management, attacks
│   │   ├── forge/              # Weapon crafting, runes
│   │   └── lore/               # Single-player lore explorer
│   ├── hooks/                  # Custom React hooks
│   │   ├── useAuth.ts
│   │   ├── useWebSocket.ts
│   │   ├── useVillageData.ts
│   │   └── ...
│   ├── services/               # API + WebSocket service layer
│   │   ├── api.ts              # REST API client (axios or fetch wrapper)
│   │   ├── websocket.ts        # WebSocket connection manager
│   │   └── auth.ts             # Auth-specific API calls
│   ├── stores/                 # Zustand stores
│   │   ├── authStore.ts
│   │   ├── villageStore.ts
│   │   ├── mapStore.ts
│   │   └── gameStore.ts
│   ├── styles/                 # Global styles
│   │   ├── globals.css         # Reset, variables, fonts
│   │   ├── themes.css          # Kingdom theme variables
│   │   └── typography.css      # Font definitions
│   ├── types/                  # TypeScript interfaces
│   │   ├── api.ts              # API response/request types
│   │   ├── game.ts             # Game entity types
│   │   ├── websocket.ts        # WebSocket message types
│   │   └── village.ts          # Village-related types
│   ├── utils/                  # Pure utility functions
│   │   ├── format.ts           # Number/date formatting
│   │   ├── calculations.ts     # Display-only calculations (resource ETA, etc.)
│   │   └── constants.ts        # Client-side constants
│   ├── App.tsx                 # Root component with routing
│   ├── main.tsx                # Entry point
│   └── vite-env.d.ts           # Vite type declarations
├── index.html
├── package.json
├── tsconfig.json
├── vite.config.ts
└── .eslintrc.cjs
```

### State Management Strategy

| State Type | Tool | Examples |
|-----------|------|---------|
| **Server data** | React Query | Village details, player profile, building list, world map tiles |
| **Real-time game state** | Zustand | Current resources (ticking), troop positions, active events |
| **Auth state** | Zustand | Current user, token, login status |
| **UI state** | React useState | Modal open/close, selected tab, form inputs |

### WebSocket Client Flow

```
1. User logs in (REST) → receives JWT access token
2. Client opens WebSocket: ws://server/ws?token=<JWT>
3. Server validates token → accepts or rejects connection
4. Client subscribes to events: { type: "subscribe", topics: ["village:123", "map:region:5,10"] }
5. Server sends events: { type: "resource_update", data: { ... } }
6. Client updates Zustand store → React re-renders
```

---

## Server Architecture

### Technology

| Tool | Purpose |
|------|---------|
| **Go 1.22+** | Server language |
| **net/http** | HTTP server (standard library) |
| **coder/websocket** | WebSocket connections (maintained fork of nhooyr/websocket) |
| **modernc.org/sqlite** | SQLite driver (pure Go, no CGO) |
| **golang-jwt/jwt/v5** | JWT token generation and validation |
| **golang.org/x/crypto/bcrypt** | Password hashing |
| **golangci-lint** | Linting |

### Server Folder Structure

```
server/
├── cmd/
│   └── server/
│       └── main.go             # Entry point, wires everything together
├── internal/
│   ├── config/
│   │   └── config.go           # Load env vars, config file
│   ├── handler/
│   │   ├── auth_handler.go     # Login, register, OAuth endpoints
│   │   ├── village_handler.go  # Village CRUD endpoints
│   │   ├── map_handler.go      # World map endpoints
│   │   └── ...
│   ├── service/
│   │   ├── auth_service.go     # Auth business logic
│   │   ├── village_service.go  # Village business logic
│   │   ├── resource_service.go # Resource calculation logic
│   │   └── ...
│   ├── repository/
│   │   ├── interfaces.go       # Repository interfaces
│   │   ├── sqlite/
│   │   │   ├── player_repo.go  # SQLite implementation
│   │   │   ├── village_repo.go
│   │   │   └── ...
│   │   └── postgres/           # (Future) PostgreSQL implementation
│   │       └── ...
│   ├── model/
│   │   ├── player.go           # Domain structs
│   │   ├── village.go
│   │   ├── building.go
│   │   └── ...
│   ├── middleware/
│   │   ├── auth.go             # JWT validation middleware
│   │   ├── ratelimit.go        # Rate limiting
│   │   ├── logging.go          # Request logging
│   │   └── cors.go             # CORS headers
│   ├── websocket/
│   │   ├── hub.go              # Central WebSocket hub (manages clients, broadcast)
│   │   ├── client.go           # Individual client connection (read/write pumps)
│   │   ├── messages.go         # Message type definitions & data structs
│   │   └── handler.go          # HTTP upgrade handler with JWT auth via ?token=│   ├── dto/                    # Data Transfer Objects (request/response structs)
│   │   ├── auth.go
│   │   ├── village.go
│   │   └── map.go│   └── gameloop/
│       ├── ticker.go           # Game tick loop
│       ├── resource_tick.go    # Resource production per tick
│       ├── building_tick.go    # Building queue completion
│       └── combat_tick.go      # Troop movement + combat resolution
├── migrations/
│   ├── 001_create_players.sql
│   ├── 002_create_villages.sql
│   └── ...
├── go.mod
├── go.sum
└── Makefile
```

### Clean Architecture Layers

```
Handler (HTTP/WS) → Service (Business Logic) → Repository (Data Access)
```

- **Handler**: Parses HTTP requests / WebSocket messages, calls service, returns response. No business logic.
- **Service**: All business logic. Calls repository for data. Returns domain models. Handles validation, calculations, and rules.
- **Repository**: Data access only. SQL queries, marshaling/unmarshaling. Returns domain models. Defined as **interfaces** — implementation is swappable.

---

## Authentication Flow

### JWT + Email/Password

```
┌────────┐                    ┌────────┐                    ┌────────┐
│ Client │                    │ Server │                    │   DB   │
└───┬────┘                    └───┬────┘                    └───┬────┘
    │  POST /api/auth/register    │                             │
    │  {email, password, kingdom} │                             │
    ├────────────────────────────►│                             │
    │                             │  bcrypt(password)           │
    │                             │  INSERT player              │
    │                             ├────────────────────────────►│
    │                             │                             │
    │  { accessToken, refresh }   │                             │
    │◄────────────────────────────┤                             │
    │                             │                             │
    │  POST /api/auth/login       │                             │
    │  {email, password}          │                             │
    ├────────────────────────────►│                             │
    │                             │  SELECT player, bcrypt.Cmp  │
    │                             ├────────────────────────────►│
    │  { accessToken, refresh }   │                             │
    │◄────────────────────────────┤                             │
```

### OAuth (Google / Discord)

```
1. Client redirects to: /api/auth/oauth/{provider}
2. Server redirects to provider's OAuth page
3. User authenticates with provider
4. Provider redirects back to: /api/auth/oauth/{provider}/callback?code=XXX
5. Server exchanges code for provider token
6. Server fetches user profile from provider
7. Server creates/links player account
8. Server issues JWT access + refresh tokens
9. Client stores tokens and proceeds
```

### Token Strategy

| Token | Type | Lifetime | Storage |
|-------|------|----------|---------|
| **Access Token** | JWT (signed) | 15 minutes | Memory (Zustand store) |
| **Refresh Token** | Opaque UUID | 7 days | HTTP-only cookie |

- Access token is sent in `Authorization: Bearer <token>` header for REST and as query param for WebSocket connection.
- Refresh token is automatically sent via cookie. Server endpoint `/api/auth/refresh` issues new access token.

---

## WebSocket Protocol

All WebSocket messages are JSON with a `type` field.

### Client → Server Messages

```json
{ "type": "subscribe", "data": { "topics": ["village:123", "map:5,10"] } }
{ "type": "unsubscribe", "data": { "topics": ["village:123"] } }
{ "type": "build", "data": { "village_id": 123, "building_type": "barracks", "target_level": 2 } }
{ "type": "train", "data": { "village_id": 123, "unit_type": "iron_legionary", "quantity": 10 } }
{ "type": "attack", "data": { "from_village": 123, "to_x": 50, "to_y": 75, "troops": {...} } }
{ "type": "chat", "data": { "channel": "alliance", "message": "Attack at 21:00!" } }
```

### Server → Client Messages

```json
{ "type": "connection_ready", "data": { "player_id": 42, "server_time": "..." } }
{ "type": "subscription_confirmed", "data": { "topics": ["village:123"] } }
{ "type": "village_state", "data": { "village_id": 123, "buildings": [...], "resources": {...} } }
{ "type": "resource_update", "data": { "village_id": 123, "iron": 5000, "wood": 3200, ... } }
{ "type": "build_started", "data": { "village_id": 123, "building_type": "barracks", "target_level": 2, "completes_at": "..." } }
{ "type": "build_complete", "data": { "village_id": 123, "building_type": "barracks", "new_level": 2 } }
{ "type": "train_started", "data": { "village_id": 123, "troop_type": "iron_legionary", "quantity": 10, "completes_at": "..." } }
{ "type": "train_complete", "data": { "village_id": 123, "troop_type": "iron_legionary", "quantity": 10 } }
{ "type": "attack_incoming", "data": { "village_id": 123, "arrives_at": "2026-03-03T15:30:00Z" } }
{ "type": "combat_result", "data": { "attack_id": 456, "winner": "attacker", ... } }
{ "type": "world_event", "data": { "event_type": "chaos_weapon_claimed", ... } }
{ "type": "error", "data": { "code": "INSUFFICIENT_RESOURCES", "message": "Not enough iron" } }
```

### Rate Limiting (WebSocket)

- Max 30 messages per second per connection
- Max 5 build/train/attack actions per second
- Chat: max 2 messages per second
- Violations: warning → temporary mute → disconnect

---

## Game Tick System

The server runs a **tick loop** that processes time-dependent game events.

### Tick Types

| Tick | Frequency | Purpose |
|------|----------|---------|
| **Resource Tick** | On-demand (lazy) | Resources are NOT ticked periodically. Calculated on read. Written on events. |
| **Building Tick** | Every 1 second | Check if any building construction has completed. |
| **Troop Movement Tick** | Every 5 seconds | Update troop positions on map. Check for arrivals. |
| **Combat Resolution** | On arrival | When attacking troops arrive, resolve combat immediately. |
| **Moraphys Tick** | Every 1 hour | Moraphys NPC faction grows stronger, may launch raids. |
| **World Event Tick** | Configurable | Random world events (Chaos storms, rune spawns, etc.) |

### Tick Architecture

```go
// Simplified tick loop
func (g *GameLoop) Run(ctx context.Context) {
    buildingTicker := time.NewTicker(1 * time.Second)
    troopTicker := time.NewTicker(5 * time.Second)
    moraphysTicker := time.NewTicker(1 * time.Hour)

    for {
        select {
        case <-ctx.Done():
            return
        case <-buildingTicker.C:
            g.processBuildingCompletions(ctx)
        case <-troopTicker.C:
            g.processTroopMovements(ctx)
        case <-moraphysTicker.C:
            g.processMoraphysTick(ctx)
        }
    }
}
```

---

## API Routes (REST)

### Auth

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/auth/register` | Register with email/password + kingdom |
| POST | `/api/auth/login` | Login with email/password |
| POST | `/api/auth/refresh` | Refresh access token |
| POST | `/api/auth/logout` | Invalidate refresh token |
| GET | `/api/auth/oauth/{provider}` | Initiate OAuth flow |
| GET | `/api/auth/oauth/{provider}/callback` | OAuth callback |

### Game

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/villages` | List player's villages |
| GET | `/api/villages/{id}` | Get village details (buildings, resources, troops, training queue) |
| PUT | `/api/villages/{id}` | Update village (e.g. rename) |
| POST | `/api/villages/{id}/build` | Start building construction |
| DELETE | `/api/villages/{id}/build/{queueId}` | Cancel queued build |
| POST | `/api/villages/{id}/train` | Start troop training (troop_type, quantity) |
| GET | `/api/villages/{id}/train/cost` | Preview training cost (?troop_type=X&quantity=N) |
| DELETE | `/api/villages/{id}/train/{queueId}` | Cancel queued training |
| GET | `/api/villages/{id}/troops` | List troops stationed in a village |
| GET | `/api/map?x={x}&y={y}&range={r}` | Get map tiles (default range 10 = 21×21 chunk, max 20) |
| GET | `/api/map/tile?x={x}&y={y}` | Get single map tile details |
| GET | `/api/player/me` | Get current player profile (used to refresh cached data) |
| PUT | `/api/player/kingdom` | Choose kingdom (post-registration, creates first village) |
| WS | `/ws` | WebSocket connection (JWT via `?token=` query param). Hub manages clients, topic subscriptions, single-session per player. Game loop pushes `build_complete` and `train_complete` events. |

### Admin

| Method | Path | Description |
|--------|------|-------------|
| GET | `/api/admin/players` | List all players (paginated) |
| PUT | `/api/admin/players/{id}/role` | Update player role |
| GET | `/api/admin/config` | Get all world config key-value pairs |
| PUT | `/api/admin/config/{key}` | Set a world config value |
| GET | `/api/admin/stats` | Get server statistics (player/village/tile counts) |
| GET | `/api/admin/announcements` | List active announcements |
| POST | `/api/admin/announcements` | Create an announcement |
| DELETE | `/api/admin/announcements/{id}` | Delete an announcement |
| GET | `/api/admin/assets` | List all game assets |
| POST | `/api/admin/assets` | Create a game asset |
| DELETE | `/api/admin/assets/{id}` | Delete a game asset |
| POST | `/api/admin/assets/{id}/sprite` | Upload sprite for a game asset |
| DELETE | `/api/admin/assets/{id}/sprite` | Remove sprite from a game asset |
| GET | `/api/admin/resource-buildings` | List resource building configs |
| GET | `/api/admin/resource-buildings/{id}` | Get single resource building config |
| PUT | `/api/admin/resource-buildings/{id}` | Update resource building config |
| POST | `/api/admin/resource-buildings/{id}/sprite` | Upload resource building sprite |
| DELETE | `/api/admin/resource-buildings/{id}/sprite` | Remove resource building sprite |
| PUT | `/api/admin/map/terrain` | Paint terrain on the live world map |
| PUT | `/api/admin/map/zones` | Paint kingdom zones on the live world map |

### Map Templates

| Method | Path | Description |
|--------|------|-------------|
| POST | `/api/admin/templates` | Create a new map template (name, size) |
| GET | `/api/admin/templates` | List all templates |
| GET | `/api/admin/templates/{name}` | Get template (full tile data) |
| DELETE | `/api/admin/templates/{name}` | Delete a template |
| PUT | `/api/admin/templates/{name}/terrain` | Paint terrain on a template |
| PUT | `/api/admin/templates/{name}/zones` | Paint zones on a template |
| POST | `/api/admin/templates/{name}/apply` | Apply template to the live world map (delete all + reinsert) |
| POST | `/api/admin/templates/{name}/resize` | Resize a template (preserves existing tiles within bounds) |
| GET | `/api/admin/templates/{name}/export` | Export template as JSON |

---

## Map Template System

Templates decouple map editing from the live game. Admins create, paint, and preview templates before applying them.

### Architecture

- **Storage**: Templates are stored as JSON files at `server/data/templates/{sanitized_name}.json`. Each file contains the template metadata (name, size) and the full tile grid.
- **TemplateRepository**: File-based repository (not SQLite). Handles JSON read/write with file locking.
- **TemplateService**: Business logic for template CRUD, terrain/zone painting, resize, and apply-to-game.
- **TemplateHandler**: HTTP endpoints for all template operations, admin-only.
- **Apply flow**: `ApplyTemplate` does a DELETE ALL on `world_map` then INSERT of all template tiles. This handles size mismatches between the template and the existing map.

### Template Model

```go
type MapTemplate struct {
    Name  string         `json:"name"`
    Size  int            `json:"size"`  // odd number 3-201
    Tiles []TemplateTile `json:"tiles"`
}

type TemplateTile struct {
    X           int    `json:"x"`
    Y           int    `json:"y"`
    TerrainType string `json:"terrain_type"`
    KingdomZone string `json:"kingdom_zone"`
}
```

---

## Game Asset System

The admin panel supports uploading sprite images for game entities (buildings, resources, troops, village markers).

- **Upload**: Sprites are stored in `server/uploads/sprites/` with hashed filenames.
- **GameAssetRepository**: Tracks asset metadata (id, category, display_name, default_icon, sprite_path).
- **Seed data**: Default game assets are seeded on first run (building icons, resource icons, village markers per kingdom).
- **Admin UI**: Asset management page allows uploading/removing sprites, creating/deleting asset entries.

---

## API Response Conventions

### Response Envelope

All REST endpoints use a unified JSON envelope:

**Success:**
```json
{
  "data": { "id": 123, "name": "My Village", ... },
  "error": null
}
```

**Error:**
```json
{
  "data": null,
  "error": {
    "code": "INSUFFICIENT_RESOURCES",
    "message": "Not enough iron to build Barracks (need 200, have 150)"
  }
}
```

**List with pagination:**
```json
{
  "data": [ ... ],
  "meta": { "page": 1, "limit": 50, "total": 200 },
  "error": null
}
```

### Pagination

Endpoints returning lists use query parameters:
- `?page=1&limit=50` (defaults: page=1, limit=50, max limit=100)
- Response always includes `meta` with `page`, `limit`, and `total`

### Error Code Catalog

| Code | HTTP Status | Description |
|------|------------|-------------|
| `INVALID_INPUT` | 400 | Malformed request or validation failure |
| `UNAUTHORIZED` | 401 | Missing or invalid authentication |
| `FORBIDDEN` | 403 | Authenticated but insufficient permissions |
| `NOT_FOUND` | 404 | Resource does not exist |
| `RATE_LIMITED` | 429 | Too many requests |
| `INSUFFICIENT_RESOURCES` | 422 | Not enough resources for the action |
| `BUILDING_IN_PROGRESS` | 422 | Building queue is full |
| `PREREQUISITES_NOT_MET` | 422 | Required building levels not reached |
| `QUEUE_FULL` | 422 | Training or building queue at capacity |
| `MAX_LEVEL_REACHED` | 422 | Building already at maximum level |
| `INTERNAL_ERROR` | 500 | Unexpected server error |

### API Versioning

All endpoints live under `/api/` with no version prefix. If breaking changes are needed in the future, a `/api/v2/` prefix will be introduced. The original `/api/` endpoints will be maintained for backward compatibility during migration.

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of system architecture |
| 2026-03-03 | Fixed golang-jwt to v5, added dto/ package to server folder structure |
| 2026-03-03 | Replaced gorilla/websocket with coder/websocket, added missing WS messages (connection_ready, build_started, train_started/complete, village_state, subscription_confirmed), added API response envelope, error code catalog, pagination convention, map chunk spec, API versioning note |
| 2026-03-07 | Added all actually-implemented REST routes (admin panel, player, map templates). Added Map Template System section (TemplateRepository, TemplateService, template model, apply flow). Added Game Asset System section. Marked WebSocket /ws as not yet implemented. |
| 2026-03-08 | WebSocket foundation implemented: Hub (single-session, topic pub/sub, broadcast), Client (read/write pumps, ping/pong), Handler (HTTP upgrade + JWT auth via ?token=), Messages (all type constants + data structs). Game loop notification bridge: CompleteBuilds returns BuildCompletionEvent slice, tick() broadcasts build_complete via hub. 8 WebSocket tests added. PlayerService extracted (Stage 3). 163 total backend tests. |
| 2026-03-08 | Troops & Training feature: Added 4 REST endpoints (train, cost, cancel, list troops). Game loop extended with TrainCompleter + TrainCompletionNotifier for train_complete WebSocket events. Village response now includes troops + training_queue. 18 new backend tests (12 service + 6 handler). Frontend: TrainingPanel, TrainingQueue, TroopRoster components integrated into VillagePage. Arkazia kingdom (7 troop types). |

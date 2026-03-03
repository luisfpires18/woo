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
│   │   │   ├── Button.module.css
│   │   │   └── Button.mobile.css
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
│   │   ├── themes.css          # Dark/light mode variables
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
| **gorilla/websocket** or **nhooyr/websocket** | WebSocket connections |
| **modernc.org/sqlite** | SQLite driver (pure Go, no CGO) |
| **golang-jwt/jwt** | JWT token generation and validation |
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
│   │   ├── hub.go              # Central WebSocket hub
│   │   ├── client.go           # Individual client connection
│   │   ├── messages.go         # Message type definitions
│   │   └── handlers.go         # WebSocket message handlers
│   └── gameloop/
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
{ "type": "resource_update", "data": { "village_id": 123, "iron": 5000, "wood": 3200, ... } }
{ "type": "build_complete", "data": { "village_id": 123, "building_type": "barracks", "new_level": 2 } }
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
| GET | `/api/villages/{id}` | Get village details (buildings, resources) |
| POST | `/api/villages/{id}/build` | Start building construction |
| POST | `/api/villages/{id}/train` | Start troop training |
| GET | `/api/map?x={x}&y={y}&radius={r}` | Get map tiles around coordinates |
| GET | `/api/player/profile` | Get current player profile |
| WS | `/ws` | WebSocket connection (with JWT auth) |

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of system architecture |

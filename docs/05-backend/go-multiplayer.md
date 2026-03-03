# Go Multiplayer & Concurrency

> WebSocket handling, game tick loop, real-time state management, and anti-cheat patterns. Read before implementing any multiplayer features.

> **WebSocket Library**: This project uses **`coder/websocket`** (maintained fork of `nhooyr/websocket`). It supports the standard `net/http` stack, context-based cancellation, and `io.Reader`/`io.Writer` for messages. Do NOT use `gorilla/websocket` (archived). The API examples in this doc use `coder/websocket` conventions.

---

## WebSocket Architecture

### Hub Pattern

A **central Hub** manages all active WebSocket connections. It handles:
- Registering new clients on connection
- Unregistering clients on disconnect
- Broadcasting messages to subscribed clients
- Routing messages to appropriate handlers

```go
// internal/websocket/hub.go
type Hub struct {
    clients    map[*Client]bool
    register   chan *Client
    unregister chan *Client
    broadcast  chan *Message

    // Topic-based subscriptions: topic → set of clients
    topics     map[string]map[*Client]bool
    subscribe  chan *Subscription
    unsubscribe chan *Subscription

    mu sync.RWMutex
}

func NewHub() *Hub {
    return &Hub{
        clients:     make(map[*Client]bool),
        register:    make(chan *Client),
        unregister:  make(chan *Client),
        broadcast:   make(chan *Message, 256),
        topics:      make(map[string]map[*Client]bool),
        subscribe:   make(chan *Subscription),
        unsubscribe: make(chan *Subscription),
    }
}

func (h *Hub) Run(ctx context.Context) {
    for {
        select {
        case <-ctx.Done():
            return
        case client := <-h.register:
            h.clients[client] = true
        case client := <-h.unregister:
            h.removeClient(client)
        case sub := <-h.subscribe:
            h.addToTopic(sub.Topic, sub.Client)
        case sub := <-h.unsubscribe:
            h.removeFromTopic(sub.Topic, sub.Client)
        case msg := <-h.broadcast:
            h.broadcastToTopic(msg)
        }
    }
}
```

### Client Connection

Each WebSocket connection is wrapped in a `Client` struct with dedicated read/write goroutines:

```go
// internal/websocket/client.go
type Client struct {
    hub      *Hub
    conn     *websocket.Conn
    send     chan []byte
    playerID int64
    topics   map[string]bool
}

func (c *Client) ReadPump() {
    defer func() {
        c.hub.unregister <- c
        c.conn.Close()
    }()

    c.conn.SetReadLimit(maxMessageSize)
    c.conn.SetReadDeadline(time.Now().Add(pongWait))
    c.conn.SetPongHandler(func(string) error {
        c.conn.SetReadDeadline(time.Now().Add(pongWait))
        return nil
    })

    for {
        _, message, err := c.conn.ReadMessage()
        if err != nil {
            break
        }
        c.hub.handleMessage(c, message)
    }
}

func (c *Client) WritePump() {
    ticker := time.NewTicker(pingPeriod)
    defer func() {
        ticker.Stop()
        c.conn.Close()
    }()

    for {
        select {
        case msg, ok := <-c.send:
            if !ok {
                c.conn.WriteMessage(websocket.CloseMessage, []byte{})
                return
            }
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            c.conn.WriteMessage(websocket.TextMessage, msg)
        case <-ticker.C:
            c.conn.SetWriteDeadline(time.Now().Add(writeWait))
            c.conn.WriteMessage(websocket.PingMessage, nil)
        }
    }
}
```

### Connection Lifecycle

```
1. Client connects to /ws?token=<JWT>
2. Server validates JWT → extracts playerID
3. Server creates Client struct, starts ReadPump + WritePump goroutines
4. Server registers client with Hub
5. Client sends subscribe messages → Hub adds client to topics
6. Server pushes events to client via topics
7. On disconnect: Hub unregisters client, cleans up subscriptions
```

---

## Topic-Based Subscriptions

Clients subscribe to **topics** to receive relevant events. This prevents broadcasting everything to everyone.

### Topic Format

| Topic | Description | Example |
|-------|------------|---------|
| `village:{id}` | Events for a specific village | `village:123` |
| `map:{x},{y}` | Events for a map region (chunk) | `map:5,10` |
| `alliance:{id}` | Alliance chat and events | `alliance:7` |
| `world` | Global server events | `world` |
| `player:{id}` | Personal events (attacks, messages) | `player:42` |

### Security

- A player can only subscribe to topics they have access to:
  - Their own villages
  - Their alliance
  - Map regions near their villages (fog of war)
  - Their own player events
- The Hub validates subscriptions against the player's permissions.

---

## Game Tick Loop

### Design

The game loop is server-authoritative. All game state changes happen through the tick loop or through validated player actions.

```go
// internal/gameloop/ticker.go
type GameLoop struct {
    villageService  service.VillageService
    resourceService service.ResourceService
    combatService   service.CombatService
    moraphysService service.MoraphysService
    hub             *websocket.Hub
    logger          *slog.Logger
}

func (g *GameLoop) Run(ctx context.Context) {
    buildingTicker := time.NewTicker(1 * time.Second)
    troopTicker := time.NewTicker(5 * time.Second)
    moraphysTicker := time.NewTicker(1 * time.Hour)

    defer buildingTicker.Stop()
    defer troopTicker.Stop()
    defer moraphysTicker.Stop()

    for {
        select {
        case <-ctx.Done():
            g.logger.Info("game loop shutting down")
            return
        case <-buildingTicker.C:
            g.processBuildingCompletions(ctx)
        case <-troopTicker.C:
            g.processTroopMovements(ctx)
            g.processCombatArrivals(ctx)
        case <-moraphysTicker.C:
            g.processMoraphysTick(ctx)
        }
    }
}
```

### Tick Responsibilities

| Tick | Frequency | What It Does |
|------|----------|--------------|
| Building | 1s | Check `building_queue` for completed constructions. Update building levels. Notify clients. |
| Troop Movement | 5s | Check `attacks` table for arrived troops. Trigger combat resolution or reinforcement placement. |
| Combat Resolution | On arrival | Calculate attack vs defense. Distribute losses. Update resources (raid). Notify both parties. |
| Moraphys | 1h | Grow Moraphys strength. Possibly launch NPC raids. Check for Weapon of Chaos theft attempts. |
| World Events | Configurable | Spawn rune fragments, trigger chaos storms, seasonal events. |

### Resource Tick (Lazy — Not a Real Tick)

Resources are **NOT** ticked periodically. They are calculated on demand:

```go
func (s *ResourceService) GetCurrentResources(ctx context.Context, villageID int64) (*model.Resources, error) {
    stored, err := s.repo.Get(ctx, villageID)
    if err != nil {
        return nil, fmt.Errorf("get stored resources for village %d: %w", villageID, err)
    }

    elapsed := time.Since(stored.LastUpdated).Hours()

    return &model.Resources{
        Iron:  min(stored.Iron + stored.IronRate*elapsed, stored.MaxStorage),
        Wood:  min(stored.Wood + stored.WoodRate*elapsed, stored.MaxStorage),
        Stone: min(stored.Stone + stored.StoneRate*elapsed, stored.MaxStorage),
        Food:  min(stored.Food + (stored.FoodRate-stored.FoodConsumption)*elapsed, stored.MaxStorage),
        LastUpdated: time.Now(),
    }, nil
}
```

Resources are written to DB only when an event consumes or changes them (building, training, trade, attack, login).

---

## Concurrency Safety

### Principles

1. **No global mutable state.** All state is managed through services with proper synchronization.
2. **Per-village locks.** Village operations lock at the village level, not globally.
3. **Channel-based communication.** The Hub uses channels for all inter-goroutine communication.
4. **Context cancellation.** All long-running operations respect `context.Context` for graceful shutdown.

### Per-Village Mutex Pattern

```go
type VillageLockManager struct {
    mu    sync.Mutex
    locks map[int64]*sync.RWMutex
}

func NewVillageLockManager() *VillageLockManager {
    return &VillageLockManager{
        locks: make(map[int64]*sync.RWMutex),
    }
}

func (m *VillageLockManager) getOrCreate(villageID int64) *sync.RWMutex {
    m.mu.Lock()
    defer m.mu.Unlock()

    if lock, ok := m.locks[villageID]; ok {
        return lock
    }
    lock := &sync.RWMutex{}
    m.locks[villageID] = lock
    return lock
}

func (m *VillageLockManager) WithLock(villageID int64, fn func()) {
    lock := m.getOrCreate(villageID)
    lock.Lock()
    defer lock.Unlock()
    fn()
}

func (m *VillageLockManager) WithRLock(villageID int64, fn func()) {
    lock := m.getOrCreate(villageID)
    lock.RLock()
    defer lock.RUnlock()
    fn()
}
```

### When to Lock

| Operation | Lock Type | Scope |
|----------|----------|-------|
| Read resources | RLock | Per-village |
| Build / upgrade | Write Lock | Per-village |
| Train troops | Write Lock | Per-village |
| Send attack | RLock (read troops) → Write Lock (deduct troops) | Per-village (source) |
| Combat resolution | Write Lock | Per-village (both attacker + defender) |

### Deadlock Prevention

When locking multiple villages (e.g., combat between village A and B):
- **Always lock in ascending order of village ID.**
- This prevents ABBA deadlocks.

```go
func lockTwoVillages(mgr *VillageLockManager, id1, id2 int64) (unlock func()) {
    if id1 > id2 {
        id1, id2 = id2, id1
    }
    lock1 := mgr.getOrCreate(id1)
    lock2 := mgr.getOrCreate(id2)
    lock1.Lock()
    lock2.Lock()
    return func() {
        lock2.Unlock()
        lock1.Unlock()
    }
}
```

---

## Server-Authoritative Design

### Core Principle

**The client sends intents. The server decides outcomes.**

```
Client: "I want to build a Barracks in village 123"
Server: Validates → checks resources → checks prerequisites → checks queue → executes or rejects
```

### What the Client CAN Do

- Send action intents (build, train, attack, trade, chat)
- Subscribe/unsubscribe to topics
- Request current state (REST API)
- Display UI and run display-only calculations (countdown timers, resource ETAs)

### What the Client CANNOT Do

- Calculate combat results
- Determine troop arrival times (server sends definitive timestamps)
- Modify game state directly
- Bypass cooldowns or queue limits

### Validation Checklist (Every Action)

1. ✅ Is the player authenticated? (JWT valid)
2. ✅ Does the player own the village? (authorization)
3. ✅ Are resources sufficient?
4. ✅ Are prerequisites met? (building levels, etc.)
5. ✅ Is the action allowed right now? (no concurrent builds, no cooldown violations)
6. ✅ Is the input sanitized? (no SQL injection, no XSS in names)
7. ✅ Is the rate limit respected?

---

## Anti-Cheat (Server-Side)

See also: `docs/08-security/security-and-anticheat.md`

### Timestamps

- **Never trust client timestamps.** All timing (build start, attack departure, etc.) uses `time.Now()` on the server.

### Rate Limiting (WebSocket)

```go
type RateLimiter struct {
    mu      sync.Mutex
    clients map[int64]*rate.Limiter
}

func (rl *RateLimiter) Allow(playerID int64) bool {
    rl.mu.Lock()
    limiter, ok := rl.clients[playerID]
    if !ok {
        limiter = rate.NewLimiter(rate.Limit(30), 30) // 30 msg/sec burst
        rl.clients[playerID] = limiter
    }
    rl.mu.Unlock()
    return limiter.Allow()
}
```

### Action-Specific Rate Limits

| Action | Limit | Window |
|--------|-------|--------|
| General WebSocket messages | 30/sec | Burst |
| Build / Upgrade | 5/sec | Burst |
| Train troops | 5/sec | Burst |
| Attack / Send troops | 2/sec | Burst |
| Chat message | 2/sec | Burst |
| Map tile requests | 20/sec | Burst |

### Anomaly Detection

Log and flag suspicious patterns:
- Building requests faster than the queue allows
- Resource spending that exceeds current balance
- Movement commands with impossible coordinates
- Repetitive identical requests (bot behavior)

---

## Graceful Shutdown

```go
func main() {
    ctx, cancel := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
    defer cancel()

    // Start game loop
    go gameLoop.Run(ctx)

    // Start HTTP server
    go func() {
        if err := server.ListenAndServe(); err != http.ErrServerClosed {
            slog.Error("server error", "error", err)
        }
    }()

    // Wait for shutdown signal
    <-ctx.Done()
    slog.Info("shutting down...")

    // Graceful shutdown with timeout
    shutdownCtx, shutdownCancel := context.WithTimeout(context.Background(), 30*time.Second)
    defer shutdownCancel()

    server.Shutdown(shutdownCtx)
    slog.Info("server stopped")
}
```

---

## WebSocket Reconnection

The server does **not** handle reconnection logic. If a client disconnects:
1. The Hub unregisters the client and cleans up all subscriptions.
2. It is the **client's responsibility** to reconnect (see frontend guide for exponential backoff strategy).
3. On reconnect, the client must re-authenticate (JWT in query param) and re-subscribe to all topics.
4. Any events missed during disconnection are NOT replayed. The client should fetch current state via REST API after reconnecting.

---

## Online Player Tracking

Track which players are currently online via their WebSocket connections:

```go
// Using sync.Map for concurrent read/write
type OnlineTracker struct {
    players sync.Map // map[int64]*Client
}

func (t *OnlineTracker) SetOnline(playerID int64, client *Client) {
    t.players.Store(playerID, client)
}

func (t *OnlineTracker) SetOffline(playerID int64) {
    t.players.Delete(playerID)
}

func (t *OnlineTracker) IsOnline(playerID int64) bool {
    _, ok := t.players.Load(playerID)
    return ok
}

func (t *OnlineTracker) OnlineCount() int {
    count := 0
    t.players.Range(func(_, _ any) bool {
        count++
        return true
    })
    return count
}
```

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of Go multiplayer guide |
| 2026-03-03 | Added coder/websocket library note, updated lazy resource calc with max_storage and food_consumption, added WebSocket reconnection policy |

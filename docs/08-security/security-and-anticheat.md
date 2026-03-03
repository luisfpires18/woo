# Security & Anti-Cheat Guide

> All security decisions, anti-cheat patterns, and protection strategies. Read before implementing any auth, input handling, or client-server communication.

---

## Core Principle: Server-Authoritative

The **single most important security rule**:

> **All game logic runs on the server. The client is a display terminal that sends intents.**

The client can:
- Display data
- Calculate display-only values (countdown timers, resource ETAs)
- Send action requests ("I want to build X")

The client **cannot**:
- Determine combat outcomes
- Set resource values
- Bypass build timers
- Trust its own timestamps

---

## Authentication Security

### Password Handling

| Aspect | Implementation |
|--------|---------------|
| **Hashing** | bcrypt with cost factor 12 |
| **Salt** | Built into bcrypt (automatic per-hash salt) |
| **Minimum length** | 8 characters |
| **Complexity** | At least 1 uppercase, 1 lowercase, 1 digit |
| **Storage** | Only hash stored in DB. Plain password never logged or stored. |

```go
import "golang.org/x/crypto/bcrypt"

func HashPassword(password string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(password), 12)
    return string(bytes), err
}

func CheckPassword(password, hash string) bool {
    err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
    return err == nil
}
```

### JWT Security

| Token | Lifetime | Storage | Purpose |
|-------|----------|---------|---------|
| **Access Token** | 15 minutes | Client memory (Zustand) | API authentication |
| **Refresh Token** | 7 days | HTTP-only secure cookie | Silent token refresh |

**Access Token (JWT)**:
- Signed with HMAC-SHA256 (`HS256`)
- Contains: `player_id`, `kingdom`, `iat`, `exp`
- **Short-lived** (15 min) to limit damage from token theft
- Sent in `Authorization: Bearer <token>` header

**Refresh Token**:
- Opaque UUID, **not a JWT**
- Stored as SHA-256 hash in the database
- Sent via **HTTP-only, Secure, SameSite=Strict** cookie
- When used, **rotated**: old token invalidated, new token issued
- Revokable by deleting from the database

```go
// Refresh token rotation
func (s *AuthService) RefreshToken(ctx context.Context, oldToken string) (*AuthResponse, error) {
    tokenHash := sha256Hash(oldToken)

    stored, err := s.tokenRepo.GetByHash(ctx, tokenHash)
    if err != nil {
        return nil, model.ErrUnauthorized // Token not found = already used or invalid
    }

    if time.Now().After(stored.ExpiresAt) {
        s.tokenRepo.Delete(ctx, stored.ID)
        return nil, model.ErrUnauthorized
    }

    // Delete old token (rotation)
    s.tokenRepo.Delete(ctx, stored.ID)

    // Issue new tokens
    return s.issueTokens(ctx, stored.PlayerID)
}
```

### OAuth Security

- Validate `state` parameter on callback to prevent CSRF
- Exchange authorization code server-side (never expose client secret to browser)
- Validate ID token signature and claims
- Link OAuth accounts to existing email accounts when possible

---

## CORS Configuration

```go
func CORS(allowedOrigins []string) func(http.Handler) http.Handler {
    return func(next http.Handler) http.Handler {
        return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
            origin := r.Header.Get("Origin")
            for _, allowed := range allowedOrigins {
                if origin == allowed {
                    w.Header().Set("Access-Control-Allow-Origin", origin)
                    break
                }
            }
            w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
            w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")
            w.Header().Set("Access-Control-Allow-Credentials", "true")
            w.Header().Set("Access-Control-Max-Age", "86400") // 24h preflight cache

            if r.Method == "OPTIONS" {
                w.WriteHeader(http.StatusNoContent)
                return
            }
            next.ServeHTTP(w, r)
        })
    }
}
```

**Rules**:
- **Whitelist origins** — no wildcards (`*`) in production
- **Allow credentials** — needed for refresh token cookies
- **Preflight caching** — reduce OPTIONS requests

---

## Input Validation & Sanitization

### Server-Side Validation (Mandatory)

Every input is validated on the server, regardless of client-side validation:

```go
func validateUsername(username string) error {
    if len(username) < 3 || len(username) > 20 {
        return fmt.Errorf("username must be 3-20 characters: %w", model.ErrInvalidInput)
    }
    if !regexp.MustCompile(`^[a-zA-Z0-9_]+$`).MatchString(username) {
        return fmt.Errorf("username can only contain letters, numbers, and underscores: %w", model.ErrInvalidInput)
    }
    return nil
}

func validateVillageName(name string) error {
    if len(name) < 1 || len(name) > 30 {
        return fmt.Errorf("village name must be 1-30 characters: %w", model.ErrInvalidInput)
    }
    // Sanitize HTML/script tags
    name = sanitizeHTML(name)
    return nil
}
```

### XSS Prevention

- **Sanitize all user-generated text** before storage (village names, alliance names, chat messages)
- Use an HTML sanitizer library (e.g., `bluemonday` for Go)
- React auto-escapes JSX output, but **never use `dangerouslySetInnerHTML`** with user content

### SQL Injection Prevention

- **Parameterized queries ONLY** (already mandated in database guide)
- Never concatenate user input into SQL strings
- No `fmt.Sprintf` for building SQL queries

---

## Rate Limiting

### HTTP Rate Limiting

```go
type RateLimitMiddleware struct {
    limiters map[string]*rate.Limiter
    mu       sync.Mutex
    rps      int // requests per second
}

func (m *RateLimitMiddleware) Handler(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        ip := extractIP(r)
        limiter := m.getLimiter(ip)

        if !limiter.Allow() {
            w.Header().Set("Retry-After", "1")
            writeError(w, http.StatusTooManyRequests, "rate limit exceeded")
            return
        }
        next.ServeHTTP(w, r)
    })
}
```

### WebSocket Rate Limiting

Applied per-connection and per-action:

| Action Type | Limit | Enforcement |
|------------|-------|-------------|
| Any message | 30/sec | Per-connection |
| Build/Upgrade | 5/sec | Per-player |
| Train troops | 5/sec | Per-player |
| Send attack | 2/sec | Per-player |
| Chat | 2/sec | Per-player |
| Map request | 20/sec | Per-player |

### Violation Escalation

1. **Warn**: First rate limit hit → server sends `{ type: "warning", data: { message: "slow down" } }`
2. **Temporary mute**: 3 warnings in 1 minute → 30 second message suppression
3. **Disconnect**: Continued violations → WebSocket closed with code 4429
4. **Temporary ban**: Persistent abuse → 15-minute IP/account ban
5. **Permanent ban**: Confirmed bot/cheat → account flagged for review

---

## Anti-Cheat: Browser-Side

While the server is the authority, we add browser-side measures to raise the bar for casual cheaters:

### 1. WebSocket Message Integrity

- Messages use a **sequence number** that the server tracks. Out-of-order or duplicate messages are flagged.
- Server tracks message timestamps — if a player sends a command referencing a resource value the server hasn't sent yet, it's suspicious.

### 2. Action Timing Validation

- Server tracks time between player actions. Inhuman speed (e.g., 200ms between complex multi-step actions) is flagged.
- Building/training commands include the client's displayed resource value. If it differs significantly from the server's calculation, flag for review.

### 3. Behavioral Analysis

Flag accounts that exhibit:
- **Identical timing between actions** (bot-like regularity)
- **24/7 activity without breaks** (bot or automation)
- **Perfect optimization** (always training the exact right troops at the exact right time with zero waste)
- **Impossible travel times** (troops arriving faster than physics allows)

### 4. Client-Side Hardening

- **Obfuscated build** (Vite minification + tree shaking)
- **No sensitive data in client**: no damage formulas, no loot tables, no NPC stats
- **WebSocket reconnection limits**: Max 5 reconnections per minute

---

## Anti-Cheat: Server-Side

### Validation Checklist for Every Game Action

```go
func (s *VillageService) StartBuild(ctx context.Context, playerID int64, villageID int64, buildingType string, targetLevel int) error {
    // 1. Ownership check
    village, err := s.villageRepo.GetByID(ctx, villageID)
    if err != nil { return err }
    if village.PlayerID != playerID {
        return model.ErrForbidden // Not your village
    }

    // 2. Building exists and target level is valid
    building, err := s.buildingRepo.GetByType(ctx, villageID, buildingType)
    if err != nil { return err }
    if targetLevel != building.Level + 1 {
        return model.ErrInvalidInput // Can only upgrade one level at a time
    }

    // 3. Prerequisites check
    if err := s.checkPrerequisites(ctx, villageID, buildingType, targetLevel); err != nil {
        return err
    }

    // 4. Queue check (no concurrent builds)
    if s.hasPendingBuild(ctx, villageID) {
        return model.ErrBuildingInProgress
    }

    // 5. Resource check (calculate current, then deduct)
    resources, err := s.resourceService.GetCurrent(ctx, villageID)
    if err != nil { return err }
    cost := s.calculateBuildCost(buildingType, targetLevel)
    if !resources.CanAfford(cost) {
        return model.ErrInsufficientRes
    }

    // 6. Deduct resources and start build
    // ... (all validated, safe to proceed)
}
```

### Server Timestamp Authority

```go
// GOOD — server determines timing
buildQueue := &model.BuildQueue{
    VillageID:    villageID,
    BuildingType: buildingType,
    TargetLevel:  targetLevel,
    StartedAt:    time.Now().UTC(),                    // Server time
    CompletesAt:  time.Now().UTC().Add(buildDuration), // Server time
}

// BAD — trusting client timestamp
// startedAt := request.Body.StartedAt  // NEVER DO THIS
```

---

## Sensitive Data Protection

### What to Never Log

- Passwords (plain or hashed)
- JWT tokens (full)
- Refresh tokens
- OAuth secrets
- Full email addresses in general logs (use masked: `t***@example.com`)

### What to Never Send to Client

- Other players' exact troop counts (unless scouted)
- Combat formula details (client only sees results)
- Exact NPC stats (discovery-based)
- Other players' exact resource values
- Server-side configuration (rate limits, detection thresholds)

### Environment Variables

```
# .env (NEVER commit to git)
JWT_SECRET=your-secret-here
GOOGLE_CLIENT_SECRET=...
DISCORD_CLIENT_SECRET=...
DB_PATH=./data/woo.db
```

Add to `.gitignore`:
```
.env
*.db
data/
```

---

## Security Headers

Set these on all HTTP responses:

```go
func SecurityHeaders(next http.Handler) http.Handler {
    return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("X-Content-Type-Options", "nosniff")
        w.Header().Set("X-Frame-Options", "DENY")
        w.Header().Set("X-XSS-Protection", "0") // Modern browsers use CSP instead
        w.Header().Set("Referrer-Policy", "strict-origin-when-cross-origin")
        w.Header().Set("Content-Security-Policy",
            "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline' https://fonts.googleapis.com; font-src 'self' https://fonts.gstatic.com; connect-src 'self' wss://*")
        next.ServeHTTP(w, r)
    })
}
```

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of security and anti-cheat guide |

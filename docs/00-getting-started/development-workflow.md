# Development Workflow

> How agents and humans collaborate on this project.

---

## Git Branching Strategy

```
master ─────────────────────────────────── (production-ready)
  │
  └── develop ──────────────────────────── (integration branch)
        │
        ├── feature/auth-system ────────── (new features)
        ├── feature/world-map
        ├── fix/resource-calc-bug ──────── (bug fixes)
        └── docs/update-kingdom-lore ───── (documentation changes)
```

### Rules

- **`master`**: Production-ready code only. Never push directly.
- **`develop`**: Integration branch. All feature branches merge here first.
- **`feature/*`**: New features. Branch from `develop`, merge back to `develop`.
- **`fix/*`**: Bug fixes. Branch from `develop` (or `master` for hotfixes).
- **`docs/*`**: Documentation-only changes. Branch from `develop`.
- **Commit messages**: Conventional commits — `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`

---

## Development Environment Setup

### Prerequisites

| Tool | Version | Notes |
|------|---------|-------|
| **Node.js** | 20 LTS | Use `.nvmrc` at repo root (`20`) for nvm users |
| **Go** | 1.25+ | Required for the backend server |
| **golangci-lint** | latest | Install via `go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest` |
| **Git** | 2.x+ | Standard version control |

### Ports

| Service | Port | URL |
|---------|------|-----|
| Frontend (Vite dev server) | 5173 | `http://localhost:5173` |
| Backend (Go server) | 8080 | `http://localhost:8080` |

### Vite Proxy

The Vite dev server proxies API and WebSocket requests to the Go backend. Configure in `client/vite.config.ts`:

```ts
export default defineConfig({
  server: {
    proxy: {
      '/api': 'http://localhost:8080',
      '/ws': {
        target: 'http://localhost:8080',
        ws: true,
      },
    },
  },
  // ...
});
```

### Environment Variables

Copy `server/.env.example` to `server/.env` and fill in secrets. See `docs/05-backend/go-guide.md` for the full list of variables.

### How to Start

```bash
# Terminal 1 — Backend
cd server
cp .env.example .env      # first time only
go build -o server.exe ./cmd/server
.\server.exe

# Terminal 2 — Frontend
cd client
npm install               # first time only
npm run dev
```

### Dev Server Workflow

After completing any development task:

1. **Kill the running server**: `Get-Process -Name "server" -ErrorAction SilentlyContinue | Stop-Process -Force`
2. **If migrations were added or changed**: Delete the database first: `Remove-Item "server/data/woo.db*" -Force -ErrorAction SilentlyContinue`
3. **Rebuild and start the server**: `cd server; go build -o server.exe ./cmd/server; .\server.exe` (run in background)
4. **Start the client dev server**: `cd client; npm run dev` (run in background)

### Config Codegen

After changing any Go config file in `server/internal/config/`:

```bash
npm run gen-config          # from repo root
# or
cd server && go run ./cmd/genconfig
```

This regenerates `client/src/config/generated/*.json`. Commit the updated JSON files.

### Sprite Manifest Guideline

The server keeps `server/uploads/sprites/sprites.txt` auto-generated as a sprite guideline list.

- Format: `<path> | <width>x<height>`
- Example: `uploads/sprites/resources/food.png | 256x256`
- Scope: includes image sprites under `server/uploads/sprites/` (excluding `sprites.txt` itself)

`sprites.txt` is refreshed on server startup and on sprite endpoints, so adding/replacing sprite files is reflected automatically.

### .nvmrc

A `.nvmrc` file at the repo root pins the Node.js version:

```
20
```

Run `nvm use` to switch to the correct version.

---

## Agent Collaboration Protocol

Every AI agent working on this project **MUST** follow this workflow:

### 1. Read Phase (Before Writing Code)

```
Read .github/copilot-instructions.md
  → Read docs/00-getting-started/project-overview.md
  → Read ALL docs in relevant folder(s) for the task
  → Cross-reference lore docs if implementing game content
```

### 2. Plan Phase

- Identify what needs to change
- Check if existing patterns/components can be reused
- Verify naming conventions match the guides

### 3. Implement Phase

- Write code following the conventions in the relevant guide docs
- Include unit tests for all new functions/components
- Include integration tests for multiplayer/real-time features
- Use proper error handling, logging, and input validation

### 4. Update Phase (After Writing Code)

- Update any doc that was affected by the change
- Add a changelog entry with date and description
- Ensure cross-references between docs are still accurate

---

## Development Phases

### Phase 1 — Foundation ✅
**Goal**: Authentication + World Map + Basic Village

- [x] Complete all documentation
- [x] Scaffold monorepo (Vite + React, Go module, SQLite)
- [x] Implement JWT auth (register/login/refresh/logout)
- [x] Create tile-based world map (rendering + navigation)
- [x] Implement basic village view (buildings, resource display)
- [x] Resource production system (lazy calculation)
- [x] Building construction queue
- [x] Kingdom selection at registration
- [ ] OAuth (Google/Discord) — not yet implemented

### Phase 2 — Economy (Partial) ✅
**Goal**: Full resource system + building upgrades + marketplace

- [x] All resource buildings with upgrade tiers (12 resource fields, 3 per resource)
- [x] Storage buildings (Storage for lumber/stone, Provisions for food, Reservoir for water)
- [x] Resource overflow protection (capped at max_storage)
- [x] Building upgrade time scaling (exponential via config)
- [x] Admin-configurable resource building display names per kingdom
- [x] Config codegen pipeline (Go → JSON → TypeScript)
- [ ] Marketplace (player-to-player trading) — not yet implemented
- [ ] Multiple village support (founding/conquering) — not yet implemented

### Phase 3 — Military ✅
**Goal**: Troops + combat + kingdom warfare

- [x] Troop types per kingdom — 140 troops across 7 playable kingdoms (~20 per kingdom)
- [x] Troop training system (Travian-style one-at-a-time queue, speed multiplier by building level)
- [x] 5 military buildings: Barracks, Stable, Archery, Workshop, Special
- [x] Frontend: BuildingTrainingModal, TrainingQueue, TroopRoster components
- [x] Game loop: train_complete WebSocket events
- [x] Resource economy centralized in config/resources.go with codegen
- [ ] Troop movement on world map (time-based) — not yet implemented
- [ ] Combat resolution system — not yet implemented
- [ ] Raiding, wall defense, scouting/fog of war — not yet implemented

### Phase 4 — Crafting
**Goal**: Runes + Forges + Weapon crafting

- [ ] Forge building with upgrade tiers
- [ ] Rune discovery system (exploration, drops, trading)
- [ ] Weapon crafting (resource + rune requirements)
- [ ] Weapon tiers: Common → Rare → Epic → Legendary → Mythic
- [ ] Rune combination / enhancement
- [ ] Weapon equipping on troops/heroes

### Phase 5 — Endgame
**Goal**: Weapons of Chaos + Moraphys + Weapons of Order

- [ ] Weapons of Chaos spawn system on world map
- [ ] Chaos wielder debuff mechanics
- [ ] Moraphys NPC faction AI (gradual power growth)
- [ ] Endgame trigger: Moraphys collects all Weapons of Chaos
- [ ] Alliance-level Weapons of Order crafting
- [ ] Final battle mechanics against Moraphys
- [ ] Victory conditions + round reset

### Phase 6 — Lore Explorer
**Goal**: Single-player interactive lore experience

- [ ] Lore explorer UI (separate from multiplayer)
- [ ] Chapter/scene progression system
- [ ] Character interactions and dialogue
- [ ] World history exploration
- [ ] Lore achievements/unlocks that carry into multiplayer (cosmetics only)

### Phase 7 — Polish & Scale
**Goal**: Mobile optimization + performance + launch readiness

- [ ] Mobile CSS for all components
- [ ] Performance optimization (lazy loading, code splitting)
- [ ] Load testing (concurrent players, WebSocket throughput)
- [ ] PostgreSQL migration for production
- [ ] Deployment pipeline (Docker, CI/CD)
- [ ] Monitoring and logging infrastructure
- [ ] Beta testing and balance adjustments

---

## How to Start a New Feature

1. **Create a feature branch**: `git checkout -b feature/feature-name develop`
2. **Read relevant docs** (see Agent Collaboration Protocol above)
3. **Implement with tests**: No code without tests
4. **Run linting**: `gofmt` + `golangci-lint` (backend), ESLint + TypeScript strict (frontend)
5. **Run config codegen**: If any Go config changed, run `npm run gen-config` and commit updated JSON
6. **Update docs** if conventions or schemas changed
7. **Create PR** to `develop` with description referencing the relevant doc section

---

## Definition of Done

A feature is "done" when:

- [ ] Code is implemented and follows all conventions
- [ ] Unit tests pass (80%+ coverage for services)
- [ ] Integration tests pass (if applicable)
- [ ] No linting errors
- [ ] Config parity tests pass (if config changed)
- [ ] Relevant docs are updated
- [ ] PR is reviewed and approved
- [ ] Merged to `develop`

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-11 | Added Sprite Manifest guideline: `server/uploads/sprites/sprites.txt` auto-generated with sprite paths and dimensions. |
| 2026-03-03 | Initial creation of development workflow |
| 2026-03-03 | Added Development Environment Setup section (Node 20 LTS, ports 5173/8080, Vite proxy, .nvmrc, startup instructions) |
| 2026-03-10 | Major update: Go version 1.25+, Phase 1-3 marked complete, added config codegen workflow, dev server workflow, Phase 3 detailed with 140 troops across 7 kingdoms |

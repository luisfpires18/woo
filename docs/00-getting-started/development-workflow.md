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
| **Go** | 1.22+ | Required for the backend server |
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
make run

# Terminal 2 — Frontend
cd client
npm install               # first time only
npm run dev
```

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

### Phase 1 — Foundation (Current)
**Goal**: Authentication + World Map + Basic Village

- [ ] Complete all documentation
- [ ] Scaffold monorepo (Vite + React, Go module, SQLite)
- [ ] Implement JWT auth + OAuth (Google/Discord)
- [ ] Create tile-based world map (rendering + navigation)
- [ ] Implement basic village view (buildings, resource display)
- [ ] Resource production system (lazy calculation)
- [ ] Building construction queue
- [ ] Kingdom selection at registration

### Phase 2 — Economy
**Goal**: Full resource system + building upgrades + marketplace

- [ ] All resource buildings with upgrade tiers
- [ ] Warehouse capacity limits
- [ ] Resource overflow protection
- [ ] Marketplace (player-to-player trading)
- [ ] Building upgrade time scaling
- [ ] Multiple village support (founding/conquering)

### Phase 3 — Military
**Goal**: Troops + combat + kingdom warfare

- [ ] Troop types per kingdom (recruitment, stats)
- [ ] Troop movement on world map (time-based)
- [ ] Combat resolution system (attack vs defense calculations)
- [ ] Raiding (steal resources from enemy villages)
- [ ] Wall defense system
- [ ] Scouting / fog of war

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
5. **Update docs** if conventions or schemas changed
6. **Create PR** to `develop` with description referencing the relevant doc section

---

## Definition of Done

A feature is "done" when:

- [ ] Code is implemented and follows all conventions
- [ ] Unit tests pass (80%+ coverage for services)
- [ ] Integration tests pass (if applicable)
- [ ] No linting errors
- [ ] Relevant docs are updated
- [ ] PR is reviewed and approved
- [ ] Merged to `develop`

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of development workflow |\n| 2026-03-03 | Added Development Environment Setup section (Node 20 LTS, ports 5173/8080, Vite proxy, .nvmrc, startup instructions) |

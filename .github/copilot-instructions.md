# Copilot Instructions — Weapons of Order (WOO)

> **These instructions are MANDATORY for all AI agents working on this project.**
> Read this file in full before starting any task.

---

## 1. Documentation-First Workflow

### Before ANY Task

1. **Always read** `docs/00-getting-started/project-overview.md` first to understand the project vision and current state.
2. **Read all docs** in the folder(s) relevant to your task:
   - Frontend work → read everything in `docs/04-frontend/`
   - Backend work → read everything in `docs/05-backend/`
   - Game mechanics / features → read `docs/01-game-design/` and `docs/02-lore/`
   - Database work → read `docs/06-database/`
   - Testing → read `docs/07-testing/`
   - Security / anti-cheat → read `docs/08-security/`
   - Architecture decisions → read `docs/03-architecture/`
3. **Cross-reference lore**: When implementing any game content (kingdoms, troops, weapons, runes, forges), also read `docs/02-lore/` to ensure lore consistency.

### After ANY Task

1. **Update the relevant doc(s)** if your work changed conventions, patterns, schemas, or introduced new systems.
2. **Add a changelog entry** at the bottom of the affected doc with date and brief description of what changed.

### Before Creating a New Doc

1. **Search ALL existing docs** to confirm the topic is not already covered.
2. Only create a new file if the content truly does not fit in any existing document.
3. If in doubt, add a new section to an existing doc rather than creating a new file.

---

## 2. Code Conventions

### General

- **No code without tests.** Every feature must include unit tests. Multiplayer/real-time features need integration tests.
- **No hardcoded values.** Use constants, environment variables, or config files.
- **No console.log / fmt.Println in production code.** Use proper logging libraries.

### Frontend (React + TypeScript)

- Follow `docs/04-frontend/frontend-guide.md` strictly.
- **Component reuse is mandatory.** If a UI element is used in 2+ places, extract it to `client/src/components/`.
- Every component gets a `.module.css` file with mobile overrides via `@media (max-width: 768px)` inside the same file. See `docs/04-frontend/styling-guide.md`.
- Use TypeScript strict mode. All API responses must be typed.
- State management: Zustand for global game state, React Query for server data.
- **Web-first design.** Desktop layout first, then mobile adaptation.

### Backend (Go)

- Follow `docs/05-backend/go-guide.md` strictly.
- Use `internal/` package for all business logic. No exported packages except `cmd/`.
- Clean architecture: handler → service → repository. No skipping layers.
- Error wrapping: `fmt.Errorf("context: %w", err)` always.
- Context propagation: every function that does I/O takes `context.Context` as first parameter.
- Run `gofmt` and `golangci-lint` before committing.

### Database

- Follow `docs/06-database/database-guide.md` strictly.
- Repository pattern with interfaces. Business logic never touches SQL directly.
- Parameterized queries only. Never concatenate SQL strings.
- Migrations are numbered sequentially: `001_create_players.sql`, `002_create_villages.sql`, etc.

### CSS / Styling

- Follow `docs/04-frontend/styling-guide.md` strictly.
- CSS Modules for component scoping.
- CSS custom properties (variables) for theming in `:root`.
- Fonts: `Cinzel` for headings, `EB Garamond` for body text.
- Dark mode: black primary + crimson accent. Light mode: white primary + navy blue accent.
- Every `.module.css` contains responsive overrides via `@media (max-width: 768px)` inside the same file (no separate `.mobile.css`).

---

## 3. File Naming Conventions

| Type | Convention | Example |
|------|-----------|---------|
| React components | PascalCase | `VillagePanel.tsx` |
| React component styles | PascalCase + module | `VillagePanel.module.css` |
| Go files | snake_case | `village_handler.go` |
| Go test files | snake_case + _test | `village_handler_test.go` |
| SQL migrations | numbered + snake_case | `001_create_players.sql` |
| Documentation | kebab-case | `core-mechanics.md` |
| TypeScript types | PascalCase | `VillageTypes.ts` |
| Hooks | camelCase with `use` prefix | `useVillageData.ts` |
| Stores | camelCase with `Store` suffix | `villageStore.ts` |

---

## 4. Git Workflow

- **Branches**: `master` (production), `develop` (integration), `feature/*`, `fix/*`, `docs/*`
- **Never push directly to `master` or `develop`.** Always use feature branches.
- **Commit messages**: Use conventional commits: `feat:`, `fix:`, `docs:`, `refactor:`, `test:`, `chore:`
- **PR descriptions**: Reference the relevant doc section and explain what was changed.

---

## 5. Architecture Rules

- **Server-authoritative**: All game logic runs on the server. The client sends intents, the server validates and executes.
- **Never trust the client**: No client-side game state calculations that affect gameplay. Display-only calculations are fine.
- **WebSocket messages**: JSON format with a `type` field. Defined in `docs/03-architecture/system-architecture.md`.
- **Database abstraction**: Repository interfaces in Go. SQLite now, PostgreSQL later. Zero business logic changes on DB swap.

---

## 6. Lore & Game Content Rules

- The three playable kingdoms are: **Veridor** (naval/ocean), **Sylvara** (forest/jungle), **Arkazia** (mountain/gladiator).
- Additional playable kingdoms: **Draxys** (desert/gladiator), **Nordalh** (frost/viking), **Zandres** (underground/crystal), **Lumus** (light/holy).
- **Drakanith** is the 8th kingdom (dragon/volcanic) — currently NPC-only, no playable units.
- The enemy faction **Moraphys** is NPC-controlled and triggers the endgame by stealing all Weapons of Chaos.
- **Weapons of Chaos** cause debuffs and chaos to their wielders. They are powerful but dangerous. **Count is configurable per game world** — do not hardcode 7.
- **Weapons of Order** are the counter — crafted by players through alliance-level collaboration to defeat Moraphys.
- The four base resources are: **Food**, **Water**, **Lumber**, **Stone**.
- Each resource has 3 building slots per village (e.g. food_1/food_2/food_3). Display names are admin-configurable per kingdom via `resource_building_configs` table.
- Always cross-check `docs/02-lore/` and `docs/01-game-design/kingdoms.md` when implementing kingdom-specific content.

---

## 7. Performance & Security Rules

- Follow `docs/08-security/security-and-anticheat.md` for all security decisions.
- Rate limit all WebSocket messages and REST endpoints.
- Validate all inputs server-side. Sanitize all user-generated text.
- Use lazy resource calculation (calculate on read, write on events) per `docs/06-database/database-guide.md`.

---

## 8. Config Codegen Pipeline (Go → JSON → TypeScript)

Game configuration (buildings, troops, resource economy) follows a **single-source-of-truth** pipeline:

### Source of Truth

- **Go config files** in `server/internal/config/` are the authoritative source:
  - `buildings.go` — building configs (costs, scaling, prerequisites)
  - `troops.go` — troop configs (stats, costs, kingdoms)
  - `resources.go` — resource economy constants (starting values, rates, storage)

### Pipeline

1. **Edit Go config** → make changes in `server/internal/config/*.go`
2. **Run codegen** → `npm run gen-config` (from repo root) or `cd server && go run ./cmd/genconfig`
3. **Generated JSON** → committed files in `client/src/config/generated/`:
   - `buildings.json`, `troops.json`, `resources.json`
4. **TypeScript imports** → frontend reads from generated JSON:
   - `client/src/config/buildings.ts` → imports `buildings.json`
   - `client/src/config/troops.ts` → imports `troops.json`
   - `client/src/config/resources.ts` → imports `resources.json`

### Rules

- **Never edit generated JSON files directly.** Always edit the Go source and re-run genconfig.
- **Never duplicate config values** between Go and TypeScript. The TS side must import from generated JSON.
- **Parity tests** in `server/internal/config/parity_test.go` verify committed JSON matches Go config. Run `go test ./internal/config/` to check.
- **After changing any Go config**, always run `npm run gen-config` and commit the updated JSON files.
- **DTO types** in `server/internal/config/generated_types.go` define the JSON schema shared between genconfig and the parity test.

---

## 9. Dev Server Workflow

After completing any development task:

1. **Kill the running server**: `Get-Process -Name "server" -ErrorAction SilentlyContinue | Stop-Process -Force`
2. **If migrations were added or changed**: Delete the database first: `Remove-Item "d:\Workspace\WOO\server\data\woo.db*" -Force -ErrorAction SilentlyContinue`
3. **Rebuild and start the server**: `cd d:\Workspace\WOO\server; go build -o server.exe ./cmd/server; .\server.exe` (run in background)
4. **Start the client dev server**: `cd d:\Workspace\WOO\client; npm run dev` (run in background)

This ensures the latest code is always running and testable after every change.

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of copilot instructions |
| 2026-03-03 | Removed .mobile.css convention (responsive overrides in .module.css), removed mobile styles from file naming table, added configurable Weapons of Chaos count note |
| 2026-03-03 | Resources refactored: Iron/Wood/Stone/Food → Food/Water/Lumber/Stone. 4 resource buildings → 12 (3 per resource). Added resource_building_configs table for admin customisation per kingdom |
| 2026-03-07 | Added Section 8: Dev Server Workflow — mandatory kill/rebuild/restart after every dev task; delete woo.db when migrations change |
| 2026-03-09 | Added Section 8: Config Codegen Pipeline (Go → JSON → TS) for buildings, troops, and resource economy. Updated kingdom list to include all 8 kingdoms |

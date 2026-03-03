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
- Every component gets a `.module.css` + `.mobile.css` pair. See `docs/04-frontend/styling-guide.md`.
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
- Every `.module.css` has a sibling `.mobile.css` loaded via `@media (max-width: 768px)`.

---

## 3. File Naming Conventions

| Type | Convention | Example |
|------|-----------|---------|
| React components | PascalCase | `VillagePanel.tsx` |
| React component styles | PascalCase + module | `VillagePanel.module.css` |
| React mobile styles | PascalCase + mobile | `VillagePanel.mobile.css` |
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
- The enemy faction **Moraphys** is NPC-controlled and triggers the endgame by stealing all Weapons of Chaos.
- **Weapons of Chaos** cause debuffs and chaos to their wielders. They are powerful but dangerous.
- **Weapons of Order** are the counter — crafted by players through alliance-level collaboration to defeat Moraphys.
- The four base resources are: **Iron**, **Wood**, **Stone**, **Food**.
- Always cross-check `docs/02-lore/` and `docs/01-game-design/kingdoms.md` when implementing kingdom-specific content.

---

## 7. Performance & Security Rules

- Follow `docs/08-security/security-and-anticheat.md` for all security decisions.
- Rate limit all WebSocket messages and REST endpoints.
- Validate all inputs server-side. Sanitize all user-generated text.
- Use lazy resource calculation (calculate on read, write on events) per `docs/06-database/database-guide.md`.

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of copilot instructions |

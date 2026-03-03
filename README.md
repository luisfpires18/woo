# Weapons of Order (WOO)

> A multiplayer browser-based RTS game with dark fantasy lore, inspired by Travian.

---

## Quick Start

> Code scaffold not yet created. Currently in Documentation & Planning phase.

---

## Documentation

### Getting Started
- [Project Overview](docs/00-getting-started/project-overview.md) — **Read this first.** Vision, tech stack, MVP scope.
- [Development Workflow](docs/00-getting-started/development-workflow.md) — Git workflow, agent protocol, development phases.

### Game Design
- [Core Mechanics](docs/01-game-design/core-mechanics.md) — Resources, buildings, troops, weapons, runes, forges, combat.
- [Kingdoms](docs/01-game-design/kingdoms.md) — Veridor, Sylvara, Arkazia — troops, bonuses, playstyles.
- [Progression & Scaling](docs/01-game-design/progression-and-scaling.md) — Cost curves, troop scaling, weapon tiers, balance.

### Lore
- [World & History](docs/02-lore/world-and-history.md) — Aethermoor, the Sundering, Weapons of Chaos, the Prophecy of Order.
- [Kingdom Lore](docs/02-lore/kingdom-lore.md) — Founding stories, cultures, legendary heroes, motivations.

### Architecture
- [System Architecture](docs/03-architecture/system-architecture.md) — Client/server design, WebSocket protocol, auth flow, game ticks.
- [Data Models](docs/03-architecture/data-models.md) — Database schemas for all entities.

### Frontend
- [Frontend Guide](docs/04-frontend/frontend-guide.md) — React conventions, component structure, state management.
- [Styling Guide](docs/04-frontend/styling-guide.md) — CSS architecture, dark/light themes, fonts, mobile strategy.

### Backend
- [Go Guide](docs/05-backend/go-guide.md) — Project structure, clean architecture, error handling, concurrency.
- [Go Multiplayer](docs/05-backend/go-multiplayer.md) — WebSocket hub, game loop, anti-cheat, server-authoritative design.

### Database
- [Database Guide](docs/06-database/database-guide.md) — SQLite patterns, migrations, repository pattern, PostgreSQL migration path.

### Testing
- [Testing Guide](docs/07-testing/testing-guide.md) — Unit/integration test strategy for frontend and backend.

### Security
- [Security & Anti-Cheat](docs/08-security/security-and-anticheat.md) — Auth security, rate limiting, input validation, anti-cheat patterns.

### Agent Rules
- [Copilot Instructions](.github/copilot-instructions.md) — **Mandatory rules for all AI agents working on this project.**

---

## Tech Stack

| Layer | Technology |
|-------|-----------|
| Frontend | React 18+ (TypeScript, Vite) |
| Backend | Go 1.22+ |
| Database | SQLite → PostgreSQL |
| Real-time | WebSockets |
| Auth | JWT + OAuth (Google/Discord) |

---

## The Three Kingdoms

| Kingdom | Theme | Playstyle |
|---------|-------|-----------|
| **Veridor** | Naval / Ocean | Trade mastery, naval superiority |
| **Sylvara** | Forest / Jungle | Guerrilla warfare, rune mastery |
| **Arkazia** | Mountain / Gladiator | Heavy defense, forge mastery |

---

## License

TBD

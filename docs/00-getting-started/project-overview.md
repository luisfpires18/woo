# Project Overview — Weapons of Order (WOO)

> **Read this file FIRST before starting any task.**

---

## Vision

**Weapons of Order** is a multiplayer browser-based RTS game inspired by Travian, built with an original dark fantasy setting. Players choose one of **seven playable kingdoms**, build villages, gather resources, train armies, craft weapons and runes, and compete for dominance on a shared tile-based world map.

The game revolves around a unique core mechanic: **Weapons of Chaos** — immensely powerful artifacts scattered across the world that grant enormous power but inflict devastating debuffs on their wielders. The endgame is triggered when the NPC faction **Moraphys** gathers all Weapons of Chaos, forcing players to forge **Weapons of Order** through alliance-level collaboration to save the world.

The project also includes a **single-player mode**: an interactive lore explorer where players discover the world's history, characters, and backstory through an immersive narrative experience (separate from the multiplayer RTS).

---

## Tech Stack

| Layer | Technology | Notes |
|-------|-----------|-------|
| Frontend | React 19 (TypeScript) | Vite bundler, Zustand state, TanStack Query |
| Backend | Go 1.25+ | Clean architecture, WebSocket-based real-time |
| Database | SQLite (dev) → PostgreSQL (prod) | Repository pattern, UnitOfWork, migration system |
| Real-time | WebSockets | JSON protocol with `type` field |
| Auth | JWT + OAuth (Google/Discord planned) | Short-lived access + refresh tokens |
| Config Pipeline | Go → JSON → TypeScript | Codegen for buildings, troops, resources |
| Testing | Vitest (frontend), Go testing (backend) | Unit + integration tests mandatory |

---

## Monorepo Structure

```
WOO/
├── .github/
│   └── copilot-instructions.md     # AI agent rules (MANDATORY READ)
├── docs/                            # All project documentation
│   ├── 00-getting-started/          # This folder — vision, workflow
│   ├── 01-game-design/              # Mechanics, kingdoms, progression
│   ├── 02-lore/                     # World history, kingdom stories (placeholder)
│   ├── 03-architecture/             # System design, data models
│   ├── 04-frontend/                 # React guide, CSS/styling guide
│   ├── 05-backend/                  # Go guide, multiplayer/concurrency
│   ├── 06-database/                 # DB patterns, migration guide
│   ├── 07-testing/                  # Test strategy for all layers
│   └── 08-security/                 # Anti-cheat, auth, rate limiting
├── client/                          # React frontend application
│   └── src/config/generated/        # Auto-generated JSON from Go configs
├── server/                          # Go backend server
│   └── cmd/genconfig/               # Config codegen CLI tool
└── README.md                        # Quick links to all docs
```

---

## The Eight Kingdoms

### Playable Kingdoms (7)

| Kingdom | Theme | Playstyle |
|---------|-------|-----------|
| **Veridor** | Naval / Ocean / Coastal | Trade mastery, naval superiority, coastal fortifications |
| **Sylvara** | Forest / Jungle / Nature | Guerrilla tactics, fast units, nature-magic rune affinity |
| **Arkazia** | Mountain / Gladiator / Iron | Heavy infantry, superior defense, faster forge crafting |
| **Draxys** | Desert / Sand / Gladiator | Arena warfare, beast-riding, harsh frontier combat |
| **Nordalh** | Frost / Northern / Viking | Brutal melee, forge-crafted elites, cold endurance |
| **Zandres** | Underground / Crystal / Tech | Tunnel warfare, precision engineering, resonance weapons |
| **Lumus** | Light / Holy / Solar | Sacred warriors, radiant discipline, ceremonial combat |

### NPC-Only Kingdom (1)

| Kingdom | Theme | Notes |
|---------|-------|-------|
| **Drakanith** | Volcanic / Draconic | NPC-only. No playable units. Present in lore and world events. |

Each kingdom has ~20 unique troops, kingdom-themed building display names, and a full CSS theme. See `docs/01-game-design/kingdoms.md` and `docs/01-game-design/kingdoms_units_buildlings.md`.

---

## The Enemy: Moraphys

**Moraphys** is an NPC-controlled enemy faction. They do not participate in normal gameplay but grow in power throughout the game round. When Moraphys gathers **all Weapons of Chaos**, the endgame event triggers. Players must then collaborate across alliances to forge **Weapons of Order** and defeat Moraphys before the world is consumed.

> **Game round length**: TBD — draft options are 3 months or 6 months per world. To be decided during playtesting.

---

## Core Game Loop

1. **Settle** — Found your first village, choose your kingdom
2. **Gather** — Produce Food, Water, Lumber, and Stone through resource buildings
3. **Build** — Construct and upgrade buildings (Town Hall, Barracks, Forge, Rune Altar, Walls, etc.)
4. **Train** — Recruit kingdom-specific troops (140 unique troop types across 7 kingdoms)
5. **Forge** — Craft weapons using resources + runes in the Forge
6. **Expand** — Found new villages, claim territory on the world map
7. **Conquer** — Attack enemy villages, raid resources, defend your lands
8. **Discover** — Find and wield Weapons of Chaos (at your own peril)
9. **Unite** — Form alliances, trade, coordinate attacks
10. **Endgame** — Forge Weapons of Order to defeat Moraphys and win the round

---

## Four Base Resources

| Resource | Description |
|----------|------------|
| **Food** | Produced by 3 food buildings (admin-configurable per kingdom). Required to sustain troops and population. Limits army size. |
| **Water** | Collected by 3 water buildings. Used in crop irrigation, troop sustenance, and buildings. |
| **Lumber** | Harvested by 3 lumber buildings. Used in buildings, siege equipment, and ships. |
| **Stone** | Quarried by 3 stone buildings. Used in fortifications, walls, and heavy structures. |

Each resource has 3 building slots per village. Display names, descriptions, and sprites are admin-configurable per kingdom via the `resource_building_configs` table.

Resource economy constants (starting values, rates, storage) are centralized in `server/internal/config/resources.go` and exported to the frontend via the config codegen pipeline.

---

## Key Mechanics (Summary)

- **Runes**: Magical artifacts found via exploration, trading, or conquest. They modify weapon stats and grant special abilities.
- **Forges**: Village buildings where weapons are crafted. Higher-level forges unlock higher-tier weapons.
- **Weapons of Chaos**: Powerful artifacts that exist on the world map. Grant immense power but cause debuffs (resource decay, betrayal events, random disasters). Any player can claim them — but at great cost.
- **Weapons of Order**: Crafted through alliance-level collaboration during the endgame. The only way to counter Moraphys and the gathered Weapons of Chaos.

For full details, see `docs/01-game-design/core-mechanics.md`.

---

## Config Codegen Pipeline

Game configuration follows a **single-source-of-truth** pipeline: Go config files → genconfig CLI → JSON → TypeScript imports.

- **Source**: `server/internal/config/` — `buildings.go`, `troops.go`, `resources.go`
- **Generator**: `server/cmd/genconfig/main.go`
- **Output**: `client/src/config/generated/` — `buildings.json`, `troops.json`, `resources.json`
- **Frontend**: `client/src/config/` — `buildings.ts`, `troops.ts`, `resources.ts` import from generated JSON
- **Parity tests**: `server/internal/config/parity_test.go` verifies committed JSON matches Go config

Run `npm run gen-config` (from repo root) or `cd server && go run ./cmd/genconfig` after changing any Go config. See `.github/copilot-instructions.md` Section 8 for full rules.

---

## Current State

- **Phase**: Active Development (post-MVP, military system complete)
- **Backend**: Go server fully operational — auth (register/login/refresh/logout), village CRUD, building upgrades with queue, lazy resource calculation, world map (51×51 configurable), admin panel (players, config, stats, announcements, game assets, resource building configs, building display configs), map template system (create/edit/resize/apply), game loop (building + training completion ticks with WebSocket notifications), troop training system (140 troops across 7 kingdoms). SQLite with WAL mode.
- **Frontend**: React 19 + TypeScript + Vite — auth pages, kingdom selection, village view with building upgrades and resource ticker, troop training UI (BuildingTrainingModal, TrainingQueue, TroopRoster), Canvas 2D world map renderer with cel-shaded terrain, admin panel (players, config, stats, announcements, assets, map editor with template painting), theme system (kingdom-based themes for all 8 kingdoms). Zustand stores, TanStack Query, CSS Modules.
- **Config Pipeline**: Go → JSON → TypeScript codegen for buildings (21 types), troops (140 types), and resource economy. Parity tests validate committed JSON matches Go source.
- **Database**: 2 consolidated migrations (001_schema.sql + 002_seed_data.sql), SQLite with repository pattern + UnitOfWork for atomic multi-table operations. Tables: players, villages, buildings, building_queue, resources, resource_building_configs, building_display_configs, troops, training_queue, weapons, runes, alliances, alliance_members, world_map, kingdom_relations, attacks, weapons_of_chaos, refresh_tokens, world_config, announcements, game_assets, seasons, season_players.
- **WebSocket**: Foundation implemented — Hub (single-session per player, topic pub/sub, broadcast), Client (read/write pumps, ping/pong), Handler (HTTP upgrade with JWT auth via `?token=`), Messages (all type constants + typed data structs). Game loop broadcasts `build_complete` and `train_complete` events via Hub.
- **Architecture**: UnitOfWork interface encapsulates transactions — services never see `*sql.Tx`. Shared helpers: `FlushResources` (resource_calc.go), `IsValidKingdom` (kingdoms.go). Frontend uses union types (`BuildingType`, `TroopType`, `Kingdom`) for compile-time safety. Config codegen pipeline ensures Go ↔ TypeScript parity. Map utilities deduplicated into `mapUtils.ts`.
- **Not yet implemented**: Combat, weapons/runes crafting, alliances, Weapons of Chaos/Order, Moraphys NPC, marketplace, fog of war, hero system, OAuth, chat.
- **Next step**: Combat system.

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of project overview |
| 2026-03-07 | Updated Current State to reflect actual implementation status (auth, villages, buildings, map, admin panel, templates, game assets all implemented). Fixed resource names from Iron/Wood to Food/Water/Lumber/Stone. |
| 2026-03-08 | Updated Current State: WebSocket foundation implemented, consolidation complete (docs, tests, architecture, WebSocket). 163 backend tests. Next step: Troops & Training. |
| 2026-03-10 | Full codebase audit applied: migrations consolidated (001+002), UnitOfWork pattern, .Hours()→.Seconds() bug fix, shared FlushResources/IsValidKingdom, frontend union types, CSS modules + responsive, map utils extracted. |
| 2026-03-10 | Major update: 7 playable kingdoms (+ Drakanith NPC-only), 140 troops, config codegen pipeline (Go→JSON→TS), resource economy centralization, storage buildings. Tech stack updated to React 19, Go 1.25+. |

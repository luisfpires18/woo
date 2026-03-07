# Project Overview — Weapons of Order (WOO)

> **Read this file FIRST before starting any task.**

---

## Vision

**Weapons of Order** is a multiplayer browser-based RTS game inspired by Travian, built with an original dark fantasy setting. Players choose one of three kingdoms, build villages, gather resources, train armies, craft weapons and runes, and compete for dominance on a shared tile-based world map.

The game revolves around a unique core mechanic: **Weapons of Chaos** — immensely powerful artifacts scattered across the world that grant enormous power but inflict devastating debuffs on their wielders. The endgame is triggered when the NPC faction **Moraphys** gathers all Weapons of Chaos, forcing players to forge **Weapons of Order** through alliance-level collaboration to save the world.

The project also includes a **single-player mode**: an interactive lore explorer where players discover the world's history, characters, and backstory through an immersive narrative experience (separate from the multiplayer RTS).

---

## Tech Stack

| Layer | Technology | Notes |
|-------|-----------|-------|
| Frontend | React 18+ (TypeScript) | Vite bundler, Zustand state, React Query |
| Backend | Go 1.22+ | Clean architecture, WebSocket-based real-time |
| Database | SQLite (dev) → PostgreSQL (prod) | Repository pattern, migration system |
| Real-time | WebSockets | JSON protocol with `type` field |
| Auth | JWT + OAuth (Google/Discord) | Short-lived access + refresh tokens |
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
│   ├── 02-lore/                     # World history, kingdom stories
│   ├── 03-architecture/             # System design, data models
│   ├── 04-frontend/                 # React guide, CSS/styling guide
│   ├── 05-backend/                  # Go guide, multiplayer/concurrency
│   ├── 06-database/                 # DB patterns, migration guide
│   ├── 07-testing/                  # Test strategy for all layers
│   └── 08-security/                 # Anti-cheat, auth, rate limiting
├── client/                          # React frontend application
├── server/                          # Go backend server
└── README.md                        # Quick links to all docs
```

---

## The Three Kingdoms

| Kingdom | Theme | Playstyle |
|---------|-------|-----------|
| **Veridor** | Naval / Ocean / Blue | Trade mastery, naval superiority, coastal fortifications |
| **Sylvara** | Forest / Jungle / Green | Guerrilla tactics, fast units, nature-magic rune affinity |
| **Arkazia** | Mountain / Gladiator / Iron | Heavy infantry, superior defense, faster forge crafting |

Each kingdom has unique troops, bonuses, buildings, and lore. See `docs/01-game-design/kingdoms.md` and `docs/02-lore/kingdom-lore.md`.

---

## The Enemy: Moraphys

**Moraphys** is a NPC-controlled enemy faction. They do not participate in normal gameplay but grow in power throughout the game round. When Moraphys gathers **all Weapons of Chaos**, the endgame event triggers. Players must then collaborate across alliances to forge **Weapons of Order** and defeat Moraphys before the world is consumed.

> **Game round length**: TBD — draft options are 3 months or 6 months per world. To be decided during playtesting.

---

## Core Game Loop

1. **Settle** — Found your first village, choose your kingdom
2. **Gather** — Produce Iron, Wood, Stone, and Food through resource buildings
3. **Build** — Construct and upgrade buildings (Town Hall, Barracks, Forge, Rune Altar, Walls, etc.)
4. **Train** — Recruit kingdom-specific troops
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

---

## Key Mechanics (Summary)

- **Runes**: Magical artifacts found via exploration, trading, or conquest. They modify weapon stats and grant special abilities.
- **Forges**: Village buildings where weapons are crafted. Higher-level forges unlock higher-tier weapons.
- **Weapons of Chaos**: Powerful artifacts that exist on the world map. Grant immense power but cause debuffs (resource decay, betrayal events, random disasters). Any player can claim them — but at great cost.
- **Weapons of Order**: Crafted through alliance-level collaboration during the endgame. The only way to counter Moraphys and the gathered Weapons of Chaos.

For full details, see `docs/01-game-design/core-mechanics.md`.

---

## MVP Scope (Phase 1)

The initial implementation covers:

1. **Authentication**: JWT login + Google/Discord OAuth registration
2. **World Map**: Tile-based grid map with terrain types, villages, and exploration
3. **Basic Village**: Building construction, resource production, building queue
4. **Kingdom Selection**: Player chooses Veridor, Sylvara, or Arkazia at registration

Everything else (troops, combat, weapons, runes, forges, endgame) comes in later phases. See `docs/00-getting-started/development-workflow.md` for the full roadmap.

---

## Current State

- **Phase**: Documentation & Planning
- **Code**: Not yet started
- **Next step**: Complete docs, then scaffold the monorepo with basic project configuration

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of project overview |

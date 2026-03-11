# Kingdoms

> Complete reference for all kingdoms in the game. For detailed troop roster descriptions and building names, see `docs/01-game-design/kingdoms_units_buildlings.md`. For exact troop stats, see `server/internal/config/troops.go`.

---

## Playable vs NPC-Only Kingdoms

The game world has **8 kingdoms**. **7 are playable** and **1 is NPC-only**.

### Playable Kingdoms (7)

| Kingdom | Theme | Troop Count |
|---------|-------|-------------|
| Sylvara | Forest / Jungle / Nature | 19 |
| Arkazia | Mountain / Gladiator / Forge | 20 |
| Veridor | Naval / Ocean / Coastal | 20 |
| Draxys | Desert / Sand / Arena | 21 |
| Nordalh | Frost / Northern / Viking | 20 |
| Zandres | Underground / Crystal / Tech | 20 |
| Lumus | Light / Holy / Solar | 20 |

**Total**: 140 troop types across all 7 playable kingdoms.

### NPC-Only Kingdom (1)

| Kingdom | Theme | Notes |
|---------|-------|-------|
| **Drakanith** | Volcanic / Draconic | NPC-only. No playable units. Present in lore, world events, and zone painting. May become playable in future expansions. |

---

## Overview

Players choose their kingdom at registration. This choice is **permanent** for the game round and determines:
- Available troop types and their stats (~20 unique units per kingdom)
- Kingdom-specific building display names (admin-configurable)
- Visual theme (full CSS theme with unique colors per kingdom)
- Playstyle identity

Kingdom bonuses (resource production bonuses, terrain bonuses, etc.) are **designed but not yet implemented**. They will be added during the balance/combat phase.

---

## Sylvara — The Verdant Wilds

> *"The jungle does not forgive. Neither do we."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Forest / Jungle / Nature-magic |
| **CSS Theme** | Forest green (#2E7D32) accent, golden parchment background, black text |
| **Architecture** | Treehouse fortresses, vine-wrapped palisades, living-wood walls |
| **Playstyle** | Guerrilla warfare, fast strikes, rune mastery, ambush tactics |
| **Troops** | 19 units (4 barracks, 3 stable, 4 archery, 4 workshop, 4 special) |

### Military Buildings

| Building | Kingdom Display Name |
|----------|---------------------|
| Barracks | Roothall |
| Stable | Beast Hall |
| Archery | Grove Range |
| Workshop | Tree-Sapper Yard |
| Special | Spirit Glade |

### Designed Bonuses (Not Yet Implemented)

- +15% Food production (fertile jungle lands)
- +10% rune discovery rate (rune affinity)
- Ambush: Sylvara troops deal 15% bonus damage when attacking from forest tiles

---

## Arkazia — The Iron Peaks

> *"Steel is forged in fire. So are we."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Mountain / Gladiator / Knights |
| **CSS Theme** | Crimson red (#DC143C) accent, black background, white text |
| **Architecture** | Stone fortresses, mountain citadels, colosseum arenas, iron-reinforced walls |
| **Playstyle** | Heavy defense, superior crafting, brute force, siege warfare |
| **Troops** | 20 units (4 barracks, 4 stable, 4 archery, 4 workshop, 4 special) |

### Military Buildings

| Building | Kingdom Display Name |
|----------|---------------------|
| Barracks | Red Bastion |
| Stable | Arknight Stables |
| Archery | Ridge Range |
| Workshop | Stonecaller Yard |
| Special | Chapter Fortress |

### Designed Bonuses (Not Yet Implemented)

- +15% Stone production (mountain mining tradition)
- +15% Forge crafting speed (master smiths)
- Fortification: Arkazia walls provide 20% more defense bonus

---

## Veridor — The Tidal Dominion

> *"The sea provides. The sea destroys. We are both."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Naval / Ocean / Coastal |
| **CSS Theme** | Blue (#2196F3) accent, light (white/blue tint) background, black text |
| **Architecture** | Stone harbors, lighthouse towers, coral-reinforced walls |
| **Playstyle** | Trade mastery, naval superiority, economic warfare |
| **Troops** | 20 units (4 barracks, 4 stable, 4 archery, 4 workshop, 4 special) |

### Military Buildings

| Building | Kingdom Display Name |
|----------|---------------------|
| Barracks | Road Barracks |
| Stable | River Cavalry Yard |
| Archery | Chart Range |
| Workshop | Shipwright Siegeyard |
| Special | Admiralty Hall |

### Designed Bonuses (Not Yet Implemented)

- +15% Lumber production (shipbuilding tradition)
- +10% trade income at Marketplace
- Naval supremacy: Veridor troops move 20% faster on water tiles

---

## Draxys — The Scorched Frontier

> *"In the arena, there is only truth."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Desert / Sand / Gladiator Arena |
| **CSS Theme** | Yellow (#F9A825) accent, dark background, white text |
| **Architecture** | Sandstone walls, arena coliseums, scorpion-shaped siege works |
| **Playstyle** | Arena-focused, beast-riders, brutal shock warfare |
| **Troops** | 21 units (4 barracks, 4 stable, 4 archery, 4 workshop, 5 special) |

### Military Buildings

| Building | Kingdom Display Name |
|----------|---------------------|
| Barracks | Sandwall Barracks |
| Stable | Scorpion Pens |
| Archery | Oasis Range |
| Workshop | Sandwall Foundry |
| Special | Grand Arena |

### Designed Bonuses (Not Yet Implemented)

- +10% Food production (oasis farming)
- Gladiator morale: Draxys troops deal 10% bonus damage when outnumbered
- Desert endurance: No movement penalty on desert tiles

---

## Nordalh — The Frozen Holds

> *"We forge in fire. We fight in frost."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Frost / Northern / Viking |
| **CSS Theme** | Purple (#7B1FA2) accent, light (white/purple tint) background, black text |
| **Architecture** | Longhouses, stone forges, ice-reinforced palisades |
| **Playstyle** | Brutal melee, forge-crafted elites, cold endurance, wolf-riders |
| **Troops** | 20 units (4 barracks, 4 stable, 4 archery, 4 workshop, 4 special) |

### Military Buildings

| Building | Kingdom Display Name |
|----------|---------------------|
| Barracks | Hearth Barracks |
| Stable | Wolf Kennels |
| Archery | Ice Loom Range |
| Workshop | Long Forge Siegeyard |
| Special | Long Forge Hall |

### Designed Bonuses (Not Yet Implemented)

- +15% Stone production (mountain quarrying)
- +10% weapon durability (master forge-craft)
- Frost endurance: No movement penalty on swamp tiles

---

## Zandres — The Deep Lattice

> *"What you cannot see, you cannot fight."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Underground / Crystal / Technology |
| **CSS Theme** | Brown (#795548) accent, dark background, white text |
| **Architecture** | Crystal-lattice halls, bored tunnels, resonance engines |
| **Playstyle** | Precision engineering, tunnel warfare, tech-augmented troops |
| **Troops** | 20 units (4 barracks, 4 stable, 4 archery, 4 workshop, 4 special) |

### Military Buildings

| Building | Kingdom Display Name |
|----------|---------------------|
| Barracks | Doorwarden Hall |
| Stable | Crawler Pens |
| Archery | Crystal Range |
| Workshop | Resonance Works |
| Special | Circuit Archive |

### Designed Bonuses (Not Yet Implemented)

- +15% Stone production (deep mining)
- Resonance detection: Zandres watchtowers have +2 tile detection range
- Tunnel movement: Zandres troops move 15% faster on mountain tiles

---

## Lumus — The Radiant Court

> *"Light reveals all. Light endures all."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Light / Holy / Solar Ritual |
| **CSS Theme** | Golden Yellow (#FBC02D) accent, light (white/yellow tint) background, black text |
| **Architecture** | Mirror towers, prism barracks, sun-court plazas |
| **Playstyle** | Sacred warrior discipline, radiant defense, ceremonial combat |
| **Troops** | 20 units (4 barracks, 4 stable, 4 archery, 4 workshop, 4 special) |

### Military Buildings

| Building | Kingdom Display Name |
|----------|---------------------|
| Barracks | Prism Barracks |
| Stable | Sun Court Stables |
| Archery | Sunshot Range |
| Workshop | Heliostat Works |
| Special | Heliostat Sanctum |

### Designed Bonuses (Not Yet Implemented)

- +15% Water production (sacred springs)
- Radiant defense: Lumus walls glow, reducing attacker accuracy by 5%
- Holy morale: Lumus troops gain +10% defense when defending their own village

---

## Drakanith — The Ember Throne (NPC-Only)

> *"The dragon's blood is both gift and curse."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Volcanic / Draconic |
| **CSS Theme** | Orange (#FF6D00) accent, dark background, white text |
| **Status** | **NPC-only** — not selectable by players |
| **Role** | Appears in lore, world events, and as NPC villages/defenders |

Drakanith has no playable troop roster. Their data (theme, building configs, resource configs) is seeded in the database for zone painting and NPC use.

---

## Kingdom Comparison Matrix

| Aspect | Sylvara | Arkazia | Veridor | Draxys | Nordalh | Zandres | Lumus |
|--------|---------|---------|---------|--------|---------|---------|-------|
| **Total Troops** | 19 | 20 | 20 | 21 | 20 | 20 | 20 |
| **Theme** | Nature | Mountain | Naval | Desert | Frost | Underground | Holy |
| **CSS Background** | Warm parchment | Dark (black) | Light (blue tint) | Dark (black) | Light (purple tint) | Dark (black) | Light (yellow tint) |
| **Accent Color** | Green | Crimson | Blue | Yellow | Purple | Brown | Golden |

---

## Kingdom Theming (CSS)

Each kingdom has a full CSS variable theme applied via `data-kingdom` attribute on `<html>`. The theme overrides `--bg-*`, `--text-*`, `--accent*`, `--border*`, and `--shadow-*` variables. See `docs/04-frontend/styling-guide.md` and `client/src/styles/themes.css` for details.

| Kingdom | `--text-on-accent` |
|---------|-------------------|
| Arkazia | #FFFFFF |
| Draxys | #000000 |
| Drakanith | #FFFFFF |
| Zandres | #FFFFFF |
| Veridor | #FFFFFF |
| Nordalh | #FFFFFF |
| Lumus | #000000 |
| Sylvara | #FFFFFF |

---

## Balance Philosophy

- No kingdom should be strictly superior. Each excels in different situations.
- **Rock-Paper-Scissors dynamics** between kingdom playstyles.
- Balance will be refined through playtesting. Kingdom bonuses are designed but not yet implemented.
- Troop stats are balanced per building-tier across kingdoms. See `server/internal/config/troops.go` for exact values.
- Kingdom-specific UI theming is fully implemented for all 8 kingdoms (7 playable + Drakanith NPC).

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of kingdoms reference |
| 2026-03-06 | Updated kingdom color specs to match implemented kingdom themes |
| 2026-03-07 | Split kingdoms into 5 playable + 3 NPC-only |
| 2026-03-10 | Major rewrite: All 7 kingdoms now playable with full troop rosters (140 troops total). Drakanith is the only NPC-only kingdom. Added all 7 kingdom sections with building names, themes, troop counts. Removed old draft troop stat tables (actual stats live in config codegen pipeline). Marked kingdom bonuses as designed but not implemented. |

# Kingdoms

> **Superseded**: Definitive tunable values (troop stats, bonuses, costs, etc.) are in [`game-template.md`](game-template.md). Values below are **drafts** — when they conflict, `game-template.md` wins.

> Complete reference for all three playable kingdoms. Cross-reference with `docs/02-lore/kingdom-lore.md` for backstory.

---

## Overview

Players choose their kingdom at registration. This choice is **permanent** for the game round and determines:
- Available troop types and their stats
- Kingdom-specific building (1 unique building per kingdom)
- Resource production bonuses
- Visual theme (UI colors, building art, troop art)
- Playstyle identity

---

## Veridor — The Tidal Dominion

> *"The sea provides. The sea destroys. We are both."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Naval / Ocean / Coastal |
| **Colors** | Blue (#2196F3), white background, black text |
| **Architecture** | Stone harbors, lighthouse towers, coral-reinforced walls |
| **Playstyle** | Trade mastery, naval superiority, economic warfare |

### Kingdom Bonuses

- **+15% Wood production** (shipbuilding tradition)
- **+10% trade income** at Marketplace
- **Naval supremacy**: Veridor troops move 20% faster on water tiles
- **Unique Building**: **Dock** — builds naval units and enables sea trade routes

### Troop Roster

| Unit | Type | Attack | Def (Inf) | Def (Cav) | Speed | Carry | Upkeep | Notes |
|------|------|--------|-----------|-----------|-------|-------|--------|-------|
| **Tidecaller** | Infantry | 40 | 50 | 30 | 6 | 50 | 1 | Balanced frontline, good defense |
| **Harpooneer** | Ranged | 55 | 20 | 25 | 7 | 30 | 1 | High attack, fragile |
| **Wave Rider** | Cavalry | 70 | 30 | 40 | 12 | 80 | 2 | Fast raider |
| **Coral Sentinel** | Heavy Inf | 30 | 70 | 60 | 4 | 20 | 2 | Dedicated wall defender |
| **Sea Serpent** | Naval | 90 | 40 | 40 | 14 | 100 | 3 | Powerful but only on water. Dock required. |
| **Stormcaster** | Siege | 60 | 15 | 15 | 3 | 10 | 3 | Destroys walls. Slow. |
| **Gull Scout** | Scout | 10 | 5 | 5 | 18 | 0 | 1 | Fastest scout. Reveals enemy info. |

### Strengths & Weaknesses

- **Strengths**: Best economy through trade. Naval dominance on water maps. Fast scouts. Strong defensive infantry.
- **Weaknesses**: Expensive troops. Relies heavily on water tiles. Weaker siege capability on land.

---

## Sylvara — The Verdant Wilds

> *"The jungle does not forgive. Neither do we."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Forest / Jungle / Nature-magic |
| **Colors** | Forest green (#2E7D32), golden parchment background, black text |
| **Architecture** | Treehouse fortresses, vine-wrapped palisades, living-wood walls |
| **Playstyle** | Guerrilla warfare, fast strikes, rune mastery, ambush tactics |

### Kingdom Bonuses

- **+15% Food production** (fertile jungle lands)
- **+10% rune discovery rate** (rune affinity)
- **Ambush**: Sylvara troops deal 15% bonus damage when attacking from forest tiles
- **Unique Building**: **Grove Sanctum** — enhances rune altar, doubles rune combination success rate

### Troop Roster

| Unit | Type | Attack | Def (Inf) | Def (Cav) | Speed | Carry | Upkeep | Notes |
|------|------|--------|-----------|-----------|-------|-------|--------|-------|
| **Thornguard** | Infantry | 45 | 40 | 35 | 7 | 45 | 1 | Fast infantry, balanced |
| **Venomarcher** | Ranged | 60 | 15 | 20 | 8 | 25 | 1 | Poison DOT, very high attack |
| **Panther Rider** | Cavalry | 75 | 25 | 45 | 14 | 70 | 2 | Fastest cavalry in the game |
| **Rootwarden** | Heavy Inf | 25 | 60 | 55 | 5 | 15 | 2 | Defensive but can root enemies (slow debuff) |
| **Jungle Stalker** | Assassin | 85 | 10 | 10 | 10 | 40 | 2 | Massive attack, paper defense. Ambush specialist. |
| **Siege Vine** | Siege | 50 | 20 | 20 | 4 | 10 | 3 | Living siege engine. Bypasses walls partially. |
| **Hawk Eye** | Scout | 8 | 8 | 8 | 16 | 0 | 1 | Good scouting, slightly tougher than other scouts |

### Strengths & Weaknesses

- **Strengths**: Fastest troops overall. Best rune synergy. Devastating ambushes from forest tiles. Strong food economy supports large armies.
- **Weaknesses**: Fragile elite units. Poor defense if caught in open terrain. Weak on water/desert maps.

---

## Arkazia — The Iron Peaks

> *"Steel is forged in fire. So are we."*

### Identity

| Property | Value |
|----------|-------|
| **Theme** | Mountain / Gladiator / Knights |
| **Colors** | Crimson red (#DC143C), black background, white text |
| **Architecture** | Stone fortresses, mountain citadels, colosseum arenas, iron-reinforced walls |
| **Playstyle** | Heavy defense, superior crafting, brute force, siege warfare |

### Kingdom Bonuses

- **+15% Iron production** (mountain mining tradition)
- **+15% Forge crafting speed** (master smiths)
- **Fortification**: Arkazia walls provide 20% more defense bonus
- **Unique Building**: **Colosseum** — trains elite gladiator units, provides morale bonus to defending troops

### Troop Roster

| Unit | Type | Attack | Def (Inf) | Def (Cav) | Speed | Carry | Upkeep | Notes |
|------|------|--------|-----------|-----------|-------|-------|--------|-------|
| **Iron Legionary** | Infantry | 50 | 55 | 40 | 5 | 55 | 1 | Best standard infantry in the game |
| **Crossbowman** | Ranged | 50 | 30 | 25 | 5 | 35 | 1 | Balanced ranged, good defense for ranged |
| **Mountain Knight** | Cavalry | 80 | 40 | 50 | 10 | 90 | 3 | Slowest cavalry but hits hardest |
| **Shieldbearer** | Heavy Inf | 20 | 80 | 70 | 3 | 10 | 2 | Best pure defender in the game |
| **Gladiator** | Elite | 95 | 50 | 45 | 7 | 60 | 3 | Colosseum-trained. Highest attack of any standard unit. |
| **Battering Ram** | Siege | 70 | 25 | 10 | 3 | 0 | 4 | Best siege unit. Destroys walls efficiently. |
| **Mountain Scout** | Scout | 12 | 10 | 10 | 12 | 0 | 1 | Slowest scout but toughest. Hard to kill. |

### Strengths & Weaknesses

- **Strengths**: Best defensive troops. Best siege. Best forge/weapon crafting. Iron economy feeds military.
- **Weaknesses**: Slowest troops overall. Expensive elite units. Poor food production means smaller armies unless farming is prioritized.

---

## Kingdom Comparison Matrix

| Aspect | Veridor | Sylvara | Arkazia |
|--------|---------|---------|---------|
| **Resource Bonus** | +15% Wood | +15% Food | +15% Iron |
| **Special Bonus** | +10% trade | +10% rune discovery | +15% forge speed |
| **Terrain Bonus** | Water (+20% speed) | Forest (+15% damage) | Walls (+20% defense) |
| **Unique Building** | Dock | Grove Sanctum | Colosseum |
| **Speed** | Medium | Fast | Slow |
| **Attack** | Medium | High (ambush) | High (brute force) |
| **Defense** | High (naval) | Low-Medium | Very High |
| **Economy** | Best (trade) | Good (food) | Good (iron) |
| **Crafting** | Normal | Rune-focused | Forge-focused |

---

## Balance Philosophy

- No kingdom should be strictly superior. Each excels in different situations.
- **Rock-Paper-Scissors dynamics**: Veridor's economy outpaces Arkazia's slow expansion, Sylvara's speed counters Veridor's naval strategy, Arkazia's defense walls out Sylvara's hit-and-run tactics.
- Balance will be refined through playtesting. Stat values in this document are **draft** — expect iteration.
- Kingdom-specific UI theming is implemented — each kingdom has a full CSS variable theme applied via `data-kingdom` attribute. See `docs/04-frontend/styling-guide.md`.

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of kingdoms reference |
| 2026-03-05 | Marked as superseded by `game-template.md` for all tunable values |
| 2026-03-06 | Updated kingdom color specs to match implemented kingdom themes |

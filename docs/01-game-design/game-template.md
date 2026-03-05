# Game Template — Authoritative Values

> **This document is the single source of truth for ALL tunable game values.**
> Fill in every `TBD` cell. Once complete, all code implementations will be derived from this file.
>
> Values in `core-mechanics.md`, `progression-and-scaling.md`, and `kingdoms.md` are **drafts/superseded**.
> When in doubt, this file wins.

---

## A. Resources

### A1. Base Resources

| Resource | Produced By | Primary Uses |
|----------|------------|-------------|
| **Iron** | Iron Mine | TBD |
| **Wood** | Lumber Mill | TBD |
| **Stone** | Quarry | TBD |
| **Food** | Farm | TBD |

### A2. Starting Resources (New Village)

| Resource | Starting Amount |
|----------|----------------|
| Iron | TBD |
| Wood | TBD |
| Stone | TBD |
| Food | TBD |

### A3. Production Rate Curve (Resources per Hour)

All four resource buildings (Iron Mine, Lumber Mill, Quarry, Farm) use the same production curve unless overridden below.

| Level | Resources/Hour |
|-------|---------------|
| 0 | 0 |
| 1 | TBD |
| 2 | TBD |
| 3 | TBD |
| 4 | TBD |
| 5 | TBD |
| 6 | TBD |
| 7 | TBD |
| 8 | TBD |
| 9 | TBD |
| 10 | TBD |
| 11 | TBD |
| 12 | TBD |
| 13 | TBD |
| 14 | TBD |
| 15 | TBD |
| 16 | TBD |
| 17 | TBD |
| 18 | TBD |
| 19 | TBD |
| 20 | TBD |

> If a kingdom bonus applies (e.g., Veridor +X% Wood), the bonus is applied on top of this base rate.

### A4. Warehouse Capacity Curve

| Level | Max Storage per Resource |
|-------|------------------------|
| 0 | 0 |
| 1 | TBD |
| 2 | TBD |
| 3 | TBD |
| 4 | TBD |
| 5 | TBD |
| 6 | TBD |
| 7 | TBD |
| 8 | TBD |
| 9 | TBD |
| 10 | TBD |
| 11 | TBD |
| 12 | TBD |
| 13 | TBD |
| 14 | TBD |
| 15 | TBD |
| 16 | TBD |
| 17 | TBD |
| 18 | TBD |
| 19 | TBD |
| 20 | TBD |

### A5. Food Consumption Rules

| Rule | Value |
|------|-------|
| Food consumed per troop per hour | TBD (can vary per troop; see troop tables) |
| Starvation effect (when food < 0) | TBD (e.g., troops die, morale drop, etc.) |
| Starvation check interval | TBD (e.g., every X minutes) |

---

## B. Buildings (Shared / Common)

### B1. Building Definitions

| Building | Canonical ID | Function | Max Level | Prerequisites | Starting Level (New Village) |
|----------|-------------|----------|-----------|---------------|------------------------------|
| Town Hall | `town_hall` | Central building; gates access to all others | TBD | None | TBD |
| Iron Mine | `iron_mine` | Produces Iron per hour | TBD | TBD | TBD |
| Lumber Mill | `lumber_mill` | Produces Wood per hour | TBD | TBD | TBD |
| Quarry | `quarry` | Produces Stone per hour | TBD | TBD | TBD |
| Farm | `farm` | Produces Food per hour; determines population cap | TBD | TBD | TBD |
| Warehouse | `warehouse` | Stores resources; level sets max storage | TBD | TBD | TBD |
| Barracks | `barracks` | Trains infantry troops | TBD | TBD | TBD |
| Stable | `stable` | Trains mounted/fast troops | TBD | TBD | TBD |
| Forge | `forge` | Crafts weapons from resources + runes | TBD | TBD | TBD |
| Rune Altar | `rune_altar` | Combines, enhances, and stores runes | TBD | TBD | TBD |
| Walls | `walls` | Passive defense bonus for the village | TBD | TBD | TBD |
| Marketplace | `marketplace` | Trade resources with other players | TBD | TBD | TBD |
| Embassy | `embassy` | Required to form/join alliances | TBD | TBD | TBD |
| Watchtower | `watchtower` | Detects incoming attacks | TBD | TBD | TBD |

> **Prerequisites format**: `building_id level` (e.g., `town_hall 5, barracks 3`). Use `None` if no prerequisites.

### B2. Base Costs (Level 1) — All Buildings

| Building | Iron | Wood | Stone | Food | Build Time | Scaling Factor |
|----------|------|------|-------|------|-----------|----------------|
| Town Hall | TBD | TBD | TBD | TBD | TBD | TBD |
| Iron Mine | TBD | TBD | TBD | TBD | TBD | TBD |
| Lumber Mill | TBD | TBD | TBD | TBD | TBD | TBD |
| Quarry | TBD | TBD | TBD | TBD | TBD | TBD |
| Farm | TBD | TBD | TBD | TBD | TBD | TBD |
| Warehouse | TBD | TBD | TBD | TBD | TBD | TBD |
| Barracks | TBD | TBD | TBD | TBD | TBD | TBD |
| Stable | TBD | TBD | TBD | TBD | TBD | TBD |
| Forge | TBD | TBD | TBD | TBD | TBD | TBD |
| Rune Altar | TBD | TBD | TBD | TBD | TBD | TBD |
| Walls | TBD | TBD | TBD | TBD | TBD | TBD |
| Marketplace | TBD | TBD | TBD | TBD | TBD | TBD |
| Embassy | TBD | TBD | TBD | TBD | TBD | TBD |
| Watchtower | TBD | TBD | TBD | TBD | TBD | TBD |

> **Cost formula**: `cost(level) = base_cost × (scaling_factor ^ (level - 1))`
> **Time formula**: `time(level) = base_time × (time_factor ^ (level - 1))` — `time_factor` = TBD (global or per-building)

### B3. Construction Rules

| Rule | Value |
|------|-------|
| Max simultaneous constructions per village | TBD |
| Can this be upgraded (e.g., Town Hall unlocks parallel queues)? | TBD |
| Demolition policy | TBD (instant? cost? cooldown?) |
| Resources returned on demolition | TBD (0%? 50%?) |

---

## C. Buildings (Kingdom-Specific)

### C1. Veridor — Dock

| Property | Value |
|----------|-------|
| Canonical ID | `dock` |
| Function | TBD |
| Max Level | TBD |
| Prerequisites | TBD |
| Starting Level | 0 (not built) |

**Base Cost (Level 1)**

| Iron | Wood | Stone | Food | Build Time | Scaling Factor |
|------|------|-------|------|-----------|----------------|
| TBD | TBD | TBD | TBD | TBD | TBD |

### C2. Sylvara — Grove Sanctum

| Property | Value |
|----------|-------|
| Canonical ID | `grove_sanctum` |
| Function | TBD |
| Max Level | TBD |
| Prerequisites | TBD |
| Starting Level | 0 (not built) |

**Base Cost (Level 1)**

| Iron | Wood | Stone | Food | Build Time | Scaling Factor |
|------|------|-------|------|-----------|----------------|
| TBD | TBD | TBD | TBD | TBD | TBD |

### C3. Arkazia — Colosseum

| Property | Value |
|----------|-------|
| Canonical ID | `colosseum` |
| Function | TBD |
| Max Level | TBD |
| Prerequisites | TBD |
| Starting Level | 0 (not built) |

**Base Cost (Level 1)**

| Iron | Wood | Stone | Food | Build Time | Scaling Factor |
|------|------|-------|------|-----------|----------------|
| TBD | TBD | TBD | TBD | TBD | TBD |

### C4. Kingdom Bonuses

| Kingdom | Resource Bonus | Special Bonus | Terrain Bonus |
|---------|---------------|---------------|---------------|
| Veridor | TBD (e.g., +X% Wood) | TBD (e.g., +X% trade income) | TBD (e.g., water tiles +X% speed) |
| Sylvara | TBD (e.g., +X% Food) | TBD (e.g., +X% rune discovery) | TBD (e.g., forest tiles +X% damage) |
| Arkazia | TBD (e.g., +X% Iron) | TBD (e.g., +X% forge speed) | TBD (e.g., walls +X% defense) |

---

## D. Troops

### D1. Veridor Troop Roster

| # | Unit Name | Type | Attack | Def (Inf) | Def (Cav) | Speed | Carry | Upkeep | Training Building | Building Lvl Req | Iron | Wood | Stone | Food | Training Time | Notes |
|---|-----------|------|--------|-----------|-----------|-------|-------|--------|-------------------|-----------------|------|------|-------|------|--------------|-------|
| 1 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 2 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 3 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 4 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 5 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 6 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 7 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |

> **Type options**: Infantry, Ranged, Cavalry, Heavy Infantry, Naval, Siege, Scout, Assassin, Elite, Special
> **Training Building**: e.g., `barracks`, `stable`, `dock`

### D2. Sylvara Troop Roster

| # | Unit Name | Type | Attack | Def (Inf) | Def (Cav) | Speed | Carry | Upkeep | Training Building | Building Lvl Req | Iron | Wood | Stone | Food | Training Time | Notes |
|---|-----------|------|--------|-----------|-----------|-------|-------|--------|-------------------|-----------------|------|------|-------|------|--------------|-------|
| 1 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 2 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 3 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 4 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 5 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 6 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 7 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |

### D3. Arkazia Troop Roster

| # | Unit Name | Type | Attack | Def (Inf) | Def (Cav) | Speed | Carry | Upkeep | Training Building | Building Lvl Req | Iron | Wood | Stone | Food | Training Time | Notes |
|---|-----------|------|--------|-----------|-----------|-------|-------|--------|-------------------|-----------------|------|------|-------|------|--------------|-------|
| 1 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 2 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 3 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 4 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 5 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 6 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| 7 | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |

### D4. Troop Type Reference

> Fill this in if you want to define the valid troop type categories and what they mean.

| Type | Description |
|------|------------|
| Infantry | TBD |
| Ranged | TBD |
| Cavalry | TBD |
| Heavy Infantry | TBD |
| Naval | TBD |
| Siege | TBD |
| Scout | TBD |
| Assassin | TBD |
| Elite | TBD |

---

## E. Weapons

### E1. Weapon Types

| Weapon Type | Description |
|-------------|------------|
| TBD | TBD |
| TBD | TBD |
| TBD | TBD |
| TBD | TBD |
| TBD | TBD |
| TBD | TBD |

> Add or remove rows as needed (e.g., Sword, Axe, Bow, Spear, Shield, Staff, etc.).

### E2. Weapon Tiers

| Tier | Forge Level Req | Rune Tier Req | Iron | Wood | Stone | Food | Craft Time | Attack Bonus Range | Defense Bonus Range | Rune Slots | Durability |
|------|----------------|---------------|------|------|-------|------|-----------|-------------------|-------------------|-----------|-----------|
| Common | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| Rare | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| Epic | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| Legendary | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |
| Mythic | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD | TBD |

### E3. Weapon Durability & Repair

| Rule | Value |
|------|-------|
| Durability lost per battle | TBD |
| Repair cost (% of crafting cost per point) | TBD |
| Broken weapon behavior | TBD (lost forever? rune salvage chance?) |
| Rune salvage chance on break | TBD |

### E4. Weapons of Chaos Stats

| Property | Value |
|----------|-------|
| Attack Bonus | TBD |
| Defense Bonus | TBD |
| Special Effects (count) | TBD |
| Durability | TBD (infinite?) |
| Rune Slots | TBD (0? built-in?) |

**Wielder Debuffs (cost of holding a Weapon of Chaos)**

| Debuff | Value |
|--------|-------|
| Resource decay (% per period) | TBD |
| Decay period | TBD |
| Troop desertion chance (per day) | TBD |
| Chaos storm radius (tiles) | TBD |
| Chaos storm production penalty | TBD |
| Village visibility | TBD (revealed to all?) |
| Moraphys aggro increase | TBD |

### E5. Weapons of Order — Crafting Requirements

| Resource / Requirement | Amount |
|-----------------------|--------|
| Iron | TBD |
| Wood | TBD |
| Stone | TBD |
| Food | TBD |
| Primordial Runes | TBD |
| Minimum Forge Level | TBD |
| Crafting Time | TBD |
| Number needed to challenge Moraphys | TBD (= Weapons of Chaos count) |

---

## F. Runes

### F1. Rune Rarity Tiers

| Rarity | Combination Input | Success Rate | Rune Altar Level Req |
|--------|------------------|-------------|---------------------|
| Fragment | — (found via drops/exploration) | — | — |
| Minor | TBD (e.g., 3 Fragments) | TBD | TBD |
| Major | TBD (e.g., 3 Minor) | TBD | TBD |
| Grand | TBD (e.g., 3 Major) | TBD | TBD |
| Primordial | TBD (e.g., 3 Grand + condition) | TBD | TBD |

> **Sylvara bonus**: Grove Sanctum effect on success rate = TBD (e.g., doubles it).

### F2. Rune Effect Scaling

| Effect Type | Fragment | Minor | Major | Grand | Primordial |
|-------------|---------|-------|-------|-------|-----------|
| Attack Bonus | TBD | TBD | TBD | TBD | TBD |
| Defense Bonus | TBD | TBD | TBD | TBD | TBD |
| Speed Bonus | TBD | TBD | TBD | TBD | TBD |
| Resource Bonus | TBD | TBD | TBD | TBD | TBD |
| Special Effect | TBD | TBD | TBD | TBD | TBD |

### F3. How to Obtain Runes

| Method | Drop Rate / Details |
|--------|-------------------|
| Exploration (map tiles) | TBD |
| Combat drops | TBD |
| Rune Altar combination | See F1 |
| Marketplace trading | TBD |
| Quests / Events | TBD |

---

## G. Combat

### G1. Combat Resolution Formula

> Describe the formula or leave TBD. This controls how battles are calculated.

```
Attacker Power = TBD
Defender Power = TBD
Winner = TBD
Losses = TBD
```

### G2. Walls Defense Bonus Curve

| Wall Level | Defense Bonus |
|-----------|-------------|
| 0 | 0% |
| 1 | TBD |
| 2 | TBD |
| 3 | TBD |
| 4 | TBD |
| 5 | TBD |
| 6 | TBD |
| 7 | TBD |
| 8 | TBD |
| 9 | TBD |
| 10 | TBD |
| 11 | TBD |
| 12 | TBD |
| 13 | TBD |
| 14 | TBD |
| 15 | TBD |
| 16 | TBD |
| 17 | TBD |
| 18 | TBD |
| 19 | TBD |
| 20 | TBD |

### G3. Siege Rules

| Rule | Value |
|------|-------|
| Siege damage to walls formula | TBD |
| Siege unit bypass % (if any) | TBD |
| Defender advantage multiplier | TBD |

---

## H. Map & World

### H1. Map Dimensions

| Property | Value |
|----------|-------|
| Map Width (tiles) | TBD |
| Map Height (tiles) | TBD |
| Coordinate Range | TBD (e.g., -200 to +200) |
| Center Tile | (0, 0) — Moraphys Stronghold |

### H2. Terrain Types

| Terrain | DB Value | Distribution (%) | Movement Modifier |
|---------|----------|------------------|-------------------|
| Plains | `plains` | TBD | TBD |
| Forest | `forest` | TBD | TBD |
| Mountain | `mountain` | TBD | TBD |
| Water | `water` | TBD | TBD |
| Desert | `desert` | TBD | TBD |
| Swamp | `swamp` | TBD | TBD |

### H3. World Settings

| Setting | Value |
|---------|-------|
| Weapons of Chaos count (default) | TBD |
| Chaos Shrine NPC defender strength | TBD |
| Oasis frequency (% of land tiles) | TBD |
| Oasis resource bonus | TBD |
| Village minimum distance (tiles) | TBD |
| Starting village terrain requirement | TBD (e.g., plains only) |

### H4. Movement

| Rule | Value |
|------|-------|
| Distance formula | TBD (Euclidean? Manhattan?) |
| Speed unit | TBD (tiles per hour) |
| Terrain modifier applies to | TBD (path? destination only?) |
| Slowest-unit rule | TBD (group speed = slowest?) |

---

## I. Utility Building Scaling

### I1. Barracks — Training Speed Bonus

| Level | Training Speed Multiplier |
|-------|--------------------------|
| 1 | TBD |
| 5 | TBD |
| 10 | TBD |
| 15 | TBD |
| 20 | TBD |

### I2. Stable — Unlock Levels

| Stable Level | Units Unlocked |
|-------------|---------------|
| TBD | TBD |
| TBD | TBD |
| TBD | TBD |

### I3. Forge — Tier Unlocks

| Forge Level | Max Weapon Tier |
|-------------|----------------|
| TBD | Common |
| TBD | Rare |
| TBD | Epic |
| TBD | Legendary |
| TBD | Mythic |

### I4. Marketplace — Trade Capacity

| Level | Max Resources per Trade |
|-------|------------------------|
| 1 | TBD |
| 5 | TBD |
| 10 | TBD |
| 15 | TBD |

### I5. Watchtower — Detection

| Level | Warning Time | Detail Revealed |
|-------|-------------|-----------------|
| 1 | TBD | TBD |
| 3 | TBD | TBD |
| 5 | TBD | TBD |
| 8 | TBD | TBD |
| 10 | TBD | TBD |

### I6. Embassy — Alliance Size

| Level | Max Alliance Members |
|-------|---------------------|
| 1 | TBD |
| 3 | TBD |
| 5 | TBD |
| 8 | TBD |
| 10 | TBD |

---

## J. Game Round

| Setting | Value |
|---------|-------|
| Round duration | TBD |
| Endgame timer (after Moraphys collects all Chaos weapons) | TBD |
| Beginner protection duration | TBD |
| Speed modifier (configurable via admin) | TBD |

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-05 | Initial creation — all values TBD for manual definition |

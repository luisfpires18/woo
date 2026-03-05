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
| Iron | 500 |
| Wood | 500 |
| Stone | 500 |
| Food | 500 |

### A3. Production Rate Curve (Resources per Hour)

All four resource buildings (Iron Mine, Lumber Mill, Quarry, Farm) use the same production curve unless overridden below.

**Formula**: `rate = base_rate + (level × rate_per_level)` where `base_rate = 10`, `rate_per_level = 20`.

| Level | Resources/Hour |
|-------|---------------|
| 0 | 10 |
| 1 | 30 |
| 2 | 50 |
| 3 | 70 |
| 4 | 90 |
| 5 | 110 |
| 6 | 130 |
| 7 | 150 |
| 8 | 170 |
| 9 | 190 |
| 10 | 210 |
| 11 | 230 |
| 12 | 250 |
| 13 | 270 |
| 14 | 290 |
| 15 | 310 |
| 16 | 330 |
| 17 | 350 |
| 18 | 370 |
| 19 | 390 |
| 20 | 410 |

> If a kingdom bonus applies (e.g., Veridor +X% Wood), the bonus is applied on top of this base rate.

### A4. Warehouse Capacity Curve

**Formula**: `capacity = base_storage + (level × storage_per_level)` where `base_storage = 1000`, `storage_per_level = 500`.

| Level | Max Storage per Resource |
|-------|------------------------|
| 0 | 1000 |
| 1 | 1500 |
| 2 | 2000 |
| 3 | 2500 |
| 4 | 3000 |
| 5 | 3500 |
| 6 | 4000 |
| 7 | 4500 |
| 8 | 5000 |
| 9 | 5500 |
| 10 | 6000 |
| 11 | 6500 |
| 12 | 7000 |
| 13 | 7500 |
| 14 | 8000 |
| 15 | 8500 |
| 16 | 9000 |
| 17 | 9500 |
| 18 | 10000 |
| 19 | 10500 |
| 20 | 11000 |
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
| Town Hall | `town_hall` | Central building; gates access to all others | 20 | None | 0 |
| Iron Mine | `iron_mine` | Produces Iron per hour | 20 | None | 0 |
| Lumber Mill | `lumber_mill` | Produces Wood per hour | 20 | None | 0 |
| Quarry | `quarry` | Produces Stone per hour | 20 | None | 0 |
| Farm | `farm` | Produces Food per hour; determines population cap | 20 | None | 0 |
| Warehouse | `warehouse` | Stores resources; level sets max storage | 20 | None | 0 |
| Barracks | `barracks` | Trains infantry troops | 20 | town_hall 3 | 0 |
| Stable | `stable` | Trains mounted/fast troops | 15 | town_hall 5, barracks 5 | 0 |
| Forge | `forge` | Crafts weapons from resources + runes | 10 | town_hall 5, barracks 3 | 0 |
| Rune Altar | `rune_altar` | Combines, enhances, and stores runes | 10 | town_hall 7, forge 3 | 0 |
| Walls | `walls` | Passive defense bonus for the village | 20 | town_hall 2 | 0 |
| Marketplace | `marketplace` | Trade resources with other players | 15 | town_hall 5, warehouse 3 | 0 |
| Embassy | `embassy` | Required to form/join alliances | 10 | town_hall 8 | 0 |
| Watchtower | `watchtower` | Detects incoming attacks | 10 | town_hall 3, walls 1 | 0 |

> **Prerequisites format**: `building_id level` (e.g., `town_hall 5, barracks 3`). Use `None` if no prerequisites.

### B2. Base Costs (Level 1) — All Buildings

| Building | Iron | Wood | Stone | Food | Build Time | Scaling Factor | Time Factor |
|----------|------|------|-------|------|-----------|----------------|-------------|
| Town Hall | 200 | 200 | 200 | 100 | 5 min | 1.7 | 1.7 |
| Iron Mine | 100 | 80 | 50 | 30 | 2 min | 1.5 | 1.5 |
| Lumber Mill | 80 | 100 | 50 | 30 | 2 min | 1.5 | 1.5 |
| Quarry | 80 | 50 | 100 | 30 | 2 min | 1.5 | 1.5 |
| Farm | 50 | 80 | 50 | 20 | 2 min | 1.5 | 1.5 |
| Warehouse | 120 | 120 | 100 | 50 | 3 min | 1.6 | 1.6 |
| Barracks | 200 | 150 | 100 | 80 | 5 min | 1.8 | 1.8 |
| Stable | 300 | 200 | 150 | 120 | 8 min | 1.8 | 1.8 |
| Forge | 250 | 180 | 200 | 100 | 8 min | 1.8 | 1.8 |
| Rune Altar | 300 | 250 | 250 | 150 | 10 min | 1.9 | 1.9 |
| Walls | 150 | 100 | 200 | 50 | 4 min | 1.6 | 1.6 |
| Marketplace | 180 | 180 | 120 | 80 | 5 min | 1.6 | 1.6 |
| Embassy | 200 | 200 | 200 | 100 | 8 min | 1.7 | 1.7 |
| Watchtower | 150 | 100 | 150 | 60 | 4 min | 1.6 | 1.6 |

> **Cost formula**: `cost(level) = base_cost × (scaling_factor ^ (level - 1))`
> **Time formula**: `time(level) = base_time × (time_factor ^ (level - 1))`

### B3. Construction Rules

| Rule | Value |
|------|-------|
| Max simultaneous constructions per village | 1 |
| Can this be upgraded (e.g., Town Hall unlocks parallel queues)? | Not in v1 (future consideration) |
| Demolition policy | Instant, no cooldown |
| Resources returned on demolition | 0% (no refund) |

---

## C. Buildings (Kingdom-Specific)

### C1. Veridor — Dock

| Property | Value |
|----------|-------|
| Canonical ID | `dock` |
| Function | Naval operations, sea-based trade and troop transport |
| Max Level | 15 |
| Prerequisites | town_hall 6 |
| Starting Level | 0 (not built) |

**Base Cost (Level 1)**

| Iron | Wood | Stone | Food | Build Time | Scaling Factor |
|------|------|-------|------|-----------|----------------|
| 250 | 300 | 150 | 100 | 8 min | 1.8 |

### C2. Sylvara — Grove Sanctum

| Property | Value |
|----------|-------|
| Canonical ID | `grove_sanctum` |
| Function | Nature magic, enhanced rune crafting and forest-based bonuses |
| Max Level | 15 |
| Prerequisites | town_hall 6 |
| Starting Level | 0 (not built) |

**Base Cost (Level 1)**

| Iron | Wood | Stone | Food | Build Time | Scaling Factor |
|------|------|-------|------|-----------|----------------|
| 200 | 300 | 200 | 150 | 8 min | 1.8 |

### C3. Arkazia — Colosseum

| Property | Value |
|----------|-------|
| Canonical ID | `colosseum` |
| Function | Gladiatorial combat, troop morale boost and elite unit training |
| Max Level | 15 |
| Prerequisites | town_hall 6 |
| Starting Level | 0 (not built) |

**Base Cost (Level 1)**

| Iron | Wood | Stone | Food | Build Time | Scaling Factor |
|------|------|-------|------|-----------|----------------|
| 300 | 200 | 300 | 100 | 8 min | 1.8 |

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
| 2026-03-05 | Filled sections A2–A4 (starting resources, production rates, warehouse capacity), B1–B3 (building definitions, base costs, construction rules), C1–C3 (kingdom-specific buildings) with draft values matching `config/buildings.go` implementation |

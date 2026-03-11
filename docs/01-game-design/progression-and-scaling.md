# Progression & Scaling

> How buildings, troops, weapons, and runes scale in cost, time, and power. Essential reference for balancing.
>
> **Authoritative values** for building costs, troop stats, and resource economy live in `server/internal/config/*.go` and are exported via the config codegen pipeline. This document describes design philosophy and formulas — for exact numbers, check the Go source.

---

## General Scaling Philosophy

- **Exponential cost, linear power**: Upgrading costs more each level, but the power gain per level is relatively linear. This prevents runaway snowballing.
- **Time as a balancer**: Higher-tier upgrades take significantly more real time. Active players don't leap too far ahead of casual players.
- **Lumber + Stone primary**: Building construction costs primarily Lumber and Stone (the two construction resources), with thematic Food/Water additions for certain buildings.
- **Diminishing returns at cap**: The last few levels of any building/troop upgrade are expensive relative to their benefit, preventing must-max-everything pressure.

---

## Building Scaling

### Cost Formula

```
cost(level) = base_cost × (scaling_factor ^ (level - 1))
time(level) = base_time × (time_factor ^ (level - 1))
```

### Base Costs — All Buildings (Level 1)

Values from `server/internal/config/buildings.go`:

| Building | Food | Water | Lumber | Stone | Build Time | Scaling | Time Factor | Max Level |
|----------|------|-------|--------|-------|-----------|---------|-------------|-----------|
| Town Hall | 0 | 0 | 250 | 250 | 5 min | 1.7 | 1.7 | 20 |
| Resource Fields (all 12) | 0 | 0 | 90 | 50 | 2 min | 1.5 | 1.5 | 20 |
| Barracks | 50 | 0 | 200 | 150 | 5 min | 1.8 | 1.8 | 20 |
| Stable | 40 | 60 | 250 | 200 | 8 min | 1.8 | 1.8 | 15 |
| Archery | 0 | 0 | 200 | 150 | 5 min | 1.8 | 1.8 | 15 |
| Workshop | 0 | 0 | 350 | 300 | 10 min | 1.8 | 1.8 | 15 |
| Special | 40 | 40 | 300 | 350 | 15 min | 1.8 | 1.8 | 15 |
| Storage | 0 | 0 | 150 | 120 | 3 min | 1.6 | 1.6 | 20 |
| Provisions | 40 | 0 | 120 | 100 | 3 min | 1.6 | 1.6 | 20 |
| Reservoir | 0 | 0 | 100 | 150 | 3 min | 1.6 | 1.6 | 20 |

> Resource field costs are identical for all 12 fields (food_1–3, water_1–3, lumber_1–3, stone_1–3). Military buildings cost primarily Lumber + Stone with small thematic Food/Water additions.

### Example: Resource Field Scaling (Lumber + Stone only)

| Level | Lumber | Stone | Time |
|-------|--------|-------|------|
| 1 | 90 | 50 | 2 min |
| 2 | 135 | 75 | 3 min |
| 3 | 203 | 113 | 5 min |
| 5 | 456 | 253 | 10 min |
| 10 | 3,462 | 1,923 | 1.3 hr |
| 15 | 26,301 | 14,612 | 10 hr |
| 20 | 199,796 | 111,009 | 3.3 days |

### Production Rates — All Resource Buildings

All 12 resource buildings share the same production formula. Values from `server/internal/config/resources.go`:

- **BaseResourceRate**: 1.0/s (passive, always present)
- **RatePerLevel**: 2.0/s per building level
- **Total rate for a resource**: BaseResourceRate + RatePerLevel × sum(all 3 building levels for that resource)

| Combined 3-building Level | Total Rate/s |
|--------------------------|-------------|
| 3 (all at 1) | 7.0/s |
| 6 (all at 2) | 13.0/s |
| 9 (all at 3) | 19.0/s |
| 15 (all at 5) | 31.0/s |
| 30 (all at 10) | 61.0/s |
| 45 (all at 15) | 91.0/s |
| 60 (all at 20) | 121.0/s |

### Storage Capacity

Three specialized storage buildings control resource caps:

| Building | Resources Stored | Base Cost (L1) |
|----------|-----------------|----------------|
| **Storage** | Lumber, Stone | 150 Lumber, 120 Stone |
| **Provisions** | Food | 40 Food, 120 Lumber, 100 Stone |
| **Reservoir** | Water | 100 Lumber, 150 Stone |

**Storage formula**: `BaseStorage(1200) + sum(storage_building_level × StoragePerLevel(400))`

| Storage Building Level | Max Capacity |
|-----------------------|-------------|
| 0 (no storage building) | 1,200 (base) |
| 1 | 1,600 |
| 5 | 3,200 |
| 10 | 5,200 |
| 15 | 7,200 |
| 20 | 9,200 |

> These capacities are per resource type. A village with Storage Lv10 can hold 5,200 Lumber and 5,200 Stone. Provisions and Reservoir work the same way for Food and Water respectively.

### Training Building Speed Bonus

Military buildings grant a training speed multiplier (applied to all troops trained in that building):

| Building Level | Speed Multiplier |
|---------------|-----------------|
| 1 | 1.0× (base) |
| 5 | 1.25× |
| 10 | 1.6× |
| 15 | 2.0× |
| 20 | 2.5× |

Interpolated linearly between these breakpoints.

---

## Troop Scaling

Troops do **not** have individual levels. Their effectiveness scales through:

1. **Building levels**: Higher Barracks/Stable/etc. level → unlock better troops + faster training.
2. **Weapons**: Equipping weapons adds combat bonuses.
3. **Runes**: Socketed in weapons for additional effects.
4. **Hero bonuses**: Hero skills can buff specific troop types.
5. **Research** (if added later): Tech tree upgrades for troop stats.

### Troop Cost Ranges

Troop costs vary by kingdom and building tier. General pattern (from `server/internal/config/troops.go`):

| Building | Tier | Approx Total Cost/Unit | Approx Train Time |
|----------|------|----------------------|-------------------|
| Barracks T1 | Level 1 req | 200–250 resources | 2 min |
| Barracks T4 | Level 8 req | 400–500 resources | 3.5 min |
| Stable T1 | Level 1 req | 300–350 resources | 2.5 min |
| Stable T3+ | Level 5+ req | 500–700 resources | 4–5 min |
| Archery T1 | Level 1 req | 200–250 resources | 2 min |
| Workshop T1 | Level 1 req | 300–400 resources | 3 min |
| Workshop T4 | Level 8 req | 600–800 resources | 5 min |
| Special T1 | Level 1 req | 350–450 resources | 3.5 min |
| Special T4–5 | Level 8–10 req | 600–900 resources | 5–6 min |

> Training times above are base times. Higher building levels reduce these via the speed multiplier.

### Troop Food Upkeep

Every troop type has a food_upkeep value (food consumed per hour per unit). Ranges from 1 (basic infantry) to 4 (heavy siege/elite). This creates army size limits based on food economy.

---

## Weapon Scaling

### Crafting Costs by Tier

| Tier | Food | Water | Lumber | Stone | Runes | Craft Time | Attack Bonus Range | Defense Bonus Range |
|------|------|-------|--------|-------|-------|------------|-------------------|-------------------|
| Common | 50 | 50 | 200 | 100 | 0 | 10 min | +5 to +15 | +3 to +10 |
| Rare | 150 | 150 | 600 | 300 | 1 Minor | 1 hr | +15 to +30 | +10 to +20 |
| Epic | 500 | 500 | 2,000 | 1,000 | 2 Major | 6 hr | +30 to +55 | +20 to +40 |
| Legendary | 2,000 | 2,000 | 8,000 | 4,000 | 3 Grand | 24 hr | +55 to +80 | +40 to +65 |
| Mythic | 7,500 | 7,500 | 30,000 | 15,000 | 5 Grand | 3 days | +80 to +120 | +65 to +95 |

### Weapon Durability

- Weapons lose 1 durability per battle.
- Durability by tier: Common (10), Rare (20), Epic (35), Legendary (50), Mythic (100).
- Repair cost = 10% of crafting cost per durability point.
- At 0 durability, the weapon breaks and is lost forever. Runes are salvaged (50% chance per rune).

---

## Rune Scaling

### Combination Requirements

| Target Rune | Input Required | Success Rate | Rune Altar Level |
|-------------|---------------|-------------|-----------------|
| Minor | 3 Fragments | 100% | 1 |
| Major | 3 Minor | 80% | 3 |
| Grand | 3 Major | 60% | 6 |
| Primordial | 3 Grand + special event | 40% | 10 |

- **Failed combinations**: Input runes are lost. This creates scarcity and drives trading.

### Rune Effect Scaling

Rune effects scale with rarity:

| Effect Type | Fragment | Minor | Major | Grand | Primordial |
|-------------|---------|-------|-------|-------|-----------|
| Attack Bonus | +2 | +5 | +12 | +25 | +50 |
| Defense Bonus | +2 | +5 | +12 | +25 | +50 |
| Speed Bonus | — | +5% | +10% | +15% | +25% |
| Resource Bonus | — | +3% | +7% | +12% | +20% |
| Special Effect | — | Minor | Moderate | Strong | Legendary |

---

## Weapons of Chaos Power Level

Weapons of Chaos are **far more powerful** than any player-crafted weapon:

| Property | Weapons of Chaos | Mythic (Best Player Weapon) |
|----------|-----------------|---------------------------|
| Attack Bonus | +200 to +300 | +80 to +120 |
| Defense Bonus | +150 to +250 | +65 to +95 |
| Special Effects | 2-3 legendary effects | 1 based on runes |
| Durability | Infinite | 100 |
| Rune Slots | 0 (built-in effects) | 5 |

This massive power gap is intentional — Weapons of Chaos are game-changing, which is why wielding them comes with severe debuffs.

---

## Endgame: Weapons of Order Crafting

### Requirements per Weapon of Order

| Resource | Amount | Notes |
|----------|--------|-------|
| Food | 200,000 | Alliance-pooled |
| Water | 200,000 | Alliance-pooled |
| Lumber | 500,000 | Alliance-pooled |
| Stone | 300,000 | Alliance-pooled |
| Primordial Runes | 5 | Extremely rare |
| Forge Level | 10 (max) | At least one alliance member |
| Crafting Time | 7 days | Cannot be sped up |

### Number Required

- To challenge Moraphys, the alliance needs to forge a number of Weapons of Order **equal to the number of Weapons of Chaos** configured for the game world.
- Multiple alliances can cooperate if they share a diplomacy pact.

---

## Game Round Duration

> **Draft — To Be Decided**

| Option | Duration | Pros | Cons |
|--------|---------|------|------|
| **Short** | 3 months | Fast-paced, multiple rounds per year, easier testing | May feel rushed for casual players |
| **Long** | 6 months | More strategic depth, slower pacing, bigger empires | Commitment fatigue, latecomers disadvantaged |

The final decision will be made during playtesting. The system should be configurable per game world.

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of progression and scaling |
| 2026-03-03 | Added all building base costs, production rates, warehouse capacity, utility building scaling |
| 2026-03-10 | Major update: Replaced draft values with actual implemented values from Go config. Updated building costs to Lumber+Stone primary design. Added storage buildings (Storage, Provisions, Reservoir). Updated resource economy to match resources.go constants. Updated weapon costs to use Food/Water/Lumber/Stone. Removed game-template.md superseded reference. |

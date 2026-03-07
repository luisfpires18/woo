# Progression & Scaling

> **Superseded**: Definitive tunable values (costs, rates, scaling factors, etc.) are in [`game-template.md`](game-template.md). Values below are **drafts** — when they conflict, `game-template.md` wins.

> How buildings, troops, weapons, and runes scale in cost, time, and power. Essential reference for balancing.

---

## General Scaling Philosophy

- **Exponential cost, linear power**: Upgrading costs more each level, but the power gain per level is relatively linear. This prevents runaway snowballing.
- **Time as a balancer**: Higher-tier upgrades take significantly more real time. Active players don't leap too far ahead of casual players.
- **Resource diversity**: Higher levels require all 4 resources in increasing proportions, forcing diversified economies.
- **Diminishing returns at cap**: The last few levels of any building/troop upgrade are expensive relative to their benefit, preventing must-max-everything pressure.

---

## Building Scaling

### Cost Formula (Draft)

```
cost(level) = base_cost × (scaling_factor ^ (level - 1))
time(level) = base_time × (time_factor ^ (level - 1))
```

| Parameter | Value | Notes |
|-----------|-------|-------|
| `scaling_factor` | 1.5 – 2.0 | Varies by building type. Resource buildings scale slower (1.5×), military buildings faster (1.8×). |
| `time_factor` | 1.4 – 1.7 | Construction time scaling |
| `max_level` | 20 (most buildings) | Town Hall caps at 20. Some buildings cap at 10 or 15. |

### Example: Food Field Scaling

| Level | Food | Water | Lumber | Stone | Time |
|-------|------|-------|--------|-------|------|
| 1 | 60 | 40 | 80 | 50 | 2 min |
| 2 | 90 | 60 | 120 | 75 | 3 min |
| 3 | 135 | 90 | 180 | 112 | 5 min |
| 5 | 304 | 202 | 405 | 253 | 12 min |
| 10 | 2,307 | 1,538 | 3,075 | 1,922 | 1.5 hr |
| 15 | 17,515 | 11,677 | 23,354 | 14,596 | 12 hr |
| 20 | 133,028 | 88,685 | 177,370 | 110,856 | 4 days |

> These are **draft values**. Will be tuned during playtesting. The principle is clear: early levels are quick and cheap, late levels are a major investment.

### Production Rates — All Resource Buildings

All 12 resource buildings (3 per resource type: food_1/2/3, water_1/2/3, lumber_1/2/3, stone_1/2/3) share the same production curve. Production rate = `BaseResourceRate(1.0) + RatePerLevel(2.0) × sum(levels of all 3 buildings for that resource)`:

| Level | Resources/Second |
|-------|------------------|
| 1 | 3 |
| 2 | 5 |
| 3 | 7 |
| 5 | 11 |
| 10 | 21 |
| 15 | 31 |
| 20 | 41 |

> Each resource has 3 building slots. The rate shown is for a single slot at that level. Total rate for a resource is the base rate (1.0/s) plus 2.0 × the combined levels of all 3 slots.

Production increases roughly linearly to keep late-game upgrades worthwhile but not game-breaking.

### Warehouse Capacity

| Level | Max Storage per Resource |
|-------|------------------------|
| 1 | 1,000 |
| 2 | 1,500 |
| 3 | 2,200 |
| 5 | 5,000 |
| 10 | 15,000 |
| 15 | 40,000 |
| 20 | 80,000 |

### Base Costs — All Buildings (Level 1)

| Building | Food | Water | Lumber | Stone | Build Time | Scaling Factor |
|----------|------|-------|--------|-------|-----------|----------------|
| Town Hall | 100 | 200 | 200 | 200 | 5 min | 1.7 |
| Resource Fields (all 12) | 60 | 40 | 80 | 50 | 2 min | 1.5 |
| Warehouse | 50 | 120 | 120 | 100 | 3 min | 1.6 |
| Barracks | 80 | 200 | 150 | 100 | 5 min | 1.8 |
| Stable | 120 | 300 | 200 | 150 | 8 min | 1.8 |
| Forge | 100 | 250 | 180 | 200 | 8 min | 1.8 |
| Rune Altar | 150 | 300 | 250 | 250 | 10 min | 1.9 |
| Walls | 50 | 150 | 100 | 200 | 4 min | 1.6 |
| Marketplace | 80 | 180 | 180 | 120 | 5 min | 1.6 |
| Embassy | 100 | 200 | 200 | 200 | 8 min | 1.7 |
| Watchtower | 60 | 150 | 100 | 150 | 4 min | 1.6 |
| Dock | 100 | 250 | 300 | 150 | 8 min | 1.8 |

> Apply the cost formula: `cost(level) = base_cost × (scaling_factor ^ (level - 1))`. Scaling factor varies per building (resource buildings 1.5×, military 1.8×, etc.).

### Utility Building Scaling

#### Barracks — Training Speed Bonus

| Level | Training Speed Bonus |
|-------|---------------------|
| 1 | Base (1.0×) |
| 5 | 1.25× |
| 10 | 1.6× |
| 15 | 2.0× |
| 20 | 2.5× |

#### Forge — Tier Unlocks

| Forge Level | Max Weapon Tier |
|-------------|----------------|
| 1 | Common |
| 3 | Rare |
| 5 | Epic |
| 8 | Legendary |
| 10 | Mythic |

#### Walls — Defense Bonus

| Level | Defense Bonus |
|-------|-------------|
| 1 | +5% |
| 5 | +15% |
| 10 | +30% |
| 15 | +50% |
| 20 | +75% |

#### Marketplace — Trade Capacity

| Level | Max Resources per Trade |
|-------|------------------------|
| 1 | 500 |
| 5 | 2,000 |
| 10 | 8,000 |
| 15 | 20,000 |

#### Watchtower — Detection Range

| Level | Warning Time | Detail |
|-------|-------------|--------|
| 1 | 30 sec | "Incoming attack" only |
| 3 | 1 min | Approximate troop count |
| 5 | 2 min | Troop types visible |
| 8 | 5 min | Exact troop count + types |
| 10 | 10 min | Full detail + origin village |

---

## Troop Scaling

Troops do **not** have individual levels. Their effectiveness scales through:

1. **Building levels**: Higher Barracks/Stable level → unlock better troops + faster training.
2. **Weapons**: Equipping weapons adds combat bonuses.
3. **Runes**: Socketed in weapons for additional effects.
4. **Hero bonuses**: Hero skills can buff specific troop types.
5. **Research** (if added later): Tech tree upgrades for troop stats.

### Training Cost Examples (Arkazia)

| Unit | Iron | Wood | Stone | Food | Training Time |
|------|------|------|-------|------|--------------|
| Iron Legionary | 80 | 40 | 30 | 50 | 3 min |
| Crossbowman | 60 | 70 | 20 | 40 | 4 min |
| Mountain Knight | 150 | 50 | 80 | 100 | 10 min |
| Shieldbearer | 100 | 30 | 120 | 60 | 8 min |
| Gladiator | 200 | 60 | 100 | 150 | 15 min |
| Battering Ram | 300 | 200 | 150 | 50 | 20 min |
| Mountain Scout | 40 | 30 | 20 | 30 | 2 min |

> Training times decrease with higher Barracks/Stable levels (approximately -5% per level).

---

## Weapon Scaling

### Crafting Costs by Tier

| Tier | Iron | Wood | Stone | Food | Runes | Craft Time | Attack Bonus Range | Defense Bonus Range |
|------|------|------|-------|------|-------|------------|-------------------|-------------------|
| Common | 200 | 150 | 100 | 50 | 0 | 10 min | +5 to +15 | +3 to +10 |
| Rare | 600 | 450 | 300 | 150 | 1 Minor | 1 hr | +15 to +30 | +10 to +20 |
| Epic | 2,000 | 1,500 | 1,000 | 500 | 2 Major | 6 hr | +30 to +55 | +20 to +40 |
| Legendary | 8,000 | 6,000 | 4,000 | 2,000 | 3 Grand | 24 hr | +55 to +80 | +40 to +65 |
| Mythic | 30,000 | 22,500 | 15,000 | 7,500 | 5 Grand | 3 days | +80 to +120 | +65 to +95 |

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

- **Sylvara bonus**: Grove Sanctum doubles success rate (e.g., Grand becomes 100% instead of 60%).
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
| Iron | 500,000 | Alliance-pooled |
| Wood | 400,000 | Alliance-pooled |
| Stone | 300,000 | Alliance-pooled |
| Food | 200,000 | Alliance-pooled |
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
| 2026-03-03 | Added all building base costs, production rates for all resource buildings, Warehouse capacity table, utility building scaling (Barracks, Forge tiers, Walls defense, Marketplace, Watchtower) — re-applied after previous silent failure |
| 2026-03-05 | Marked as superseded by `game-template.md` for all tunable values |

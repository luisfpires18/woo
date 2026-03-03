# Core Mechanics

> All fundamental game systems. Read this before implementing any gameplay feature.

---

## Resources

The game has **four base resources**. Every village produces all four, but production rates depend on buildings and kingdom bonuses.

| Resource | Produced By | Primary Uses |
|----------|------------|-------------|
| **Iron** | Iron Mine | Weapons, troop equipment, advanced buildings |
| **Wood** | Lumber Mill | Buildings, siege equipment, ships (Veridor) |
| **Stone** | Quarry | Fortifications, walls, heavy structures |
| **Food** | Farm | Troop upkeep, population sustenance, army size cap |

### Resource Rules

- Resources accumulate over time based on building production rates.
- **Lazy calculation**: Do NOT update resources in the DB every tick. Store `last_updated` timestamp and calculate current value on read: `current = stored + (rate_per_hour × hours_elapsed)`. Write to DB only on events (build, trade, attack, etc.).
- Each village has a **Warehouse** that caps resource storage. Overflow is lost.
- **Food** is special: it is consumed by troops and population. If food production < consumption, troops start dying (starvation mechanic).
- Resources can be traded at the Marketplace between players.
- Resources can be raided from enemy villages via attacks.

---

## Buildings

Buildings are constructed inside a village and provide various functions. Each building has levels (starting at 0 = not built, max level TBD during balancing).

### Building Types

| Building | Function | Unlocks |
|----------|---------|---------|
| **Town Hall** | Central building. Its level determines what other buildings can be built. | All other buildings |
| **Iron Mine** | Produces Iron per hour. Higher level = more production. | — |
| **Lumber Mill** | Produces Wood per hour. | — |
| **Quarry** | Produces Stone per hour. | — |
| **Farm** | Produces Food per hour. Determines population cap. | — |
| **Warehouse** | Stores resources. Level determines max storage per resource. | — |
| **Barracks** | Trains infantry troops. Higher level = faster training, more unit types. | Troop types by level |
| **Stable** | Trains mounted/fast troops. | Advanced troop types |
| **Forge** | Crafts weapons from resources + runes. Higher level = higher weapon tiers. | Weapon tiers |
| **Rune Altar** | Combines, enhances, and stores runes. | Rune combinations |
| **Walls** | Passive defense bonus for the village. Higher level = stronger defense. | — |
| **Marketplace** | Trade resources with other players. Level affects trade capacity. | — |
| **Embassy** | Required to form/join alliances. Level affects alliance size. | Alliance features |
| **Watchtower** | Detects incoming attacks. Higher level = earlier warning + more detail. | — |
| **Dock** | (Veridor-only) Builds naval units. Enables sea-based attacks and trade routes. | Naval troops |

### Construction Rules

- Only one building can be under construction at a time per village (upgradeable via Town Hall to allow parallel queues).
- Construction requires resources and time. Both scale with building level.
- Buildings have **prerequisites**: e.g., Barracks requires Town Hall level 3, Forge requires Town Hall level 5 + Barracks level 3.
- Destroying (demolishing) a building is instant but returns zero resources.

---

## Troops

Each kingdom has a unique set of troop types. See `docs/01-game-design/kingdoms.md` for kingdom-specific unit rosters.

### Universal Troop Properties

| Property | Description |
|----------|------------|
| **Attack** | Offensive power in combat |
| **Defense (Infantry)** | Defense against infantry attacks |
| **Defense (Cavalry)** | Defense against cavalry/mounted attacks |
| **Speed** | Movement speed on the world map (tiles per hour) |
| **Carry Capacity** | Max resources this unit can carry when raiding |
| **Food Upkeep** | Food consumed per hour per unit |
| **Training Time** | Base time to train one unit |
| **Training Cost** | Resource cost per unit (Iron, Wood, Stone, Food) |

### Troop Actions

- **Attack**: Send troops to an enemy village. Combat resolves on arrival.
- **Raid**: Attack focused on stealing resources, not destroying buildings.
- **Defend**: Station troops in your own or an ally's village.
- **Scout**: Send a scout to reveal enemy village details (troop count, resources, buildings).
- **Reinforce**: Send troops to an ally's village as additional defense.

### Combat Resolution

Combat uses a **point-based system**:

1. Attacker's total attack power = sum of all attacking troops' attack values + weapon bonuses + hero bonuses.
2. Defender's total defense power = sum of all defending troops' relevant defense values (infantry vs infantry, cavalry vs cavalry) + wall bonus + weapon bonuses.
3. The side with higher total power wins. Losses are proportional to the ratio of powers.
4. Detailed formula TBD during Phase 3 balancing. Will be documented here when finalized.

---

## Weapons

Weapons are the game's distinguishing mechanic. They are crafted in Forges, enhanced with Runes, and equipped on troops or heroes.

### Weapon Tiers

| Tier | Forge Level Required | Runes Required | Power Level |
|------|---------------------|---------------|-------------|
| **Common** | 1 | 0 | Low |
| **Rare** | 3 | 1 | Medium |
| **Epic** | 5 | 2 | High |
| **Legendary** | 8 | 3 | Very High |
| **Mythic** | 10 | 5 (rare+) | Extreme |

### Weapon Properties

- **Type**: Sword, Axe, Bow, Spear, Shield, Staff (more TBD)
- **Attack Bonus**: Added to equipped troop/hero attack
- **Defense Bonus**: Added to equipped troop/hero defense
- **Special Effect**: Derived from socketed runes (e.g., "Burns target for 5% extra damage", "Heals wielder 2% per hit")
- **Rune Slots**: Number of runes that can be socketed (based on weapon tier)
- **Durability**: Weapons degrade over battles. Must be repaired at the Forge.

### Crafting

1. Player selects weapon type at the Forge.
2. Pays resource cost (scales with tier).
3. Optionally sockets runes during or after crafting.
4. Crafting takes time (scales with tier).
5. Result: a weapon with randomized bonus stats within a range for the tier.

---

## Runes

Runes are magical artifacts that modify weapon stats and grant special abilities.

### How to Obtain Runes

- **Exploration**: Discover rune fragments on unexplored map tiles.
- **Combat Drops**: Chance to find runes after winning battles.
- **Rune Altar**: Combine lower-tier rune fragments into complete runes.
- **Trading**: Buy/sell runes via the Marketplace.
- **Quests/Events**: World events and seasonal quests may reward runes.

### Rune Rarity

| Rarity | Drop Rate | Power |
|--------|----------|-------|
| **Fragment** | Common | Must combine 3 to make a Minor rune |
| **Minor** | Uncommon | Small stat bonus or minor effect |
| **Major** | Rare | Significant stat bonus or notable effect |
| **Grand** | Very Rare | Powerful effect, required for Legendary+ weapons |
| **Primordial** | Ultra Rare | Required for Weapons of Order crafting |

### Rune Effects (Examples)

- **Rune of Fire**: +10% attack, burns target for DOT
- **Rune of Iron**: +15% defense, reduces incoming damage
- **Rune of Swiftness**: +20% troop speed when equipped on hero
- **Rune of Harvest**: +10% resource production in wielder's village
- **Rune of Chaos**: Massive power boost but random negative events

---

## Weapons of Chaos

These are pre-existing, map-spawned legendary weapons of immense power. They are **NOT crafted by players** — they exist in the world from the start of each game round.

### Rules

- A fixed number of Weapons of Chaos exist per game world (e.g., 7).
- They are located in special "Chaos Shrine" map tiles, guarded by NPC defenders.
- Any player can send troops to claim one by defeating the NPC guardians.
- **Wielding a Weapon of Chaos grants enormous combat bonuses** but also:
  - Random resource decay (lose X% of a random resource periodically)
  - Betrayal events (small chance each day that some of your troops desert)
  - Chaos storms in your territory (reduce neighbor production)
  - Your village location is revealed to all players (no fog of war protection)
  - Increased aggression from Moraphys NPC faction
- A Weapon of Chaos can be **stolen** by another player attacking the wielder's village.
- They **cannot be destroyed** by players — only transferred.
- Wielders can voluntarily drop a Weapon of Chaos, returning it to a random Chaos Shrine.

### Endgame Trigger

When **Moraphys** (NPC faction) successfully gathers **all Weapons of Chaos**, the endgame event begins:
- Moraphys announces dominion over the world
- A countdown timer starts (TBD: days/weeks)
- Players must forge Weapons of Order to challenge Moraphys before the timer expires
- If the timer expires without Moraphys being defeated, the round ends in darkness (no winner)

---

## Weapons of Order

The ultimate crafted artifacts. Created by players working together to counter the Weapons of Chaos and defeat Moraphys.

### Crafting Requirements

- **Alliance-level collaboration**: Requires multiple alliance members contributing resources, runes, and forge capacity.
- **Primordial Runes**: Each Weapon of Order requires multiple Primordial runes (ultra rare).
- **Massive resource cost**: The entire alliance pools resources.
- **High-level Forge**: At least one alliance member needs a max-level Forge.
- **Time**: Crafting takes significant time (days of in-game time, TBD).

### How They Work

- Weapons of Order **neutralize** the power of Weapons of Chaos during the endgame battle.
- They provide massive combat bonuses specifically against Moraphys forces.
- The final battle is an alliance-coordinated attack on Moraphys's stronghold.
- Victory condition: defeat Moraphys while they hold all Weapons of Chaos.

---

## World Map

The multiplayer world is a **tile-based grid map** (similar to Travian's coordinate system).

### Map Properties

- **Grid**: Square or hexagonal tiles (TBD during Phase 1 implementation). Each tile has (x, y) coordinates.
- **Terrain Types**: Plains, Forest, Mountain, Water, Desert, Swamp. Terrain affects movement speed and resource availability.
- **Village Tiles**: Where player villages are located.
- **Chaos Shrine Tiles**: Where Weapons of Chaos are guarded.
- **Moraphys Stronghold**: Central/special tile where Moraphys resides. Grows in power over the game round.
- **Fog of War**: Players can only see tiles near their villages + allied territory. Scouting reveals more.
- **Oases**: Special tiles that provide resource bonuses to adjacent villages.

### Movement

- Troops move across the map from tile to tile.
- Movement time = distance / troop speed (slowest unit in the group).
- Terrain modifies movement speed (e.g., mountains = slower, plains = normal, roads = faster).

---

## Alliances

Players can form alliances for cooperative play.

### Features

- **Shared Map Vision**: Alliance members share fog of war.
- **Reinforcements**: Send troops to defend ally villages.
- **Resource Trading**: Reduced marketplace fees between allies.
- **Alliance Chat**: Dedicated communication channel.
- **Coordinated Attacks**: Plan attacks with timed waves.
- **Weapons of Order**: Alliance-level crafting (endgame requirement).

### Limits

- Alliance size is limited by the leader's Embassy level.
- A player can only be in one alliance at a time.
- Alliance diplomacy: alliances can declare war, peace, or NAP (Non-Aggression Pact) with other alliances.

---

## Hero System (Future)

> Planned for Phase 3-4. Placeholder for design reference.

- Each player has one hero, tied to their kingdom.
- Heroes lead armies, equip weapons and runes, and gain XP from battles.
- Heroes have a skill tree with kingdom-specific abilities.
- Heroes can die in battle (revive at village after cooldown).

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of core mechanics |

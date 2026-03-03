# Data Models

> Entity schemas for the database. Agents use these to create SQL migrations and Go model structs.

---

## Conventions

- All tables use `id` as primary key (INTEGER AUTOINCREMENT for SQLite, SERIAL for PostgreSQL).
- Timestamps use ISO 8601 format stored as TEXT (SQLite) or TIMESTAMPTZ (PostgreSQL).
- Foreign keys are enforced (`PRAGMA foreign_keys = ON` in SQLite).
- Soft deletes are NOT used — if a record is removed, it's deleted. Historical data uses separate audit/log tables if needed.
- JSON columns use TEXT in SQLite (parsed in Go), JSONB in PostgreSQL.

---

## Entity Relationship Overview

```
players ──1:N──► villages ──1:N──► buildings
   │                 │
   │                 ├──1:1──► resources
   │                 │
   │                 ├──1:N──► troops
   │                 │
   │                 ├──1:1──► forge (building, but with extra state)
   │                 │
   │                 └──1:N──► building_queue   │                 │
   │                 └──1:N──► training_queue   │
   ├──1:N──► weapons
   │
   ├──1:N──► runes
   │
   └──1:N──► alliance_members ──N:1──► alliances

world_map (grid of tiles, some linked to villages)

attacks (links attacker village → defender village, with troops and timing)

weapons_of_chaos (special singleton weapons on the map)
```

---

## Table Schemas

### players

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Unique player ID |
| username | TEXT | NOT NULL, UNIQUE | Display name |
| email | TEXT | NOT NULL, UNIQUE | Login email |
| password_hash | TEXT | NULL | Bcrypt hash (NULL if OAuth-only) |
| kingdom | TEXT | NOT NULL | 'veridor', 'sylvara', or 'arkazia' |
| oauth_provider | TEXT | NULL | 'google', 'discord', or NULL |
| oauth_id | TEXT | NULL | Provider-specific user ID |
| created_at | TEXT | NOT NULL | ISO 8601 timestamp |
| last_login_at | TEXT | NULL | Last login timestamp |

**Indexes**: `username`, `email`, `(oauth_provider, oauth_id)`

---

### villages

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Unique village ID |
| player_id | INTEGER | FK → players.id, NOT NULL | Owner |
| name | TEXT | NOT NULL | Village name (player-chosen) |
| x | INTEGER | NOT NULL | World map X coordinate |
| y | INTEGER | NOT NULL | World map Y coordinate |
| is_capital | INTEGER | NOT NULL, DEFAULT 0 | 1 if this is the player's capital |
| created_at | TEXT | NOT NULL | ISO 8601 timestamp |

**Indexes**: `player_id`, `(x, y)` UNIQUE

---

### buildings

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Unique building ID |
| village_id | INTEGER | FK → villages.id, NOT NULL | Village this building belongs to |
| building_type | TEXT | NOT NULL | Building type (e.g., 'town_hall', 'iron_mine', 'barracks'). See `docs/01-game-design/core-mechanics.md` for canonical constants. |
| level | INTEGER | NOT NULL, DEFAULT 0 | Current level (0 = not built) |

**Indexes**: `village_id`, `(village_id, building_type)` UNIQUE

---

### building_queue

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Queue entry ID |
| village_id | INTEGER | FK → villages.id, NOT NULL | Village |
| building_type | TEXT | NOT NULL | Which building is being upgraded |
| target_level | INTEGER | NOT NULL | Level being upgraded to |
| started_at | TEXT | NOT NULL | When construction started |
| completes_at | TEXT | NOT NULL | When construction finishes |

**Indexes**: `village_id`, `completes_at`

---

### resources

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| village_id | INTEGER | PK, FK → villages.id | One row per village |
| iron | REAL | NOT NULL, DEFAULT 0 | Current stored iron |
| wood | REAL | NOT NULL, DEFAULT 0 | Current stored wood |
| stone | REAL | NOT NULL, DEFAULT 0 | Current stored stone |
| food | REAL | NOT NULL, DEFAULT 0 | Current stored food |
| iron_rate | REAL | NOT NULL, DEFAULT 0 | Iron production per hour |
| wood_rate | REAL | NOT NULL, DEFAULT 0 | Wood production per hour |
| stone_rate | REAL | NOT NULL, DEFAULT 0 | Stone production per hour |
| food_rate | REAL | NOT NULL, DEFAULT 0 | Food production per hour (gross) |
| food_consumption | REAL | NOT NULL, DEFAULT 0 | Food consumed per hour by troops/population |
| max_storage | REAL | NOT NULL, DEFAULT 1000 | Max storage per resource (from Warehouse level) |
| last_updated | TEXT | NOT NULL | Timestamp of last resource write |

**Note**: Current resources are calculated lazily: `current = min(stored + (rate × hours_since_last_updated), max_storage)`. Net food rate = `food_rate - food_consumption`. The `stored` values are only written to DB on events (build, trade, attack, login). The `max_storage` value is updated when the Warehouse is upgraded.

---

### troops

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Unique troop stack ID |
| village_id | INTEGER | FK → villages.id, NOT NULL | Home village |
| type | TEXT | NOT NULL | Unit type (e.g., 'iron_legionary', 'wave_rider') |
| quantity | INTEGER | NOT NULL, DEFAULT 0 | Number of this unit type |
| status | TEXT | NOT NULL, DEFAULT 'stationed' | 'stationed', 'marching', 'defending', 'returning' |

**Indexes**: `village_id`, `(village_id, type)` UNIQUE  
**Note**: Troops in transit are tracked by the `attacks` table. Village troops represent what's currently stationed.

---

### training_queue

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Queue entry ID |
| village_id | INTEGER | FK → villages.id, NOT NULL | Village |
| troop_type | TEXT | NOT NULL | Unit type being trained (e.g., 'iron_legionary') |
| quantity | INTEGER | NOT NULL | Number of units in this batch |
| started_at | TEXT | NOT NULL | When training started |
| completes_at | TEXT | NOT NULL | When training finishes |

**Indexes**: `village_id`, `completes_at`

---

### weapons

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Unique weapon ID |
| player_id | INTEGER | FK → players.id, NOT NULL | Owner |
| name | TEXT | NOT NULL | Weapon name (generated or player-named) |
| weapon_type | TEXT | NOT NULL | 'sword', 'axe', 'bow', 'spear', 'shield', 'staff' |
| tier | TEXT | NOT NULL | 'common', 'rare', 'epic', 'legendary', 'mythic' |
| attack_bonus | INTEGER | NOT NULL, DEFAULT 0 | Attack stat bonus |
| defense_bonus | INTEGER | NOT NULL, DEFAULT 0 | Defense stat bonus |
| rune_slots | INTEGER | NOT NULL, DEFAULT 0 | Max runes that can be socketed |
| durability | INTEGER | NOT NULL | Current durability |
| max_durability | INTEGER | NOT NULL | Max durability |
| equipped_on | TEXT | NULL | What this weapon is equipped on (troop type, hero, etc.) |
| stats_json | TEXT | NULL | Additional stats/effects as JSON |
| created_at | TEXT | NOT NULL | ISO 8601 timestamp |

**Indexes**: `player_id`, `tier`

---

### runes

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Unique rune ID |
| player_id | INTEGER | FK → players.id, NOT NULL | Owner |
| rune_type | TEXT | NOT NULL | Effect type (e.g., 'fire', 'iron', 'swiftness', 'harvest') |
| rarity | TEXT | NOT NULL | 'fragment', 'minor', 'major', 'grand', 'primordial' |
| effect_json | TEXT | NOT NULL | JSON describing the rune's specific effects |
| weapon_id | INTEGER | FK → weapons.id, NULL | If socketed in a weapon, which one |
| created_at | TEXT | NOT NULL | ISO 8601 timestamp |

**Indexes**: `player_id`, `weapon_id`, `rarity`

---

### alliances

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Alliance ID |
| name | TEXT | NOT NULL, UNIQUE | Alliance name |
| tag | TEXT | NOT NULL, UNIQUE | Short tag (3-5 chars) |
| leader_id | INTEGER | FK → players.id, NOT NULL | Alliance leader |
| max_members | INTEGER | NOT NULL, DEFAULT 10 | Max size (based on leader's Embassy level) |
| created_at | TEXT | NOT NULL | ISO 8601 timestamp |

---

### alliance_members

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| alliance_id | INTEGER | FK → alliances.id, NOT NULL | Alliance |
| player_id | INTEGER | FK → players.id, NOT NULL, UNIQUE | Player (can only be in one alliance) |
| role | TEXT | NOT NULL, DEFAULT 'member' | 'leader', 'officer', 'member' |
| joined_at | TEXT | NOT NULL | ISO 8601 timestamp |

**Primary Key**: `(alliance_id, player_id)`

---

### world_map

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| x | INTEGER | NOT NULL | X coordinate |
| y | INTEGER | NOT NULL | Y coordinate |
| terrain_type | TEXT | NOT NULL | 'plains', 'forest', 'mountain', 'water', 'desert', 'swamp' |
| owner_player_id | INTEGER | FK → players.id, NULL | Player who controls this tile (NULL if unclaimed) |
| village_id | INTEGER | FK → villages.id, NULL | Village on this tile (NULL if no village) |
| has_oasis | INTEGER | NOT NULL, DEFAULT 0 | Whether this tile is an oasis (resource bonus) |
| has_chaos_shrine | INTEGER | NOT NULL, DEFAULT 0 | Whether this tile has a Chaos Shrine |

**Primary Key**: `(x, y)`  
**Indexes**: `owner_player_id`, `village_id`, `terrain_type`

---

### attacks

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Attack ID |
| attacker_player_id | INTEGER | FK → players.id, NOT NULL | Attacker |
| attacker_village_id | INTEGER | FK → villages.id, NOT NULL | Origin village |
| target_x | INTEGER | NOT NULL | Target tile X |
| target_y | INTEGER | NOT NULL | Target tile Y |
| attack_type | TEXT | NOT NULL | 'attack', 'raid', 'scout', 'reinforce' |
| troops_json | TEXT | NOT NULL | JSON: { "iron_legionary": 50, "crossbowman": 20, ... } |
| weapons_json | TEXT | NULL | JSON: weapons carried by hero/troops |
| departed_at | TEXT | NOT NULL | When troops left the village |
| arrives_at | TEXT | NOT NULL | When troops arrive at target |
| status | TEXT | NOT NULL, DEFAULT 'marching' | 'marching', 'arrived', 'returning', 'completed' |
| result_json | TEXT | NULL | Combat result (filled after resolution) |

**Indexes**: `attacker_player_id`, `arrives_at`, `status`

---

### weapons_of_chaos

> **Configurable per game world.** The number and identity of Weapons of Chaos are set by game administrators at world creation. The schema supports any count.

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Weapon of Chaos ID |
| name | TEXT | NOT NULL, UNIQUE | e.g., 'Voidfang', 'The Tide Ender' |
| weapon_type | TEXT | NOT NULL | Type of weapon |
| attack_bonus | INTEGER | NOT NULL | Attack stat bonus |
| defense_bonus | INTEGER | NOT NULL | Defense stat bonus |
| effects_json | TEXT | NOT NULL | JSON: special effects and debuffs |
| location_x | INTEGER | NULL | Current map X (if at shrine) |
| location_y | INTEGER | NULL | Current map Y (if at shrine) |
| wielder_player_id | INTEGER | FK → players.id, NULL | Player currently wielding (NULL if at shrine or held by Moraphys) |
| held_by_moraphys | INTEGER | NOT NULL, DEFAULT 0 | 1 if Moraphys has stolen this weapon |
| claimed_at | TEXT | NULL | When it was last claimed |

**Indexes**: `wielder_player_id`

---

### refresh_tokens

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| id | INTEGER | PK, AUTOINCREMENT | Token ID |
| player_id | INTEGER | FK → players.id, NOT NULL | Owner |
| token_hash | TEXT | NOT NULL, UNIQUE | SHA-256 hash of the refresh token |
| expires_at | TEXT | NOT NULL | Expiry timestamp |
| created_at | TEXT | NOT NULL | Creation timestamp |

**Indexes**: `player_id`, `token_hash`

---

### schema_migrations

| Column | Type | Constraints | Description |
|--------|------|------------|-------------|
| version | INTEGER | PK | Migration number |
| applied_at | TEXT | NOT NULL | When the migration was run |

Used by the migration system to track which migrations have been applied.

---

## Notes for Implementation

1. **SQLite → PostgreSQL**: When migrating, change `INTEGER AUTOINCREMENT` to `SERIAL`, `TEXT` timestamps to `TIMESTAMPTZ`, and `TEXT` JSON columns to `JSONB`. The Go repository interface stays the same.
2. **Lazy resource calculation**: The `resources` table stores a snapshot. Actual current resources = `min(snapshot + (rate × elapsed_time), max_storage)`. Net food = `food_rate - food_consumption`. See `docs/06-database/database-guide.md`.
3. **Troop movement**: While troops are marching, they exist in the `attacks` table. The `troops` table only shows what's currently stationed in a village.
4. **Weapon durability**: Decremented after each battle. When 0, weapon is deleted (with rune salvage logic).
5. **Column naming**: All tables use `building_type` / `troop_type` (snake_case descriptive names). The `buildings` and `building_queue` tables both use `building_type` for consistency.

---

## Changelog

| Date | Change |
|------|--------|
| 2026-03-03 | Initial creation of data models |
| 2026-03-03 | Added max_storage and food_consumption to resources, added training_queue table, renamed buildings.type to building_type, made weapons_of_chaos configurable, added column naming note |

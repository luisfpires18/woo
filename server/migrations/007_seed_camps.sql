-- 007_seed_camps.sql — Seed data for beast templates, camp templates,
-- spawn rules, reward tables, and default battle tuning.

-- ── Beast Templates (Forest Biome) ──────────────────────────────────────────
-- Regular beasts (weaker, tier 1-3)
INSERT INTO beast_templates (name, sprite_key, hp, attack_power, attack_interval, defense_percent, crit_chance_percent, description, created_at, updated_at)
VALUES
    ('Forest Bee',      'enemies/forest/bee',      30,  8,  3,  5,   10,  'Small but aggressive stinging insect.',     datetime('now'), datetime('now')),
    ('Forest Beetle',   'enemies/forest/beetle',   50,  6,  4,  20,  5,   'Tough-shelled beetle with strong defenses.', datetime('now'), datetime('now')),
    ('Forest Spider',   'enemies/forest/spider',   40,  12, 2,  5,   15,  'Fast-striking venomous spider.',             datetime('now'), datetime('now')),
    ('Forest Snake',    'enemies/forest/snake',    35,  15, 3,  8,   20,  'Venomous serpent with deadly precision.',     datetime('now'), datetime('now')),
    ('Forest Monkey',   'enemies/forest/monkey',   45,  10, 2,  10,  12,  'Agile primate that attacks in groups.',      datetime('now'), datetime('now')),
    ('Forest Toucan',   'enemies/forest/toucan',   25,  7,  2,  3,   8,   'Colorful bird with a sharp beak.',           datetime('now'), datetime('now')),
    ('Forest Eagle',    'enemies/forest/eagle',    55,  14, 3,  12,  18,  'Majestic raptor with powerful talons.',       datetime('now'), datetime('now')),
    ('Forest Cheetah',  'enemies/forest/cheetah',  70,  18, 2,  15,  25,  'Swift predator that hunts with precision.',  datetime('now'), datetime('now')),
    ('Forest Stag',     'enemies/forest/stag',     80,  12, 3,  25,  5,   'Proud stag with hardened antlers.',           datetime('now'), datetime('now')),
    ('Forest Panther',  'enemies/forest/panther',  90,  20, 2,  18,  30,  'Stealthy jungle cat with lethal claws.',     datetime('now'), datetime('now'));

-- Boss beasts (stronger, tier 4-10)
INSERT INTO beast_templates (name, sprite_key, hp, attack_power, attack_interval, defense_percent, crit_chance_percent, description, created_at, updated_at)
VALUES
    ('Forest Bear',         'enemies/forest/boss_1_bear',      200, 25, 3, 30, 15, 'Massive bear — guardian of the forest.',       datetime('now'), datetime('now')),
    ('Forest Gorilla',      'enemies/forest/boss_2_gorilla',   250, 30, 3, 35, 10, 'Immense silver-back gorilla.',                 datetime('now'), datetime('now')),
    ('Forest Tiger',        'enemies/forest/boss_3_tiger',     300, 35, 2, 25, 30, 'Striped hunter with devastating attacks.',      datetime('now'), datetime('now')),
    ('Forest Crocodile',    'enemies/forest/boss_4_croc',      350, 28, 4, 45, 12, 'Ancient crocodile with near-impenetrable hide.', datetime('now'), datetime('now')),
    ('Forest Rhino',        'enemies/forest/boss_5_rhino',     400, 32, 4, 50, 8,  'Armored rhinoceros that charges relentlessly.', datetime('now'), datetime('now')),
    ('Forest Elephant',     'enemies/forest/boss_6_elephant',  500, 22, 5, 55, 5,  'Colossal elephant with tremendous stamina.',    datetime('now'), datetime('now')),
    ('Forest Hydra',        'enemies/forest/boss_7_hydra',     450, 40, 3, 30, 25, 'Multi-headed serpent of the deep forest.',      datetime('now'), datetime('now')),
    ('Forest Wyrm',         'enemies/forest/boss_8_wyrm',      550, 45, 3, 35, 28, 'Wingless dragon-kin lurking in the canopy.',    datetime('now'), datetime('now')),
    ('Forest Ancient Treant','enemies/forest/boss_9_treant',   700, 20, 6, 65, 3,  'Sentient ancient tree spirit.',                 datetime('now'), datetime('now')),
    ('Forest Lion King',    'enemies/forest/boss_10_lion',     600, 50, 2, 40, 35, 'Legendary lion — apex predator of the forest.', datetime('now'), datetime('now'));

-- ── Reward Tables ────────────────────────────────────────────────────────────
-- Tier 1 rewards (small, for beginner camps)
INSERT INTO reward_tables (name, created_at, updated_at) VALUES ('tier_1_loot', datetime('now'), datetime('now'));
INSERT INTO reward_table_entries (reward_table_id, reward_type, min_amount, max_amount, drop_chance_pct, created_at)
VALUES
    ((SELECT id FROM reward_tables WHERE name = 'tier_1_loot'), 'food',   50,  150, 80, datetime('now')),
    ((SELECT id FROM reward_tables WHERE name = 'tier_1_loot'), 'lumber', 50,  150, 80, datetime('now')),
    ((SELECT id FROM reward_tables WHERE name = 'tier_1_loot'), 'stone',  30,  100, 60, datetime('now')),
    ((SELECT id FROM reward_tables WHERE name = 'tier_1_loot'), 'water',  30,  100, 60, datetime('now'));

-- Tier 2 rewards (medium)
INSERT INTO reward_tables (name, created_at, updated_at) VALUES ('tier_2_loot', datetime('now'), datetime('now'));
INSERT INTO reward_table_entries (reward_table_id, reward_type, min_amount, max_amount, drop_chance_pct, created_at)
VALUES
    ((SELECT id FROM reward_tables WHERE name = 'tier_2_loot'), 'food',   100, 300, 90, datetime('now')),
    ((SELECT id FROM reward_tables WHERE name = 'tier_2_loot'), 'lumber', 100, 300, 90, datetime('now')),
    ((SELECT id FROM reward_tables WHERE name = 'tier_2_loot'), 'stone',  80,  250, 75, datetime('now')),
    ((SELECT id FROM reward_tables WHERE name = 'tier_2_loot'), 'water',  80,  250, 75, datetime('now'));

-- Tier 3 rewards (large, for boss camps)
INSERT INTO reward_tables (name, created_at, updated_at) VALUES ('tier_3_loot', datetime('now'), datetime('now'));
INSERT INTO reward_table_entries (reward_table_id, reward_type, min_amount, max_amount, drop_chance_pct, created_at)
VALUES
    ((SELECT id FROM reward_tables WHERE name = 'tier_3_loot'), 'food',   200, 600,  95, datetime('now')),
    ((SELECT id FROM reward_tables WHERE name = 'tier_3_loot'), 'lumber', 200, 600,  95, datetime('now')),
    ((SELECT id FROM reward_tables WHERE name = 'tier_3_loot'), 'stone',  150, 500,  85, datetime('now')),
    ((SELECT id FROM reward_tables WHERE name = 'tier_3_loot'), 'water',  150, 500,  85, datetime('now'));

-- ── Camp Templates ───────────────────────────────────────────────────────────
-- Tier 1: Small camps with weak beasts
INSERT INTO camp_templates (name, tier, sprite_key, description, reward_table_id, created_at, updated_at)
VALUES ('Insect Nest', 1, 'camp_tier_1', 'A small nest of forest insects.', (SELECT id FROM reward_tables WHERE name = 'tier_1_loot'), datetime('now'), datetime('now'));

INSERT INTO camp_beast_slots (camp_template_id, beast_template_id, min_count, max_count)
VALUES
    ((SELECT id FROM camp_templates WHERE name = 'Insect Nest'), (SELECT id FROM beast_templates WHERE name = 'Forest Bee'),    2, 4),
    ((SELECT id FROM camp_templates WHERE name = 'Insect Nest'), (SELECT id FROM beast_templates WHERE name = 'Forest Beetle'), 1, 2);

-- Tier 2: Medium camps with mixed beasts
INSERT INTO camp_templates (name, tier, sprite_key, description, reward_table_id, created_at, updated_at)
VALUES ('Spider Den', 2, 'camp_tier_2', 'A dark den crawling with spiders and snakes.', (SELECT id FROM reward_tables WHERE name = 'tier_1_loot'), datetime('now'), datetime('now'));

INSERT INTO camp_beast_slots (camp_template_id, beast_template_id, min_count, max_count)
VALUES
    ((SELECT id FROM camp_templates WHERE name = 'Spider Den'), (SELECT id FROM beast_templates WHERE name = 'Forest Spider'), 2, 4),
    ((SELECT id FROM camp_templates WHERE name = 'Spider Den'), (SELECT id FROM beast_templates WHERE name = 'Forest Snake'),  1, 3);

-- Tier 3: Pack camps with agile beasts
INSERT INTO camp_templates (name, tier, sprite_key, description, reward_table_id, created_at, updated_at)
VALUES ('Monkey Troop', 3, 'camp_tier_3', 'A band of mischievous monkeys and eagles.', (SELECT id FROM reward_tables WHERE name = 'tier_2_loot'), datetime('now'), datetime('now'));

INSERT INTO camp_beast_slots (camp_template_id, beast_template_id, min_count, max_count)
VALUES
    ((SELECT id FROM camp_templates WHERE name = 'Monkey Troop'), (SELECT id FROM beast_templates WHERE name = 'Forest Monkey'), 3, 5),
    ((SELECT id FROM camp_templates WHERE name = 'Monkey Troop'), (SELECT id FROM beast_templates WHERE name = 'Forest Eagle'),  1, 2);

-- Tier 4: Predator camps
INSERT INTO camp_templates (name, tier, sprite_key, description, reward_table_id, created_at, updated_at)
VALUES ('Predator Lair', 4, 'camp_tier_4', 'A lair of dangerous predators.', (SELECT id FROM reward_tables WHERE name = 'tier_2_loot'), datetime('now'), datetime('now'));

INSERT INTO camp_beast_slots (camp_template_id, beast_template_id, min_count, max_count)
VALUES
    ((SELECT id FROM camp_templates WHERE name = 'Predator Lair'), (SELECT id FROM beast_templates WHERE name = 'Forest Cheetah'), 1, 3),
    ((SELECT id FROM camp_templates WHERE name = 'Predator Lair'), (SELECT id FROM beast_templates WHERE name = 'Forest Panther'), 1, 2),
    ((SELECT id FROM camp_templates WHERE name = 'Predator Lair'), (SELECT id FROM beast_templates WHERE name = 'Forest Stag'),    1, 2);

-- Tier 5: Boss camp
INSERT INTO camp_templates (name, tier, sprite_key, description, reward_table_id, created_at, updated_at)
VALUES ('Bear Cave', 5, 'camp_tier_5', 'A dark cave guarded by a massive bear.', (SELECT id FROM reward_tables WHERE name = 'tier_3_loot'), datetime('now'), datetime('now'));

INSERT INTO camp_beast_slots (camp_template_id, beast_template_id, min_count, max_count)
VALUES
    ((SELECT id FROM camp_templates WHERE name = 'Bear Cave'), (SELECT id FROM beast_templates WHERE name = 'Forest Bear'),    1, 1),
    ((SELECT id FROM camp_templates WHERE name = 'Bear Cave'), (SELECT id FROM beast_templates WHERE name = 'Forest Panther'), 2, 3);

-- ── Spawn Rules ──────────────────────────────────────────────────────────────
-- Forest wilderness spawn rule: spawns tier 1-3 camps in forest tiles
INSERT INTO spawn_rules (
    name, terrain_types_json, zone_types_json, camp_template_pool_json,
    max_camps, spawn_interval_sec, despawn_after_sec,
    min_camp_distance, min_village_distance, enabled, created_at, updated_at
) VALUES (
    'Wilderness Camps',
    '["plains"]',
    '["wilderness"]',
    '[' ||
        '{"camp_template_id":' || (SELECT id FROM camp_templates WHERE name = 'Insect Nest')   || ',"weight":40},' ||
        '{"camp_template_id":' || (SELECT id FROM camp_templates WHERE name = 'Spider Den')    || ',"weight":30},' ||
        '{"camp_template_id":' || (SELECT id FROM camp_templates WHERE name = 'Monkey Troop')  || ',"weight":20},' ||
        '{"camp_template_id":' || (SELECT id FROM camp_templates WHERE name = 'Predator Lair') || ',"weight":8},' ||
        '{"camp_template_id":' || (SELECT id FROM camp_templates WHERE name = 'Bear Cave')     || ',"weight":2}' ||
    ']',
    15,    -- max 15 camps from this rule
    60,    -- spawn check every 60 seconds
    3600,  -- camps despawn after 1 hour
    3,     -- minimum 3 tiles apart
    2,     -- minimum 2 tiles from villages
    1,     -- enabled
    datetime('now'),
    datetime('now')
);

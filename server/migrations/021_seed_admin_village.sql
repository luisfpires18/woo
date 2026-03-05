-- Seed a starting village for the admin user (who was seeded in 019 without one).
-- Uses (0, 0) as admin spawn — center of the map.

INSERT OR IGNORE INTO villages (player_id, name, x, y, is_capital)
SELECT id, 'Admin''s Village', 0, 0, 1
FROM players WHERE email = 'admin@woo.local';

-- Starter buildings (all level 0 = slot exists). Veridor kingdom → includes dock.
INSERT OR IGNORE INTO buildings (village_id, building_type, level)
SELECT v.id, b.type, 0
FROM villages v
JOIN (
    SELECT 'town_hall' AS type UNION ALL
    SELECT 'iron_mine' UNION ALL
    SELECT 'lumber_mill' UNION ALL
    SELECT 'quarry' UNION ALL
    SELECT 'farm' UNION ALL
    SELECT 'warehouse' UNION ALL
    SELECT 'barracks' UNION ALL
    SELECT 'stable' UNION ALL
    SELECT 'forge' UNION ALL
    SELECT 'rune_altar' UNION ALL
    SELECT 'walls' UNION ALL
    SELECT 'marketplace' UNION ALL
    SELECT 'embassy' UNION ALL
    SELECT 'watchtower' UNION ALL
    SELECT 'dock'
) b ON 1=1
WHERE v.player_id = (SELECT id FROM players WHERE email = 'admin@woo.local')
  AND v.is_capital = 1;

-- Starter resources
INSERT OR IGNORE INTO resources (village_id, iron, wood, stone, food, iron_rate, wood_rate, stone_rate, food_rate, food_consumption, max_storage)
SELECT v.id, 500, 500, 500, 500, 30, 30, 30, 30, 0, 1000
FROM villages v
WHERE v.player_id = (SELECT id FROM players WHERE email = 'admin@woo.local')
  AND v.is_capital = 1;

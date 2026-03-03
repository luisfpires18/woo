-- 005_create_resources.sql
CREATE TABLE IF NOT EXISTS resources (
    village_id INTEGER PRIMARY KEY REFERENCES villages(id),
    iron REAL NOT NULL DEFAULT 500,
    wood REAL NOT NULL DEFAULT 500,
    stone REAL NOT NULL DEFAULT 500,
    food REAL NOT NULL DEFAULT 500,
    iron_rate REAL NOT NULL DEFAULT 30,
    wood_rate REAL NOT NULL DEFAULT 30,
    stone_rate REAL NOT NULL DEFAULT 30,
    food_rate REAL NOT NULL DEFAULT 30,
    food_consumption REAL NOT NULL DEFAULT 0,
    max_storage REAL NOT NULL DEFAULT 1000,
    last_updated TEXT NOT NULL DEFAULT (datetime('now'))
);

-- 003_create_buildings.sql
CREATE TABLE IF NOT EXISTS buildings (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    village_id INTEGER NOT NULL REFERENCES villages(id),
    building_type TEXT NOT NULL,
    level INTEGER NOT NULL DEFAULT 0
);

CREATE INDEX IF NOT EXISTS idx_buildings_village_id ON buildings(village_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_buildings_village_type ON buildings(village_id, building_type);

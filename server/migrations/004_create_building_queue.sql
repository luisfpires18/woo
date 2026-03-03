-- 004_create_building_queue.sql
CREATE TABLE IF NOT EXISTS building_queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    village_id INTEGER NOT NULL REFERENCES villages(id),
    building_type TEXT NOT NULL,
    target_level INTEGER NOT NULL,
    started_at TEXT NOT NULL DEFAULT (datetime('now')),
    completes_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_building_queue_village_id ON building_queue(village_id);
CREATE INDEX IF NOT EXISTS idx_building_queue_completes_at ON building_queue(completes_at);

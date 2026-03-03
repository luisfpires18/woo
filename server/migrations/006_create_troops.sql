-- 006_create_troops.sql
CREATE TABLE IF NOT EXISTS troops (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    village_id INTEGER NOT NULL REFERENCES villages(id),
    type TEXT NOT NULL,
    quantity INTEGER NOT NULL DEFAULT 0,
    status TEXT NOT NULL DEFAULT 'stationed' CHECK (status IN ('stationed', 'marching', 'defending', 'returning'))
);

CREATE INDEX IF NOT EXISTS idx_troops_village_id ON troops(village_id);
CREATE UNIQUE INDEX IF NOT EXISTS idx_troops_village_type ON troops(village_id, type);

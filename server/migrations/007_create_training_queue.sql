-- 007_create_training_queue.sql
CREATE TABLE IF NOT EXISTS training_queue (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    village_id INTEGER NOT NULL REFERENCES villages(id),
    troop_type TEXT NOT NULL,
    quantity INTEGER NOT NULL,
    started_at TEXT NOT NULL DEFAULT (datetime('now')),
    completes_at TEXT NOT NULL
);

CREATE INDEX IF NOT EXISTS idx_training_queue_village_id ON training_queue(village_id);
CREATE INDEX IF NOT EXISTS idx_training_queue_completes_at ON training_queue(completes_at);

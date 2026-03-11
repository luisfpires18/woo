-- Per-resource storage caps: replace single max_storage with max_food, max_water, max_lumber, max_stone.
-- Each storage building type now only increases its corresponding resource cap.
-- Storage → lumber + stone, Provisions → food, Reservoir → water.

ALTER TABLE resources ADD COLUMN max_food   REAL NOT NULL DEFAULT 1200;
ALTER TABLE resources ADD COLUMN max_water  REAL NOT NULL DEFAULT 1200;
ALTER TABLE resources ADD COLUMN max_lumber REAL NOT NULL DEFAULT 1200;
ALTER TABLE resources ADD COLUMN max_stone  REAL NOT NULL DEFAULT 1200;

-- Migrate: copy existing max_storage value into all four new columns
UPDATE resources SET
    max_food   = max_storage,
    max_water  = max_storage,
    max_lumber = max_storage,
    max_stone  = max_storage;

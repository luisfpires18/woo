-- Migration 004: Add population tracking to resources table
-- pop_used tracks how much population is consumed by trained troops.
-- pop_cap is computed on read from buildings (not stored).

ALTER TABLE resources ADD COLUMN pop_used INTEGER NOT NULL DEFAULT 0;

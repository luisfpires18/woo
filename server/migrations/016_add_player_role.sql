-- Add role column to players table.
-- Default role is 'player'. Admin accounts get 'admin'.
ALTER TABLE players ADD COLUMN role TEXT NOT NULL DEFAULT 'player' CHECK (role IN ('player', 'admin'));

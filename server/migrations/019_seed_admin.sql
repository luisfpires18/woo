-- Seed admin accounts.
-- Password: Woo123! for all accounts
-- Hash: $2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC
-- Note: kingdom is left empty ('') so users see kingdom picker on first login
INSERT OR IGNORE INTO players (username, email, password_hash, kingdom, role, created_at)
VALUES 
  ('wright', 'wright@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
  ('coyote', 'coyote@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
  ('veridor', 'veridor@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
  ('sylvaine', 'sylvaine@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
  ('taldros', 'taldros@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
  ('blake', 'blake@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
  ('hawkes', 'hawkes@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now')),
  ('drakt', 'drakt@woo.local', '$2a$12$qh1rnZcYztG9ZjSA2xqOveVy4To13dV.eVRWnLs6Du0HhUJAj4zzC', '', 'admin', datetime('now'));

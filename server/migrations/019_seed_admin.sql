-- Seed default admin account.
-- ⚠️  CHANGE THIS PASSWORD IMMEDIATELY AFTER FIRST LOGIN!
-- Default credentials: admin@woo.local / Admin@WOO2026
INSERT OR IGNORE INTO players (username, email, password_hash, kingdom, role, created_at)
VALUES ('admin', 'admin@woo.local', '$2a$12$fBCDizCk9PMxYRIPtv7LjeA5qrEv4QtesIp0ZtXnA0BqWdQ4ZWt52', 'veridor', 'admin', datetime('now'));

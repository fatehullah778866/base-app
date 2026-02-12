ALTER TABLE users ADD COLUMN role TEXT DEFAULT 'user';
CREATE INDEX IF NOT EXISTS idx_users_role ON users(role);

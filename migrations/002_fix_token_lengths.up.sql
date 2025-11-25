-- Fix token column lengths to accommodate JWT tokens
-- JWT tokens can be 200-400+ characters, so VARCHAR(255) is insufficient

ALTER TABLE sessions 
  ALTER COLUMN token TYPE TEXT,
  ALTER COLUMN refresh_token TYPE TEXT;

-- Revert token column lengths
ALTER TABLE sessions 
  ALTER COLUMN token TYPE VARCHAR(255),
  ALTER COLUMN refresh_token TYPE VARCHAR(255);

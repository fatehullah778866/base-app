-- Rollback migration
DROP INDEX IF EXISTS idx_dashboard_items_created_at;
DROP INDEX IF EXISTS idx_dashboard_items_category;
DROP INDEX IF EXISTS idx_dashboard_items_status;
DROP INDEX IF EXISTS idx_dashboard_items_user_id;
DROP TABLE IF EXISTS dashboard_items;

DROP INDEX IF EXISTS idx_user_settings_username;
DROP TABLE IF EXISTS user_settings_comprehensive;

DROP INDEX IF EXISTS idx_password_reset_tokens_expires;
DROP INDEX IF EXISTS idx_password_reset_tokens_token;
DROP INDEX IF EXISTS idx_password_reset_tokens_user_id;
DROP TABLE IF EXISTS password_reset_tokens;


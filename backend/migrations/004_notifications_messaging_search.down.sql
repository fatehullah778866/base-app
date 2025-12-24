-- Rollback migration
DROP TRIGGER IF EXISTS messages_fts_update;
DROP TRIGGER IF EXISTS messages_fts_delete;
DROP TRIGGER IF EXISTS messages_fts_insert;
DROP TRIGGER IF EXISTS dashboard_items_fts_update;
DROP TRIGGER IF EXISTS dashboard_items_fts_delete;
DROP TRIGGER IF EXISTS dashboard_items_fts_insert;

DROP TABLE IF EXISTS messages_fts;
DROP TABLE IF EXISTS dashboard_items_fts;

DROP INDEX IF EXISTS idx_search_history_created_at;
DROP INDEX IF EXISTS idx_search_history_user_id;
DROP TABLE IF EXISTS search_history;

DROP INDEX IF EXISTS idx_account_switches_created_at;
DROP INDEX IF EXISTS idx_account_switches_user_id;
DROP TABLE IF EXISTS account_switches;

DROP INDEX IF EXISTS idx_conversations_last_message_at;
DROP INDEX IF EXISTS idx_conversations_participant2;
DROP INDEX IF EXISTS idx_conversations_participant1;
DROP TABLE IF EXISTS conversations;

DROP INDEX IF EXISTS idx_messages_created_at;
DROP INDEX IF EXISTS idx_messages_is_read;
DROP INDEX IF EXISTS idx_messages_recipient_id;
DROP INDEX IF EXISTS idx_messages_sender_id;
DROP TABLE IF EXISTS messages;

DROP INDEX IF EXISTS idx_notifications_created_at;
DROP INDEX IF EXISTS idx_notifications_type;
DROP INDEX IF EXISTS idx_notifications_is_read;
DROP INDEX IF EXISTS idx_notifications_user_id;
DROP TABLE IF EXISTS notifications;


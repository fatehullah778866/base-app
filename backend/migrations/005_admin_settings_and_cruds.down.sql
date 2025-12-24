-- Rollback migration
DROP INDEX IF EXISTS idx_admin_permissions_name;
DROP INDEX IF EXISTS idx_admin_permissions_admin_id;
DROP TABLE IF EXISTS admin_permissions;

DROP INDEX IF EXISTS idx_user_management_actions_action_type;
DROP INDEX IF EXISTS idx_user_management_actions_user_id;
DROP INDEX IF EXISTS idx_user_management_actions_admin_id;
DROP TABLE IF EXISTS user_management_actions;

DROP INDEX IF EXISTS idx_admin_activity_logs_created_at;
DROP INDEX IF EXISTS idx_admin_activity_logs_entity_type;
DROP INDEX IF EXISTS idx_admin_activity_logs_action;
DROP INDEX IF EXISTS idx_admin_activity_logs_admin_id;
DROP TABLE IF EXISTS admin_activity_logs;

DROP INDEX IF EXISTS idx_custom_crud_data_deleted_at;
DROP INDEX IF EXISTS idx_custom_crud_data_created_by;
DROP INDEX IF EXISTS idx_custom_crud_data_entity_id;
DROP TABLE IF EXISTS custom_crud_data;

DROP INDEX IF EXISTS idx_custom_crud_entities_active;
DROP INDEX IF EXISTS idx_custom_crud_entities_created_by;
DROP TABLE IF EXISTS custom_crud_entities;

DROP TABLE IF EXISTS admin_settings;


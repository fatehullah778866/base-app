-- Rollback CRUD Templates Migration

DROP INDEX IF EXISTS idx_crud_templates_created_by;
DROP INDEX IF EXISTS idx_crud_templates_active;
DROP INDEX IF EXISTS idx_crud_templates_category;
DROP INDEX IF EXISTS idx_crud_templates_name;

DROP TABLE IF EXISTS crud_templates;



-- CRUD Templates Migration
-- This allows admins to create and manage CRUD templates in the database

CREATE TABLE IF NOT EXISTS crud_templates (
    id TEXT PRIMARY KEY,
    name TEXT NOT NULL UNIQUE, -- e.g., "portfolio", "visa", "products"
    display_name TEXT NOT NULL,
    description TEXT,
    schema TEXT NOT NULL, -- JSON schema defining fields
    icon TEXT, -- Icon identifier
    category TEXT, -- e.g., "business", "travel", "ecommerce"
    created_by TEXT NOT NULL, -- Admin who created the template
    is_active INTEGER DEFAULT 1,
    is_system INTEGER DEFAULT 0, -- System templates cannot be deleted
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_crud_templates_name ON crud_templates(name);
CREATE INDEX IF NOT EXISTS idx_crud_templates_category ON crud_templates(category);
CREATE INDEX IF NOT EXISTS idx_crud_templates_active ON crud_templates(is_active);
CREATE INDEX IF NOT EXISTS idx_crud_templates_created_by ON crud_templates(created_by);



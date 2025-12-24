-- Admin Settings and Flexible CRUD System Migration

-- Admin Settings table
CREATE TABLE IF NOT EXISTS admin_settings (
    admin_id TEXT PRIMARY KEY,
    dashboard_layout TEXT, -- JSON for dashboard customization
    default_permissions TEXT, -- JSON array of default permissions
    notification_preferences TEXT, -- JSON for notification settings
    theme_preferences TEXT, -- JSON for admin theme
    admin_verification_code TEXT DEFAULT 'Kompasstech2025@', -- Verification code for admin creation
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(admin_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Custom CRUD Entities table (for admin-created CRUDs)
CREATE TABLE IF NOT EXISTS custom_crud_entities (
    id TEXT PRIMARY KEY,
    created_by TEXT NOT NULL,
    entity_name TEXT NOT NULL UNIQUE, -- e.g., "products", "orders", "inventory"
    display_name TEXT NOT NULL,
    description TEXT,
    schema TEXT NOT NULL, -- JSON schema defining fields
    is_active INTEGER DEFAULT 1,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_custom_crud_entities_created_by ON custom_crud_entities(created_by);
CREATE INDEX IF NOT EXISTS idx_custom_crud_entities_active ON custom_crud_entities(is_active);

-- Custom CRUD Data table (stores actual data for custom entities)
CREATE TABLE IF NOT EXISTS custom_crud_data (
    id TEXT PRIMARY KEY,
    entity_id TEXT NOT NULL,
    data TEXT NOT NULL, -- JSON data
    created_by TEXT NOT NULL,
    updated_by TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY(entity_id) REFERENCES custom_crud_entities(id) ON DELETE CASCADE,
    FOREIGN KEY(created_by) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(updated_by) REFERENCES users(id) ON DELETE SET NULL
);

CREATE INDEX IF NOT EXISTS idx_custom_crud_data_entity_id ON custom_crud_data(entity_id);
CREATE INDEX IF NOT EXISTS idx_custom_crud_data_created_by ON custom_crud_data(created_by);
CREATE INDEX IF NOT EXISTS idx_custom_crud_data_deleted_at ON custom_crud_data(deleted_at);

-- Admin Activity Logs (enhanced)
CREATE TABLE IF NOT EXISTS admin_activity_logs (
    id TEXT PRIMARY KEY,
    admin_id TEXT NOT NULL,
    action TEXT NOT NULL, -- create, update, delete, view, etc.
    entity_type TEXT NOT NULL, -- user, custom_crud, settings, etc.
    entity_id TEXT,
    details TEXT, -- JSON for action details
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(admin_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_admin_activity_logs_admin_id ON admin_activity_logs(admin_id);
CREATE INDEX IF NOT EXISTS idx_admin_activity_logs_action ON admin_activity_logs(action);
CREATE INDEX IF NOT EXISTS idx_admin_activity_logs_entity_type ON admin_activity_logs(entity_type);
CREATE INDEX IF NOT EXISTS idx_admin_activity_logs_created_at ON admin_activity_logs(created_at);

-- User Management Actions (for tracking user CRUD operations)
CREATE TABLE IF NOT EXISTS user_management_actions (
    id TEXT PRIMARY KEY,
    admin_id TEXT NOT NULL,
    user_id TEXT NOT NULL,
    action_type TEXT NOT NULL, -- create, update, delete, activate, deactivate, suspend
    changes TEXT, -- JSON of what changed
    reason TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(admin_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_management_actions_admin_id ON user_management_actions(admin_id);
CREATE INDEX IF NOT EXISTS idx_user_management_actions_user_id ON user_management_actions(user_id);
CREATE INDEX IF NOT EXISTS idx_user_management_actions_action_type ON user_management_actions(action_type);

-- Admin Permissions table
CREATE TABLE IF NOT EXISTS admin_permissions (
    id TEXT PRIMARY KEY,
    admin_id TEXT NOT NULL,
    permission_name TEXT NOT NULL, -- manage_users, manage_cruds, manage_settings, etc.
    granted_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    granted_by TEXT,
    FOREIGN KEY(admin_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(granted_by) REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(admin_id, permission_name)
);

CREATE INDEX IF NOT EXISTS idx_admin_permissions_admin_id ON admin_permissions(admin_id);
CREATE INDEX IF NOT EXISTS idx_admin_permissions_name ON admin_permissions(permission_name);


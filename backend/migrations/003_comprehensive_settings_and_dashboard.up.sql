-- Comprehensive Settings and Dashboard Migration

-- Password Reset Tokens table
CREATE TABLE IF NOT EXISTS password_reset_tokens (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    expires_at DATETIME NOT NULL,
    used_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_user_id ON password_reset_tokens(user_id);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_token ON password_reset_tokens(token);
CREATE INDEX IF NOT EXISTS idx_password_reset_tokens_expires ON password_reset_tokens(expires_at);

-- Enhanced User Settings table (comprehensive)
CREATE TABLE IF NOT EXISTS user_settings_comprehensive (
    user_id TEXT PRIMARY KEY,
    -- Profile Settings
    username TEXT UNIQUE,
    display_name TEXT,
    bio TEXT,
    date_of_birth TEXT,
    
    -- Security Settings
    two_factor_enabled INTEGER DEFAULT 0,
    two_factor_secret TEXT,
    two_factor_backup_codes TEXT,
    security_questions TEXT, -- JSON array
    password_last_changed DATETIME,
    
    -- Privacy Settings
    profile_visibility TEXT DEFAULT 'public', -- public, private, friends
    email_visibility TEXT DEFAULT 'private', -- public, private, friends
    phone_visibility TEXT DEFAULT 'private',
    allow_messaging TEXT DEFAULT 'everyone', -- everyone, friends, none
    search_visibility INTEGER DEFAULT 1, -- 1 = visible, 0 = hidden
    data_sharing_enabled INTEGER DEFAULT 0,
    
    -- Notification Settings
    email_notifications INTEGER DEFAULT 1,
    sms_notifications INTEGER DEFAULT 0,
    push_notifications INTEGER DEFAULT 1,
    notification_messages INTEGER DEFAULT 1,
    notification_alerts INTEGER DEFAULT 1,
    notification_promotions INTEGER DEFAULT 0,
    notification_security INTEGER DEFAULT 1,
    
    -- Account Preferences
    language TEXT DEFAULT 'en',
    timezone TEXT DEFAULT 'UTC',
    theme TEXT DEFAULT 'light', -- light, dark, auto
    font_size TEXT DEFAULT 'medium', -- small, medium, large
    high_contrast INTEGER DEFAULT 0,
    reduced_motion INTEGER DEFAULT 0,
    screen_reader INTEGER DEFAULT 0,
    
    -- Connected Accounts (JSON array)
    connected_accounts TEXT, -- JSON: [{"provider": "google", "email": "...", "connected_at": "..."}]
    
    -- Data & Account Control
    account_deletion_requested INTEGER DEFAULT 0,
    account_deletion_scheduled_at DATETIME,
    account_deactivated INTEGER DEFAULT 0,
    account_deactivated_at DATETIME,
    
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_settings_username ON user_settings_comprehensive(username);

-- Dashboard Items table (CRUD with title and description)
CREATE TABLE IF NOT EXISTS dashboard_items (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    title TEXT NOT NULL,
    description TEXT,
    category TEXT, -- Optional category for grouping
    status TEXT DEFAULT 'active', -- active, archived, deleted
    priority INTEGER DEFAULT 0, -- For ordering
    metadata TEXT, -- JSON for extensibility
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    deleted_at DATETIME,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_dashboard_items_user_id ON dashboard_items(user_id);
CREATE INDEX IF NOT EXISTS idx_dashboard_items_status ON dashboard_items(status);
CREATE INDEX IF NOT EXISTS idx_dashboard_items_category ON dashboard_items(category);
CREATE INDEX IF NOT EXISTS idx_dashboard_items_created_at ON dashboard_items(created_at);

-- User Sessions (enhanced for security settings)
-- Already exists, but ensure it has all needed fields
-- Sessions table already has: device_type, device_name, device_id, os, browser, ip_address, etc.


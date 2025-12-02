-- SQLite-compatible schema

PRAGMA foreign_keys = ON;

-- Users table
CREATE TABLE IF NOT EXISTS users (
    id TEXT PRIMARY KEY,
    email TEXT UNIQUE NOT NULL,
    email_verified INTEGER DEFAULT 0,
    email_verification_token TEXT,
    email_verification_expires_at DATETIME,
    password_hash TEXT NOT NULL,
    password_changed_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    name TEXT,
    first_name TEXT,
    last_name TEXT,
    photo_url TEXT,
    bio TEXT,
    phone TEXT,
    phone_verified INTEGER DEFAULT 0,
    phone_verification_code TEXT,
    phone_verification_expires_at DATETIME,
    date_of_birth TEXT,
    gender TEXT,
    locale TEXT DEFAULT 'en',
    timezone TEXT,
    status TEXT DEFAULT 'active',
    status_changed_at DATETIME,
    status_reason TEXT,
    signup_source TEXT,
    signup_platform TEXT,
    referrer_url TEXT,
    signup_campaign TEXT,
    last_login_at DATETIME,
    last_login_ip TEXT,
    last_active_at DATETIME,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP
);

CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);
CREATE INDEX IF NOT EXISTS idx_users_status ON users(status);
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at);

-- Sessions table
CREATE TABLE IF NOT EXISTS sessions (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    token TEXT UNIQUE NOT NULL,
    refresh_token TEXT UNIQUE,
    refresh_token_expires_at DATETIME,
    device_type TEXT,
    device_name TEXT,
    device_id TEXT,
    os TEXT,
    browser TEXT,
    ip_address TEXT,
    user_agent TEXT,
    location_country TEXT,
    location_city TEXT,
    is_active INTEGER DEFAULT 1,
    revoked_at DATETIME,
    revoked_reason TEXT,
    expires_at DATETIME NOT NULL,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    last_used_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_sessions_user_id ON sessions(user_id);
CREATE INDEX IF NOT EXISTS idx_sessions_token ON sessions(token);
CREATE INDEX IF NOT EXISTS idx_sessions_refresh_token ON sessions(refresh_token);

-- User devices table
CREATE TABLE IF NOT EXISTS user_devices (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    device_id TEXT UNIQUE NOT NULL,
    device_name TEXT,
    device_type TEXT,
    os TEXT,
    browser TEXT,
    ip_address TEXT,
    location_country TEXT,
    location_city TEXT,
    is_trusted INTEGER DEFAULT 0,
    trusted_at DATETIME,
    last_used_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_user_devices_user_id ON user_devices(user_id);
CREATE INDEX IF NOT EXISTS idx_user_devices_device_id ON user_devices(device_id);

-- User settings table
CREATE TABLE IF NOT EXISTS user_settings (
    user_id TEXT PRIMARY KEY,
    email_notifications_enabled INTEGER DEFAULT 1,
    push_notifications_enabled INTEGER DEFAULT 1,
    sms_notifications_enabled INTEGER DEFAULT 0,
    notification_email_marketing INTEGER DEFAULT 0,
    notification_email_product INTEGER DEFAULT 1,
    notification_email_security INTEGER DEFAULT 1,
    notification_push_product INTEGER DEFAULT 1,
    notification_push_security INTEGER DEFAULT 1,
    privacy_level TEXT DEFAULT 'public',
    profile_visibility TEXT DEFAULT 'public',
    show_email INTEGER DEFAULT 0,
    show_phone INTEGER DEFAULT 0,
    theme TEXT DEFAULT 'light',
    language TEXT DEFAULT 'en',
    currency TEXT DEFAULT 'USD',
    date_format TEXT DEFAULT 'MM/DD/YYYY',
    time_format TEXT DEFAULT '12h',
    preferred_contact_method TEXT DEFAULT 'email',
    kompassui_theme TEXT DEFAULT 'auto',
    kompassui_contrast TEXT DEFAULT 'standard',
    kompassui_text_direction TEXT DEFAULT 'auto',
    kompassui_brand TEXT,
    theme_synced_at DATETIME,
    theme_sync_enabled INTEGER DEFAULT 1,
    accessibility_high_contrast INTEGER DEFAULT 0,
    accessibility_reduced_motion INTEGER DEFAULT 0,
    accessibility_screen_reader INTEGER DEFAULT 0,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

-- Product theme overrides table
CREATE TABLE IF NOT EXISTS product_theme_preferences (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    product_name TEXT NOT NULL,
    theme TEXT,
    contrast TEXT,
    text_direction TEXT,
    brand TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(user_id, product_name),
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_product_theme_preferences_user_id ON product_theme_preferences(user_id);
CREATE INDEX IF NOT EXISTS idx_product_theme_preferences_product ON product_theme_preferences(product_name);

-- Webhook subscriptions table
CREATE TABLE IF NOT EXISTS webhook_subscriptions (
    id TEXT PRIMARY KEY,
    user_id TEXT,
    subscription_name TEXT NOT NULL,
    webhook_url TEXT NOT NULL,
    webhook_secret TEXT NOT NULL,
    event_types TEXT,
    is_active INTEGER DEFAULT 1,
    is_verified INTEGER DEFAULT 0,
    rate_limit_per_minute INTEGER DEFAULT 60,
    max_retries INTEGER DEFAULT 3,
    retry_backoff_multiplier REAL DEFAULT 2.0,
    description TEXT,
    metadata TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_user_id ON webhook_subscriptions(user_id);
CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_active ON webhook_subscriptions(is_active);
CREATE INDEX IF NOT EXISTS idx_webhook_subscriptions_url ON webhook_subscriptions(webhook_url);

-- Webhook events table
CREATE TABLE IF NOT EXISTS webhook_events (
    id TEXT PRIMARY KEY,
    event_type TEXT NOT NULL,
    event_version TEXT NOT NULL,
    event_source TEXT NOT NULL,
    user_id TEXT NOT NULL,
    payload TEXT NOT NULL,
    payload_hash TEXT NOT NULL,
    webhook_url TEXT NOT NULL,
    webhook_secret TEXT,
    status TEXT DEFAULT 'pending',
    delivery_attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    scheduled_at DATETIME NOT NULL,
    processed_at DATETIME,
    delivered_at DATETIME,
    next_retry_at DATETIME,
    last_response_status INTEGER,
    last_response_body TEXT,
    last_error_message TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_webhook_events_status ON webhook_events(status);
CREATE INDEX IF NOT EXISTS idx_webhook_events_user_id ON webhook_events(user_id);
CREATE INDEX IF NOT EXISTS idx_webhook_events_scheduled ON webhook_events(scheduled_at);


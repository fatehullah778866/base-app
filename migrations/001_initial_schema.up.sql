-- Enable UUID extension
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";

-- Users table
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    email VARCHAR(255) UNIQUE NOT NULL,
    email_verified BOOLEAN DEFAULT FALSE,
    email_verification_token VARCHAR(255),
    email_verification_expires_at TIMESTAMP,
    
    password_hash VARCHAR(255) NOT NULL,
    password_changed_at TIMESTAMP DEFAULT NOW(),
    
    name VARCHAR(255),
    first_name VARCHAR(100),
    last_name VARCHAR(100),
    photo_url TEXT,
    bio TEXT,
    
    phone VARCHAR(20),
    phone_verified BOOLEAN DEFAULT FALSE,
    phone_verification_code VARCHAR(10),
    phone_verification_expires_at TIMESTAMP,
    
    date_of_birth DATE,
    gender VARCHAR(50),
    locale VARCHAR(10) DEFAULT 'en',
    timezone VARCHAR(50),
    
    status VARCHAR(50) DEFAULT 'active',
    status_changed_at TIMESTAMP,
    status_reason TEXT,
    
    signup_source VARCHAR(100),
    signup_platform VARCHAR(100),
    referrer_url TEXT,
    signup_campaign VARCHAR(255),
    
    last_login_at TIMESTAMP,
    last_login_ip INET,
    last_active_at TIMESTAMP,
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT email_format CHECK (email ~* '^[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}$'),
    CONSTRAINT status_valid CHECK (status IN ('active', 'suspended', 'deleted', 'pending'))
);

CREATE INDEX idx_users_email ON users(email);
CREATE INDEX idx_users_status ON users(status);
CREATE INDEX idx_users_created_at ON users(created_at);
CREATE INDEX idx_users_email_verification_token ON users(email_verification_token) WHERE email_verification_token IS NOT NULL;

-- Sessions table
CREATE TABLE sessions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    refresh_token VARCHAR(255) UNIQUE,
    refresh_token_expires_at TIMESTAMP,
    
    device_type VARCHAR(50),
    device_name VARCHAR(255),
    device_id VARCHAR(255),
    os VARCHAR(100),
    browser VARCHAR(100),
    ip_address INET,
    user_agent TEXT,
    location_country VARCHAR(2),
    location_city VARCHAR(100),
    
    is_active BOOLEAN DEFAULT TRUE,
    revoked_at TIMESTAMP,
    revoked_reason TEXT,
    
    expires_at TIMESTAMP NOT NULL,
    created_at TIMESTAMP DEFAULT NOW(),
    last_used_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_sessions_user_id ON sessions(user_id);
CREATE INDEX idx_sessions_token ON sessions(token);
CREATE INDEX idx_sessions_refresh_token ON sessions(refresh_token);
CREATE INDEX idx_sessions_user_active ON sessions(user_id, is_active) WHERE is_active = TRUE;
CREATE INDEX idx_sessions_expires_at ON sessions(expires_at) WHERE is_active = TRUE;

-- User devices table
CREATE TABLE user_devices (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    device_id VARCHAR(255) UNIQUE NOT NULL,
    device_name VARCHAR(255),
    device_type VARCHAR(50),
    os VARCHAR(100),
    browser VARCHAR(100),
    ip_address INET,
    location_country VARCHAR(2),
    location_city VARCHAR(100),
    is_trusted BOOLEAN DEFAULT FALSE,
    trusted_at TIMESTAMP,
    last_used_at TIMESTAMP DEFAULT NOW(),
    created_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_user_devices_user_id ON user_devices(user_id);
CREATE INDEX idx_user_devices_device_id ON user_devices(device_id);
CREATE INDEX idx_user_devices_user_trusted ON user_devices(user_id, is_trusted);

-- User settings table
CREATE TABLE user_settings (
    user_id UUID PRIMARY KEY REFERENCES users(id) ON DELETE CASCADE,
    
    -- Notifications
    email_notifications_enabled BOOLEAN DEFAULT TRUE,
    push_notifications_enabled BOOLEAN DEFAULT TRUE,
    sms_notifications_enabled BOOLEAN DEFAULT FALSE,
    
    notification_email_marketing BOOLEAN DEFAULT FALSE,
    notification_email_product BOOLEAN DEFAULT TRUE,
    notification_email_security BOOLEAN DEFAULT TRUE,
    notification_push_product BOOLEAN DEFAULT TRUE,
    notification_push_security BOOLEAN DEFAULT TRUE,
    
    -- Privacy
    privacy_level VARCHAR(50) DEFAULT 'public',
    profile_visibility VARCHAR(50) DEFAULT 'public',
    show_email BOOLEAN DEFAULT FALSE,
    show_phone BOOLEAN DEFAULT FALSE,
    
    -- Preferences
    theme VARCHAR(50) DEFAULT 'light',
    language VARCHAR(10) DEFAULT 'en',
    currency VARCHAR(10) DEFAULT 'USD',
    date_format VARCHAR(50) DEFAULT 'MM/DD/YYYY',
    time_format VARCHAR(50) DEFAULT '12h',
    
    preferred_contact_method VARCHAR(50) DEFAULT 'email',
    
    -- KompassUI Theme Preferences
    kompassui_theme VARCHAR(50) DEFAULT 'auto',
    kompassui_contrast VARCHAR(50) DEFAULT 'standard',
    kompassui_text_direction VARCHAR(10) DEFAULT 'auto',
    kompassui_brand VARCHAR(100),
    theme_synced_at TIMESTAMP,
    theme_sync_enabled BOOLEAN DEFAULT TRUE,
    
    -- Accessibility
    accessibility_high_contrast BOOLEAN DEFAULT FALSE,
    accessibility_reduced_motion BOOLEAN DEFAULT FALSE,
    accessibility_screen_reader BOOLEAN DEFAULT FALSE,
    
    updated_at TIMESTAMP DEFAULT NOW()
);

-- Product theme overrides table
CREATE TABLE product_theme_preferences (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    product_name VARCHAR(100) NOT NULL,
    
    theme VARCHAR(50),
    contrast VARCHAR(50),
    text_direction VARCHAR(10),
    brand VARCHAR(100),
    
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    UNIQUE(user_id, product_name)
);

CREATE INDEX idx_product_theme_preferences_user_id ON product_theme_preferences(user_id);
CREATE INDEX idx_product_theme_preferences_product ON product_theme_preferences(product_name);

-- Webhook subscriptions table
CREATE TABLE webhook_subscriptions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID REFERENCES users(id) ON DELETE CASCADE,
    subscription_name VARCHAR(255) NOT NULL,
    webhook_url TEXT NOT NULL,
    webhook_secret VARCHAR(255) NOT NULL,
    event_types TEXT[] NOT NULL,
    is_active BOOLEAN DEFAULT TRUE,
    is_verified BOOLEAN DEFAULT FALSE,
    rate_limit_per_minute INTEGER DEFAULT 60,
    max_retries INTEGER DEFAULT 3,
    retry_backoff_multiplier DECIMAL(3,2) DEFAULT 2.0,
    description TEXT,
    metadata JSONB,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW()
);

CREATE INDEX idx_webhook_subscriptions_user_id ON webhook_subscriptions(user_id);
CREATE INDEX idx_webhook_subscriptions_active ON webhook_subscriptions(is_active) WHERE is_active = TRUE;
CREATE INDEX idx_webhook_subscriptions_url ON webhook_subscriptions(webhook_url);

-- Webhook events table (outbox pattern)
CREATE TABLE webhook_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    event_type VARCHAR(100) NOT NULL,
    event_version VARCHAR(50) NOT NULL,
    event_source VARCHAR(100) NOT NULL,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    payload JSONB NOT NULL,
    payload_hash VARCHAR(64) NOT NULL,
    webhook_url TEXT NOT NULL,
    webhook_secret VARCHAR(255),
    status VARCHAR(50) DEFAULT 'pending',
    delivery_attempts INTEGER DEFAULT 0,
    max_attempts INTEGER DEFAULT 3,
    scheduled_at TIMESTAMP NOT NULL,
    processed_at TIMESTAMP,
    delivered_at TIMESTAMP,
    next_retry_at TIMESTAMP,
    last_response_status INTEGER,
    last_response_body TEXT,
    last_error_message TEXT,
    created_at TIMESTAMP DEFAULT NOW(),
    updated_at TIMESTAMP DEFAULT NOW(),
    
    CONSTRAINT status_valid CHECK (status IN ('pending', 'processing', 'retrying', 'delivered', 'failed'))
);

CREATE INDEX idx_webhook_events_status ON webhook_events(status);
CREATE INDEX idx_webhook_events_user_id ON webhook_events(user_id);
CREATE INDEX idx_webhook_events_scheduled ON webhook_events(scheduled_at);
CREATE INDEX idx_webhook_events_pending ON webhook_events(status, scheduled_at) WHERE status IN ('pending', 'retrying');
CREATE INDEX idx_webhook_events_next_retry ON webhook_events(next_retry_at) WHERE next_retry_at IS NOT NULL;


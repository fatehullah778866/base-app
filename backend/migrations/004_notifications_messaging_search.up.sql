-- Notifications, Messaging, Account Switching, and Search Migration

-- Notifications table
CREATE TABLE IF NOT EXISTS notifications (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    type TEXT NOT NULL, -- message, alert, promotion, security, system
    title TEXT NOT NULL,
    message TEXT NOT NULL,
    link TEXT,
    is_read INTEGER DEFAULT 0,
    read_at DATETIME,
    metadata TEXT, -- JSON for additional data
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_notifications_user_id ON notifications(user_id);
CREATE INDEX IF NOT EXISTS idx_notifications_is_read ON notifications(is_read);
CREATE INDEX IF NOT EXISTS idx_notifications_type ON notifications(type);
CREATE INDEX IF NOT EXISTS idx_notifications_created_at ON notifications(created_at);

-- Messages table (for user-to-user messaging)
CREATE TABLE IF NOT EXISTS messages (
    id TEXT PRIMARY KEY,
    sender_id TEXT NOT NULL,
    recipient_id TEXT NOT NULL,
    subject TEXT,
    content TEXT NOT NULL,
    is_read INTEGER DEFAULT 0,
    read_at DATETIME,
    is_archived INTEGER DEFAULT 0,
    archived_at DATETIME,
    metadata TEXT, -- JSON for attachments, etc.
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(sender_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(recipient_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_messages_sender_id ON messages(sender_id);
CREATE INDEX IF NOT EXISTS idx_messages_recipient_id ON messages(recipient_id);
CREATE INDEX IF NOT EXISTS idx_messages_is_read ON messages(is_read);
CREATE INDEX IF NOT EXISTS idx_messages_created_at ON messages(created_at);

-- Conversations table (for grouping messages)
CREATE TABLE IF NOT EXISTS conversations (
    id TEXT PRIMARY KEY,
    participant1_id TEXT NOT NULL,
    participant2_id TEXT NOT NULL,
    last_message_id TEXT,
    last_message_at DATETIME,
    participant1_unread_count INTEGER DEFAULT 0,
    participant2_unread_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    updated_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(participant1_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(participant2_id) REFERENCES users(id) ON DELETE CASCADE,
    FOREIGN KEY(last_message_id) REFERENCES messages(id) ON DELETE SET NULL,
    UNIQUE(participant1_id, participant2_id)
);

CREATE INDEX IF NOT EXISTS idx_conversations_participant1 ON conversations(participant1_id);
CREATE INDEX IF NOT EXISTS idx_conversations_participant2 ON conversations(participant2_id);
CREATE INDEX IF NOT EXISTS idx_conversations_last_message_at ON conversations(last_message_at);

-- Account switching (for users with multiple accounts/roles)
CREATE TABLE IF NOT EXISTS account_switches (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    switched_to_user_id TEXT, -- If switching to another user account
    switched_to_role TEXT, -- If switching role context
    switched_from_role TEXT,
    reason TEXT,
    ip_address TEXT,
    user_agent TEXT,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_account_switches_user_id ON account_switches(user_id);
CREATE INDEX IF NOT EXISTS idx_account_switches_created_at ON account_switches(created_at);

-- Search history (for tracking user searches)
CREATE TABLE IF NOT EXISTS search_history (
    id TEXT PRIMARY KEY,
    user_id TEXT NOT NULL,
    query TEXT NOT NULL,
    search_type TEXT, -- users, dashboard_items, messages, all
    results_count INTEGER DEFAULT 0,
    created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
    FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
);

CREATE INDEX IF NOT EXISTS idx_search_history_user_id ON search_history(user_id);
CREATE INDEX IF NOT EXISTS idx_search_history_created_at ON search_history(created_at);

-- Full-text search support (for SQLite FTS5)
CREATE VIRTUAL TABLE IF NOT EXISTS dashboard_items_fts USING fts5(
    id UNINDEXED,
    user_id UNINDEXED,
    title,
    description,
    category,
    content='dashboard_items',
    content_rowid='rowid'
);

CREATE VIRTUAL TABLE IF NOT EXISTS messages_fts USING fts5(
    id UNINDEXED,
    sender_id UNINDEXED,
    recipient_id UNINDEXED,
    subject,
    content,
    content='messages',
    content_rowid='rowid'
);

-- Triggers to keep FTS tables in sync
CREATE TRIGGER IF NOT EXISTS dashboard_items_fts_insert AFTER INSERT ON dashboard_items BEGIN
    INSERT INTO dashboard_items_fts(rowid, id, user_id, title, description, category)
    VALUES (new.rowid, new.id, new.user_id, new.title, new.description, new.category);
END;

CREATE TRIGGER IF NOT EXISTS dashboard_items_fts_delete AFTER DELETE ON dashboard_items BEGIN
    INSERT INTO dashboard_items_fts(dashboard_items_fts, rowid, id, user_id, title, description, category)
    VALUES ('delete', old.rowid, old.id, old.user_id, old.title, old.description, old.category);
END;

CREATE TRIGGER IF NOT EXISTS dashboard_items_fts_update AFTER UPDATE ON dashboard_items BEGIN
    INSERT INTO dashboard_items_fts(dashboard_items_fts, rowid, id, user_id, title, description, category)
    VALUES ('delete', old.rowid, old.id, old.user_id, old.title, old.description, old.category);
    INSERT INTO dashboard_items_fts(rowid, id, user_id, title, description, category)
    VALUES (new.rowid, new.id, new.user_id, new.title, new.description, new.category);
END;

CREATE TRIGGER IF NOT EXISTS messages_fts_insert AFTER INSERT ON messages BEGIN
    INSERT INTO messages_fts(rowid, id, sender_id, recipient_id, subject, content)
    VALUES (new.rowid, new.id, new.sender_id, new.recipient_id, new.subject, new.content);
END;

CREATE TRIGGER IF NOT EXISTS messages_fts_delete AFTER DELETE ON messages BEGIN
    INSERT INTO messages_fts(messages_fts, rowid, id, sender_id, recipient_id, subject, content)
    VALUES ('delete', old.rowid, old.id, old.sender_id, old.recipient_id, old.subject, old.content);
END;

CREATE TRIGGER IF NOT EXISTS messages_fts_update AFTER UPDATE ON messages BEGIN
    INSERT INTO messages_fts(messages_fts, rowid, id, sender_id, recipient_id, subject, content)
    VALUES ('delete', old.rowid, old.id, old.sender_id, old.recipient_id, old.subject, old.content);
    INSERT INTO messages_fts(rowid, id, sender_id, recipient_id, subject, content)
    VALUES (new.rowid, new.id, new.sender_id, new.recipient_id, new.subject, new.content);
END;


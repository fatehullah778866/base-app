// Admin Dashboard functionality
let users = [];
let currentUserId = null;
let currentSettings = null;
let conversations = [];
let activeConversationId = null;
let activeRecipientId = null;

document.addEventListener('DOMContentLoaded', async () => {
    if (!requireAuth()) return;

    const user = JSON.parse(localStorage.getItem('user') || '{}');
    currentUserId = user.id;
    if (user.role !== 'admin') {
        redirect('dashboard.html');
        return;
    }

    setupNavbar();
    await loadSettingsForUI();
    await loadStats();
    await loadUsers();
    await loadSettings();
    setupEventListeners();
    await loadConversations();
    await refreshBadges();
});

function setupNavbar() {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const adminNameEl = document.getElementById('admin-name');
    const adminIcon = document.getElementById('admin-icon');
    const dropdown = document.getElementById('admin-dropdown');

    if (adminNameEl) adminNameEl.textContent = user.name || user.email || 'Admin';
    applyAvatar(adminIcon, user.name || user.email || 'Admin', user.photo_url);

    adminIcon?.addEventListener('click', () => {
        dropdown?.classList.toggle('show');
    });

    document.addEventListener('click', (e) => {
        if (!dropdown?.contains(e.target) && !adminIcon?.contains(e.target)) {
            dropdown?.classList.remove('show');
        }
        const notificationsMenu = document.getElementById('notifications-menu');
        const notificationsButton = document.getElementById('notifications-button');
        if (!notificationsMenu?.contains(e.target) && !notificationsButton?.contains(e.target)) {
            notificationsMenu?.classList.remove('show');
        }
    });

    document.getElementById('admin-logout-link')?.addEventListener('click', async (e) => {
        e.preventDefault();
        localStorage.clear();
        redirect('index.html');
    });

    document.getElementById('messages-button')?.addEventListener('click', () => {
        document.getElementById('messages-section')?.scrollIntoView({ behavior: 'smooth' });
    });

    document.getElementById('notifications-button')?.addEventListener('click', (e) => {
        e.stopPropagation();
        document.getElementById('notifications-menu')?.classList.toggle('show');
        loadNotificationsList();
    });

    document.getElementById('mark-notifications-read')?.addEventListener('click', async (e) => {
        e.preventDefault();
        await notificationAPI.markAllAsRead();
        await refreshBadges();
        loadNotificationsList();
    });
}

function applyAvatar(imageEl, name, photoUrl) {
    if (!imageEl) return;
    const resolved = resolvePhotoUrl(photoUrl);
    if (resolved) {
        imageEl.src = resolved;
        return;
    }
    imageEl.src = createInitialAvatar(name);
}

function createInitialAvatar(name) {
    const safeName = (name || 'A').trim();
    const initial = safeName ? safeName[0].toUpperCase() : 'A';
    const svg = `<svg xmlns="http://www.w3.org/2000/svg" width="80" height="80">
        <rect width="100%" height="100%" fill="#3498db"/>
        <text x="50%" y="55%" font-size="36" text-anchor="middle" fill="#ffffff" font-family="Arial, sans-serif">${initial}</text>
    </svg>`;
    return `data:image/svg+xml;charset=UTF-8,${encodeURIComponent(svg)}`;
}

function getApiOrigin() {
    return API_BASE_URL.replace(/\/v1\/?$/, '');
}

function resolvePhotoUrl(photoUrl) {
    if (!photoUrl) return '';
    if (photoUrl.startsWith('http://') || photoUrl.startsWith('https://')) {
        return photoUrl;
    }
    if (photoUrl.startsWith('/uploads/')) {
        return `${getApiOrigin()}${photoUrl}`;
    }
    return photoUrl;
}

function setupEventListeners() {
    document.getElementById('add-user-btn')?.addEventListener('click', () => {
        openUserModal();
    });

    document.getElementById('user-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveUser();
    });

    document.getElementById('cancel-user-btn')?.addEventListener('click', closeUserModal);
    document.querySelector('.close')?.addEventListener('click', closeUserModal);

    document.getElementById('user-search')?.addEventListener('input', (e) => {
        filterUsers(e.target.value);
    });

    document.getElementById('save-verification-code-btn')?.addEventListener('click', async () => {
        await saveVerificationCode();
    });

    window.addEventListener('click', (e) => {
        if (e.target.classList.contains('modal')) {
            closeUserModal();
        }
    });

    document.getElementById('message-search-btn')?.addEventListener('click', async () => {
        await searchUsersForMessaging();
    });

    document.getElementById('message-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        await sendMessage();
    });
}

async function loadStats() {
    try {
        const response = await adminAPI.getUsers();
        if (response.success && response.data) {
            const allUsers = response.data;
            const activeUsers = allUsers.filter(u => u.status === 'active');
            const admins = allUsers.filter(u => u.role === 'admin');

            document.getElementById('total-users').textContent = allUsers.length;
            document.getElementById('active-users').textContent = activeUsers.length;
            document.getElementById('total-admins').textContent = admins.length;
        }
    } catch (error) {
        console.error('Failed to load stats:', error);
    }
}

async function loadUsers() {
    try {
        const response = await adminAPI.getUsers();
        if (response.success) {
            users = response.data || [];
            renderUsers(users);
        }
    } catch (error) {
        showMessage('Failed to load users', 'error');
    }
}

function renderUsers(usersToRender) {
    const container = document.getElementById('users-list');
    if (!container) return;

    if (usersToRender.length === 0) {
        container.innerHTML = '<div class="loading">No users found.</div>';
        return;
    }

    container.innerHTML = `
        <table>
            <thead>
                <tr>
                    <th>Name</th>
                    <th>Email</th>
                    <th>Role</th>
                    <th>Status</th>
                    <th>Actions</th>
                </tr>
            </thead>
            <tbody>
                ${usersToRender.map(user => `
                    <tr>
                        <td>${escapeHtml(user.name || '')}</td>
                        <td>${escapeHtml(user.email || '')}</td>
                        <td>${escapeHtml(user.role || 'user')}</td>
                        <td>${escapeHtml(user.status || 'active')}</td>
                        <td>
                            <button class="btn btn-primary" onclick="editUser('${user.id}')">Edit</button>
                            <button class="btn btn-secondary" onclick="deleteUser('${user.id}')">Delete</button>
                        </td>
                    </tr>
                `).join('')}
            </tbody>
        </table>
    `;
}

function filterUsers(searchTerm) {
    const filtered = users.filter(user =>
        user.name?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        user.email?.toLowerCase().includes(searchTerm.toLowerCase())
    );
    renderUsers(filtered);
}

function openUserModal(userId = null) {
    const modal = document.getElementById('user-modal');
    const form = document.getElementById('user-form');
    const title = document.getElementById('modal-title');

    if (userId) {
        const user = users.find(u => u.id === userId);
        if (user) {
            document.getElementById('user-id').value = user.id;
            document.getElementById('user-email').value = user.email || '';
            document.getElementById('user-name').value = user.name || '';
            document.getElementById('user-role').value = user.role || 'user';
            document.getElementById('user-status').value = user.status || 'active';
            document.getElementById('user-password').value = '';
            title.textContent = 'Edit User';
        }
    } else {
        form.reset();
        document.getElementById('user-id').value = '';
        title.textContent = 'Add User';
    }

    modal.style.display = 'block';
}

function closeUserModal() {
    document.getElementById('user-modal').style.display = 'none';
    document.getElementById('user-form').reset();
}

async function saveUser() {
    const id = document.getElementById('user-id').value;
    const email = document.getElementById('user-email').value;
    const name = document.getElementById('user-name').value;
    const role = document.getElementById('user-role').value;
    const status = document.getElementById('user-status').value;
    const password = document.getElementById('user-password').value;

    const data = { email, name, role, status };
    if (password) data.password = password;

    try {
        if (id) {
            await adminAPI.updateUser(id, data);
            showMessage('User updated successfully', 'success');
        } else {
            await adminAPI.createUser(data);
            showMessage('User created successfully', 'success');
        }
        closeUserModal();
        await loadUsers();
        await loadStats();
    } catch (error) {
        showMessage(error.message || 'Failed to save user', 'error');
    }
}

async function editUser(id) {
    openUserModal(id);
}

async function deleteUser(id) {
    if (!confirm('Are you sure you want to delete this user?')) return;

    try {
        await adminAPI.deleteUser(id);
        showMessage('User deleted successfully', 'success');
        await loadUsers();
        await loadStats();
    } catch (error) {
        showMessage(error.message || 'Failed to delete user', 'error');
    }
}

async function loadSettings() {
    try {
        const response = await adminAPI.getSettings();
        if (response.success && response.data && response.data.admin_verification_code) {
            document.getElementById('verification-code-setting').value = response.data.admin_verification_code;
        }
    } catch (error) {
        console.error('Failed to load settings:', error);
    }
}

async function saveVerificationCode() {
    const code = document.getElementById('verification-code-setting').value;
    if (!code) {
        showMessage('Please enter a verification code', 'error');
        return;
    }

    try {
        await adminAPI.updateSettings({ admin_verification_code: code });
        showMessage('Verification code updated successfully', 'success');
    } catch (error) {
        showMessage(error.message || 'Failed to update verification code', 'error');
    }
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

async function loadSettingsForUI() {
    try {
        const response = await settingsAPI.getSettings();
        if (response.success && response.data) {
            currentSettings = response.data;
            applyThemeFromSettings(currentSettings);
            toggleNotificationsVisibility(currentSettings);
        }
    } catch (error) {
        console.error('Failed to load settings:', error);
    }
}

function applyThemeFromSettings(settings) {
    const theme = settings?.theme || 'light';
    const highContrast = Boolean(settings?.high_contrast);
    document.body.classList.remove('theme-dark', 'theme-light', 'high-contrast');
    if (theme === 'dark') {
        document.body.classList.add('theme-dark');
    } else if (theme === 'system') {
        const prefersDark = window.matchMedia && window.matchMedia('(prefers-color-scheme: dark)').matches;
        document.body.classList.add(prefersDark ? 'theme-dark' : 'theme-light');
    } else {
        document.body.classList.add('theme-light');
    }
    if (highContrast) {
        document.body.classList.add('high-contrast');
    }
}

function notificationsEnabled(settings) {
    if (!settings) return false;
    return Boolean(
        settings.email_notifications ||
        settings.sms_notifications ||
        settings.push_notifications ||
        settings.notification_messages ||
        settings.notification_alerts ||
        settings.notification_promotions ||
        settings.notification_security
    );
}

function toggleNotificationsVisibility(settings) {
    const button = document.getElementById('notifications-button');
    if (!button) return;
    button.style.display = notificationsEnabled(settings) ? 'inline-flex' : 'none';
}

async function refreshBadges() {
    await Promise.all([refreshMessageBadge(), refreshNotificationBadge()]);
}

async function refreshMessageBadge() {
    try {
        const response = await messagingAPI.getUnreadCount();
        const count = response.count || 0;
        updateBadge(document.getElementById('messages-badge'), count);
    } catch (error) {
        updateBadge(document.getElementById('messages-badge'), 0);
    }
}

async function refreshNotificationBadge() {
    if (!notificationsEnabled(currentSettings)) return;
    try {
        const response = await notificationAPI.getUnreadCount();
        const count = response.count || 0;
        updateBadge(document.getElementById('notifications-badge'), count);
    } catch (error) {
        updateBadge(document.getElementById('notifications-badge'), 0);
    }
}

function updateBadge(element, count) {
    if (!element) return;
    if (count > 0) {
        element.textContent = count;
        element.style.display = 'inline-flex';
    } else {
        element.textContent = '';
        element.style.display = 'none';
    }
}

async function loadNotificationsList() {
    const list = document.getElementById('notifications-list');
    if (!list || !notificationsEnabled(currentSettings)) return;
    try {
        const response = await notificationAPI.getNotifications({ unreadOnly: false, limit: 5 });
        const notifications = response.data || [];
        if (!notifications.length) {
            list.innerHTML = '<div class="muted-text">No notifications.</div>';
            return;
        }
        list.innerHTML = notifications.map((notification) => {
            const title = notification.title || 'Notification';
            const content = notification.message || '';
            const unreadClass = notification.is_read ? '' : 'unread';
            return `
                <div class="notification-item ${unreadClass}" data-id="${notification.id}">
                    <strong>${escapeHtml(title)}</strong>
                    <div class="muted-text">${escapeHtml(content)}</div>
                </div>
            `;
        }).join('');
        list.querySelectorAll('.notification-item').forEach((item) => {
            item.addEventListener('click', async () => {
                const id = item.getAttribute('data-id');
                if (id) {
                    await notificationAPI.markAsRead(id);
                    await refreshNotificationBadge();
                    loadNotificationsList();
                }
            });
        });
    } catch (error) {
        list.innerHTML = '<div class="muted-text">Failed to load notifications.</div>';
    }
}

async function loadConversations() {
    try {
        const response = await messagingAPI.getConversations();
        conversations = response.data || [];
        renderConversations();
    } catch (error) {
        renderConversations();
    }
}

function renderConversations() {
    const list = document.getElementById('conversation-list');
    if (!list) return;
    if (!conversations.length) {
        list.innerHTML = '<div class="muted-text">No conversations yet.</div>';
        return;
    }

    list.innerHTML = conversations.map((conversation) => {
        const otherId = getOtherParticipantId(conversation);
        const label = otherId ? getUserLabel(otherId) : 'Unknown user';
        const isActive = activeConversationId === conversation.id;
        return `
            <button class="message-user ${isActive ? 'active' : ''}" data-id="${conversation.id}">
                <span>${escapeHtml(label)}</span>
                <span class="muted-text">${conversation.last_message_at ? formatDate(conversation.last_message_at) : ''}</span>
            </button>
        `;
    }).join('');

    list.querySelectorAll('.message-user').forEach((button) => {
        button.addEventListener('click', async () => {
            const conversationId = button.getAttribute('data-id');
            if (!conversationId) return;
            activeConversationId = conversationId;
            activeRecipientId = getOtherParticipantId(conversations.find((c) => c.id === conversationId));
            renderConversations();
            await loadMessages(conversationId);
        });
    });
}

function getOtherParticipantId(conversation) {
    if (!conversation || !currentUserId) return null;
    return conversation.participant1_id === currentUserId
        ? conversation.participant2_id
        : conversation.participant1_id;
}

function getUserLabel(userId) {
    const user = users.find((u) => u.id === userId);
    if (user) return user.name || user.email || userId;
    return `User ${userId.slice(0, 8)}`;
}

async function loadMessages(conversationId) {
    const thread = document.getElementById('message-thread');
    if (!thread) return;
    thread.innerHTML = '<div class="muted-text">Loading messages...</div>';
    try {
        const response = await messagingAPI.getMessages(conversationId, 50);
        const messages = response.data || [];
        if (!messages.length) {
            thread.innerHTML = '<div class="muted-text">No messages yet.</div>';
            return;
        }
        thread.innerHTML = messages.reverse().map((message) => {
            const isSent = message.sender_id === currentUserId;
            return `
                <div class="message-bubble ${isSent ? 'sent' : 'received'}">
                    <div>${escapeHtml(message.content || '')}</div>
                    <div class="message-meta">${formatDate(message.created_at)}</div>
                </div>
            `;
        }).join('');

        for (const message of messages) {
            if (!message.is_read && message.recipient_id === currentUserId) {
                await messagingAPI.markAsRead(message.id);
            }
        }
        await refreshMessageBadge();
    } catch (error) {
        thread.innerHTML = '<div class="muted-text">Failed to load messages.</div>';
    }
}

async function sendMessage() {
    const contentEl = document.getElementById('message-content');
    const recipientId = document.getElementById('message-recipient-id').value || activeRecipientId;
    const content = contentEl.value.trim();
    if (!recipientId) {
        showMessage('Select a user to message', 'error');
        return;
    }
    if (!content) {
        showMessage('Message cannot be empty', 'error');
        return;
    }

    try {
        await messagingAPI.sendMessage({
            recipient_id: recipientId,
            content,
        });
        contentEl.value = '';
        await loadConversations();
        const conversation = conversations.find((conv) =>
            [conv.participant1_id, conv.participant2_id].includes(recipientId)
        );
        if (conversation) {
            activeConversationId = conversation.id;
            activeRecipientId = recipientId;
            renderConversations();
            await loadMessages(conversation.id);
        }
        await refreshMessageBadge();
    } catch (error) {
        showMessage(error.message || 'Failed to send message', 'error');
    }
}

async function searchUsersForMessaging() {
    const query = document.getElementById('message-search-input').value.trim();
    const resultsEl = document.getElementById('message-search-results');
    if (!resultsEl) return;
    if (!query) {
        resultsEl.innerHTML = '';
        return;
    }

    try {
        const response = await searchAPI.searchUsers(query, 5);
        const results = response.data?.data?.results || [];
        const matched = results.filter((item) => item.type === 'user').map((item) => item.data);
        if (!matched.length) {
            resultsEl.innerHTML = '<div class="muted-text">No users found.</div>';
            return;
        }
        resultsEl.innerHTML = matched.map((user) => `
            <div class="message-user">
                <span>${escapeHtml(user.name || user.email || user.id)}</span>
                <button class="btn btn-secondary" data-id="${user.id}" type="button">Message</button>
            </div>
        `).join('');
        resultsEl.querySelectorAll('button[data-id]').forEach((button) => {
            button.addEventListener('click', () => {
                const recipientId = button.getAttribute('data-id');
                if (!recipientId) return;
                activeRecipientId = recipientId;
                document.getElementById('message-recipient-id').value = recipientId;
                showMessage('Recipient selected. Write your message below.', 'success');
            });
        });
    } catch (error) {
        resultsEl.innerHTML = '<div class="muted-text">Failed to search users.</div>';
    }
}

function formatDate(value) {
    if (!value) return '';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return '';
    return date.toLocaleString();
}


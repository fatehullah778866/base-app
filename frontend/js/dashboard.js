// Dashboard functionality
let items = [];
let currentUserId = null;
let currentSettings = null;
let conversations = [];
let activeConversationId = null;
let activeRecipientId = null;

document.addEventListener('DOMContentLoaded', async () => {
    if (!requireAuth()) return;

    const user = JSON.parse(localStorage.getItem('user') || '{}');
    currentUserId = user.id;
    if (user.role === 'admin') {
        redirect('admin-dashboard.html');
        return;
    }

    setupNavbar();
    await loadSettingsForUI();
    await loadItems();
    setupEventListeners();
    await loadConversations();
    await refreshBadges();
});

function setupNavbar() {
    const user = JSON.parse(localStorage.getItem('user') || '{}');
    const userNameEl = document.getElementById('user-name');
    const profileIcon = document.getElementById('profile-icon');
    const dropdown = document.getElementById('dropdown-menu');

    if (userNameEl) userNameEl.textContent = user.name || user.email || 'User';
    applyAvatar(profileIcon, user.name || user.email || 'User', user.photo_url);

    profileIcon?.addEventListener('click', () => {
        dropdown?.classList.toggle('show');
    });

    document.addEventListener('click', (e) => {
        if (!dropdown?.contains(e.target) && !profileIcon?.contains(e.target)) {
            dropdown?.classList.remove('show');
        }
        const notificationsMenu = document.getElementById('notifications-menu');
        const notificationsButton = document.getElementById('notifications-button');
        if (!notificationsMenu?.contains(e.target) && !notificationsButton?.contains(e.target)) {
            notificationsMenu?.classList.remove('show');
        }
    });

    document.getElementById('logout-link')?.addEventListener('click', async (e) => {
        e.preventDefault();
        localStorage.clear();
        redirect('index.html');
    });

    document.getElementById('profile-link')?.addEventListener('click', (e) => {
        e.preventDefault();
        redirect('settings.html');
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
    const safeName = (name || 'U').trim();
    const initial = safeName ? safeName[0].toUpperCase() : 'U';
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
    document.getElementById('add-item-btn')?.addEventListener('click', () => {
        openModal();
    });

    document.getElementById('item-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveItem();
    });

    document.getElementById('cancel-btn')?.addEventListener('click', closeModal);
    document.querySelector('.close')?.addEventListener('click', closeModal);

    // Enhanced search with backend integration
    let searchTimeout;
    document.getElementById('search-input')?.addEventListener('input', (e) => {
        const query = e.target.value.trim();
        clearTimeout(searchTimeout);
        
        if (query.length === 0) {
            document.getElementById('search-results').style.display = 'none';
            renderItems(items);
            return;
        }
        
        // Debounce search
        searchTimeout = setTimeout(() => {
            performSearch(query);
        }, 300);
    });
    
    document.getElementById('search-input')?.addEventListener('keypress', (e) => {
        if (e.key === 'Enter') {
            e.preventDefault();
            const query = e.target.value.trim();
            if (query) {
                performSearch(query);
            }
        }
    });
    
    document.getElementById('advanced-search-btn')?.addEventListener('click', () => {
        const filters = document.getElementById('search-filters');
        filters.style.display = filters.style.display === 'none' ? 'block' : 'none';
    });
    
    document.getElementById('apply-search-filters')?.addEventListener('click', () => {
        const query = document.getElementById('search-input').value.trim();
        performAdvancedSearch(query);
    });
    
    document.getElementById('clear-search-filters')?.addEventListener('click', () => {
        document.getElementById('search-type').value = 'all';
        document.getElementById('search-location').value = '';
        document.getElementById('search-country').value = '';
        document.getElementById('search-city').value = '';
        document.getElementById('search-date-from').value = '';
        document.getElementById('search-date-to').value = '';
        document.getElementById('search-input').value = '';
        document.getElementById('search-results').style.display = 'none';
        renderItems(items);
    });
    
    document.getElementById('close-search-results')?.addEventListener('click', () => {
        document.getElementById('search-results').style.display = 'none';
        document.getElementById('search-input').value = '';
        renderItems(items);
    });

    window.addEventListener('click', (e) => {
        if (e.target.classList.contains('modal')) {
            closeModal();
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

async function loadItems() {
    try {
        const response = await dashboardAPI.getItems();
        if (response.success) {
            items = response.data || [];
            renderItems(items);
        }
    } catch (error) {
        showMessage('Failed to load items', 'error');
    }
}

function renderItems(itemsToRender) {
    const container = document.getElementById('items-list');
    if (!container) return;

    if (itemsToRender.length === 0) {
        container.innerHTML = '<div class="loading">No items found. Click "Add Item" to create one.</div>';
        return;
    }

    container.innerHTML = itemsToRender.map(item => `
        <div class="item-card">
            <h3>${escapeHtml(item.title || 'Untitled')}</h3>
            <p>${escapeHtml(item.description || '')}</p>
            <div class="item-actions">
                <button class="btn btn-primary" onclick="editItem('${item.id}')">Edit</button>
                <button class="btn btn-secondary" onclick="deleteItem('${item.id}')">Delete</button>
            </div>
        </div>
    `).join('');
}

// Backend-powered search
async function performSearch(query) {
    if (!query || query.length < 2) {
        return;
    }
    
    try {
        const response = await searchAPI.search({
            query: query,
            type: 'all',
            limit: 50
        });
        
        if (response.success && response.data) {
            displaySearchResults(response.data);
        }
    } catch (error) {
        console.error('Search error:', error);
        // Fallback to client-side filtering
        filterItems(query);
    }
}

async function performAdvancedSearch(query) {
    const type = document.getElementById('search-type').value;
    const location = document.getElementById('search-location').value.trim();
    const country = document.getElementById('search-country').value.trim();
    const city = document.getElementById('search-city').value.trim();
    const dateFrom = document.getElementById('search-date-from').value;
    const dateTo = document.getElementById('search-date-to').value;
    
    const params = {
        type: type || 'all',
        limit: 50
    };
    
    if (query) params.query = query;
    if (location) params.location = location;
    if (country) params.country = country;
    if (city) params.city = city;
    if (dateFrom) params.dateFrom = dateFrom;
    if (dateTo) params.dateTo = dateTo;
    
    try {
        const response = await searchAPI.searchAdvanced(params);
        if (response.success && response.data) {
            displaySearchResults(response.data);
        }
    } catch (error) {
        showMessage('Search failed: ' + (error.message || 'Unknown error'), 'error');
    }
}

function displaySearchResults(searchData) {
    const resultsContainer = document.getElementById('search-results');
    const resultsContent = document.getElementById('search-results-content');
    
    if (!resultsContainer || !resultsContent) return;
    
    const results = searchData.data?.results || [];
    const count = searchData.data?.count || 0;
    
    if (count === 0) {
        resultsContent.innerHTML = '<div class="muted-text">No results found.</div>';
        resultsContainer.style.display = 'block';
        return;
    }
    
    // Group results by type
    const grouped = {};
    results.forEach(result => {
        if (!grouped[result.type]) {
            grouped[result.type] = [];
        }
        grouped[result.type].push(result);
    });
    
    let html = `<div class="search-summary">Found ${count} result(s)</div>`;
    
    // Dashboard items
    if (grouped.dashboard_item) {
        html += '<div class="search-group"><h4>Dashboard Items</h4>';
        html += grouped.dashboard_item.map(item => `
            <div class="search-result-item" onclick="openItemFromSearch('${item.id}')">
                <strong>${escapeHtml(item.title || 'Untitled')}</strong>
                <div class="muted-text">${escapeHtml(item.description || '')}</div>
            </div>
        `).join('');
        html += '</div>';
    }
    
    // Messages
    if (grouped.message) {
        html += '<div class="search-group"><h4>Messages</h4>';
        html += grouped.message.map(item => `
            <div class="search-result-item">
                <strong>${escapeHtml(item.title || 'Message')}</strong>
                <div class="muted-text">${escapeHtml(item.description || '')}</div>
            </div>
        `).join('');
        html += '</div>';
    }
    
    // Users
    if (grouped.user) {
        html += '<div class="search-group"><h4>Users</h4>';
        html += grouped.user.map(item => `
            <div class="search-result-item" onclick="selectUserFromSearch('${item.id}')">
                <strong>${escapeHtml(item.title || 'User')}</strong>
            </div>
        `).join('');
        html += '</div>';
    }
    
    // Notifications
    if (grouped.notification) {
        html += '<div class="search-group"><h4>Notifications</h4>';
        html += grouped.notification.map(item => `
            <div class="search-result-item">
                <strong>${escapeHtml(item.title || 'Notification')}</strong>
                <div class="muted-text">${escapeHtml(item.description || '')}</div>
            </div>
        `).join('');
        html += '</div>';
    }
    
    resultsContent.innerHTML = html;
    resultsContainer.style.display = 'block';
}

function openItemFromSearch(itemId) {
    const item = items.find(i => i.id === itemId);
    if (item) {
        openModal(itemId);
        document.getElementById('search-results').style.display = 'none';
    }
}

function selectUserFromSearch(userId) {
    activeRecipientId = userId;
    document.getElementById('message-recipient-id').value = userId;
    document.getElementById('messages-section').scrollIntoView({ behavior: 'smooth' });
    document.getElementById('search-results').style.display = 'none';
    showMessage('User selected. You can now send a message.', 'success');
}

// Fallback client-side filtering
function filterItems(searchTerm) {
    const filtered = items.filter(item =>
        item.title?.toLowerCase().includes(searchTerm.toLowerCase()) ||
        item.description?.toLowerCase().includes(searchTerm.toLowerCase())
    );
    renderItems(filtered);
}

function openModal(itemId = null) {
    const modal = document.getElementById('item-modal');
    const form = document.getElementById('item-form');
    const title = document.getElementById('modal-title');

    if (itemId) {
        const item = items.find(i => i.id === itemId);
        if (item) {
            document.getElementById('item-id').value = item.id;
            document.getElementById('item-title').value = item.title || '';
            document.getElementById('item-description').value = item.description || '';
            title.textContent = 'Edit Item';
        }
    } else {
        form.reset();
        document.getElementById('item-id').value = '';
        title.textContent = 'Add Item';
    }

    modal.style.display = 'block';
}

function closeModal() {
    document.getElementById('item-modal').style.display = 'none';
    document.getElementById('item-form').reset();
}

async function saveItem() {
    const id = document.getElementById('item-id').value;
    const title = document.getElementById('item-title').value;
    const description = document.getElementById('item-description').value;

    try {
        if (id) {
            await dashboardAPI.updateItem(id, { title, description });
            showMessage('Item updated successfully', 'success');
        } else {
            await dashboardAPI.createItem({ title, description });
            showMessage('Item created successfully', 'success');
        }
        closeModal();
        await loadItems();
    } catch (error) {
        showMessage(error.message || 'Failed to save item', 'error');
    }
}

async function editItem(id) {
    openModal(id);
}

async function deleteItem(id) {
    if (!confirm('Are you sure you want to delete this item?')) return;

    try {
        await dashboardAPI.deleteItem(id);
        showMessage('Item deleted successfully', 'success');
        await loadItems();
    } catch (error) {
        showMessage(error.message || 'Failed to delete item', 'error');
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
        const label = otherId ? `User ${otherId.slice(0, 8)}` : 'Unknown user';
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
        const users = results.filter((item) => item.type === 'user').map((item) => item.data);
        if (!users.length) {
            resultsEl.innerHTML = '<div class="muted-text">No users found.</div>';
            return;
        }
        resultsEl.innerHTML = users.map((user) => `
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

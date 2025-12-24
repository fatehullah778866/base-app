// Settings functionality
let currentUser = null;
let currentSettings = null;

document.addEventListener('DOMContentLoaded', async () => {
    if (!requireAuth()) return;

    setupNavbar();
    setupSidebarNavigation();
    await loadProfile();
    await loadSettings();
    setupEventListeners();
    await renderActiveSessions();
    syncConnectedButtons();
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
    });

    document.getElementById('logout-link')?.addEventListener('click', async (e) => {
        e.preventDefault();
        localStorage.clear();
        redirect('index.html');
    });
}

async function loadProfile() {
    try {
        const response = await userAPI.getCurrentUser();
        if (response.success && response.data) {
            currentUser = response.data;
            document.getElementById('profile-name').value = currentUser.name || '';
            document.getElementById('profile-email').value = currentUser.email || '';
            document.getElementById('profile-phone').value = currentUser.phone || '';
            const avatarNote = document.getElementById('profile-avatar-current');
            if (avatarNote) {
                avatarNote.textContent = currentUser.photo_url ? `Current: ${currentUser.photo_url}` : 'No profile picture uploaded.';
            }
        }
    } catch (error) {
        console.error('Failed to load profile:', error);
    }
}

async function loadSettings() {
    try {
        const response = await settingsAPI.getSettings();
        if (response.success && response.data) {
            currentSettings = response.data;
            hydrateSettingsForm(currentSettings);
            applyThemeFromSettings(currentSettings);
        }
    } catch (error) {
        console.error('Failed to load settings:', error);
    }
}

function hydrateSettingsForm(settings) {
    document.getElementById('profile-display-name').value = settings.display_name || '';
    document.getElementById('profile-username').value = settings.username || '';
    document.getElementById('profile-bio').value = settings.bio || '';
    document.getElementById('profile-dob').value = settings.date_of_birth || '';

    const securityData = parseSecurityQuestions(settings.security_questions);
    document.getElementById('security-2fa').checked = Boolean(settings.two_factor_enabled);
    document.getElementById('security-question').value = securityData.question || '';
    document.getElementById('security-answer').value = securityData.answer || '';

    document.getElementById('privacy-visibility').value = settings.profile_visibility || 'public';
    document.getElementById('privacy-email').value = settings.email_visibility || 'private';
    document.getElementById('privacy-phone').value = settings.phone_visibility || 'private';
    document.getElementById('privacy-messaging').value = settings.allow_messaging || 'everyone';
    document.getElementById('privacy-search').value = settings.search_visibility ? 'visible' : 'hidden';
    document.getElementById('privacy-sharing').value = settings.data_sharing_enabled ? 'standard' : 'none';

    document.getElementById('notify-email').checked = Boolean(settings.email_notifications);
    document.getElementById('notify-sms').checked = Boolean(settings.sms_notifications);
    document.getElementById('notify-push').checked = Boolean(settings.push_notifications);
    document.getElementById('notify-messages').checked = Boolean(settings.notification_messages);
    document.getElementById('notify-alerts').checked = Boolean(settings.notification_alerts);
    document.getElementById('notify-promotions').checked = Boolean(settings.notification_promotions);

    document.getElementById('pref-language').value = settings.language || 'en';
    document.getElementById('pref-timezone').value = settings.timezone || '';
    document.getElementById('pref-theme').value = settings.theme || 'light';
    document.getElementById('pref-font-size').value = settings.font_size || 'medium';
    document.getElementById('pref-contrast').value = settings.high_contrast ? 'high' : 'standard';
}

function setupSidebarNavigation() {
    const nav = document.getElementById('settings-nav');
    if (!nav) return;

    nav.addEventListener('click', (e) => {
        const target = e.target.closest('.settings-nav-item');
        if (!target) return;
        const section = target.getAttribute('data-section');
        if (!section) return;
        switchSection(section);
    });
}

function switchSection(section) {
    document.querySelectorAll('.settings-nav-item').forEach((item) => {
        item.classList.toggle('active', item.getAttribute('data-section') === section);
    });
    document.querySelectorAll('.settings-panel').forEach((panel) => {
        panel.classList.toggle('active', panel.getAttribute('data-section') === section);
    });
}

function setupEventListeners() {
    document.getElementById('profile-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveProfile();
    });

    document.getElementById('password-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        await changePassword();
    });

    document.getElementById('security-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveSecuritySettings();
    });

    document.getElementById('privacy-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        await savePrivacySettings();
    });

    document.getElementById('notifications-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        await saveNotificationSettings();
    });

    document.getElementById('preferences-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        await savePreferenceSettings();
    });

    document.querySelectorAll('.connect-btn').forEach((button) => {
        button.addEventListener('click', async () => toggleConnectedAccount(button));
    });

    document.getElementById('download-data-btn')?.addEventListener('click', downloadUserData);
    document.getElementById('deactivate-account-btn')?.addEventListener('click', () => handleAccountAction('deactivate'));
    document.getElementById('delete-account-btn')?.addEventListener('click', () => handleAccountAction('delete'));

    document.querySelectorAll('.support-link').forEach((link) => {
        link.addEventListener('click', (e) => {
            e.preventDefault();
            showMessage('Support request noted. A support agent will contact you.', 'success');
        });
    });

    document.getElementById('logout-all-btn')?.addEventListener('click', async () => {
        try {
            await settingsAPI.logoutAllDevices();
            localStorage.clear();
            showMessage('All sessions cleared. Please log in again.', 'success');
            setTimeout(() => redirect('index.html'), 800);
        } catch (error) {
            showMessage(error.message || 'Failed to log out from all devices', 'error');
        }
    });
}

async function saveProfile() {
    const name = document.getElementById('profile-name').value;
    const phone = document.getElementById('profile-phone').value;
    const avatarFile = document.getElementById('profile-avatar').files?.[0];

    const displayName = document.getElementById('profile-display-name').value;
    const username = document.getElementById('profile-username').value;
    const bio = document.getElementById('profile-bio').value;
    const dob = document.getElementById('profile-dob').value;

    try {
        let avatarUrl = currentUser?.photo_url || null;
        if (avatarFile) {
            avatarUrl = await uploadAvatarFile(avatarFile);
        }

        await userAPI.updateProfile({
            name,
            phone,
            photo_url: avatarUrl || null,
        });

        await settingsAPI.updateProfile({
            display_name: displayName,
            username,
            bio,
            date_of_birth: dob,
        });

        showMessage('Profile updated successfully', 'success');
        const user = JSON.parse(localStorage.getItem('user') || '{}');
        user.name = name;
        user.phone = phone;
        if (avatarUrl) {
            user.photo_url = avatarUrl;
        }
        localStorage.setItem('user', JSON.stringify(user));
        setupNavbar();
    } catch (error) {
        showMessage(error.message || 'Failed to update profile', 'error');
    }
}

async function changePassword() {
    const currentPassword = document.getElementById('current-password').value;
    const newPassword = document.getElementById('new-password').value;
    const confirmPassword = document.getElementById('confirm-password').value;

    if (newPassword !== confirmPassword) {
        showMessage('Passwords do not match', 'error');
        return;
    }

    if (newPassword.length < 8) {
        showMessage('Password must be at least 8 characters', 'error');
        return;
    }

    try {
        await settingsAPI.changePassword({
            current_password: currentPassword,
            new_password: newPassword,
        });
        showMessage('Password changed successfully', 'success');
        document.getElementById('password-form').reset();
    } catch (error) {
        showMessage(error.message || 'Failed to change password', 'error');
    }
}

async function saveSecuritySettings() {
    const securityQuestions = JSON.stringify({
        question: document.getElementById('security-question').value,
        answer: document.getElementById('security-answer').value,
    });

    try {
        await settingsAPI.updateSecurity({
            two_factor_enabled: document.getElementById('security-2fa').checked,
            security_questions: securityQuestions,
        });
        showMessage('Security settings saved', 'success');
    } catch (error) {
        showMessage(error.message || 'Failed to save security settings', 'error');
    }
}

async function savePrivacySettings() {
    try {
        await settingsAPI.updatePrivacy({
            profile_visibility: document.getElementById('privacy-visibility').value,
            email_visibility: document.getElementById('privacy-email').value,
            phone_visibility: document.getElementById('privacy-phone').value,
            allow_messaging: document.getElementById('privacy-messaging').value,
            search_visibility: document.getElementById('privacy-search').value === 'visible',
            data_sharing_enabled: document.getElementById('privacy-sharing').value !== 'none',
        });
        showMessage('Privacy settings saved', 'success');
    } catch (error) {
        showMessage(error.message || 'Failed to save privacy settings', 'error');
    }
}

async function saveNotificationSettings() {
    try {
        await settingsAPI.updateNotifications({
            email_notifications: document.getElementById('notify-email').checked,
            sms_notifications: document.getElementById('notify-sms').checked,
            push_notifications: document.getElementById('notify-push').checked,
            notification_messages: document.getElementById('notify-messages').checked,
            notification_alerts: document.getElementById('notify-alerts').checked,
            notification_promotions: document.getElementById('notify-promotions').checked,
        });
        showMessage('Notification settings saved', 'success');
    } catch (error) {
        showMessage(error.message || 'Failed to save notification settings', 'error');
    }
}

async function savePreferenceSettings() {
    try {
        const response = await settingsAPI.updatePreferences({
            language: document.getElementById('pref-language').value,
            timezone: document.getElementById('pref-timezone').value,
            theme: document.getElementById('pref-theme').value,
            font_size: document.getElementById('pref-font-size').value,
            high_contrast: document.getElementById('pref-contrast').value === 'high',
        });
        if (response.success) {
            applyThemeFromSettings({
                theme: document.getElementById('pref-theme').value,
                high_contrast: document.getElementById('pref-contrast').value === 'high',
            });
        }
        showMessage('Preferences saved', 'success');
    } catch (error) {
        showMessage(error.message || 'Failed to save preferences', 'error');
    }
}

async function toggleConnectedAccount(button) {
    const provider = button.getAttribute('data-provider');
    if (!provider) return;

    try {
        const connectedProviders = getConnectedProviders(currentSettings);
        const isConnected = connectedProviders.includes(provider);

        if (isConnected) {
            await settingsAPI.removeConnectedAccount(provider);
        } else {
            await settingsAPI.addConnectedAccount({
                provider,
                email: currentUser?.email || '',
                connected_at: new Date().toISOString(),
            });
        }

        await loadSettings();
        syncConnectedButtons();
        showMessage(isConnected ? `${provider} disconnected` : `${provider} connected`, 'success');
    } catch (error) {
        showMessage(error.message || 'Failed to update connected account', 'error');
    }
}

function syncConnectedButtons() {
    const connectedProviders = getConnectedProviders(currentSettings);
    document.querySelectorAll('.connect-btn').forEach((button) => {
        const provider = button.getAttribute('data-provider');
        const connected = provider ? connectedProviders.includes(provider) : false;
        button.textContent = connected ? 'Disconnect' : 'Connect';
    });
}

async function renderActiveSessions() {
    const sessionsEl = document.getElementById('active-sessions');
    if (!sessionsEl) return;

    sessionsEl.innerHTML = '';
    try {
        const response = await settingsAPI.getSessions();
        if (!response.success || !Array.isArray(response.data)) {
            sessionsEl.innerHTML = '<div class="muted-text">No active sessions found.</div>';
            return;
        }

        response.data.forEach((session) => {
            const div = document.createElement('div');
            div.className = 'settings-item';
            const deviceLabel = session.device_name || session.browser || session.os || 'Unknown device';
            const lastUsed = session.last_used_at || session.created_at;
            div.innerHTML = `<div><strong>${deviceLabel}</strong><p class="muted-text">Last active: ${formatDate(lastUsed)}</p></div>`;
            sessionsEl.appendChild(div);
        });
    } catch (error) {
        sessionsEl.innerHTML = '<div class="muted-text">Failed to load sessions.</div>';
    }
}

async function downloadUserData() {
    try {
        const response = await settingsAPI.exportData();
        if (!response.success) {
            showMessage('Failed to export data', 'error');
            return;
        }
        const blob = new Blob([JSON.stringify(response.data, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const link = document.createElement('a');
        link.href = url;
        link.download = 'baseapp-data.json';
        document.body.appendChild(link);
        link.click();
        link.remove();
        URL.revokeObjectURL(url);
        showMessage('Your data download has started', 'success');
    } catch (error) {
        showMessage(error.message || 'Failed to export data', 'error');
    }
}

async function handleAccountAction(action) {
    const message = action === 'delete'
        ? 'This will schedule account deletion in 5 days. Continue?'
        : 'This will deactivate your account. Continue?';
    if (!window.confirm(message)) return;

    try {
        if (action === 'delete') {
            await settingsAPI.requestAccountDeletion(5);
        } else {
            await settingsAPI.deactivateAccount();
        }
        localStorage.clear();
        showMessage('Account action completed. Redirecting...', 'success');
        setTimeout(() => redirect('index.html'), 800);
    } catch (error) {
        showMessage(error.message || 'Account action failed', 'error');
    }
}

function getConnectedProviders(settings) {
    if (!settings || !settings.connected_accounts) return [];
    try {
        const accounts = JSON.parse(settings.connected_accounts);
        if (!Array.isArray(accounts)) return [];
        return accounts.map((account) => account.provider);
    } catch {
        return [];
    }
}

function parseSecurityQuestions(value) {
    if (!value) return { question: '', answer: '' };
    try {
        const parsed = JSON.parse(value);
        return {
            question: parsed.question || '',
            answer: parsed.answer || '',
        };
    } catch {
        return { question: '', answer: '' };
    }
}

function formatDate(value) {
    if (!value) return 'Unknown';
    const date = new Date(value);
    if (Number.isNaN(date.getTime())) return 'Unknown';
    return date.toLocaleString();
}

function applyThemeFromSettings(settings) {
    const theme = settings?.theme || 'light';
    const highContrast = Boolean(settings?.high_contrast);
    document.body.classList.remove('theme-dark', 'theme-light', 'theme-system', 'high-contrast');

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

async function uploadAvatarFile(file) {
    const formData = new FormData();
    formData.append('file', file);

    const token = localStorage.getItem('access_token');
    const response = await fetch(`${API_BASE_URL}/files/upload/image`, {
        method: 'POST',
        headers: token ? { Authorization: `Bearer ${token}` } : {},
        body: formData,
    });

    const data = await response.json();
    if (!response.ok || !data.success || !data.data) {
        const errorMsg = data?.error?.message || data?.message || 'Failed to upload image';
        throw new Error(errorMsg);
    }

    if (!data.data.stored_name) {
        throw new Error('Upload completed, but image path is missing');
    }

    return `${getApiOrigin()}/uploads/${data.data.stored_name}`;
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

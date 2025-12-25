// Settings page functionality
let currentSettings = null;

window.addEventListener('DOMContentLoaded', async () => {
    await loadAllSettings();
    showSettingsSection('profile'); // Show profile by default
});

// Show specific settings section
function showSettingsSection(section) {
    // Hide all sections
    document.querySelectorAll('.settings-section').forEach(el => el.classList.remove('active'));
    document.querySelectorAll('.settings-menu-item').forEach(el => el.classList.remove('active'));
    
    // Show selected section
    const sectionEl = document.getElementById(`${section}-section`);
    const menuItem = event?.currentTarget || document.querySelector(`[onclick="showSettingsSection('${section}')"]`);
    
    if (sectionEl) {
        sectionEl.classList.add('active');
    }
    if (menuItem) {
        menuItem.classList.add('active');
    }
}

// Load all settings
async function loadAllSettings() {
    try {
        // Load user profile
        const userResponse = await api.get('/users/me');
        const user = userResponse.data || userResponse;
        
        // Load profile picture
        loadProfilePicturePreview(user);
        
        // Populate profile form
        document.getElementById('profile-name').value = user.name || '';
        document.getElementById('profile-email').value = user.email || '';
        document.getElementById('profile-phone').value = user.phone || '';
        document.getElementById('profile-username').value = user.username || '';
        document.getElementById('profile-bio').value = user.bio || user.about_me || '';
        if (user.date_of_birth) {
            document.getElementById('profile-dob').value = user.date_of_birth.split('T')[0];
        }
        
        // Load all settings
        const settingsResponse = await api.get('/users/me/settings');
        currentSettings = settingsResponse.data || settingsResponse || {};
        
        // Load privacy settings
        loadPrivacySettings(currentSettings.privacy || {});
        
        // Load notification settings
        loadNotificationSettings(currentSettings.notifications || {});
        
        // Load account preferences
        loadAccountPreferences(currentSettings.preferences || {});
        
        // Load connected accounts
        await loadConnectedAccounts();
        
        // Load sessions
        await loadSessions();
    } catch (error) {
        console.error('Failed to load settings:', error);
        const errorMsg = error instanceof Error ? error.message : 'Failed to load settings';
        showMessage(errorMsg, 'error');
    }
}

function loadPrivacySettings(privacy) {
    document.getElementById('privacy-profile-visibility').value = privacy.profile_visibility || 'public';
    document.getElementById('privacy-email-visible').checked = privacy.email_visible || false;
    document.getElementById('privacy-phone-visible').checked = privacy.phone_visible || false;
    document.getElementById('privacy-message-who').value = privacy.message_who || 'everyone';
    document.getElementById('privacy-search-visible').checked = privacy.search_visible !== false;
    document.getElementById('privacy-data-sharing').checked = privacy.data_sharing || false;
}

function loadNotificationSettings(notifications) {
    document.getElementById('notif-email').checked = notifications.email_notifications !== false;
    document.getElementById('notif-sms').checked = notifications.sms_notifications || false;
    document.getElementById('notif-push').checked = notifications.push_notifications || false;
    const messageNotifEnabled = notifications.message_notifications !== false;
    document.getElementById('notif-messages').checked = messageNotifEnabled;
    document.getElementById('notif-alerts').checked = notifications.alert_notifications !== false;
    document.getElementById('notif-promotions').checked = notifications.promotion_notifications || false;
    
    // Store in localStorage for polling control
    localStorage.setItem('notifications_enabled', messageNotifEnabled ? 'true' : 'false');
    localStorage.setItem('messaging_enabled', messageNotifEnabled ? 'true' : 'false');
}

function loadAccountPreferences(preferences) {
    document.getElementById('pref-language').value = preferences.language || 'en';
    document.getElementById('pref-timezone').value = preferences.timezone || 'UTC';
    document.getElementById('pref-theme').value = preferences.theme || 'light';
    document.getElementById('pref-font-size').value = preferences.font_size || 'medium';
    document.getElementById('pref-high-contrast').checked = preferences.high_contrast || false;
}

// Profile Picture Functions
let profilePictureFile = null;

function loadProfilePicturePreview(user) {
    const previewEl = document.getElementById('profile-picture-preview');
    const initialEl = document.getElementById('profile-picture-initial');
    
    if (!previewEl || !initialEl) return;
    
    if (user.photo_url) {
        previewEl.innerHTML = `<img src="${escapeHtml(user.photo_url)}" alt="Profile" style="width: 100%; height: 100%; object-fit: cover;">`;
    } else {
        const initial = (user.name || 'U').charAt(0).toUpperCase();
        initialEl.textContent = initial;
        previewEl.innerHTML = `<span id="profile-picture-initial" style="font-size: 2rem; font-weight: bold; color: var(--primary);">${initial}</span>`;
    }
}

async function handleProfilePictureChange(event) {
    const file = event.target.files[0];
    if (!file) return;
    
    // Validate file type
    if (!file.type.startsWith('image/')) {
        showMessage('Please select an image file', 'error');
        return;
    }
    
    // Validate file size (5MB)
    if (file.size > 5 * 1024 * 1024) {
        showMessage('Image size must be less than 5MB', 'error');
        return;
    }
    
    profilePictureFile = file;
    
    // Show preview
    const reader = new FileReader();
    reader.onload = (e) => {
        const previewEl = document.getElementById('profile-picture-preview');
        if (previewEl) {
            previewEl.innerHTML = `<img src="${e.target.result}" alt="Profile" style="width: 100%; height: 100%; object-fit: cover;">`;
        }
    };
    reader.readAsDataURL(file);
    
    // Upload immediately
    await uploadProfilePicture(file);
}

async function uploadProfilePicture(file) {
    try {
        const formData = new FormData();
        formData.append('file', file);
        
        const token = localStorage.getItem('access_token');
        const response = await fetch('http://localhost:8080/v1/files/upload/image', {
            method: 'POST',
            headers: {
                'Authorization': `Bearer ${token}`
            },
            body: formData
        });
        
        if (!response.ok) {
            const error = await response.json();
            throw new Error(error.error?.message || error.message || 'Upload failed');
        }
        
        const data = await response.json();
        const fileInfo = data.data || data;
        
        // Construct the URL - backend serves files at /uploads/
        let fileUrl = `/uploads/${fileInfo.stored_name}`;
        if (fileInfo.path && !fileInfo.path.startsWith('/')) {
            // If path is relative, use it directly
            fileUrl = `/uploads/${fileInfo.stored_name}`;
        } else if (fileInfo.path) {
            // Extract just the filename from full path
            const pathParts = fileInfo.path.split('/');
            fileUrl = `/uploads/${pathParts[pathParts.length - 1]}`;
        }
        
        // Update user profile with photo URL
        await api.put('/users/me', {
            photo_url: fileUrl
        });
        
        showMessage('Profile picture updated successfully', 'success');
        profilePictureFile = null;
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to upload profile picture';
        showMessage(errorMsg, 'error');
    }
}

// Update Profile Settings
async function updateProfileSettings(e) {
    e.preventDefault();
    const updates = {
        name: document.getElementById('profile-name').value,
        email: document.getElementById('profile-email').value,
        phone: document.getElementById('profile-phone').value,
        username: document.getElementById('profile-username').value,
        bio: document.getElementById('profile-bio').value,
        date_of_birth: document.getElementById('profile-dob').value
    };

    try {
        // Upload profile picture if changed
        if (profilePictureFile) {
            await uploadProfilePicture(profilePictureFile);
        }
        
        // Update via profile endpoint
        await api.put('/users/me', {
            name: updates.name,
            email: updates.email,
            phone: updates.phone
        });
        
        // Update via settings endpoint for additional fields
        await api.put('/users/me/settings/profile', {
            username: updates.username,
            bio: updates.bio,
            date_of_birth: updates.date_of_birth
        });
        
        showMessage('Profile updated successfully', 'success');
        profilePictureFile = null;
        
        // Refresh navbar avatar if on dashboard
        if (typeof window.updateAvatar === 'function') {
            // Reload user data and update avatar
            try {
                const userResponse = await api.get('/users/me');
                const user = userResponse.data || userResponse;
                localStorage.setItem('user', JSON.stringify(user));
                if (typeof window.initializeNavbar === 'function') {
                    await window.initializeNavbar();
                }
            } catch (e) {
                console.error('Failed to refresh navbar:', e);
            }
        }
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to update profile';
        showMessage(errorMsg, 'error');
    }
}

// Change Password
async function changePassword(e) {
    e.preventDefault();
    const currentPassword = document.getElementById('current-password').value;
    const newPassword = document.getElementById('new-password').value;
    const confirmPassword = document.getElementById('confirm-password').value;

    if (newPassword !== confirmPassword) {
        showMessage('New passwords do not match', 'error');
        return;
    }

    if (newPassword.length < 8) {
        showMessage('Password must be at least 8 characters long', 'error');
        return;
    }

    try {
        await api.put('/users/me/password', {
            current_password: currentPassword,
            new_password: newPassword
        });
        showMessage('Password changed successfully', 'success');
        document.getElementById('current-password').value = '';
        document.getElementById('new-password').value = '';
        document.getElementById('confirm-password').value = '';
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to change password';
        showMessage(errorMsg, 'error');
    }
}

// Update 2FA
async function update2FA(e) {
    e.preventDefault();
    const enabled = document.getElementById('2fa-enabled').checked;

    try {
        await api.put('/users/me/settings/security', {
            two_factor_enabled: enabled
        });
        showMessage('2FA settings updated', 'success');
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to update 2FA';
        showMessage(errorMsg, 'error');
    }
}

// Update Privacy Settings
async function updatePrivacySettings(e) {
    e.preventDefault();
    const updates = {
        profile_visibility: document.getElementById('privacy-profile-visibility').value,
        email_visible: document.getElementById('privacy-email-visible').checked,
        phone_visible: document.getElementById('privacy-phone-visible').checked,
        message_who: document.getElementById('privacy-message-who').value,
        search_visible: document.getElementById('privacy-search-visible').checked,
        data_sharing: document.getElementById('privacy-data-sharing').checked
    };

    try {
        await api.put('/users/me/settings/privacy', updates);
        showMessage('Privacy settings updated', 'success');
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to update privacy settings';
        showMessage(errorMsg, 'error');
    }
}

// Update Notification Settings
async function updateNotificationSettings(e) {
    e.preventDefault();
    const updates = {
        email_notifications: document.getElementById('notif-email').checked,
        sms_notifications: document.getElementById('notif-sms').checked,
        push_notifications: document.getElementById('notif-push').checked,
        message_notifications: document.getElementById('notif-messages').checked,
        alert_notifications: document.getElementById('notif-alerts').checked,
        promotion_notifications: document.getElementById('notif-promotions').checked
    };

    try {
        await api.put('/users/me/settings/notifications', updates);
        
        // Store in localStorage for real-time polling control
        localStorage.setItem('notifications_enabled', updates.message_notifications ? 'true' : 'false');
        localStorage.setItem('messaging_enabled', updates.message_notifications ? 'true' : 'false');
        
        // Restart polling if enabled
        if (typeof window.startNotificationPolling === 'function') {
            if (updates.message_notifications) {
                window.startNotificationPolling();
            } else {
                if (typeof window.stopNotificationPolling === 'function') {
                    window.stopNotificationPolling();
                }
            }
        }
        
        if (typeof window.startMessagePolling === 'function') {
            if (updates.message_notifications) {
                window.startMessagePolling();
            } else {
                if (typeof window.stopMessagePolling === 'function') {
                    window.stopMessagePolling();
                }
            }
        }
        
        showMessage('Notification settings updated', 'success');
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to update notification settings';
        showMessage(errorMsg, 'error');
    }
}

// Update Account Preferences
async function updateAccountPreferences(e) {
    e.preventDefault();
    const updates = {
        language: document.getElementById('pref-language').value,
        timezone: document.getElementById('pref-timezone').value,
        theme: document.getElementById('pref-theme').value,
        font_size: document.getElementById('pref-font-size').value,
        high_contrast: document.getElementById('pref-high-contrast').checked
    };

    try {
        await api.put('/users/me/settings/preferences', updates);
        showMessage('Preferences updated', 'success');
        
        // Apply theme if changed
        if (updates.theme === 'dark') {
            document.body.classList.add('dark-theme');
        } else {
            document.body.classList.remove('dark-theme');
        }
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to update preferences';
        showMessage(errorMsg, 'error');
    }
}

// Load Connected Accounts
async function loadConnectedAccounts() {
    const listEl = document.getElementById('connected-accounts-list');
    if (!listEl) return;

    try {
        const response = await api.get('/users/me/settings');
        const settings = response.data || response || {};
        const accounts = settings.connected_accounts || [];
        
        if (accounts.length === 0) {
            listEl.innerHTML = '<div class="empty-state"><p>No connected accounts</p></div>';
            return;
        }

        listEl.innerHTML = accounts.map(account => {
            const provider = account.provider || 'unknown';
            const email = account.email || '';
            return `
                <div class="connected-account-item" style="padding: 1rem; border: 1px solid var(--border); border-radius: 6px; margin-bottom: 0.5rem; display: flex; justify-content: space-between; align-items: center;">
                    <div>
                        <strong>${escapeHtml(provider.charAt(0).toUpperCase() + provider.slice(1))}</strong>
                        ${email ? `<p style="margin: 0.25rem 0 0 0; color: var(--text-light); font-size: 0.85rem;">${escapeHtml(email)}</p>` : ''}
                    </div>
                    <button class="btn btn-danger btn-sm" onclick="disconnectAccount('${provider}')">Disconnect</button>
                </div>
            `;
        }).join('');
    } catch (error) {
        listEl.innerHTML = '<div class="empty-state"><p>Failed to load connected accounts</p></div>';
    }
}

async function connectAccount(provider) {
    showMessage(`Connecting ${provider} account... (Feature coming soon)`, 'info');
    // Implementation would depend on OAuth flow
}

async function disconnectAccount(provider) {
    if (!confirm(`Are you sure you want to disconnect your ${provider} account?`)) return;

    try {
        await api.delete('/users/me/settings/connected-accounts', {
            body: { provider }
        });
        showMessage(`${provider} account disconnected`, 'success');
        await loadConnectedAccounts();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to disconnect account';
        showMessage(errorMsg, 'error');
    }
}

// Load Sessions
async function loadSessions() {
    const listEl = document.getElementById('sessions-list');
    if (!listEl) return;

    try {
        listEl.innerHTML = '<div class="loading">Loading...</div>';
        const response = await api.get('/users/me/settings/sessions');
        const sessions = response.data || response || [];
        renderSessions(sessions);
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load sessions';
        listEl.innerHTML = `<div class="empty-state"><p>${escapeHtml(errorMsg)}</p></div>`;
    }
}

function renderSessions(sessions) {
    const listEl = document.getElementById('sessions-list');
    
    if (!sessions || sessions.length === 0) {
        listEl.innerHTML = '<div class="empty-state"><p>No active sessions</p></div>';
        return;
    }

    listEl.innerHTML = sessions.map(session => {
        const device = session.device_name || session.browser || 'Unknown Device';
        const location = session.location_city ? `${session.location_city}, ${session.location_country}` : (session.ip_address || 'Unknown Location');
        const lastActive = session.last_used_at || session.last_active_at || session.created_at || '';
        const isCurrent = session.is_current || false;

        return `
            <div class="session-item" style="padding: 1rem; border-bottom: 1px solid var(--border);">
                <div style="display: flex; justify-content: space-between; align-items: start;">
                    <div>
                        <strong>${escapeHtml(device)}</strong>
                        ${isCurrent ? '<span class="badge" style="margin-left: 0.5rem;">Current</span>' : ''}
                        <p style="margin: 0.5rem 0 0 0; color: var(--text-light); font-size: 0.85rem;">
                            ${escapeHtml(location)} • Last active: ${formatDate(lastActive)}
                        </p>
                    </div>
                </div>
            </div>
        `;
    }).join('');
}

async function logoutAllDevices() {
    if (!confirm('Are you sure you want to logout from all devices?')) return;

    try {
        await api.post('/users/me/settings/sessions/logout-all');
        showMessage('Logged out from all devices', 'success');
        setTimeout(() => {
            logout();
        }, 2000);
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to logout from all devices';
        showMessage(errorMsg, 'error');
    }
}

async function exportData() {
    try {
        const response = await api.get('/users/me/export');
        const blob = new Blob([JSON.stringify(response, null, 2)], { type: 'application/json' });
        const url = URL.createObjectURL(blob);
        const a = document.createElement('a');
        a.href = url;
        a.download = `user-data-${new Date().toISOString()}.json`;
        a.click();
        URL.revokeObjectURL(url);
        showMessage('Data exported successfully', 'success');
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to export data';
        showMessage(errorMsg, 'error');
    }
}

async function deactivateAccount() {
    if (!confirm('Are you sure you want to deactivate your account? You can reactivate it later.')) return;

    try {
        await api.post('/users/me/settings/account/deactivate');
        showMessage('Account deactivated', 'success');
        setTimeout(() => logout(), 2000);
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to deactivate account';
        showMessage(errorMsg, 'error');
    }
}

async function deleteAccount() {
    if (!confirm('Are you sure you want to delete your account? This action cannot be undone.')) return;
    
    const confirmText = prompt('Type "DELETE" to confirm account deletion:');
    if (confirmText !== 'DELETE') {
        showMessage('Account deletion cancelled', 'info');
        return;
    }

    try {
        await api.post('/users/me/settings/account/delete', {
            days_until_deletion: 7
        });
        showMessage('Account deletion requested. You will receive a confirmation email.', 'success');
        setTimeout(() => logout(), 3000);
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to request account deletion';
        showMessage(errorMsg, 'error');
    }
}

function showHelp(topic) {
    const helpMessages = {
        'change-email': 'To change your email, go to Profile Settings and update your email address. You will receive a verification email.',
        'password': 'If you forgot your password, use the "Forgot Password" link on the login page. Passwords must be at least 8 characters long.',
        '2fa': 'Two-Factor Authentication adds an extra layer of security. Enable it in Security Settings.'
    };
    showMessage(helpMessages[topic] || 'Help information coming soon', 'info');
}

function contactSupport() {
    showMessage('Please email support@baseapp.com for assistance', 'info');
}

function reportProblem() {
    showMessage('Please email bugs@baseapp.com to report issues', 'info');
}

function escapeHtml(text) {
    if (!text) return '';
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

function formatDate(dateString) {
    if (!dateString) return '';
    try {
        const date = new Date(dateString);
        return date.toLocaleString();
    } catch (e) {
        return dateString;
    }
}

// Export functions
window.showSettingsSection = showSettingsSection;
window.updateProfileSettings = updateProfileSettings;
window.handleProfilePictureChange = handleProfilePictureChange;
window.changePassword = changePassword;
window.update2FA = update2FA;
window.updatePrivacySettings = updatePrivacySettings;
window.updateNotificationSettings = updateNotificationSettings;
window.updateAccountPreferences = updateAccountPreferences;
window.connectAccount = connectAccount;
window.disconnectAccount = disconnectAccount;
window.logoutAllDevices = logoutAllDevices;
window.exportData = exportData;
window.deactivateAccount = deactivateAccount;
window.deleteAccount = deleteAccount;
window.showHelp = showHelp;
window.contactSupport = contactSupport;
window.reportProblem = reportProblem;

const API_BASE_URL = (typeof CONFIG !== 'undefined' && CONFIG.API_BASE_URL) || 'http://localhost:8080/v1';

class APIClient {
    constructor(baseURL) {
        this.baseURL = baseURL;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        let token = localStorage.getItem('access_token');

        const config = {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...(token && { 'Authorization': `Bearer ${token}` }),
                ...options.headers,
            },
        };

        try {
            // Check if running from file:// protocol (which causes CORS issues)
            if (window.location.protocol === 'file:') {
                throw new Error('CORS_ERROR: Frontend is opened as file://. Please use a web server. Run: python -m http.server 8000 in the frontend folder, then open http://localhost:8000');
            }
            
            let response = await fetch(url, config);
            
            // Check content type before parsing JSON
            const contentType = response.headers.get('content-type') || '';
            let data;
            
            if (contentType.includes('application/json')) {
                try {
                    const text = await response.text();
                    if (!text || text.trim() === '') {
                        throw new Error('Empty response from server');
                    }
                    data = JSON.parse(text);
                } catch (parseError) {
                    console.error('JSON parse error:', parseError);
                    console.error('Response text:', await response.clone().text());
                    throw new Error('Invalid response from server. Please try again.');
                }
            } else {
                const text = await response.text();
                throw new Error(`Server error: ${text || 'Unknown error'}`);
            }

            if (response.status === 401 && token && endpoint !== '/auth/refresh') {
                const refreshed = await this.refreshToken();
                if (refreshed) {
                    token = localStorage.getItem('access_token');
                    config.headers['Authorization'] = `Bearer ${token}`;
                    response = await fetch(url, config);
                    const retryContentType = response.headers.get('content-type') || '';
                    if (retryContentType.includes('application/json')) {
                        const retryText = await response.text();
                        data = JSON.parse(retryText);
                    } else {
                        throw new Error('Invalid response format');
                    }
                } else {
                    localStorage.removeItem('access_token');
                    localStorage.removeItem('refresh_token');
                    localStorage.removeItem('user');
                    if (!window.location.pathname.includes('index.html')) {
                        window.location.href = 'index.html';
                    }
                    throw new Error('Session expired');
                }
            }

            if (!response.ok) {
                const errorMsg = data?.error?.message || data?.error?.code || data?.message || 'Request failed';
                throw new Error(errorMsg);
            }

            return data;
        } catch (error) {
            // Check for file:// protocol error first
            if (error.message && error.message.includes('CORS_ERROR')) {
                throw error;
            }
            
            if (error.name === 'TypeError' && error.message.includes('fetch')) {
                // Check if it's a CORS error
                if (error.message.includes('CORS') || error.message.includes('Failed to fetch')) {
                    throw new Error('CORS error. If opening HTML file directly, use a web server. Run: python -m http.server 8000 in frontend folder, then open http://localhost:8000. Backend URL: ' + url);
                }
                throw new Error('Network error. Please check: 1) Backend server is running at ' + url + ', 2) If opening HTML directly, use a web server (python -m http.server 8000)');
            }
            if (error.message.includes('JSON')) {
                throw new Error('Server response error. Please check the backend server logs.');
            }
            throw error;
        }
    }

    async refreshToken() {
        const refreshToken = localStorage.getItem('refresh_token');
        if (!refreshToken) return false;

        try {
            const response = await fetch(`${this.baseURL}/auth/refresh`, {
                method: 'POST',
                headers: { 'Content-Type': 'application/json' },
                body: JSON.stringify({ refresh_token: refreshToken }),
            });

            const data = await response.json();
            if (response.ok && data.success && data.data) {
                const accessToken = data.data.access_token || data.data.token;
                if (accessToken) localStorage.setItem('access_token', accessToken);
                if (data.data.refresh_token) localStorage.setItem('refresh_token', data.data.refresh_token);
                return true;
            }
            return false;
        } catch {
            return false;
        }
    }

    get(endpoint, options = {}) {
        return this.request(endpoint, { ...options, method: 'GET' });
    }

    post(endpoint, data, options = {}) {
        return this.request(endpoint, {
            ...options,
            method: 'POST',
            body: JSON.stringify(data),
        });
    }

    put(endpoint, data, options = {}) {
        return this.request(endpoint, {
            ...options,
            method: 'PUT',
            body: JSON.stringify(data),
        });
    }

    delete(endpoint, options = {}) {
        return this.request(endpoint, { ...options, method: 'DELETE' });
    }
}

const api = new APIClient(API_BASE_URL);

const authAPI = {
    signup: (data) => api.post('/auth/signup', data),
    login: (data) => api.post('/auth/login', data),
    logout: () => api.post('/auth/logout'),
    forgotPassword: (data) => api.post('/auth/forgot-password', data),
};

const userAPI = {
    getCurrentUser: () => api.get('/users/me'),
    updateProfile: (data) => api.put('/users/me', data),
};

const settingsAPI = {
    getSettings: () => api.get('/users/me/settings'),
    getSessions: () => api.get('/users/me/settings/sessions'),
    logoutAllDevices: () => api.post('/users/me/settings/sessions/logout-all', {}),
    updateProfile: (data) => api.put('/users/me/settings/profile', data),
    updateSecurity: (data) => api.put('/users/me/settings/security', data),
    updatePrivacy: (data) => api.put('/users/me/settings/privacy', data),
    updateNotifications: (data) => api.put('/users/me/settings/notifications', data),
    updatePreferences: (data) => api.put('/users/me/settings/preferences', data),
    addConnectedAccount: (data) => api.post('/users/me/settings/connected-accounts', data),
    removeConnectedAccount: (provider) => api.delete('/users/me/settings/connected-accounts', {
        body: JSON.stringify({ provider }),
    }),
    deactivateAccount: () => api.post('/users/me/settings/account/deactivate', {}),
    reactivateAccount: () => api.post('/users/me/settings/account/reactivate', {}),
    requestAccountDeletion: (daysUntilDeletion) => api.post('/users/me/settings/account/delete', {
        days_until_deletion: daysUntilDeletion,
    }),
    exportData: () => api.get('/users/me/export'),
    changePassword: (data) => api.put('/users/me/password', data),
};

const dashboardAPI = {
    createItem: (data) => api.post('/dashboard/items', data),
    getItems: () => api.get('/dashboard/items'),
    updateItem: (id, data) => api.put(`/dashboard/items/${id}`, data),
    deleteItem: (id) => api.delete(`/dashboard/items/${id}`),
};

const adminAPI = {
    login: (data) => api.post('/admin/login', data),
    createAdminPublic: async (data) => {
        // Public endpoint - no auth token needed
        const url = `${API_BASE_URL}/admin/create`;
        try {
            const response = await fetch(url, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify(data),
            });
            
            const contentType = response.headers.get('content-type') || '';
            let responseData;
            
            if (contentType.includes('application/json')) {
                const text = await response.text();
                if (!text || text.trim() === '') {
                    throw new Error('Empty response from server');
                }
                responseData = JSON.parse(text);
            } else {
                const text = await response.text();
                throw new Error(`Server error: ${text || 'Unknown error'}`);
            }
            
            if (!response.ok) {
                const errorMsg = responseData?.error?.message || responseData?.error?.code || responseData?.message || 'Request failed';
                throw new Error(errorMsg);
            }
            
            return responseData;
        } catch (error) {
            if (error.message.includes('JSON')) {
                throw new Error('Invalid server response. Please check backend logs.');
            }
            throw error;
        }
    },
    getUsers: () => api.get('/admin/users'),
    createUser: (data) => api.post('/admin/users', data),
    updateUser: (id, data) => api.put(`/admin/users/${id}`, data),
    deleteUser: (id) => api.delete(`/admin/users/${id}`),
    getSettings: () => api.get('/admin/settings'),
    updateSettings: (data) => api.put('/admin/settings', data),
};

const messagingAPI = {
    sendMessage: (data) => api.post('/messages', data),
    getConversations: () => api.get('/messages/conversations'),
    getMessages: (conversationId, limit = 50) =>
        api.get(`/messages?conversation_id=${conversationId}&limit=${limit}`),
    markAsRead: (messageId) => api.post('/messages/read', { message_id: messageId }),
    getUnreadCount: () => api.get('/messages/unread-count'),
};

const notificationAPI = {
    getNotifications: (options = {}) => {
        const unreadOnly = options.unreadOnly ? 'true' : 'false';
        const limit = options.limit || 10;
        return api.get(`/notifications?unread_only=${unreadOnly}&limit=${limit}`);
    },
    getUnreadCount: () => api.get('/notifications/unread-count'),
    markAsRead: (id) => api.post('/notifications/read', { id }),
    markAllAsRead: () => api.post('/notifications/read-all', {}),
    deleteNotification: (id) => api.delete('/notifications', { body: JSON.stringify({ id }) }),
};

const searchAPI = {
    // Global search - searches across all entities
    search: (params) => {
        const queryParams = new URLSearchParams();
        if (params.query) queryParams.append('q', params.query);
        if (params.type) queryParams.append('type', params.type);
        if (params.limit) queryParams.append('limit', params.limit);
        if (params.offset) queryParams.append('offset', params.offset);
        if (params.location) queryParams.append('location', params.location);
        if (params.country) queryParams.append('country', params.country);
        if (params.city) queryParams.append('city', params.city);
        if (params.dateFrom) queryParams.append('date_from', params.dateFrom);
        if (params.dateTo) queryParams.append('date_to', params.dateTo);
        if (params.category) queryParams.append('category', params.category);
        if (params.status) queryParams.append('status', params.status);
        if (params.entityId) queryParams.append('entity_id', params.entityId);
        
        return api.get(`/search?${queryParams.toString()}`);
    },
    
    // Search with JSON body (POST)
    searchAdvanced: (params) => api.post('/search', params),
    
    // Quick search methods
    searchUsers: (query, limit = 10) =>
        api.get(`/search?type=users&q=${encodeURIComponent(query)}&limit=${limit}`),
    
    searchDashboard: (query, limit = 20) =>
        api.get(`/search?type=dashboard_items&q=${encodeURIComponent(query)}&limit=${limit}`),
    
    searchMessages: (query, limit = 20) =>
        api.get(`/search?type=messages&q=${encodeURIComponent(query)}&limit=${limit}`),
    
    searchNotifications: (query, limit = 20) =>
        api.get(`/search?type=notifications&q=${encodeURIComponent(query)}&limit=${limit}`),
    
    searchByLocation: (location, country, city, limit = 20) => {
        const params = { type: 'locations', limit };
        if (location) params.location = location;
        if (country) params.country = country;
        if (city) params.city = city;
        return api.get(`/search?${new URLSearchParams(params).toString()}`);
    },
    
    // Search history
    getSearchHistory: (limit = 50) => api.get(`/search/history?limit=${limit}`),
    clearSearchHistory: () => api.delete('/search/history'),
};

function showMessage(message, type = 'success') {
    const messageEl = document.getElementById('message');
    if (messageEl) {
        messageEl.textContent = message;
        messageEl.className = `message ${type}`;
        setTimeout(() => {
            messageEl.className = 'message';
            messageEl.textContent = '';
        }, 5000);
    }
}

function redirect(url) {
    window.location.href = url;
}

function requireAuth() {
    if (!localStorage.getItem('access_token')) {
        redirect('index.html');
        return false;
    }
    return true;
}

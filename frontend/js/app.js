const API_BASE_URL = window.API_BASE_URL || `${window.location.origin}/v1`;

// API Client
class API {
    constructor() {
        this.baseURL = API_BASE_URL;
    }

    async request(endpoint, options = {}) {
        const url = `${this.baseURL}${endpoint}`;
        const token = localStorage.getItem('access_token');

        const config = {
            ...options,
            headers: {
                'Content-Type': 'application/json',
                ...(token && { 'Authorization': `Bearer ${token}` }),
                ...options.headers,
            },
        };

        if (options.body && typeof options.body === 'object') {
            config.body = JSON.stringify(options.body);
        }

        try {
            const response = await fetch(url, config);
            let data;
            
            // Handle 401 Unauthorized - but don't redirect for login/signup endpoints
            if (response.status === 401) {
                const path = window.location.pathname;
                const isAuthEndpoint = endpoint.includes('/auth/login') || 
                                     endpoint.includes('/auth/signup') ||
                                     endpoint.includes('/admin/login');
                
                // Only clear tokens and redirect if not on login page and not an auth endpoint
                // For auth endpoints, let the error be handled normally (show error message)
                if (!isAuthEndpoint && path !== '/' && path !== '/index.html') {
                    localStorage.removeItem('access_token');
                    localStorage.removeItem('user');
                    window.location.href = '/';
                    throw new Error('UNAUTHORIZED'); // Throw to stop execution
                }
            }
            
            // Try to parse JSON, but handle non-JSON responses
            const contentType = response.headers.get('content-type');
            if (contentType && contentType.includes('application/json')) {
                try {
                    data = await response.json();
                } catch (e) {
                    const text = await response.text();
                    if (!response.ok) {
                        throw new Error(text || 'Invalid response from server');
                    }
                    data = {};
                }
            } else {
                const text = await response.text();
                if (!response.ok) {
                    throw new Error(text || 'Invalid response from server');
                }
                data = {};
            }

            if (!response.ok) {
                // Handle different error formats
                let errorMessage = 'Request failed';
                if (data && data.error) {
                    if (typeof data.error === 'string') {
                        errorMessage = data.error;
                    } else if (data.error.message) {
                        errorMessage = data.error.message;
                    } else if (data.error.code) {
                        errorMessage = data.error.code;
                    }
                } else if (data && data.message) {
                    errorMessage = data.message;
                }
                throw new Error(errorMessage);
            }

            return data;
        } catch (error) {
            console.error('API Error:', error);
            // Ensure we always throw an Error object with a string message
            if (error instanceof Error) {
                throw error;
            } else if (typeof error === 'string') {
                throw new Error(error);
            } else {
                throw new Error('An unexpected error occurred');
            }
        }
    }

    get(endpoint) {
        return this.request(endpoint, { method: 'GET' });
    }

    post(endpoint, body) {
        return this.request(endpoint, { method: 'POST', body });
    }

    put(endpoint, body) {
        return this.request(endpoint, { method: 'PUT', body });
    }

    delete(endpoint) {
        return this.request(endpoint, { method: 'DELETE' });
    }
}

const api = new API();

// Auth Functions
async function handleLogin(e) {
    e.preventDefault();
    const email = document.getElementById('login-email').value.trim();
    const password = document.getElementById('login-password').value;

    if (!email || !password) {
        showMessage('Please enter both email and password', 'error');
        return;
    }

    try {
        // Show loading state
        const submitBtn = e.target.querySelector('button[type="submit"]');
        const originalText = submitBtn ? submitBtn.textContent : 'Login';
        if (submitBtn) {
            submitBtn.disabled = true;
            submitBtn.textContent = 'Logging in...';
        }

        const response = await api.post('/auth/login', { email, password });
        
        // Extract data from nested response structure
        // Backend returns: { success: true, data: { user: {...}, session: {...} } }
        const data = response.data || response;
        const user = data.user || {};
        const session = data.session || {};
        
        // Store token (try multiple possible locations)
        const token = session.token || data.token || response.access_token || response.token;
        if (!token) {
            throw new Error('No authentication token received from server');
        }
        localStorage.setItem('access_token', token);
        
        // Store user info
        if (user && Object.keys(user).length > 0) {
            localStorage.setItem('user', JSON.stringify(user));
        } else {
            throw new Error('No user data received from server');
        }
        
        // Check if admin and redirect
        const role = user.role || 'user';
        showMessage('Login successful! Redirecting...', 'success');
        
        // Small delay to show success message
        setTimeout(() => {
            if (role === 'admin') {
                window.location.href = '/admin-dashboard';
            } else {
                window.location.href = '/dashboard';
            }
        }, 500);
    } catch (error) {
        // Restore button state
        const submitBtn = e.target.querySelector('button[type="submit"]');
        if (submitBtn) {
            submitBtn.disabled = false;
            submitBtn.textContent = 'Login';
        }
        
        const errorMsg = getErrorMessage(error);
        console.error('Login error:', error, errorMsg);
        showMessage(errorMsg, 'error');
    }
}

async function handleSignup(e) {
    e.preventDefault();
    const name = document.getElementById('signup-name').value;
    const email = document.getElementById('signup-email').value;
    const password = document.getElementById('signup-password').value;
    const termsAccepted = document.getElementById('signup-terms').checked;

    if (!termsAccepted) {
        showMessage('Please accept the terms and conditions', 'error');
        return;
    }

    try {
        const response = await api.post('/auth/signup', {
            name,
            email,
            password,
            terms_accepted: true,
            terms_version: '1.0'
        });
        
        showMessage('Account created successfully! Please login.', 'success');
        setTimeout(() => switchTab('login'), 2000);
    } catch (error) {
        showMessage(getErrorMessage(error), 'error');
    }
}

async function handleAdminLogin(e) {
    e.preventDefault();
    const email = document.getElementById('admin-email').value.trim();
    const password = document.getElementById('admin-password').value;
    const code = document.getElementById('admin-code').value.trim();

    if (!email || !password || !code) {
        showMessage('Please fill in all fields', 'error');
        return;
    }

    try {
        // First verify admin code (backend expects 'verification_code')
        await api.post('/admin/verify-code', { verification_code: code });
        
        // Then login
        const response = await api.post('/admin/login', { email, password });
        
        // Extract data from nested response structure
        const data = response.data || response;
        const admin = data.admin || data.user || {};
        const session = data.session || {};
        
        // Store token (try multiple possible locations)
        const token = session.token || data.token || response.access_token || response.token;
        if (token) {
            localStorage.setItem('access_token', token);
        }
        
        // Store user info (use admin data if available, otherwise user)
        const userData = admin.id ? admin : user;
        if (userData && Object.keys(userData).length > 0) {
            localStorage.setItem('user', JSON.stringify(userData));
        }
        
        window.location.href = '/admin-dashboard';
    } catch (error) {
        showMessage(getErrorMessage(error), 'error');
    }
}

async function handleAdminVerify(e) {
    e.preventDefault();
    const code = document.getElementById('verify-code').value.trim();

    if (!code) {
        showMessage('Please enter a verification code', 'error');
        return;
    }

    try {
        // Backend expects 'verification_code' not 'code'
        await api.post('/admin/verify-code', { verification_code: code });
        showMessage('Verification code accepted! Please fill in your admin details.', 'success');
        // Store verified code for signup
        sessionStorage.setItem('admin_verified_code', code);
        setTimeout(() => switchTab('admin-signup'), 1500);
    } catch (error) {
        showMessage(getErrorMessage(error), 'error');
    }
}

async function handleAdminSignup(e) {
    e.preventDefault();
    const name = document.getElementById('admin-signup-name').value.trim();
    const email = document.getElementById('admin-signup-email').value.trim();
    const password = document.getElementById('admin-signup-password').value;
    const termsAccepted = document.getElementById('admin-signup-terms').checked;
    const verifiedCode = sessionStorage.getItem('admin_verified_code');
    
    if (!name || !email || !password) {
        showMessage('Please fill in all required fields', 'error');
        return;
    }

    if (!termsAccepted) {
        showMessage('Please accept the terms and conditions', 'error');
        return;
    }

    if (!verifiedCode) {
        showMessage('Please verify the admin code first', 'error');
        switchTab('admin-verify');
        return;
    }

    try {
        const response = await api.post('/admin/create', {
            name,
            email,
            password,
            verification_code: verifiedCode,
            terms_accepted: true,
            terms_version: '1.0'
        });
        
        showMessage('Admin account created successfully! Please login.', 'success');
        sessionStorage.removeItem('admin_verified_code');
        setTimeout(() => switchTab('admin-login'), 2000);
    } catch (error) {
        showMessage(getErrorMessage(error), 'error');
    }
}

// Helper function to extract error message
function getErrorMessage(error) {
    if (error instanceof Error) {
        return error.message;
    } else if (typeof error === 'string') {
        return error;
    } else if (error && typeof error === 'object') {
        if (error.message) {
            return error.message;
        } else if (error.error) {
            if (typeof error.error === 'string') {
                return error.error;
            } else if (error.error.message) {
                return error.error.message;
            }
        }
    }
    return 'An unexpected error occurred';
}

function logout() {
    localStorage.removeItem('access_token');
    localStorage.removeItem('user');
    window.location.href = '/';
}

// Tab Switching
function switchTab(tab) {
    document.querySelectorAll('.form-container').forEach(el => el.classList.remove('active'));
    document.querySelectorAll('.tab-btn').forEach(el => el.classList.remove('active'));

    if (tab === 'login') {
        document.getElementById('login-form').classList.add('active');
        document.querySelectorAll('.tab-btn')[0].classList.add('active');
    } else if (tab === 'signup') {
        document.getElementById('signup-form').classList.add('active');
        document.querySelectorAll('.tab-btn')[1].classList.add('active');
    } else if (tab === 'admin-login') {
        document.getElementById('admin-login-form').classList.add('active');
    } else if (tab === 'admin-verify') {
        document.getElementById('admin-verify-form').classList.add('active');
    } else if (tab === 'admin-signup') {
        document.getElementById('admin-signup-form').classList.add('active');
    } else if (tab === 'forgot-password') {
        document.getElementById('forgot-password-form').classList.add('active');
    } else if (tab === 'reset-password') {
        document.getElementById('reset-password-form').classList.add('active');
    }
}

// Forgot Password
async function handleForgotPassword(e) {
    e.preventDefault();
    const email = document.getElementById('forgot-email').value;

    try {
        await api.post('/auth/forgot-password', { email });
        showMessage('If an account exists with this email, a password reset link has been sent. Please check your email.', 'success');
        setTimeout(() => switchTab('login'), 3000);
    } catch (error) {
        showMessage(getErrorMessage(error), 'error');
    }
}

// Reset Password
async function handleResetPassword(e) {
    e.preventDefault();
    const token = document.getElementById('reset-token').value;
    const newPassword = document.getElementById('reset-new-password').value;
    const confirmPassword = document.getElementById('reset-confirm-password').value;

    if (newPassword !== confirmPassword) {
        showMessage('Passwords do not match', 'error');
        return;
    }

    try {
        await api.post('/auth/reset-password', {
            token: token,
            new_password: newPassword
        });
        showMessage('Password reset successfully! Please login with your new password.', 'success');
        setTimeout(() => switchTab('login'), 2000);
    } catch (error) {
        showMessage(getErrorMessage(error), 'error');
    }
}

// Check for reset token in URL
window.addEventListener('DOMContentLoaded', () => {
    const urlParams = new URLSearchParams(window.location.search);
    const token = urlParams.get('token');
    if (token) {
        document.getElementById('reset-token').value = token;
        switchTab('reset-password');
    }
});

// Message Display
function showMessage(text, type = 'info') {
    const messageEl = document.getElementById('message');
    // Ensure text is always a string
    const messageText = typeof text === 'string' ? text : String(text);
    messageEl.textContent = messageText;
    messageEl.className = `message ${type} active`;
    
    setTimeout(() => {
        messageEl.classList.remove('active');
    }, 3000);
}

// Modal Functions
function openModal(id) {
    const modal = document.getElementById(id);
    if (modal) {
        modal.classList.add('active');
    }
}

function closeModal(id) {
    const modal = document.getElementById(id);
    if (modal) {
        modal.classList.remove('active');
        
        // Clean up map if advanced search modal is closed
        if (id === 'advanced-search-modal' && window._searchMapInstance) {
            try {
                window._searchMapInstance.remove();
                window._searchMapInstance = null;
            } catch (e) {
                console.error('Error cleaning up map:', e);
            }
        }
    }
}

// Close modal on outside click
document.addEventListener('click', (e) => {
    if (e.target.classList.contains('modal')) {
        e.target.classList.remove('active');
    }
});

// Check auth on page load
window.addEventListener('DOMContentLoaded', () => {
    const token = localStorage.getItem('access_token');
    const userStr = localStorage.getItem('user');
    let user = null;
    
    try {
        user = userStr ? JSON.parse(userStr) : null;
    } catch (e) {
        console.error('Failed to parse user from localStorage:', e);
        localStorage.removeItem('user');
    }
    
    const path = window.location.pathname;
    
    if (token && user && user.role) {
        if (path === '/' || path === '/index.html') {
            if (user.role === 'admin') {
                window.location.href = '/admin-dashboard';
            } else {
                window.location.href = '/dashboard';
            }
        } else if (path === '/dashboard' && user.role === 'admin') {
            window.location.href = '/admin-dashboard';
        } else if (path === '/admin-dashboard' && user.role !== 'admin') {
            window.location.href = '/dashboard';
        }
    } else if (path !== '/' && path !== '/index.html') {
        // No valid auth, redirect to login
        window.location.href = '/';
    }
});

// Export API for other scripts
window.api = api;
window.showMessage = showMessage;
window.openModal = openModal;
window.closeModal = closeModal;
window.logout = logout;
window.getErrorMessage = getErrorMessage;
window.switchTab = switchTab;
window.handleAdminVerify = handleAdminVerify;
window.handleAdminSignup = handleAdminSignup;
window.handleForgotPassword = handleForgotPassword;
window.handleResetPassword = handleResetPassword;


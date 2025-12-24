// Main App Logic for index.html
document.addEventListener('DOMContentLoaded', () => {
    console.log('DOM loaded, setting up...');
    setupAuthNavigation();
    setupForms();
    console.log('Setup complete');
});

function setupAuthNavigation() {
    // Show login by default
    showSection('login');

    // Navigation links - use both onclick and addEventListener for maximum compatibility
    function attachLinkHandlers() {
        const showSignupLink = document.getElementById('show-signup');
        const showLoginLink = document.getElementById('show-login');
        const showForgotLink = document.getElementById('show-forgot');
        const showAdminLink = document.getElementById('show-admin');
        
        if (showSignupLink) {
            showSignupLink.onclick = function(e) {
                e.preventDefault();
                e.stopPropagation();
                console.log('Signup clicked');
                showSection('signup');
                return false;
            };
            showSignupLink.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                showSection('signup');
            });
        }

        if (showLoginLink) {
            showLoginLink.onclick = function(e) {
                e.preventDefault();
                e.stopPropagation();
                showSection('login');
                return false;
            };
            showLoginLink.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                showSection('login');
            });
        }

        if (showForgotLink) {
            showForgotLink.onclick = function(e) {
                e.preventDefault();
                e.stopPropagation();
                console.log('Forgot password clicked');
                showSection('forgot');
                return false;
            };
            showForgotLink.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                showSection('forgot');
            });
        }

        if (showAdminLink) {
            showAdminLink.onclick = function(e) {
                e.preventDefault();
                e.stopPropagation();
                console.log('Admin clicked');
                showSection('admin');
                return false;
            };
            showAdminLink.addEventListener('click', function(e) {
                e.preventDefault();
                e.stopPropagation();
                showSection('admin');
            });
        }
    }
    
    // Attach immediately
    attachLinkHandlers();
    
    // Also attach after a short delay to ensure DOM is ready
    setTimeout(attachLinkHandlers, 100);

    document.getElementById('back-to-login')?.addEventListener('click', (e) => {
        e.preventDefault();
        showSection('login');
    });

    document.getElementById('back-to-login-from-admin')?.addEventListener('click', (e) => {
        e.preventDefault();
        showSection('login');
    });

    // Admin buttons - use onclick for reliability
    setTimeout(() => {
        const adminLoginBtn = document.getElementById('admin-login-btn');
        const createAdminBtn = document.getElementById('create-admin-btn');
        
        if (adminLoginBtn) {
            adminLoginBtn.onclick = function(e) {
                e.preventDefault();
                e.stopPropagation();
                document.getElementById('admin-login-form').style.display = 'block';
                document.getElementById('verification-form').style.display = 'none';
                document.getElementById('create-admin-form').style.display = 'none';
            };
        }
        
        if (createAdminBtn) {
            createAdminBtn.onclick = function(e) {
                e.preventDefault();
                e.stopPropagation();
                document.getElementById('verification-form').style.display = 'block';
                document.getElementById('admin-login-form').style.display = 'none';
                document.getElementById('create-admin-form').style.display = 'none';
                const codeInput = document.getElementById('verification-code');
                if (codeInput) codeInput.value = '';
                sessionStorage.removeItem('admin_verification_code');
            };
        }
    }, 100);

    document.getElementById('cancel-admin-login')?.addEventListener('click', () => {
        document.getElementById('admin-login-form').style.display = 'none';
    });

    document.getElementById('cancel-verification')?.addEventListener('click', () => {
        document.getElementById('verification-form').style.display = 'none';
        document.getElementById('verification-code').value = '';
        sessionStorage.removeItem('admin_verification_code');
    });

    document.getElementById('cancel-create-admin')?.addEventListener('click', () => {
        document.getElementById('create-admin-form').style.display = 'none';
        document.getElementById('verification-form').style.display = 'block';
        // Clear signup form
        document.getElementById('new-admin-name').value = '';
        document.getElementById('new-admin-email').value = '';
        document.getElementById('new-admin-password').value = '';
    });
}

function showSection(sectionName) {
    document.querySelectorAll('.auth-section').forEach(section => {
        section.classList.remove('active');
    });
    const section = document.getElementById(`${sectionName}-section`);
    if (section) {
        section.classList.add('active');
    }
    // Hide admin forms when switching
    if (sectionName !== 'admin') {
        const adminLoginForm = document.getElementById('admin-login-form');
        const verificationForm = document.getElementById('verification-form');
        const createAdminForm = document.getElementById('create-admin-form');
        if (adminLoginForm) adminLoginForm.style.display = 'none';
        if (verificationForm) verificationForm.style.display = 'none';
        if (createAdminForm) createAdminForm.style.display = 'none';
        sessionStorage.removeItem('admin_verification_code');
    } else {
        // When showing admin section, reset to initial state and setup buttons
        const adminLoginForm = document.getElementById('admin-login-form');
        const verificationForm = document.getElementById('verification-form');
        const createAdminForm = document.getElementById('create-admin-form');
        const codeInput = document.getElementById('verification-code');
        if (adminLoginForm) adminLoginForm.style.display = 'none';
        if (verificationForm) verificationForm.style.display = 'none';
        if (createAdminForm) createAdminForm.style.display = 'none';
        if (codeInput) codeInput.value = '';
        
        // Setup buttons when admin section is shown
        setTimeout(() => {
            const adminLoginBtn = document.getElementById('admin-login-btn');
            const createAdminBtn = document.getElementById('create-admin-btn');
            
            if (adminLoginBtn) {
                adminLoginBtn.onclick = function(e) {
                    e.preventDefault();
                    e.stopPropagation();
                    if (adminLoginForm) adminLoginForm.style.display = 'block';
                    if (verificationForm) verificationForm.style.display = 'none';
                    if (createAdminForm) createAdminForm.style.display = 'none';
                };
            }
            
            if (createAdminBtn) {
                createAdminBtn.onclick = function(e) {
                    e.preventDefault();
                    e.stopPropagation();
                    if (verificationForm) verificationForm.style.display = 'block';
                    if (adminLoginForm) adminLoginForm.style.display = 'none';
                    if (createAdminForm) createAdminForm.style.display = 'none';
                    if (codeInput) codeInput.value = '';
                    sessionStorage.removeItem('admin_verification_code');
                };
            }
        }, 50);
    }
}

function setupForms() {
    // Login form
    const loginForm = document.getElementById('login-form');
    if (loginForm) {
        loginForm.addEventListener('submit', async (e) => {
            e.preventDefault();
            e.stopPropagation();
            const email = document.getElementById('login-email')?.value;
            const password = document.getElementById('login-password')?.value;

            if (!email || !password) {
                showMessage('Please fill in all fields', 'error');
                return false;
            }

            try {
                const response = await authAPI.login({ email, password });
                if (response.success && response.data) {
                    localStorage.setItem('access_token', response.data.session.token);
                    if (response.data.session.refresh_token) {
                        localStorage.setItem('refresh_token', response.data.session.refresh_token);
                    }
                    localStorage.setItem('user', JSON.stringify(response.data.user));
                    
                    const role = response.data.user.role || 'user';
                    showMessage('Login successful! Redirecting...', 'success');
                    setTimeout(() => {
                        if (role === 'admin') {
                            redirect('admin-dashboard.html');
                        } else {
                            redirect('dashboard.html');
                        }
                    }, 1000);
                }
            } catch (error) {
                showMessage(error.message || 'Login failed', 'error');
            }
            return false;
        });
    }

    // Signup form
    document.getElementById('signup-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const name = document.getElementById('signup-name').value;
        const email = document.getElementById('signup-email').value;
        const password = document.getElementById('signup-password').value;
        const termsAccepted = document.getElementById('signup-terms')?.checked;

        if (!termsAccepted) {
            showMessage('Please accept the Terms of Service to continue', 'error');
            return;
        }

        try {
            const response = await authAPI.signup({
                name,
                email,
                password,
                terms_accepted: true,
                terms_version: 'v1',
                marketing_consent: false,
            });
            if (response.success) {
                showMessage('Account created! Please login.', 'success');
                setTimeout(() => showSection('login'), 2000);
            }
        } catch (error) {
            showMessage(error.message || 'Signup failed', 'error');
        }
    });

    // Forgot password form
    document.getElementById('forgot-form')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const email = document.getElementById('forgot-email').value;

        try {
            await authAPI.forgotPassword({ email });
            showMessage('Password reset link sent to your email', 'success');
        } catch (error) {
            showMessage(error.message || 'Failed to send reset link', 'error');
        }
    });

    // Admin login form
    document.getElementById('admin-login-form-submit')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const email = document.getElementById('admin-email').value;
        const password = document.getElementById('admin-password').value;

        try {
            const response = await adminAPI.login({ email, password });
            if (response.success && response.data) {
                localStorage.setItem('access_token', response.data.session.token);
                if (response.data.session.refresh_token) {
                    localStorage.setItem('refresh_token', response.data.session.refresh_token);
                }
                localStorage.setItem('user', JSON.stringify(response.data.admin));
                showMessage('Admin login successful! Redirecting...', 'success');
                setTimeout(() => redirect('admin-dashboard.html'), 1000);
            }
        } catch (error) {
            showMessage(error.message || 'Admin login failed', 'error');
        }
    });

    // Verification form (Step 1) - Verify code first
    document.getElementById('verification-form-submit')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const verificationCode = document.getElementById('verification-code').value.trim();

        if (!verificationCode) {
            showMessage('Please enter verification code', 'error');
            return;
        }

        try {
            // Verify the code with backend
            const response = await fetch(`${API_BASE_URL}/admin/verify-code`, {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({ verification_code: verificationCode }),
            });

            const contentType = response.headers.get('content-type') || '';
            let data;
            
            if (contentType.includes('application/json')) {
                const text = await response.text();
                if (!text || text.trim() === '') {
                    throw new Error('Empty response from server');
                }
                try {
                    data = JSON.parse(text);
                } catch (parseError) {
                    throw new Error(`Server response error: ${text.substring(0, 100)}`);
                }
            } else {
                const text = await response.text();
                throw new Error(`Server error: ${text || 'Unknown error'}`);
            }

            if (response.ok && data.success) {
                // Code is valid - store and show signup form
                sessionStorage.setItem('admin_verification_code', verificationCode);
                showMessage('Verification successful! Please fill the form below.', 'success');
                document.getElementById('verification-form').style.display = 'none';
                document.getElementById('create-admin-form').style.display = 'block';
            } else {
                // Invalid code - don't show signup form
                const errorMsg = data.error?.message || 'Invalid verification code';
                showMessage(errorMsg, 'error');
                // Keep verification form visible
            }
        } catch (error) {
            if (error.message && error.message.includes('JSON')) {
                showMessage('Server response error. Please check backend.', 'error');
            } else {
                showMessage(error.message || 'Failed to verify code', 'error');
            }
        }
    });

    // Create admin form (Step 2)
    document.getElementById('create-admin-form-submit')?.addEventListener('submit', async (e) => {
        e.preventDefault();
        const verificationCode = sessionStorage.getItem('admin_verification_code');
        const name = document.getElementById('new-admin-name').value;
        const email = document.getElementById('new-admin-email').value;
        const password = document.getElementById('new-admin-password').value;

        if (!verificationCode || !name || !email || !password) {
            showMessage('Please fill all fields', 'error');
            return;
        }

        try {
            const response = await adminAPI.createAdminPublic({
                name,
                email,
                password,
                verification_code: verificationCode,
            });

            if (response.success) {
                showMessage('Admin created successfully! Please login.', 'success');
                sessionStorage.removeItem('admin_verification_code');
                document.getElementById('create-admin-form').style.display = 'none';
                document.getElementById('admin-login-form').style.display = 'block';
                document.getElementById('admin-email').value = email;
                // Clear form
                document.getElementById('new-admin-name').value = '';
                document.getElementById('new-admin-email').value = '';
                document.getElementById('new-admin-password').value = '';
            }
        } catch (error) {
            showMessage(error.message || 'Failed to create admin', 'error');
            // If verification failed, go back to verification form
            if (error.message && error.message.includes('verification')) {
                document.getElementById('create-admin-form').style.display = 'none';
                document.getElementById('verification-form').style.display = 'block';
                sessionStorage.removeItem('admin_verification_code');
            }
        }
    });
}

function showMessage(message, type) {
    const messageEl = document.getElementById('message');
    if (messageEl) {
        messageEl.textContent = message;
        messageEl.className = `message ${type}`;
        messageEl.style.display = 'block';
        setTimeout(() => {
            messageEl.style.display = 'none';
        }, 5000);
    }
}

function redirect(url) {
    window.location.href = url;
}


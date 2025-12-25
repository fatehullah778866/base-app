// Admin Dashboard Functions
let users = [];
let templates = [];
let adminCRUDs = [];

// Load data on page load
window.addEventListener('DOMContentLoaded', async () => {
    // Check if user is authenticated
    const token = localStorage.getItem('access_token');
    const user = JSON.parse(localStorage.getItem('user') || 'null');
    
    if (!token || !user) {
        // Not authenticated, redirect to login
        window.location.href = '/';
        return;
    }
    
    // Check if user is admin
    if (user.role !== 'admin') {
        window.location.href = '/dashboard';
        return;
    }
    
    if (user) {
        const adminNameEl = document.getElementById('admin-name');
        if (adminNameEl) {
            adminNameEl.textContent = user.name || user.email;
        }
    }
    
    // Initialize navbar if available
    if (typeof initializeNavbar === 'function') {
        try {
            await initializeNavbar();
        } catch (error) {
            console.error('Failed to initialize navbar:', error);
        }
    }
    
    // Load data with error handling - same approach as user dashboard
    try {
        await Promise.all([loadUsers(), loadTemplates(), loadCRUDs(), loadAdminSettings()]);
    } catch (error) {
        console.error('Failed to load admin dashboard data:', error);
        // If it's an auth error, redirect will happen in API client
    }
});

// Load Users
async function loadUsers() {
    const grid = document.getElementById('users-grid');
    if (!grid) return;
    
    try {
        grid.innerHTML = '<div class="loading">Loading...</div>';
        const response = await api.get('/admin/users');
        users = response.data || response || [];
        renderUsers();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load users';
        // Don't show error if it's an auth error (redirect will happen)
        if (errorMsg.includes('UNAUTHORIZED') || errorMsg.includes('401')) {
            return; // Will redirect to login
        }
        console.error('Failed to load users:', error);
        grid.innerHTML = `<div class="empty-state"><p>Failed to load users. Please refresh the page.</p></div>`;
        users = [];
    }
}

function renderUsers() {
    const grid = document.getElementById('users-grid');
    if (!grid) return;
    
    if (!users || users.length === 0) {
        grid.innerHTML = '<div class="empty-state"><p>No users yet</p></div>';
        return;
    }

    try {
        grid.innerHTML = users.map(user => {
            // Handle different ID formats (UUID string, object with String() method, etc.)
            let id = '';
            if (user.id) {
                id = typeof user.id === 'string' ? user.id : (user.id.String ? user.id.String() : String(user.id));
            } else if (user.ID) {
                id = typeof user.ID === 'string' ? user.ID : (user.ID.String ? user.ID.String() : String(user.ID));
            } else if (user.user_id) {
                id = typeof user.user_id === 'string' ? user.user_id : String(user.user_id);
            }
            
            const name = user.name || user.Name || 'Unnamed User';
            const email = user.email || user.Email || 'No email';
            const phone = user.phone || user.Phone || '';
            const role = user.role || user.Role || 'user';
            const status = user.status || user.Status || 'active';
            const createdAt = user.created_at || user.CreatedAt || user.createdAt || '';
            
            // Format date if available
            let dateStr = '';
            if (createdAt) {
                try {
                    const date = new Date(createdAt);
                    dateStr = date.toLocaleDateString();
                } catch (e) {
                    dateStr = '';
                }
            }
            
            const statusClass = status === 'active' ? 'badge-success' : status === 'disabled' ? 'badge-danger' : 'badge-secondary';
            const roleClass = role === 'admin' ? 'badge-primary' : 'badge-secondary';
            
            return `
                <div class="card">
                    <div class="card-header">
                        <h4 class="card-title">${escapeHtml(name)}</h4>
                        <div style="display: flex; gap: 0.5rem; align-items: center;">
                            <span class="card-badge ${statusClass}">${escapeHtml(status)}</span>
                            <span class="card-badge ${roleClass}">${escapeHtml(role)}</span>
                        </div>
                    </div>
                    <div class="card-body" style="padding: 1rem;">
                        <p class="card-description" style="margin-bottom: 0.5rem;">
                            <strong>Email:</strong> ${escapeHtml(email)}
                        </p>
                        ${phone ? `<p class="card-description" style="margin-bottom: 0.5rem;"><strong>Phone:</strong> ${escapeHtml(phone)}</p>` : ''}
                        ${dateStr ? `<p class="card-description" style="margin-bottom: 0.5rem; font-size: 0.85rem; color: var(--text-light);"><strong>Joined:</strong> ${escapeHtml(dateStr)}</p>` : ''}
                    </div>
                    <div class="card-actions" style="display: flex; gap: 0.5rem; flex-wrap: wrap;">
                        <button class="btn btn-primary btn-sm" onclick="viewUser('${id}')" title="View Details">👁️ View</button>
                        <button class="btn btn-secondary btn-sm" onclick="editUser('${id}')" title="Edit User">✏️ Edit</button>
                        <button class="btn btn-warning btn-sm" onclick="toggleUserStatus('${id}', '${status}')" title="Toggle Status">
                            ${status === 'active' ? '⏸️ Disable' : '▶️ Enable'}
                        </button>
                        <button class="btn btn-danger btn-sm" onclick="deleteUser('${id}')" title="Delete User">🗑️ Delete</button>
                    </div>
                </div>
            `;
        }).join('');
        console.log('Users rendered successfully');
    } catch (error) {
        console.error('Error rendering users:', error);
        grid.innerHTML = `<div class="empty-state"><p>Error rendering users: ${escapeHtml(error.message)}</p></div>`;
    }
}

// Load Templates
async function loadTemplates() {
    const grid = document.getElementById('templates-grid');
    if (!grid) return;
    
    try {
        grid.innerHTML = '<div class="loading">Loading...</div>';
        const response = await api.get('/admin/cruds/templates');
        templates = response.data || response || [];
        renderTemplates();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load templates';
        // Don't show error if it's an auth error (redirect will happen)
        if (errorMsg.includes('UNAUTHORIZED') || errorMsg.includes('401')) {
            return; // Will redirect to login
        }
        console.error('Failed to load templates:', error);
        grid.innerHTML = `<div class="empty-state"><p>Failed to load templates. Please refresh the page.</p></div>`;
        templates = [];
    }
}

function renderTemplates() {
    const grid = document.getElementById('templates-grid');
    if (!grid) return;
    
    if (!templates || templates.length === 0) {
        grid.innerHTML = '<div class="empty-state"><p>No templates yet. Create one to get started!</p></div>';
        return;
    }

    grid.innerHTML = templates.map(template => {
        const id = template.id || template.ID || '';
        const name = template.display_name || template.name || 'Unnamed';
        const desc = template.description || 'No description';
        const category = template.category || 'general';
        return `
            <div class="card">
                <div class="card-header">
                    <h4 class="card-title">${escapeHtml(name)}</h4>
                    <span class="card-badge">${escapeHtml(category)}</span>
                </div>
                <p class="card-description">${escapeHtml(desc)}</p>
                <div class="card-actions">
                    <button class="btn btn-secondary btn-sm" onclick="editTemplate('${id}')">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteTemplate('${id}')">Delete</button>
                </div>
            </div>
        `;
    }).join('');
}

// Load CRUDs
async function loadCRUDs() {
    const grid = document.getElementById('cruds-grid');
    if (!grid) return;
    
    try {
        grid.innerHTML = '<div class="loading">Loading...</div>';
        const response = await api.get('/admin/cruds/entities');
        adminCRUDs = response.data || response || [];
        renderCRUDs();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load CRUDs';
        // Don't show error if it's an auth error (redirect will happen)
        if (errorMsg.includes('UNAUTHORIZED') || errorMsg.includes('401')) {
            return; // Will redirect to login
        }
        console.error('Failed to load CRUDs:', error);
        grid.innerHTML = `<div class="empty-state"><p>Failed to load CRUDs. Please refresh the page.</p></div>`;
        adminCRUDs = [];
    }
}

function renderCRUDs() {
    const grid = document.getElementById('cruds-grid');
    if (!grid) return;
    
    if (!adminCRUDs || adminCRUDs.length === 0) {
        grid.innerHTML = '<div class="empty-state"><p>No custom CRUDs yet. Create one to get started!</p></div>';
        return;
    }

    grid.innerHTML = adminCRUDs.map(crud => {
        const id = crud.id || crud.ID || '';
        const name = crud.display_name || crud.entity_name || 'Unnamed';
        const desc = crud.description || 'No description';
        return `
            <div class="card">
                <div class="card-header">
                    <h4 class="card-title">${escapeHtml(name)}</h4>
                    <span class="card-badge">CRUD</span>
                </div>
                <p class="card-description">${escapeHtml(desc)}</p>
                <div class="card-actions">
                    <button class="btn btn-secondary btn-sm" onclick="viewCRUD('${id}')">View</button>
                    <button class="btn btn-secondary btn-sm" onclick="editCRUD('${id}')">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteCRUD('${id}')">Delete</button>
                </div>
            </div>
        `;
    }).join('');
}

// Create User
async function createUser(e) {
    e.preventDefault();
    const name = document.getElementById('user-name').value;
    const email = document.getElementById('user-email').value;
    const password = document.getElementById('user-password').value;

    try {
        await api.post('/admin/users', {
            name: name,
            email: email,
            password: password
        });

        showMessage('User created successfully!', 'success');
        closeModal('create-user-modal');
        document.getElementById('user-name').value = '';
        document.getElementById('user-email').value = '';
        document.getElementById('user-password').value = '';
        await loadUsers();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'An error occurred';
        showMessage(errorMsg, 'error');
    }
}

// Template field management
let templateFields = [];

function addTemplateField() {
    const fieldsList = document.getElementById('template-fields-list');
    if (!fieldsList) return;
    
    const fieldId = 'field_' + Date.now();
    
    const fieldHTML = `
        <div class="template-field-item" data-field-id="${fieldId}" style="padding: 1rem; border: 1px solid var(--border); border-radius: 6px; margin-bottom: 0.5rem; background: var(--light);">
            <div style="display: grid; grid-template-columns: 2fr 1fr 1fr auto; gap: 0.5rem; align-items: end;">
                <div class="form-group" style="margin: 0;">
                    <label style="font-size: 0.85rem;">Field Name *</label>
                    <input type="text" class="field-name" placeholder="e.g., title, description" required>
                </div>
                <div class="form-group" style="margin: 0;">
                    <label style="font-size: 0.85rem;">Type *</label>
                    <select class="field-type" required>
                        <option value="string">String</option>
                        <option value="number">Number</option>
                        <option value="boolean">Boolean</option>
                        <option value="date">Date</option>
                        <option value="email">Email</option>
                        <option value="url">URL</option>
                    </select>
                </div>
                <div class="form-group" style="margin: 0;">
                    <label style="font-size: 0.85rem;">Required</label>
                    <input type="checkbox" class="field-required">
                </div>
                <button type="button" class="btn btn-danger btn-sm" onclick="removeTemplateField('${fieldId}')">×</button>
            </div>
            <div class="form-group" style="margin-top: 0.5rem;">
                <input type="text" class="field-description" placeholder="Field description (optional)" style="width: 100%;">
            </div>
        </div>
    `;
    
    fieldsList.insertAdjacentHTML('beforeend', fieldHTML);
    templateFields.push(fieldId);
    updateTemplateSchema();
}

function removeTemplateField(fieldId) {
    const fieldEl = document.querySelector(`[data-field-id="${fieldId}"]`);
    if (fieldEl) {
        fieldEl.remove();
        templateFields = templateFields.filter(id => id !== fieldId);
    }
    updateTemplateSchema();
}

function updateTemplateSchema() {
    const fields = {};
    const required = [];
    
    document.querySelectorAll('.template-field-item').forEach(item => {
        const name = item.querySelector('.field-name')?.value.trim();
        const type = item.querySelector('.field-type')?.value;
        const description = item.querySelector('.field-description')?.value.trim();
        const isRequired = item.querySelector('.field-required')?.checked;
        
        if (name) {
            const field = {
                type: type,
                description: description || name
            };
            
            // Add format for specific types
            if (type === 'email') {
                field.format = 'email';
                field.type = 'string';
            } else if (type === 'url') {
                field.format = 'uri';
                field.type = 'string';
            } else if (type === 'date') {
                field.format = 'date';
                field.type = 'string';
            }
            
            fields[name] = field;
            
            if (isRequired) {
                required.push(name);
            }
        }
    });
    
    const schema = {
        type: 'object',
        properties: fields,
        required: required
    };
    
    const schemaInput = document.getElementById('template-schema');
    if (schemaInput) {
        schemaInput.value = JSON.stringify(schema, null, 2);
    }
}

// Template field management
let templateFields = [];

function addTemplateField() {
    const fieldsList = document.getElementById('template-fields-list');
    if (!fieldsList) return;
    
    const fieldId = 'field_' + Date.now();
    
    const fieldHTML = `
        <div class="template-field-item" data-field-id="${fieldId}" style="padding: 1rem; border: 1px solid var(--border); border-radius: 6px; margin-bottom: 0.5rem; background: var(--light);">
            <div style="display: grid; grid-template-columns: 2fr 1fr 1fr auto; gap: 0.5rem; align-items: end;">
                <div class="form-group" style="margin: 0;">
                    <label style="font-size: 0.85rem;">Field Name *</label>
                    <input type="text" class="field-name" placeholder="e.g., title, description" required>
                </div>
                <div class="form-group" style="margin: 0;">
                    <label style="font-size: 0.85rem;">Type *</label>
                    <select class="field-type" required>
                        <option value="string">String</option>
                        <option value="number">Number</option>
                        <option value="boolean">Boolean</option>
                        <option value="date">Date</option>
                        <option value="email">Email</option>
                        <option value="url">URL</option>
                    </select>
                </div>
                <div class="form-group" style="margin: 0;">
                    <label style="font-size: 0.85rem;">Required</label>
                    <input type="checkbox" class="field-required">
                </div>
                <button type="button" class="btn btn-danger btn-sm" onclick="removeTemplateField('${fieldId}')">×</button>
            </div>
            <div class="form-group" style="margin-top: 0.5rem;">
                <input type="text" class="field-description" placeholder="Field description (optional)" style="width: 100%;">
            </div>
        </div>
    `;
    
    fieldsList.insertAdjacentHTML('beforeend', fieldHTML);
    templateFields.push(fieldId);
    updateTemplateSchema();
}

function removeTemplateField(fieldId) {
    const fieldEl = document.querySelector(`[data-field-id="${fieldId}"]`);
    if (fieldEl) {
        fieldEl.remove();
        templateFields = templateFields.filter(id => id !== fieldId);
    }
    updateTemplateSchema();
}

function updateTemplateSchema() {
    const fields = {};
    const required = [];
    
    document.querySelectorAll('.template-field-item').forEach(item => {
        const name = item.querySelector('.field-name')?.value.trim();
        const type = item.querySelector('.field-type')?.value;
        const description = item.querySelector('.field-description')?.value.trim();
        const isRequired = item.querySelector('.field-required')?.checked;
        
        if (name) {
            const field = {
                type: type,
                description: description || name
            };
            
            // Add format for specific types
            if (type === 'email') {
                field.format = 'email';
                field.type = 'string';
            } else if (type === 'url') {
                field.format = 'uri';
                field.type = 'string';
            } else if (type === 'date') {
                field.format = 'date';
                field.type = 'string';
            }
            
            fields[name] = field;
            
            if (isRequired) {
                required.push(name);
            }
        }
    });
    
    const schema = {
        type: 'object',
        properties: fields,
        required: required
    };
    
    const schemaInput = document.getElementById('template-schema');
    if (schemaInput) {
        schemaInput.value = JSON.stringify(schema, null, 2);
    }
}

// Create Template
async function createTemplate(e) {
    e.preventDefault();
    const name = document.getElementById('template-name').value;
    const displayName = document.getElementById('template-display-name').value;
    const description = document.getElementById('template-description').value;
    const category = document.getElementById('template-category').value;
    const icon = document.getElementById('template-icon')?.value || '📋';
    
    // Generate schema from fields
    updateTemplateSchema();
    const schemaText = document.getElementById('template-schema').value;

    try {
        let schema;
        try {
            schema = JSON.parse(schemaText);
        } catch (err) {
            showMessage('Invalid schema. Please add at least one field.', 'error');
            return;
        }

        // Validate at least one field
        if (!schema.properties || Object.keys(schema.properties).length === 0) {
            showMessage('Please add at least one field to the template', 'error');
            return;
        }

        const templateName = name.toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_]/g, '');
        
        await api.post('/admin/cruds/templates', {
            name: templateName,
            display_name: displayName,
            description: description,
            category: category,
            icon: icon,
            schema: schema
        });

        showMessage('Template created successfully!', 'success');
        closeModal('create-template-modal');
        
        // Reset form
        document.getElementById('template-name').value = '';
        document.getElementById('template-display-name').value = '';
        document.getElementById('template-description').value = '';
        document.getElementById('template-category').value = 'general';
        if (document.getElementById('template-icon')) {
            document.getElementById('template-icon').value = '';
        }
        const fieldsList = document.getElementById('template-fields-list');
        if (fieldsList) {
            fieldsList.innerHTML = '';
        }
        templateFields = [];
        
        await loadTemplates();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to create template';
        showMessage(errorMsg, 'error');
    }
}

// Create CRUD
async function createCRUD(e) {
    e.preventDefault();
    const displayName = document.getElementById('crud-display-name').value;
    const description = document.getElementById('crud-description').value;
    const schemaText = document.getElementById('crud-schema').value;

    try {
        let schema;
        try {
            schema = JSON.parse(schemaText);
        } catch (err) {
            showMessage('Invalid JSON schema', 'error');
            return;
        }

        const entityName = displayName.toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_]/g, '');
        
        await api.post('/admin/cruds/entities', {
            entity_name: entityName,
            display_name: displayName,
            description: description,
            schema: schema
        });

        showMessage('CRUD created successfully!', 'success');
        closeModal('create-crud-modal');
        document.getElementById('crud-display-name').value = '';
        document.getElementById('crud-description').value = '';
        await loadCRUDs();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'An error occurred';
        showMessage(errorMsg, 'error');
    }
}

// Delete Functions
async function deleteUser(id) {
    if (!confirm('Are you sure you want to delete this user?')) return;

    try {
        await api.delete(`/admin/users/${id}`);
        showMessage('User deleted successfully!', 'success');
        await loadUsers();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'An error occurred';
        showMessage(errorMsg, 'error');
    }
}

async function deleteTemplate(id) {
    if (!confirm('Are you sure you want to delete this template?')) return;

    try {
        await api.delete(`/admin/cruds/templates/id/${id}`);
        showMessage('Template deleted successfully!', 'success');
        await loadTemplates();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'An error occurred';
        showMessage(errorMsg, 'error');
    }
}

async function deleteCRUD(id) {
    if (!confirm('Are you sure you want to delete this CRUD?')) return;

    try {
        await api.delete(`/admin/cruds/entities/${id}`);
        showMessage('CRUD deleted successfully!', 'success');
        await loadCRUDs();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'An error occurred';
        showMessage(errorMsg, 'error');
    }
}

// Modal Functions
function openCreateUserModal() {
    openModal('create-user-modal');
}

function openCreateTemplateModal() {
    // Reset fields
    templateFields = [];
    document.getElementById('template-fields-list').innerHTML = '';
    // Add one default field
    addTemplateField();
    openModal('create-template-modal');
}

function openCreateCRUDModal() {
    openModal('create-crud-modal');
}

// Utility
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// User action functions
async function viewUser(id) {
    try {
        const response = await api.get(`/admin/users/${id}`);
        const user = response.data || response;
        
        const modal = document.getElementById('view-user-modal');
        if (!modal) {
            // Create modal if it doesn't exist
            const modalHTML = `
                <div id="view-user-modal" class="modal">
                    <div class="modal-content" style="max-width: 600px;">
                        <span class="close" onclick="closeModal('view-user-modal')">&times;</span>
                        <h3>User Details</h3>
                        <div id="view-user-content"></div>
                    </div>
                </div>
            `;
            document.body.insertAdjacentHTML('beforeend', modalHTML);
        }
        
        const content = document.getElementById('view-user-content');
        if (content) {
            const name = user.name || user.Name || 'Unnamed';
            const email = user.email || user.Email || '';
            const phone = user.phone || user.Phone || '';
            const role = user.role || user.Role || 'user';
            const status = user.status || user.Status || 'active';
            const createdAt = user.created_at || user.CreatedAt || user.createdAt || '';
            
            content.innerHTML = `
                <div class="card">
                    <div class="card-body">
                        <p><strong>Name:</strong> ${escapeHtml(name)}</p>
                        <p><strong>Email:</strong> ${escapeHtml(email)}</p>
                        ${phone ? `<p><strong>Phone:</strong> ${escapeHtml(phone)}</p>` : ''}
                        <p><strong>Role:</strong> ${escapeHtml(role)}</p>
                        <p><strong>Status:</strong> ${escapeHtml(status)}</p>
                        ${createdAt ? `<p><strong>Created:</strong> ${escapeHtml(new Date(createdAt).toLocaleString())}</p>` : ''}
                    </div>
                </div>
            `;
        }
        
        openModal('view-user-modal');
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load user details';
        showMessage(errorMsg, 'error');
    }
}

async function toggleUserStatus(id, currentStatus) {
    const newStatus = currentStatus === 'active' ? 'disabled' : 'active';
    const action = newStatus === 'active' ? 'enable' : 'disable';
    
    if (!confirm(`Are you sure you want to ${action} this user?`)) return;
    
    try {
        await api.post(`/admin/users/${id}/status`, { status: newStatus });
        showMessage(`User ${action}d successfully!`, 'success');
        await loadUsers();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'An error occurred';
        showMessage(errorMsg, 'error');
    }
}

async function editUser(id) {
    try {
        const response = await api.get(`/admin/users/${id}`);
        const user = response.data || response;
        
        // For now, show a message - can be enhanced with a proper edit modal
        showMessage(`Edit user: ${user.name || user.Name || 'User'} (ID: ${id})`, 'info');
        // TODO: Open edit modal with user data
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load user';
        showMessage(errorMsg, 'error');
    }
}

function editTemplate(id) {
    showMessage('Edit template functionality coming soon', 'info');
}

function viewCRUD(id) {
    showMessage('View CRUD functionality coming soon', 'info');
}

function editCRUD(id) {
    showMessage('Edit CRUD functionality coming soon', 'info');
}

// Admin Settings
async function loadAdminSettings() {
    try {
        const response = await api.get('/admin/settings');
        const settings = response.data || response || {};
        
        // Load current verification code
        const currentCodeEl = document.getElementById('current-verification-code');
        if (currentCodeEl) {
            const code = settings.admin_verification_code || 'Kompasstech2025@';
            currentCodeEl.value = code;
        }
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load admin settings';
        // Don't show error if it's an auth error (redirect will happen)
        if (!errorMsg.includes('UNAUTHORIZED') && !errorMsg.includes('401')) {
            console.error('Failed to load admin settings:', error);
        }
    }
}

async function updateVerificationCode(e) {
    e.preventDefault();
    const newCode = document.getElementById('new-verification-code').value.trim();
    const confirmCode = document.getElementById('confirm-verification-code').value.trim();
    
    if (!newCode) {
        showMessage('Please enter a new verification code', 'error');
        return;
    }
    
    if (newCode !== confirmCode) {
        showMessage('Verification codes do not match', 'error');
        return;
    }
    
    if (newCode.length < 6) {
        showMessage('Verification code must be at least 6 characters long', 'error');
        return;
    }
    
    try {
        await api.put('/admin/settings', {
            admin_verification_code: newCode
        });
        
        showMessage('Verification code updated successfully', 'success');
        
        // Update current code display
        const currentCodeEl = document.getElementById('current-verification-code');
        if (currentCodeEl) {
            currentCodeEl.value = newCode;
        }
        
        // Clear form
        document.getElementById('new-verification-code').value = '';
        document.getElementById('confirm-verification-code').value = '';
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to update verification code';
        showMessage(errorMsg, 'error');
    }
}

// Export functions for retry buttons and onclick handlers
window.loadUsers = loadUsers;
window.loadTemplates = loadTemplates;
window.loadCRUDs = loadCRUDs;

// Export user action functions
window.viewUser = viewUser;
window.editUser = editUser;
window.deleteUser = deleteUser;
window.toggleUserStatus = toggleUserStatus;

// Export template and CRUD action functions
window.editTemplate = editTemplate;
window.deleteTemplate = deleteTemplate;
window.viewCRUD = viewCRUD;
window.editCRUD = editCRUD;
window.deleteCRUD = deleteCRUD;

// Export creation functions
window.createUser = createUser;
window.createTemplate = createTemplate;
window.createCRUD = createCRUD;
window.openCreateUserModal = openCreateUserModal;
window.openCreateTemplateModal = openCreateTemplateModal;
window.openCreateCRUDModal = openCreateCRUDModal;
window.addTemplateField = addTemplateField;
window.removeTemplateField = removeTemplateField;
window.updateVerificationCode = updateVerificationCode;


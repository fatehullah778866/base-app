// Dashboard Functions
let userCRUDs = [];
let dashboardItems = [];
let availableTemplates = [];

function ensureArray(value) {
    return Array.isArray(value) ? value : [];
}

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
    
    if (user) {
        const userNameEl = document.getElementById('user-name');
        if (userNameEl) {
            userNameEl.textContent = user.name || user.email;
        }
    }
    
    // Load data with error handling
    try {
        await Promise.all([loadTemplates(), loadCRUDs(), loadDashboardItems()]);
    } catch (error) {
        console.error('Failed to load dashboard data:', error);
        // If it's an auth error, redirect will happen in API client
    }
});

// Load Templates
async function loadTemplates() {
    const grid = document.getElementById('templates-grid');
    if (!grid) return;
    
    try {
        grid.innerHTML = '<div class="loading">Loading...</div>';
        // Users can access templates via /cruds/templates (protected route, not admin-only)
        const response = await api.get('/cruds/templates?active_only=true');
        const activeTemplatesPayload = response.data || response;
        availableTemplates = ensureArray(activeTemplatesPayload);
        // Filter to only active templates
        availableTemplates = availableTemplates.filter(t => t.is_active !== false);
        renderTemplates();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load templates';
        // Don't show error if it's an auth error (redirect will happen)
        if (errorMsg.includes('UNAUTHORIZED') || errorMsg.includes('401')) {
            return; // Will redirect to login
        }
        console.error('Failed to load templates:', error);
        grid.innerHTML = `<div class="empty-state"><p>Failed to load templates. Please refresh the page.</p></div>`;
        availableTemplates = [];
    }
}

function renderTemplates() {
    const grid = document.getElementById('templates-grid');
    if (!grid) return;
    
    if (!availableTemplates || availableTemplates.length === 0) {
        grid.innerHTML = '<div class="empty-state"><p>No templates available. Create a custom CRUD instead.</p></div>';
        return;
    }
    
    grid.innerHTML = availableTemplates.map(template => {
        const name = template.display_name || template.name || 'Unnamed';
        const desc = template.description || 'No description';
        const category = template.category || 'general';
        const icon = template.icon || '📋';
        
        return `
            <div class="card" onclick="createCRUDFromTemplate('${template.name || template.id}')">
                <div class="card-header">
                    <h4 class="card-title">${icon} ${escapeHtml(name)}</h4>
                    <span class="card-badge">${escapeHtml(category)}</span>
                </div>
                <p class="card-description">${escapeHtml(desc)}</p>
                <div class="card-actions">
                    <button class="btn btn-primary btn-sm" onclick="event.stopPropagation(); createCRUDFromTemplate('${template.name || template.id}')">Use Template</button>
                </div>
            </div>
        `;
    }).join('');
}

async function createCRUDFromTemplate(templateName) {
    // Get template details first
    try {
        const templateResponse = await api.get(`/cruds/templates/${templateName}`);
        const template = templateResponse.data || templateResponse;
        
        if (!template || !template.schema) {
            showMessage('Template not found or invalid', 'error');
            return;
        }
        
        // Open modal to create CRUD from template
        const displayName = prompt('Enter a display name for your CRUD:', template.display_name || '');
        if (!displayName) return;
        
        const description = prompt('Enter a description (optional):', template.description || '') || '';
        
        // Use the template create endpoint which handles everything
        try {
            const response = await api.post(`/cruds/templates/${templateName}/create`, {
                display_name: displayName,
                description: description
            });
            
            showMessage('CRUD created from template successfully!', 'success');
            await loadCRUDs();
        } catch (error) {
            // Fallback: create manually using template schema
            const entityName = displayName.toLowerCase().replace(/\s+/g, '_').replace(/[^a-z0-9_]/g, '');
            
            const response = await api.post('/cruds/entities', {
                entity_name: entityName,
                display_name: displayName,
                description: description,
                schema: template.schema
            });
            
            showMessage('CRUD created from template successfully!', 'success');
            await loadCRUDs();
        }
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to create CRUD from template';
        showMessage(errorMsg, 'error');
    }
}

// Load CRUDs
async function loadCRUDs() {
    const grid = document.getElementById('cruds-grid');
    if (!grid) return;
    
    try {
        grid.innerHTML = '<div class="loading">Loading...</div>';
        // Users create their own CRUDs via /cruds/entities
        const response = await api.get('/cruds/entities');
        const userCRUDPayload = response.data || response;
        userCRUDs = ensureArray(userCRUDPayload);
        renderCRUDs();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load CRUDs';
        // Don't show error if it's an auth error (redirect will happen)
        if (errorMsg.includes('UNAUTHORIZED') || errorMsg.includes('401')) {
            return; // Will redirect to login
        }
        console.error('Failed to load CRUDs:', error);
        grid.innerHTML = `<div class="empty-state"><p>Failed to load CRUDs. Please refresh the page.</p></div>`;
        userCRUDs = [];
    }
}

function renderCRUDs() {
    const grid = document.getElementById('cruds-grid');
    if (!userCRUDs || userCRUDs.length === 0) {
        grid.innerHTML = '<div class="empty-state"><p>No CRUDs yet. Create one to get started!</p></div>';
        return;
    }

    grid.innerHTML = userCRUDs.map(crud => {
        const id = crud.id || crud.ID;
        const name = crud.display_name || crud.entity_name || 'Unnamed';
        const desc = crud.description || 'No description';
        return `
            <div class="card">
                <div class="card-header">
                    <h4 class="card-title">${escapeHtml(name)}</h4>
                </div>
                <p class="card-description">${escapeHtml(desc)}</p>
                <div class="card-actions">
                    <button class="btn btn-primary btn-sm" onclick="viewCRUD('${id}')">View</button>
                    <button class="btn btn-secondary btn-sm" onclick="editCRUD('${id}')">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteCRUD('${id}')">Delete</button>
                </div>
            </div>
        `;
    }).join('');
}

// Load Dashboard Items
async function loadDashboardItems() {
    const grid = document.getElementById('items-grid');
    if (!grid) return;
    
    try {
        grid.innerHTML = '<div class="loading">Loading...</div>';
        const response = await api.get('/dashboard/items');
        const itemsPayload = response.data || response;
        dashboardItems = ensureArray(itemsPayload);
        renderDashboardItems();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load items';
        // Don't show error if it's an auth error (redirect will happen)
        if (errorMsg.includes('UNAUTHORIZED') || errorMsg.includes('401')) {
            return; // Will redirect to login
        }
        console.error('Failed to load items:', error);
        grid.innerHTML = `<div class="empty-state"><p>Failed to load items. Please refresh the page.</p></div>`;
        dashboardItems = [];
    }
}

function renderDashboardItems() {
    const grid = document.getElementById('items-grid');
    if (!dashboardItems || dashboardItems.length === 0) {
        grid.innerHTML = '<div class="empty-state"><p>No items yet. Add one to get started!</p></div>';
        return;
    }

    grid.innerHTML = dashboardItems.map(item => {
        const id = item.id || item.ID;
        const title = item.title || 'Untitled';
        const desc = item.description || '';
        return `
            <div class="card">
                <div class="card-header">
                    <h4 class="card-title">${escapeHtml(title)}</h4>
                </div>
                <p class="card-description">${escapeHtml(desc)}</p>
                <div class="card-actions">
                    <button class="btn btn-secondary btn-sm" onclick="editItem('${id}')">Edit</button>
                    <button class="btn btn-danger btn-sm" onclick="deleteItem('${id}')">Delete</button>
                </div>
            </div>
        `;
    }).join('');
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
        
        await api.post('/cruds/entities', {
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

// Create Item
async function createItem(e) {
    e.preventDefault();
    const title = document.getElementById('item-title').value;
    const description = document.getElementById('item-description').value;

    try {
        await api.post('/dashboard/items', {
            title: title,
            description: description
        });

        showMessage('Item added successfully!', 'success');
        closeModal('create-item-modal');
        document.getElementById('item-title').value = '';
        document.getElementById('item-description').value = '';
        await loadDashboardItems();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'An error occurred';
        showMessage(errorMsg, 'error');
    }
}

// Delete CRUD
async function deleteCRUD(id) {
    if (!confirm('Are you sure you want to delete this CRUD?')) return;

    try {
        await api.delete(`/cruds/entities/${id}`);
        showMessage('CRUD deleted successfully!', 'success');
        await loadCRUDs();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'An error occurred';
        showMessage(errorMsg, 'error');
    }
}

// Delete Item
async function deleteItem(id) {
    if (!confirm('Are you sure you want to delete this item?')) return;

    try {
        await api.delete(`/dashboard/items/${id}`);
        showMessage('Item deleted successfully!', 'success');
        await loadDashboardItems();
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'An error occurred';
        showMessage(errorMsg, 'error');
    }
}

// Modal Functions
function openCreateCRUDModal() {
    openModal('create-crud-modal');
}

function openCreateItemModal() {
    openModal('create-item-modal');
}

// Utility
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// Placeholder functions
function viewCRUD(id) {
    showMessage('View CRUD functionality coming soon', 'info');
}

function editCRUD(id) {
    showMessage('Edit CRUD functionality coming soon', 'info');
}

function editItem(id) {
    showMessage('Edit item functionality coming soon', 'info');
}

// Export
window.createCRUD = createCRUD;
window.createItem = createItem;
window.openCreateCRUDModal = openCreateCRUDModal;
window.openCreateItemModal = openCreateItemModal;
window.createCRUDFromTemplate = createCRUDFromTemplate;


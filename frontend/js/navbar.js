// Navbar functionality
let currentUser = null;
let unreadNotifications = 0;
let unreadMessages = 0;

// Initialize navbar on page load
window.addEventListener('DOMContentLoaded', async () => {
    await initializeNavbar();
    await loadNotificationCount();
    await loadMessageCount();
    
    // Start real-time polling
    startMessagePolling();
    startNotificationPolling();
});

async function initializeNavbar() {
    // Check if user is authenticated
    const token = localStorage.getItem('access_token');
    if (!token) {
        // Not authenticated, don't try to load user data
        return;
    }
    
    try {
        // Always try to fetch fresh user data to get latest photo_url
        try {
            const response = await api.get('/users/me');
            const data = response.data || response;
            currentUser = data;
            localStorage.setItem('user', JSON.stringify(data));
            updateAvatar();
        } catch (fetchError) {
            // If it's an auth error, don't fallback - let redirect happen
            const errorMsg = fetchError instanceof Error ? fetchError.message : '';
            if (errorMsg.includes('UNAUTHORIZED') || errorMsg.includes('401')) {
                return; // Will redirect to login
            }
            // Fallback to localStorage if fetch fails for other reasons
            const userStr = localStorage.getItem('user');
            if (userStr) {
                currentUser = JSON.parse(userStr);
                updateAvatar();
            }
        }
    } catch (error) {
        console.error('Failed to load user:', error);
    }
}

function updateAvatar() {
    if (!currentUser) return;
    
    const avatarEl = document.getElementById('user-avatar');
    if (!avatarEl) return;
    
    const initialEl = document.getElementById('avatar-initial');
    
    // If user has a photo URL, show image
    if (currentUser.photo_url) {
        avatarEl.style.background = 'none';
        avatarEl.style.padding = '0';
        avatarEl.innerHTML = `<img src="${escapeHtml(currentUser.photo_url)}" alt="Avatar" style="width: 100%; height: 100%; border-radius: 50%; object-fit: cover;">`;
    } else if (initialEl && currentUser.name) {
        // Show initial
        avatarEl.style.background = 'var(--primary)';
        avatarEl.style.padding = '';
        const initial = currentUser.name.charAt(0).toUpperCase();
        initialEl.textContent = initial;
        if (!avatarEl.contains(initialEl)) {
            avatarEl.innerHTML = `<span id="avatar-initial">${initial}</span>`;
        }
    }
}

function toggleAvatarMenu() {
    const menu = document.getElementById('avatar-menu');
    if (menu) {
        menu.classList.toggle('active');
    }
    
    // Close menu when clicking outside
    document.addEventListener('click', function closeMenu(e) {
        if (!e.target.closest('.avatar-container')) {
            menu.classList.remove('active');
            document.removeEventListener('click', closeMenu);
        }
    });
}

async function viewProfile() {
    closeAvatarMenu();
    await openProfileModal();
}

function closeAvatarMenu() {
    const menu = document.getElementById('avatar-menu');
    if (menu) {
        menu.classList.remove('active');
    }
}

// Profile Modal
async function openProfileModal() {
    openModal('profile-modal');
    await loadProfileCard();
}

async function loadProfileCard() {
    const contentEl = document.getElementById('profile-card-content');
    if (!contentEl) return;
    
    try {
        contentEl.innerHTML = '<div class="loading">Loading profile...</div>';
        const response = await api.get('/users/me');
        const user = response.data || response;
        
        const avatarUrl = user.photo_url || '';
        const avatarDisplay = avatarUrl ? 
            `<img src="${escapeHtml(avatarUrl)}" alt="Avatar" style="width: 100px; height: 100px; border-radius: 50%; object-fit: cover; border: 3px solid var(--primary);">` :
            `<div style="width: 100px; height: 100px; border-radius: 50%; background: var(--primary); color: var(--white); display: flex; align-items: center; justify-content: center; font-size: 2.5rem; font-weight: bold; margin: 0 auto;">${(user.name || 'U').charAt(0).toUpperCase()}</div>`;
        
        contentEl.innerHTML = `
            <div style="text-align: center; padding: 1.5rem;">
                ${avatarDisplay}
                <h3 style="margin: 1rem 0 0.5rem 0; color: var(--text);">${escapeHtml(user.name || 'User')}</h3>
                <p style="color: var(--text-light); margin-bottom: 1.5rem;">${escapeHtml(user.email || '')}</p>
                <div style="text-align: left; background: var(--light); padding: 1rem; border-radius: 6px; margin-top: 1rem;">
                    <div style="margin-bottom: 0.75rem;">
                        <strong style="color: var(--text-light); font-size: 0.85rem;">Phone:</strong>
                        <p style="margin: 0.25rem 0 0 0;">${escapeHtml(user.phone || 'Not provided')}</p>
                    </div>
                    <div style="margin-bottom: 0.75rem;">
                        <strong style="color: var(--text-light); font-size: 0.85rem;">Role:</strong>
                        <p style="margin: 0.25rem 0 0 0; text-transform: capitalize;">${escapeHtml(user.role || 'user')}</p>
                    </div>
                    <div>
                        <strong style="color: var(--text-light); font-size: 0.85rem;">Member Since:</strong>
                        <p style="margin: 0.25rem 0 0 0;">${formatDate(user.created_at || '')}</p>
                    </div>
                </div>
                <div style="margin-top: 1.5rem;">
                    <a href="/settings" class="btn btn-primary" style="text-decoration: none; display: inline-block;">Edit Profile</a>
                </div>
            </div>
        `;
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load profile';
        contentEl.innerHTML = `<div class="empty-state"><p>${escapeHtml(errorMsg)}</p></div>`;
    }
}

// Live Search functionality (separate from advanced search)
let searchTimeout;
let liveSearchResults = null;

document.addEventListener('DOMContentLoaded', () => {
    const searchInput = document.getElementById('search-input');
    const searchContainer = document.querySelector('.search-container');
    
    if (searchInput && searchContainer) {
        // Create live search results dropdown
        const resultsDropdown = document.createElement('div');
        resultsDropdown.id = 'live-search-results';
        resultsDropdown.className = 'live-search-results';
        searchContainer.style.position = 'relative';
        searchContainer.appendChild(resultsDropdown);
        
        searchInput.addEventListener('input', (e) => {
            clearTimeout(searchTimeout);
            const query = e.target.value.trim();
            if (query.length > 2) {
                searchTimeout = setTimeout(() => performLiveSearch(query), 500);
            } else {
                hideLiveSearchResults();
            }
        });
        
        // Hide results when clicking outside
        document.addEventListener('click', (e) => {
            if (!searchContainer.contains(e.target)) {
                hideLiveSearchResults();
            }
        });
    }
});

function hideLiveSearchResults() {
    const resultsEl = document.getElementById('live-search-results');
    if (resultsEl) {
        resultsEl.style.display = 'none';
    }
}

async function performLiveSearch(query) {
    const resultsEl = document.getElementById('live-search-results');
    if (!resultsEl) return;
    
    try {
        resultsEl.style.display = 'block';
        resultsEl.innerHTML = '<div class="loading" style="padding: 1rem;">Searching...</div>';
        
        const response = await api.post('/search', { query, limit: 5 });
        const data = response.data || response;
        const results = data.results || data || [];
        
        if (results.length === 0) {
            resultsEl.innerHTML = '<div class="empty-state" style="padding: 1rem;"><p>No results found</p></div>';
            return;
        }
        
        resultsEl.innerHTML = results.map(result => {
            const type = result.type || 'unknown';
            const title = result.title || result.name || 'Untitled';
            const description = result.description || result.content || '';
            const resultId = result.id || result.ID || '';
            const userId = result.user_id || result.userId || '';
            
            return `
                <div class="live-search-item" onclick="handleLiveSearchResult('${type}', '${resultId}', '${userId}')">
                    <div style="display: flex; justify-content: space-between; align-items: start;">
                        <div style="flex: 1;">
                            <strong style="font-size: 0.9rem;">${escapeHtml(title)}</strong>
                            <span class="card-badge" style="margin-left: 0.5rem; font-size: 0.7rem;">${escapeHtml(type)}</span>
                        </div>
                    </div>
                    <p style="margin: 0.25rem 0 0 0; color: var(--text-light); font-size: 0.85rem;">${escapeHtml(description.substring(0, 60))}${description.length > 60 ? '...' : ''}</p>
                </div>
            `;
        }).join('');
        
        liveSearchResults = results;
    } catch (error) {
        console.error('Live search error:', error);
        resultsEl.innerHTML = '<div class="empty-state" style="padding: 1rem;"><p>Search failed</p></div>';
    }
}

function handleLiveSearchResult(type, id, userId) {
    hideLiveSearchResults();
    // Clear search input
    const searchInput = document.getElementById('search-input');
    if (searchInput) searchInput.value = '';
    
    if (type === 'user' && userId) {
        // Open messages with this user
        openMessages().then(() => {
            setTimeout(() => selectConversation(userId), 500);
        });
    } else {
        showMessage(`Opening ${type} ${id}`, 'info');
    }
}

// Map search functionality
let searchMap = null;
let mapMarker = null;
let mapClickLatLng = null;
let searchNearMeLocation = null; // Store user's location for "Search Near Me"

// Expose to window for cleanup
window.searchMap = () => searchMap;
window.mapMarker = () => mapMarker;
window.mapClickLatLng = () => mapClickLatLng;

function openAdvancedSearch() {
    openModal('advanced-search-modal');
    // Clear live search results
    hideLiveSearchResults();
    const searchInput = document.getElementById('search-input');
    if (searchInput) searchInput.value = '';
    
    // Reset map state
    mapClickLatLng = null;
    if (mapMarker) {
        if (searchMap) {
            searchMap.removeLayer(mapMarker);
        }
        mapMarker = null;
    }
    if (window.resultMarkersGroup && searchMap) {
        searchMap.removeLayer(window.resultMarkersGroup);
        window.resultMarkersGroup = null;
    }
    
    // Reset Search Near Me state (keep location if already set)
    const searchNearMeCheckbox = document.getElementById('search-near-me');
    if (searchNearMeCheckbox && !searchNearMeCheckbox.checked) {
        searchNearMeLocation = null;
    }
    
    // Check if map search is enabled
    const useMapSearch = document.getElementById('use-map-search');
    if (useMapSearch && useMapSearch.checked) {
        // Initialize map if not already initialized
        setTimeout(() => {
            initializeSearchMap();
        }, 100);
    }
}

function initializeSearchMap() {
    const mapEl = document.getElementById('search-map');
    if (!mapEl) return;
    
    // Check if Leaflet is loaded
    if (typeof L === 'undefined') {
        console.error('Leaflet library not loaded');
        mapEl.innerHTML = '<div style="padding: 2rem; text-align: center; color: var(--text-light);">Map library loading...</div>';
        return;
    }
    
    // Destroy existing map if it exists
    if (searchMap) {
        searchMap.remove();
        searchMap = null;
    }
    
    // Clear map container
    mapEl.innerHTML = '';
    
    // Initialize map centered on world view
    searchMap = L.map('search-map').setView([20, 0], 2);
    
    // Store reference globally for cleanup
    window._searchMapInstance = searchMap;
    
    // Add OpenStreetMap tiles
    L.tileLayer('https://{s}.tile.openstreetmap.org/{z}/{x}/{y}.png', {
        attribution: '© OpenStreetMap contributors',
        maxZoom: 19
    }).addTo(searchMap);
    
    // Add click handler to set location
    searchMap.on('click', async (e) => {
        mapClickLatLng = e.latlng;
        
        // Remove existing marker
        if (mapMarker) {
            searchMap.removeLayer(mapMarker);
        }
        
        // Add new marker
        mapMarker = L.marker([e.latlng.lat, e.latlng.lng]).addTo(searchMap);
        
        // Update coordinates display
        const coordsEl = document.getElementById('map-coordinates');
        if (coordsEl) {
            coordsEl.textContent = `Location: ${e.latlng.lat.toFixed(4)}, ${e.latlng.lng.toFixed(4)}`;
        }
        
        // Reverse geocode to get city/country
        await reverseGeocode(e.latlng.lat, e.latlng.lng);
    });
    
    // Invalidate size after a short delay to ensure proper rendering
    setTimeout(() => {
        if (searchMap) {
            searchMap.invalidateSize();
        }
    }, 200);
}

async function reverseGeocode(lat, lng) {
    try {
        // Use Nominatim (OpenStreetMap's geocoding service) - free, no API key needed
        const response = await fetch(`https://nominatim.openstreetmap.org/reverse?format=json&lat=${lat}&lon=${lng}&zoom=18&addressdetails=1`);
        const data = await response.json();
        
        if (data.address) {
            const city = data.address.city || data.address.town || data.address.village || data.address.municipality || '';
            const country = data.address.country || '';
            const location = data.address.display_name || '';
            
            // Update form fields
            if (city) document.getElementById('search-city').value = city;
            if (country) document.getElementById('search-country').value = country;
            if (location) document.getElementById('search-location').value = location;
            
            // Update coordinates display with address
            const coordsEl = document.getElementById('map-coordinates');
            if (coordsEl) {
                const addressText = location || `${city}, ${country}`.replace(/^,\s*|,\s*$/g, '');
                coordsEl.textContent = addressText || `Location: ${lat.toFixed(4)}, ${lng.toFixed(4)}`;
            }
        }
    } catch (error) {
        console.error('Reverse geocoding failed:', error);
        // Keep coordinates display even if geocoding fails
    }
}

async function getCurrentLocation() {
    if (!navigator.geolocation) {
        showMessage('Geolocation is not supported by your browser', 'error');
        return;
    }
    
    showMessage('Getting your location...', 'info');
    
    navigator.geolocation.getCurrentPosition(
        async (position) => {
            const lat = position.coords.latitude;
            const lng = position.coords.longitude;
            
            // Initialize map if not already done
            if (!searchMap) {
                initializeSearchMap();
                await new Promise(resolve => setTimeout(resolve, 200));
            }
            
            // Center map on user location
            searchMap.setView([lat, lng], 13);
            
            // Remove existing marker
            if (mapMarker) {
                searchMap.removeLayer(mapMarker);
            }
            
            // Add marker at user location
            mapMarker = L.marker([lat, lng]).addTo(searchMap);
            
            // Reverse geocode
            await reverseGeocode(lat, lng);
            
            showMessage('Location set successfully', 'success');
        },
        (error) => {
            showMessage('Failed to get your location: ' + error.message, 'error');
        }
    );
}

function toggleMapSearch() {
    const useMapSearch = document.getElementById('use-map-search').checked;
    const mapContainer = document.getElementById('map-search-container');
    
    if (mapContainer) {
        mapContainer.style.display = useMapSearch ? 'block' : 'none';
    }
    
    if (useMapSearch) {
        // Initialize map when enabled
        setTimeout(() => {
            initializeSearchMap();
        }, 200);
    } else {
        // Clean up map when disabled
        if (searchMap) {
            searchMap.remove();
            searchMap = null;
            mapMarker = null;
            mapClickLatLng = null;
        }
        if (window.resultMarkersGroup) {
            window.resultMarkersGroup = null;
        }
    }
}

function toggleSearchNearMe() {
    const searchNearMe = document.getElementById('search-near-me').checked;
    const nearMeOptions = document.getElementById('near-me-options');
    
    if (nearMeOptions) {
        nearMeOptions.style.display = searchNearMe ? 'block' : 'none';
    }
    
    // If enabled and location not set, try to get it automatically
    if (searchNearMe && !searchNearMeLocation) {
        setSearchNearMeLocation();
    }
}

async function setSearchNearMeLocation() {
    if (!navigator.geolocation) {
        showMessage('Geolocation is not supported by your browser', 'error');
        return;
    }
    
    showMessage('Getting your location...', 'info');
    
    navigator.geolocation.getCurrentPosition(
        async (position) => {
            const lat = position.coords.latitude;
            const lng = position.coords.longitude;
            
            searchNearMeLocation = { lat, lng };
            
            // Reverse geocode to fill location fields
            await reverseGeocode(lat, lng);
            
            // Also update map if it's visible
            const useMapSearch = document.getElementById('use-map-search');
            if (useMapSearch && useMapSearch.checked && searchMap) {
                searchMap.setView([lat, lng], 13);
                if (mapMarker) {
                    searchMap.removeLayer(mapMarker);
                }
                mapMarker = L.marker([lat, lng]).addTo(searchMap);
            }
            
            showMessage('Location set successfully! You can now search near you.', 'success');
        },
        (error) => {
            showMessage('Failed to get your location: ' + error.message, 'error');
        },
        {
            enableHighAccuracy: true,
            timeout: 10000,
            maximumAge: 0
        }
    );
}

function displaySearchResultsOnMap(results) {
    if (!searchMap || !results || results.length === 0) return;
    
    // Clear existing result markers (keep the location marker)
    // We'll add result markers in a group
    if (window.resultMarkersGroup) {
        searchMap.removeLayer(window.resultMarkersGroup);
    }
    
    window.resultMarkersGroup = L.layerGroup().addTo(searchMap);
    
    // Add markers for results that have location data
    results.forEach((result, index) => {
        // Try to extract location from result
        // This depends on your data structure - adjust as needed
        const lat = result.latitude || result.lat;
        const lng = result.longitude || result.lng || result.lon;
        
        if (lat && lng) {
            const title = result.title || result.name || 'Result';
            const marker = L.marker([lat, lng]).addTo(window.resultMarkersGroup);
            marker.bindPopup(`<strong>${escapeHtml(title)}</strong><br>${escapeHtml(result.description || '')}`);
        }
    });
    
    // Fit map to show all markers if we have any
    if (window.resultMarkersGroup.getLayers().length > 0) {
        const bounds = window.resultMarkersGroup.getBounds();
        searchMap.fitBounds(bounds, { padding: [50, 50] });
    }
}

async function performAdvancedSearch(e) {
    e.preventDefault();
    const query = document.getElementById('advanced-query').value;
    const type = document.getElementById('search-type').value;
    const location = document.getElementById('search-location').value;
    const country = document.getElementById('search-country').value;
    const city = document.getElementById('search-city').value;
    const resultsEl = document.getElementById('search-results');
    
    // Check if "Search Near Me" is enabled
    const searchNearMe = document.getElementById('search-near-me').checked;
    const useMapSearch = document.getElementById('use-map-search').checked;
    
    // If "Search Near Me" is enabled but location not set, try to get it
    if (searchNearMe && !searchNearMeLocation) {
        await setSearchNearMeLocation();
        if (!searchNearMeLocation) {
            showMessage('Please allow location access to search near you', 'error');
            return;
        }
    }
    
    // If map search is enabled and we have coordinates, use them
    if (useMapSearch && mapClickLatLng) {
        // Use map coordinates for search
        // The location fields should already be filled by reverse geocoding
    }
    
    // Validate search parameters
    if (!query.trim() && !location && !country && !city && !mapClickLatLng && !searchNearMeLocation) {
        showMessage('Please enter a search query, location, or enable "Search Near Me"', 'error');
        return;
    }
    
    try {
        resultsEl.innerHTML = '<div class="loading">Searching...</div>';
        const searchParams = {
            query: query || '',
            type: type || 'all',
            limit: 50
        };
        
        // Add location filters if provided
        if (location) searchParams.location = location;
        if (country) searchParams.country = country;
        if (city) searchParams.city = city;
        
        // Priority: Search Near Me > Map Click > Manual location
        if (searchNearMe && searchNearMeLocation) {
            // Use "Search Near Me" location with radius
            searchParams.latitude = searchNearMeLocation.lat;
            searchParams.longitude = searchNearMeLocation.lng;
            const radius = document.getElementById('search-radius')?.value || '5';
            searchParams.radius = parseFloat(radius); // Radius in kilometers
        } else if (mapClickLatLng) {
            // Use map click coordinates
            searchParams.latitude = mapClickLatLng.lat;
            searchParams.longitude = mapClickLatLng.lng;
        }
        
        const response = await api.post('/search', searchParams);
        
        const data = response.data || response;
        const results = data.results || data || [];
        renderSearchResults(results);
        
        // Display results on map if map is available or if Search Near Me is enabled
        if ((useMapSearch && searchMap) || (searchNearMe && searchNearMeLocation)) {
            // Initialize map if needed for Search Near Me
            if (searchNearMe && !searchMap) {
                const mapCheckbox = document.getElementById('use-map-search');
                if (mapCheckbox) {
                    mapCheckbox.checked = true;
                    toggleMapSearch();
                    await new Promise(resolve => setTimeout(resolve, 300));
                }
            }
            if (searchMap) {
                displaySearchResultsOnMap(results);
                // Center map on user location if Search Near Me
                if (searchNearMe && searchNearMeLocation) {
                    searchMap.setView([searchNearMeLocation.lat, searchNearMeLocation.lng], 12);
                }
            }
        }
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Search failed';
        resultsEl.innerHTML = `<div class="empty-state"><p>${escapeHtml(errorMsg)}</p></div>`;
    }
}

function renderSearchResults(results) {
    const resultsEl = document.getElementById('search-results');
    
    if (!results || results.length === 0) {
        resultsEl.innerHTML = '<div class="empty-state"><p>No results found</p></div>';
        return;
    }
    
    resultsEl.innerHTML = results.map(result => {
        const type = result.type || 'unknown';
        const title = result.title || result.name || 'Untitled';
        const description = result.description || result.content || '';
        const resultId = result.id || result.ID || '';
        const userId = result.user_id || result.userId || '';
        
        // Show message button for user results
        const messageButton = type === 'user' && userId ? 
            `<button class="btn btn-primary btn-sm" onclick="messageUserFromSearch('${userId}', '${escapeHtml(title)}')" style="margin-top: 0.5rem;">💬 Message</button>` : '';
        
        return `
            <div class="search-result-item">
                <div style="display: flex; justify-content: space-between; align-items: start;">
                    <div style="flex: 1;">
                        <strong>${escapeHtml(title)}</strong>
                        <span class="card-badge" style="margin-left: 0.5rem;">${escapeHtml(type)}</span>
                    </div>
                </div>
                <p class="card-description">${escapeHtml(description.substring(0, 100))}${description.length > 100 ? '...' : ''}</p>
                ${messageButton}
            </div>
        `;
    }).join('');
}

function handleSearchResult(type, id) {
    showMessage(`Opening ${type} ${id}`, 'info');
    // Implement navigation based on type
}

async function messageUserFromSearch(userId, userName) {
    // Close search modal
    closeModal('advanced-search-modal');
    
    // Open messages modal
    await openMessages();
    
    // Select conversation with this user
    setTimeout(async () => {
        await selectConversation(userId);
    }, 500);
}

// Notifications
async function loadNotificationCount() {
    // Check if authenticated
    const token = localStorage.getItem('access_token');
    if (!token) return;
    
    try {
        const response = await api.get('/notifications/unread-count');
        const count = response.data?.count || response.count || 0;
        unreadNotifications = count;
        updateNotificationBadge(count);
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : '';
        // Don't log auth errors - redirect will happen
        if (!errorMsg.includes('UNAUTHORIZED') && !errorMsg.includes('401')) {
            console.error('Failed to load notification count:', error);
        }
    }
}

function updateNotificationBadge(count) {
    const badge = document.getElementById('notifications-badge');
    if (badge) {
        if (count > 0) {
            badge.textContent = count > 99 ? '99+' : count;
            badge.style.display = 'block';
        } else {
            badge.style.display = 'none';
        }
    }
}

async function openNotifications() {
    openModal('notifications-modal');
    await loadNotifications();
}

async function loadNotifications() {
    const listEl = document.getElementById('notifications-list');
    if (!listEl) return;
    
    try {
        listEl.innerHTML = '<div class="loading">Loading...</div>';
        const response = await api.get('/notifications');
        const notifications = response.data || response || [];
        renderNotifications(notifications);
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load notifications';
        listEl.innerHTML = `<div class="empty-state"><p>${escapeHtml(errorMsg)}</p></div>`;
    }
}

function renderNotifications(notifications) {
    const listEl = document.getElementById('notifications-list');
    
    if (!notifications || notifications.length === 0) {
        listEl.innerHTML = '<div class="empty-state"><p>No notifications</p></div>';
        return;
    }
    
    listEl.innerHTML = notifications.map(notif => {
        const isRead = notif.read || notif.is_read;
        const title = notif.title || 'Notification';
        const message = notif.message || notif.content || '';
        const createdAt = notif.created_at || notif.createdAt || '';
        
        return `
            <div class="notification-item ${isRead ? 'read' : 'unread'}" onclick="handleNotification('${notif.id || ''}')">
                <div style="display: flex; justify-content: space-between;">
                    <div>
                        <strong>${escapeHtml(title)}</strong>
                        ${!isRead ? '<span class="badge" style="margin-left: 0.5rem;">New</span>' : ''}
                    </div>
                    <small style="color: var(--text-light);">${formatDate(createdAt)}</small>
                </div>
                <p style="margin: 0.5rem 0 0 0; color: var(--text-light);">${escapeHtml(message)}</p>
            </div>
        `;
    }).join('');
}

async function handleNotification(id) {
    try {
        await api.post('/notifications/read', { notification_id: id });
        await loadNotifications();
        await loadNotificationCount();
    } catch (error) {
        console.error('Failed to mark notification as read:', error);
    }
}

async function markAllNotificationsRead() {
    try {
        await api.post('/notifications/read-all');
        await loadNotifications();
        await loadNotificationCount();
        showMessage('All notifications marked as read', 'success');
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to mark all as read';
        showMessage(errorMsg, 'error');
    }
}

// Messages
async function loadMessageCount() {
    // Check if authenticated
    const token = localStorage.getItem('access_token');
    if (!token) return;
    
    try {
        const response = await api.get('/messages/unread-count');
        const count = response.data?.count || response.count || 0;
        unreadMessages = count;
        updateMessageBadge(count);
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : '';
        // Don't log auth errors - redirect will happen
        if (!errorMsg.includes('UNAUTHORIZED') && !errorMsg.includes('401')) {
            console.error('Failed to load message count:', error);
        }
    }
}

function updateMessageBadge(count) {
    const badge = document.getElementById('messages-badge');
    if (badge) {
        if (count > 0) {
            badge.textContent = count > 99 ? '99+' : count;
            badge.style.display = 'block';
        } else {
            badge.style.display = 'none';
        }
    }
}

let currentConversationId = null;
let messagePollInterval = null;
let notificationPollInterval = null;

async function openMessages() {
    openModal('messages-modal');
    await loadConversations();
    
    // Start polling for new messages if enabled in settings
    startMessagePolling();
}

// User search for messaging
let userSearchTimeout;
async function searchUsersForMessage(event) {
    const query = event.target.value.trim();
    const resultsEl = document.getElementById('message-user-search-results');
    
    if (!resultsEl) return;
    
    clearTimeout(userSearchTimeout);
    
    if (query.length < 2) {
        resultsEl.style.display = 'none';
        return;
    }
    
    userSearchTimeout = setTimeout(async () => {
        try {
            resultsEl.style.display = 'block';
            resultsEl.innerHTML = '<div class="loading" style="padding: 0.5rem;">Searching...</div>';
            
            const response = await api.post('/search', {
                query: query,
                type: 'user',
                limit: 10
            });
            
            // Handle nested response structure from backend
            // Backend returns: { success: true, data: { type: "search_results", data: { results: [...] } } }
            let results = [];
            
            if (response.data) {
                // Check if data.data.results exists (nested structure)
                if (response.data.data && response.data.data.results && Array.isArray(response.data.data.results)) {
                    results = response.data.data.results;
                } 
                // Check if data.results exists
                else if (response.data.results && Array.isArray(response.data.results)) {
                    results = response.data.results;
                }
                // Check if data is directly an array
                else if (Array.isArray(response.data)) {
                    results = response.data;
                }
                // Check if data.data is an array
                else if (response.data.data && Array.isArray(response.data.data)) {
                    results = response.data.data;
                }
            } 
            // Check if response is directly an array
            else if (Array.isArray(response)) {
                results = response;
            }
            // Check if response.results exists
            else if (response.results && Array.isArray(response.results)) {
                results = response.results;
            }
            
            // Filter to only user type results and extract user data
            results = results
                .filter(item => {
                    const type = item.type || '';
                    return type === 'user' || type === 'User';
                })
                .map(item => {
                    // Extract user data from Data field if it exists
                    if (item.data && typeof item.data === 'object') {
                        return {
                            ...item.data,
                            id: item.id || item.ID || item.data.id || item.data.ID,
                            type: 'user'
                        };
                    }
                    return item;
                })
                .filter(item => {
                    const id = item.id || item.ID || '';
                    return id; // Only include items with valid IDs
                });
            
            if (results.length === 0) {
                resultsEl.innerHTML = '<div class="empty-state" style="padding: 0.5rem;"><p>No users found</p></div>';
                return;
            }
            
            resultsEl.innerHTML = results.map(user => {
                const userId = user.id || user.ID || '';
                const userName = user.name || user.email || 'Unknown';
                const userEmail = user.email || '';
                const userRole = user.role || 'user';
                
                return `
                    <div class="message-search-result-item" onclick="startConversationWithUser('${userId}', '${escapeHtml(userName)}')">
                        <div style="display: flex; align-items: center; gap: 0.5rem;">
                            <div style="width: 32px; height: 32px; border-radius: 50%; background: var(--primary); color: var(--white); display: flex; align-items: center; justify-content: center; font-weight: bold;">
                                ${(userName || 'U').charAt(0).toUpperCase()}
                            </div>
                            <div style="flex: 1;">
                                <strong style="font-size: 0.9rem;">${escapeHtml(userName)}</strong>
                                <p style="margin: 0.25rem 0 0 0; color: var(--text-light); font-size: 0.8rem;">${escapeHtml(userEmail)}</p>
                            </div>
                            <span class="card-badge" style="font-size: 0.7rem;">${escapeHtml(userRole)}</span>
                        </div>
                    </div>
                `;
            }).join('');
        } catch (error) {
            console.error('User search error:', error);
            const errorMsg = error instanceof Error ? error.message : 'Search failed';
            
            // Don't show error if it's an auth error (redirect will happen)
            if (errorMsg.includes('UNAUTHORIZED') || errorMsg.includes('401')) {
                resultsEl.style.display = 'none';
                return;
            }
            
            resultsEl.innerHTML = `<div class="empty-state" style="padding: 0.5rem;"><p>${escapeHtml(errorMsg)}</p></div>`;
        }
    }, 300);
}

// Hide search results when clicking outside
document.addEventListener('click', (e) => {
    const searchContainer = document.querySelector('.messages-search-container');
    const resultsEl = document.getElementById('message-user-search-results');
    if (searchContainer && resultsEl && !searchContainer.contains(e.target)) {
        resultsEl.style.display = 'none';
    }
});

async function startConversationWithUser(userId, userName) {
    // Hide search results
    const resultsEl = document.getElementById('message-user-search-results');
    if (resultsEl) {
        resultsEl.style.display = 'none';
    }
    
    // Clear search input
    const searchInput = document.getElementById('message-user-search');
    if (searchInput) searchInput.value = '';
    
    // Select or create conversation with this user
    await selectConversation(userId);
}

async function loadConversations() {
    const listEl = document.getElementById('conversations-list');
    if (!listEl) return;
    
    try {
        listEl.innerHTML = '<div class="loading">Loading...</div>';
        const response = await api.get('/messages/conversations');
        
        // Handle different response structures
        let conversations = [];
        if (Array.isArray(response)) {
            conversations = response;
        } else if (response.data) {
            conversations = Array.isArray(response.data) ? response.data : [];
        } else if (response.conversations) {
            conversations = Array.isArray(response.conversations) ? response.conversations : [];
        }
        
        renderConversations(conversations);
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load conversations';
        // Don't show error if it's an auth error (redirect will happen)
        if (errorMsg.includes('UNAUTHORIZED') || errorMsg.includes('401')) {
            return; // Will redirect to login
        }
        listEl.innerHTML = `<div class="empty-state"><p>${escapeHtml(errorMsg)}</p></div>`;
    }
}

function renderConversations(conversations) {
    const listEl = document.getElementById('conversations-list');
    if (!listEl) return;
    
    // Ensure conversations is an array
    if (!Array.isArray(conversations)) {
        console.error('Conversations is not an array:', conversations);
        listEl.innerHTML = '<div class="empty-state"><p>No conversations</p></div>';
        return;
    }
    
    if (conversations.length === 0) {
        listEl.innerHTML = '<div class="empty-state"><p>No conversations</p></div>';
        return;
    }
    
    listEl.innerHTML = conversations.map(conv => {
        const name = conv.other_user_name || conv.name || 'Unknown';
        const lastMessage = conv.last_message || '';
        const unread = conv.unread_count || 0;
        
        const userId = conv.other_user_id || conv.id;
        return `
            <div class="conversation-item" onclick="selectConversation('${userId}', this)">
                <div style="display: flex; justify-content: space-between; align-items: start;">
                    <div>
                        <strong>${escapeHtml(name)}</strong>
                        ${unread > 0 ? `<span class="badge">${unread}</span>` : ''}
                    </div>
                </div>
                <p style="margin: 0.5rem 0 0 0; color: var(--text-light); font-size: 0.85rem;">${escapeHtml(lastMessage.substring(0, 50))}${lastMessage.length > 50 ? '...' : ''}</p>
            </div>
        `;
    }).join('');
}

async function selectConversation(userId, element) {
    currentConversationId = userId;
    
    // Update active conversation
    document.querySelectorAll('.conversation-item').forEach(el => el.classList.remove('active'));
    if (element) {
        element.classList.add('active');
    } else if (event && event.currentTarget) {
        event.currentTarget.classList.add('active');
    }
    
    await loadMessages(userId);
}

async function loadMessages(userId) {
    const viewEl = document.getElementById('messages-view');
    if (!viewEl) return;
    
    try {
        viewEl.innerHTML = '<div class="loading">Loading messages...</div>';
        const response = await api.get(`/messages?user_id=${userId}`);
        const messages = response.data || response || [];
        renderMessages(messages, userId);
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to load messages';
        viewEl.innerHTML = `<div class="empty-state"><p>${escapeHtml(errorMsg)}</p></div>`;
    }
}

function renderMessages(messages, userId) {
    const viewEl = document.getElementById('messages-view');
    
    if (!messages || messages.length === 0) {
        viewEl.innerHTML = '<div class="empty-state"><p>No messages yet</p></div>';
        return;
    }
    
    const currentUserId = currentUser?.id || currentUser?.ID;
    
    const messagesHTML = messages.map(msg => {
        const isSent = msg.sender_id === currentUserId || msg.sender_id?.toString() === currentUserId?.toString();
        const content = msg.content || msg.message || '';
        const createdAt = msg.created_at || msg.createdAt || '';
        
        return `
            <div class="message-item ${isSent ? 'sent' : 'received'}">
                <p style="margin: 0;">${escapeHtml(content)}</p>
                <small style="opacity: 0.7; font-size: 0.75rem;">${formatDate(createdAt)}</small>
            </div>
        `;
    }).join('');
    
    viewEl.innerHTML = `
        <div class="messages-list" id="messages-list">
            ${messagesHTML}
        </div>
        <div class="message-input-area">
            <input type="text" id="message-input" placeholder="Type a message..." onkeypress="handleMessageKeyPress(event, '${userId}')">
            <button class="btn btn-primary" onclick="sendMessage('${userId}')">Send</button>
        </div>
    `;
    
    // Scroll to bottom
    const messagesList = document.getElementById('messages-list');
    if (messagesList) {
        messagesList.scrollTop = messagesList.scrollHeight;
    }
}

function handleMessageKeyPress(e, userId) {
    if (e.key === 'Enter') {
        sendMessage(userId);
    }
}

async function sendMessage(userId) {
    const input = document.getElementById('message-input');
    const content = input?.value.trim();
    
    if (!content) return;
    
    try {
        await api.post('/messages', {
            recipient_id: userId,
            content: content
        });
        
        input.value = '';
        await loadMessages(userId);
        await loadMessageCount();
        await loadConversations();
        
        // Refresh message count for recipient (they'll see it on next poll)
    } catch (error) {
        const errorMsg = error instanceof Error ? error.message : 'Failed to send message';
        showMessage(errorMsg, 'error');
    }
}

// Real-time polling for messages and notifications
function startMessagePolling() {
    // Check if messaging notifications are enabled
    const messagingEnabled = localStorage.getItem('messaging_enabled') !== 'false';
    
    if (!messagingEnabled) return;
    
    // Clear existing interval
    if (messagePollInterval) {
        clearInterval(messagePollInterval);
    }
    
    // Poll every 10 seconds for new messages
    messagePollInterval = setInterval(async () => {
        try {
            await loadMessageCount();
            // Reload conversations if messages modal is open
            const messagesModal = document.getElementById('messages-modal');
            if (messagesModal && messagesModal.classList.contains('active')) {
                await loadConversations();
                // Reload current conversation messages if one is selected
                if (currentConversationId) {
                    await loadMessages(currentConversationId);
                }
            }
        } catch (error) {
            console.error('Message polling error:', error);
        }
    }, 10000); // Poll every 10 seconds
}

function startNotificationPolling() {
    // Check if notifications are enabled
    const notificationsEnabled = localStorage.getItem('notifications_enabled') !== 'false';
    
    if (!notificationsEnabled) return;
    
    // Clear existing interval
    if (notificationPollInterval) {
        clearInterval(notificationPollInterval);
    }
    
    // Poll every 15 seconds for new notifications
    notificationPollInterval = setInterval(async () => {
        try {
            await loadNotificationCount();
        } catch (error) {
            console.error('Notification polling error:', error);
        }
    }, 15000); // Poll every 15 seconds
}

function stopMessagePolling() {
    if (messagePollInterval) {
        clearInterval(messagePollInterval);
        messagePollInterval = null;
    }
}

function stopNotificationPolling() {
    if (notificationPollInterval) {
        clearInterval(notificationPollInterval);
        notificationPollInterval = null;
    }
}

// Utility functions
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
        const now = new Date();
        const diff = now - date;
        const minutes = Math.floor(diff / 60000);
        const hours = Math.floor(diff / 3600000);
        const days = Math.floor(diff / 86400000);
        
        if (minutes < 1) return 'Just now';
        if (minutes < 60) return `${minutes}m ago`;
        if (hours < 24) return `${hours}h ago`;
        if (days < 7) return `${days}d ago`;
        
        return date.toLocaleDateString();
    } catch (e) {
        return dateString;
    }
}

// Export functions
window.toggleAvatarMenu = toggleAvatarMenu;
window.viewProfile = viewProfile;
window.openProfileModal = openProfileModal;
window.openAdvancedSearch = openAdvancedSearch;
window.performAdvancedSearch = performAdvancedSearch;
window.openNotifications = openNotifications;
window.markAllNotificationsRead = markAllNotificationsRead;
window.openMessages = openMessages;
window.selectConversation = selectConversation;
window.sendMessage = sendMessage;
window.handleMessageKeyPress = handleMessageKeyPress;
window.messageUserFromSearch = messageUserFromSearch;
window.handleLiveSearchResult = handleLiveSearchResult;
window.initializeNavbar = initializeNavbar;
window.updateAvatar = updateAvatar;
window.toggleMapSearch = toggleMapSearch;
window.getCurrentLocation = getCurrentLocation;
window.toggleSearchNearMe = toggleSearchNearMe;
window.setSearchNearMeLocation = setSearchNearMeLocation;
window.searchUsersForMessage = searchUsersForMessage;
window.startConversationWithUser = startConversationWithUser;
window.startMessagePolling = startMessagePolling;
window.stopMessagePolling = stopMessagePolling;
window.startNotificationPolling = startNotificationPolling;
window.stopNotificationPolling = stopNotificationPolling;


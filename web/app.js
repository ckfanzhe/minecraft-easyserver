// API base URL
const API_BASE = '/api';

// DOM elements
const elements = {
    serverStatus: document.getElementById('server-status'),
    startBtn: document.getElementById('start-btn'),
    stopBtn: document.getElementById('stop-btn'),
    restartBtn: document.getElementById('restart-btn'),
    refreshBtn: document.getElementById('refresh-btn'),
    configForm: document.getElementById('config-form'),
    newPlayerInput: document.getElementById('new-player'),
    addPlayerBtn: document.getElementById('add-player-btn'),
    allowlistContainer: document.getElementById('allowlist-container'),
    permissionPlayer: document.getElementById('permission-player'),
    addPermissionBtn: document.getElementById('add-permission-btn'),
    permissionsContainer: document.getElementById('permissions-container'),
    uploadBtn: document.getElementById('upload-btn'),
    worldUpload: document.getElementById('world-upload'),
    worldsContainer: document.getElementById('worlds-container'),
    toast: document.getElementById('toast'),
    toastMessage: document.getElementById('toast-message'),
    permissionModal: document.getElementById('permission-modal'),
    modalPlayerName: document.getElementById('modal-player-name'),
    closeModalBtn: document.getElementById('close-modal-btn'),
    cancelModalBtn: document.getElementById('cancel-modal-btn')
};

// Initialize application
document.addEventListener('DOMContentLoaded', function() {
    // Initialize i18n first
    if (window.i18n) {
        window.i18n.init();
    }
    
    initializeApp();
    bindEvents();
});

// Initialize application data
async function initializeApp() {
    await loadServerStatus();
    await loadServerConfig();
    await loadAllowlist();
    await loadPermissions();
    await loadWorlds();
}

// Bind event listeners
function bindEvents() {
    // Server control buttons
    if (elements.startBtn) elements.startBtn.addEventListener('click', () => controlServer('start'));
    if (elements.stopBtn) elements.stopBtn.addEventListener('click', () => controlServer('stop'));
    if (elements.restartBtn) elements.restartBtn.addEventListener('click', () => controlServer('restart'));
    if (elements.refreshBtn) elements.refreshBtn.addEventListener('click', initializeApp);

    // Configuration form
    if (elements.configForm) elements.configForm.addEventListener('submit', saveServerConfig);

    // Allowlist management
    if (elements.addPlayerBtn) elements.addPlayerBtn.addEventListener('click', addToAllowlist);
    if (elements.newPlayerInput) {
        elements.newPlayerInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                addToAllowlist();
            }
        });
    }

    // Permission management
    if (elements.addPermissionBtn) elements.addPermissionBtn.addEventListener('click', showPermissionModal);
    if (elements.permissionPlayer) {
        elements.permissionPlayer.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                showPermissionModal();
            }
        });
    }

    // Modal events - add existence check
    if (elements.closeModalBtn) {
        elements.closeModalBtn.addEventListener('click', hidePermissionModal);
    }
    if (elements.cancelModalBtn) {
        elements.cancelModalBtn.addEventListener('click', hidePermissionModal);
    }
    if (elements.permissionModal) {
        elements.permissionModal.addEventListener('click', function(e) {
            if (e.target === elements.permissionModal) {
                hidePermissionModal();
            }
        });
    }

    // Permission option click events
    document.addEventListener('click', function(e) {
        if (e.target.closest('.permission-option')) {
            const level = e.target.closest('.permission-option').dataset.level;
            setPlayerPermission(level);
        }
    });

    // World upload
    if (elements.uploadBtn) elements.uploadBtn.addEventListener('click', () => elements.worldUpload.click());
    if (elements.worldUpload) elements.worldUpload.addEventListener('change', uploadWorld);
}

// API request wrapper
async function apiRequest(endpoint, options = {}) {
    try {
        const response = await fetch(`${API_BASE}${endpoint}`, {
            headers: {
                'Content-Type': 'application/json',
                ...options.headers
            },
            ...options
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        return await response.json();
    } catch (error) {
        console.error('API request failed:', error);
        showToast(window.i18n ? window.i18n.t('message.request-failed') : 'Request failed: ' + error.message, 'error');
        throw error;
    }
}

// Show toast message
function showToast(message, type = 'success') {
    elements.toastMessage.textContent = message;
    elements.toast.className = `fixed top-4 right-4 px-6 py-3 rounded-lg shadow-lg transform transition-transform duration-300 ${
        type === 'error' ? 'bg-red-500' : 'bg-green-500'
    } text-white`;
    
    // Show toast
    elements.toast.style.transform = 'translateX(0)';
    
    // Hide after 3 seconds
    setTimeout(() => {
        elements.toast.style.transform = 'translateX(100%)';
    }, 3000);
}

// Load server status
async function loadServerStatus() {
    try {
        const data = await apiRequest('/status');
        updateServerStatus(data.status);
    } catch (error) {
        updateServerStatus('unknown');
    }
}

// Update server status display
function updateServerStatus(status) {
    const statusElement = elements.serverStatus;
    statusElement.className = 'px-3 py-1 rounded-full text-sm';
    
    switch (status) {
        case 'running':
            statusElement.textContent = window.i18n ? window.i18n.t('nav.status.running') : 'Running';
            statusElement.classList.add('bg-green-500');
            break;
        case 'stopped':
            statusElement.textContent = window.i18n ? window.i18n.t('nav.status.stopped') : 'Stopped';
            statusElement.classList.add('bg-red-500');
            break;
        default:
            statusElement.textContent = window.i18n ? window.i18n.t('nav.status.unknown') : 'Unknown';
            statusElement.classList.add('bg-gray-500');
    }
}

// Server control
async function controlServer(action) {
    try {
        const data = await apiRequest(`/${action}`, { method: 'POST' });
        showToast(data.message);
        
        // Delay status refresh
        setTimeout(loadServerStatus, 2000);
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Load server configuration
async function loadServerConfig() {
    try {
        const data = await apiRequest('/config');
        if (data.config) {
            populateConfigForm(data.config);
        }
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Populate configuration form
function populateConfigForm(config) {
    document.getElementById('server-name').value = config['server-name'] || '';
    document.getElementById('gamemode').value = config.gamemode || 'survival';
    document.getElementById('difficulty').value = config.difficulty || 'easy';
    document.getElementById('max-players').value = config['max-players'] || 10;
    document.getElementById('server-port').value = config['server-port'] || 19132;
    document.getElementById('allow-cheats').checked = config['allow-cheats'] === 'true';
    document.getElementById('allow-list').checked = config['allow-list'] === 'true';
}

// Save server configuration
async function saveServerConfig(e) {
    e.preventDefault();
    
    const config = {
        'server-name': document.getElementById('server-name').value,
        'gamemode': document.getElementById('gamemode').value,
        'difficulty': document.getElementById('difficulty').value,
        'max-players': parseInt(document.getElementById('max-players').value),
        'server-port': parseInt(document.getElementById('server-port').value),
        'allow-cheats': document.getElementById('allow-cheats').checked,
        'allow-list': document.getElementById('allow-list').checked
    };
    
    try {
        const data = await apiRequest('/config', {
            method: 'PUT',
            body: JSON.stringify({ config })
        });
        showToast(data.message);
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Load allowlist
async function loadAllowlist() {
    try {
        const data = await apiRequest('/allowlist');
        renderAllowlist(data.allowlist || []);
    } catch (error) {
        renderAllowlist([]);
    }
}

// Render allowlist
function renderAllowlist(allowlist) {
    elements.allowlistContainer.innerHTML = '';
    
    if (allowlist.length === 0) {
        const emptyMessage = window.i18n ? 
            window.i18n.t('allowlist.empty') : 
            'No allowlist users';
        elements.allowlistContainer.innerHTML = `<p class="text-gray-500 text-center py-4">${emptyMessage}</p>`;
        return;
    }
    
    allowlist.forEach(player => {
        const playerElement = createPlayerElement(player, 'allowlist');
        elements.allowlistContainer.appendChild(playerElement);
    });
}

// Add to allowlist
async function addToAllowlist() {
    const playerName = elements.newPlayerInput.value.trim();
    if (!playerName) {
        const errorMessage = window.i18n ? 
            window.i18n.t('allowlist.error.empty-name') : 
            'Please enter player name';
        showToast(errorMessage, 'error');
        return;
    }
    
    try {
        const data = await apiRequest('/allowlist', {
            method: 'POST',
            body: JSON.stringify({ name: playerName })
        });
        showToast(data.message);
        elements.newPlayerInput.value = '';
        await loadAllowlist();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Remove from allowlist
async function removeFromAllowlist(playerName) {
    try {
        const data = await apiRequest(`/allowlist/${encodeURIComponent(playerName)}`, {
            method: 'DELETE'
        });
        showToast(data.message);
        await loadAllowlist();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Load permissions
async function loadPermissions() {
    try {
        const data = await apiRequest('/permissions');
        renderPermissions(data.permissions || []);
    } catch (error) {
        renderPermissions([]);
    }
}

// Render permissions
function renderPermissions(permissions) {
    elements.permissionsContainer.innerHTML = '';
    
    if (permissions.length === 0) {
        const emptyMessage = window.i18n ? 
            window.i18n.t('permission.empty') : 
            'No permission settings';
        elements.permissionsContainer.innerHTML = `<p class="text-gray-500 text-center py-4">${emptyMessage}</p>`;
        return;
    }
    
    permissions.forEach(permission => {
        const permissionElement = createPermissionElement(permission);
        elements.permissionsContainer.appendChild(permissionElement);
    });
}

// Show permission selection modal
function showPermissionModal() {
    const playerName = elements.permissionPlayer.value.trim();
    
    if (!playerName) {
        const errorMessage = window.i18n ? 
            window.i18n.t('permission.error.empty-name') : 
            'Please enter player name';
        showToast(errorMessage, 'error');
        return;
    }
    
    elements.modalPlayerName.textContent = playerName;
    elements.permissionModal.classList.remove('hidden');
}

// Hide permission selection modal
function hidePermissionModal() {
    elements.permissionModal.classList.add('hidden');
}

// Set player permission
async function setPlayerPermission(level) {
    const playerName = elements.permissionPlayer.value.trim();
    
    if (!playerName) {
        const errorMessage = window.i18n ? 
            window.i18n.t('permission.error.empty-name') : 
            'Please enter player name';
        showToast(errorMessage, 'error');
        return;
    }
    
    try {
        const data = await apiRequest('/permissions', {
            method: 'PUT',
            body: JSON.stringify({ name: playerName, level })
        });
        showToast(data.message);
        elements.permissionPlayer.value = '';
        hidePermissionModal();
        await loadPermissions();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Load worlds list
async function loadWorlds() {
    try {
        const data = await apiRequest('/worlds');
        renderWorlds(data.worlds || []);
    } catch (error) {
        renderWorlds([]);
    }
}

// Render worlds list
function renderWorlds(worlds) {
    elements.worldsContainer.innerHTML = '';
    
    if (worlds.length === 0) {
        const emptyMessage = window.i18n ? 
            window.i18n.t('world.empty') : 
            'No world files';
        elements.worldsContainer.innerHTML = `<p class="text-gray-500 text-center py-4">${emptyMessage}</p>`;
        return;
    }
    
    worlds.forEach(world => {
        const worldElement = createWorldElement(world);
        elements.worldsContainer.appendChild(worldElement);
    });
}

// Upload world
async function uploadWorld() {
    const file = elements.worldUpload.files[0];
    if (!file) return;
    
    const formData = new FormData();
    formData.append('world', file);
    
    try {
        const response = await fetch(`${API_BASE}/worlds/upload`, {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        showToast(data.message);
        elements.worldUpload.value = '';
        await loadWorlds();
    } catch (error) {
        const errorMessage = window.i18n ? 
            window.i18n.t('world.upload.error') : 
            'Upload failed: ';
        showToast(errorMessage + error.message, 'error');
    }
}

// Create player element
function createPlayerElement(playerName, type) {
    const div = document.createElement('div');
    div.className = 'flex items-center justify-between bg-gray-50 px-3 py-2 rounded';
    
    // Escape special characters in player name
    const escapedName = playerName.replace(/'/g, "\\'").replace(/"/g, '\\"');
    
    div.innerHTML = `
        <span class="font-medium">${playerName}</span>
        <button onclick="${type === 'allowlist' ? 'removeFromAllowlist' : 'removePermission'}('${escapedName}')" 
                class="text-red-500 hover:text-red-700 transition duration-200">
            <i class="fas fa-trash"></i>
        </button>
    `;
    
    return div;
}

// Create permission element
function createPermissionElement(permission) {
    const div = document.createElement('div');
    div.className = 'flex items-center justify-between bg-gray-50 px-3 py-2 rounded';
    
    const levelText = {
        'visitor': window.i18n ? window.i18n.t('permission.level.visitor') : 'Visitor',
        'member': window.i18n ? window.i18n.t('permission.level.member') : 'Member',
        'operator': window.i18n ? window.i18n.t('permission.level.operator') : 'Operator'
    };
    
    const levelColor = {
        'visitor': 'text-gray-600',
        'member': 'text-blue-600',
        'operator': 'text-red-600'
    };
    
    // Escape special characters in permission name
    const escapedName = permission.name.replace(/'/g, "\\'").replace(/"/g, '\\"');
    
    div.innerHTML = `
        <div>
            <span class="font-medium">${permission.name}</span>
            <span class="ml-2 px-2 py-1 text-xs rounded ${levelColor[permission.level]} bg-gray-200">
                ${levelText[permission.level]}
            </span>
        </div>
        <button onclick="removePermission('${escapedName}')" 
                class="text-red-500 hover:text-red-700 transition duration-200">
            <i class="fas fa-trash"></i>
        </button>
    `;
    
    return div;
}

// Create world element
function createWorldElement(world) {
    const div = document.createElement('div');
    div.className = 'flex items-center justify-between bg-gray-50 px-3 py-2 rounded';
    
    // Escape special characters in world name
    const escapedName = world.name.replace(/'/g, "\\'").replace(/"/g, '\\"');
    
    const currentWorldText = window.i18n ? window.i18n.t('world.current') : 'Current World';
    
    div.innerHTML = `
        <div>
            <span class="font-medium">${world.name}</span>
            ${world.active ? `<span class="ml-2 px-2 py-1 text-xs rounded bg-green-200 text-green-800">${currentWorldText}</span>` : ''}
        </div>
        <div class="space-x-2">
            ${!world.active ? `<button onclick="activateWorld('${escapedName}')" 
                class="text-blue-500 hover:text-blue-700 transition duration-200">
                <i class="fas fa-play"></i>
            </button>` : ''}
            <button onclick="deleteWorld('${escapedName}')" 
                    class="text-red-500 hover:text-red-700 transition duration-200">
                <i class="fas fa-trash"></i>
            </button>
        </div>
    `;
    
    return div;
}

// Delete world
async function deleteWorld(worldName) {
    const confirmMessage = window.i18n ? 
        window.i18n.t('world.deleteConfirm', { worldName }) : 
        `Are you sure you want to delete world "${worldName}"? This action cannot be undone!`;
    
    if (!confirm(confirmMessage)) {
        return;
    }
    
    try {
        const data = await apiRequest(`/worlds/${encodeURIComponent(worldName)}`, {
            method: 'DELETE'
        });
        showToast(data.message);
        await loadWorlds();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Activate world
async function activateWorld(worldName) {
    try {
        const data = await apiRequest(`/worlds/${encodeURIComponent(worldName)}/activate`, {
            method: 'PUT'
        });
        showToast(data.message);
        await loadWorlds();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Remove permission
async function removePermission(playerName) {
    try {
        const data = await apiRequest(`/permissions/${encodeURIComponent(playerName)}`, {
            method: 'DELETE'
        });
        showToast(data.message);
        await loadPermissions();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// add function to Global area
window.removeFromAllowlist = removeFromAllowlist;
window.removePermission = removePermission;
window.deleteWorld = deleteWorld;
window.activateWorld = activateWorld;
// API base URL
const API_BASE = '/api';

// DOM elements
const elements = {
    serverStatus: document.getElementById('server-status'),
    startBtn: document.getElementById('start-btn'),
    stopBtn: document.getElementById('stop-btn'),
    restartBtn: document.getElementById('restart-btn'),
    navStartBtn: document.getElementById('nav-start-btn'),
    navStopBtn: document.getElementById('nav-stop-btn'),
    navRestartBtn: document.getElementById('nav-restart-btn'),
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
    resourcepackUploadBtn: document.getElementById('resourcepack-upload-btn'),
    resourcepackUpload: document.getElementById('resourcepack-upload'),
    resourcepacksContainer: document.getElementById('resourcepacks-container'),
    serverVersionsContainer: document.getElementById('server-versions-container'),
    toast: document.getElementById('toast'),
    toastMessage: document.getElementById('toast-message'),
    permissionModal: document.getElementById('permission-modal'),
    modalPlayerName: document.getElementById('modal-player-name'),
    closeModalBtn: document.getElementById('close-modal-btn'),
    cancelModalBtn: document.getElementById('cancel-modal-btn'),
    // Logs elements
    logsContainer: document.getElementById('logs-container'),
    logsContent: document.getElementById('logs-content'),
    logsRefreshBtn: document.getElementById('logs-refresh-btn'),
    logsClearBtn: document.getElementById('logs-clear-btn'),
    logsAutoScroll: document.getElementById('logs-auto-scroll'),
    logsConnectionStatus: document.getElementById('logs-connection-status'),
    // Interaction elements
    interactionStatus: document.getElementById('interaction-status'),
    commandInput: document.getElementById('command-input'),
    sendCommandBtn: document.getElementById('send-command-btn'),
    commandHistory: document.getElementById('command-history'),
    clearHistoryBtn: document.getElementById('clear-history-btn'),
    // Commands elements
    commandCategories: document.getElementById('command-categories'),
    quickCommandsContainer: document.getElementById('quick-commands-container')
};

// WebSocket connection for logs
let logsWebSocket = null;

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
    await loadResourcePacks();
    await loadServerVersions();
    await loadLogs();
    await loadInteractionStatus();
    await loadQuickCommands();
    initializeLogsWebSocket();
}

// Bind event listeners
function bindEvents() {
    // Server control buttons
    if (elements.startBtn) elements.startBtn.addEventListener('click', () => controlServer('start'));
    if (elements.stopBtn) elements.stopBtn.addEventListener('click', () => controlServer('stop'));
    if (elements.restartBtn) elements.restartBtn.addEventListener('click', () => controlServer('restart'));
    
    // Navigation server control buttons
    if (elements.navStartBtn) elements.navStartBtn.addEventListener('click', () => controlServer('start'));
    if (elements.navStopBtn) elements.navStopBtn.addEventListener('click', () => controlServer('stop'));
    if (elements.navRestartBtn) elements.navRestartBtn.addEventListener('click', () => controlServer('restart'));
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

    // Logs events
    if (elements.logsRefreshBtn) elements.logsRefreshBtn.addEventListener('click', loadLogs);
    if (elements.logsClearBtn) elements.logsClearBtn.addEventListener('click', clearLogs);

    // Interaction events
    if (elements.sendCommandBtn) elements.sendCommandBtn.addEventListener('click', sendCommand);
    if (elements.commandInput) {
        elements.commandInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                sendCommand();
            }
        });
        elements.commandInput.addEventListener('input', function() {
            const hasValue = this.value.trim().length > 0;
            if (elements.sendCommandBtn) {
                elements.sendCommandBtn.disabled = !hasValue;
            }
        });
    }
    if (elements.clearHistoryBtn) elements.clearHistoryBtn.addEventListener('click', clearCommandHistory);
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

    // Resource pack upload
    if (elements.resourcepackUploadBtn) elements.resourcepackUploadBtn.addEventListener('click', () => elements.resourcepackUpload.click());
    if (elements.resourcepackUpload) elements.resourcepackUpload.addEventListener('change', uploadResourcePack);

    // Server version update
    const updateVersionsBtn = document.getElementById('update-versions-btn');
    if (updateVersionsBtn) updateVersionsBtn.addEventListener('click', updateServerVersions);
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
            'Please enter player xuid';
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
    const playerXuid = elements.permissionPlayer.value.trim();
    
    if (!playerXuid) {
        const errorMessage = window.i18n ? 
            window.i18n.t('permission.error.empty-name') : 
            'Please enter player xuid';
        showToast(errorMessage, 'error');
        return;
    }
    
    elements.modalPlayerName.textContent = playerXuid;
    elements.permissionModal.classList.remove('hidden');
}

// Hide permission selection modal
function hidePermissionModal() {
    elements.permissionModal.classList.add('hidden');
}

// Set player permission
async function setPlayerPermission(level) {
    const playerXuid = elements.permissionPlayer.value.trim();
    
    if (!playerXuid) {
        const errorMessage = window.i18n ? 
            window.i18n.t('permission.error.empty-name') : 
            'Please enter player xuid';
        showToast(errorMessage, 'error');
        return;
    }
    
    try {
        const data = await apiRequest('/permissions', {
            method: 'PUT',
            body: JSON.stringify({ xuid: playerXuid, level })
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
    
    // Escape special characters in player xuid
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
    
    // Use permission.permission instead of permission.level
    const permissionLevel = permission.permission || 'visitor';
    
    // Escape special characters in permission xuid
    const escapedXuid = (permission.xuid || permission.name || '').replace(/'/g, "\\'").replace(/"/g, '\\"');
    const displayXuid = permission.xuid || permission.name || '';
    
    div.innerHTML = `
        <div>
            <span class="font-medium">${displayXuid}</span>
            <span class="ml-2 px-2 py-1 text-xs rounded ${levelColor[permissionLevel]} bg-gray-200">
                ${levelText[permissionLevel] || permissionLevel}
            </span>
        </div>
        <button onclick="removePermission('${escapedXuid}')" 
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
async function removePermission(playerXuid) {
    try {
        const data = await apiRequest(`/permissions/${encodeURIComponent(playerXuid)}`, {
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
window.activateResourcePack = activateResourcePack;
window.deactivateResourcePack = deactivateResourcePack;
window.deleteResourcePack = deleteResourcePack;
window.downloadServerVersion = downloadServerVersion;
window.activateServerVersion = activateServerVersion;

// Resource pack management functions

// Load resource packs list
async function loadResourcePacks() {
    try {
        const data = await apiRequest('/resource-packs');
        renderResourcePacks(data.resource_packs || []);
    } catch (error) {
        renderResourcePacks([]);
    }
}

// Render resource packs list
function renderResourcePacks(resourcePacks) {
    elements.resourcepacksContainer.innerHTML = '';
    
    if (resourcePacks.length === 0) {
        const emptyMessage = window.i18n ? 
            window.i18n.t('resourcepack.empty') : 
            'No resource packs';
        elements.resourcepacksContainer.innerHTML = `<p class="text-gray-500 text-center py-4">${emptyMessage}</p>`;
        return;
    }
    
    resourcePacks.forEach(pack => {
        const packElement = createResourcePackElement(pack);
        elements.resourcepacksContainer.appendChild(packElement);
    });
}

// Upload resource pack
async function uploadResourcePack() {
    const file = elements.resourcepackUpload.files[0];
    if (!file) return;
    
    const formData = new FormData();
    formData.append('resource_pack', file);
    
    try {
        const response = await fetch(`${API_BASE}/resource-packs/upload`, {
            method: 'POST',
            body: formData
        });
        
        if (!response.ok) {
            throw new Error(`HTTP error! status: ${response.status}`);
        }
        
        const data = await response.json();
        showToast(data.message);
        elements.resourcepackUpload.value = '';
        await loadResourcePacks();
    } catch (error) {
        const errorMessage = window.i18n ? 
            window.i18n.t('resourcepack.upload.error') : 
            'Upload failed: ';
        showToast(errorMessage + error.message, 'error');
    }
}

// Create resource pack element
function createResourcePackElement(pack) {
    const div = document.createElement('div');
    div.className = 'flex items-center justify-between bg-gray-50 px-3 py-2 rounded';
    
    // Escape special characters in pack uuid
    const escapedUuid = pack.uuid.replace(/'/g, "\\'").replace(/"/g, '\\"');
    
    const activeText = window.i18n ? window.i18n.t('resourcepack.active') : 'Active';
    const activateText = window.i18n ? window.i18n.t('resourcepack.activate') : 'Activate';
    const deactivateText = window.i18n ? window.i18n.t('resourcepack.deactivate') : 'Deactivate';
    const deleteText = window.i18n ? window.i18n.t('resourcepack.delete') : 'Delete';
    
    div.innerHTML = `
        <div>
            <span class="font-medium">${pack.name}</span>
            <span class="ml-2 text-sm text-gray-500">v${pack.version.join('.')}</span>
            ${pack.active ? `<span class="ml-2 px-2 py-1 text-xs rounded bg-green-200 text-green-800">${activeText}</span>` : ''}
        </div>
        <div class="space-x-2">
            ${!pack.active ? `<button onclick="activateResourcePack('${escapedUuid}')" 
                class="text-blue-500 hover:text-blue-700 transition duration-200" title="${activateText}">
                <i class="fas fa-play"></i>
            </button>` : `<button onclick="deactivateResourcePack('${escapedUuid}')" 
                class="text-orange-500 hover:text-orange-700 transition duration-200" title="${deactivateText}">
                <i class="fas fa-pause"></i>
            </button>`}
            <button onclick="deleteResourcePack('${escapedUuid}')" 
                    class="text-red-500 hover:text-red-700 transition duration-200" title="${deleteText}">
                <i class="fas fa-trash"></i>
            </button>
        </div>
    `;
    
    return div;
}

// Activate resource pack
async function activateResourcePack(uuid) {
    try {
        const data = await apiRequest(`/resource-packs/${encodeURIComponent(uuid)}/activate`, {
            method: 'PUT'
        });
        showToast(data.message);
        await loadResourcePacks();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Server version management functions

// Load server versions list
async function loadServerVersions() {
    try {
        const data = await apiRequest('/server-versions');
        renderServerVersions(data.data || []);
    } catch (error) {
        renderServerVersions([]);
    }
}

// Render server versions list
function renderServerVersions(versions) {
    elements.serverVersionsContainer.innerHTML = '';
    
    if (versions.length === 0) {
        const emptyMessage = window.i18n ? 
            window.i18n.t('server.versions.empty') : 
            'No server versions available';
        elements.serverVersionsContainer.innerHTML = `<p class="text-gray-500 text-center py-4">${emptyMessage}</p>`;
        return;
    }
    
    versions.forEach(version => {
        const versionElement = createServerVersionElement(version);
        elements.serverVersionsContainer.appendChild(versionElement);
    });
}

// Create server version element
function createServerVersionElement(version) {
    const div = document.createElement('div');
    div.className = 'border rounded-lg p-4 bg-gray-50';
    
    const downloadedText = window.i18n ? window.i18n.t('server.versions.downloaded') : 'Downloaded';
    const activeText = window.i18n ? window.i18n.t('server.versions.active') : 'Active';
    const downloadText = window.i18n ? window.i18n.t('server.versions.download') : 'Download';
    const activateText = window.i18n ? window.i18n.t('server.versions.activate') : 'Activate';
    const downloadingText = window.i18n ? window.i18n.t('server.versions.downloading') : 'Downloading...';
    
    let statusBadge = '';
    if (version.active) {
        statusBadge = `<span class="px-2 py-1 text-xs rounded bg-green-200 text-green-800">${activeText}</span>`;
    } else if (version.downloaded) {
        statusBadge = `<span class="px-2 py-1 text-xs rounded bg-blue-200 text-blue-800">${downloadedText}</span>`;
    }
    
    let actionButton = '';
    if (!version.downloaded) {
        actionButton = `
            <button onclick="downloadServerVersion('${version.version}')" 
                    class="bg-blue-500 hover:bg-blue-600 text-white px-4 py-2 rounded transition duration-200">
                <i class="fas fa-download mr-2"></i>${downloadText}
            </button>
        `;
    } else if (!version.active) {
        actionButton = `
            <button onclick="activateServerVersion('${version.version}')" 
                    class="bg-green-500 hover:bg-green-600 text-white px-4 py-2 rounded transition duration-200">
                <i class="fas fa-play mr-2"></i>${activateText}
            </button>
        `;
    }
    
    div.innerHTML = `
        <div class="flex items-center justify-between">
            <div>
                <h3 class="text-lg font-semibold">Bedrock Server ${version.version}</h3>
                <p class="text-sm text-gray-600 mt-1">Minecraft Bedrock Dedicated Server</p>
                <div class="mt-2">
                    ${statusBadge}
                </div>
            </div>
            <div class="flex items-center space-x-2">
                <div id="progress-${version.version}" class="hidden">
                    <div class="w-48 bg-gray-200 rounded-full h-2">
                        <div class="bg-blue-600 h-2 rounded-full transition-all duration-300" style="width: 0%"></div>
                    </div>
                    <p class="text-xs text-gray-600 mt-1">${downloadingText}</p>
                </div>
                ${actionButton}
            </div>
        </div>
    `;
    
    return div;
}

// Download server version
async function downloadServerVersion(version) {
    try {
        const data = await apiRequest(`/server-versions/${version}/download`, {
            method: 'POST'
        });
        showToast(data.message);
        
        // Show progress bar
        const progressContainer = document.getElementById(`progress-${version}`);
        if (progressContainer) {
            progressContainer.classList.remove('hidden');
        }
        
        // Start polling for progress
        pollDownloadProgress(version);
        
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Poll download progress
async function pollDownloadProgress(version) {
    try {
        const data = await apiRequest(`/server-versions/${version}/progress`);
        const progress = data.data;
        
        // Update progress bar
        const progressContainer = document.getElementById(`progress-${version}`);
        if (progressContainer) {
            const progressBar = progressContainer.querySelector('.bg-blue-600');
            if (progressBar) {
                progressBar.style.width = `${progress.progress}%`;
            }
            
            const progressText = progressContainer.querySelector('.text-xs');
            if (progressText) {
                progressText.textContent = progress.message;
            }
        }
        
        // Continue polling if still downloading
        if (progress.status === 'downloading' || progress.status === 'extracting') {
            setTimeout(() => pollDownloadProgress(version), 1000);
        } else {
            // Download completed or failed, refresh the list
            setTimeout(() => {
                loadServerVersions();
            }, 1000);
        }
        
    } catch (error) {
        // Progress not found, download might be completed
        setTimeout(() => {
            loadServerVersions();
        }, 1000);
    }
}

// Activate server version
async function activateServerVersion(version) {
    const confirmMessage = window.i18n ? 
        window.i18n.t('server.versions.activate-confirm', { version }) : 
        `Are you sure you want to activate server version ${version}? This will change the active server configuration.`;
    
    if (!confirm(confirmMessage)) {
        return;
    }
    
    try {
        const data = await apiRequest(`/server-versions/${version}/activate`, {
            method: 'PUT'
        });
        showToast(data.message);
        await loadServerVersions();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Deactivate resource pack
async function deactivateResourcePack(uuid) {
    try {
        const data = await apiRequest(`/resource-packs/${encodeURIComponent(uuid)}/deactivate`, {
            method: 'PUT'
        });
        showToast(data.message);
        await loadResourcePacks();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Delete resource pack
async function deleteResourcePack(uuid) {
    const confirmMessage = window.i18n ? 
        window.i18n.t('resourcepack.deleteConfirm') : 
        'Are you sure you want to delete this resource pack? This action cannot be undone!';
    
    if (!confirm(confirmMessage)) {
        return;
    }
    
    try {
        const data = await apiRequest(`/resource-packs/${encodeURIComponent(uuid)}`, {
            method: 'DELETE'
        });
        showToast(data.message);
        await loadResourcePacks();
    } catch (error) {
        // Error already handled in apiRequest
    }
}

// Update server versions from GitHub
async function updateServerVersions() {
    const updateBtn = document.getElementById('update-versions-btn');
    const originalText = updateBtn.innerHTML;
    
    // Show loading state
    updateBtn.disabled = true;
    updateBtn.innerHTML = '<i class="fas fa-spinner fa-spin mr-2"></i>' + 
        (window.i18n ? window.i18n.t('server.versions.updating') : 'Updating...');
    
    try {
        const data = await apiRequest('/server-versions/update-config', {
            method: 'POST'
        });
        
        showToast(data.message);
        
        // Update the versions list with new data
        if (data.data) {
            renderServerVersions(data.data);
        } else {
            // Reload versions if no data returned
            await loadServerVersions();
        }
        
    } catch (error) {
        // Error already handled in apiRequest
    } finally {
        // Restore button state
        updateBtn.disabled = false;
        updateBtn.innerHTML = originalText;
    }
}

// Quick action functions for dashboard buttons
async function startServer() {
    await controlServer('start');
}

async function stopServer() {
    await controlServer('stop');
}

// ===== LOGS FUNCTIONALITY =====

// Load logs from server
async function loadLogs() {
    try {
        const response = await apiRequest('/logs');
        if (response.logs) {
            renderLogs(response.logs);
        }
    } catch (error) {
        console.error('Failed to load logs:', error);
        const errorMessage = window.i18n ? window.i18n.t('logs.load-failed') : 'Failed to load logs';
        showToast(errorMessage, 'error');
    }
}

// Render logs in the container
function renderLogs(logs) {
    if (!elements.logsContent) return;
    
    elements.logsContent.innerHTML = '';
    
    if (logs.length === 0) {
        const noLogsText = window.i18n ? window.i18n.t('logs.no-logs') : 'No logs available';
        elements.logsContent.innerHTML = `<div class="text-gray-500">${noLogsText}</div>`;
        return;
    }
    
    logs.forEach(log => {
        const logElement = createLogElement(log);
        elements.logsContent.appendChild(logElement);
    });
    
    // Auto scroll to bottom if enabled
    if (elements.logsAutoScroll && elements.logsAutoScroll.checked) {
        elements.logsContainer.scrollTop = elements.logsContainer.scrollHeight;
    }
}

// Create a log element
function createLogElement(log) {
    const div = document.createElement('div');
    div.className = 'log-entry mb-1';
    
    const timestamp = new Date(log.timestamp).toLocaleTimeString();
    const levelClass = getLevelClass(log.level);
    
    div.innerHTML = `
        <span class="text-gray-400">[${timestamp}]</span>
        <span class="${levelClass} font-semibold">[${log.level}]</span>
        <span>${escapeHtml(log.message)}</span>
    `;
    
    return div;
}

// Get CSS class for log level
function getLevelClass(level) {
    switch (level.toUpperCase()) {
        case 'ERROR': return 'text-red-400';
        case 'WARN': return 'text-yellow-400';
        case 'INFO': return 'text-blue-400';
        case 'DEBUG': return 'text-gray-400';
        default: return 'text-green-400';
    }
}

// Clear logs
async function clearLogs() {
    try {
        await apiRequest('/logs', { method: 'DELETE' });
        if (elements.logsContent) {
            const logsClearedText = window.i18n ? window.i18n.t('logs.cleared') : 'Logs cleared';
            elements.logsContent.innerHTML = `<div class="text-gray-500">${logsClearedText}</div>`;
        }
        const successMessage = window.i18n ? window.i18n.t('logs.clear-success') : 'Logs cleared successfully';
        showToast(successMessage);
    } catch (error) {
        console.error('Failed to clear logs:', error);
        const errorMessage = window.i18n ? window.i18n.t('logs.clear-failed') : 'Failed to clear logs';
        showToast(errorMessage, 'error');
    }
}

// Initialize WebSocket connection for real-time logs
function initializeLogsWebSocket() {
    if (!elements.logsConnectionStatus) return;
    
    const protocol = window.location.protocol === 'https:' ? 'wss:' : 'ws:';
    const wsUrl = `${protocol}//${window.location.host}/api/logs/ws`;
    
    try {
        logsWebSocket = new WebSocket(wsUrl);
        
        logsWebSocket.onopen = function() {
            const connectedText = window.i18n ? window.i18n.t('logs.connected') : '已连接';
            elements.logsConnectionStatus.textContent = connectedText;
            elements.logsConnectionStatus.className = 'font-medium text-green-600';
        };
        
        logsWebSocket.onmessage = function(event) {
            try {
                const log = JSON.parse(event.data);
                appendLogEntry(log);
            } catch (error) {
                console.error('Failed to parse log message:', error);
            }
        };
        
        logsWebSocket.onclose = function() {
            const disconnectedText = window.i18n ? window.i18n.t('logs.disconnected') : '已断开';
            elements.logsConnectionStatus.textContent = disconnectedText;
            elements.logsConnectionStatus.className = 'font-medium text-red-600';
            
            // Attempt to reconnect after 3 seconds
            setTimeout(initializeLogsWebSocket, 3000);
        };
        
        logsWebSocket.onerror = function(error) {
            console.error('WebSocket error:', error);
            const errorText = window.i18n ? window.i18n.t('logs.connection-error') : '连接错误';
            elements.logsConnectionStatus.textContent = errorText;
            elements.logsConnectionStatus.className = 'font-medium text-red-600';
        };
    } catch (error) {
        console.error('Failed to create WebSocket:', error);
        const failedText = window.i18n ? window.i18n.t('logs.connection-failed') : '连接失败';
        elements.logsConnectionStatus.textContent = failedText;
        elements.logsConnectionStatus.className = 'font-medium text-red-600';
    }
}

// Append a single log entry to the container
function appendLogEntry(log) {
    if (!elements.logsContent) return;
    
    const logElement = createLogElement(log);
    elements.logsContent.appendChild(logElement);
    
    // Auto scroll to bottom if enabled
    if (elements.logsAutoScroll && elements.logsAutoScroll.checked) {
        elements.logsContainer.scrollTop = elements.logsContainer.scrollHeight;
    }
    
    // Remove old entries if too many (keep last 1000)
    const entries = elements.logsContent.children;
    if (entries.length > 1000) {
        elements.logsContent.removeChild(entries[0]);
    }
}

// ===== INTERACTION FUNCTIONALITY =====

// Load interaction status
async function loadInteractionStatus() {
    try {
        const response = await apiRequest('/interaction/status');
        updateInteractionStatus(response.enabled);
        if (response.enabled) {
            await loadCommandHistory();
        }
    } catch (error) {
        console.error('Failed to load interaction status:', error);
        updateInteractionStatus(false);
    }
}

// Update interaction status display
function updateInteractionStatus(enabled) {
    if (!elements.interactionStatus) return;
    
    if (enabled) {
        const enabledText = window.i18n ? window.i18n.t('interaction.enabled') : '命令交互已启用';
        elements.interactionStatus.innerHTML = `
            <div class="bg-green-100 border border-green-400 text-green-700 px-3 py-2 rounded">
                <i class="fas fa-check-circle mr-2"></i>
                ${enabledText}
            </div>
        `;
        if (elements.sendCommandBtn) elements.sendCommandBtn.disabled = false;
        if (elements.commandInput) elements.commandInput.disabled = false;
    } else {
        const disabledText = window.i18n ? window.i18n.t('interaction.disabled') : '命令交互在当前平台不可用';
        elements.interactionStatus.innerHTML = `
            <div class="bg-yellow-100 border border-yellow-400 text-yellow-700 px-3 py-2 rounded">
                <i class="fas fa-exclamation-triangle mr-2"></i>
                ${disabledText}
            </div>
        `;
        if (elements.sendCommandBtn) elements.sendCommandBtn.disabled = true;
        if (elements.commandInput) elements.commandInput.disabled = true;
    }
}

// Send command to server
async function sendCommand() {
    const command = elements.commandInput?.value.trim();
    if (!command) return;
    
    try {
        const response = await apiRequest('/interaction/command', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ command })
        });
        
        if (elements.commandInput) elements.commandInput.value = '';
        if (elements.sendCommandBtn) elements.sendCommandBtn.disabled = true;
        
        const successText = window.i18n ? window.i18n.t('interaction.command-sent') : 'Command sent successfully';
        showToast(successText);
        await loadCommandHistory();
    } catch (error) {
        console.error('Failed to send command:', error);
        const failedText = window.i18n ? window.i18n.t('interaction.send-failed') : 'Failed to send command';
        showToast(error.message || failedText, 'error');
    }
}

// Load command history
async function loadCommandHistory() {
    try {
        const response = await apiRequest('/interaction/history');
        if (response.history) {
            renderCommandHistory(response.history);
        }
    } catch (error) {
        console.error('Failed to load command history:', error);
    }
}

// Render command history
function renderCommandHistory(history) {
    if (!elements.commandHistory) return;
    
    elements.commandHistory.innerHTML = '';
    
    if (history.length === 0) {
        const noHistoryText = window.i18n ? window.i18n.t('interaction.no-history') : 'No command history';
        elements.commandHistory.innerHTML = `<div class="text-gray-500 text-sm">${noHistoryText}</div>`;
        return;
    }
    
    history.slice(-10).reverse().forEach(entry => {
        const historyElement = createCommandHistoryElement(entry);
        elements.commandHistory.appendChild(historyElement);
    });
}

// Create command history element
function createCommandHistoryElement(entry) {
    const div = document.createElement('div');
    div.className = 'bg-gray-50 p-2 rounded text-sm cursor-pointer hover:bg-gray-100 transition-colors duration-200';
    
    const timestamp = new Date(entry.timestamp).toLocaleTimeString();
    
    div.innerHTML = `
        <div class="flex justify-between items-start">
            <div class="flex-1">
                <div class="font-mono text-blue-600">${escapeHtml(entry.command)}</div>
                <div class="text-gray-600 mt-1">${escapeHtml(entry.response)}</div>
            </div>
            <div class="text-xs text-gray-400 ml-2">${timestamp}</div>
        </div>
    `;
    
    // Add click event to fill command input
    div.addEventListener('click', function() {
        if (elements.commandInput) {
            elements.commandInput.value = entry.command;
            elements.commandInput.focus();
            // Trigger input event to enable send button
            elements.commandInput.dispatchEvent(new Event('input'));
        }
    });
    
    return div;
}

// Clear command history
async function clearCommandHistory() {
    try {
        await apiRequest('/interaction/history', { method: 'DELETE' });
        if (elements.commandHistory) {
            const clearedText = window.i18n ? window.i18n.t('interaction.history-cleared') : 'Command history cleared';
            elements.commandHistory.innerHTML = `<div class="text-gray-500 text-sm">${clearedText}</div>`;
        }
        const successText = window.i18n ? window.i18n.t('interaction.clear-history-success') : 'Command history cleared successfully';
        showToast(successText);
    } catch (error) {
        console.error('Failed to clear command history:', error);
        const failedText = window.i18n ? window.i18n.t('interaction.clear-history-failed') : 'Failed to clear command history';
        showToast(failedText, 'error');
    }
}

// ===== QUICK COMMANDS FUNCTIONALITY =====

// Load quick commands
async function loadQuickCommands() {
    try {
        const [commandsResponse, categoriesResponse] = await Promise.all([
            apiRequest('/commands'),
            apiRequest('/commands/categories')
        ]);
        
        if (categoriesResponse.categories) {
            renderCommandCategories(categoriesResponse.categories);
        }
        
        if (commandsResponse.commands) {
            renderQuickCommands(commandsResponse.commands);
        }
    } catch (error) {
        console.error('Failed to load quick commands:', error);
        const failedText = window.i18n ? window.i18n.t('commands.load-failed') : 'Failed to load quick commands';
        showToast(failedText, 'error');
    }
}

// Render command categories
function renderCommandCategories(categories) {
    if (!elements.commandCategories) return;
    
    elements.commandCategories.innerHTML = '';
    
    // Add "All" category
    const allBtn = document.createElement('button');
    allBtn.className = 'px-3 py-1 bg-blue-500 text-white rounded text-sm hover:bg-blue-600 transition duration-200';
    const allText = window.i18n ? window.i18n.t('commands.all') : '全部';
    allBtn.textContent = allText;
    allBtn.addEventListener('click', () => filterCommandsByCategory(''));
    elements.commandCategories.appendChild(allBtn);
    
    categories.forEach(category => {
        const btn = document.createElement('button');
        btn.className = 'px-3 py-1 bg-gray-200 text-gray-700 rounded text-sm hover:bg-gray-300 transition duration-200';
        btn.textContent = getCategoryDisplayName(category);
        btn.addEventListener('click', () => filterCommandsByCategory(category));
        elements.commandCategories.appendChild(btn);
    });
}

// Get display name for category
function getCategoryDisplayName(category) {
    if (!window.i18n) {
        const categoryNames = {
            'time': '时间',
            'weather': '天气',
            'gamemode': '游戏模式',
            'difficulty': '难度'
        };
        return categoryNames[category] || category;
    }
    
    const categoryKey = `commands.${category}`;
    return window.i18n.t(categoryKey) || category;
}

// Filter commands by category
async function filterCommandsByCategory(category) {
    try {
        const url = category ? `/commands?category=${encodeURIComponent(category)}` : '/commands';
        const response = await apiRequest(url);
        
        if (response.commands) {
            renderQuickCommands(response.commands);
        }
        
        // Update active category button
        const buttons = elements.commandCategories?.querySelectorAll('button');
        buttons?.forEach((btn, index) => {
            if ((index === 0 && !category) || (index > 0 && btn.textContent === getCategoryDisplayName(category))) {
                btn.className = 'px-3 py-1 bg-blue-500 text-white rounded text-sm hover:bg-blue-600 transition duration-200';
            } else {
                btn.className = 'px-3 py-1 bg-gray-200 text-gray-700 rounded text-sm hover:bg-gray-300 transition duration-200';
            }
        });
    } catch (error) {
        console.error('Failed to filter commands:', error);
        const failedText = window.i18n ? window.i18n.t('commands.filter-failed') : 'Failed to filter commands';
        showToast(failedText, 'error');
    }
}

// Render quick commands
function renderQuickCommands(commands) {
    if (!elements.quickCommandsContainer) return;
    
    elements.quickCommandsContainer.innerHTML = '';
    
    if (commands.length === 0) {
        const noCommandsText = window.i18n ? window.i18n.t('commands.no-commands') : 'No commands available';
        elements.quickCommandsContainer.innerHTML = `<div class="col-span-full text-gray-500 text-center py-8">${noCommandsText}</div>`;
        return;
    }
    
    commands.forEach(command => {
        const commandElement = createQuickCommandElement(command);
        elements.quickCommandsContainer.appendChild(commandElement);
    });
}

// Create quick command element
function createQuickCommandElement(command) {
    const div = document.createElement('div');
    div.className = 'bg-white border border-gray-200 rounded-lg p-4 hover:shadow-md transition duration-200';
    
    div.innerHTML = `
        <div class="flex justify-between items-start mb-2">
            <h4 class="font-semibold text-gray-800">${escapeHtml(command.name)}</h4>
            <span class="text-xs bg-gray-100 text-gray-600 px-2 py-1 rounded">${getCategoryDisplayName(command.category)}</span>
        </div>
        <p class="text-sm text-gray-600 mb-3">${escapeHtml(command.description)}</p>
        <div class="flex justify-between items-center">
            <code class="text-xs bg-gray-100 px-2 py-1 rounded font-mono">${escapeHtml(command.command)}</code>
            <button class="bg-green-500 hover:bg-green-600 text-white px-3 py-1 rounded text-sm transition duration-200" 
                    onclick="executeQuickCommand('${command.id}')">
                <i class="fas fa-play mr-1"></i>${window.i18n ? window.i18n.t('commands.execute') : '执行'}
            </button>
        </div>
    `;
    
    return div;
}

// Execute quick command
async function executeQuickCommand(commandId) {
    try {
        const response = await apiRequest(`/commands/${commandId}/execute`, {
            method: 'POST'
        });
        
        const executedText = window.i18n ? window.i18n.t('commands.executed') : 'Command executed';
        showToast(`${executedText}: ${response.command}`);
        await loadCommandHistory();
    } catch (error) {
        console.error('Failed to execute command:', error);
        const failedText = window.i18n ? window.i18n.t('commands.execute-failed') : 'Failed to execute command';
        showToast(error.message || failedText, 'error');
    }
}

// Utility function to escape HTML
function escapeHtml(text) {
    const div = document.createElement('div');
    div.textContent = text;
    return div.innerHTML;
}

// ===== PERFORMANCE MONITORING FUNCTIONALITY =====

// Performance monitoring variables
let performanceMonitoringInterval = null;

// Load performance monitoring data
async function loadPerformanceMonitoring() {
    try {
        const data = await apiRequest('/monitor/performance');
        updatePerformanceDisplay(data);
    } catch (error) {
        console.error('Failed to load performance monitoring data:', error);
        // Reset display on error
        updatePerformanceDisplay(null);
    }
}

// Update performance display
function updatePerformanceDisplay(data) {
    const systemCpuElement = document.getElementById('system-cpu');
    const systemMemoryElement = document.getElementById('system-memory');
    const bedrockCpuElement = document.getElementById('bedrock-cpu');
    const bedrockMemoryElement = document.getElementById('bedrock-memory');
    const bedrockStatusElement = document.getElementById('bedrock-status');

    if (!data) {
        // Reset all displays
        if (systemCpuElement) systemCpuElement.textContent = '--';
        if (systemMemoryElement) systemMemoryElement.textContent = '--';
        if (bedrockCpuElement) bedrockCpuElement.textContent = '--';
        if (bedrockMemoryElement) bedrockMemoryElement.textContent = '--';
        if (bedrockStatusElement) {
            bedrockStatusElement.textContent = window.i18n ? window.i18n.t('dashboard.performance.bedrock-stopped') : 'Bedrock服务器未运行';
        }
        return;
    }

    // Update system performance
    if (systemCpuElement) {
        systemCpuElement.textContent = `${data.system.cpu_usage.toFixed(1)}%`;
    }
    if (systemMemoryElement) {
        systemMemoryElement.textContent = `${data.system.memory_usage.toFixed(1)}%`;
    }

    // Update bedrock process performance
    if (data.bedrock.pid > 0) {
        if (bedrockCpuElement) {
            bedrockCpuElement.textContent = `${data.bedrock.cpu_usage.toFixed(1)}%`;
        }
        if (bedrockMemoryElement) {
            bedrockMemoryElement.textContent = `${data.bedrock.memory_mb.toFixed(1)}MB`;
        }
        if (bedrockStatusElement) {
            const statusText = window.i18n ? 
                window.i18n.t('dashboard.performance.bedrock-running', { pid: data.bedrock.pid }) : 
                `PID: ${data.bedrock.pid}`;
            bedrockStatusElement.textContent = statusText;
        }
    } else {
        if (bedrockCpuElement) bedrockCpuElement.textContent = '--';
        if (bedrockMemoryElement) bedrockMemoryElement.textContent = '--';
        if (bedrockStatusElement) {
            bedrockStatusElement.textContent = window.i18n ? 
                window.i18n.t('dashboard.performance.bedrock-stopped') : 
                'Bedrock服务器未运行';
        }
    }
}

// Start performance monitoring
function startPerformanceMonitoring() {
    // Load initial data
    loadPerformanceMonitoring();
    
    // Set up interval for periodic updates (every 5 seconds)
    if (performanceMonitoringInterval) {
        clearInterval(performanceMonitoringInterval);
    }
    performanceMonitoringInterval = setInterval(loadPerformanceMonitoring, 5000);
}

// Stop performance monitoring
function stopPerformanceMonitoring() {
    if (performanceMonitoringInterval) {
        clearInterval(performanceMonitoringInterval);
        performanceMonitoringInterval = null;
    }
}

// Initialize performance monitoring when dashboard is active
function initializePerformanceMonitoring() {
    const dashboardSection = document.getElementById('dashboard');
    if (dashboardSection && dashboardSection.classList.contains('active')) {
        startPerformanceMonitoring();
    }
}

// Add performance monitoring to navigation handling
const originalShowSection = window.showSection;
window.showSection = function(sectionId) {
    if (originalShowSection) {
        originalShowSection(sectionId);
    }
    
    // Start/stop performance monitoring based on active section
    if (sectionId === 'dashboard') {
        startPerformanceMonitoring();
    } else {
        stopPerformanceMonitoring();
    }
};

// Initialize performance monitoring on page load
document.addEventListener('DOMContentLoaded', function() {
    // Delay initialization to ensure other components are ready
    setTimeout(initializePerformanceMonitoring, 1000);
});

// Make functions globally available
window.executeQuickCommand = executeQuickCommand;
window.loadPerformanceMonitoring = loadPerformanceMonitoring;
window.startPerformanceMonitoring = startPerformanceMonitoring;
window.stopPerformanceMonitoring = stopPerformanceMonitoring;

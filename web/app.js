// API基础URL
const API_BASE = '/api';

// DOM元素
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

// 初始化应用
document.addEventListener('DOMContentLoaded', function() {
    initializeApp();
    bindEvents();
});

// 初始化应用数据
async function initializeApp() {
    await loadServerStatus();
    await loadServerConfig();
    await loadAllowlist();
    await loadPermissions();
    await loadWorlds();
}

// 绑定事件监听器
function bindEvents() {
    // 服务器控制按钮
    if (elements.startBtn) elements.startBtn.addEventListener('click', () => controlServer('start'));
    if (elements.stopBtn) elements.stopBtn.addEventListener('click', () => controlServer('stop'));
    if (elements.restartBtn) elements.restartBtn.addEventListener('click', () => controlServer('restart'));
    if (elements.refreshBtn) elements.refreshBtn.addEventListener('click', initializeApp);

    // 配置表单
    if (elements.configForm) elements.configForm.addEventListener('submit', saveServerConfig);

    // 白名单管理
    if (elements.addPlayerBtn) elements.addPlayerBtn.addEventListener('click', addToAllowlist);
    if (elements.newPlayerInput) {
        elements.newPlayerInput.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                addToAllowlist();
            }
        });
    }

    // 权限管理
    if (elements.addPermissionBtn) elements.addPermissionBtn.addEventListener('click', showPermissionModal);
    if (elements.permissionPlayer) {
        elements.permissionPlayer.addEventListener('keypress', function(e) {
            if (e.key === 'Enter') {
                showPermissionModal();
            }
        });
    }

    // 弹窗事件 - 添加存在性检查
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

    // 权限选项点击事件
    document.addEventListener('click', function(e) {
        if (e.target.closest('.permission-option')) {
            const level = e.target.closest('.permission-option').dataset.level;
            setPlayerPermission(level);
        }
    });

    // 世界上传
    if (elements.uploadBtn) elements.uploadBtn.addEventListener('click', () => elements.worldUpload.click());
    if (elements.worldUpload) elements.worldUpload.addEventListener('change', uploadWorld);
}

// API请求封装
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
        console.error('API请求失败:', error);
        showToast('请求失败: ' + error.message, 'error');
        throw error;
    }
}

// 显示提示消息
function showToast(message, type = 'success') {
    elements.toastMessage.textContent = message;
    elements.toast.className = `fixed top-4 right-4 px-6 py-3 rounded-lg shadow-lg transform transition-transform duration-300 ${
        type === 'error' ? 'bg-red-500' : 'bg-green-500'
    } text-white`;
    
    // 显示toast
    elements.toast.style.transform = 'translateX(0)';
    
    // 3秒后隐藏
    setTimeout(() => {
        elements.toast.style.transform = 'translateX(100%)';
    }, 3000);
}

// 加载服务器状态
async function loadServerStatus() {
    try {
        const data = await apiRequest('/status');
        updateServerStatus(data.status);
    } catch (error) {
        updateServerStatus('unknown');
    }
}

// 更新服务器状态显示
function updateServerStatus(status) {
    const statusElement = elements.serverStatus;
    statusElement.className = 'px-3 py-1 rounded-full text-sm';
    
    switch (status) {
        case 'running':
            statusElement.textContent = '运行中';
            statusElement.classList.add('bg-green-500');
            break;
        case 'stopped':
            statusElement.textContent = '已停止';
            statusElement.classList.add('bg-red-500');
            break;
        default:
            statusElement.textContent = '未知';
            statusElement.classList.add('bg-gray-500');
    }
}

// 服务器控制
async function controlServer(action) {
    try {
        const data = await apiRequest(`/${action}`, { method: 'POST' });
        showToast(data.message);
        
        // 延迟刷新状态
        setTimeout(loadServerStatus, 2000);
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 加载服务器配置
async function loadServerConfig() {
    try {
        const data = await apiRequest('/config');
        if (data.config) {
            populateConfigForm(data.config);
        }
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 填充配置表单
function populateConfigForm(config) {
    document.getElementById('server-name').value = config['server-name'] || '';
    document.getElementById('gamemode').value = config.gamemode || 'survival';
    document.getElementById('difficulty').value = config.difficulty || 'easy';
    document.getElementById('max-players').value = config['max-players'] || 10;
    document.getElementById('server-port').value = config['server-port'] || 19132;
    document.getElementById('allow-cheats').checked = config['allow-cheats'] === 'true';
    document.getElementById('allow-list').checked = config['allow-list'] === 'true';
}

// 保存服务器配置
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
        // 错误已在apiRequest中处理
    }
}

// 加载白名单
async function loadAllowlist() {
    try {
        const data = await apiRequest('/allowlist');
        renderAllowlist(data.allowlist || []);
    } catch (error) {
        renderAllowlist([]);
    }
}

// 渲染白名单
function renderAllowlist(allowlist) {
    elements.allowlistContainer.innerHTML = '';
    
    if (allowlist.length === 0) {
        elements.allowlistContainer.innerHTML = '<p class="text-gray-500 text-center py-4">暂无白名单用户</p>';
        return;
    }
    
    allowlist.forEach(player => {
        const playerElement = createPlayerElement(player, 'allowlist');
        elements.allowlistContainer.appendChild(playerElement);
    });
}

// 添加到白名单
async function addToAllowlist() {
    const playerName = elements.newPlayerInput.value.trim();
    if (!playerName) {
        showToast('请输入玩家名称', 'error');
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
        // 错误已在apiRequest中处理
    }
}

// 从白名单移除
async function removeFromAllowlist(playerName) {
    try {
        const data = await apiRequest(`/allowlist/${encodeURIComponent(playerName)}`, {
            method: 'DELETE'
        });
        showToast(data.message);
        await loadAllowlist();
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 加载权限
async function loadPermissions() {
    try {
        const data = await apiRequest('/permissions');
        renderPermissions(data.permissions || []);
    } catch (error) {
        renderPermissions([]);
    }
}

// 渲染权限
function renderPermissions(permissions) {
    elements.permissionsContainer.innerHTML = '';
    
    if (permissions.length === 0) {
        elements.permissionsContainer.innerHTML = '<p class="text-gray-500 text-center py-4">暂无权限设置</p>';
        return;
    }
    
    permissions.forEach(permission => {
        const permissionElement = createPermissionElement(permission);
        elements.permissionsContainer.appendChild(permissionElement);
    });
}

// 显示权限选择弹窗
function showPermissionModal() {
    const playerName = elements.permissionPlayer.value.trim();
    
    if (!playerName) {
        showToast('请输入玩家名称', 'error');
        return;
    }
    
    elements.modalPlayerName.textContent = playerName;
    elements.permissionModal.classList.remove('hidden');
}

// 隐藏权限选择弹窗
function hidePermissionModal() {
    elements.permissionModal.classList.add('hidden');
}

// 设置玩家权限
async function setPlayerPermission(level) {
    const playerName = elements.permissionPlayer.value.trim();
    
    if (!playerName) {
        showToast('请输入玩家名称', 'error');
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
        // 错误已在apiRequest中处理
    }
}

// 加载世界列表
async function loadWorlds() {
    try {
        const data = await apiRequest('/worlds');
        renderWorlds(data.worlds || []);
    } catch (error) {
        renderWorlds([]);
    }
}

// 渲染世界列表
function renderWorlds(worlds) {
    elements.worldsContainer.innerHTML = '';
    
    if (worlds.length === 0) {
        elements.worldsContainer.innerHTML = '<p class="text-gray-500 text-center py-4">暂无世界文件</p>';
        return;
    }
    
    worlds.forEach(world => {
        const worldElement = createWorldElement(world);
        elements.worldsContainer.appendChild(worldElement);
    });
}

// 上传世界
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
        showToast('上传失败: ' + error.message, 'error');
    }
}

// 创建玩家元素
function createPlayerElement(playerName, type) {
    const div = document.createElement('div');
    div.className = 'flex items-center justify-between bg-gray-50 px-3 py-2 rounded';
    
    // 转义玩家名称中的特殊字符
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

// 创建权限元素
function createPermissionElement(permission) {
    const div = document.createElement('div');
    div.className = 'flex items-center justify-between bg-gray-50 px-3 py-2 rounded';
    
    const levelText = {
        'visitor': '访客',
        'member': '成员',
        'operator': '管理员'
    };
    
    const levelColor = {
        'visitor': 'text-gray-600',
        'member': 'text-blue-600',
        'operator': 'text-red-600'
    };
    
    // 转义权限名称中的特殊字符
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

// 创建世界元素
function createWorldElement(world) {
    const div = document.createElement('div');
    div.className = 'flex items-center justify-between bg-gray-50 px-3 py-2 rounded';
    
    // 转义世界名称中的特殊字符
    const escapedName = world.name.replace(/'/g, "\\'").replace(/"/g, '\\"');
    
    div.innerHTML = `
        <div>
            <span class="font-medium">${world.name}</span>
            ${world.active ? '<span class="ml-2 px-2 py-1 text-xs rounded bg-green-200 text-green-800">当前世界</span>' : ''}
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

// 删除世界
async function deleteWorld(worldName) {
    if (!confirm(`确定要删除世界 "${worldName}" 吗？此操作不可撤销！`)) {
        return;
    }
    
    try {
        const data = await apiRequest(`/worlds/${encodeURIComponent(worldName)}`, {
            method: 'DELETE'
        });
        showToast(data.message);
        await loadWorlds();
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 激活世界
async function activateWorld(worldName) {
    try {
        const data = await apiRequest(`/worlds/${encodeURIComponent(worldName)}/activate`, {
            method: 'PUT'
        });
        showToast(data.message);
        await loadWorlds();
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 移除权限
async function removePermission(playerName) {
    try {
        const data = await apiRequest(`/permissions/${encodeURIComponent(playerName)}`, {
            method: 'DELETE'
        });
        showToast(data.message);
        await loadPermissions();
    } catch (error) {
        // 错误已在apiRequest中处理
    }
}

// 将函数添加到全局作用域，以便onclick可以访问
window.removeFromAllowlist = removeFromAllowlist;
window.removePermission = removePermission;
window.deleteWorld = deleteWorld;
window.activateWorld = activateWorld;
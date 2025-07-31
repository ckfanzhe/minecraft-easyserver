// Internationalization (i18n) support
class I18n {
    constructor() {
        this.currentLanguage = localStorage.getItem('language') || 'zh';
        this.translations = {
            zh: {
                // Navigation
                'nav.title': 'Minecraft 控制面板',
                'nav.status.offline': '离线',
                'nav.status.online': '在线',
                'nav.status.unknown': '未知',
                'nav.status.running': '运行中',
                'nav.status.stopped': '已停止',
                
                // Server Control
                'server.control.title': '服务器控制',
                'server.control.start': '启动服务器',
                'server.control.stop': '停止服务器',
                'server.control.restart': '重启服务器',
                
                // Server Configuration
                'config.title': '服务器配置',
                'config.server-name': '服务器名称',
                'config.gamemode': '游戏模式',
                'config.gamemode.survival': '生存模式',
                'config.gamemode.creative': '创造模式',
                'config.gamemode.adventure': '冒险模式',
                'config.difficulty': '难度',
                'config.difficulty.peaceful': '和平',
                'config.difficulty.easy': '简单',
                'config.difficulty.normal': '普通',
                'config.difficulty.hard': '困难',
                'config.max-players': '最大玩家数',
                'config.server-port': '服务器端口',
                'config.allow-cheats': '允许作弊',
                'config.allow-list': '启用白名单',
                'config.save': '保存配置',
                
                // Allowlist Management
                'allowlist.title': '白名单管理',
                'allowlist.placeholder': '输入玩家名称',
                'allowlist.add': '添加',
                'allowlist.remove': '移除',
                'allowlist.empty': '暂无白名单用户',
                'allowlist.error.empty-name': '请输入玩家名称',
                
                // Permission Management
                'permission.title': '权限管理',
                'permission.placeholder': '输入玩家名称',
                'permission.add': '添加权限',
                'permission.empty': '暂无权限设置',
                'permission.error.empty-name': '请输入玩家名称',
                'permission.modal.title': '选择权限级别',
                'permission.modal.description': '为玩家',
                'permission.modal.description2': '选择权限级别：',
                'permission.level.visitor': '访客',
                'permission.level.visitor.desc': '只能查看，无法修改',
                'permission.level.member': '成员',
                'permission.level.member.desc': '可以进行基本操作',
                'permission.level.operator': '管理员',
                'permission.level.operator.desc': '拥有完全管理权限',
                'permission.modal.cancel': '取消',
                'permission.remove': '移除',
                
                // World Management
                'world.title': '世界管理',
                'world.upload': '上传世界文件',
                'world.upload.desc': '支持 .zip 和 .mcworld 格式，自动解压并删除压缩包',
                'world.upload.note': '上传后将自动解压到世界目录，原压缩文件会被删除',
                'world.activate': '激活',
                'world.delete': '删除',
                'world.current': '当前世界',
                'world.no-worlds': '暂无世界文件',
                'world.upload-failed': '上传失败',
                'world.delete-confirm': '确定要删除世界 "{worldName}" 吗？此操作不可撤销！',
                
                // Messages
                'message.request-failed': '请求失败',
                'message.config-saved': '配置已保存',
                'message.player-added': '玩家已添加',
                'message.player-removed': '玩家已移除',
                'message.permission-updated': '权限已更新',
                'message.world-uploaded': '世界已上传',
                'message.world-activated': '世界已激活',
                'message.world-deleted': '世界已删除',
                
                // Language
                'language.switch': '切换语言',
                'language.chinese': '中文',
                'language.english': 'English'
            },
            en: {
                // Navigation
                'nav.title': 'Minecraft EasyServer',
                'nav.status.offline': 'Offline',
                'nav.status.online': 'Online',
                'nav.status.unknown': 'Unknown',
                'nav.status.running': 'Running',
                'nav.status.stopped': 'Stopped',
                
                // Server Control
                'server.control.title': 'Server Control',
                'server.control.start': 'Start Server',
                'server.control.stop': 'Stop Server',
                'server.control.restart': 'Restart Server',
                
                // Server Configuration
                'config.title': 'Server Configuration',
                'config.server-name': 'Server Name',
                'config.gamemode': 'Game Mode',
                'config.gamemode.survival': 'Survival',
                'config.gamemode.creative': 'Creative',
                'config.gamemode.adventure': 'Adventure',
                'config.difficulty': 'Difficulty',
                'config.difficulty.peaceful': 'Peaceful',
                'config.difficulty.easy': 'Easy',
                'config.difficulty.normal': 'Normal',
                'config.difficulty.hard': 'Hard',
                'config.max-players': 'Max Players',
                'config.server-port': 'Server Port',
                'config.allow-cheats': 'Allow Cheats',
                'config.allow-list': 'Enable Allowlist',
                'config.save': 'Save Configuration',
                
                // Allowlist Management
                'allowlist.title': 'Allowlist Management',
                'allowlist.placeholder': 'Enter player name',
                'allowlist.add': 'Add',
                'allowlist.remove': 'Remove',
                'allowlist.empty': 'No allowlist users',
                'allowlist.error.empty-name': 'Please enter player name',
                
                // Permission Management
                'permission.title': 'Permission Management',
                'permission.placeholder': 'Enter player name',
                'permission.add': 'Add Permission',
                'permission.empty': 'No permission settings',
                'permission.error.empty-name': 'Please enter player name',
                'permission.modal.title': 'Select Permission Level',
                'permission.modal.description': 'For player',
                'permission.modal.description2': 'select permission level:',
                'permission.level.visitor': 'Visitor',
                'permission.level.visitor.desc': 'View only, cannot modify',
                'permission.level.member': 'Member',
                'permission.level.member.desc': 'Can perform basic operations',
                'permission.level.operator': 'Operator',
                'permission.level.operator.desc': 'Full administrative permissions',
                'permission.modal.cancel': 'Cancel',
                'permission.remove': 'Remove',
                
                // World Management
                'world.title': 'World Management',
                'world.upload': 'Upload World File',
                'world.upload.desc': 'Supports .zip and .mcworld formats, auto-extract and delete archive',
                'world.upload.note': 'Files will be auto-extracted to worlds directory, original archive will be deleted',
                'world.activate': 'Activate',
                'world.delete': 'Delete',
                'world.current': 'Current World',
                'world.no-worlds': 'No world files',
                'world.upload-failed': 'Upload failed',
                'world.delete-confirm': 'Are you sure you want to delete world "{worldName}"? This action cannot be undone!',
                
                // Messages
                'message.request-failed': 'Request failed',
                'message.config-saved': 'Configuration saved',
                'message.player-added': 'Player added',
                'message.player-removed': 'Player removed',
                'message.permission-updated': 'Permission updated',
                'message.world-uploaded': 'World uploaded',
                'message.world-activated': 'World activated',
                'message.world-deleted': 'World deleted',
                
                // Language
                'language.switch': 'Switch Language',
                'language.chinese': '中文',
                'language.english': 'English'
            }
        };
    }

    // Get translation for a key
    t(key, params = {}) {
        const translation = this.translations[this.currentLanguage][key] || key;
        
        // Replace parameters in translation
        let result = translation;
        for (const [param, value] of Object.entries(params)) {
            result = result.replace(`{${param}}`, value);
        }
        
        return result;
    }

    // Set current language
    setLanguage(language) {
        if (this.translations[language]) {
            this.currentLanguage = language;
            localStorage.setItem('language', language);
            this.updatePageTexts();
            this.updatePageLanguage();
        }
    }

    // Get current language
    getCurrentLanguage() {
        return this.currentLanguage;
    }

    // Update all texts on the page
    updatePageTexts() {
        // Update elements with data-i18n attribute
        document.querySelectorAll('[data-i18n]').forEach(element => {
            const key = element.getAttribute('data-i18n');
            element.textContent = this.t(key);
        });

        // Update elements with data-i18n-placeholder attribute
        document.querySelectorAll('[data-i18n-placeholder]').forEach(element => {
            const key = element.getAttribute('data-i18n-placeholder');
            element.placeholder = this.t(key);
        });

        // Update page title
        document.title = this.t('nav.title') + ' 管理面板';
        
        // Update HTML lang attribute
        document.documentElement.lang = this.currentLanguage === 'zh' ? 'zh-CN' : 'en';
    }

    // Update page language attribute
    updatePageLanguage() {
        document.documentElement.lang = this.currentLanguage === 'zh' ? 'zh-CN' : 'en';
    }

    // Initialize i18n
    init() {
        this.updatePageTexts();
        this.createLanguageToggle();
    }

    // Create language toggle button
    createLanguageToggle() {
        const nav = document.querySelector('nav .container .flex');
        if (nav) {
            const rightSection = nav.querySelector('.flex.items-center.space-x-4');
            if (rightSection) {
                // Remove existing language toggle if present
                const existingToggle = rightSection.querySelector('.language-toggle');
                if (existingToggle) {
                    existingToggle.remove();
                }

                const languageToggle = document.createElement('div');
                languageToggle.className = 'language-toggle flex items-center space-x-2';
                languageToggle.innerHTML = `
                    <button id="language-btn" class="p-2 hover:bg-white hover:bg-opacity-20 rounded text-sm">
                        ${this.currentLanguage === 'zh' ? 'EN' : '中'}
                    </button>
                `;

                rightSection.insertBefore(languageToggle, rightSection.firstChild);

                // Add click event
                document.getElementById('language-btn').addEventListener('click', () => {
                    const newLanguage = this.currentLanguage === 'zh' ? 'en' : 'zh';
                    this.setLanguage(newLanguage);
                    document.getElementById('language-btn').textContent = newLanguage === 'zh' ? 'EN' : '中';
                });
            }
        }
    }
}

// Create global i18n instance
window.i18n = new I18n();
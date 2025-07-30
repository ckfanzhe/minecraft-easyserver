# Bedrock Web管理面板

一个基于 Go 语言开发的 **轻量级** Minecraft Bedrock 服务器 Web 管理面板，提供现代化的用户界面和完整的服务器管理功能。

**目前仅在Windows下做测试运行**

## 🚀 功能特性

### 🎮 服务器控制
- **一键启动/停止/重启** Minecraft Bedrock 服务器
 -**实时状态监控** 显示服务器运行状态和进程信息
- **安全的进程管理** 确保服务器进程的稳定运行

### ⚙️ 配置管理
- **支持所有主要配置项**：
  - 服务器名称和描述
  - 游戏模式（生存/创造/冒险）
  - 难度设置（和平/简单/普通/困难）
  - 最大玩家数量
  - 服务器端口配置
  - 作弊和白名单开关
- **配置验证** 确保输入的配置值有效
- **实时保存** 修改后立即保存到配置文件

### 👥 白名单管理
- **添加/删除玩家** 管理允许加入服务器的玩家列表
- **实时同步** 修改后立即更新 `allowlist.json` 文件

### 🛡️ 权限管理
- **三级权限系统**：
  - **访客 (Visitor)** - 基础游戏权限
  - **成员 (Member)** - 标准玩家权限
  - **管理员 (Operator)** - 完整管理权限
- **玩家权限设置** 为特定玩家分配权限级别
- **权限文件管理** 自动维护 `permissions.json` 文件

### 🌍 世界管理
- **世界文件上传** 支持 `.zip` 和 `.mcworld` 格式
- **世界切换** 一键激活不同的世界
- **世界删除** 安全删除不需要的世界文件
- **当前世界标识** 清晰显示正在使用的世界

## 👀 管理端预览
![管理端预览](resources/screenshot.png)

## 📋 系统要求

### 服务器环境 Windows
- **操作系统**: Windows 10 或更高版本
- **Go 语言**: 1.21 或更高版本
- **内存**: 至少 2GB RAM
- **存储**: 至少 2GB 可用空间
- **网络**: 开放端口 8080（管理面板）和 19132（Minecraft 服务器）

### Minecraft Bedrock 服务器
- 已下载并解压的 Minecraft Bedrock Dedicated Server
- 服务器文件应放置在 `./bedrock-server/bedrock-server-1.21.95.1/` 目录下

## 🛠️ 安装指南

### 1. 环境准备

#### 安装 Go 语言
1. 访问 [Go 官网](https://golang.org/dl/) 下载 Windows 版本
2. 运行安装程序并按照提示完成安装
3. 验证安装：
   ```powershell
   go version
   ```

#### 下载 Minecraft Bedrock 服务器
1. 访问 [Minecraft 官网](https://www.minecraft.net/en-us/download/server/bedrock)
2. 下载 Bedrock Dedicated Server
3. 解压到项目目录的 `bedrock-server` 文件夹中

### 2. 项目部署

#### 克隆或下载项目
```powershell
# 如果使用 Git
git clone <repository-url>
cd bedrock-easyserver

# 或者直接下载并解压项目文件
```

#### 安装依赖
```powershell
go mod tidy
```

#### 构建项目
```powershell
# 构建可执行文件
go build -o bedrock-easyserver.exe

# 或者直接运行
go run main.go
```

### 3. 目录结构确认

确保你的项目目录结构如下：
```
bedrock-easyserver/
├── main.go                    # 主程序文件
├── go.mod                     # Go 模块文件
├── go.sum                     # Go 依赖校验文件
├── config.yml                 # 应用配置文件
├── readme.md                  # 项目说明文档
├── .gitignore                 # Git 忽略文件配置
├── config/                    # 配置模块
│   └── config.go             # 配置处理逻辑
├── handlers/                  # HTTP 处理器
│   ├── handlers.go           # API 路由处理
│   └── handlers_test.go      # 处理器单元测试
├── models/                    # 数据模型
│   └── models.go             # 数据结构定义
├── services/                  # 业务逻辑服务
│   ├── services.go           # 核心业务逻辑
│   └── services_test.go      # 服务层单元测试
├── web/                       # 前端文件
│   ├── index.html            # 主页面
│   └── app.js                # JavaScript 逻辑
└── bedrock-server/           # Bedrock 服务器目录
    └── bedrock-server-1.21.95.1/
        ├── bedrock_server.exe
        ├── server.properties
        ├── allowlist.json
        ├── permissions.json
        └── worlds/
```

## 🚀 使用指南

### 启动管理面板

1. **命令行启动**：
   ```powershell
   # 方式一：直接运行源码
   go run main.go
   
   # 方式二：运行编译后的程序
   ./bedrock-easyserver.exe
   ```

2. **访问管理界面**：
   - 打开浏览器访问：`http://localhost:8080`
   - 管理面板将自动加载

### 防火墙配置
确保以下端口在防火墙中开放：
- **8080**: 管理面板访问端口
- **19132**: Minecraft Bedrock 服务器默认端口
- **19133**: Minecraft Bedrock 服务器 IPv6 端口

## 其他

### TODO计划功能
- 🔄 支持一键导入mcpackage模组
- 🔄 支持Linux操作系统
- 🔄 bedrock服务器日志实时查看
- 🔄 直接通过页面执行命令到Bedrock服务器
- 🔄 玩家在线状态监控
- 🔄 服务器性能监控
- 🔄 世界自动备份功能
- 🔄 多语言界面支持

## 🤝 贡献指南

欢迎提交问题报告、功能建议和代码贡献！

### 开发环境设置
1. Fork 项目仓库
2. 创建功能分支：`git checkout -b feature/new-feature`
3. 提交更改：`git commit -am 'Add new feature'`
4. 推送分支：`git push origin feature/new-feature`
5. 创建 Pull Request

### 代码规范
- 使用 Go 标准代码格式
- 添加适当的注释和文档
- 确保代码通过测试
- 遵循项目的架构模式

## 📄 许可证

本项目采用 MIT 许可证 - 查看 [LICENSE](LICENSE) 文件了解详情。

## 🙏 致谢

- [Gin Web Framework](https://gin-gonic.com/) - 高性能的 Go Web 框架
- [Tailwind CSS](https://tailwindcss.com/) - 实用优先的 CSS 框架
- [Font Awesome](https://fontawesome.com/) - 图标库
- [Minecraft Bedrock](https://www.minecraft.net/) - 游戏服务器

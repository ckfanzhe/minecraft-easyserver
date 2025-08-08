# Minecraft Easy Server API 文档

## 概述

Minecraft Easy Server 提供了一套完整的 RESTful API 来管理 Minecraft Bedrock 服务器。本文档详细描述了所有可用的 API 端点、请求参数和响应格式。

## 基础信息

- **基础 URL**: `http://localhost:8080/api`
- **内容类型**: `application/json`
- **WebSocket 端点**: `ws://localhost:8080/api/logs/ws`
- **认证方式**: JWT Bearer Token

## 认证

除了登录接口外，所有 API 端点都需要 JWT 认证。请在请求头中包含有效的 Bearer Token：

```http
Authorization: Bearer <your-jwt-token>
```

## API 端点

### 0. 认证

#### 0.1 用户登录

```http
POST /api/auth/login
```

**请求体**:
```json
{
  "password": "your-password"
}
```

**响应示例**:
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "message": "Login successful",
  "requirePasswordChange": false
}
```

**错误响应**:
```json
{
  "error": "invalid password"
}
```

**说明**:
- 密码在服务器配置文件 `config/config.yml` 中的 `auth.password` 字段设置
- 返回的 token 有效期为 24 小时
- 后续请求需要在 Authorization 头中携带此 token
- `requirePasswordChange` 字段表示是否需要强制修改密码（当使用默认密码 `admin123` 时为 `true`）

#### 0.2 修改密码

```http
POST /api/auth/change-password
```

**请求头**:
```http
Authorization: Bearer <your-jwt-token>
```

**请求体**:
```json
{
  "current_password": "admin123",
  "new_password": "NewSecure123!"
}
```

**响应示例**:
```json
{
  "message": "密码修改成功",
  "success": true
}
```

**错误响应**:
```json
{
  "message": "当前密码不正确",
  "success": false
}
```

**说明**:
- 需要提供当前密码进行验证
- 新密码必须满足强度要求：至少8位，包含大小写字母、数字和特殊字符
- 修改成功后需要重新登录

### 1. 服务器控制

#### 1.1 获取服务器状态

```http
GET /api/status
```

**响应示例**:
```json
{
  "status": "running",
  "message": "Server is running",
  "pid": 12345
}
```

#### 1.2 启动服务器

```http
POST /api/start
```

**响应示例**:
```json
{
  "status": "starting",
  "message": "Server is starting"
}
```

#### 1.3 停止服务器

```http
POST /api/stop
```

**响应示例**:
```json
{
  "status": "stopped",
  "message": "Server stopped successfully"
}
```

#### 1.4 重启服务器

```http
POST /api/restart
```

**响应示例**:
```json
{
  "status": "restarting",
  "message": "Server is restarting"
}
```

### 2. 服务器配置

#### 2.1 获取服务器配置

```http
GET /api/config
```

**响应示例**:
```json
{
  "server-name": "My Minecraft Server",
  "gamemode": "survival",
  "difficulty": "normal",
  "max-players": 10,
  "server-port": 19132,
  "allow-cheats": false,
  "allow-list": true,
  "online-mode": true,
  "level-name": "Bedrock level",
  "default-player-permission-level": "member"
}
```

#### 2.2 更新服务器配置

```http
PUT /api/config
```

**请求体**:
```json
{
  "server-name": "My Minecraft Server",
  "gamemode": "survival",
  "difficulty": "normal",
  "max-players": 10,
  "server-port": 19132,
  "allow-cheats": false,
  "allow-list": true,
  "online-mode": true,
  "level-name": "Bedrock level",
  "default-player-permission-level": "member"
}
```

**响应示例**:
```json
{
  "message": "Configuration updated successfully"
}
```

### 3. 白名单管理

#### 3.1 获取白名单

```http
GET /api/allowlist
```

**响应示例**:
```json
{
  "allowlist": [
    {
      "name": "player1",
      "ignoresPlayerLimit": false
    },
    {
      "name": "admin",
      "ignoresPlayerLimit": true
    }
  ]
}
```

#### 3.2 添加白名单条目

```http
POST /api/allowlist
```

**请求体**:
```json
{
  "name": "player1",
  "ignoresPlayerLimit": false
}
```

**响应示例**:
```json
{
  "message": "Player added to allowlist successfully"
}
```

#### 3.3 删除白名单条目

```http
DELETE /api/allowlist/{name}
```

**响应示例**:
```json
{
  "message": "Player removed from allowlist successfully"
}
```

### 4. 权限管理

#### 4.1 获取权限列表

```http
GET /api/permissions
```

**响应示例**:
```json
{
  "permissions": [
    {
      "xuid": "2535428692891648",
      "level": "operator"
    }
  ]
}
```

#### 4.2 更新用户权限

```http
PUT /api/permissions
```

**请求体**:
```json
{
  "xuid": "2535428692891648",
  "level": "operator"
}
```

**响应示例**:
```json
{
  "message": "Permission updated successfully"
}
```

#### 4.3 删除用户权限

```http
DELETE /api/permissions/{xuid}
```

**响应示例**:
```json
{
  "message": "Permission removed successfully"
}
```

### 5. 世界管理

#### 5.1 获取世界列表

```http
GET /api/worlds
```

**响应示例**:
```json
{
  "worlds": [
    {
      "name": "Bedrock level",
      "active": true
    },
    {
      "name": "world2",
      "active": false
    }
  ]
}
```

#### 5.2 上传世界文件

```http
POST /api/worlds/upload
```

**请求**: `multipart/form-data`
- `world`: 世界文件 (ZIP 格式)

**响应示例**:
```json
{
  "message": "World uploaded successfully"
}
```

#### 5.3 删除世界

```http
DELETE /api/worlds/{name}
```

**响应示例**:
```json
{
  "message": "World deleted successfully"
}
```

#### 5.4 激活世界

```http
PUT /api/worlds/{name}/activate
```

**响应示例**:
```json
{
  "message": "World activated successfully"
}
```

### 6. 资源包管理

#### 6.1 获取资源包列表

```http
GET /api/resource-packs
```

**响应示例**:
```json
{
  "resource_packs": [
    {
      "name": "My Resource Pack",
      "uuid": "12345678-1234-1234-1234-123456789012",
      "version": [1, 0, 0],
      "description": "A custom resource pack",
      "folder_name": "my_resource_pack",
      "active": true
    }
  ]
}
```

#### 6.2 上传资源包

```http
POST /api/resource-packs/upload
```

**请求**: `multipart/form-data`
- `resource_pack`: 资源包文件 (ZIP 或 MCPACK 格式)

**响应示例**:
```json
{
  "message": "Resource pack uploaded successfully"
}
```

#### 6.3 激活资源包

```http
PUT /api/resource-packs/{uuid}/activate
```

**响应示例**:
```json
{
  "message": "Resource pack activated successfully"
}
```

#### 6.4 停用资源包

```http
PUT /api/resource-packs/{uuid}/deactivate
```

**响应示例**:
```json
{
  "message": "Resource pack deactivated successfully"
}
```

#### 6.5 删除资源包

```http
DELETE /api/resource-packs/{uuid}
```

**响应示例**:
```json
{
  "message": "Resource pack deleted successfully"
}
```

### 7. 服务器版本管理

#### 7.1 获取可用版本列表

```http
GET /api/server-versions
```

**响应示例**:
```json
{
  "success": true,
  "data": [
    {
      "version": "1.20.1.02",
      "download_url": "https://example.com/bedrock-server-1.20.1.02.zip",
      "active": true,
      "downloaded": true,
      "path": "/path/to/server",
      "release_date": "2023-06-07",
      "description": "Minecraft Bedrock Server 1.20.1.02"
    }
  ]
}
```

#### 7.2 下载服务器版本

```http
POST /api/server-versions/{version}/download
```

**响应示例**:
```json
{
  "success": true,
  "message": "Download started successfully"
}
```

#### 7.3 获取下载进度

```http
GET /api/server-versions/{version}/progress
```

**响应示例**:
```json
{
  "success": true,
  "data": {
    "version": "1.20.1.02",
    "progress": 75.5,
    "status": "downloading",
    "message": "Downloading...",
    "total_bytes": 104857600,
    "downloaded_bytes": 79165440
  }
}
```

#### 7.4 激活服务器版本

```http
POST /api/server-versions/{version}/activate
```

**响应示例**:
```json
{
  "success": true,
  "message": "Version activated successfully"
}
```

#### 7.5 更新版本配置

```http
POST /api/server-versions/update-config
```

**请求体**:
```json
{
  "description": "Updated description"
}
```

**响应示例**:
```json
{
  "success": true,
  "message": "Version configuration updated successfully"
}
```

### 8. 日志管理

#### 8.1 获取服务器日志

```http
GET /api/logs?limit=100
```

**查询参数**:
- `limit`: 返回的日志条数 (默认: 100)

**响应示例**:
```json
{
  "logs": [
    {
      "timestamp": "2023-06-07T10:30:00Z",
      "level": "INFO",
      "message": "Server started successfully"
    }
  ],
  "count": 1
}
```

#### 8.2 清空日志

```http
DELETE /api/logs
```

**响应示例**:
```json
{
  "message": "Logs cleared successfully"
}
```

#### 8.3 WebSocket 实时日志

```
ws://localhost:8080/api/logs/ws
```

### 9. 服务器交互

#### 9.1 获取交互状态

```http
GET /api/interaction/status
```

**响应示例**:
```json
{
  "enabled": true,
  "platform": "linux"
}
```

#### 9.2 发送命令

```http
POST /api/interaction/command
```

**请求体**:
```json
{
  "command": "say Hello World!",
  "timestamp": "2023-06-07T10:30:00Z"
}
```

**响应示例**:
```json
{
  "command": "say Hello World!",
  "response": "Command executed successfully",
  "timestamp": "2023-06-07T10:30:00Z",
  "success": true
}
```

#### 9.3 获取命令历史

```http
GET /api/interaction/history?limit=50
```

**查询参数**:
- `limit`: 返回的历史记录条数 (默认: 50)

**响应示例**:
```json
{
  "history": [
    {
      "command": "say Hello World!",
      "response": "Command executed successfully",
      "timestamp": "2023-06-07T10:30:00Z",
      "success": true
    }
  ],
  "count": 1
}
```

#### 9.4 清空命令历史

```http
DELETE /api/interaction/history
```

**响应示例**:
```json
{
  "message": "Command history cleared successfully"
}
```

### 10. 快捷命令

#### 10.1 获取快捷命令列表

```http
GET /api/commands?category=admin
```

**查询参数**:
- `category`: 命令分类 (可选)

**响应示例**:
```json
{
  "commands": [
    {
      "id": "cmd_001",
      "name": "重载配置",
      "description": "重新加载服务器配置",
      "command": "reload",
      "category": "admin"
    }
  ],
  "count": 1
}
```

#### 10.2 获取命令分类

```http
GET /api/commands/categories
```

**响应示例**:
```json
{
  "categories": ["admin", "player", "world"],
  "count": 3
}
```

#### 10.3 执行快捷命令

```http
POST /api/commands/{id}/execute
```

**响应示例**:
```json
{
  "command": "reload",
  "response": "Configuration reloaded successfully",
  "timestamp": "2023-06-07T10:30:00Z",
  "success": true
}
```

#### 10.4 添加快捷命令

```http
POST /api/commands
```

**请求体**:
```json
{
  "name": "自定义命令",
  "description": "这是一个自定义命令",
  "command": "custom command",
  "category": "custom"
}
```

**响应示例**:
```json
{
  "message": "Quick command added successfully",
  "id": "cmd_002"
}
```

#### 10.5 删除快捷命令

```http
DELETE /api/commands/{id}
```

**响应示例**:
```json
{
  "message": "Quick command removed successfully"
}
```

### 11. 性能监控

#### 11.1 获取性能监控数据

```http
GET /api/monitor/performance
```

**响应示例**:
```json
{
  "system": {
    "cpu_usage": 25.5,
    "memory_usage": 60.2,
    "timestamp": "2023-06-07T10:30:00Z"
  },
  "bedrock": {
    "pid": 12345,
    "cpu_usage": 15.3,
    "memory_usage": 45.8,
    "memory_mb": 512.5,
    "timestamp": "2023-06-07T10:30:00Z"
  }
}
```

## 数据模型

### ServerConfig
```json
{
  "server-name": "string",
  "gamemode": "string",
  "difficulty": "string",
  "max-players": "integer",
  "server-port": "integer",
  "allow-cheats": "boolean",
  "allow-list": "boolean",
  "online-mode": "boolean",
  "level-name": "string",
  "default-player-permission-level": "string"
}
```

### AllowlistEntry
```json
{
  "name": "string",
  "ignoresPlayerLimit": "boolean"
}
```

### PermissionEntry
```json
{
  "xuid": "string",
  "level": "string"
}
```

### WorldInfo
```json
{
  "name": "string",
  "active": "boolean"
}
```

### ServerStatus
```json
{
  "status": "string",
  "message": "string",
  "pid": "integer"
}
```

### ResourcePackInfo
```json
{
  "name": "string",
  "uuid": "string",
  "version": ["integer", "integer", "integer"],
  "description": "string",
  "folder_name": "string",
  "active": "boolean"
}
```

### ServerVersion
```json
{
  "version": "string",
  "download_url": "string",
  "active": "boolean",
  "downloaded": "boolean",
  "path": "string",
  "release_date": "string",
  "description": "string"
}
```

### DownloadProgress
```json
{
  "version": "string",
  "progress": "number",
  "status": "string",
  "message": "string",
  "total_bytes": "integer",
  "downloaded_bytes": "integer"
}
```

### ServerLogEntry
```json
{
  "timestamp": "string",
  "level": "string",
  "message": "string"
}
```

### ServerCommand
```json
{
  "command": "string",
  "timestamp": "string"
}
```

### ServerCommandResponse
```json
{
  "command": "string",
  "response": "string",
  "timestamp": "string",
  "success": "boolean"
}
```

### QuickCommand
```json
{
  "id": "string",
  "name": "string",
  "description": "string",
  "command": "string",
  "category": "string"
}
```

### PerformanceMonitoringData
```json
{
  "system": {
    "cpu_usage": "number",
    "memory_usage": "number",
    "timestamp": "string"
  },
  "bedrock": {
    "pid": "integer",
    "cpu_usage": "number",
    "memory_usage": "number",
    "memory_mb": "number",
    "timestamp": "string"
  }
}
```

## 错误响应

所有 API 端点在发生错误时都会返回以下格式的错误响应:

```json
{
  "error": "错误描述信息"
}
```

常见的 HTTP 状态码:
- `200`: 成功
- `400`: 请求参数错误
- `404`: 资源未找到
- `500`: 服务器内部错误
- `503`: 服务不可用

## WebSocket 连接

### 实时日志
连接到 `ws://localhost:8080/api/logs/ws` 可以接收实时的服务器日志信息。

消息格式:
```json
{
  "timestamp": "2023-06-07T10:30:00Z",
  "level": "INFO",
  "message": "Server started successfully"
}
```

## 注意事项

1. 所有的文件上传接口都使用 `multipart/form-data` 格式
2. 时间戳格式为 ISO 8601 (RFC 3339)
3. 某些操作需要服务器处于特定状态才能执行
4. WebSocket 连接需要处理重连逻辑
5. 性能监控数据会定期更新，建议适当的轮询间隔
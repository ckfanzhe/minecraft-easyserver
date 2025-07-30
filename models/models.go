package models

// ServerConfig 服务器配置结构
type ServerConfig struct {
	ServerName              string `json:"server-name"`
	Gamemode                string `json:"gamemode"`
	Difficulty              string `json:"difficulty"`
	MaxPlayers              int    `json:"max-players"`
	ServerPort              int    `json:"server-port"`
	AllowCheats             bool   `json:"allow-cheats"`
	AllowList               bool   `json:"allow-list"`
	OnlineMode              bool   `json:"online-mode"`
	LevelName               string `json:"level-name"`
	DefaultPlayerPermission string `json:"default-player-permission-level"`
}

// AllowlistEntry 白名单条目
type AllowlistEntry struct {
	Name               string `json:"name"`
	IgnoresPlayerLimit bool   `json:"ignoresPlayerLimit"`
}

// PermissionEntry 权限条目
type PermissionEntry struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

// WorldInfo 世界信息
type WorldInfo struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// ServerStatus 服务器状态
type ServerStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	PID     int    `json:"pid,omitempty"`
}
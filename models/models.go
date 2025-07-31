package models

// ServerConfig server configuration structure
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

// AllowlistEntry allowlist entry
type AllowlistEntry struct {
	Name               string `json:"name"`
	IgnoresPlayerLimit bool   `json:"ignoresPlayerLimit"`
}

// PermissionEntry permission entry
type PermissionEntry struct {
	Name  string `json:"name"`
	Level string `json:"level"`
}

// WorldInfo world information
type WorldInfo struct {
	Name   string `json:"name"`
	Active bool   `json:"active"`
}

// ServerStatus server status
type ServerStatus struct {
	Status  string `json:"status"`
	Message string `json:"message"`
	PID     int    `json:"pid,omitempty"`
}
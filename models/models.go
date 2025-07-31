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

// ResourcePackManifest resource pack manifest structure
type ResourcePackManifest struct {
	FormatVersion int                      `json:"format_version"`
	Header        ResourcePackHeader       `json:"header"`
	Modules       []ResourcePackModule     `json:"modules"`
	Subpacks      []ResourcePackSubpack    `json:"subpacks,omitempty"`
	Capabilities  []string                 `json:"capabilities,omitempty"`
}

// ResourcePackHeader resource pack header information
type ResourcePackHeader struct {
	Description      string    `json:"description"`
	Name             string    `json:"name"`
	UUID             string    `json:"uuid"`
	Version          [3]int    `json:"version"`
	MinEngineVersion [3]int    `json:"min_engine_version"`
}

// ResourcePackModule resource pack module information
type ResourcePackModule struct {
	Description string `json:"description"`
	Type        string `json:"type"`
	UUID        string `json:"uuid"`
	Version     [3]int `json:"version"`
}

// ResourcePackSubpack resource pack subpack information
type ResourcePackSubpack struct {
	FolderName  string `json:"folder_name"`
	Name        string `json:"name"`
	MemoryTier  int    `json:"memory_tier"`
}

// ResourcePackInfo resource pack information for API response
type ResourcePackInfo struct {
	Name        string `json:"name"`
	UUID        string `json:"uuid"`
	Version     [3]int `json:"version"`
	Description string `json:"description"`
	FolderName  string `json:"folder_name"`
	Active      bool   `json:"active"`
}

// WorldResourcePack world resource pack configuration entry
type WorldResourcePack struct {
	PackID  string `json:"pack_id"`
	Version [3]int `json:"version"`
}
package models

// LoginRequest login request structure
type LoginRequest struct {
	Password string `json:"password" binding:"required"`
}

// LoginResponse login response structure
type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// JWTClaims JWT claims structure
type JWTClaims struct {
	Authorized bool   `json:"authorized"`
	Exp        int64  `json:"exp"`
	Iat        int64  `json:"iat"`
}

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
	Xuid  string `json:"xuid"`
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

// ServerVersion server version information
type ServerVersion struct {
	Version     string `json:"version"`
	DownloadURL string `json:"download_url"`
	Active      bool   `json:"active"`
	Downloaded  bool   `json:"downloaded"`
	Path        string `json:"path"`
	ReleaseDate string `json:"release_date,omitempty"`
	Description string `json:"description,omitempty"`
}

// ServerVersionConfig configuration structure for server versions
type ServerVersionConfig struct {
	Versions []ServerVersionInfo `json:"versions"`
}

// ServerVersionInfo version information from config file
type ServerVersionInfo struct {
	Version            string `json:"version"`
	DownloadURL        string `json:"download_url,omitempty"`        // For backward compatibility
	DownloadURLWindows string `json:"download_url_windows,omitempty"` // Windows download URL
	DownloadURLLinux   string `json:"download_url_linux,omitempty"`   // Linux download URL
	ReleaseDate        string `json:"release_date"`
	Description        string `json:"description"`
}

// DownloadProgress download progress information
type DownloadProgress struct {
	Version     string  `json:"version"`
	Progress    float64 `json:"progress"`
	Status      string  `json:"status"`
	Message     string  `json:"message"`
	TotalBytes  int64   `json:"total_bytes"`
	DownloadedBytes int64 `json:"downloaded_bytes"`
}

// ServerLogEntry server log entry
type ServerLogEntry struct {
	Timestamp string `json:"timestamp"`
	Level     string `json:"level"`
	Message   string `json:"message"`
}

// ServerCommand server command structure
type ServerCommand struct {
	Command   string `json:"command"`
	Timestamp string `json:"timestamp"`
}

// ServerCommandResponse server command response
type ServerCommandResponse struct {
	Command   string `json:"command"`
	Response  string `json:"response"`
	Timestamp string `json:"timestamp"`
	Success   bool   `json:"success"`
}

// QuickCommand quick command structure
type QuickCommand struct {
	ID          string `json:"id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Command     string `json:"command"`
	Category    string `json:"category"`
}

// SystemPerformance system performance monitoring data
type SystemPerformance struct {
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	Timestamp   string  `json:"timestamp"`
}

// ProcessPerformance process performance monitoring data
type ProcessPerformance struct {
	PID         int     `json:"pid"`
	CPUUsage    float64 `json:"cpu_usage"`
	MemoryUsage float64 `json:"memory_usage"`
	MemoryMB    float64 `json:"memory_mb"`
	Timestamp   string  `json:"timestamp"`
}

// PerformanceMonitoringData combined performance monitoring data
type PerformanceMonitoringData struct {
	System  SystemPerformance  `json:"system"`
	Bedrock ProcessPerformance `json:"bedrock"`
}
package services

import (
	"archive/zip"
	"encoding/json"
	"fmt"
	"io"
	"minecraft-easyserver/models"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"

	"gopkg.in/yaml.v2"
)

// ServerVersionService handles server version management
type ServerVersionService struct {
	downloadProgress map[string]*models.DownloadProgress
	progressMutex    sync.RWMutex
}

// NewServerVersionService creates a new server version service
func NewServerVersionService() *ServerVersionService {
	return &ServerVersionService{
		downloadProgress: make(map[string]*models.DownloadProgress),
	}
}

// GetAvailableVersions returns available server versions from config file
func (s *ServerVersionService) GetAvailableVersions() []models.ServerVersion {
	// Load versions from config file
	configPath := "./config/server_versions.json"
	versionConfig, err := s.loadVersionConfig(configPath)
	if err != nil {
		// Fallback to hardcoded versions if config file is not available
		return s.getFallbackVersions()
	}

	var versions []models.ServerVersion
	for _, versionInfo := range versionConfig.Versions {
		version := models.ServerVersion{
			Version:     versionInfo.Version,
			DownloadURL: versionInfo.DownloadURL,
			Active:      false,
			Downloaded:  false,
			Path:        fmt.Sprintf("./bedrock-server/bedrock-server-%s", versionInfo.Version),
			ReleaseDate: versionInfo.ReleaseDate,
			Description: versionInfo.Description,
		}
		versions = append(versions, version)
	}

	// Check which versions are downloaded and which is active
	for i := range versions {
		versions[i].Downloaded = s.isVersionDownloaded(versions[i].Version)
		versions[i].Active = s.isVersionActive(versions[i].Version)
	}

	return versions
}

// DownloadVersion downloads a specific server version
func (s *ServerVersionService) DownloadVersion(version string) error {
	versions := s.GetAvailableVersions()
	var targetVersion *models.ServerVersion

	for _, v := range versions {
		if v.Version == version {
			targetVersion = &v
			break
		}
	}

	if targetVersion == nil {
		return fmt.Errorf("version %s not found", version)
	}

	if targetVersion.Downloaded {
		return fmt.Errorf("version %s is already downloaded", version)
	}

	// Initialize progress tracking
	s.progressMutex.Lock()
	s.downloadProgress[version] = &models.DownloadProgress{
		Version:  version,
		Progress: 0,
		Status:   "downloading",
		Message:  "Starting download...",
	}
	s.progressMutex.Unlock()

	go s.downloadAndExtract(targetVersion)
	return nil
}

// GetDownloadProgress returns download progress for a version
func (s *ServerVersionService) GetDownloadProgress(version string) (*models.DownloadProgress, bool) {
	s.progressMutex.RLock()
	defer s.progressMutex.RUnlock()
	progress, exists := s.downloadProgress[version]
	return progress, exists
}

// ActivateVersion activates a specific server version
func (s *ServerVersionService) ActivateVersion(version string) error {
	if !s.isVersionDownloaded(version) {
		return fmt.Errorf("version %s is not downloaded", version)
	}

	versionPath := fmt.Sprintf("./bedrock-server/bedrock-server-%s", version)

	// Update config.yml
	configPath := "./config.yml"
	config, err := s.loadConfig(configPath)
	if err != nil {
		return fmt.Errorf("failed to load config: %v", err)
	}

	// Update bedrock path
	config["bedrock"].(map[interface{}]interface{})["path"] = versionPath

	// Save config
	err = s.saveConfig(configPath, config)
	if err != nil {
		return fmt.Errorf("failed to save config: %v", err)
	}

	// Update the runtime bedrock path for server service
	// Convert relative path to absolute path
	absPath, err := filepath.Abs(versionPath)
	if err != nil {
		return fmt.Errorf("failed to get absolute path: %v", err)
	}
	
	// Update the server service bedrock path
	SetBedrockPath(absPath)

	return nil
}

// downloadAndExtract downloads and extracts the server version
func (s *ServerVersionService) downloadAndExtract(version *models.ServerVersion) {
	// Create directory
	err := os.MkdirAll(version.Path, 0755)
	if err != nil {
		s.updateProgress(version.Version, 0, "error", fmt.Sprintf("Failed to create directory: %v", err), 0, 0)
		return
	}

	// Download file
	zipPath := filepath.Join("./bedrock-server", fmt.Sprintf("bedrock-server-%s.zip", version.Version))
	err = s.downloadFile(version.DownloadURL, zipPath, version.Version)
	if err != nil {
		s.updateProgress(version.Version, 0, "error", fmt.Sprintf("Download failed: %v", err), 0, 0)
		return
	}

	// Extract file
	s.updateProgress(version.Version, 90, "extracting", "Extracting files...", 0, 0)
	err = s.extractZip(zipPath, version.Path)
	if err != nil {
		s.updateProgress(version.Version, 0, "error", fmt.Sprintf("Extraction failed: %v", err), 0, 0)
		return
	}

	// Clean up zip file
	os.Remove(zipPath)

	s.updateProgress(version.Version, 100, "completed", "Download and extraction completed", 0, 0)
}

// downloadFile downloads a file with progress tracking
func (s *ServerVersionService) downloadFile(url, filepath, version string) error {
	// Create HTTP client with HTTP/1.1 to avoid HTTP/2 stream errors
	client := &http.Client{
		Timeout: 30 * time.Minute, // 30 minutes timeout for large files
		Transport: &http.Transport{
			ForceAttemptHTTP2:     false, // Force HTTP/1.1
			MaxIdleConns:          10,
			IdleConnTimeout:       30 * time.Second,
			DisableCompression:    true, // Disable compression for better progress tracking
			TLSHandshakeTimeout:   10 * time.Second,
			ResponseHeaderTimeout: 10 * time.Second,
		},
	}

	// Create request with proper headers
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Add user agent to avoid blocking
	req.Header.Set("User-Agent", "Bedrock-EasyServer/1.0")
	req.Header.Set("Accept", "*/*")

	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to download: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("bad status: %s", resp.Status)
	}

	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Get total size
	totalSize := resp.ContentLength
	s.updateProgress(version, 0, "downloading", "Downloading...", totalSize, 0)

	// Create progress reader
	reader := &progressReader{
		Reader: resp.Body,
		total:  totalSize,
		onProgress: func(downloaded int64) {
			progress := float64(downloaded) / float64(totalSize) * 80 // 80% for download, 20% for extraction
			s.updateProgress(version, progress, "downloading", fmt.Sprintf("Downloaded %d/%d bytes", downloaded, totalSize), totalSize, downloaded)
		},
	}

	_, err = io.Copy(out, reader)
	return err
}

// extractZip extracts a zip file to destination
func (s *ServerVersionService) extractZip(src, dest string) error {
	r, err := zip.OpenReader(src)
	if err != nil {
		return err
	}
	defer r.Close()

	for _, f := range r.File {
		rc, err := f.Open()
		if err != nil {
			return err
		}

		path := filepath.Join(dest, f.Name)
		if f.FileInfo().IsDir() {
			os.MkdirAll(path, f.FileInfo().Mode())
			rc.Close()
			continue
		}

		// Create directory for file
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			rc.Close()
			return err
		}

		// Create file
		outFile, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, f.FileInfo().Mode())
		if err != nil {
			rc.Close()
			return err
		}

		_, err = io.Copy(outFile, rc)
		outFile.Close()
		rc.Close()

		if err != nil {
			return err
		}
	}

	return nil
}

// isVersionDownloaded checks if a version is downloaded
func (s *ServerVersionService) isVersionDownloaded(version string) bool {
	path := fmt.Sprintf("./bedrock-server/bedrock-server-%s", version)
	executablePath := filepath.Join(path, "bedrock_server.exe")
	_, err := os.Stat(executablePath)
	return err == nil
}

// isVersionActive checks if a version is currently active
func (s *ServerVersionService) isVersionActive(version string) bool {
	config, err := s.loadConfig("./config.yml")
	if err != nil {
		return false
	}

	bedrock, ok := config["bedrock"].(map[interface{}]interface{})
	if !ok {
		return false
	}

	path, ok := bedrock["path"].(string)
	if !ok {
		return false
	}

	expectedPath := fmt.Sprintf("./bedrock-server/bedrock-server-%s", version)
	return strings.Contains(path, version) || path == expectedPath
}

// loadConfig loads configuration from YAML file
func (s *ServerVersionService) loadConfig(path string) (map[interface{}]interface{}, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var config map[interface{}]interface{}
	err = yaml.Unmarshal(data, &config)
	return config, err
}

// saveConfig saves configuration to YAML file
func (s *ServerVersionService) saveConfig(path string, config map[interface{}]interface{}) error {
	data, err := yaml.Marshal(config)
	if err != nil {
		return err
	}

	return os.WriteFile(path, data, 0644)
}

// updateProgress updates download progress
func (s *ServerVersionService) updateProgress(version string, progress float64, status, message string, totalBytes, downloadedBytes int64) {
	s.progressMutex.Lock()
	defer s.progressMutex.Unlock()

	if s.downloadProgress[version] == nil {
		s.downloadProgress[version] = &models.DownloadProgress{}
	}

	s.downloadProgress[version].Version = version
	s.downloadProgress[version].Progress = progress
	s.downloadProgress[version].Status = status
	s.downloadProgress[version].Message = message
	s.downloadProgress[version].TotalBytes = totalBytes
	s.downloadProgress[version].DownloadedBytes = downloadedBytes

	// Clean up completed downloads after 30 seconds
	if status == "completed" || status == "error" {
		go func() {
			time.Sleep(30 * time.Second)
			s.progressMutex.Lock()
			delete(s.downloadProgress, version)
			s.progressMutex.Unlock()
		}()
	}
}

// progressReader wraps an io.Reader to track progress
type progressReader struct {
	io.Reader
	total      int64
	downloaded int64
	onProgress func(int64)
}

func (pr *progressReader) Read(p []byte) (int, error) {
	n, err := pr.Reader.Read(p)
	pr.downloaded += int64(n)
	if pr.onProgress != nil {
		pr.onProgress(pr.downloaded)
	}
	return n, err
}

// loadVersionConfig loads server version configuration from JSON file
func (s *ServerVersionService) loadVersionConfig(configPath string) (*models.ServerVersionConfig, error) {
	data, err := os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file: %v", err)
	}

	var config models.ServerVersionConfig
	err = json.Unmarshal(data, &config)
	if err != nil {
		return nil, fmt.Errorf("failed to parse config file: %v", err)
	}

	return &config, nil
}

// getFallbackVersions returns hardcoded versions as fallback
func (s *ServerVersionService) getFallbackVersions() []models.ServerVersion {
	versions := []models.ServerVersion{
		{
			Version:     "1.21.94.1",
			DownloadURL: "https://www.minecraft.net/bedrockdedicatedserver/bin-win/bedrock-server-1.21.94.1.zip",
			Active:      false,
			Downloaded:  false,
			Path:        "./bedrock-server/bedrock-server-1.21.94.1",
			ReleaseDate: "2024-12-10",
			Description: "Minecraft Bedrock Server 1.21.94.1",
		},
		{
			Version:     "1.21.95.1",
			DownloadURL: "https://www.minecraft.net/bedrockdedicatedserver/bin-win/bedrock-server-1.21.95.1.zip",
			Active:      false,
			Downloaded:  false,
			Path:        "./bedrock-server/bedrock-server-1.21.95.1",
			ReleaseDate: "2024-12-17",
			Description: "Minecraft Bedrock Server 1.21.95.1",
		},
	}

	// Check which versions are downloaded and which is active
	for i := range versions {
		versions[i].Downloaded = s.isVersionDownloaded(versions[i].Version)
		versions[i].Active = s.isVersionActive(versions[i].Version)
	}

	return versions
}

// UpdateVersionConfigFromGitHub downloads the latest server_versions.json from GitHub
func (s *ServerVersionService) UpdateVersionConfigFromGitHub() error {
	// GitHub raw URL for the server_versions.json file
	// TODO: Replace with your actual GitHub repository URL
	githubURL := "https://raw.githubusercontent.com/ckfanzhe/minecraft-easy-server/refs/heads/feature/multi-version/config/server_versions.json"

	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
		Transport: &http.Transport{
			ForceAttemptHTTP2: false, // Force HTTP/1.1
		},
	}

	// Create request
	req, err := http.NewRequest("GET", githubURL, nil)
	if err != nil {
		return fmt.Errorf("failed to create request: %v", err)
	}

	// Set headers
	req.Header.Set("User-Agent", "Bedrock-EasyServer/1.0")
	req.Header.Set("Accept", "application/json")

	// Make request
	resp, err := client.Do(req)
	if err != nil {
		return fmt.Errorf("failed to fetch from GitHub: %v", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("GitHub returned status %d", resp.StatusCode)
	}

	// Read response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return fmt.Errorf("failed to read response: %v", err)
	}

	// Validate JSON format
	var config models.ServerVersionConfig
	err = json.Unmarshal(body, &config)
	if err != nil {
		return fmt.Errorf("invalid JSON format from GitHub: %v", err)
	}

	// Ensure config directory exists
	configDir := "./config"
	if err := os.MkdirAll(configDir, 0755); err != nil {
		return fmt.Errorf("failed to create config directory: %v", err)
	}

	// Write to local file
	configPath := filepath.Join(configDir, "server_versions.json")
	err = os.WriteFile(configPath, body, 0644)
	if err != nil {
		return fmt.Errorf("failed to write config file: %v", err)
	}

	return nil
}

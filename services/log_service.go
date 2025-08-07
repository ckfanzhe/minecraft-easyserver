package services

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"sync"
	"time"

	"minecraft-easyserver/models"
	"github.com/gorilla/websocket"
)

type LogService struct {
	logEntries []models.ServerLogEntry
	mutex      sync.RWMutex
	clients    map[*websocket.Conn]bool
	clientsMux sync.RWMutex
	upgrader   websocket.Upgrader
	maxLogs    int
	stopChan   chan bool
	capturing  bool
}

var logService *LogService

func NewLogService() *LogService {
	if logService == nil {
		logService = &LogService{
			logEntries: make([]models.ServerLogEntry, 0),
			clients:    make(map[*websocket.Conn]bool),
			upgrader: websocket.Upgrader{
				CheckOrigin: func(r *http.Request) bool {
					return true // Allow all origins for development
				},
			},
			maxLogs: 1000, // Keep last 1000 log entries
			stopChan: make(chan bool),
			capturing: false,
		}
	}
	return logService
}

// AddLogEntry adds a new log entry
func (ls *LogService) AddLogEntry(level, message string) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()

	entry := models.ServerLogEntry{
		Timestamp: time.Now().Format("2006-01-02 15:04:05"),
		Level:     level,
		Message:   message,
	}

	ls.logEntries = append(ls.logEntries, entry)

	// Keep only the last maxLogs entries
	if len(ls.logEntries) > ls.maxLogs {
		ls.logEntries = ls.logEntries[len(ls.logEntries)-ls.maxLogs:]
	}

	// Broadcast to all connected clients
	ls.broadcastLogEntry(entry)
}

// GetLogs returns recent log entries
func (ls *LogService) GetLogs(limit int) []models.ServerLogEntry {
	ls.mutex.RLock()
	defer ls.mutex.RUnlock()

	if limit <= 0 || limit > len(ls.logEntries) {
		return ls.logEntries
	}

	start := len(ls.logEntries) - limit
	return ls.logEntries[start:]
}

// ClearLogs clears all log entries
func (ls *LogService) ClearLogs() {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	ls.logEntries = make([]models.ServerLogEntry, 0)
}

// HandleWebSocket handles WebSocket connections for real-time logs
func (ls *LogService) HandleWebSocket(w http.ResponseWriter, r *http.Request) {
	conn, err := ls.upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Printf("WebSocket upgrade failed: %v", err)
		return
	}
	defer conn.Close()

	log.Printf("WebSocket connection established from %s", r.RemoteAddr)

	// Add client
	ls.clientsMux.Lock()
	ls.clients[conn] = true
	ls.clientsMux.Unlock()

	// Send recent logs to new client
	ls.sendRecentLogs(conn)

	// Remove client when connection closes
	defer func() {
		ls.clientsMux.Lock()
		delete(ls.clients, conn)
		ls.clientsMux.Unlock()
	}()

	// Keep connection alive
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			break
		}
	}
}

// sendRecentLogs sends recent log entries to a specific client
func (ls *LogService) sendRecentLogs(conn *websocket.Conn) {
	logs := ls.GetLogs(100) // Send last 100 logs
	for _, entry := range logs {
		data, _ := json.Marshal(entry)
		conn.WriteMessage(websocket.TextMessage, data)
	}
}

// broadcastLogEntry broadcasts a log entry to all connected clients
func (ls *LogService) broadcastLogEntry(entry models.ServerLogEntry) {
	ls.clientsMux.RLock()
	defer ls.clientsMux.RUnlock()

	data, err := json.Marshal(entry)
	if err != nil {
		return
	}

	for client := range ls.clients {
		err := client.WriteMessage(websocket.TextMessage, data)
		if err != nil {
			// Remove disconnected client
			delete(ls.clients, client)
			client.Close()
		}
	}
}

// StartLogCapture starts capturing logs from server process
// StopLogCapture stops the log capture
func (ls *LogService) StopLogCapture() {
	if ls.capturing {
		ls.capturing = false
		// Signal all capture goroutines to stop
		select {
		case ls.stopChan <- true:
		default:
		}
		select {
		case ls.stopChan <- true:
		default:
		}
	}
}

func (ls *LogService) StartLogCapture(stdout, stderr io.ReadCloser) {
	ls.mutex.Lock()
	defer ls.mutex.Unlock()
	
	// Stop any existing capture
	if ls.capturing {
		ls.StopLogCapture()
	}
	
	ls.stopChan = make(chan bool, 2) // Buffer for 2 goroutines
	ls.capturing = true
	
	if stdout != nil {
		go ls.captureOutput(stdout, "INFO")
	}
	if stderr != nil {
		go ls.captureOutput(stderr, "ERROR")
	}
}

// captureOutput captures output from a reader and adds to logs
func (ls *LogService) captureOutput(reader io.ReadCloser, level string) {
	defer reader.Close()
	scanner := bufio.NewScanner(reader)

	for scanner.Scan() {
		// Check if we should stop
		select {
		case <-ls.stopChan:
			return
		default:
		}
		
		line := scanner.Text()
		if line != "" {
			ls.AddLogEntry(level, line)
		}
	}

	if err := scanner.Err(); err != nil {
		ls.AddLogEntry("ERROR", fmt.Sprintf("Log capture error: %v", err))
	}
}
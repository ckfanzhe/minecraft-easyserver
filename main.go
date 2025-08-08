package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"minecraft-easyserver/config"
	"minecraft-easyserver/routes"
	"minecraft-easyserver/services"
	"strings"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

//go:embed web/dist/*
var webFS embed.FS

func main() {	
	// Load configuration
	if err := config.LoadConfig(config.DefaultConfigPath); err != nil {
		log.Fatalf("Failed to load configuration file 'config/config.yml': %v\nPlease ensure the config/config.yml file exists and is properly formatted.", err)
	}
	log.Printf("Configuration loaded successfully. Server will run on port %d, debug mode: %v", 
		config.AppConfig.Server.Port, config.AppConfig.Server.Debug)

	// Set Gin mode based on configuration
	if !config.AppConfig.Server.Debug {
		log.Println("Setting Gin to release mode")
		gin.SetMode(gin.ReleaseMode)
	} else {
		log.Println("Running in debug mode")
	}

	// Initialize bedrock path (don't fail if path doesn't exist yet)
	// Users can now select and activate versions through the web interface
	bedrockPath := config.AppConfig.GetBedrockPath()
	if bedrockPath != "" {
		log.Printf("Configured bedrock path: %s", bedrockPath)
		// Try to initialize, but don't fail if it doesn't exist
		if err := services.InitBedrockPath(bedrockPath); err != nil {
			log.Printf("Warning: Bedrock path initialization failed (%s). Users can download and activate versions through the web interface. Error: %v", bedrockPath, err)
			// Set an empty path to indicate no version is currently active
			services.SetBedrockPath("")
		} else {
			log.Printf("Bedrock path initialized successfully: %s", bedrockPath)
		}
	} else {
		log.Printf("No bedrock path configured. Users can download and activate versions through the web interface.")
		services.SetBedrockPath("")
	}

	// Create Gin engine
	r := gin.Default()

	// Setup CORS middleware
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3001", "http://127.0.0.1:3001"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Accept", "Authorization", "X-Requested-With", "Upgrade", "Connection", "Sec-WebSocket-Key", "Sec-WebSocket-Version", "Sec-WebSocket-Protocol"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	// Setup embedded static files
	webSubFS, err := fs.Sub(webFS, "web/dist")
	if err != nil {
		log.Fatalf("Failed to create web sub filesystem: %v\nThis indicates an issue with embedded files compilation.", err)
	}

	// Static file service using embedded files
	// Serve all static files (including bundle.js, images, etc.)
	// r.StaticFS("/static", http.FS(webSubFS))

	// Serve bundle.js directly
	r.GET("/bundle.js", func(c *gin.Context) {
		bundleFile, err := webSubFS.Open("bundle.js")
		if err != nil {
			log.Printf("Error serving bundle.js: %v", err)
			c.String(404, "bundle.js not found")
			return
		}
		defer bundleFile.Close()
		c.DataFromReader(200, -1, "application/javascript", bundleFile, nil)
	})


	// Load HTML template from embedded files
	tmpl, err := webSubFS.Open("index.html")
	if err != nil {
		log.Fatalf("Failed to open index.html: %v\nThis indicates missing index.html in embedded files.", err)
	}
	tmpl.Close()

	// Home page - serve embedded index.html
	r.GET("/", func(c *gin.Context) {
		indexFile, err := webSubFS.Open("index.html")
		if err != nil {
			log.Printf("Error serving index.html: %v", err)
			c.String(500, "Failed to load index.html")
			return
		}
		defer indexFile.Close()
		c.DataFromReader(200, -1, "text/html; charset=utf-8", indexFile, nil)
	})

	// Setup API routes
	routes.SetupRoutes(r)

	// Catch-all route for SPA - serve index.html for all non-API routes
	// This must be placed after API routes to avoid conflicts
	r.NoRoute(func(c *gin.Context) {
		// Skip API routes
		if strings.HasPrefix(c.Request.URL.Path, "/api/") {
			c.JSON(404, gin.H{"error": "API endpoint not found"})
			return
		}
		
		// Serve index.html for all other routes (SPA routing)
		indexFile, err := webSubFS.Open("index.html")
		if err != nil {
			log.Printf("Error serving index.html for SPA route %s: %v", c.Request.URL.Path, err)
			c.String(500, "Failed to load index.html")
			return
		}
		defer indexFile.Close()
		c.DataFromReader(200, -1, "text/html; charset=utf-8", indexFile, nil)
	})

	// Start server
	serverAddr := config.AppConfig.GetServerAddress()
	log.Printf("Starting HTTP server on %s...", serverAddr)
	log.Printf("Server will be accessible at http://%s", serverAddr)
	
	port := fmt.Sprintf(":%d", config.AppConfig.Server.Port)
	
	if err := r.Run(port); err != nil {
		log.Fatalf("Failed to start server on port %s: %v\nPossible causes:\n- Port %d is already in use\n- Insufficient permissions\n- Invalid port number", port, err, config.AppConfig.Server.Port)
	}
}

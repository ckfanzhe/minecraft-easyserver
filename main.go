package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"

	"minecraft-easyserver/config"
	"minecraft-easyserver/routes"
	"minecraft-easyserver/services"

	"github.com/gin-gonic/gin"
)

//go:embed web/*
var webFS embed.FS

func main() {
	// Load configuration file
	if err := config.LoadConfig("config.yml"); err != nil {
		log.Fatal("Failed to load configuration file:", err)
	}

	// Set Gin mode based on configuration
	if !config.AppConfig.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize bedrock path (don't fail if path doesn't exist yet)
	// Users can now select and activate versions through the web interface
	bedrockPath := config.AppConfig.GetBedrockPath()
	if bedrockPath != "" {
		// Try to initialize, but don't fail if it doesn't exist
		if err := services.InitBedrockPath(bedrockPath); err != nil {
			log.Printf("Warning: Bedrock path not found (%s). Users can download and activate versions through the web interface.", err)
			// Set an empty path to indicate no version is currently active
			services.SetBedrockPath("")
		}
	} else {
		log.Printf("No bedrock path configured. Users can download and activate versions through the web interface.")
		services.SetBedrockPath("")
	}

	// Create Gin engine
	r := gin.Default()

	// Setup embedded static files
	webSubFS, err := fs.Sub(webFS, "web")
	if err != nil {
		log.Fatal("Failed to create web sub filesystem:", err)
	}

	// Static file service using embedded files
	// Create a sub filesystem for assets
	assetsFS, err := fs.Sub(webSubFS, "assets")
	if err != nil {
		log.Fatal("Failed to create assets sub filesystem:", err)
	}
	r.StaticFS("/assets", http.FS(assetsFS))

	// Load HTML template from embedded files
	tmpl, err := webSubFS.Open("index.html")
	if err != nil {
		log.Fatal("Failed to open index.html:", err)
	}
	tmpl.Close()

	// Home page - serve embedded index.html
	r.GET("/", func(c *gin.Context) {
		indexFile, err := webSubFS.Open("index.html")
		if err != nil {
			c.String(500, "Failed to load index.html")
			return
		}
		defer indexFile.Close()
		c.DataFromReader(200, -1, "text/html; charset=utf-8", indexFile, nil)
	})

	// Handle vite dev server client requests to avoid 404 errors
	r.GET("/@vite/client", func(c *gin.Context) {
		c.Status(204) // Return 204 No Content
	})

	// Setup API routes
	routes.SetupRoutes(r)

	// Start server
	serverAddr := config.AppConfig.GetServerAddress()
	log.Printf("Server started at http://%s", serverAddr)
	r.Run(":" + fmt.Sprintf("%d", config.AppConfig.Server.Port))
}

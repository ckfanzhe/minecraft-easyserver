package main

import (
	"fmt"
	"log"

	"bedrock-easyserver/config"
	"bedrock-easyserver/handlers"
	"bedrock-easyserver/services"

	"github.com/gin-gonic/gin"
)

func main() {
	// Load configuration file
	if err := config.LoadConfig("config.yml"); err != nil {
		log.Fatal("Failed to load configuration file:", err)
	}

	// Set Gin mode based on configuration
	if !config.AppConfig.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// Initialize bedrock path
	if err := services.InitBedrockPath(config.AppConfig.GetBedrockPath()); err != nil {
		log.Fatal("Failed to initialize bedrock path:", err)
	}

	// Create Gin engine
	r := gin.Default()

	// Static file service
	r.Static("/static", config.AppConfig.Web.StaticDir)
	r.Static("/webfonts", config.AppConfig.Web.WebfontsDir)
	r.LoadHTMLFiles(config.AppConfig.Web.TemplateFile)

	// Home page
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// Handle vite dev server client requests to avoid 404 errors
	r.GET("/@vite/client", func(c *gin.Context) {
		c.Status(204) // Return 204 No Content
	})

	// Create handler instances
	serverHandler := handlers.NewServerHandler()
	configHandler := handlers.NewConfigHandler()
	allowlistHandler := handlers.NewAllowlistHandler()
	permissionHandler := handlers.NewPermissionHandler()
	worldHandler := handlers.NewWorldHandler()
	resourcePackHandler := handlers.NewResourcePackHandler()

	// API routes
	api := r.Group("/api")
	{
		// Server control
		api.GET("/status", serverHandler.GetStatus)
		api.POST("/start", serverHandler.StartServer)
		api.POST("/stop", serverHandler.StopServer)
		api.POST("/restart", serverHandler.RestartServer)

		// Configuration management
		api.GET("/config", configHandler.GetConfig)
		api.PUT("/config", configHandler.UpdateConfig)

		// Allowlist management
		api.GET("/allowlist", allowlistHandler.GetAllowlist)
		api.POST("/allowlist", allowlistHandler.AddToAllowlist)
		api.DELETE("/allowlist/:name", allowlistHandler.RemoveFromAllowlist)

		// Permission management
		api.GET("/permissions", permissionHandler.GetPermissions)
		api.PUT("/permissions", permissionHandler.UpdatePermission)
		api.DELETE("/permissions/:name", permissionHandler.RemovePermission)

		// World management
		api.GET("/worlds", worldHandler.GetWorlds)
		api.POST("/worlds/upload", worldHandler.UploadWorld)
		api.DELETE("/worlds/:name", worldHandler.DeleteWorld)
		api.PUT("/worlds/:name/activate", worldHandler.ActivateWorld)

		// Resource pack management
		api.GET("/resource-packs", resourcePackHandler.GetResourcePacks)
		api.POST("/resource-packs/upload", resourcePackHandler.UploadResourcePack)
		api.PUT("/resource-packs/:uuid/activate", resourcePackHandler.ActivateResourcePack)
		api.PUT("/resource-packs/:uuid/deactivate", resourcePackHandler.DeactivateResourcePack)
		api.DELETE("/resource-packs/:uuid", resourcePackHandler.DeleteResourcePack)
	}

	// Start server
	serverAddr := config.AppConfig.GetServerAddress()
	log.Printf("Server started at http://%s", serverAddr)
	r.Run(":" + fmt.Sprintf("%d", config.AppConfig.Server.Port))
}

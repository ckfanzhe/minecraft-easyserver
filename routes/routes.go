package routes

import (
	"minecraft-easyserver/handlers"
	"minecraft-easyserver/middleware"

	"github.com/gin-gonic/gin"
)

// SetupRoutes sets up all API routes
func SetupRoutes(r *gin.Engine) {
	// Create handler instances
	serverHandler := handlers.NewServerHandler()
	configHandler := handlers.NewConfigHandler()
	allowlistHandler := handlers.NewAllowlistHandler()
	permissionHandler := handlers.NewPermissionHandler()
	worldHandler := handlers.NewWorldHandler()
	resourcePackHandler := handlers.NewResourcePackHandler()
	serverVersionHandler := handlers.NewServerVersionHandler()
	logHandler := handlers.NewLogHandler()
	interactionHandler := handlers.NewInteractionHandler()
	commandHandler := handlers.NewCommandHandler()
	performanceMonitoringHandler := handlers.NewPerformanceMonitoringHandler()
	authHandler := handlers.NewAuthHandler()

	// API routes
	api := r.Group("/api")
	{
		// Public routes (no authentication required)
		api.POST("/auth/login", authHandler.Login)
		
		// Auth routes (authentication required)
		auth := api.Group("/auth")
		auth.Use(middleware.JWTAuthMiddleware())
		{
			auth.POST("/change-password", authHandler.ChangePassword)
		}
		
		// WebSocket route with built-in authentication (must be before protected routes)
		api.GET("/websocket/logs", logHandler.HandleWebSocketWithAuth)
		
		// Protected routes (authentication required)
		protected := api.Group("/")
		protected.Use(middleware.JWTAuthMiddleware())
		{
			// Server control routes
			setupServerRoutes(protected, serverHandler)
			
			// Configuration routes
			setupConfigRoutes(protected, configHandler)
			
			// Allowlist routes
			setupAllowlistRoutes(protected, allowlistHandler)
			
			// Permission routes
			setupPermissionRoutes(protected, permissionHandler)
			
			// World routes
			setupWorldRoutes(protected, worldHandler)
			
			// Resource pack routes
			setupResourcePackRoutes(protected, resourcePackHandler)
			
			// Server version routes
			setupServerVersionRoutes(protected, serverVersionHandler)
			
			// Log routes
			setupLogRoutes(protected, logHandler)
			
			// Interaction routes
			setupInteractionRoutes(protected, interactionHandler)
			
			// Command routes
			setupCommandRoutes(protected, commandHandler)

			// Performace monitoring routes
			setupPerformanceMonitoringRoutes(protected, performanceMonitoringHandler)
		}
	}
}

// setupServerRoutes sets up server control routes
func setupServerRoutes(api *gin.RouterGroup, handler *handlers.ServerHandler) {
	api.GET("/status", handler.GetStatus)
	api.POST("/start", handler.StartServer)
	api.POST("/stop", handler.StopServer)
	api.POST("/restart", handler.RestartServer)
}

// setupConfigRoutes sets up configuration routes
func setupConfigRoutes(api *gin.RouterGroup, handler *handlers.ConfigHandler) {
	api.GET("/config", handler.GetConfig)
	api.PUT("/config", handler.UpdateConfig)
}

// setupAllowlistRoutes sets up allowlist routes
func setupAllowlistRoutes(api *gin.RouterGroup, handler *handlers.AllowlistHandler) {
	api.GET("/allowlist", handler.GetAllowlist)
	api.POST("/allowlist", handler.AddToAllowlist)
	api.DELETE("/allowlist/:name", handler.RemoveFromAllowlist)
}

// setupPermissionRoutes sets up permission routes
func setupPermissionRoutes(api *gin.RouterGroup, handler *handlers.PermissionHandler) {
	api.GET("/permissions", handler.GetPermissions)
	api.PUT("/permissions", handler.UpdatePermission)
	api.DELETE("/permissions/:xuid", handler.RemovePermission)
}

// setupWorldRoutes sets up world routes
func setupWorldRoutes(api *gin.RouterGroup, handler *handlers.WorldHandler) {
	api.GET("/worlds", handler.GetWorlds)
	api.POST("/worlds/upload", handler.UploadWorld)
	api.DELETE("/worlds/:name", handler.DeleteWorld)
	api.PUT("/worlds/:name/activate", handler.ActivateWorld)
}

// setupResourcePackRoutes sets up resource pack routes
func setupResourcePackRoutes(api *gin.RouterGroup, handler *handlers.ResourcePackHandler) {
	api.GET("/resource-packs", handler.GetResourcePacks)
	api.POST("/resource-packs/upload", handler.UploadResourcePack)
	api.PUT("/resource-packs/:uuid/activate", handler.ActivateResourcePack)
	api.PUT("/resource-packs/:uuid/deactivate", handler.DeactivateResourcePack)
	api.DELETE("/resource-packs/:uuid", handler.DeleteResourcePack)
}

// setupServerVersionRoutes sets up server version routes
func setupServerVersionRoutes(api *gin.RouterGroup, handler *handlers.ServerVersionHandler) {
	api.GET("/server-versions", handler.GetVersions)
	api.POST("/server-versions/:version/download", handler.DownloadVersion)
	api.GET("/server-versions/:version/progress", handler.GetDownloadProgress)
	api.PUT("/server-versions/:version/activate", handler.ActivateVersion)
	api.POST("/server-versions/update-config", handler.UpdateVersionConfig)
}

// setupLogRoutes sets up log routes
func setupLogRoutes(api *gin.RouterGroup, handler *handlers.LogHandler) {
	api.GET("/logs", handler.GetLogs)
	api.DELETE("/logs", handler.ClearLogs)
}

// setupInteractionRoutes sets up interaction routes
func setupInteractionRoutes(api *gin.RouterGroup, handler *handlers.InteractionHandler) {
	api.GET("/interaction/status", handler.GetStatus)
	api.POST("/interaction/command", handler.SendCommand)
	api.GET("/interaction/history", handler.GetCommandHistory)
	api.DELETE("/interaction/history", handler.ClearHistory)
}

// setupCommandRoutes sets up command routes
func setupCommandRoutes(api *gin.RouterGroup, handler *handlers.CommandHandler) {
	api.GET("/commands", handler.GetQuickCommands)
	api.GET("/commands/categories", handler.GetCategories)
	api.POST("/commands/:id/execute", handler.ExecuteQuickCommand)
	api.POST("/commands", handler.AddQuickCommand)
	api.DELETE("/commands/:id", handler.RemoveQuickCommand)
}

// setupPerformanceMonitoringRoutes sets up performance monitoring routes
func setupPerformanceMonitoringRoutes(api *gin.RouterGroup, handler *handlers.PerformanceMonitoringHandler) {
	api.GET("/monitor/performance", handler.GetPerformanceMonitoringData)
}
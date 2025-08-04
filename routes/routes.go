package routes

import (
	"minecraft-easyserver/handlers"

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

	// API routes
	api := r.Group("/api")
	{
		// Server control routes
		setupServerRoutes(api, serverHandler)
		
		// Configuration routes
		setupConfigRoutes(api, configHandler)
		
		// Allowlist routes
		setupAllowlistRoutes(api, allowlistHandler)
		
		// Permission routes
		setupPermissionRoutes(api, permissionHandler)
		
		// World routes
		setupWorldRoutes(api, worldHandler)
		
		// Resource pack routes
		setupResourcePackRoutes(api, resourcePackHandler)
		
		// Server version routes
		setupServerVersionRoutes(api, serverVersionHandler)
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
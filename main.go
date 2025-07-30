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
	// 加载配置文件
	if err := config.LoadConfig("config.yml"); err != nil {
		log.Fatal("加载配置文件失败:", err)
	}

	// 根据配置设置Gin模式
	if !config.AppConfig.Server.Debug {
		gin.SetMode(gin.ReleaseMode)
	}

	// 初始化bedrock路径
	if err := services.InitBedrockPath(config.AppConfig.GetBedrockPath()); err != nil {
		log.Fatal("初始化bedrock路径失败:", err)
	}

	// 创建Gin引擎
	r := gin.Default()

	// 静态文件服务
	r.Static("/static", config.AppConfig.Web.StaticDir)
	r.LoadHTMLFiles(config.AppConfig.Web.TemplateFile)

	// 主页
	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})

	// 处理vite开发服务器客户端请求，避免404错误
	r.GET("/@vite/client", func(c *gin.Context) {
		c.Status(204) // 返回204 No Content
	})

	// 创建处理器实例
	serverHandler := handlers.NewServerHandler()
	configHandler := handlers.NewConfigHandler()
	allowlistHandler := handlers.NewAllowlistHandler()
	permissionHandler := handlers.NewPermissionHandler()
	worldHandler := handlers.NewWorldHandler()

	// API路由
	api := r.Group("/api")
	{
		// 服务器控制
		api.GET("/status", serverHandler.GetStatus)
		api.POST("/start", serverHandler.StartServer)
		api.POST("/stop", serverHandler.StopServer)
		api.POST("/restart", serverHandler.RestartServer)

		// 配置管理
		api.GET("/config", configHandler.GetConfig)
		api.PUT("/config", configHandler.UpdateConfig)

		// 白名单管理
		api.GET("/allowlist", allowlistHandler.GetAllowlist)
		api.POST("/allowlist", allowlistHandler.AddToAllowlist)
		api.DELETE("/allowlist/:name", allowlistHandler.RemoveFromAllowlist)

		// 权限管理
		api.GET("/permissions", permissionHandler.GetPermissions)
		api.PUT("/permissions", permissionHandler.UpdatePermission)
		api.DELETE("/permissions/:name", permissionHandler.RemovePermission)

		// 世界管理
		api.GET("/worlds", worldHandler.GetWorlds)
		api.POST("/worlds/upload", worldHandler.UploadWorld)
		api.DELETE("/worlds/:name", worldHandler.DeleteWorld)
		api.PUT("/worlds/:name/activate", worldHandler.ActivateWorld)
	}

	// 启动服务器
	serverAddr := config.AppConfig.GetServerAddress()
	log.Printf("服务器启动在 http://%s", serverAddr)
	r.Run(":" + fmt.Sprintf("%d", config.AppConfig.Server.Port))
}

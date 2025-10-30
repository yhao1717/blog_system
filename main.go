package main

import (
	"blog_system/config"
	"blog_system/database"
	"blog_system/handlers"
	"blog_system/middleware"
	"flag"
	"log"

	"github.com/gin-gonic/gin"
)

func main() {
	// 解析命令行参数
	var configPath string
	flag.StringVar(&configPath, "config", "config.yaml", "Path to config file")
	flag.Parse()

	// 加载配置
	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		log.Fatalf("Failed to load configuration: %v", err)
	}

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库
	if err := database.InitDB(cfg); err != nil {
		log.Fatal("Failed to initialize database:", err)
	}

	// 初始化认证中间件
	middleware.InitAuth(cfg)

	router := gin.Default()

	// 中间件
	router.Use(middleware.LoggerMiddleware())

	// 全局中间件：添加应用信息到响应头
	router.Use(func(c *gin.Context) {
		c.Header("X-App-Name", cfg.App.Name)
		c.Header("X-App-Version", cfg.App.Version)
		c.Next()
	})

	// 健康检查端点
	router.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status":  "ok",
			"app":     cfg.App.Name,
			"version": cfg.App.Version,
			"config": gin.H{
				"driver": cfg.Database.Driver,
				"port":   cfg.Server.Port,
				"mode":   cfg.Server.Mode,
			},
		})
	})

	// 路由分组
	api := router.Group("/api/v1")
	{
		// 公开路由
		auth := api.Group("/auth")
		{
			auth.POST("/register", handlers.Register)
			auth.POST("/login", handlers.Login)
		}

		// 需要认证的路由
		authenticated := api.Group("")
		authenticated.Use(middleware.AuthMiddleware())
		{
			// 文章管理
			posts := authenticated.Group("/posts")
			{
				posts.POST("", handlers.CreatePost)
				posts.PUT("/:id", handlers.UpdatePost)
				posts.DELETE("/:id", handlers.DeletePost)
			}

			// 评论管理
			comments := authenticated.Group("/comments")
			{
				comments.POST("", handlers.CreateComment)
			}
		}

		// 公开的读取路由
		api.GET("/posts", handlers.GetPosts)
		api.GET("/posts/:id", handlers.GetPost)
		api.GET("/posts/:postId/comments", handlers.GetCommentsByPost)
	}

	log.Printf("%s v%s starting on :%s", cfg.App.Name, cfg.App.Version, cfg.Server.Port)
	if err := router.Run(":" + cfg.Server.Port); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}

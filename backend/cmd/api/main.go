package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/unifocus/backend/internal/api/handlers"
	"github.com/unifocus/backend/internal/api/middleware"
	"github.com/unifocus/backend/internal/config"
	"github.com/unifocus/backend/internal/repository/postgres"
	"github.com/unifocus/backend/internal/repository/redis"
	"github.com/unifocus/backend/internal/service"
	"github.com/unifocus/backend/pkg/jwt"
	"github.com/unifocus/backend/pkg/logger"
)

func main() {
	// 加载配置
	cfg, err := config.Load("")
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 初始化日志
	if err := logger.Init(&cfg.Log); err != nil {
		log.Fatalf("Failed to initialize logger: %v", err)
	}
	defer logger.Sync()

	logger.Info("Starting UniFocus API Server...")
	logger.Infof("Environment: %s", cfg.Server.Mode)

	// 设置Gin模式
	gin.SetMode(cfg.Server.Mode)

	// 初始化数据库连接
	db, err := postgres.NewDatabase(&cfg.Database)
	if err != nil {
		logger.Fatalf("Failed to initialize database: %v", err)
	}
	defer func() {
		if err := db.Close(); err != nil {
			logger.Errorf("Failed to close database: %v", err)
		}
	}()

	// 初始化Redis连接
	rdb, err := redis.NewClient(&cfg.Redis, "unifocus")
	if err != nil {
		logger.Fatalf("Failed to initialize redis: %v", err)
	}
	defer func() {
		if err := rdb.Close(); err != nil {
			logger.Errorf("Failed to close redis: %v", err)
		}
	}()

	// 初始化服务层
	userRepo := postgres.NewUserRepository(db)
	oppRepo := postgres.NewOpportunityRepository(db)
	profileRepo := postgres.NewProfileRepository(db)
	jwtMgr := jwt.NewManager(&cfg.JWT)
	authService := service.NewAuthService(userRepo, jwtMgr)
	oppService := service.NewOpportunityService(oppRepo)
	profileService := service.NewProfileService(profileRepo, nil) // NLP客户端待集成

	// 创建路由（传入数据库和Redis实例供后续使用）
	router := setupRouter(cfg, db, rdb, authService, oppService, profileService)

	// 创建HTTP服务器
	srv := &http.Server{
		Addr:           fmt.Sprintf(":%d", cfg.Server.Port),
		Handler:        router,
		ReadTimeout:    time.Duration(cfg.Server.ReadTimeout) * time.Second,
		WriteTimeout:   time.Duration(cfg.Server.WriteTimeout) * time.Second,
		MaxHeaderBytes: 1 << 20, // 1 MB
	}

	// 启动服务器（优雅关闭）
	go func() {
		logger.Infof("Server is running on port %d", cfg.Server.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Fatalf("Failed to start server: %v", err)
		}
	}()

	// 等待中断信号
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Info("Shutting down server...")

	// 优雅关闭（5秒超时）
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := srv.Shutdown(ctx); err != nil {
		logger.Fatalf("Server forced to shutdown: %v", err)
	}

	logger.Info("Server exited")
}

// setupRouter 设置路由
// db: PostgreSQL数据库连接实例
// rdb: Redis客户端实例
// authService: 认证服务实例
// oppService: 机会服务实例
// profileService: 用户画像服务实例
func setupRouter(cfg *config.Config, db *postgres.DB, rdb *redis.Client, authService *service.AuthService, oppService *service.OpportunityService, profileService *service.ProfileService) *gin.Engine {
	router := gin.New()

	// 中间件
	router.Use(gin.Recovery())
	router.Use(corsMiddleware())
	router.Use(middleware.MetricsMiddleware())
	router.Use(loggerMiddleware())

	// 初始化handlers
	authHandler := handlers.NewAuthHandler(authService)
	oppHandler := handlers.NewOpportunityHandler(oppService)
	profileHandler := handlers.NewProfileHandler(profileService)
	metricsHandler := handlers.NewMetricsHandler()

	// 健康检查
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"version": "1.0.0",
			"time":    time.Now().Unix(),
		})
	})

	// 监控指标
	router.GET("/api/v1/metrics", metricsHandler.GetMetrics)

	// API路由组
	v1 := router.Group("/api/v1")
	{
		// 认证路由（无需JWT）
		auth := v1.Group("/auth")
		{
			auth.POST("/register", authHandler.Register)
			auth.POST("/login", authHandler.Login)
			auth.POST("/refresh", authHandler.RefreshToken)
		}

		// 公开的机会查询路由（无需认证）
		v1.GET("/opportunities", oppHandler.List)
		v1.GET("/opportunities/:id", oppHandler.GetByID)

		// 需要认证的路由
		authorized := v1.Group("")
		authorized.Use(middleware.AuthMiddleware(authService))
		{
			// 机会管理（需要认证）
			authorized.POST("/opportunities", oppHandler.Create)
			authorized.PUT("/opportunities/:id", oppHandler.Update)
			authorized.DELETE("/opportunities/:id", oppHandler.Delete)

			// 用户画像管理
			authorized.GET("/users/me/profile", profileHandler.GetProfile)
			authorized.PUT("/users/me/profile", profileHandler.UpdateProfile)
			authorized.POST("/users/me/profile/resume", profileHandler.UploadResume)
		}
	}

	return router
}

// corsMiddleware CORS中间件
func corsMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE, PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

// loggerMiddleware 日志中间件
func loggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		start := time.Now()
		path := c.Request.URL.Path
		raw := c.Request.URL.RawQuery

		c.Next()

		latency := time.Since(start)
		clientIP := c.ClientIP()
		method := c.Request.Method
		statusCode := c.Writer.Status()

		if raw != "" {
			path = path + "?" + raw
		}

		logger.Infof("[%s] %s %s %d %v",
			method,
			path,
			clientIP,
			statusCode,
			latency,
		)
	}
}

package main

import (
	"flag"
	"fmt"
	"go-sqli-lab/src/config"
	"go-sqli-lab/src/db"
	"go-sqli-lab/src/handlers"
	"go-sqli-lab/src/logger"
	"go-sqli-lab/src/vulnerabilities"
	"os"
	"os/signal"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	configFile = flag.String("config", "config/app.yaml", "配置文件路径")
	setupDB    = flag.Bool("setup-db", false, "初始化数据库")
)

func main() {
	flag.Parse()

	// 加载配置
	cfg, err := config.Load(*configFile)
	if err != nil {
		fmt.Fprintf(os.Stderr, "加载配置失败: %v\n", err)
		os.Exit(1)
	}

	// 初始化日志
	log, err := logger.New(cfg.Log)
	if err != nil {
		fmt.Fprintf(os.Stderr, "初始化日志失败: %v\n", err)
		os.Exit(1)
	}
	defer log.Close()

	log.Info("启动 go-sqli-lab 服务器", "version", "1.0.0")

	// 初始化数据库
	database, err := db.New(cfg.Database, log)
	if err != nil {
		log.Fatal("初始化数据库失败", "error", err)
	}
	defer database.Close()

	// 如果指定了-setup-db，执行数据库初始化
	if *setupDB {
		log.Info("开始初始化数据库...")
		if err := database.Migrate(); err != nil {
			log.Fatal("数据库迁移失败", "error", err)
		}
		if err := database.Seed(); err != nil {
			log.Fatal("数据库种子数据插入失败", "error", err)
		}
		log.Info("数据库初始化完成")
		return
	}

	// 设置Gin模式
	if cfg.Server.Mode == "release" {
		gin.SetMode(gin.ReleaseMode)
	}

	// 创建Gin路由
	r := gin.New()
	r.Use(gin.Recovery())
	r.Use(logger.GinLogger(log))

	// 注册路由
	registerRoutes(r, database, log)

	// 启动服务器
	addr := fmt.Sprintf("%s:%d", cfg.Server.Host, cfg.Server.Port)
	log.Info("服务器启动", "address", addr)

	// 优雅关闭
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		if err := r.Run(addr); err != nil {
			log.Fatal("服务器启动失败", "error", err)
		}
	}()

	<-quit
	log.Info("服务器正在关闭...")
}

// registerRoutes 注册所有路由
func registerRoutes(r *gin.Engine, database db.Database, log logger.Logger) {
	// 加载HTML模板
	r.LoadHTMLGlob("./src/web/*.html")

	// 静态资源 - 从 src/web 目录提供
	r.Static("/static", "./src/web")
	r.StaticFile("/", "./src/web/index.html")

	// 健康检查
	r.GET("/health", handlers.HealthHandler())

	// API路由组
	api := r.Group("/api")
	{
		// 数据库状态
		api.GET("/db-status", handlers.DBStatusHandler(database))
		// AJAX重置数据库API
		api.POST("/reset-db", handlers.ResetDBAPIHandler(database, log))
	}

	// 注册漏洞路由
	vulnRegistry := vulnerabilities.NewRegistry(database, log)
	vulnRegistry.RegisterAll(r)

	// 设置数据库路由
	r.GET("/setup-db", handlers.SetupDBHandler(database, log))
	// 加载重置进度页面
	r.StaticFile("/loading", "./src/web/loading.html")
	// 重置成功/错误页面
	r.GET("/setup-success", handlers.SetupSuccessHandler())
	r.GET("/setup-error", handlers.SetupErrorHandler())
}

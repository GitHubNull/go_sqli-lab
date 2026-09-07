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
	"path/filepath"
	"strings"
	"syscall"

	"github.com/gin-gonic/gin"
)

var (
	configFile = flag.String("config", "config/app.yaml", "配置文件路径")
	setupDB    = flag.Bool("setup-db", false, "初始化数据库")
)

func main() {
	flag.Parse()

	// 在 chdir 前将 --config 相对路径转为绝对路径，防止 chdir 后路径解析错误
	if !filepath.IsAbs(*configFile) {
		if absPath, err := filepath.Abs(*configFile); err == nil {
			*configFile = absPath
		}
	}

	// 自动定位项目根目录并切换工作目录
	if err := ensureProjectRoot(); err != nil {
		fmt.Fprintf(os.Stderr, "错误: 无法定位项目根目录: %v\n", err)
		fmt.Fprintf(os.Stderr, "请确保从项目根目录运行程序，或项目根目录下存在 go.mod 文件\n")
		os.Exit(1)
	}

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

	log.Info("启动 go-sqli-lab 服务器", "version", "1.5.0")

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

// ensureProjectRoot 自动定位项目根目录并切换工作目录
// 通过向上查找 go.mod 和 config/app.yaml 来确定项目根目录
func ensureProjectRoot() error {
	// 如果当前目录下已有 config/app.yaml，说明已经在项目根目录
	if isProjectRoot("") {
		return nil
	}

	// 尝试通过可执行文件路径推断项目根目录
	// 可执行文件通常在 bin/ 目录下，项目根目录是其父目录
	exePath, err := os.Executable()
	if err != nil {
		return fmt.Errorf("获取可执行文件路径失败: %w", err)
	}
	// 解析符号链接，归一化路径
	exePath, err = filepath.EvalSymlinks(exePath)
	if err != nil {
		return fmt.Errorf("解析可执行文件路径失败: %w", err)
	}
	exeDir := filepath.Dir(exePath)

	// 如果可执行文件在 bin/ 目录下，项目根目录是其父目录（大小写不敏感，兼容 Windows）
	if strings.EqualFold(filepath.Base(exeDir), "bin") {
		projectRoot := filepath.Dir(exeDir)
		if isProjectRoot(projectRoot) {
			fmt.Fprintf(os.Stderr, "检测到项目根目录: %s\n", projectRoot)
			if err := os.Chdir(projectRoot); err != nil {
				return fmt.Errorf("切换到项目根目录失败: %w", err)
			}
			return nil
		}
	}

	// 向上逐级查找 go.mod 和 config/app.yaml
	dir, err := os.Getwd()
	if err != nil {
		return err
	}
	for {
		if isProjectRoot(dir) {
			fmt.Fprintf(os.Stderr, "检测到项目根目录: %s\n", dir)
			if err := os.Chdir(dir); err != nil {
				return fmt.Errorf("切换到项目根目录失败: %w", err)
			}
			return nil
		}
		parent := filepath.Dir(dir)
		if parent == dir {
			break
		}
		dir = parent
	}

	return fmt.Errorf("未找到项目根目录")
}

// isProjectRoot 检查指定目录是否为本项目的根目录
// 同时验证 go.mod 和 config/app.yaml 的存在性，防止误定位到其他 Go 项目
func isProjectRoot(dir string) bool {
	modPath := filepath.Join(dir, "go.mod")
	cfgPath := filepath.Join(dir, "config", "app.yaml")
	if dir == "" {
		modPath = "go.mod"
		cfgPath = "config/app.yaml"
	}
	_, modErr := os.Stat(modPath)
	_, cfgErr := os.Stat(cfgPath)
	return modErr == nil && cfgErr == nil
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

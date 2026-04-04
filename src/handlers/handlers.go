package handlers

import (
	"net/http"
	"time"

	"go-sqli-lab/src/db"
	"go-sqli-lab/src/logger"

	"github.com/gin-gonic/gin"
)

// HealthHandler 健康检查处理器
func HealthHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "go-sqli-lab is running",
		})
	}
}

// DBStatusHandler 数据库状态处理器
func DBStatusHandler(database db.Database) gin.HandlerFunc {
	return func(c *gin.Context) {
		dbConn := database.GetDB()
		if err := dbConn.Ping(); err != nil {
			c.JSON(http.StatusServiceUnavailable, gin.H{
				"status":  "error",
				"message": "数据库连接失败",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "ok",
			"message": "数据库连接正常",
		})
	}
}

// SetupDBHandler 数据库初始化和重置处理器（向后兼容）
func SetupDBHandler(database db.Database, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()

		// 检查是否请求重置
		isReset := c.Query("reset") == "true"

		if isReset {
			log.Info("开始重置数据库")

			// 执行重置操作
			if err := database.Reset(); err != nil {
				log.Error("数据库重置失败", "error", err)
				c.HTML(http.StatusInternalServerError, "setup_error.html", gin.H{
					"ErrorMessage": err.Error(),
					"IsReset":      true,
				})
				return
			}

			executionTime := time.Since(startTime).Round(time.Millisecond).String()
			log.Info("数据库重置成功", "duration", executionTime)

			c.HTML(http.StatusOK, "setup_success.html", gin.H{
				"Message":         "数据库重置成功",
				"ClearedTables":   16,                           // 4个基础表 + 12个挑战表
				"InsertedRecords": 16 + len(challengeRecords()), // 基础数据 + 挑战数据
				"ExecutionTime":   executionTime,
				"DatabaseType":    "SQLite/MySQL/PostgreSQL",
				"IsReset":         true,
			})
			return
		}

		// 普通初始化流程
		log.Info("开始初始化数据库")

		// 执行迁移
		if err := database.Migrate(); err != nil {
			log.Error("数据库迁移失败", "error", err)
			c.HTML(http.StatusInternalServerError, "setup_error.html", gin.H{
				"ErrorMessage": err.Error(),
				"IsReset":      false,
			})
			return
		}

		// 插入种子数据
		if err := database.Seed(); err != nil {
			log.Error("种子数据插入失败", "error", err)
			c.HTML(http.StatusInternalServerError, "setup_error.html", gin.H{
				"ErrorMessage": err.Error(),
				"IsReset":      false,
			})
			return
		}

		executionTime := time.Since(startTime).Round(time.Millisecond).String()
		log.Info("数据库初始化成功", "duration", executionTime)

		c.HTML(http.StatusOK, "setup_success.html", gin.H{
			"Message":         "数据库初始化成功",
			"ClearedTables":   0,
			"InsertedRecords": 16, // 8 users + 8 emails
			"ExecutionTime":   executionTime,
			"DatabaseType":    "SQLite/MySQL/PostgreSQL",
			"IsReset":         false,
		})
	}
}

// ResetDBAPIHandler AJAX重置数据库API
func ResetDBAPIHandler(database db.Database, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		startTime := time.Now()
		log.Info("API: 开始重置数据库")

		// 执行重置操作
		if err := database.Reset(); err != nil {
			log.Error("数据库重置失败", "error", err)
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   err.Error(),
			})
			return
		}

		executionTime := time.Since(startTime).Round(time.Millisecond).String()
		log.Info("数据库重置成功", "duration", executionTime)

		c.JSON(http.StatusOK, gin.H{
			"success":         true,
			"clearedTables":   16,
			"insertedRecords": 16 + len(challengeRecords()),
			"executionTime":   executionTime,
		})
	}
}

// SetupSuccessHandler 重置成功页面
func SetupSuccessHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "setup_success.html", gin.H{
			"Message":         "数据库重置成功",
			"ClearedTables":   c.Query("cleared"),
			"InsertedRecords": c.Query("inserted"),
			"ExecutionTime":   c.Query("time"),
			"DatabaseType":    "SQLite/MySQL/PostgreSQL",
			"IsReset":         true,
		})
	}
}

// SetupErrorHandler 重置错误页面
func SetupErrorHandler() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "setup_error.html", gin.H{
			"ErrorMessage": c.Query("message"),
			"IsReset":      true,
		})
	}
}

// challengeRecords 返回挑战关卡数据数量
func challengeRecords() []struct{} {
	// 12个挑战表，每个表有1-5条数据，平均约36条
	return make([]struct{}, 36)
}

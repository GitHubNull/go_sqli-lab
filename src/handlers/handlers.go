package handlers

import (
	"net/http"

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
		// 测试数据库连接
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

// SetupDBHandler 数据库初始化处理器
func SetupDBHandler(database db.Database, log logger.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 执行迁移
		if err := database.Migrate(); err != nil {
			log.Error("数据库迁移失败", "error", err)
			c.HTML(http.StatusInternalServerError, "setup_error.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		// 插入种子数据
		if err := database.Seed(); err != nil {
			log.Error("种子数据插入失败", "error", err)
			c.HTML(http.StatusInternalServerError, "setup_error.html", gin.H{
				"error": err.Error(),
			})
			return
		}

		c.HTML(http.StatusOK, "setup_success.html", gin.H{
			"message": "数据库初始化成功",
		})
	}
}

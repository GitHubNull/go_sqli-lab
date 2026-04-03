package vulnerabilities

import (
	"fmt"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"go-sqli-lab/src/db"
	"go-sqli-lab/src/logger"

	"github.com/gin-gonic/gin"
)

// BaseLesson 基础课程结构
type BaseLesson struct {
	db     db.Database
	logger logger.Logger
}

// logRequest 记录请求到result.txt
func (b *BaseLesson) logRequest(c *gin.Context, lesson string) {
	// 确保目录存在
	resultDir := filepath.Join(".", "tmp", lesson)
	os.MkdirAll(resultDir, 0755)

	resultFile := filepath.Join(resultDir, "result.txt")
	f, err := os.OpenFile(resultFile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return
	}
	defer f.Close()

	// 记录请求信息
	fmt.Fprintf(f, "[%s] %s %s\n", time.Now().Format("2006-01-02 15:04:05"), c.Request.Method, c.Request.URL.String())
	for key, values := range c.Request.URL.Query() {
		for _, value := range values {
			fmt.Fprintf(f, "  %s=%s\n", key, value)
		}
	}
}

// renderHTML 渲染HTML响应
func (b *BaseLesson) renderHTML(c *gin.Context, lesson, result, error string) {
	html := fmt.Sprintf(`<!DOCTYPE html>
<html>
<head>
	<meta charset="UTF-8">
	<title>%s</title>
	<link rel="stylesheet" href="/static/css/style.css">
</head>
<body>
	<div class="container">
		<h1>%s</h1>
		<div class="content">
			%s
			%s
		</div>
		<div class="nav">
			<a href="/">返回首页</a>
		</div>
	</div>
</body>
</html>`, lesson, lesson, renderResult(result), renderError(error))
	c.Header("Content-Type", "text/html; charset=utf-8")
	c.String(http.StatusOK, html)
}

func renderResult(result string) string {
	if result == "" {
		return ""
	}
	return fmt.Sprintf(`<div class="result success">%s</div>`, result)
}

func renderError(err string) string {
	if err == "" {
		return ""
	}
	return fmt.Sprintf(`<div class="error">%s</div>`, err)
}

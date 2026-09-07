package vulnerabilities

import (
	"github.com/gin-gonic/gin"
)

// Less71 Less-71: 文件上传注入 - 上传 XLSX（单元格内容可 SQL 注入）
type Less71 struct {
	BaseLesson
}

func (l *Less71) ID() string {
	return "less-71"
}

func (l *Less71) Name() string {
	return "Less-71: Upload XLSX Injection"
}

func (l *Less71) Description() string {
	return "文件上传注入 - 上传 XLSX 文件批量验证，单元格内容存在 SQL 注入"
}

func (l *Less71) Category() string {
	return "Upload Injection"
}

func (l *Less71) Route(r *gin.RouterGroup) {
	group := r.Group("/less-71")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.GET("/template", l.handleTemplate)
		group.POST("/upload", l.handleUpload)
	}
}

func (l *Less71) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-71")
	renderUploadPage(c, l.ID(), l.Name(), l.Description(), "XLSX", ".xlsx", nil, "", false)
}

func (l *Less71) handleTemplate(c *gin.Context) {
	l.logRequest(c, "less-71")
	serveUploadTemplate(c, "XLSX")
}

func (l *Less71) handleUpload(c *gin.Context) {
	handleUploadFile(c, &l.BaseLesson, l.ID(), l.Name(), l.Description(), "XLSX", ".xlsx", parseXLSX)
}

func (r *Registry) registerLess71() {
	r.Register(&Less71{BaseLesson{db: r.db, logger: r.logger}})
}

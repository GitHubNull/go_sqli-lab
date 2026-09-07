package vulnerabilities

import (
	"github.com/gin-gonic/gin"
)

// Less70 Less-70: 文件上传注入 - 上传 XLS（单元格内容可 SQL 注入）
type Less70 struct {
	BaseLesson
}

func (l *Less70) ID() string {
	return "less-70"
}

func (l *Less70) Name() string {
	return "Less-70: Upload XLS Injection"
}

func (l *Less70) Description() string {
	return "文件上传注入 - 上传 XLS 文件批量验证，单元格内容存在 SQL 注入"
}

func (l *Less70) Category() string {
	return "Upload Injection"
}

func (l *Less70) Route(r *gin.RouterGroup) {
	group := r.Group("/less-70")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.GET("/template", l.handleTemplate)
		group.POST("/upload", l.handleUpload)
	}
}

func (l *Less70) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-70")
	renderUploadPage(c, l.ID(), l.Name(), l.Description(), "XLS", ".xls", nil, "", false)
}

func (l *Less70) handleTemplate(c *gin.Context) {
	l.logRequest(c, "less-70")
	serveUploadTemplate(c, "XLS")
}

func (l *Less70) handleUpload(c *gin.Context) {
	handleUploadFile(c, &l.BaseLesson, l.ID(), l.Name(), l.Description(), "XLS", ".xls", parseXLS)
}

func (r *Registry) registerLess70() {
	r.Register(&Less70{BaseLesson{db: r.db, logger: r.logger}})
}

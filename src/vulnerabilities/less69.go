package vulnerabilities

import (
	"github.com/gin-gonic/gin"
)

// Less69 Less-69: 简历导出注入 - 导出复杂布局 XLSX（查询关键字可 SQL 注入）
type Less69 struct {
	BaseLesson
}

func (l *Less69) ID() string {
	return "less-69"
}

func (l *Less69) Name() string {
	return "Less-69: Resume Export XLSX Injection"
}

func (l *Less69) Description() string {
	return "简历导出注入 - 导出复杂布局 XLSX 文件（简历登记表），查询关键字存在 SQL 注入"
}

func (l *Less69) Category() string {
	return "Export Injection"
}

func (l *Less69) Route(r *gin.RouterGroup) {
	group := r.Group("/less-69")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.GET("/export", l.handleExport)
	}
}

func (l *Less69) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-69")

	keyword, searched := c.GetQuery("keyword")
	if !searched {
		renderResumeExportPage(c, l.ID(), l.Name(), l.Description(), "XLSX", "", nil, nil, "", "", false)
		return
	}

	cols, rows, rawSQL, err := runResumeExportQuery(l.db, keyword)
	if err != nil {
		renderResumeExportPage(c, l.ID(), l.Name(), l.Description(), "XLSX", keyword, nil, nil, rawSQL, err.Error(), true)
		return
	}
	renderResumeExportPage(c, l.ID(), l.Name(), l.Description(), "XLSX", keyword, cols, rows, rawSQL, "", true)
}

func (l *Less69) handleExport(c *gin.Context) {
	l.logRequest(c, "less-69")

	keyword := c.Query("keyword")
	cols, rows, rawSQL, err := runResumeExportQuery(l.db, keyword)
	if err != nil {
		serveExportError(c, rawSQL, err)
		return
	}

	data, buildErr := buildResumeXLSX(cols, rows)
	if buildErr != nil {
		serveExportError(c, rawSQL, buildErr)
		return
	}
	serveExport(c, "resume_export.xlsx", contentTypeXLSX, data)
}

func (r *Registry) registerLess69() {
	r.Register(&Less69{BaseLesson{db: r.db, logger: r.logger}})
}

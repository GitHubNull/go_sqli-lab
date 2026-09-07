package vulnerabilities

import (
	"github.com/gin-gonic/gin"
)

// Less68 Less-68: 简历导出注入 - 导出复杂布局 XLS（查询关键字可 SQL 注入）
type Less68 struct {
	BaseLesson
}

func (l *Less68) ID() string {
	return "less-68"
}

func (l *Less68) Name() string {
	return "Less-68: Resume Export XLS Injection"
}

func (l *Less68) Description() string {
	return "简历导出注入 - 导出复杂布局 XLS 文件（简历登记表），查询关键字存在 SQL 注入"
}

func (l *Less68) Category() string {
	return "Export Injection"
}

func (l *Less68) Route(r *gin.RouterGroup) {
	group := r.Group("/less-68")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.GET("/export", l.handleExport)
	}
}

func (l *Less68) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-68")

	keyword, searched := c.GetQuery("keyword")
	if !searched {
		renderResumeExportPage(c, l.ID(), l.Name(), l.Description(), "XLS", "", nil, nil, "", "", false)
		return
	}

	cols, rows, rawSQL, err := runResumeExportQuery(l.db, keyword)
	if err != nil {
		renderResumeExportPage(c, l.ID(), l.Name(), l.Description(), "XLS", keyword, nil, nil, rawSQL, err.Error(), true)
		return
	}
	renderResumeExportPage(c, l.ID(), l.Name(), l.Description(), "XLS", keyword, cols, rows, rawSQL, "", true)
}

func (l *Less68) handleExport(c *gin.Context) {
	l.logRequest(c, "less-68")

	keyword := c.Query("keyword")
	cols, rows, rawSQL, err := runResumeExportQuery(l.db, keyword)
	if err != nil {
		serveExportError(c, rawSQL, err)
		return
	}

	data, buildErr := buildResumeXLS(cols, rows)
	if buildErr != nil {
		serveExportError(c, rawSQL, buildErr)
		return
	}
	serveExport(c, "resume_export.xls", contentTypeXLS, data)
}

func (r *Registry) registerLess68() {
	r.Register(&Less68{BaseLesson{db: r.db, logger: r.logger}})
}

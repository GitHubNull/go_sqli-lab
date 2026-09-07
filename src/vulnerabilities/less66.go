package vulnerabilities

import (
	"github.com/gin-gonic/gin"
)

// Less66 Less-66: 文件导出注入 - 导出 XLS（查询关键字可 SQL 注入）
type Less66 struct {
	BaseLesson
}

func (l *Less66) ID() string {
	return "less-66"
}

func (l *Less66) Name() string {
	return "Less-66: Export XLS Injection"
}

func (l *Less66) Description() string {
	return "文件导出注入 - 导出 XLS 文件，查询关键字存在 SQL 注入"
}

func (l *Less66) Category() string {
	return "Export Injection"
}

func (l *Less66) Route(r *gin.RouterGroup) {
	group := r.Group("/less-66")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.GET("/export", l.handleExport)
	}
}

func (l *Less66) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-66")

	keyword, searched := c.GetQuery("keyword")
	if !searched {
		renderExportPage(c, l.ID(), l.Name(), l.Description(), "XLS", "", nil, nil, "", "", false)
		return
	}

	cols, rows, rawSQL, err := runExportQuery(l.db, keyword)
	if err != nil {
		renderExportPage(c, l.ID(), l.Name(), l.Description(), "XLS", keyword, nil, nil, rawSQL, err.Error(), true)
		return
	}
	renderExportPage(c, l.ID(), l.Name(), l.Description(), "XLS", keyword, cols, rows, rawSQL, "", true)
}

func (l *Less66) handleExport(c *gin.Context) {
	l.logRequest(c, "less-66")

	keyword := c.Query("keyword")
	cols, rows, rawSQL, err := runExportQuery(l.db, keyword)
	if err != nil {
		serveExportError(c, rawSQL, err)
		return
	}

	serveExport(c, "users_export.xls", contentTypeXLS, buildXLS(cols, rows))
}

func (r *Registry) registerLess66() {
	r.Register(&Less66{BaseLesson{db: r.db, logger: r.logger}})
}

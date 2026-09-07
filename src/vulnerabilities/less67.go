package vulnerabilities

import (
	"github.com/gin-gonic/gin"
)

// Less67 Less-67: 文件导出注入 - 导出 XLSX（查询关键字可 SQL 注入）
type Less67 struct {
	BaseLesson
}

func (l *Less67) ID() string {
	return "less-67"
}

func (l *Less67) Name() string {
	return "Less-67: Export XLSX Injection"
}

func (l *Less67) Description() string {
	return "文件导出注入 - 导出 XLSX 文件，查询关键字存在 SQL 注入"
}

func (l *Less67) Category() string {
	return "Export Injection"
}

func (l *Less67) Route(r *gin.RouterGroup) {
	group := r.Group("/less-67")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
		group.GET("/export", l.handleExport)
	}
}

func (l *Less67) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-67")

	keyword, searched := c.GetQuery("keyword")
	if !searched {
		renderExportPage(c, l.ID(), l.Name(), l.Description(), "XLSX", "", nil, nil, "", "", false)
		return
	}

	cols, rows, rawSQL, err := runExportQuery(l.db, keyword)
	if err != nil {
		renderExportPage(c, l.ID(), l.Name(), l.Description(), "XLSX", keyword, nil, nil, rawSQL, err.Error(), true)
		return
	}
	renderExportPage(c, l.ID(), l.Name(), l.Description(), "XLSX", keyword, cols, rows, rawSQL, "", true)
}

func (l *Less67) handleExport(c *gin.Context) {
	l.logRequest(c, "less-67")

	keyword := c.Query("keyword")
	cols, rows, rawSQL, err := runExportQuery(l.db, keyword)
	if err != nil {
		serveExportError(c, rawSQL, err)
		return
	}

	data, buildErr := buildXLSX(cols, rows)
	if buildErr != nil {
		serveExportError(c, rawSQL, buildErr)
		return
	}
	serveExport(c, "users_export.xlsx", contentTypeXLSX, data)
}

func (r *Registry) registerLess67() {
	r.Register(&Less67{BaseLesson{db: r.db, logger: r.logger}})
}

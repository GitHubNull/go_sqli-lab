package vulnerabilities

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Less49 Less-49: ORDER BY Error
type Less49 struct {
	BaseLesson
}

func (l *Less49) ID() string {
	return "less-49"
}

func (l *Less49) Name() string {
	return "Less-49: ORDER BY Error"
}

func (l *Less49) Description() string {
	return "ORDER BY错误注入"
}

func (l *Less49) Category() string {
	return "Advanced"
}

func (l *Less49) Route(r *gin.RouterGroup) {
	group := r.Group("/less-49")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less49) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-49")

	sort := c.Query("sort")
	if sort == "" {
		sort = "id"
	}

	query := fmt.Sprintf("SELECT * FROM users ORDER BY %s", sort)

	rows, err := l.db.Query(query)
	if err != nil {
		l.renderHTML(c, "less-49", "", err.Error())
		return
	}
	defer rows.Close()

	result := "You are in..........."
	l.renderHTML(c, "less-49", result, "")
}

func (r *Registry) registerLess49() {
	r.Register(&Less49{BaseLesson{db: r.db, logger: r.logger}})
}

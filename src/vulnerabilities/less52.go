package vulnerabilities

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Less52 Less-52: ORDER BY Stacked
type Less52 struct {
	BaseLesson
}

func (l *Less52) ID() string {
	return "less-52"
}

func (l *Less52) Name() string {
	return "Less-52: ORDER BY Stacked"
}

func (l *Less52) Description() string {
	return "ORDER BY堆叠查询注入"
}

func (l *Less52) Category() string {
	return "Advanced"
}

func (l *Less52) Route(r *gin.RouterGroup) {
	group := r.Group("/less-52")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less52) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-52")

	sort := c.Query("sort")
	if sort == "" {
		sort = "id"
	}

	query := fmt.Sprintf("SELECT * FROM users ORDER BY %s", sort)

	rows, err := l.db.Query(query)
	if err != nil {
		l.renderHTML(c, "less-52", "", "")
		return
	}
	defer rows.Close()

	result := "You are in..........."
	l.renderHTML(c, "less-52", result, "")
}

func (r *Registry) registerLess52() {
	r.Register(&Less52{BaseLesson{db: r.db, logger: r.logger}})
}

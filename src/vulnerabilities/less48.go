package vulnerabilities

import (
	"fmt"

	"github.com/gin-gonic/gin"
)

// Less48 Less-48: ORDER BY Blind
type Less48 struct {
	BaseLesson
}

func (l *Less48) ID() string {
	return "less-48"
}

func (l *Less48) Name() string {
	return "Less-48: ORDER BY Blind"
}

func (l *Less48) Description() string {
	return "ORDER BY盲注"
}

func (l *Less48) Category() string {
	return "Advanced"
}

func (l *Less48) Route(r *gin.RouterGroup) {
	group := r.Group("/less-48")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less48) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-48")

	sort := c.Query("sort")
	if sort == "" {
		sort = "id"
	}

	query := fmt.Sprintf("SELECT * FROM users ORDER BY %s", sort)

	rows, err := l.db.Query(query)
	if err != nil {
		l.renderHTML(c, "less-48", "", "")
		return
	}
	defer rows.Close()

	result := "You are in..........."
	l.renderHTML(c, "less-48", result, "")
}

func (r *Registry) registerLess48() {
	r.Register(&Less48{BaseLesson{db: r.db, logger: r.logger}})
}

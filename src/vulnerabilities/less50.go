package vulnerabilities

import (
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less50 Less-50: ORDER BY Integer
type Less50 struct {
	BaseLesson
}

func (l *Less50) ID() string {
	return "less-50"
}

func (l *Less50) Name() string {
	return "Less-50: ORDER BY Integer"
}

func (l *Less50) Description() string {
	return "ORDER BY整数类型注入"
}

func (l *Less50) Category() string {
	return "Advanced"
}

func (l *Less50) Route(r *gin.RouterGroup) {
	group := r.Group("/less-50")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less50) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-50")

	sort := c.Query("sort")
	if sort == "" {
		sort = "id"
	}

	query := fmt.Sprintf("SELECT * FROM users ORDER BY %s", sort)

	rows, err := l.db.Query(query)
	if err != nil {
		l.renderHTML(c, "less-50", "", err.Error())
		return
	}
	defer rows.Close()

	result := "<table border='1'><tr><th>ID</th><th>Username</th><th>Password</th></tr>"
	for rows.Next() {
		var user models.User
		if err := rows.Scan(&user.ID, &user.Username, &user.Password); err == nil {
			result += fmt.Sprintf("<tr><td>%d</td><td>%s</td><td>%s</td></tr>", user.ID, user.Username, user.Password)
		}
	}
	result += "</table>"

	l.renderHTML(c, "less-50", result, "")
}

func (r *Registry) registerLess50() {
	r.Register(&Less50{BaseLesson{db: r.db, logger: r.logger}})
}

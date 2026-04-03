package vulnerabilities

import (
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less51 Less-51: ORDER BY String
type Less51 struct {
	BaseLesson
}

func (l *Less51) ID() string {
	return "less-51"
}

func (l *Less51) Name() string {
	return "Less-51: ORDER BY String"
}

func (l *Less51) Description() string {
	return "ORDER BY字符串类型注入"
}

func (l *Less51) Category() string {
	return "Advanced"
}

func (l *Less51) Route(r *gin.RouterGroup) {
	group := r.Group("/less-51")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less51) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-51")

	sort := c.Query("sort")
	if sort == "" {
		sort = "id"
	}

	query := fmt.Sprintf("SELECT * FROM users ORDER BY %s", sort)

	rows, err := l.db.Query(query)
	if err != nil {
		l.renderHTML(c, "less-51", "", err.Error())
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

	l.renderHTML(c, "less-51", result, "")
}

func (r *Registry) registerLess51() {
	r.Register(&Less51{BaseLesson{db: r.db, logger: r.logger}})
}

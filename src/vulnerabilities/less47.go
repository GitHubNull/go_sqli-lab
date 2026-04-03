package vulnerabilities

import (
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less47 Less-47: ORDER BY
type Less47 struct {
	BaseLesson
}

func (l *Less47) ID() string {
	return "less-47"
}

func (l *Less47) Name() string {
	return "Less-47: ORDER BY"
}

func (l *Less47) Description() string {
	return "ORDER BY子句注入"
}

func (l *Less47) Category() string {
	return "Advanced"
}

func (l *Less47) Route(r *gin.RouterGroup) {
	group := r.Group("/less-47")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less47) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-47")

	sort := c.Query("sort")
	if sort == "" {
		sort = "id"
	}

	query := fmt.Sprintf("SELECT * FROM users ORDER BY %s", sort)

	rows, err := l.db.Query(query)
	if err != nil {
		l.renderHTML(c, "less-47", "", err.Error())
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

	l.renderHTML(c, "less-47", result, "")
}

func (r *Registry) registerLess47() {
	r.Register(&Less47{BaseLesson{db: r.db, logger: r.logger}})
}

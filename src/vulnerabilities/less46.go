package vulnerabilities

import (
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less46 GET - Error Based - ORDER BY CLAUSE
type Less46 struct {
	BaseLesson
}

func (l *Less46) ID() string {
	return "less-46"
}

func (l *Less46) Name() string {
	return "Less-46: GET - Error Based - ORDER BY CLAUSE"
}

func (l *Less46) Description() string {
	return "GET请求 - ORDER BY子句注入"
}

func (l *Less46) Category() string {
	return "Advanced"
}

func (l *Less46) Route(r *gin.RouterGroup) {
	group := r.Group("/less-46")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less46) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-46")

	sort := c.Query("sort")
	if sort == "" {
		sort = "id"
	}

	query := fmt.Sprintf("SELECT * FROM users ORDER BY %s", sort)

	rows, err := l.db.Query(query)
	if err != nil {
		l.renderHTML(c, "less-46", "", err.Error())
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

	l.renderHTML(c, "less-46", result, "")
}

func (r *Registry) registerLess46() {
	r.Register(&Less46{BaseLesson{db: r.db, logger: r.logger}})
}

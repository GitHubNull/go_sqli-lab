package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less53 Less-53: LIMIT Clause
type Less53 struct {
	BaseLesson
}

func (l *Less53) ID() string {
	return "less-53"
}

func (l *Less53) Name() string {
	return "Less-53: LIMIT Clause"
}

func (l *Less53) Description() string {
	return "LIMIT子句注入"
}

func (l *Less53) Category() string {
	return "Advanced"
}

func (l *Less53) Route(r *gin.RouterGroup) {
	group := r.Group("/less-53")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less53) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-53")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	query := fmt.Sprintf("SELECT * FROM users LIMIT 1 OFFSET %s", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-53", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-53", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-53", result, "")
}

func (r *Registry) registerLess53() {
	r.Register(&Less53{BaseLesson{db: r.db, logger: r.logger}})
}

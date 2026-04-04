package vulnerabilities

import (
	"database/sql"
	"fmt"
	"strings"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less61 Less-61: Challenge 8 - Integer Only
type Less61 struct {
	BaseLesson
}

func (l *Less61) ID() string {
	return "less-61"
}

func (l *Less61) Name() string {
	return "Less-61: Challenge 8 - Integer Only"
}

func (l *Less61) Description() string {
	return "挑战关卡8：纯整数类型注入"
}

func (l *Less61) Category() string {
	return "Challenge"
}

func (l *Less61) Route(r *gin.RouterGroup) {
	group := r.Group("/less-61")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less61) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-61")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	// Remove quotes if present
	id = strings.ReplaceAll(id, "'", "")
	id = strings.ReplaceAll(id, "\"", "")

	query := fmt.Sprintf("SELECT * FROM challenge8 WHERE id=%s LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-61", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-61", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-61", result, "")
}

func (r *Registry) registerLess61() {
	r.Register(&Less61{BaseLesson{db: r.db, logger: r.logger}})
}

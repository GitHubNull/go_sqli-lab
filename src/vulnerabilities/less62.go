package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less62 Less-62: Challenge 9 - Parenthesis Filter
type Less62 struct {
	BaseLesson
}

func (l *Less62) ID() string {
	return "less-62"
}

func (l *Less62) Name() string {
	return "Less-62: Challenge 9 - Parenthesis Filter"
}

func (l *Less62) Description() string {
	return "挑战关卡9：复杂括号过滤注入"
}

func (l *Less62) Category() string {
	return "Challenge"
}

func (l *Less62) Route(r *gin.RouterGroup) {
	group := r.Group("/less-62")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less62) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-62")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	// Multiple parentheses challenge
	query := fmt.Sprintf("SELECT * FROM challenge9 WHERE id=(('%s')) LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-62", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-62", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-62", result, "")
}

func (r *Registry) registerLess62() {
	r.Register(&Less62{BaseLesson{db: r.db, logger: r.logger}})
}

package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less54 Less-54: Challenge 1 - Custom Table Union
type Less54 struct {
	BaseLesson
}

func (l *Less54) ID() string {
	return "less-54"
}

func (l *Less54) Name() string {
	return "Less-54: Challenge 1 - Custom Table"
}

func (l *Less54) Description() string {
	return "挑战关卡1：自定义表Union注入"
}

func (l *Less54) Category() string {
	return "Challenge"
}

func (l *Less54) Route(r *gin.RouterGroup) {
	group := r.Group("/less-54")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less54) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-54")

	id := c.Query("id")
	if id == "" {
		id = "1"
	}

	query := fmt.Sprintf("SELECT * FROM challenge1 WHERE id='%s' LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-54", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-54", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-54", result, "")
}

func (r *Registry) registerLess54() {
	r.Register(&Less54{BaseLesson{db: r.db, logger: r.logger}})
}

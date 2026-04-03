package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less30 GET - Error Based - WAF 2
type Less30 struct {
	BaseLesson
}

func (l *Less30) ID() string {
	return "less-30"
}

func (l *Less30) Name() string {
	return "Less-30: GET - Error Based - WAF 2"
}

func (l *Less30) Description() string {
	return "GET请求 - WAF绕过 2"
}

func (l *Less30) Category() string {
	return "WAF Bypass"
}

func (l *Less30) Route(r *gin.RouterGroup) {
	group := r.Group("/less-30")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less30) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-30")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-30", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-30", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-30", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-30", result, "")
}

func (r *Registry) registerLess30() {
	r.Register(&Less30{BaseLesson{db: r.db, logger: r.logger}})
}

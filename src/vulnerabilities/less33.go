package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less33 GET - Bypass addslashes() - UTF-8
type Less33 struct {
	BaseLesson
}

func (l *Less33) ID() string {
	return "less-33"
}

func (l *Less33) Name() string {
	return "Less-33: GET - Bypass addslashes() - UTF-8"
}

func (l *Less33) Description() string {
	return "GET请求 - UTF-8绕过addslashes()"
}

func (l *Less33) Category() string {
	return "WAF Bypass"
}

func (l *Less33) Route(r *gin.RouterGroup) {
	group := r.Group("/less-33")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less33) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-33")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-33", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-33", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-33", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-33", result, "")
}

func (r *Registry) registerLess33() {
	r.Register(&Less33{BaseLesson{db: r.db, logger: r.logger}})
}

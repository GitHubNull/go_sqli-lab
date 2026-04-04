package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less36 GET - Bypass MySQL Real Escape String
type Less36 struct {
	BaseLesson
}

func (l *Less36) ID() string {
	return "less-36"
}

func (l *Less36) Name() string {
	return "Less-36: GET - Bypass MySQL Real Escape String"
}

func (l *Less36) Description() string {
	return "GET请求 - 绕过mysql_real_escape_string"
}

func (l *Less36) Category() string {
	return "WAF Bypass"
}

func (l *Less36) Route(r *gin.RouterGroup) {
	group := r.Group("/less-36")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less36) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-36")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-36", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-36", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-36", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-36", result, "")
}

func (r *Registry) registerLess36() {
	r.Register(&Less36{BaseLesson{db: r.db, logger: r.logger}})
}

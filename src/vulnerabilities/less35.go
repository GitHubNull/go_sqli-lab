package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less35 GET - Integer - Bypass addslashes()
type Less35 struct {
	BaseLesson
}

func (l *Less35) ID() string {
	return "less-35"
}

func (l *Less35) Name() string {
	return "Less-35: GET - Integer - Bypass addslashes()"
}

func (l *Less35) Description() string {
	return "GET请求 - 整数类型绕过addslashes()"
}

func (l *Less35) Category() string {
	return "WAF Bypass"
}

func (l *Less35) Route(r *gin.RouterGroup) {
	group := r.Group("/less-35")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less35) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-35")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-35", "请输入ID参数", "")
		return
	}

	// 整数类型，无引号
	query := fmt.Sprintf("SELECT * FROM users WHERE id=%s LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-35", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-35", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-35", result, "")
}

func (r *Registry) registerLess35() {
	r.Register(&Less35{BaseLesson{db: r.db, logger: r.logger}})
}

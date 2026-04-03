package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less29 GET - Error Based - WAF 1
type Less29 struct {
	BaseLesson
}

func (l *Less29) ID() string {
	return "less-29"
}

func (l *Less29) Name() string {
	return "Less-29: GET - Error Based - WAF 1"
}

func (l *Less29) Description() string {
	return "GET请求 - WAF绕过 1"
}

func (l *Less29) Category() string {
	return "WAF Bypass"
}

func (l *Less29) Route(r *gin.RouterGroup) {
	group := r.Group("/less-29")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less29) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-29")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-29", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-29", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-29", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-29", result, "")
}

func (r *Registry) registerLess29() {
	r.Register(&Less29{BaseLesson{db: r.db, logger: r.logger}})
}

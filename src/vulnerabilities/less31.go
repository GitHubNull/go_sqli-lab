package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less31 GET - Error Based - WAF 3
type Less31 struct {
	BaseLesson
}

func (l *Less31) ID() string {
	return "less-31"
}

func (l *Less31) Name() string {
	return "Less-31: GET - Error Based - WAF 3"
}

func (l *Less31) Description() string {
	return "GET请求 - WAF绕过 3"
}

func (l *Less31) Category() string {
	return "WAF Bypass"
}

func (l *Less31) Route(r *gin.RouterGroup) {
	group := r.Group("/less-31")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less31) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-31")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-31", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-31", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-31", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-31", result, "")
}

func (r *Registry) registerLess31() {
	r.Register(&Less31{BaseLesson{db: r.db, logger: r.logger}})
}

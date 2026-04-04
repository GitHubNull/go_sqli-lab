package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less28 GET - Error Based - Impedance Mismatch
type Less28 struct {
	BaseLesson
}

func (l *Less28) ID() string {
	return "less-28"
}

func (l *Less28) Name() string {
	return "Less-28: GET - Error Based - Impedance Mismatch"
}

func (l *Less28) Description() string {
	return "GET请求 - 阻抗不匹配的SQL注入"
}

func (l *Less28) Category() string {
	return "WAF Bypass"
}

func (l *Less28) Route(r *gin.RouterGroup) {
	group := r.Group("/less-28")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less28) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-28")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-28", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-28", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-28", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-28", result, "")
}

func (r *Registry) registerLess28() {
	r.Register(&Less28{BaseLesson{db: r.db, logger: r.logger}})
}

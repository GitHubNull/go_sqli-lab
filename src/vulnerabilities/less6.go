package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less6 Less-6漏洞: Double Query - Double Quotes
type Less6 struct {
	BaseLesson
}

func (l *Less6) ID() string {
	return "less-6"
}

func (l *Less6) Name() string {
	return "Less-6: Double Query - Double Quotes"
}

func (l *Less6) Description() string {
	return "双查询注入 - 双引号字符串类型"
}

func (l *Less6) Category() string {
	return "Double Injection"
}

func (l *Less6) Route(r *gin.RouterGroup) {
	group := r.Group("/less-6")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less6) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-6")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-6", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id=\"%s\" LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-6", "", "You are in...........")
		} else {
			l.renderHTML(c, "less-6", "", err.Error())
		}
		return
	}

	l.renderHTML(c, "less-6", "You are in...........", "")
}

func (r *Registry) registerLess6() {
	r.Register(&Less6{BaseLesson{db: r.db, logger: r.logger}})
}

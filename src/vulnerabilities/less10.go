package vulnerabilities

import (
	"database/sql"
	"fmt"
	"time"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less10 Less-10漏洞: Blind - Time Based - Double Quotes
type Less10 struct {
	BaseLesson
}

func (l *Less10) ID() string {
	return "less-10"
}

func (l *Less10) Name() string {
	return "Less-10: Blind - Time Based - Double Quotes"
}

func (l *Less10) Description() string {
	return "盲注 - 基于时间的SQL注入 (双引号)"
}

func (l *Less10) Category() string {
	return "Time Based"
}

func (l *Less10) Route(r *gin.RouterGroup) {
	group := r.Group("/less-10")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less10) handleIndex(c *gin.Context) {
	start := time.Now()
	l.logRequest(c, "less-10")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-10", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id=\"%s\" LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	elapsed := time.Since(start)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-10", "You are in...........", "")
		} else {
			l.renderHTML(c, "less-10", "You are in...........", "")
		}
		return
	}

	c.Header("X-Response-Time", fmt.Sprintf("%dms", elapsed.Milliseconds()))
	l.renderHTML(c, "less-10", "You are in...........", "")
}

func (r *Registry) registerLess10() {
	r.Register(&Less10{BaseLesson{db: r.db, logger: r.logger}})
}

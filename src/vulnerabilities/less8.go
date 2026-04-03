package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less8 Less-8漏洞: Blind - Boolean Based - Single Quotes
type Less8 struct {
	BaseLesson
}

func (l *Less8) ID() string {
	return "less-8"
}

func (l *Less8) Name() string {
	return "Less-8: Blind - Boolean Based - Single Quotes"
}

func (l *Less8) Description() string {
	return "盲注 - 基于布尔的SQL注入 (单引号)"
}

func (l *Less8) Category() string {
	return "Blind Injection"
}

func (l *Less8) Route(r *gin.RouterGroup) {
	group := r.Group("/less-8")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less8) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-8")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-8", "请输入ID参数", "")
		return
	}

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			// 盲注特点：不显示错误信息
			l.renderHTML(c, "less-8", "", "")
		} else {
			// 盲注：不显示SQL错误
			l.renderHTML(c, "less-8", "", "")
		}
		return
	}

	// 只有成功时才显示信息
	l.renderHTML(c, "less-8", "You are in...........", "")
}

func (r *Registry) registerLess8() {
	r.Register(&Less8{BaseLesson{db: r.db, logger: r.logger}})
}

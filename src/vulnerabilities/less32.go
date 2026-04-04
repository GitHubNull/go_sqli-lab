package vulnerabilities

import (
	"database/sql"
	"fmt"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less32 GET - Bypass addslashes() - Wide Char
type Less32 struct {
	BaseLesson
}

func (l *Less32) ID() string {
	return "less-32"
}

func (l *Less32) Name() string {
	return "Less-32: GET - Bypass addslashes() - Wide Char"
}

func (l *Less32) Description() string {
	return "GET请求 - 绕过addslashes()宽字符注入"
}

func (l *Less32) Category() string {
	return "WAF Bypass"
}

func (l *Less32) Route(r *gin.RouterGroup) {
	group := r.Group("/less-32")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less32) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-32")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-32", "请输入ID参数", "")
		return
	}

	// 模拟addslashes处理
	id = addslashes(id)

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-32", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-32", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-32", result, "")
}

func addslashes(s string) string {
	result := ""
	for _, ch := range s {
		switch ch {
		case '\'', '"', '\\':
			result += "\\" + string(ch)
		default:
			result += string(ch)
		}
	}
	return result
}

func (r *Registry) registerLess32() {
	r.Register(&Less32{BaseLesson{db: r.db, logger: r.logger}})
}

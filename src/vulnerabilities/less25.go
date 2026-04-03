package vulnerabilities

import (
	"database/sql"
	"fmt"
	"strings"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less25 GET - Error Based - OR & AND Bypass
type Less25 struct {
	BaseLesson
}

func (l *Less25) ID() string {
	return "less-25"
}

func (l *Less25) Name() string {
	return "Less-25: GET - Error Based - OR & AND Bypass"
}

func (l *Less25) Description() string {
	return "GET请求 - 过滤OR和AND的SQL注入"
}

func (l *Less25) Category() string {
	return "WAF Bypass"
}

func (l *Less25) Route(r *gin.RouterGroup) {
	group := r.Group("/less-25")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less25) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-25")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-25", "请输入ID参数", "")
		return
	}

	// 过滤OR和AND（不区分大小写）
	id = strings.ReplaceAll(id, "OR", "")
	id = strings.ReplaceAll(id, "or", "")
	id = strings.ReplaceAll(id, "AND", "")
	id = strings.ReplaceAll(id, "and", "")

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 0,1", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-25", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-25", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-25", result, "")
}

func (r *Registry) registerLess25() {
	r.Register(&Less25{BaseLesson{db: r.db, logger: r.logger}})
}

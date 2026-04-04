package vulnerabilities

import (
	"database/sql"
	"fmt"
	"strings"

	"go-sqli-lab/src/models"

	"github.com/gin-gonic/gin"
)

// Less26 GET - Error Based - Spaces and Comments
type Less26 struct {
	BaseLesson
}

func (l *Less26) ID() string {
	return "less-26"
}

func (l *Less26) Name() string {
	return "Less-26: GET - Error Based - Spaces and Comments"
}

func (l *Less26) Description() string {
	return "GET请求 - 过滤空格和注释的SQL注入"
}

func (l *Less26) Category() string {
	return "WAF Bypass"
}

func (l *Less26) Route(r *gin.RouterGroup) {
	group := r.Group("/less-26")
	{
		group.GET("", l.handleIndex)
		group.GET("/", l.handleIndex)
	}
}

func (l *Less26) handleIndex(c *gin.Context) {
	l.logRequest(c, "less-26")

	id := c.Query("id")
	if id == "" {
		l.renderHTML(c, "less-26", "请输入ID参数", "")
		return
	}

	// 过滤空格和注释
	id = strings.ReplaceAll(id, " ", "")
	id = strings.ReplaceAll(id, "--", "")
	id = strings.ReplaceAll(id, "#", "")
	id = strings.ReplaceAll(id, "/*", "")
	id = strings.ReplaceAll(id, "*/", "")

	query := fmt.Sprintf("SELECT * FROM users WHERE id='%s' LIMIT 1 OFFSET 0", id)

	var user models.User
	err := l.db.QueryRow(query).Scan(&user.ID, &user.Username, &user.Password)

	if err != nil {
		if err == sql.ErrNoRows {
			l.renderHTML(c, "less-26", "", "没有找到记录")
		} else {
			l.renderHTML(c, "less-26", "", err.Error())
		}
		return
	}

	result := fmt.Sprintf("Your Login name: %s<br>Your Password: %s", user.Username, user.Password)
	l.renderHTML(c, "less-26", result, "")
}

func (r *Registry) registerLess26() {
	r.Register(&Less26{BaseLesson{db: r.db, logger: r.logger}})
}
